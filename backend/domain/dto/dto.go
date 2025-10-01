package dto

import "time"

type RequestInfoDTO struct {
	RequestId string
	Timestamp time.Time
	RealIP    *string
	UserAgent *string
	Headers   map[string][]string
	Path      string
	Method    string
	Query     map[string][]string
	Body      map[string]any
}

type ShortLinkUpdateDTO struct {
	OriginalURL *string
	Slug        *string
	ExpiresAt   *time.Time
}
