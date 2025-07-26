package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Name     string `gorm:"type:varchar(255);not null"`
	Email    string `gorm:"type:varchar(255);uniqueIndex;not null"`
	Password string `gorm:"type:varchar(255);not null"`
	Role     string `gorm:"type:ENUM('superadmin','admin','user');default:'user'"`
}
