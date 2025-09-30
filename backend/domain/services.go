package domain

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/ProRocketeers/url-shortener/domain/dto"
	"github.com/ProRocketeers/url-shortener/domain/model"
	"github.com/ProRocketeers/url-shortener/storage"
	"github.com/rs/zerolog/log"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"github.com/segmentio/ksuid"
)

type ShortLinkService struct {
	Repository *storage.ShortLinkRepository
	BaseUrl    string
}

func (s *ShortLinkService) Create(ctx context.Context, originalUrl string, slug *string, expiresAt *time.Time) (model.ShortLink, error) {
	usedSlug := createSlug()
	if slug != nil {
		// first check if it exists or not
		_, err := s.Repository.FindBySlug(ctx, *slug)
		if err == nil || !errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ShortLink{}, &ShortLinkError{ErrorCodeSlugConflict}
		}
		usedSlug = *slug
	}

	link := model.ShortLink{
		// URL should be already validated via validator struct tag
		OriginalURL: originalUrl,
		Slug:        usedSlug,
		ExpiresAt:   expiresAt,
	}

	log.Info().Str("url", originalUrl).Str("slug", usedSlug).Msg("creating new link")

	err := s.Repository.Create(ctx, &link)
	if err != nil {
		log.Error().
			Err(err).
			Str("originalUrl", originalUrl).
			Str("slug", usedSlug).
			Msg("creating link error")
		return model.ShortLink{}, &ShortLinkError{ErrorCodeLinkCreate}
	}

	return link, nil
}

func (s *ShortLinkService) FindBySlug(ctx context.Context, slug string) (model.ShortLink, error) {
	link, err := s.Repository.FindBySlug(ctx, slug)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return model.ShortLink{}, &ShortLinkError{ErrorCodeLinkNotFound}
	} else if err != nil {
		log.Error().
			Err(err).
			Str("slug", slug).
			Msg("getting link other error")
		return model.ShortLink{}, &ShortLinkError{ErrorCodeLinkGetOther}
	}
	if link.ExpiresAt != nil && time.Now().UTC().After(*link.ExpiresAt) {
		return model.ShortLink{}, &ShortLinkError{ErrorCodeLinkExpired}
	}
	return link, nil
}

func (s *ShortLinkService) GetShortUrl(link model.ShortLink) string {
	return fmt.Sprintf("%s/%s", s.BaseUrl, link.Slug)
}

func createSlug() string {
	return ksuid.New().String()[:8]
}

// ---------------------

type RequestInfoService struct {
	Repository *storage.RequestInfoRepository
}

func (s *RequestInfoService) Create(ctx context.Context, requestInfoDto dto.RequestInfoDTO) {
	headers, headersErr := jsonFromMap(requestInfoDto.Headers)
	query, queryErr := jsonFromMap(requestInfoDto.Query)
	body, bodyErr := jsonFromMap(requestInfoDto.Body)
	if headersErr != nil || queryErr != nil || bodyErr != nil {
		log.Warn().
			Errs("errors", []error{headersErr, queryErr, bodyErr}).
			Str("path", requestInfoDto.Path).
			Msg("errors serializing request info to json")
	}

	info := model.RequestInfo{
		RequestId: requestInfoDto.RequestId,
		Timestamp: requestInfoDto.Timestamp,
		RealIP:    requestInfoDto.RealIP,
		UserAgent: requestInfoDto.UserAgent,
		Path:      requestInfoDto.Path,
		Method:    requestInfoDto.Method,
		Headers:   headers,
		Query:     query,
		Body:      body,
	}

	err := s.Repository.Create(ctx, &info)
	if err != nil {
		log.Warn().
			Err(err).
			Str("path", requestInfoDto.Path).
			Msg("saving request info error")
	}
}

func jsonFromMap[T any](data map[string]T) (datatypes.JSON, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return datatypes.JSON{}, err
	}
	return datatypes.JSON(bytes), nil
}
