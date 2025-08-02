package storage

import (
	"context"

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
