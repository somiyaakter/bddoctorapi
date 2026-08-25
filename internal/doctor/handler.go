package doctor

import (
	"datalab_api/internal/httpx"
	"errors"
	"log"
	"net/http"
	"strconv"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	if page < 1 {
		page = 1
	}

	if pageSize < 1 {
		pageSize = 20
	}

	doctors, total, err := h.service.List(
		r.Context(),
		page,
		pageSize,
	)

	if err != nil {
		// DEBUG: show the real database/repository error
		log.Printf("LIST DOCTORS ERROR: %v", err)

		httpx.WriteError(
			w,
			http.StatusInternalServerError,
			"failed to list doctors",
		)
		return
	}

	if pageSize > 100 {
		pageSize = 100
	}

	totalPages := (total + pageSize - 1) / pageSize

	httpx.WriteJSON(
		w,
		http.StatusOK,
		map[string]interface{}{
			"data": doctors,
			"pagination": httpx.Pagination{
				Page:       page,
				PageSize:   pageSize,
				TotalItems: total,
				TotalPages: totalPages,
			},
		},
	)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(
		r.PathValue("id"),
		10,
		64,
	)

	if err != nil {
		httpx.WriteError(
			w,
			http.StatusBadRequest,
			"invalid doctor id",
		)
		return
	}

	d, err := h.service.GetByID(
		r.Context(),
		id,
	)

	if errors.Is(err, ErrNotFound) {
		httpx.WriteError(
			w,
			http.StatusNotFound,
			"doctor not found",
		)
		return
	}

	if err != nil {
		log.Printf(
			"GET DOCTOR ERROR id=%d: %v",
			id,
			err,
		)

		httpx.WriteError(
			w,
			http.StatusInternalServerError,
			"failed to fetch doctor",
		)
		return
	}

	httpx.WriteJSON(
		w,
		http.StatusOK,
		map[string]interface{}{
			"data": d,
		},
	)
}