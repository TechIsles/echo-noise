package database

import (
    "fmt"
    "log"
    "os"

    "github.com/lin-snow/ech0/config"
    "github.com/lin-snow/ech0/internal/models"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

var DB *gorm.DB

// InitDB 初始化数据库
func InitDB() error {
    // 读取数据库类型和路径
    dbType := config.Config.Database.Type
    dbPath := config.Config.Database.Path

    // 确保数据库目录存在
    dir := fmt.Sprintf("%s", dbPath[:len(dbPath)-len("/noise.db")]) 
    if err := os.MkdirAll(dir, os.ModePerm); err != nil {
        log.Fatalf("Failed to create database directory: %v", err)
        return err
    }

    var err error
    // 根据数据库类型选择不同的数据库驱动
    if dbType == "sqlite" {
        DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
    } else {
        log.Fatalf(models.DatabaseTypeMessage+": %s", dbType)
        return err
    }

    if err != nil {
        log.Fatalf(models.DatabaseConnectionError+": %v", err)
        return err
    }

    // 设置全局数据库连接
    models.SetDB(DB)

    // 执行数据库迁移和初始化
    if err = models.MigrateDB(DB); err != nil {
        log.Fatal(models.DatabaseMigrationError+":", err)
    }

    return nil
}
// GetSetting 获取系统设置
func GetSetting() (*models.Setting, error) {
    var setting models.Setting
    if err := DB.First(&setting).Error; err != nil {
        return nil, err
    }
    return &setting, nil
}

// UpdateSetting 更新系统设置
func UpdateSetting(allowRegistration bool) error {
    return DB.Model(&models.Setting{}).Where("1 = 1").Update("allow_registration", allowRegistration).Error
}

// GetDB 获取数据库连接
func GetDB() *gorm.DB {
    return DB
}

// GetSystemStatus 获取系统完整状态
func GetSystemStatus() (map[string]interface{}, error) {
    var setting models.Setting
    if err := DB.First(&setting).Error; err != nil {
        return nil, err
    }

    // 获取管理员信息
    var adminUser models.User
    if err := DB.Where("is_admin = ?", true).First(&adminUser).Error; err != nil {
        return nil, err
    }

    // 获取消息总数
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
// ReconnectDB 重新连接数据库
func ReconnectDB() error {
    if DB != nil {
        sqlDB, err := DB.DB()
        if err == nil {
            sqlDB.Close() // 关闭现有连接
        }
    }
    return InitDB() // 重新初始化
}