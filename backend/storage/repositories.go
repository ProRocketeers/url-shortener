package storage

import (
	"context"
	"fmt"
	"strings"

	"github.com/ProRocketeers/url-shortener/domain/dto"
	"github.com/ProRocketeers/url-shortener/domain/model"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

type ShortLinkRepository struct {
	Repository
}

func (r *ShortLinkRepository) Create(ctx context.Context, link *model.ShortLink) error {
	return gorm.G[model.ShortLink](r.DB).Create(ctx, link)
}

func (r *ShortLinkRepository) FindBySlug(ctx context.Context, slug string) (model.ShortLink, error) {
	return gorm.G[model.ShortLink](r.DB).Where("slug = ?", slug).First(ctx)
}

func (r *ShortLinkRepository) FindById(ctx context.Context, id uint) (model.ShortLink, error) {
	return gorm.G[model.ShortLink](r.DB).Where("id = ?", id).First(ctx)
}

func (r *ShortLinkRepository) UpdateById(ctx context.Context, id uint, d dto.ShortLinkUpdateDTO) (model.ShortLink, error) {
	// TODO: try gorm.io/gen ?
	// should be a client-side tool to code-gen type-safe DB queries
	// because the docs suck
	setClauses := []string{}
	parameters := []any{}

	if d.Slug != nil {
		setClauses = append(setClauses, "slug = ?")
		parameters = append(parameters, *d.Slug)
	}
	if d.OriginalURL != nil {
		setClauses = append(setClauses, "original_url = ?")
		parameters = append(parameters, *d.OriginalURL)
	}
	// no real way to determine if "expire = nil" means "not specified" or "delete expiration", so we'll go with consistency and meaning "not specified"
	if d.ExpiresAt != nil {
		setClauses = append(setClauses, "expires_at = ?")
		parameters = append(parameters, *d.ExpiresAt)
	}
	parameters = append(parameters, id)

	return gorm.G[model.ShortLink](r.DB).Raw(
		fmt.Sprintf("UPDATE short_links SET %v WHERE id = ? RETURNING *", strings.Join(setClauses, ", ")),
		parameters...,
	).First(ctx)
}

func (r *ShortLinkRepository) DeleteById(ctx context.Context, id uint) (int, error) {
	return gorm.G[model.ShortLink](r.DB).Where("id = ?", id).Delete(ctx)
}

// ------------------------

type RequestInfoRepository struct {
	Repository
}

func (r *RequestInfoRepository) Create(ctx context.Context, info *model.RequestInfo) error {
	return gorm.G[model.RequestInfo](r.DB).Create(ctx, info)
}
