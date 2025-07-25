package database

import (
	"github.com/andarpratama/nimbus-drive/internal/models"

	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.File{},
		&models.Folder{},
		&models.SharedFile{},
		&models.Session{},
	)
}
