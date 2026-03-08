package v1

import (
	"errors"
	"net/http"
	"time"

	"github.com/ProRocketeers/url-shortener/domain"
	"github.com/ProRocketeers/url-shortener/domain/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type ApiHandler struct {
	ShortLinkService   *services.ShortLinkService
	RequestInfoService *services.RequestInfoService

	validate *validator.Validate
}

func NewApiHandler(shortLinkService *services.ShortLinkService, requestInfoService *services.RequestInfoService) *ApiHandler {
	validate := validator.New(validator.WithRequiredStructEnabled())
	return &ApiHandler{shortLinkService, requestInfoService, validate}
}

// ShortenUrl godoc
//
//	@Summary		Shorten an URL
//	@Description	Returns a shortened link for the given URL
//	@Accept			json
//	@Produce		json
//	@Param			body		body		shortenUrlRequest	true	"Request body"
//	@Success		200			{object}	shortenUrlResponse
//	@Failure		400			{object}	genericErrorResponse	"slug already used"
//	@Failure		500			{object}	genericErrorResponse
//	@Router			/v1/shorten	[post]
func (h *ApiHandler) ShortenUrl(w http.ResponseWriter, r *http.Request) {
	defer func() {
		h.RequestInfoService.Create(r.Context(), getInfoFromRequest(r))
	}()

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
		Slug:     link.Slug,
	}
	sendJsonBody(w, res)
}

// RedirectSlug godoc
//
//	@Summary		Redirect from short link
//	@Description	Redirects the user to the original URL in the link
//	@Param			slug	path	string	true	"Slug"
//	@Success		307		"Temporary redirect to URL"
//	@Failure		404		{object}	genericErrorResponse	"link not found"
//	@Failure		404		{object}	genericErrorResponse	"link expired"
//	@Failure		500		{object}	genericErrorResponse	"internal error"
//	@Router			/v1/{slug} [get]
func (h *ApiHandler) RedirectSlug(w http.ResponseWriter, r *http.Request) {
	defer func() {
		h.RequestInfoService.Create(r.Context(), getInfoFromRequest(r))
	}()

	slug := chi.URLParam(r, "slug")
	link, err := h.ShortLinkService.FindBySlug(r.Context(), slug, true)
	if err != nil {
		var e *domain.ShortLinkError
		if errors.As(err, &e) {
			switch e.Code {
			case domain.ErrorCodeLinkNotFound:
				sendJsonError(w, "link not found", http.StatusNotFound)
			case domain.ErrorCodeLinkOther:
				sendJsonError(w, "internal error", http.StatusInternalServerError)
			case domain.ErrorCodeLinkExpired:
				sendJsonError(w, "link expired", http.StatusNotFound)
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

// GetShortLinkInfoBySlug godoc
//
//	@Summary		Get short link info by slug
//	@Description	Returns basic information about the short link including original URL and click count
//	@Produce		json
//	@Param			slug	path		string	true	"Slug"
//	@Success		200		{object}	shortLinkInfoResponse
//	@Failure		404		{object}	genericErrorResponse	"link not found"
//	@Failure		500		{object}	genericErrorResponse	"internal error"
//	@Router			/v1/info/{slug} [get]
func (h *ApiHandler) GetShortLinkInfoBySlug(w http.ResponseWriter, r *http.Request) {
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

	clickCount, err := h.RequestInfoService.CountBySlug(r.Context(), slug)
	if err != nil {
		sendJsonError(w, "internal error", http.StatusInternalServerError)
		return
	}

	sendJsonBody(w, shortLinkInfoResponse{
		OriginalURL: link.OriginalURL,
		ClickCount:  clickCount,
	})
}
