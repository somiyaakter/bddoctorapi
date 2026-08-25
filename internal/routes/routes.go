package routes

import (
	"net/http"

	"datalab_api/internal/doctor"
	"datalab_api/internal/taxonomy"
)

func Setup(doctorHandler *doctor.Handler, taxonomyHandler *taxonomy.Handler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/v1/doctors", doctorHandler.List)
	mux.HandleFunc("GET /api/v1/doctors/{id}", doctorHandler.Get)
	mux.HandleFunc("GET /api/v1/locations", taxonomyHandler.ListLocations)
	mux.HandleFunc("GET /api/v1/specialties", taxonomyHandler.ListSpecialties)

	return mux
}
