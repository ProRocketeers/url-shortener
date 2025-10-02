package api

import (
	"time"

	"github.com/ProRocketeers/url-shortener/domain/dto"
)

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

type requestInfoDto struct {
	ID        uint                `json:"id"`
	RequestId string              `json:"requestId"`
	Timestamp time.Time           `json:"timestamp"`
	RealIP    *string             `json:"realIp"`
	UserAgent *string             `json:"userAgent"`
	Headers   map[string][]string `json:"headers"`
	Path      string              `json:"path"`
	Method    string              `json:"method"`
	Query     map[string][]string `json:"query"`
	Body      map[string]any      `json:"body"`
}

type listRequestInfoResponse struct {
	Data       []requestInfoDto       `json:"data"`
	Pagination *dto.PaginationInfoDTO `json:"pagination"`
}

type genericErrorResponse struct {
	Error string `json:"error"`
}
