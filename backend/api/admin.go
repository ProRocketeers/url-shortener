package api

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/ProRocketeers/url-shortener/domain"
	"github.com/ProRocketeers/url-shortener/domain/dto"
	"github.com/ProRocketeers/url-shortener/domain/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type AdminApiHandler struct {
	ShortLinkService   *services.ShortLinkService
	RequestInfoService *services.RequestInfoService

	validate *validator.Validate
}

func NewAdminApiHandler(shortLinkService *services.ShortLinkService, requestInfoService *services.RequestInfoService) *AdminApiHandler {
	validate := validator.New(validator.WithRequiredStructEnabled())
	return &AdminApiHandler{shortLinkService, requestInfoService, validate}
}

// ShortenUrl godoc
//
//	@Summary		Create a short link
//	@Description	Returns a shortened link for the given URL
//	@Tags			admin
//	@Accept			json
//	@Produce		json
//	@Param			body		body		createShortLinkRequest	true	"Request body"
//	@Success		200			{object}	shortLinkDto
//	@Failure		400			{object}	genericErrorResponse	"slug already used"
//	@Failure		500			{object}	genericErrorResponse
//	@Router			/admin/link	[post]
func (h *AdminApiHandler) CreateShortLink(w http.ResponseWriter, r *http.Request) {
	req, err := parseJsonBody[createShortLinkRequest](r)
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

	res := createShortLinkDto(link, h.ShortLinkService.GetShortUrl(link))
	sendJsonBody(w, res)
}

// ShortenUrl godoc
//
//	@Summary		Get a short link by ID
//	@Description	Returns a shortened link
//	@Tags			admin
//	@Accept			json
//	@Produce		json
//	@Param			id					path		int	true	"ID"
//	@Success		200					{object}	shortLinkDto
//	@Failure		400					{object}	genericErrorResponse	"invalid link ID"
//	@Failure		404					{object}	genericErrorResponse	"link not found"
//	@Failure		500					{object}	genericErrorResponse
//	@Router			/admin/link/id/{id}	[get]
func (h *AdminApiHandler) GetShortLinkById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	uid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		sendJsonError(w, "invalid link ID", http.StatusBadRequest)
		return
	}

	link, err := h.ShortLinkService.FindById(r.Context(), uint(uid))
	if err != nil {
		var e *domain.ShortLinkError
		if errors.As(err, &e) {
			switch e.Code {
			case domain.ErrorCodeLinkNotFound:
				sendJsonError(w, "link not found", http.StatusNotFound)
			case domain.ErrorCodeLinkOther:
				sendJsonError(w, "internal error", http.StatusInternalServerError)
			default:
				sendJsonError(w, "internal error", http.StatusInternalServerError)
			}
		} else {
			sendJsonError(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	res := createShortLinkDto(link, h.ShortLinkService.GetShortUrl(link))
	sendJsonBody(w, res)
}

// ShortenUrl godoc
//
//	@Summary		Get a short link by slug
//	@Description	Returns a shortened link
//	@Tags			admin
//	@Accept			json
//	@Produce		json
//	@Param			slug					path		string	true	"slug"
//	@Success		200						{object}	shortLinkDto
//	@Failure		404						{object}	genericErrorResponse	"link not found"
//	@Failure		500						{object}	genericErrorResponse
//	@Router			/admin/link/slug/{slug}	[get]
func (h *AdminApiHandler) GetShortLinkBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	link, err := h.ShortLinkService.FindBySlug(r.Context(), slug, false)
	if err != nil {
		var e *domain.ShortLinkError
		if errors.As(err, &e) {
			switch e.Code {
			case domain.ErrorCodeLinkNotFound:
				sendJsonError(w, "link not found", http.StatusNotFound)
			case domain.ErrorCodeLinkOther:
				sendJsonError(w, "internal error", http.StatusInternalServerError)
			default:
				sendJsonError(w, "internal error", http.StatusInternalServerError)
			}
		} else {
			sendJsonError(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	res := createShortLinkDto(link, h.ShortLinkService.GetShortUrl(link))
	sendJsonBody(w, res)
}

// ShortenUrl godoc
//
//	@Summary		Update a short link by ID
//	@Description	Returns updated shortened link
//	@Tags			admin
//	@Accept			json
//	@Produce		json
//	@Param			id					path		int						true	"ID"
//	@Param			body				body		updateShortLinkRequest	true	"fields to update"
//	@Success		200					{object}	shortLinkDto
//	@Failure		400					{object}	genericErrorResponse	"invalid request parameters"
//	@Failure		404					{object}	genericErrorResponse	"link not found"
//	@Failure		500					{object}	genericErrorResponse
//	@Router			/admin/link/id/{id}	[put]
func (h *AdminApiHandler) UpdateShortLinkById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	uid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		sendJsonError(w, "invalid link ID", http.StatusBadRequest)
		return
	}

	req, err := parseJsonBody[updateShortLinkRequest](r)
	if err != nil {
		sendJsonError(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	// validate fields with struct tags
	if err := h.validate.Struct(req); err != nil {
		sendJsonError(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if req.ExpiresAt == nil && req.Slug == nil && req.OriginalURL == nil {
		sendJsonError(w, "no fields to update", http.StatusBadRequest)
		return
	}

	link, err := h.ShortLinkService.UpdateById(r.Context(), uint(uid), dto.ShortLinkUpdateDTO{
		OriginalURL: req.OriginalURL,
		Slug:        req.Slug,
		ExpiresAt:   req.ExpiresAt,
	})
	if err != nil {
		var e *domain.ShortLinkError
		if errors.As(err, &e) {
			switch e.Code {
			case domain.ErrorCodeLinkNotFound:
				sendJsonError(w, "link not found", http.StatusNotFound)
			case domain.ErrorCodeLinkOther:
				sendJsonError(w, "internal error", http.StatusInternalServerError)
			default:
				sendJsonError(w, "internal error", http.StatusInternalServerError)
			}
		} else {
			sendJsonError(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	res := createShortLinkDto(link, h.ShortLinkService.GetShortUrl(link))
	sendJsonBody(w, res)
}

// ShortenUrl godoc
//
//	@Summary		Delete a short link by ID
//	@Description	Deletes a shortened link
//	@Tags			admin
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"ID"
//	@Success		200
//	@Failure		400					{object}	genericErrorResponse	"invalid link ID"
//	@Failure		404					{object}	genericErrorResponse	"link not found"
//	@Failure		500					{object}	genericErrorResponse
//	@Router			/admin/link/id/{id}	[delete]
func (h *AdminApiHandler) DeleteShortLinkById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	uid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		sendJsonError(w, "invalid link ID", http.StatusBadRequest)
		return
	}

	err = h.ShortLinkService.DeleteById(r.Context(), uint(uid))
	if err != nil {
		var e *domain.ShortLinkError
		if errors.As(err, &e) {
			switch e.Code {
			case domain.ErrorCodeLinkNotFound:
				sendJsonError(w, "link not found", http.StatusNotFound)
			case domain.ErrorCodeLinkOther:
				sendJsonError(w, "internal error", http.StatusInternalServerError)
			default:
				sendJsonError(w, "internal error", http.StatusInternalServerError)
			}
		} else {
			sendJsonError(w, "internal error", http.StatusInternalServerError)
		}
		return
	}
}
