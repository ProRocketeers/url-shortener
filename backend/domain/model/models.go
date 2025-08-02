package model

import (
	"time"
)

type ShortLink struct {
	// original `gorm.Model`, but without `DeletedAt` which automatically soft deletes
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	OriginalURL string `gorm:"not null"`
	Slug        string `gorm:"not null;unique"`
	ExpiresAt   *time.Time
}
