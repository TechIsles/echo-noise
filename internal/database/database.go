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

    // 自动迁移新增的表
    if err = DB.AutoMigrate(
        &models.Setting{},
        &models.User{},
        &models.Message{},
        &models.SiteConfig{},
    ); err != nil {
        log.Fatal("数据库迁移失败:", err)
        return err
    }

    // 初始化前端配置
    var config models.SiteConfig
    if err := DB.First(&config).Error; err != nil {
        defaultBg := `https://s2.loli.net/2025/03/27/KJ1trnU2ksbFEYM.jpg,
https://s2.loli.net/2025/03/27/MZqaLczCvwjSmW7.jpg,
https://s2.loli.net/2025/03/27/UMijKXwJ9yTqSeE.jpg,
https://s2.loli.net/2025/03/27/WJQIlkXvBg2afcR.jpg,
https://s2.loli.net/2025/03/27/oHNQtf4spkq2iln.jpg,
https://s2.loli.net/2025/03/27/PMRuX5loc6Uaimw.jpg,
https://s2.loli.net/2025/03/27/U2WIslbNyTLt4rD.jpg,
https://s2.loli.net/2025/03/27/xu1jZL5Og4pqT9d.jpg,
https://s2.loli.net/2025/03/27/OXqwzZ6v3PVIns9.jpg,
https://s2.loli.net/2025/03/27/HGuqlE6apgNywbh.jpg,
https://s2.loli.net/2025/03/26/d7iyuPYA8cRqD1K.jpg,
https://s2.loli.net/2025/03/27/7Zck3y6XTzhYPs5.jpg,
https://s2.loli.net/2025/03/27/y67m2k5xcSdTsHN.jpg`
        
        newConfig := &models.SiteConfig{
            SiteTitle:    "Noise的说说笔记",
            SubtitleText: "欢迎访问，点击头像可更换封面背景！",
            AvatarURL:    "https://s2.loli.net/2025/03/24/HnSXKvibAQlosIW.png",
            Username:     "Noise",
            Description:  "执迷不悟",
            Backgrounds:  defaultBg,
        }
        
        if err := DB.Create(newConfig).Error; err != nil {
            log.Printf("创建默认前端配置失败: %v", err)
            return err
        }
    }

    // 初始化系统设置
    var setting models.Setting
    result := DB.First(&setting)
    if result.RowsAffected == 0 {
        defaultSetting := models.Setting{
            AllowRegistration: true,
        }
        if err := DB.Create(&defaultSetting).Error; err != nil {
            log.Fatal("Failed to create default setting:", err)
            return err
        }
    }

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