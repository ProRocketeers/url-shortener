package v1

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
	Slug     string `json:"slug"`
}

type shortLinkInfoResponse struct {
	OriginalURL string     `json:"originalUrl"`
	ClickCount  int64      `json:"clickCount"`
	ExpiresAt   *time.Time `json:"expiresAt" example:"2036-03-09T12:00:00Z"`
}

// --------ADMIN----------

// same as for /shorten request, but separated
type createShortLinkRequest struct {
	OriginalURL string  `json:"originalUrl" validate:"required,http_url" example:"https://example.com/very/long/url/that/needs/shortening"`
	Slug        *string `json:"slug" example:"myl1nk"`
	// needs to be RFC3339 with timezone (or UTC)
	ExpiresAt *time.Time `json:"expiresAt" example:"2036-03-09T12:00:00Z"`
}

type updateShortLinkRequest struct {
	OriginalURL *string `json:"originalUrl" validate:"omitnil,http_url"`
	Slug        *string `json:"slug" example:"myl1nk"`
	// needs to be RFC3339 with timezone (or UTC)
	ExpiresAt *time.Time `json:"expiresAt" example:"2036-03-09T12:00:00Z"`
}
type shortLinkDto struct {
	ID          uint       `json:"id" example:"1"`
	OriginalURL string     `json:"originalUrl" example:"https://example.com/very/long/url/that/needs/shortening"`
	ShortURL    string     `json:"shortUrl" example:"https://short.link/v1/myl1nk"`
	Slug        *string    `json:"slug" example:"myl1nk"`
	ExpiresAt   *time.Time `json:"expiresAt" example:"2036-03-09T12:00:00Z"`
}

type listShortLinksResponse struct {
	Data       []shortLinkDto         `json:"data"`
	Pagination *dto.PaginationInfoDTO `json:"pagination"`
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
