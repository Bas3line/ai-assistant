package models

import (
	"time"
)

type User struct {
	ID            string     `json:"id" db:"id"`
	Email         string     `json:"email" db:"email"`
	Name          *string    `json:"name" db:"name"`
	Image         *string    `json:"image" db:"image"`
	EmailVerified *time.Time `json:"emailVerified" db:"email_verified"`
	CreatedAt     time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time  `json:"updatedAt" db:"updated_at"`
}

type Account struct {
	ID                string  `json:"id" db:"id"`
	UserID            string  `json:"userId" db:"user_id"`
	Type              string  `json:"type" db:"type"`
	Provider          string  `json:"provider" db:"provider"`
	ProviderAccountID string  `json:"providerAccountId" db:"provider_account_id"`
	RefreshToken      *string `json:"refresh_token" db:"refresh_token"`
	AccessToken       *string `json:"access_token" db:"access_token"`
	ExpiresAt         *int    `json:"expires_at" db:"expires_at"`
	TokenType         *string `json:"token_type" db:"token_type"`
	Scope             *string `json:"scope" db:"scope"`
	IDToken           *string `json:"id_token" db:"id_token"`
	SessionState      *string `json:"session_state" db:"session_state"`
}

type Session struct {
	ID           string    `json:"id" db:"id"`
	SessionToken string    `json:"sessionToken" db:"session_token"`
	UserID       string    `json:"userId" db:"user_id"`
	Expires      time.Time `json:"expires" db:"expires"`
}

type Email struct {
	ID        string    `json:"id" db:"id"`
	MessageID string    `json:"messageId" db:"message_id"`
	ThreadID  *string   `json:"threadId" db:"thread_id"`
	Subject   *string   `json:"subject" db:"subject"`
	From      string    `json:"from" db:"from"`
	To        []string  `json:"to" db:"to"`
	Body      *string   `json:"body" db:"body"`
	HTMLBody  *string   `json:"htmlBody" db:"html_body"`
	IsRead    bool      `json:"isRead" db:"is_read"`
	Labels    []string  `json:"labels" db:"labels"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UserID    string    `json:"userId" db:"user_id"`
}

type AIConversation struct {
	ID        string    `json:"id" db:"id"`
	EmailID   *string   `json:"emailId" db:"email_id"`
	Prompt    string    `json:"prompt" db:"prompt"`
	Response  string    `json:"response" db:"response"`
	Provider  string    `json:"provider" db:"provider"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

type AIRequest struct {
	Prompt   string `json:"prompt" binding:"required"`
	Provider string `json:"provider,omitempty"`
}

type AIResponse struct {
	Response string `json:"response"`
	Provider string `json:"provider"`
}

type AuthUser struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
	Image string `json:"image,omitempty"`
}