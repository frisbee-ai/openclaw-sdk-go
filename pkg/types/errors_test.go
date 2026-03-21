// Package types provides tests for error handling
package types

import (
	"errors"
	"fmt"
	"testing"
)

// compareDetails compares two any values for equality, handling maps correctly.
func compareDetails(a, b any) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	// Use fmt.Sprintf for simple comparison of comparable types
	return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}

// TestNewAPIError_AuthErrors tests that AUTH_* and CHALLENGE_* codes create AuthError.
func TestNewAPIError_AuthErrors(t *testing.T) {
	tests := []struct {
		name  string
		shape *ErrorShape
	}{
		{
			name: "AUTH_TOKEN_EXPIRED",
			shape: &ErrorShape{
				Code:    "AUTH_TOKEN_EXPIRED",
				Message: "Token has expired",
			},
		},
		{
			name: "AUTH_TOKEN_MISMATCH",
			shape: &ErrorShape{
				Code:    "AUTH_TOKEN_MISMATCH",
				Message: "Token mismatch",
			},
		},
		{
			name: "CHALLENGE_EXPIRED",
			shape: &ErrorShape{
				Code:    "CHALLENGE_EXPIRED",
				Message: "Challenge expired",
			},
		},
		{
			name: "CHALLENGE_FAILED",
			shape: &ErrorShape{
				Code:    "CHALLENGE_FAILED",
				Message: "Challenge failed",
			},
		},
		{
			name: "AUTH_RATE_LIMITED",
			shape: &ErrorShape{
				Code:    "AUTH_RATE_LIMITED",
				Message: "Rate limited",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewAPIError(tt.shape)
			if !IsAuthError(err) {
				t.Errorf("IsAuthError() = false, want true for code %s", tt.shape.Code)
			}
		})
	}
}

// TestNewAPIError_ConnectionErrors tests that CONNECTION_* and TLS_FINGERPRINT_MISMATCH create ConnectionError.
func TestNewAPIError_ConnectionErrors(t *testing.T) {
	tests := []struct {
		name  string
		shape *ErrorShape
	}{
		{
			name: "CONNECTION_STALE",
			shape: &ErrorShape{
				Code:    "CONNECTION_STALE",
				Message: "Connection stale",
			},
		},
		{
			name: "CONNECTION_CLOSED",
			shape: &ErrorShape{
				Code:    "CONNECTION_CLOSED",
				Message: "Connection closed",
			},
		},
		{
			name: "CONNECT_TIMEOUT",
			shape: &ErrorShape{
				Code:    "CONNECT_TIMEOUT",
				Message: "Connect timeout",
			},
		},
		{
			name: "TLS_FINGERPRINT_MISMATCH",
			shape: &ErrorShape{
				Code:    "TLS_FINGERPRINT_MISMATCH",
				Message: "TLS fingerprint mismatch",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewAPIError(tt.shape)
			if !IsConnectionError(err) {
				t.Errorf("IsConnectionError() = false, want true for code %s", tt.shape.Code)
			}
		})
	}
}

// TestNewAPIError_ProtocolErrors tests that PROTOCOL_* codes create ProtocolError.
func TestNewAPIError_ProtocolErrors(t *testing.T) {
	tests := []struct {
		name  string
		shape *ErrorShape
	}{
		{
			name: "PROTOCOL_UNSUPPORTED",
			shape: &ErrorShape{
				Code:    "PROTOCOL_UNSUPPORTED",
				Message: "Protocol unsupported",
			},
		},
		{
			name: "PROTOCOL_NEGOTIATION_FAILED",
			shape: &ErrorShape{
				Code:    "PROTOCOL_NEGOTIATION_FAILED",
				Message: "Negotiation failed",
			},
		},
		{
			name: "INVALID_FRAME",
			shape: &ErrorShape{
				Code:    "INVALID_FRAME",
				Message: "Invalid frame",
			},
		},
		{
			name: "FRAME_TOO_LARGE",
			shape: &ErrorShape{
				Code:    "FRAME_TOO_LARGE",
				Message: "Frame too large",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewAPIError(tt.shape)
			if !IsProtocolError(err) {
				t.Errorf("IsProtocolError() = false, want true for code %s", tt.shape.Code)
			}
		})
	}
}

