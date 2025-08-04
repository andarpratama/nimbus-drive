package models
import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)
type Folder struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey;default:(UUID())"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name     string     `gorm:"not null"`
	ParentID *uuid.UUID `gorm:"type:char(36)"` 
	Parent   *Folder    `gorm:"foreignKey:ParentID"` 
	UserID   uuid.UUID  `gorm:"type:char(36);not null"`
	Files    []File     `gorm:"foreignKey:FolderID"` 
}
