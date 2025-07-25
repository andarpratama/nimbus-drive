package models

import "gorm.io/gorm"

type Session struct {
	gorm.Model
	UserID uint
	Token  string `gorm:"type:varchar(255);uniqueIndex"`
}