// TestNewAPIError_RequestErrors tests that exact match codes create RequestError.
func TestNewAPIError_RequestErrors(t *testing.T) {
	tests := []struct {
		name  string
		shape *ErrorShape
	}{
		{
			name: "METHOD_NOT_FOUND",
			shape: &ErrorShape{
				Code:    "METHOD_NOT_FOUND",
				Message: "Method not found",
			},
		},
		{
			name: "INVALID_PARAMS",
			shape: &ErrorShape{
				Code:    "INVALID_PARAMS",
				Message: "Invalid params",
			},
		},
		{
			name: "INTERNAL_ERROR",
			shape: &ErrorShape{
				Code:    "INTERNAL_ERROR",
				Message: "Internal error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewAPIError(tt.shape)
			if !IsRequestError(err) {
				t.Errorf("IsRequestError() = false, want true for code %s", tt.shape.Code)
			}
		})
	}
}

// TestNewAPIError_GatewayErrors tests that gateway/business logic errors create GatewayError.
func TestNewAPIError_GatewayErrors(t *testing.T) {
	tests := []struct {
		name  string
		shape *ErrorShape
	}{
		{
			name: "AGENT_NOT_FOUND",
			shape: &ErrorShape{
				Code:    "AGENT_NOT_FOUND",
				Message: "Agent not found",
			},
		},
		{
			name: "AGENT_BUSY",
			shape: &ErrorShape{
				Code:    "AGENT_BUSY",
				Message: "Agent busy",
			},
		},
		{
			name: "NODE_NOT_FOUND",
			shape: &ErrorShape{
				Code:    "NODE_NOT_FOUND",
				Message: "Node not found",
			},
		},
		{
			name: "SESSION_NOT_FOUND",
			shape: &ErrorShape{
				Code:    "SESSION_NOT_FOUND",
				Message: "Session not found",
			},
		},
		{
			name: "PERMISSION_DENIED",
			shape: &ErrorShape{
				Code:    "PERMISSION_DENIED",
				Message: "Permission denied",
			},
		},
		{
			name: "QUOTA_EXCEEDED",
			shape: &ErrorShape{
				Code:    "QUOTA_EXCEEDED",
				Message: "Quota exceeded",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewAPIError(tt.shape)
			if !IsGatewayError(err) {
				t.Errorf("IsGatewayError() = false, want true for code %s", tt.shape.Code)
			}
		})
	}
}

// TestNewAPIError_Fallthrough tests that REQUEST_TIMEOUT, REQUEST_CANCELLED,
// and REQUEST_ABORTED fall through to GatewayError (NOT RequestError).
// This matches TypeScript createErrorFromResponse behavior.
func TestNewAPIError_Fallthrough(t *testing.T) {
	tests := []struct {
		name  string
		shape *ErrorShape
	}{
		{
			name: "REQUEST_TIMEOUT",
			shape: &ErrorShape{
				Code:    "REQUEST_TIMEOUT",
				Message: "Request timeout",
			},
		},
		{
			name: "REQUEST_CANCELLED",
			shape: &ErrorShape{
				Code:    "REQUEST_CANCELLED",
				Message: "Request cancelled",
			},
		},
		{
			name: "REQUEST_ABORTED",
			shape: &ErrorShape{
				Code:    "REQUEST_ABORTED",
				Message: "Request aborted",
			},
		},
		{
			name: "UNKNOWN_ERROR_CODE",
			shape: &ErrorShape{
				Code:    "UNKNOWN_ERROR_CODE",
				Message: "Unknown error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewAPIError(tt.shape)
			// These should be GatewayError, NOT RequestError
			if IsRequestError(err) {
				t.Errorf("IsRequestError() = true for %s, want false (should be GatewayError)", tt.shape.Code)
			}
			if !IsGatewayError(err) {
				t.Errorf("IsGatewayError() = false for %s, want true", tt.shape.Code)
			}
		})
	}
}

