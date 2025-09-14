package errors

import (
	"net/http"
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(code int, message, details string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

func ErrBadRequest(message string) *AppError {
	return NewAppError(http.StatusBadRequest, message, "")
}

func ErrUnauthorized(message string) *AppError {
	return NewAppError(http.StatusUnauthorized, message, "")
}

func ErrForbidden(message string) *AppError {
	return NewAppError(http.StatusForbidden, message, "")
}

func ErrNotFound(message string) *AppError {
	return NewAppError(http.StatusNotFound, message, "")
}

func ErrConflict(message string) *AppError {
	return NewAppError(http.StatusConflict, message, "")
}

func ErrInternalServerError(message string) *AppError {
	return NewAppError(http.StatusInternalServerError, message, "")
}

func ErrServiceUnavailable(message string) *AppError {
	return NewAppError(http.StatusServiceUnavailable, message, "")
}

var (
	ErrInvalidCredentials = ErrUnauthorized("Invalid credentials")
	ErrTokenExpired      = ErrUnauthorized("Token expired")
	ErrTokenInvalid      = ErrUnauthorized("Invalid token")
	ErrUserNotFound      = ErrNotFound("User not found")
	ErrEmailExists       = ErrConflict("Email already exists")
	ErrInvalidInput      = ErrBadRequest("Invalid input")
	ErrDatabaseError     = ErrInternalServerError("Database error occurred")
	ErrExternalService   = ErrServiceUnavailable("External service unavailable")
)