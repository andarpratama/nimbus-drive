package models
import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)
type SharedFile struct {
	ID              uuid.UUID `gorm:"type:char(36);primaryKey;default:(UUID())"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
	FileID          uuid.UUID `gorm:"type:char(36);not null"`
	SharedWithEmail string    `gorm:"type:varchar(255)"`
	Token           string    `gorm:"type:varchar(255);uniqueIndex"`
	ExpiredAt       *int64
}
