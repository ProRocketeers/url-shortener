package services

import (
	"context"
	"encoding/json"

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

func jsonFromMap[T any](data map[string]T) (datatypes.JSON, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return datatypes.JSON{}, err
	}
	return datatypes.JSON(bytes), nil
}
