package model

import (
	"time"

	"gorm.io/datatypes"
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

type RequestInfo struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	RequestId string    `gorm:"not null"`
	Timestamp time.Time `gorm:"not null"`
	RealIP    *string
	UserAgent *string
	Headers   datatypes.JSON `gorm:"type:jsonb"`
	Path      string         `gorm:"not null"`
	Method    string         `gorm:"not null"`
	Query     datatypes.JSON `gorm:"type:jsonb"`
	Body      datatypes.JSON `gorm:"type:jsonb"`
}
