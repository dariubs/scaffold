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
	GitHubID    string `gorm:"uniqueIndex"`        // GitHub OAuth ID
	LinkedInID  string `gorm:"uniqueIndex"`        // LinkedIn OAuth ID
	XID         string `gorm:"uniqueIndex"`        // X (Twitter) OAuth ID
	LoginMethod string `gorm:"default:'password'"` // 'password', 'google', 'github', 'linkedin', 'x'
	IsAdmin     bool   `gorm:"default:false"`      // Admin flag
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
