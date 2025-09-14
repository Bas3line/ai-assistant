package routes

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"ai-assistant/internal/handlers"
	internalMiddleware "ai-assistant/internal/middleware"
	"ai-assistant/internal/services/auth"
	"ai-assistant/pkg/cache"
)

// SetupRoutes configures and returns the router with all application routes
func SetupRoutes(
	authHandler *handlers.AuthHandler,
	aiHandler *handlers.AIHandler,
	emailHandler *handlers.EmailHandler,
	authService *auth.AuthService,
	redisService *cache.RedisService,
) chi.Router {
	router := chi.NewRouter()

	// Global middleware
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(internalMiddleware.CORS())
	router.Use(internalMiddleware.ErrorHandler())

	// Root endpoint
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"message": "AI Assistant API",
			"version": "1.0.0",
			"status":  "running",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	// Health check endpoint
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		healthStatus := map[string]interface{}{
			"status":    "healthy",
			"gemini":    "available",
			"claude":    "not configured",
			"redis":     "disconnected",
			"database":  "connected",
		}

		// Check Redis status
		if redisService != nil {
			if err := redisService.Connect(); err == nil {
				healthStatus["redis"] = "connected"
			}
		}

		statusCode := http.StatusOK
		if healthStatus["status"] == "unhealthy" {
			statusCode = http.StatusServiceUnavailable
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(healthStatus)
	})

	// API routes
	router.Route("/api", func(r chi.Router) {
		// Authentication routes (public)
		authHandler.RegisterRoutes(r)
		
		// AI routes (protected)
		r.Route("/ai", func(r chi.Router) {
			r.Use(authService.RequireAuth())
			aiHandler.RegisterRoutes(r)
		})

		// Email routes (protected)
		r.Route("/emails", func(r chi.Router) {
			r.Use(authService.RequireAuth())
			emailHandler.RegisterRoutes(r)
		})
	})

	// 404 handler
	router.NotFound(internalMiddleware.NotFoundHandler())

	return router
}