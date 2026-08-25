package taxonomy

import (
	"net/http"
	"strconv"

	"datalab_api/internal/httpx"
)

type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) ListLocations(w http.ResponseWriter, r *http.Request) {
	locations, err := h.repo.ListLocations(r.Context())
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to list locations")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]interface{}{"data": locations})
}

func (h *Handler) ListSpecialties(w http.ResponseWriter, r *http.Request) {
	var locationID *int64
	if raw := r.URL.Query().Get("location_id"); raw != "" {
		id, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			httpx.WriteError(w, http.StatusBadRequest, "invalid location_id")
			return
		}
		locationID = &id
	}

	specialties, err := h.repo.ListSpecialties(r.Context(), locationID)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to list specialties")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]interface{}{"data": specialties})
}