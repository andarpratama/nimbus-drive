package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Starred struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey;default:(UUID())"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	UserID   uuid.UUID `gorm:"type:char(36);not null"`
	FileID   *uuid.UUID `gorm:"type:char(36)"` // Optional: can star a file
	FolderID *uuid.UUID `gorm:"type:char(36)"` // Optional: can star a folder
	
	// Relationships
	User   User   `gorm:"foreignKey:UserID"`
	File   *File  `gorm:"foreignKey:FileID"`
	Folder *Folder `gorm:"foreignKey:FolderID"`
} 