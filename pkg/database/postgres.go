package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
	"ai-assistant/internal/app/config"
	"ai-assistant/pkg/logger"
)

type DB struct {
	*sql.DB
	logger *logger.Logger
}

func New(cfg *config.Config) (*DB, error) {
	db, err := sql.Open("postgres", cfg.Database.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger := logger.New()
	logger.Info("Successfully connected to database")

	return &DB{
		DB:     db,
		logger: logger,
	}, nil
}

func (db *DB) Close() error {
	db.logger.Info("Closing database connection")
	return db.DB.Close()
}

func (db *DB) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return db.PingContext(ctx)
}