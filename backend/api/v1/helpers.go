package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ProRocketeers/url-shortener/domain/dto"
	"github.com/ProRocketeers/url-shortener/domain/model"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
)

func parseJsonBody[T any](r *http.Request) (T, error) {
	var zero T

	ctx := r.Context()
	bodyBytes := ctx.Value("body").([]byte)

	if err := json.Unmarshal(bodyBytes, &zero); err != nil {
		return zero, err
	}
	return zero, nil
}

func sendJsonBody[T any](w http.ResponseWriter, data T) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func sendJsonError(w http.ResponseWriter, message string, status int) {
	data := genericErrorResponse{
		Error: message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func ptrTo[T any](value T) *T {
	return &value
}

func getInfoFromRequest(r *http.Request) dto.RequestInfoDTO {
	ctx := r.Context()

	requestId := middleware.GetReqID(ctx)
	timestamp := time.Now().UTC()
	realIp := r.RemoteAddr
	userAgent := r.UserAgent()
	headers := r.Header
	path := r.URL.Path
	method := r.Method
	query := r.URL.Query()
	body := parseBody(r)

	return dto.RequestInfoDTO{
		RequestId: requestId,
		Timestamp: timestamp,
		RealIP:    &realIp,
		UserAgent: &userAgent,
		Headers:   headers,
		Path:      path,
		Method:    method,
		Query:     query,
		Body:      body,
	}
}

func parseBody(r *http.Request) map[string]any {
	ctx := r.Context()
	bodyBytes := ctx.Value("body").([]byte)

	if len(bodyBytes) == 0 {
		return map[string]any{}
	}

	var jsonBody map[string]any
	if err := json.Unmarshal(bodyBytes, &jsonBody); err == nil {
		return jsonBody
	}

	return map[string]any{"text": string(bodyBytes)}
}

func createShortLinkDto(l model.ShortLink, shortUrl string) shortLinkDto {
	return shortLinkDto{
		ID:          l.ID,
		OriginalURL: l.OriginalURL,
		Slug:        &l.Slug,
		ShortURL:    shortUrl,
		ExpiresAt:   l.ExpiresAt,
		CreatedAt:   l.CreatedAt,
		UpdatedAt:   l.UpdatedAt,
	}
}

func createRequestInfoDto(i model.RequestInfo) requestInfoDto {
	dto := requestInfoDto{
		ID:        i.ID,
		RequestId: i.RequestId,
		Timestamp: i.Timestamp,
		RealIP:    i.RealIP,
		UserAgent: i.UserAgent,
		Path:      i.Path,
		Method:    i.Method,
	}

	// probably always passes, since it's stored as `{}` and not `null`, but just to be sure
	if len(i.Headers) > 0 {
		if err := json.Unmarshal(i.Headers, &dto.Headers); err != nil {
			log.Warn().Err(err).Uint("id", i.ID).Msg("error unmarshalling request info headers")
		}
	}
	if len(i.Query) > 0 {
		if err := json.Unmarshal(i.Query, &dto.Query); err != nil {
			log.Warn().Err(err).Uint("id", i.ID).Msg("error unmarshalling request info query")
		}
	}
	if len(i.Body) > 0 {
		if err := json.Unmarshal(i.Body, &dto.Body); err != nil {
			log.Warn().Err(err).Uint("id", i.ID).Msg("error unmarshalling request info headers")
		}
	}
	return dto
}

func parseOptionalIntString(s string) (*int, error) {
	if s == "" {
		return nil, nil
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		return nil, err
	}
	return &i, nil
}

func resolveRequestInfoListParams(sizeI, pageI, offsetI, limitI string) (retOffset *int, retLimit *int, err error) {
	if (sizeI == "" && pageI != "") || (sizeI != "" && pageI == "") {
		return nil, nil, fmt.Errorf("both size and page must be supplied")
	}
	if (offsetI == "" && limitI != "") || (offsetI != "" && limitI == "") {
		return nil, nil, fmt.Errorf("both offset and limit must be supplied")
	}

	if sizeI != "" && pageI != "" {
		var size, page *int
		size, err = parseOptionalIntString(sizeI)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid size")
		}
		if *size < 1 {
			return nil, nil, fmt.Errorf("size must be > 0")
		}
		page, err = parseOptionalIntString(pageI)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid page")
		}
		if *page < 1 {
			return nil, nil, fmt.Errorf("page must be > 0")
		}

		retLimit = size
		o := (*page - 1) * (*size)
		retOffset = &o
	} else if offsetI != "" && limitI != "" {
		retOffset, err = parseOptionalIntString(offsetI)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid offset")
		}
		if *retOffset < 0 {
			return nil, nil, fmt.Errorf("offset must be >= 0")
		}
		retLimit, err = parseOptionalIntString(limitI)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid limit")
		}
		if *retLimit < 1 {
			return nil, nil, fmt.Errorf("limit must be > 0")
		}
	}
	return
}
