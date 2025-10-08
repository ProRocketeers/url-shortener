package services

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/ProRocketeers/url-shortener/domain"
	"github.com/ProRocketeers/url-shortener/domain/dto"
	"github.com/ProRocketeers/url-shortener/domain/model"
	"github.com/ProRocketeers/url-shortener/domain/storage"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	"github.com/segmentio/ksuid"
)

type ShortLinkService struct {
	Repository *storage.ShortLinkRepository
	BaseUrl    url.URL
}

func (s *ShortLinkService) Create(ctx context.Context, originalUrl string, slug *string, expiresAt *time.Time) (model.ShortLink, error) {
	usedSlug := createSlug()
	if slug != nil {
		// first check if it exists or not
		_, err := s.Repository.FindBySlug(ctx, *slug)
		if err == nil || !errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ShortLink{}, &domain.ShortLinkError{Code: domain.ErrorCodeSlugConflict}
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
		return model.ShortLink{}, &domain.ShortLinkError{Code: domain.ErrorCodeLinkCreate}
	}

	return link, nil
}

func (s *ShortLinkService) FindBySlug(ctx context.Context, slug string, checkExpire bool) (model.ShortLink, error) {
	link, err := s.Repository.FindBySlug(ctx, slug)
	if link == nil {
		return model.ShortLink{}, &domain.ShortLinkError{Code: domain.ErrorCodeLinkNotFound}
	} else if err != nil {
		log.Error().
			Err(err).
			Str("slug", slug).
			Msg("getting link other error")
		return model.ShortLink{}, &domain.ShortLinkError{Code: domain.ErrorCodeLinkOther}
	}
	if checkExpire && link.ExpiresAt != nil && time.Now().UTC().After(*link.ExpiresAt) {
		return model.ShortLink{}, &domain.ShortLinkError{Code: domain.ErrorCodeLinkExpired}
	}
	return *link, nil
}

func (s *ShortLinkService) FindById(ctx context.Context, id uint) (model.ShortLink, error) {
	link, err := s.Repository.FindById(ctx, id)
	if link == nil {
		return model.ShortLink{}, &domain.ShortLinkError{Code: domain.ErrorCodeLinkNotFound}
	} else if err != nil {
		log.Error().
			Err(err).
			Uint("id", id).
			Msg("[admin] getting link other error")
		return model.ShortLink{}, &domain.ShortLinkError{Code: domain.ErrorCodeLinkOther}
	}
	return *link, nil
}

func (s *ShortLinkService) UpdateById(ctx context.Context, id uint, d dto.ShortLinkUpdateDTO) (model.ShortLink, error) {
	link, err := s.Repository.UpdateById(ctx, id, d)
	if link == nil {
		return model.ShortLink{}, &domain.ShortLinkError{Code: domain.ErrorCodeLinkNotFound}
	} else if err != nil {
		log.Error().
			Err(err).
			Uint("id", id).
			Msg("[admin] updating link other error")
		return model.ShortLink{}, &domain.ShortLinkError{Code: domain.ErrorCodeLinkOther}
	}
	return *link, nil
}

func (s *ShortLinkService) DeleteById(ctx context.Context, id uint) error {
	// does not return error on record not found
	rows, err := s.Repository.DeleteById(ctx, id)
	if rows == 0 {
		return &domain.ShortLinkError{Code: domain.ErrorCodeLinkNotFound}
	} else if err != nil {
		log.Error().
			Err(err).
			Uint("id", id).
			Msg("[admin] deleting link other error")
		return &domain.ShortLinkError{Code: domain.ErrorCodeLinkOther}
	}
	return nil
}

func (s *ShortLinkService) GetShortUrl(link model.ShortLink) string {
	return fmt.Sprintf("%s/v1/%s", s.BaseUrl.String(), link.Slug)
}

func createSlug() string {
	return ksuid.New().String()[:8]
}