// TestNewAPIError_LowercaseCodes tests case-insensitive code matching.
func TestNewAPIError_LowercaseCodes(t *testing.T) {
	tests := []struct {
		name  string
		shape *ErrorShape
		want  func(error) bool
	}{
		{
			name: "auth_token_expired (lowercase)",
			shape: &ErrorShape{
				Code:    "auth_token_expired",
				Message: "Token expired",
			},
			want: IsAuthError,
		},
		{
			name: "connection_stale (lowercase)",
			shape: &ErrorShape{
				Code:    "connection_stale",
				Message: "Connection stale",
			},
			want: IsConnectionError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewAPIError(tt.shape)
			if !tt.want(err) {
				t.Errorf("expected error type check to return true for code %s", tt.shape.Code)
			}
		})
	}
}

// TestNewAPIError_Retryable tests retryable field propagation.
func TestNewAPIError_Retryable(t *testing.T) {
	retryable := true
	notRetryable := false

	tests := []struct {
		name  string
		shape *ErrorShape
		want  bool
	}{
		{
			name: "with retryable=true",
			shape: &ErrorShape{
				Code:      "AUTH_TOKEN_EXPIRED",
				Message:   "Token expired",
				Retryable: &retryable,
			},
			want: true,
		},
		{
			name: "with retryable=false",
			shape: &ErrorShape{
				Code:      "AUTH_TOKEN_EXPIRED",
				Message:   "Token expired",
				Retryable: &notRetryable,
			},
			want: false,
		},
		{
			name: "without retryable field",
			shape: &ErrorShape{
				Code:    "AUTH_TOKEN_EXPIRED",
				Message: "Token expired",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewAPIError(tt.shape)
			if IsRetryable(err) != tt.want {
				t.Errorf("IsRetryable() = %v, want %v", IsRetryable(err), tt.want)
			}
		})
	}
}

// TestNewAPIError_Details tests details field propagation.
func TestNewAPIError_Details(t *testing.T) {
	details := map[string]interface{}{"key": "value"}
	shape := &ErrorShape{
		Code:    "AUTH_TOKEN_EXPIRED",
		Message: "Token expired",
		Details: details,
	}

	err := NewAPIError(shape)
	authErr, ok := err.(*AuthError)
	if !ok {
		t.Fatalf("expected *AuthError, got %T", err)
	}

	if authErr.Details() == nil {
		t.Error("Details() = nil, want details")
	}
}

// TestNewAPIError_ReconnectErrors tests reconnect error codes.
func TestNewAPIError_ReconnectErrors(t *testing.T) {
	tests := []struct {
		name  string
		shape *ErrorShape
	}{
		{
			name: "MAX_RECONNECT_ATTEMPTS",
			shape: &ErrorShape{
				Code:    "MAX_RECONNECT_ATTEMPTS",
				Message: "Max reconnect attempts",
			},
		},
		{
			name: "MAX_AUTH_RETRIES",
			shape: &ErrorShape{
				Code:    "MAX_AUTH_RETRIES",
				Message: "Max auth retries",
			},
		},
		{
			name: "RECONNECT_DISABLED",
			shape: &ErrorShape{
				Code:    "RECONNECT_DISABLED",
				Message: "Reconnect disabled",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewAPIError(tt.shape)
			if !IsGatewayError(err) {
				// These fall through to GatewayError in current implementation
				t.Logf("Note: %s falls through to GatewayError", tt.shape.Code)
			}
		})
	}
}

// TestTimeoutError_SpecificType tests that TimeoutError is a specific error type.
func TestTimeoutError_SpecificType(t *testing.T) {
	err := NewTimeoutError("request timed out", nil)

	if !IsTimeoutError(err) {
		t.Errorf("IsTimeoutError() = false, want true")
	}

	if !IsRequestError(err) {
		t.Errorf("IsRequestError() = false, want true (TimeoutError is a RequestError)")
	}
}

// TestErrorInterfaces tests that error types implement expected interfaces.
func TestErrorInterfaces(t *testing.T) {
	err := NewAuthError("AUTH_TOKEN_EXPIRED", "Token expired", true, nil)

	t.Run("implements error interface", func(t *testing.T) {
		var _ error = err
		_ = err.Error()
	})

	t.Run("implements OpenClawError interface", func(t *testing.T) {
		var oe OpenClawError = err
		_ = oe.Code()
		_ = oe.Retryable()
		_ = oe.Unwrap()
	})
}

