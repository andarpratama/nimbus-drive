package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type File struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey;default:(UUID())"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Name     string     `gorm:"not null"`
	Size     int64      `gorm:"default:0"`
	Path     string     `gorm:"not null"`
	MimeType string     `gorm:"type:varchar(255)"`
	UserID   uuid.UUID  `gorm:"type:char(36);not null"`
	FolderID *uuid.UUID `gorm:"type:char(36)"` // Optional: file can be in a folder
	Folder   *Folder    `gorm:"foreignKey:FolderID"` // Relationship to folder
}
