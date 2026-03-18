// Package types provides shared types for the OpenClaw SDK.
//
// This package contains error types following Go best practices:
//   - Custom error types with error code and cause chain
//   - Standard errors.Is() and errors.As() compatibility
//   - Typed errors for different failure categories
package types

import "errors"

// ErrorCode represents an error code for categorizing errors.
// Used with errors.Is() to check error categories.
type ErrorCode string

const (
	ErrCodeConnection ErrorCode = "CONNECTION_ERROR"
	ErrCodeAuth       ErrorCode = "AUTH_ERROR"
	ErrCodeTimeout    ErrorCode = "TIMEOUT"
	ErrCodeProtocol   ErrorCode = "PROTOCOL_ERROR"
	ErrCodeValidation ErrorCode = "VALIDATION_ERROR"
	ErrCodeTransport  ErrorCode = "TRANSPORT_ERROR"
	ErrCodeUnknown    ErrorCode = "UNKNOWN"
)

// OpenClawError is the base error interface for all OpenClaw SDK errors.
// Provides error code and error chain (cause) for proper error handling.
type OpenClawError interface {
	error
	Code() ErrorCode
	Unwrap() error
}

// BaseError is the base error struct implementing OpenClawError.
// Contains error code, message, and optional cause (wrapped error).
type BaseError struct {
	code    ErrorCode
	message string
	err     error
}

func (e *BaseError) Error() string   { return e.message }
func (e *BaseError) Code() ErrorCode { return e.code }
func (e *BaseError) Unwrap() error   { return e.err }

// ConnectionError represents a connection error.
// Returned when network connection fails or is lost.
type ConnectionError struct {
	*BaseError
}

// AuthError represents an authentication error.
// Returned when authentication fails or credentials are invalid.
type AuthError struct {
	*BaseError
}

// TimeoutError represents a timeout error.
// Returned when operations exceed their time limits.
type TimeoutError struct {
	*BaseError
}

// ProtocolError represents a protocol error.
// Returned when protocol-level errors occur (e.g., invalid frame format).
type ProtocolError struct {
	*BaseError
}

// ValidationError represents a validation error.
// Returned when input validation fails (e.g., missing required fields).
type ValidationError struct {
	*BaseError
}

// TransportError represents a transport error.
// Returned when the underlying transport (e.g., WebSocket) fails.
type TransportError struct {
	*BaseError
}

// NewError creates a new error with the given code, message, and cause
func NewError(code ErrorCode, message string, err error) OpenClawError {
	return &BaseError{code: code, message: message, err: err}
}

// NewConnectionError creates a new connection error
func NewConnectionError(message string, err error) OpenClawError {
	return &ConnectionError{&BaseError{ErrCodeConnection, message, err}}
}

// NewAuthError creates a new authentication error
func NewAuthError(message string, err error) OpenClawError {
	return &AuthError{&BaseError{ErrCodeAuth, message, err}}
}

// NewTimeoutError creates a new timeout error
func NewTimeoutError(message string, err error) OpenClawError {
	return &TimeoutError{&BaseError{ErrCodeTimeout, message, err}}
}

// NewProtocolError creates a new protocol error
func NewProtocolError(message string, err error) OpenClawError {
	return &ProtocolError{&BaseError{ErrCodeProtocol, message, err}}
}

// NewValidationError creates a new validation error
func NewValidationError(message string, err error) OpenClawError {
	return &ValidationError{&BaseError{ErrCodeValidation, message, err}}
}

// NewTransportError creates a new transport error
func NewTransportError(message string, err error) OpenClawError {
	return &TransportError{&BaseError{ErrCodeTransport, message, err}}
}

// Is checks if the error matches the given code
// Uses standard library errors.Is() with custom unwrap
func Is(err error, code ErrorCode) bool {
	var e OpenClawError
	if As(err, &e) {
		return e.Code() == code
	}
	return false
}

// As casts the error to OpenClawError
// Uses standard library errors.As() for proper type matching
func As(err error, target *OpenClawError) bool {
	return errors.As(err, target)
}
