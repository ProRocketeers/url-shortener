package services

import (
	"context"
	"encoding/json"
	"math"

	"github.com/ProRocketeers/url-shortener/domain"
	"github.com/ProRocketeers/url-shortener/domain/dto"
	"github.com/ProRocketeers/url-shortener/domain/model"
	"github.com/ProRocketeers/url-shortener/domain/storage"
	"github.com/rs/zerolog/log"
	"gorm.io/datatypes"
)

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

func (s *RequestInfoService) FindByIdOrRequestId(ctx context.Context, id uint, requestId string) (model.RequestInfo, error) {
	var (
		info *model.RequestInfo
		err  error
	)
	if id != 0 {
		info, err = s.Repository.FindById(ctx, id)
	} else {
		info, err = s.Repository.FindByRequestId(ctx, requestId)
	}

	if info == nil {
		return model.RequestInfo{}, &domain.RequestInfoError{Code: domain.ErrorCodeInfoNotFound}
	} else if err != nil {
		log.Error().
			Err(err).
			Uint("id", id).
			Msg("[admin] getting request info other error")
		return model.RequestInfo{}, &domain.RequestInfoError{Code: domain.ErrorCodeInfoOther}
	}
	return *info, nil
}

func (s *RequestInfoService) ListRequestInfos(ctx context.Context, offset, limit *int) ([]model.RequestInfo, *dto.PaginationInfoDTO, error) {
	var (
		infos      []model.RequestInfo
		total      int64
		err        error
		pagination *dto.PaginationInfoDTO
	)

	if offset == nil {
		infos, _, err = s.Repository.List(ctx)
		if err != nil {
			log.Error().
				Err(err).
				Msg("[admin] listing request info other error")
			return []model.RequestInfo{}, nil, &domain.RequestInfoError{Code: domain.ErrorCodeInfoOther}
		}
	} else {
		infos, total, err = s.Repository.PaginatedList(ctx, *offset, *limit)
		if err != nil {
			log.Error().
				Err(err).
				Int("offset", *offset).
				Int("limit", *limit).
				Msg("[admin] listing request info other error")
			return []model.RequestInfo{}, nil, &domain.RequestInfoError{Code: domain.ErrorCodeInfoOther}
		}

		currentPage := (*offset / *limit) + 1
		totalPages := int(math.Ceil(float64(total) / float64(*limit)))

		pagination = &dto.PaginationInfoDTO{
			TotalRecords: total,
			TotalPages:   totalPages,
			CurrentPage:  currentPage,
			PreviousPage: func() *int {
				if currentPage == 1 {
					return nil
				}
				p := currentPage - 1
				return &p
			}(),
			NextPage: func() *int {
				if currentPage == totalPages {
					return nil
				}
				p := currentPage + 1
				return &p
			}(),
			Offset: *offset,
			Limit:  *limit,
		}
	}
	return infos, pagination, nil
}

func (s *RequestInfoService) ListRequestInfosBySlug(ctx context.Context, slug string, offset, limit *int) ([]model.RequestInfo, *dto.PaginationInfoDTO, error) {
	var (
		infos      []model.RequestInfo
		total      int64
		err        error
		pagination *dto.PaginationInfoDTO
	)

	if offset == nil {
		infos, _, err = s.Repository.ListBySlug(ctx, slug)
		if err != nil {
			log.Error().
				Err(err).
				Str("slug", slug).
				Msg("[admin] listing request info by slug other error")
			return []model.RequestInfo{}, nil, &domain.RequestInfoError{Code: domain.ErrorCodeInfoOther}
		}
	} else {
		infos, total, err = s.Repository.PaginatedListBySlug(ctx, slug, *offset, *limit)
		if err != nil {
			log.Error().
				Err(err).
				Str("slug", slug).
				Int("offset", *offset).
				Int("limit", *limit).
				Msg("[admin] listing request info by slug other error")
			return []model.RequestInfo{}, nil, &domain.RequestInfoError{Code: domain.ErrorCodeInfoOther}
		}

		currentPage := (*offset / *limit) + 1
		totalPages := int(math.Ceil(float64(total) / float64(*limit)))

		pagination = &dto.PaginationInfoDTO{
			TotalRecords: total,
			TotalPages:   totalPages,
			CurrentPage:  currentPage,
			PreviousPage: func() *int {
				if currentPage == 1 {
					return nil
				}
				p := currentPage - 1
				return &p
			}(),
			NextPage: func() *int {
				if currentPage == totalPages {
					return nil
				}
				p := currentPage + 1
				return &p
			}(),
			Offset: *offset,
			Limit:  *limit,
		}
	}
	return infos, pagination, nil
}

func jsonFromMap[T any](data map[string]T) (datatypes.JSON, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return datatypes.JSON{}, err
	}
	return datatypes.JSON(bytes), nil
}