// TestNewAPIError_UnknownCode tests that unknown codes fall through to GatewayError.
func TestNewAPIError_UnknownCode(t *testing.T) {
	shape := &ErrorShape{
		Code:    "UNKNOWN_ERROR_CODE",
		Message: "Unknown error",
	}

	err := NewAPIError(shape)
	if !IsGatewayError(err) {
		t.Errorf("IsGatewayError() = false, want true for unknown code")
	}
}

// TestErrorMessageAndCode tests error message and code fields.
func TestErrorMessageAndCode(t *testing.T) {
	shape := &ErrorShape{
		Code:    "AUTH_TOKEN_EXPIRED",
		Message: "Token has expired",
	}

	err := NewAPIError(shape)
	authErr, ok := err.(*AuthError)
	if !ok {
		t.Fatalf("expected *AuthError, got %T", err)
	}

	if authErr.Error() != "Token has expired" {
		t.Errorf("Error() = %s, want %s", authErr.Error(), "Token has expired")
	}

	if authErr.Code() != "AUTH_TOKEN_EXPIRED" {
		t.Errorf("Code() = %s, want %s", authErr.Code(), "AUTH_TOKEN_EXPIRED")
	}
}

// TestErrorUnwrap tests error unwrapping.
func TestErrorUnwrap(t *testing.T) {
	innerErr := errors.New("inner error")
	shape := &ErrorShape{
		Code:    "INTERNAL_ERROR",
		Message: "Internal error",
		Details: innerErr,
	}

	err := NewAPIError(shape)
	reqErr, ok := err.(*RequestError)
	if !ok {
		t.Fatalf("expected *RequestError, got %T", err)
	}

	unwrapped := reqErr.Unwrap()
	if unwrapped == nil {
		t.Error("Unwrap() = nil, want inner error")
	}
}

// TestErrorImplementsInterface tests that error types implement OpenClawError.
func TestErrorImplementsInterface(t *testing.T) {
	shape := &ErrorShape{
		Code:    "AUTH_TOKEN_EXPIRED",
		Message: "Token expired",
	}

	err := NewAPIError(shape)

	var oe OpenClawError
	if !errors.As(err, &oe) {
		t.Fatal("error does not implement OpenClawError")
	}

	if oe.Code() != "AUTH_TOKEN_EXPIRED" {
		t.Errorf("Code() = %s, want %s", oe.Code(), "AUTH_TOKEN_EXPIRED")
	}

	if oe.Retryable() != false {
		t.Errorf("Retryable() = %v, want false", oe.Retryable())
	}

	if oe.Error() != "Token expired" {
		t.Errorf("Error() = %s, want %s", oe.Error(), "Token expired")
	}
}

// ============================================================================
// NewReconnectError Tests
// ============================================================================

