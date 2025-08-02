package api

import "time"

// TODO: Swagger annotations
type shortenUrlRequest struct {
	OriginalURL string  `json:"originalUrl" validate:"required,http_url"`
	Slug        *string `json:"slug"`
	// needs to be RFC3339 with timezone (or UTC)
	ExpiresAt *time.Time `json:"expiresAt"`
}

type shortenUrlResponse struct {
	ShortURL string `json:"shortUrl"`
}
