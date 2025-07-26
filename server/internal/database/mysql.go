package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectMySQL() error {
	// Use environment variable APP_DEBUG to determine log level
	var logLevel logger.LogLevel
	if os.Getenv("APP_DEBUG") == "true" {
		logLevel = logger.Info // or logger.Warn for fewer logs
	} else {
		logLevel = logger.Silent
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  logLevel,
				IgnoreRecordNotFoundError: true,
			},
		),
	})
	if err != nil {
		return err
	}

	DB = db
	return nil
}
