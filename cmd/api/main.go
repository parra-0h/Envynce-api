package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/hans/config-service/internal/config"
	"github.com/hans/config-service/internal/handler"
	internal_middleware "github.com/hans/config-service/internal/middleware"
	"github.com/hans/config-service/internal/repository"
	"github.com/hans/config-service/internal/service"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize Database
	db, err := repository.InitDB(cfg.DSN)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize Repositories
	appRepo := repository.NewApplicationRepository(db)
	envRepo := repository.NewEnvironmentRepository(db)
	configRepo := repository.NewConfigurationRepository(db)

	// Initialize Services
	appSvc := service.NewApplicationService(appRepo)
	envSvc := service.NewEnvironmentService(envRepo)
	configSvc := service.NewConfigurationService(configRepo, appRepo, envRepo)

	// Initialize Handlers
	h := handler.NewBaseHandler(appSvc, envSvc, configSvc)

	// Setup Router
	r := chi.NewRouter()

	// Standard Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Routes
	r.Get("/health", h.HealthCheck)

	r.Group(func(r chi.Router) {
		r.Use(internal_middleware.APIKeyAuth(cfg))

		r.Route("/applications", func(r chi.Router) {
			r.Post("/", h.CreateApplication)
			r.Get("/", h.ListApplications)
		})

		r.Route("/environments", func(r chi.Router) {
			r.Post("/", h.CreateEnvironment)
			r.Get("/", h.ListEnvironments)
		})

		r.Route("/configs", func(r chi.Router) {
			r.Post("/", h.CreateConfiguration)
			r.Get("/", h.ListConfigurations)
		})

		r.Get("/audit-logs", h.GetAuditLogs)
	})

	// Start Server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	// Graceful shutdown
	go func() {
		log.Printf("Starting server on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
