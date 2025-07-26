package models

import "gorm.io/gorm"

type SharedFile struct {
	gorm.Model
	FileID          uint
	SharedWithEmail string `gorm:"type:varchar(255)"`
	Token           string `gorm:"type:varchar(255);uniqueIndex"`
	ExpiredAt       *int64
}
