package models

import "gorm.io/gorm"

type File struct {
	gorm.Model        // includes ID, CreatedAt, UpdatedAt, DeletedAt
	Name       string `gorm:"not null"`
	Size       int64
	Path       string `gorm:"not null"`
	MimeType   string
	UserID     uint
	FolderID   *uint
}