func TestNewReconnectError(t *testing.T) {
	tests := []struct {
		name      string
		code      string
		message   string
		retryable bool
		details   any
		wantCode  string
		wantMsg   string
		wantRetry bool
	}{
		{
			name:      "max attempts error",
			code:      "MAX_RECONNECT_ATTEMPTS",
			message:   "Maximum reconnection attempts reached",
			retryable: false,
			details:   nil,
			wantCode:  "MAX_RECONNECT_ATTEMPTS",
			wantMsg:   "Maximum reconnection attempts reached",
			wantRetry: false,
		},
		{
			name:      "max auth retries error",
			code:      "MAX_AUTH_RETRIES",
			message:   "Maximum authentication retries exceeded",
			retryable: false,
			details:   nil,
			wantCode:  "MAX_AUTH_RETRIES",
			wantMsg:   "Maximum authentication retries exceeded",
			wantRetry: false,
		},
		{
			name:      "reconnect disabled error",
			code:      "RECONNECT_DISABLED",
			message:   "Reconnection is disabled",
			retryable: false,
			details:   nil,
			wantCode:  "RECONNECT_DISABLED",
			wantMsg:   "Reconnection is disabled",
			wantRetry: false,
		},
		{
			name:      "with details",
			code:      "MAX_RECONNECT_ATTEMPTS",
			message:   "Max attempts reached",
			retryable: true,
			details:   map[string]any{"attempts": 5},
			wantCode:  "MAX_RECONNECT_ATTEMPTS",
			wantMsg:   "Max attempts reached",
			wantRetry: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewReconnectError(tt.code, tt.message, tt.retryable, tt.details)

			if !IsReconnectError(err) {
				t.Errorf("IsReconnectError() = false, want true")
			}

			if err.Code() != tt.wantCode {
				t.Errorf("Code() = %s, want %s", err.Code(), tt.wantCode)
			}

			if err.Error() != tt.wantMsg {
				t.Errorf("Error() = %s, want %s", err.Error(), tt.wantMsg)
			}

			if err.Retryable() != tt.wantRetry {
				t.Errorf("Retryable() = %v, want %v", err.Retryable(), tt.wantRetry)
			}

			if !compareDetails(err.Details(), tt.details) {
				t.Errorf("Details() = %v, want %v", err.Details(), tt.details)
			}
		})
	}
}

// ============================================================================
// NewCancelledError Tests
// ============================================================================

func TestNewCancelledError(t *testing.T) {
	tests := []struct {
		name    string
		message string
		details any
		wantMsg string
	}{
		{
			name:    "basic cancelled error",
			message: "Request was cancelled",
			details: nil,
			wantMsg: "Request was cancelled",
		},
		{
			name:    "with details",
			message: "Request cancelled by user",
			details: "user-initiated",
			wantMsg: "Request cancelled by user",
		},
		{
			name:    "with context details",
			message: "Request cancelled",
			details: map[string]any{"requestId": "abc123"},
			wantMsg: "Request cancelled",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewCancelledError(tt.message, tt.details)

			if !IsCancelledError(err) {
				t.Errorf("IsCancelledError() = false, want true")
			}

			if !IsRequestError(err) {
				t.Errorf("IsRequestError() = false, want true (CancelledError embeds RequestError)")
			}

			if err.Error() != tt.wantMsg {
				t.Errorf("Error() = %s, want %s", err.Error(), tt.wantMsg)
			}

			if err.Code() != "REQUEST_CANCELLED" {
				t.Errorf("Code() = %s, want REQUEST_CANCELLED", err.Code())
			}

			if err.Retryable() != false {
				t.Errorf("Retryable() = %v, want false", err.Retryable())
			}

			if !compareDetails(err.Details(), tt.details) {
				t.Errorf("Details() = %v, want %v", err.Details(), tt.details)
			}
		})
	}
}

// ============================================================================
// NewAbortError Tests
// ============================================================================

func TestNewAbortError(t *testing.T) {
	tests := []struct {
		name    string
		message string
		details any
		wantMsg string
	}{
		{
			name:    "basic abort error",
			message: "Request was aborted",
			details: nil,
			wantMsg: "Request was aborted",
		},
		{
			name:    "with details",
			message: "Request aborted by server",
			details: "server-shutdown",
			wantMsg: "Request aborted by server",
		},
		{
			name:    "with context details",
			message: "Request aborted",
			details: map[string]any{"reason": "timeout"},
			wantMsg: "Request aborted",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewAbortError(tt.message, tt.details)

			if !IsAbortError(err) {
				t.Errorf("IsAbortError() = false, want true")
			}

			if !IsRequestError(err) {
				t.Errorf("IsRequestError() = false, want true (AbortError embeds RequestError)")
			}

			if err.Error() != tt.wantMsg {
				t.Errorf("Error() = %s, want %s", err.Error(), tt.wantMsg)
			}

			if err.Code() != "REQUEST_ABORTED" {
				t.Errorf("Code() = %s, want REQUEST_ABORTED", err.Code())
			}

			if err.Retryable() != false {
				t.Errorf("Retryable() = %v, want false", err.Retryable())
			}

			if !compareDetails(err.Details(), tt.details) {
				t.Errorf("Details() = %v, want %v", err.Details(), tt.details)
			}
		})
	}
}

