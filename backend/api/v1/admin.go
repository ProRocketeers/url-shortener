package v1

import (
	"errors"
	"fmt"
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
//	@Tags			admin,links
//	@Accept			json
//	@Produce		json
//	@Param			body			body		createShortLinkRequest	true	"Request body"
//	@Success		200				{object}	shortLinkDto
//	@Failure		400				{object}	genericErrorResponse	"slug already used"
//	@Failure		500				{object}	genericErrorResponse
//	@Router			/v1/admin/link	[post]
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
//	@Tags			admin,links
//	@Produce		json
//	@Param			id						path		int	true	"ID"
//	@Success		200						{object}	shortLinkDto
//	@Failure		400						{object}	genericErrorResponse	"invalid link ID"
//	@Failure		404						{object}	genericErrorResponse	"link not found"
//	@Failure		500						{object}	genericErrorResponse
//	@Router			/v1/admin/link/id/{id}	[get]
func (h *AdminApiHandler) GetShortLinkById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	uid, err := strconv.ParseUint(id, 10, 64)
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
//	@Tags			admin,links
//	@Produce		json
//	@Param			slug						path		string	true	"slug"
//	@Success		200							{object}	shortLinkDto
//	@Failure		404							{object}	genericErrorResponse	"link not found"
//	@Failure		500							{object}	genericErrorResponse
//	@Router			/v1/admin/link/slug/{slug}	[get]
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
//	@Tags			admin,links
//	@Accept			json
//	@Produce		json
//	@Param			id						path		int						true	"ID"
//	@Param			body					body		updateShortLinkRequest	true	"fields to update"
//	@Success		200						{object}	shortLinkDto
//	@Failure		400						{object}	genericErrorResponse	"invalid request parameters"
//	@Failure		404						{object}	genericErrorResponse	"link not found"
//	@Failure		500						{object}	genericErrorResponse
//	@Router			/v1/admin/link/id/{id}	[put]
func (h *AdminApiHandler) UpdateShortLinkById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	uid, err := strconv.ParseUint(id, 10, 64)
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
//	@Tags			admin,links
//	@Produce		json
//	@Param			id	path	int	true	"ID"
//	@Success		200
//	@Failure		400						{object}	genericErrorResponse	"invalid link ID"
//	@Failure		404						{object}	genericErrorResponse	"link not found"
//	@Failure		500						{object}	genericErrorResponse
//	@Router			/v1/admin/link/id/{id}	[delete]
func (h *AdminApiHandler) DeleteShortLinkById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	uid, err := strconv.ParseUint(id, 10, 64)
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

// ShortenUrl godoc
//
//	@Summary		Lists short links
//	@Description	Supports either combination (offset + limit) or (page size + page) or no pagination
//	@Description	Either pair must have either both set, or both unset
//	@Description	If both pairs are supplied, page size + page is used
//	@Tags			admin,links
//	@Accept			json
//	@Produce		json
//	@Param			size				query		integer	false	"size"
//	@Param			page				query		integer	false	"page"
//	@Param			offset				query		integer	false	"offset"
//	@Param			limit				query		integer	false	"limit"
//	@Success		200					{object}	listShortLinksResponse
//	@Failure		400					{object}	genericErrorResponse	"invalid request parameters"
//	@Failure		500					{object}	genericErrorResponse
//	@Router			/v1/admin/link/list	[get]
func (h *AdminApiHandler) ListShortLinks(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()

	sizeQ := queryValues.Get("size")
	pageQ := queryValues.Get("page")
	offsetQ := queryValues.Get("offset")
	limitQ := queryValues.Get("limit")

	offset, limit, err := resolveRequestInfoListParams(sizeQ, pageQ, offsetQ, limitQ)
	if err != nil {
		sendJsonError(w, fmt.Sprintf("input validation error: %s", err.Error()), http.StatusBadRequest)
		return
	}

	links, pagination, err := h.ShortLinkService.ListShortLinks(r.Context(), offset, limit)
	if err != nil {
		sendJsonError(w, "internal error", http.StatusInternalServerError)
		return
	}

	sendJsonBody(w, listShortLinksResponse{
		Data: func() []shortLinkDto {
			ret := []shortLinkDto{}
			for _, link := range links {
				ret = append(ret, createShortLinkDto(link, h.ShortLinkService.GetShortUrl(link)))
			}
			return ret
		}(),
		Pagination: pagination,
	})
}

// ShortenUrl godoc
//
//	@Summary		Finds a request info by ID or request ID
//	@Description	If both params are supplied, ID is used and request ID is ignored
//	@Tags			admin,info
//	@Produce		json
//	@Param			id				query		int		false	"ID"
//	@Param			requestId		query		string	false	"request ID"
//	@Success		200				{object}	dto.RequestInfoDTO
//	@Failure		400				{object}	genericErrorResponse	"invalid request parameters"
//	@Failure		404				{object}	genericErrorResponse	"request info not found"
//	@Failure		500				{object}	genericErrorResponse
//	@Router			/v1/admin/info	[get]
func (h *AdminApiHandler) FindSingleRequestInfo(w http.ResponseWriter, r *http.Request) {
	var (
		uid uint64
		err error
	)
	id := r.URL.Query().Get("id")
	if id != "" {
		uid, err = strconv.ParseUint(id, 10, 64)
		if err != nil {
			sendJsonError(w, "invalid request info ID", http.StatusBadRequest)
			return
		}
	}
	requestId := r.URL.Query().Get("requestId")

	if uid == 0 && requestId == "" {
		sendJsonError(w, "must pass either 'id' or 'requestId' query parameter", http.StatusBadRequest)
		return
	}

	info, err := h.RequestInfoService.FindByIdOrRequestId(r.Context(), uint(uid), requestId)
	if err != nil {
		var e *domain.RequestInfoError
		if errors.As(err, &e) {
			switch e.Code {
			case domain.ErrorCodeInfoNotFound:
				sendJsonError(w, "request info not found", http.StatusNotFound)
			case domain.ErrorCodeInfoOther:
				sendJsonError(w, "internal error", http.StatusInternalServerError)
			default:
				sendJsonError(w, "internal error", http.StatusInternalServerError)
			}
		} else {
			sendJsonError(w, "internal error", http.StatusInternalServerError)
		}
		return
	}
	sendJsonBody(w, createRequestInfoDto(info))
}

// ShortenUrl godoc
//
//	@Summary		Lists request infos
//	@Description	Supports either combination (offset + limit) or (page size + page) or no pagination
//	@Description	Either pair must have either both set, or both unset
//	@Description	If both pairs are supplied, page size + page is used
//	@Tags			admin,info
//	@Accept			json
//	@Produce		json
//	@Param			size				query		integer	false	"size"
//	@Param			page				query		integer	false	"page"
//	@Param			offset				query		integer	false	"offset"
//	@Param			limit				query		integer	false	"limit"
//	@Success		200					{object}	listRequestInfoResponse
//	@Failure		400					{object}	genericErrorResponse	"invalid request parameters"
//	@Failure		500					{object}	genericErrorResponse
//	@Router			/v1/admin/info/list	[get]
func (h *AdminApiHandler) ListRequestInfos(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()

	sizeQ := queryValues.Get("size")
	pageQ := queryValues.Get("page")
	offsetQ := queryValues.Get("offset")
	limitQ := queryValues.Get("limit")

	offset, limit, err := resolveRequestInfoListParams(sizeQ, pageQ, offsetQ, limitQ)

	if err != nil {
		sendJsonError(w, fmt.Sprintf("input validation error: %s", err.Error()), http.StatusBadRequest)
		return
	}

	infos, pagination, err := h.RequestInfoService.ListRequestInfos(r.Context(), offset, limit)

	if err != nil {
		sendJsonError(w, "internal error", http.StatusInternalServerError)
		return
	}
	sendJsonBody(w, listRequestInfoResponse{
		Data: func() []requestInfoDto {
			ret := []requestInfoDto{}
			for _, info := range infos {
				ret = append(ret, createRequestInfoDto(info))
			}
			return ret
		}(),
		Pagination: pagination,
	})
}

// ShortenUrl godoc
//
//	@Summary		Lists request infos for a specific slug
//	@Description	Supports either combination (offset + limit) or (page size + page) or no pagination
//	@Description	Either pair must have either both set, or both unset
//	@Description	If both pairs are supplied, page size + page is used
//	@Tags			admin,info
//	@Accept			json
//	@Produce		json
//	@Param			slug						path		string	true	"slug"
//	@Param			size						query		integer	false	"size"
//	@Param			page						query		integer	false	"page"
//	@Param			offset						query		integer	false	"offset"
//	@Param			limit						query		integer	false	"limit"
//	@Success		200							{object}	listRequestInfoResponse
//	@Failure		400							{object}	genericErrorResponse	"invalid request parameters"
//	@Failure		500							{object}	genericErrorResponse
//	@Router			/v1/admin/info/list/{slug}	[get]
func (h *AdminApiHandler) ListRequestInfosBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	queryValues := r.URL.Query()

	sizeQ := queryValues.Get("size")
	pageQ := queryValues.Get("page")
	offsetQ := queryValues.Get("offset")
	limitQ := queryValues.Get("limit")

	offset, limit, err := resolveRequestInfoListParams(sizeQ, pageQ, offsetQ, limitQ)

	if err != nil {
		sendJsonError(w, fmt.Sprintf("input validation error: %s", err.Error()), http.StatusBadRequest)
		return
	}

	infos, pagination, err := h.RequestInfoService.ListRequestInfosBySlug(r.Context(), slug, offset, limit)

	if err != nil {
		sendJsonError(w, "internal error", http.StatusInternalServerError)
		return
	}
	sendJsonBody(w, listRequestInfoResponse{
		Data: func() []requestInfoDto {
			ret := []requestInfoDto{}
			for _, info := range infos {
				ret = append(ret, createRequestInfoDto(info))
			}
			return ret
		}(),
		Pagination: pagination,
	})
}
