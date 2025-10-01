package api

import "time"

type shortenUrlRequest struct {
	OriginalURL string  `json:"originalUrl" validate:"required,http_url"`
	Slug        *string `json:"slug"`
	// needs to be RFC3339 with timezone (or UTC)
	ExpiresAt *time.Time `json:"expiresAt"`
}

type shortenUrlResponse struct {
	ShortURL string `json:"shortUrl"`
}

// --------ADMIN----------

// same as for /shorten request, but separated
type createShortLinkRequest struct {
	OriginalURL string  `json:"originalUrl" validate:"required,http_url"`
	Slug        *string `json:"slug"`
	// needs to be RFC3339 with timezone (or UTC)
	ExpiresAt *time.Time `json:"expiresAt"`
}

type updateShortLinkRequest struct {
	OriginalURL *string `json:"originalUrl" validate:"omitnil,http_url"`
	Slug        *string `json:"slug"`
	// needs to be RFC3339 with timezone (or UTC)
	ExpiresAt *time.Time `json:"expiresAt"`
}
type shortLinkDto struct {
	ID          uint       `json:"id"`
	OriginalURL string     `json:"originalUrl"`
	ShortURL    string     `json:"shortUrl"`
	Slug        *string    `json:"slug"`
	ExpiresAt   *time.Time `json:"expiresAt"`
}

type genericErrorResponse struct {
	Error string `json:"error"`
}
