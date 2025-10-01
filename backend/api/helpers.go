package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ProRocketeers/url-shortener/domain/dto"
	"github.com/ProRocketeers/url-shortener/domain/model"
	"github.com/go-chi/chi/v5/middleware"
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
	}
}
