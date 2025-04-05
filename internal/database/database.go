package database

import (
   "errors"
    "fmt"
    "os"

    "github.com/lin-snow/ech0/config"
    "github.com/lin-snow/ech0/internal/models"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() error {
    if DB != nil {
        return nil
    }

    dbType := config.Config.Database.Type
    dbPath := config.Config.Database.Path

    dir := fmt.Sprintf("%s", dbPath[:len(dbPath)-len("/noise.db")]) 
    if err := os.MkdirAll(dir, os.ModePerm); err != nil {
        return fmt.Errorf("创建数据库目录失败: %v", err)
    }

    var err error
    if dbType == "sqlite" {
        DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
    } else {
        return fmt.Errorf("不支持的数据库类型: %s", dbType)
    }

    if err != nil {
        return fmt.Errorf("数据库连接失败: %v", err)
    }

    models.SetDB(DB)

    if err = models.MigrateDB(DB); err != nil {
        return fmt.Errorf("数据库迁移失败: %v", err)
    }

    return nil
}

func GetSetting() (*models.Setting, error) {
    if DB == nil {
        return nil, errors.New("数据库未初始化")
    }

    var setting models.Setting
    if err := DB.First(&setting).Error; err != nil {
        return nil, err
    }
    return &setting, nil
}

func UpdateSetting(allowRegistration bool) error {
    if DB == nil {
        return errors.New("数据库未初始化")
    }

    return DB.Model(&models.Setting{}).Where("1 = 1").Update("allow_registration", allowRegistration).Error
}

func GetDB() (*gorm.DB, error) {
    if DB == nil {
        return nil, errors.New("数据库未初始化")
    }
    return DB, nil
}

func GetSystemStatus() (map[string]interface{}, error) {
    if DB == nil {
        return nil, errors.New("数据库未初始化")
    }

    var setting models.Setting
    if err := DB.First(&setting).Error; err != nil {
        return nil, err
    }

    var adminUser models.User
    if err := DB.Where("is_admin = ?", true).First(&adminUser).Error; err != nil {
        return nil, err
    }

    var totalMessages int64
    if err := DB.Model(&models.Message{}).Count(&totalMessages).Error; err != nil {
        return nil, err
    }

    return map[string]interface{}{
        "username":          adminUser.Username,
        "isAdmin":          adminUser.IsAdmin,
        "total_messages":    totalMessages,
        "allowRegistration": setting.AllowRegistration,
    }, nil
}

func ReconnectDB() error {
    if DB != nil {
        sqlDB, err := DB.DB()
        if err == nil {
            if err := sqlDB.Close(); err != nil {
                return fmt.Errorf("关闭数据库连接失败: %v", err)
            }
        }
        DB = nil
    }
    return InitDB()
}