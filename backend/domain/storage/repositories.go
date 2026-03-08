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

func (r *ShortLinkRepository) List(ctx context.Context) ([]model.ShortLink, int64, error) {
	ret := []model.ShortLink{}
	links, err := r.DB.WithContext(ctx).ShortLink.Order(r.DB.ShortLink.ID.Desc()).Find()

	for _, link := range links {
		ret = append(ret, *link)
	}
	return ret, int64(len(ret)), err
}

func (r *ShortLinkRepository) PaginatedList(ctx context.Context, offset, limit int) ([]model.ShortLink, int64, error) {
	ret := []model.ShortLink{}
	links, totalCount, err := r.DB.WithContext(ctx).ShortLink.Order(r.DB.ShortLink.ID.Desc()).FindByPage(offset, limit)

	for _, link := range links {
		ret = append(ret, *link)
	}
	return ret, totalCount, err
}

// ------------------------

type RequestInfoRepository struct {
	Repository
}

func (r *RequestInfoRepository) Create(ctx context.Context, info *model.RequestInfo) error {
	return r.DB.WithContext(ctx).RequestInfo.Create(info)
}

func (r *RequestInfoRepository) FindById(ctx context.Context, id uint) (*model.RequestInfo, error) {
	return r.DB.WithContext(ctx).RequestInfo.Where(r.DB.RequestInfo.ID.Eq(id)).First()
}

func (r *RequestInfoRepository) FindByRequestId(ctx context.Context, requestId string) (*model.RequestInfo, error) {
	return r.DB.WithContext(ctx).RequestInfo.Where(r.DB.RequestInfo.RequestId.Eq(requestId)).First()
}

func (r *RequestInfoRepository) List(ctx context.Context) ([]model.RequestInfo, int64, error) {
	ret := []model.RequestInfo{}
	infos, err := r.DB.WithContext(ctx).RequestInfo.Find()

	for _, info := range infos {
		ret = append(ret, *info)
	}
	return ret, int64(len(ret)), err
}

func (r *RequestInfoRepository) PaginatedList(ctx context.Context, offset, limit int) ([]model.RequestInfo, int64, error) {
	ret := []model.RequestInfo{}
	infos, totalCount, err := r.DB.WithContext(ctx).RequestInfo.FindByPage(offset, limit)

	for _, info := range infos {
		ret = append(ret, *info)
	}
	return ret, totalCount, err
}

func (r *RequestInfoRepository) ListBySlug(ctx context.Context, slug string) ([]model.RequestInfo, int64, error) {
	ret := []model.RequestInfo{}
	filter := r.DB.RequestInfo.Path.Like("%/v1/" + slug)
	infos, err := r.DB.WithContext(ctx).RequestInfo.Where(filter).Find()

	for _, info := range infos {
		ret = append(ret, *info)
	}
	return ret, int64(len(ret)), err
}

func (r *RequestInfoRepository) PaginatedListBySlug(ctx context.Context, slug string, offset, limit int) ([]model.RequestInfo, int64, error) {
	ret := []model.RequestInfo{}
	filter := r.DB.RequestInfo.Path.Like("%/v1/" + slug)
	infos, totalCount, err := r.DB.WithContext(ctx).RequestInfo.Where(filter).FindByPage(offset, limit)

	for _, info := range infos {
		ret = append(ret, *info)
	}
	return ret, totalCount, err
}

func (r *RequestInfoRepository) CountBySlug(ctx context.Context, slug string) (int64, error) {
	pathFilter := r.DB.RequestInfo.Path.Like("%/v1/" + slug)
	methodFilter := r.DB.RequestInfo.Method.Eq("GET")

	return r.DB.WithContext(ctx).RequestInfo.Where(pathFilter, methodFilter).Count()
}
