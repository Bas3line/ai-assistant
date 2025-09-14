package repository

import (
	"context"
	"database/sql"

	"ai-assistant/internal/models"
	"ai-assistant/pkg/database"
)

type AIConversationRepository struct {
	db *database.DB
}

func NewAIConversationRepository(db *database.DB) *AIConversationRepository {
	return &AIConversationRepository{db: db}
}

func (r *AIConversationRepository) Create(ctx context.Context, conversation *models.AIConversation) error {
	query := `
		INSERT INTO ai_conversations (id, email_id, prompt, response, provider, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, query,
		conversation.ID, conversation.EmailID, conversation.Prompt,
		conversation.Response, conversation.Provider, conversation.CreatedAt)
	return err
}

func (r *AIConversationRepository) GetByID(ctx context.Context, id string) (*models.AIConversation, error) {
	conversation := &models.AIConversation{}
	query := `
		SELECT id, email_id, prompt, response, provider, created_at
		FROM ai_conversations WHERE id = $1
	`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&conversation.ID, &conversation.EmailID, &conversation.Prompt,
		&conversation.Response, &conversation.Provider, &conversation.CreatedAt)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return conversation, err
}

func (r *AIConversationRepository) GetByProvider(ctx context.Context, provider string, limit int) ([]*models.AIConversation, error) {
	query := `
		SELECT id, email_id, prompt, response, provider, created_at
		FROM ai_conversations 
		WHERE provider = $1 
		ORDER BY created_at DESC 
		LIMIT $2
	`
	
	rows, err := r.db.QueryContext(ctx, query, provider, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversations []*models.AIConversation
	for rows.Next() {
		conversation := &models.AIConversation{}
		err := rows.Scan(
			&conversation.ID, &conversation.EmailID, &conversation.Prompt,
			&conversation.Response, &conversation.Provider, &conversation.CreatedAt)
		if err != nil {
			return nil, err
		}
		conversations = append(conversations, conversation)
	}
	
	return conversations, rows.Err()
}

func (r *AIConversationRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM ai_conversations WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

type EmailRepository struct {
	db *database.DB
}

func NewEmailRepository(db *database.DB) *EmailRepository {
	return &EmailRepository{db: db}
}

func (r *EmailRepository) Create(ctx context.Context, email *models.Email) error {
	query := `
		INSERT INTO emails (id, message_id, thread_id, subject, "from", "to", body, html_body, is_read, labels, created_at, user_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	_, err := r.db.ExecContext(ctx, query,
		email.ID, email.MessageID, email.ThreadID, email.Subject,
		email.From, email.To, email.Body, email.HTMLBody,
		email.IsRead, email.Labels, email.CreatedAt, email.UserID)
	return err
}

func (r *EmailRepository) GetByUserID(ctx context.Context, userID string, limit int, offset int) ([]*models.Email, error) {
	query := `
		SELECT id, message_id, thread_id, subject, "from", "to", body, html_body, is_read, labels, created_at, user_id
		FROM emails 
		WHERE user_id = $1 
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3
	`
	
	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var emails []*models.Email
	for rows.Next() {
		email := &models.Email{}
		err := rows.Scan(
			&email.ID, &email.MessageID, &email.ThreadID, &email.Subject,
			&email.From, &email.To, &email.Body, &email.HTMLBody,
			&email.IsRead, &email.Labels, &email.CreatedAt, &email.UserID)
		if err != nil {
			return nil, err
		}
		emails = append(emails, email)
	}
	
	return emails, rows.Err()
}

func (r *EmailRepository) UpdateReadStatus(ctx context.Context, id string, isRead bool) error {
	query := `UPDATE emails SET is_read = $2 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id, isRead)
	return err
}