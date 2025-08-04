package models
import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)
type Session struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey;default:(UUID())"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	UserID uuid.UUID `gorm:"type:char(36);not null"`
	Token  string    `gorm:"type:varchar(255);uniqueIndex"`
}
