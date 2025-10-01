package storage

import (
	"context"

	"github.com/ProRocketeers/url-shortener/domain/dto"
	"github.com/ProRocketeers/url-shortener/domain/model"
	"github.com/ProRocketeers/url-shortener/domain/query"
	"gorm.io/gen/field"
	"gorm.io/gorm"
)

type Repository struct {
	DB *query.Query
}

type ShortLinkRepository struct {
	Repository
}

func (r *ShortLinkRepository) Create(ctx context.Context, link *model.ShortLink) error {
	return r.DB.WithContext(ctx).ShortLink.Create(link)
}

func (r *ShortLinkRepository) FindBySlug(ctx context.Context, slug string) (*model.ShortLink, error) {
	return r.DB.WithContext(ctx).ShortLink.Where(r.DB.ShortLink.Slug.Eq(slug)).First()
}

func (r *ShortLinkRepository) FindById(ctx context.Context, id uint) (*model.ShortLink, error) {
	return r.DB.WithContext(ctx).ShortLink.Where(r.DB.ShortLink.ID.Eq(id)).First()
}

func (r *ShortLinkRepository) UpdateById(ctx context.Context, id uint, d dto.ShortLinkUpdateDTO) (*model.ShortLink, error) {
	updates := []field.AssignExpr{}
	if d.Slug != nil {
		updates = append(updates, r.DB.ShortLink.Slug.Value(*d.Slug))
	}
	if d.OriginalURL != nil {
		updates = append(updates, r.DB.ShortLink.OriginalURL.Value(*d.OriginalURL))
	}
	// no real way to determine if "expire = nil" means "not specified" or "delete expiration", so we'll go with consistency and meaning "not specified"
	if d.ExpiresAt != nil {
		updates = append(updates, r.DB.ShortLink.ExpiresAt.Value(*d.ExpiresAt))
	}

	// first check if the link exists
	count, err := r.DB.WithContext(ctx).ShortLink.Where(r.DB.ShortLink.ID.Eq(id)).Count()
	if count == 0 || err != nil {
		return nil, gorm.ErrRecordNotFound
	}

	var link model.ShortLink
	_, err = r.DB.WithContext(ctx).ShortLink.Returning(&link).Where(r.DB.ShortLink.ID.Eq(id)).UpdateSimple(updates...)
	return &link, err
}

func (r *ShortLinkRepository) DeleteById(ctx context.Context, id uint) (int64, error) {
	info, _ := r.DB.WithContext(ctx).ShortLink.Where(r.DB.ShortLink.ID.Eq(id)).Delete()
	return info.RowsAffected, info.Error
}

// ------------------------

type RequestInfoRepository struct {
	Repository
}

func (r *RequestInfoRepository) Create(ctx context.Context, info *model.RequestInfo) error {
	return r.DB.WithContext(ctx).RequestInfo.Create(info)
}
