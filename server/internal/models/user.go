package models
import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)
type User struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey;default:(UUID())"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name     string `gorm:"type:varchar(255);not null"`
	Username string `gorm:"type:varchar(255);uniqueIndex;not null"`
	Email    string `gorm:"type:varchar(255);uniqueIndex;not null"`
	Password string `gorm:"type:varchar(255);not null"`
	Role     string `gorm:"type:ENUM('superadmin','admin','user');default:'user'"`
}
