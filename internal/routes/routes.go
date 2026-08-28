package routes

import (
	"net/http"

	"datalab_api/internal/auth"
	"datalab_api/internal/doctor"
	"datalab_api/internal/taxonomy"
)

func Setup(doctorHandler *doctor.Handler, taxonomyHandler *taxonomy.Handler, authMiddleware *auth.Middleware) http.Handler {
	mux := http.NewServeMux()

	// Public — no API key required, used for uptime/health checks
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Everything else requires auth
	protected := http.NewServeMux()
	protected.HandleFunc("GET /api/v1/doctors", doctorHandler.List)
	protected.HandleFunc("GET /api/v1/doctors/{id}", doctorHandler.Get)
	protected.HandleFunc("GET /api/v1/locations", taxonomyHandler.ListLocations)
	protected.HandleFunc("GET /api/v1/specialties", taxonomyHandler.ListSpecialties)

	mux.Handle("/api/v1/", authMiddleware.Authenticate(protected))

	return mux
}
