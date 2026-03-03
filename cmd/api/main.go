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
	chi_middleware "github.com/go-chi/chi/v5/middleware"
	"github.com/hans/config-service/internal/config"
	"github.com/hans/config-service/internal/handler"
	"github.com/hans/config-service/internal/middleware"
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

	// ─── Repositories ───────────────────────────────────────────────────────────
	userRepo := repository.NewUserRepository(db)
	appRepo := repository.NewApplicationRepository(db)
	envRepo := repository.NewEnvironmentRepository(db)
	configRepo := repository.NewConfigurationRepository(db)
	versionRepo := repository.NewConfigVersionRepository(db)
	apiKeyRepo := repository.NewAPIKeyRepository(db)

	// ─── Services ────────────────────────────────────────────────────────────────
	authSvc := service.NewAuthService(userRepo, cfg.JWTSecret)
	userSvc := service.NewUserService(userRepo)
	appSvc := service.NewApplicationService(appRepo)
	envSvc := service.NewEnvironmentService(envRepo)
	configSvc := service.NewConfigurationService(configRepo, appRepo, envRepo, versionRepo)
	apiKeySvc := service.NewAPIKeyService(apiKeyRepo)

	// ─── Handlers ────────────────────────────────────────────────────────────────
	authH := handler.NewAuthHandler(authSvc)
	userH := handler.NewUserHandler(userSvc)
	baseH := handler.NewBaseHandler(appSvc, envSvc, configSvc, apiKeySvc)

	// ─── Router ──────────────────────────────────────────────────────────────────
	r := chi.NewRouter()

	// Global middleware
	r.Use(chi_middleware.RequestID)
	r.Use(chi_middleware.RealIP)
	r.Use(chi_middleware.Logger)
	r.Use(chi_middleware.Recoverer)
	r.Use(chi_middleware.Timeout(60 * time.Second))

	// CORS
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Authorization")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	// ─── Routes ──────────────────────────────────────────────────────────────────
	r.Get("/health", baseH.HealthCheck)

	r.Route("/api/v1", func(r chi.Router) {
		// Public: Auth
		r.Post("/auth/register", authH.Register)
		r.Post("/auth/login", authH.Login)

		// Protected routes (JWT required)
		r.Group(func(r chi.Router) {
			r.Use(middleware.JWTAuth(cfg.JWTSecret))

			r.Get("/dashboard/stats", baseH.GetDashboardStats)
			r.Get("/audit-logs", baseH.GetAuditLogs)

			// ── Users (ADMIN only) ──────────────────────────────────────────────
			r.Route("/users", func(r chi.Router) {
				r.Use(middleware.RequireRole("admin"))
				r.Get("/", userH.ListUsers)
				r.Get("/{id}", userH.GetUser)
				r.Put("/{id}", userH.UpdateUser)
				r.Delete("/{id}", userH.DeleteUser)
			})

			// ── Applications ────────────────────────────────────────────────────
			r.Route("/applications", func(r chi.Router) {
				r.Get("/", baseH.ListApplications)
				r.Get("/{id}", baseH.GetApplication)
				r.With(middleware.RequireRole("admin", "editor")).Post("/", baseH.CreateApplication)
				r.With(middleware.RequireRole("admin", "editor")).Put("/{id}", baseH.UpdateApplication)
				r.With(middleware.RequireRole("admin")).Delete("/{id}", baseH.DeleteApplication)
			})

			// ── Environments ─────────────────────────────────────────────────────
			r.Route("/environments", func(r chi.Router) {
				r.Get("/", baseH.ListEnvironments)
				r.Get("/{id}", baseH.GetEnvironment)
				r.With(middleware.RequireRole("admin", "editor")).Post("/", baseH.CreateEnvironment)
				r.With(middleware.RequireRole("admin", "editor")).Put("/{id}", baseH.UpdateEnvironment)
				r.With(middleware.RequireRole("admin")).Delete("/{id}", baseH.DeleteEnvironment)
			})

			// ── Configurations ───────────────────────────────────────────────────
			r.Route("/configs", func(r chi.Router) {
				r.Get("/", baseH.ListConfigurations)
				r.Get("/{id}", baseH.GetConfiguration)
				r.Get("/{id}/versions", baseH.GetConfigVersions)
				r.With(middleware.RequireRole("admin", "editor")).Post("/", baseH.CreateConfiguration)
				r.With(middleware.RequireRole("admin", "editor")).Put("/{id}", baseH.UpdateConfiguration)
				r.With(middleware.RequireRole("admin")).Delete("/{id}", baseH.DeleteConfiguration)
			})

			// ── API Keys ────────────────────────────────────────────────────────
			r.Route("/api-keys", func(r chi.Router) {
				r.Get("/", baseH.ListAPIKeys)
				r.Post("/", baseH.CreateAPIKey)
				r.Delete("/{id}", baseH.RevokeAPIKey)
			})
		})
	})

	// ─── Start Server ─────────────────────────────────────────────────────────
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("🚀 Envynce API running on port %s [%s]", cfg.Port, cfg.AppEnv)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server exited")
}
