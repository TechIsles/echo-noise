package models

import (
    "gorm.io/gorm"
)

func MigrateDB(db *gorm.DB) error {
    // 自动迁移所有模型
    err := db.AutoMigrate(&User{}, &Message{}, &Setting{})
    if err != nil {
        return err
    }

    // 初始化系统设置
    var setting Setting
    result := db.First(&setting)
    if result.RowsAffected == 0 {
        defaultSetting := Setting{
            AllowRegistration: true,
        }
        if err := db.Create(&defaultSetting).Error; err != nil {
            return err
        }
    }

    // 初始化管理员账户
    var adminUser User
    result = db.Where("is_admin = ?", true).First(&adminUser)
    if result.RowsAffected == 0 {
        defaultAdmin := User{
            Username: "admin",
            Password: HashPassword("admin"),
            IsAdmin:  true,
        }
        if err := db.Create(&defaultAdmin).Error; err != nil {
            return err
        }
    }

    return nil
}