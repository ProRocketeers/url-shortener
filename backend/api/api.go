package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/ProRocketeers/url-shortener/domain"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type ApiHandler struct {
	ShortLinkService *domain.ShortLinkService

	validate *validator.Validate
}

func NewApiHandler(service *domain.ShortLinkService) *ApiHandler {
	validate := validator.New(validator.WithRequiredStructEnabled())
	return &ApiHandler{service, validate}
}

func (h *ApiHandler) ShortenUrl(w http.ResponseWriter, r *http.Request) {
	req, err := parseJsonBody[shortenUrlRequest](r)
	if err != nil {
		sendJsonError(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	// validate required fields with struct tags
	if err := h.validate.Struct(req); err != nil {
		sendJsonError(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	if req.ExpiresAt != nil {
		req.ExpiresAt = ptrTo(req.ExpiresAt.In(time.UTC))
	}

	link, err := h.ShortLinkService.Create(r.Context(), req.OriginalURL, req.Slug, req.ExpiresAt)
	if err != nil {
		var e *domain.ShortLinkError
		if errors.As(err, &e) {
			switch e.Code {
			case domain.ErrorCodeLinkCreate:
				// error happened while creating the link but we don't want to expose internals
				sendJsonError(w, "internal error", http.StatusInternalServerError)
			case domain.ErrorCodeSlugConflict:
				sendJsonError(w, "slug is already used", http.StatusBadRequest)
			default:
				sendJsonError(w, "internal error", http.StatusInternalServerError)
			}
		} else {
			sendJsonError(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	res := shortenUrlResponse{
		ShortURL: h.ShortLinkService.GetShortUrl(link),
	}
	sendJsonBody(w, res)
}

func (h *ApiHandler) RedirectSlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	link, err := h.ShortLinkService.FindBySlug(r.Context(), slug)
	if err != nil {
		var e *domain.ShortLinkError
		if errors.As(err, &e) {
			switch e.Code {
			case domain.ErrorCodeLinkNotFound:
				sendJsonError(w, "link not found", http.StatusNotFound)
			case domain.ErrorCodeLinkGetOther:
				sendJsonError(w, "internal error", http.StatusBadRequest)
			case domain.ErrorCodeLinkExpired:
				sendJsonError(w, "link expired", http.StatusBadRequest)
			default:
				sendJsonError(w, "internal error", http.StatusInternalServerError)
			}
		} else {
			sendJsonError(w, "internal error", http.StatusInternalServerError)
		}
		return
	}
	http.Redirect(w, r, link.OriginalURL, http.StatusTemporaryRedirect)
}
