package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/lin-snow/ech0/config"
	"github.com/lin-snow/ech0/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() error {
	dbType := config.Config.Database.Type
	dbPath := config.Config.Database.Path

	dir := fmt.Sprintf("%s", dbPath[:len(dbPath)-len("/ech0.db")])
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatalf("Failed to create database directory: %v", err)
		return err
	}

	// 配置 GORM 日志
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,  // 修改为 Info 级别
			IgnoreRecordNotFoundError: false,
			Colorful:                  true,
		},
	)

	var err error
	if dbType == "sqlite" {
		DB, err = gorm.Open(sqlite.Open(dbPath+"?_timeout=5000"), &gorm.Config{
			Logger:                 newLogger,
			PrepareStmt:           true,
			QueryFields:           true,
			AllowGlobalUpdate:     true,
			CreateBatchSize:       1000,
			SkipDefaultTransaction: false,
			DryRun:                false,
			NowFunc: func() time.Time {
				return time.Now().Local()
			},
		})
	} else {
		log.Fatalf(models.DatabaseTypeMessage+": %s", dbType)
		return err
	}

	if err != nil {
		log.Fatalf(models.DatabaseConnectionError+": %v", err)
		return err
	}

	// 设置连接池
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	// 设置数据库连接池参数
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 设置数据库文件权限
	if err := os.Chmod(dbPath, 0666); err != nil {
		log.Printf("Warning: Failed to set database file permissions: %v", err)
	}

	// 确保数据库目录权限正确
	if err := os.Chmod(dir, 0755); err != nil {
		log.Printf("Warning: Failed to set database directory permissions: %v", err)
	}

	if err = models.MigrateDB(DB); err != nil {
		log.Fatal(models.DatabaseMigrationError+":", err)
	}

	return nil
}