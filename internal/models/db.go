package models

import "gorm.io/gorm"

var db *gorm.DB

// SetDB 设置数据库连接
func SetDB(database *gorm.DB) {
    db = database
}

// GetDB 获取数据库连接
func GetDB() *gorm.DB {
    return db
}