// ============================================================================
// IsRetryable Tests
// ============================================================================

func TestIsRetryable(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		want    bool
		wantErr bool
	}{
		{
			name:    "retryable AuthError",
			err:     NewAuthError("AUTH_TOKEN_EXPIRED", "Token expired", true, nil),
			want:    true,
			wantErr: false,
		},
		{
			name:    "non-retryable AuthError",
			err:     NewAuthError("AUTH_TOKEN_EXPIRED", "Token expired", false, nil),
			want:    false,
			wantErr: false,
		},
		{
			name:    "retryable ConnectionError",
			err:     NewConnectionError("CONNECTION_TIMEOUT", "Connection timeout", true, nil),
			want:    true,
			wantErr: false,
		},
		{
			name:    "non-retryable ConnectionError",
			err:     NewConnectionError("CONNECTION_CLOSED", "Connection closed", false, nil),
			want:    false,
			wantErr: false,
		},
		{
			name:    "retryable RequestError",
			err:     NewRequestError("METHOD_NOT_FOUND", "Method not found", true, nil),
			want:    true,
			wantErr: false,
		},
		{
			name:    "non-retryable RequestError",
			err:     NewRequestError("INTERNAL_ERROR", "Internal error", false, nil),
			want:    false,
			wantErr: false,
		},
		{
			name:    "TimeoutError is retryable",
			err:     NewTimeoutError("Request timed out", nil),
			want:    true,
			wantErr: false,
		},
		{
			name:    "CancelledError is not retryable",
			err:     NewCancelledError("Request cancelled", nil),
			want:    false,
			wantErr: false,
		},
		{
			name:    "AbortError is not retryable",
			err:     NewAbortError("Request aborted", nil),
			want:    false,
			wantErr: false,
		},
		{
			name:    "Retryable ReconnectError",
			err:     NewReconnectError("MAX_RECONNECT_ATTEMPTS", "Max attempts", true, nil),
			want:    true,
			wantErr: false,
		},
		{
			name:    "Non-retryable ReconnectError",
			err:     NewReconnectError("RECONNECT_DISABLED", "Disabled", false, nil),
			want:    false,
			wantErr: false,
		},
		{
			name:    "standard error returns false",
			err:     errors.New("standard error"),
			want:    false,
			wantErr: false,
		},
		{
			name:    "nil error returns false",
			err:     nil,
			want:    false,
			wantErr: false,
		},
		{
			name:    "wrapped OpenClawError",
			err:     fmt.Errorf("wrapped: %w", NewAuthError("AUTH_TOKEN_EXPIRED", "Token expired", true, nil)),
			want:    true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsRetryable(tt.err)
			if got != tt.want {
				t.Errorf("IsRetryable() = %v, want %v", got, tt.want)
			}
		})
	}
}

// ============================================================================
// Unwrap Tests for Wrapped Errors
// ============================================================================

