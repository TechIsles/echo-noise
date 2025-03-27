package models

import "gorm.io/gorm"

// Setting 系统设置模型
type Setting struct {
    gorm.Model
    AllowRegistration bool `gorm:"default:true"`
}