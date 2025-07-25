package models

import "gorm.io/gorm"

type Folder struct {
	gorm.Model
	Name     string  `gorm:"not null"`
	ParentID *uint   // Self-referencing key for nested folders
	Parent   *Folder `gorm:"foreignKey:ParentID"` // Optional: explicitly define relationship
	UserID   uint
	Files    []File `gorm:"foreignKey:FolderID"` // Links files to this folder
}
