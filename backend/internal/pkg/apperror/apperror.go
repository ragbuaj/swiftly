package apperror

import (
	"errors"
	"net/http"
)

// AppError is a custom error type for the entire application
type AppError struct {
	Code    string `json:"code"`    // Unique error identifier (e.g. ERR_NOT_FOUND)
	Message string `json:"message"` // User-friendly message
	Status  int    `json:"-"`       // Internal HTTP status code (not exported to JSON)
}

func (e *AppError) Error() string {
	return e.Message
}

// Predefined Domain Errors
var (
	ErrNotFound       = &AppError{Code: "NOT_FOUND", Message: "The requested resource was not found", Status: http.StatusNotFound}
	ErrUnauthorized   = &AppError{Code: "UNAUTHORIZED", Message: "You must be authenticated to perform this action", Status: http.StatusUnauthorized}
	ErrForbidden      = &AppError{Code: "FORBIDDEN", Message: "You don't have permission to access this resource", Status: http.StatusForbidden}
	ErrBadRequest     = &AppError{Code: "BAD_REQUEST", Message: "The request payload is invalid", Status: http.StatusBadRequest}
	ErrInternalServer = &AppError{Code: "INTERNAL_ERROR", Message: "An unexpected internal error occurred", Status: http.StatusInternalServerError}
	ErrConflict       = &AppError{Code: "CONFLICT", Message: "The resource already exists", Status: http.StatusConflict}
	ErrRateLimit      = &AppError{Code: "RATE_LIMIT", Message: "Too many requests. Please try again later", Status: http.StatusTooManyRequests}
)

// New creates a custom domain-specific error
func New(status int, code, message string) *AppError {
	return &AppError{
		Status:  status,
		Code:    code,
		Message: message,
	}
}

// MapToHTTP analyzes an error and returns the appropriate HTTP status and formatted error data
func MapToHTTP(err error) (int, string, interface{}) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Status, appErr.Message, map[string]string{
			"code": appErr.Code,
		}
	}

	// Fallback for unknown errors (log these internally, but don't leak details to user)
	return http.StatusInternalServerError, "An unexpected error occurred", map[string]string{
		"code": "INTERNAL_ERROR",
	}
}
