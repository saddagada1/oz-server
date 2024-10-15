package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email            string `gorm:"unique"`
	Username         string `gorm:"unique"`
	Password         string
	HasLinkedEbay    bool
	EbayAccessToken  *string
	EbayRefreshToken *string
	TokenVersion     uint
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