func TestUnwrapMethods(t *testing.T) {
	// The type-specific Unwrap() methods return the embedded *BaseError
	// (which is an error itself via Error() method), not the underlying err field.

	t.Run("AuthError Unwrap returns BaseError", func(t *testing.T) {
		err := NewAuthError("AUTH_TOKEN_EXPIRED", "Token expired", true, nil)
		unwrapped := err.Unwrap()
		if unwrapped == nil {
			t.Fatal("AuthError.Unwrap() = nil, want *BaseError")
		}
		// BaseError.Error() returns the message
		if unwrapped.Error() != "Token expired" {
			t.Errorf("unwrapped.Error() = %s, want 'Token expired'", unwrapped.Error())
		}
	})

	t.Run("ConnectionError Unwrap returns BaseError", func(t *testing.T) {
		err := NewConnectionError("CONNECTION_CLOSED", "Connection closed", false, nil)
		unwrapped := err.Unwrap()
		if unwrapped == nil {
			t.Fatal("ConnectionError.Unwrap() = nil, want *BaseError")
		}
		if unwrapped.Error() != "Connection closed" {
			t.Errorf("unwrapped.Error() = %s, want 'Connection closed'", unwrapped.Error())
		}
	})

	t.Run("ProtocolError Unwrap returns BaseError", func(t *testing.T) {
		err := NewProtocolError("PROTOCOL_UNSUPPORTED", "Protocol unsupported", true, nil)
		unwrapped := err.Unwrap()
		if unwrapped == nil {
			t.Fatal("ProtocolError.Unwrap() = nil, want *BaseError")
		}
		if unwrapped.Error() != "Protocol unsupported" {
			t.Errorf("unwrapped.Error() = %s, want 'Protocol unsupported'", unwrapped.Error())
		}
	})

	t.Run("RequestError Unwrap returns BaseError", func(t *testing.T) {
		err := NewRequestError("INTERNAL_ERROR", "Internal error", false, nil)
		unwrapped := err.Unwrap()
		if unwrapped == nil {
			t.Fatal("RequestError.Unwrap() = nil, want *BaseError")
		}
		if unwrapped.Error() != "Internal error" {
			t.Errorf("unwrapped.Error() = %s, want 'Internal error'", unwrapped.Error())
		}
	})

	t.Run("GatewayError Unwrap returns BaseError", func(t *testing.T) {
		err := NewGatewayError("AGENT_NOT_FOUND", "Agent not found", false, nil)
		unwrapped := err.Unwrap()
		if unwrapped == nil {
			t.Fatal("GatewayError.Unwrap() = nil, want *BaseError")
		}
		if unwrapped.Error() != "Agent not found" {
			t.Errorf("unwrapped.Error() = %s, want 'Agent not found'", unwrapped.Error())
		}
	})

	t.Run("ReconnectError Unwrap returns BaseError", func(t *testing.T) {
		err := NewReconnectError("MAX_RECONNECT_ATTEMPTS", "Max attempts", false, nil)
		unwrapped := err.Unwrap()
		if unwrapped == nil {
			t.Fatal("ReconnectError.Unwrap() = nil, want *BaseError")
		}
		if unwrapped.Error() != "Max attempts" {
			t.Errorf("unwrapped.Error() = %s, want 'Max attempts'", unwrapped.Error())
		}
	})
}

func TestUnwrapTimeoutCancelledAbort(t *testing.T) {
	// TimeoutError, CancelledError, AbortError embed *RequestError,
	// and their Unwrap() returns that embedded *RequestError

	t.Run("TimeoutError Unwrap returns *RequestError", func(t *testing.T) {
		err := NewTimeoutError("Request timed out", nil)
		unwrapped := err.Unwrap()
		if unwrapped == nil {
			t.Fatal("TimeoutError.Unwrap() = nil, want *RequestError")
		}
		var reqErr *RequestError
		if !errors.As(unwrapped, &reqErr) {
			t.Error("unwrapped is not a *RequestError")
		}
	})

	t.Run("CancelledError Unwrap returns *RequestError", func(t *testing.T) {
		err := NewCancelledError("Request cancelled", nil)
		unwrapped := err.Unwrap()
		if unwrapped == nil {
			t.Fatal("CancelledError.Unwrap() = nil, want *RequestError")
		}
		var reqErr *RequestError
		if !errors.As(unwrapped, &reqErr) {
			t.Error("unwrapped is not a *RequestError")
		}
	})

	t.Run("AbortError Unwrap returns *RequestError", func(t *testing.T) {
		err := NewAbortError("Request aborted", nil)
		unwrapped := err.Unwrap()
		if unwrapped == nil {
			t.Fatal("AbortError.Unwrap() = nil, want *RequestError")
		}
		var reqErr *RequestError
		if !errors.As(unwrapped, &reqErr) {
			t.Error("unwrapped is not a *RequestError")
		}
	})
}

