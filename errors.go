package actorhub

import (
	"fmt"
)

// ActorHubError is the base error type for ActorHub SDK errors.
type ActorHubError struct {
	Message      string
	StatusCode   int
	ResponseData map[string]interface{}
	RequestID    string
}

func (e *ActorHubError) Error() string {
	parts := e.Message
	if e.StatusCode > 0 {
		parts = fmt.Sprintf("%s (HTTP %d)", parts, e.StatusCode)
	}
	if e.RequestID != "" {
		parts = fmt.Sprintf("%s [Request ID: %s]", parts, e.RequestID)
	}
	return parts
}

// AuthenticationError is raised when API key is invalid or missing.
type AuthenticationError struct {
	ActorHubError
}

// NewAuthenticationError creates a new AuthenticationError.
func NewAuthenticationError(message string, requestID string) *AuthenticationError {
	if message == "" {
		message = "Invalid or missing API key"
	}
	return &AuthenticationError{
		ActorHubError: ActorHubError{
			Message:    message,
			StatusCode: 401,
			RequestID:  requestID,
		},
	}
}

// RateLimitError is raised when rate limit is exceeded.
type RateLimitError struct {
	ActorHubError
	RetryAfter int
}

// NewRateLimitError creates a new RateLimitError.
func NewRateLimitError(message string, retryAfter int, requestID string) *RateLimitError {
	if message == "" {
		message = "Rate limit exceeded"
	}
	return &RateLimitError{
		ActorHubError: ActorHubError{
			Message:    message,
			StatusCode: 429,
			RequestID:  requestID,
		},
		RetryAfter: retryAfter,
	}
}

// ValidationError is raised when request validation fails.
type ValidationError struct {
	ActorHubError
	Errors map[string]interface{}
}

// NewValidationError creates a new ValidationError.
func NewValidationError(message string, errors map[string]interface{}, requestID string) *ValidationError {
	if message == "" {
		message = "Validation error"
	}
	return &ValidationError{
		ActorHubError: ActorHubError{
			Message:    message,
			StatusCode: 422,
			RequestID:  requestID,
		},
		Errors: errors,
	}
}

// NotFoundError is raised when requested resource is not found.
type NotFoundError struct {
	ActorHubError
}

// NewNotFoundError creates a new NotFoundError.
func NewNotFoundError(message string, requestID string) *NotFoundError {
	if message == "" {
		message = "Resource not found"
	}
	return &NotFoundError{
		ActorHubError: ActorHubError{
			Message:    message,
			StatusCode: 404,
			RequestID:  requestID,
		},
	}
}

// ServerError is raised when server returns 5xx error.
type ServerError struct {
	ActorHubError
}

// NewServerError creates a new ServerError.
func NewServerError(message string, statusCode int, requestID string) *ServerError {
	if message == "" {
		message = "Server error"
	}
	return &ServerError{
		ActorHubError: ActorHubError{
			Message:    message,
			StatusCode: statusCode,
			RequestID:  requestID,
		},
	}
}
