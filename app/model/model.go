package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username    string `gorm:"uniqueIndex;not null"`
	Email       string `gorm:"uniqueIndex;not null"`
	Password    string // Can be empty for OAuth users
	Name        string
	AvatarURL   string
	Bio         string
	GoogleID    string `gorm:"uniqueIndex"`        // Google OAuth ID
	LoginMethod string `gorm:"default:'password'"` // 'password' or 'google'
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