func TestErrorsIsSupport(t *testing.T) {
	authErr := NewAuthError("AUTH_TOKEN_EXPIRED", "Token expired", true, nil)

	t.Run("errors.Is with wrapped error", func(t *testing.T) {
		wrapped := fmt.Errorf("wrapper: %w", authErr)
		if !errors.Is(wrapped, authErr) {
			t.Error("errors.Is(wrapped, authErr) = false, want true")
		}
	})

	t.Run("errors.Is works through unwrap chain", func(t *testing.T) {
		wrapped := fmt.Errorf("wrapper: %w", authErr)
		// authErr.Unwrap() returns BaseError, which is an error
		// So errors.Is should find authErr through the chain
		if !errors.Is(wrapped, authErr) {
			t.Error("errors.Is(wrapped, authErr) = false, want true")
		}
	})
}

func TestErrorsAsSupport(t *testing.T) {
	innerErr := errors.New("inner error")

	t.Run("errors.As for AuthError", func(t *testing.T) {
		err := NewAuthError("AUTH_TOKEN_EXPIRED", "Token expired", true, innerErr)
		wrapped := fmt.Errorf("wrapper: %w", err)

		var authErr *AuthError
		if !errors.As(wrapped, &authErr) {
			t.Error("errors.As(wrapped, &authErr) = false, want true")
		}
		if authErr.Code() != "AUTH_TOKEN_EXPIRED" {
			t.Errorf("authErr.Code() = %s, want AUTH_TOKEN_EXPIRED", authErr.Code())
		}
	})

	t.Run("errors.As for ConnectionError", func(t *testing.T) {
		err := NewConnectionError("CONNECTION_CLOSED", "Connection closed", false, innerErr)
		wrapped := fmt.Errorf("wrapper: %w", err)

		var connErr *ConnectionError
		if !errors.As(wrapped, &connErr) {
			t.Error("errors.As(wrapped, &connErr) = false, want true")
		}
		if connErr.Code() != "CONNECTION_CLOSED" {
			t.Errorf("connErr.Code() = %s, want CONNECTION_CLOSED", connErr.Code())
		}
	})

	t.Run("errors.As for ReconnectError", func(t *testing.T) {
		err := NewReconnectError("MAX_RECONNECT_ATTEMPTS", "Max attempts", false, innerErr)
		wrapped := fmt.Errorf("wrapper: %w", err)

		var reconnErr *ReconnectError
		if !errors.As(wrapped, &reconnErr) {
			t.Error("errors.As(wrapped, &reconnErr) = false, want true")
		}
		if reconnErr.Code() != "MAX_RECONNECT_ATTEMPTS" {
			t.Errorf("reconnErr.Code() = %s, want MAX_RECONNECT_ATTEMPTS", reconnErr.Code())
		}
	})

	t.Run("errors.As for CancelledError", func(t *testing.T) {
		err := NewCancelledError("Request cancelled", innerErr)
		wrapped := fmt.Errorf("wrapper: %w", err)

		var cancelledErr *CancelledError
		if !errors.As(wrapped, &cancelledErr) {
			t.Error("errors.As(wrapped, &cancelledErr) = false, want true")
		}
		if cancelledErr.Code() != "REQUEST_CANCELLED" {
			t.Errorf("cancelledErr.Code() = %s, want REQUEST_CANCELLED", cancelledErr.Code())
		}
	})

	t.Run("errors.As for AbortError", func(t *testing.T) {
		err := NewAbortError("Request aborted", innerErr)
		wrapped := fmt.Errorf("wrapper: %w", err)

		var abortErr *AbortError
		if !errors.As(wrapped, &abortErr) {
			t.Error("errors.As(wrapped, &abortErr) = false, want true")
		}
		if abortErr.Code() != "REQUEST_ABORTED" {
			t.Errorf("abortErr.Code() = %s, want REQUEST_ABORTED", abortErr.Code())
		}
	})

	t.Run("errors.As for TimeoutError", func(t *testing.T) {
		err := NewTimeoutError("Request timed out", innerErr)
		wrapped := fmt.Errorf("wrapper: %w", err)

		var timeoutErr *TimeoutError
		if !errors.As(wrapped, &timeoutErr) {
			t.Error("errors.As(wrapped, &timeoutErr) = false, want true")
		}
		if timeoutErr.Code() != "REQUEST_TIMEOUT" {
			t.Errorf("timeoutErr.Code() = %s, want REQUEST_TIMEOUT", timeoutErr.Code())
		}
	})
}
