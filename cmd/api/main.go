package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ai-assistant/internal/app/config"
	"ai-assistant/internal/handlers"
	"ai-assistant/internal/repository"
	"ai-assistant/internal/routes"
	"ai-assistant/internal/usecase"
	"ai-assistant/internal/services/ai/claude"
	"ai-assistant/internal/services/ai/gemini"
	"ai-assistant/internal/services/auth"
	"ai-assistant/pkg/cache"
	"ai-assistant/pkg/database"
	"ai-assistant/pkg/logger"
)

func main() {
	cfg := config.Load()
	appLogger := logger.New()

	db, err := database.New(cfg)
	if err != nil {
		appLogger.Error("Failed to connect to database:", err)
		os.Exit(1)
	}
	defer db.Close()

	redisService := cache.NewRedisService(cfg)
	if redisService != nil {
		if err := redisService.Connect(); err != nil {
			appLogger.Warn("Failed to connect to Redis, continuing without cache:", err)
		} else {
			defer redisService.Disconnect()
		}
	}

	geminiService, err := gemini.NewGeminiService(cfg)
	if err != nil {
		appLogger.Error("Failed to initialize Gemini service:", err)
		os.Exit(1)
	}
	defer geminiService.Close()

	claudeService := claude.NewClaudeService(cfg)
	if claudeService == nil {
		appLogger.Warn("Claude service not configured (API key missing)")
	}

	userRepo := repository.NewUserRepository(db)

	aiUsecase := usecase.NewAIUsecase(geminiService, claudeService, redisService)
	authUsecase := usecase.NewAuthUsecase(userRepo)

	authService := auth.NewAuthService(cfg)

	aiHandler := handlers.NewAIHandler(aiUsecase)
	authHandler := handlers.NewAuthHandler(authService)
	emailHandler := handlers.NewEmailHandler()

	// Setup routes
	router := routes.SetupRoutes(authHandler, aiHandler, emailHandler, authService, redisService)

	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	go func() {
		appLogger.Infof("Server starting on port %s", cfg.Server.Port)
		appLogger.Infof("Server URL: %s", cfg.Server.BaseURL)
		
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Error("Failed to start server:", err)
			os.Exit(1)
		}
	}()

	_ = authUsecase

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		appLogger.Error("Server forced to shutdown:", err)
		os.Exit(1)
	}

	appLogger.Info("Server exited")
}