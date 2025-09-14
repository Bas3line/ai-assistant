package repository

import (
	"context"
	"database/sql"
	"time"

	"ai-assistant/internal/models"
	"ai-assistant/pkg/database"
)

type UserRepository struct {
	db *database.DB
}

func NewUserRepository(db *database.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (id, email, name, image, email_verified, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.ExecContext(ctx, query,
		user.ID, user.Email, user.Name, user.Image,
		user.EmailVerified, user.CreatedAt, user.UpdatedAt)
	return err
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, name, image, email_verified, created_at, updated_at
		FROM users WHERE id = $1
	`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.Name, &user.Image,
		&user.EmailVerified, &user.CreatedAt, &user.UpdatedAt)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, name, image, email_verified, created_at, updated_at
		FROM users WHERE email = $1
	`
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.Name, &user.Image,
		&user.EmailVerified, &user.CreatedAt, &user.UpdatedAt)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	user.UpdatedAt = time.Now()
	query := `
		UPDATE users 
		SET email = $2, name = $3, image = $4, email_verified = $5, updated_at = $6
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query,
		user.ID, user.Email, user.Name, user.Image,
		user.EmailVerified, user.UpdatedAt)
	return err
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}