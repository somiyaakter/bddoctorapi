package main

import (
	"context"
	"datalab_api/internal/auth"
	"datalab_api/internal/config"
	"datalab_api/internal/database"
	"datalab_api/internal/doctor"
	"datalab_api/internal/routes"
	"datalab_api/internal/taxonomy"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx := context.Background()

	// Load application configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// PostgreSQL connection
	db := database.NewPostgresDB(ctx, cfg.DatabaseURL)
	defer db.Close()

	// Doctor handler
	doctorHandler := doctor.NewHandler(
		doctor.NewService(
			doctor.NewRepository(db),
		),
	)

	// Taxonomy handler
	taxonomyHandler := taxonomy.NewHandler(
		taxonomy.NewRepository(db),
	)

	// Authentication
	authRepo := auth.NewRepository(db)
	authMiddleware := auth.NewMiddleware(authRepo)

	// Routes
	handler := routes.Setup(
		doctorHandler,
		taxonomyHandler,
		authMiddleware,
	)

	// HTTP server
	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: handler,
	}

	// Start server
	go func() {
		log.Printf("API listening on :%s", cfg.Port)

		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	log.Println("shutting down...")

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("forced shutdown: %v", err)
	}

	log.Println("server stopped")
}