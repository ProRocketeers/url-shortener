package model

import (
	"time"

	"gorm.io/gorm"
)

type ShortLink struct {
	gorm.Model

	OriginalURL string `gorm:"not null"`
	Slug        string `gorm:"not null;unique"`
	ExpiresAt   *time.Time
}
