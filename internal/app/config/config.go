package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Google   GoogleConfig
	AI       AIConfig
	Email    EmailConfig
	Auth     AuthConfig
}

type ServerConfig struct {
	Port    string
	BaseURL string
	GinMode string
}

type DatabaseConfig struct {
	URL string
}

type RedisConfig struct {
	URL string
}

type GoogleConfig struct {
	ClientID     string
	ClientSecret string
}

type AIConfig struct {
	GeminiAPIKey string
	ClaudeAPIKey string
}

type EmailConfig struct {
	ResendAPIKey   string
	ResendFromEmail string
}

type AuthConfig struct {
	JWTSecret string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	config := &Config{
		Server: ServerConfig{
			Port:    getEnv("PORT", "8080"),
			BaseURL: getEnv("BASE_URL", "http://localhost:8080"),
			GinMode: getEnv("GIN_MODE", "debug"),
		},
		Database: DatabaseConfig{
			URL: mustGetEnv("DATABASE_URL"),
		},
		Redis: RedisConfig{
			URL: getEnv("REDIS_URL", "redis://localhost:6379"),
		},
		Google: GoogleConfig{
			ClientID:     mustGetEnv("GOOGLE_CLIENT_ID"),
			ClientSecret: mustGetEnv("GOOGLE_CLIENT_SECRET"),
		},
		AI: AIConfig{
			GeminiAPIKey: mustGetEnv("GEMINI_API_KEY"),
			ClaudeAPIKey: getEnv("CLAUDE_API_KEY", ""),
		},
		Email: EmailConfig{
			ResendAPIKey:   mustGetEnv("RESEND_API_KEY"),
			ResendFromEmail: mustGetEnv("RESEND_FROM_EMAIL"),
		},
		Auth: AuthConfig{
			JWTSecret: mustGetEnv("JWT_SECRET"),
		},
	}

	return config
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func mustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s is required", key)
	}
	return value
}