# Phase 1: Project Setup and Foundation

**Files:**
- Create: `go.mod`
- Create: `types.go`
- Create: `errors.go`, `errors_test.go`
- Create: `logger.go`

---

## Task 1.1: Initialize Go Module

- [ ] **Step 1: Create go.mod**

```bash
cd /Users/linyang/workspace/my-projects/openclaw-sdk-go
go mod init github.com/i0r3k/openclaw-sdk-go
```

```go
// go.mod
module github.com/i0r3k/openclaw-sdk-go

go 1.21

require github.com/gorilla/websocket v1.5.1
```

- [ ] **Step 2: Add dependencies**

Run: `go mod tidy`

- [ ] **Step 3: Commit**

```bash
git add go.mod go.sum
git commit -m "chore: initialize go module with dependencies"
```

---

## Task 1.2: Create Basic Types

- [ ] **Step 1: Write types.go**

```go
package openclaw

import "time"

// ConnectionState represents the state of the connection
type ConnectionState string

const (
	StateDisconnected      ConnectionState = "disconnected"
	StateConnecting        ConnectionState = "connecting"
	StateConnected         ConnectionState = "connected"
	StateAuthenticating    ConnectionState = "authenticating"
	StateAuthenticated     ConnectionState = "authenticated"
	StateReconnecting      ConnectionState = "reconnecting"
	StateFailed            ConnectionState = "failed"
)

// EventType represents the type of event
type EventType string

const (
	EventConnect      EventType = "connect"
	EventDisconnect   EventType = "disconnect"
	EventError        EventType = "error"
	EventMessage      EventType = "message"
	EventRequest      EventType = "request"
	EventResponse     EventType = "response"
	EventTick         EventType = "tick"
	EventGap          EventType = "gap"
	EventStateChange  EventType = "stateChange"
)

// Event represents a generic event
type Event struct {
	Type      EventType
	Payload   interface{}
	Err       error
	Timestamp time.Time
}

// EventHandler is a function that handles events
type EventHandler func(Event)

// ReconnectConfig holds reconnection settings
type ReconnectConfig struct {
	MaxAttempts       int
	InitialDelay     time.Duration
	MaxDelay         time.Duration
	BackoffMultiplier float64
}

// DefaultReconnectConfig returns sensible defaults
func DefaultReconnectConfig() *ReconnectConfig {
	return &ReconnectConfig{
		MaxAttempts:       0,
		InitialDelay:     1 * time.Second,
		MaxDelay:         60 * time.Second,
		BackoffMultiplier: 1.618,
	}
}
```

- [ ] **Step 2: Verify it compiles**

Run: `go build ./...`

- [ ] **Step 3: Commit**

```bash
git add types.go
git commit -m "feat: add common types and constants"
```

---

## Task 1.3: Create Error Types

- [ ] **Step 1: Write errors.go**

```go
package openclaw

import "fmt"

// ErrorCode represents an error code
type ErrorCode string

const (
	ErrCodeConnection   ErrorCode = "CONNECTION_ERROR"
	ErrCodeAuth        ErrorCode = "AUTH_ERROR"
	ErrCodeTimeout     ErrorCode = "TIMEOUT"
	ErrCodeProtocol    ErrorCode = "PROTOCOL_ERROR"
	ErrCodeValidation  ErrorCode = "VALIDATION_ERROR"
	ErrCodeTransport   ErrorCode = "TRANSPORT_ERROR"
	ErrCodeUnknown     ErrorCode = "UNKNOWN"
)

// OpenClawError is the base error interface
type OpenClawError interface {
	error
	Code() ErrorCode
	Unwrap() error
}

// BaseError is the base error struct
type BaseError struct {
	code    ErrorCode
	message string
	err     error
}

func (e *BaseError) Error() string { return e.message }
func (e *BaseError) Code() ErrorCode { return e.code }
func (e *BaseError) Unwrap() error { return e.err }

// ConnectionError represents a connection error
type ConnectionError struct {
	*BaseError
}

// AuthError represents an authentication error
type AuthError struct {
	*BaseError
}

// TimeoutError represents a timeout error
type TimeoutError struct {
	*BaseError
}

// ProtocolError represents a protocol error
type ProtocolError struct {
	*BaseError
}

// ValidationError represents a validation error
type ValidationError struct {
	*BaseError
}

// TransportError represents a transport error
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
func Is(err error, code ErrorCode) bool {
	var e OpenClawError
	if As(err, &e) {
		return e.Code() == code
	}
	return false
}

// As casts the error to OpenClawError
func As(err error, target interface{}) bool {
	if target == nil {
		return false
	}
	if e, ok := err.(OpenClawError); ok {
		if t, ok := target.(**OpenClawError); ok {
			*t = &e
			return true
		}
	}
	return false
}
```

- [ ] **Step 2: Verify it compiles**

Run: `go build ./...`

- [ ] **Step 3: Write basic test**

```go
// errors_test.go
package openclaw

import (
	"testing"
)

func TestNewError(t *testing.T) {
	err := NewError(ErrCodeConnection, "test error", nil)
	if err.Error() != "test error" {
		t.Errorf("expected 'test error', got '%s'", err.Error())
	}
	if err.Code() != ErrCodeConnection {
		t.Errorf("expected CONNECTION_ERROR, got %s", err.Code())
	}
}

func TestIs(t *testing.T) {
	err := NewConnectionError("connection failed", nil)
	if !Is(err, ErrCodeConnection) {
		t.Error("expected Is to return true for matching code")
	}
	if Is(err, ErrCodeAuth) {
		t.Error("expected Is to return false for non-matching code")
	}
}
```

- [ ] **Step 4: Run tests**

Run: `go test -v ./...`

- [ ] **Step 5: Commit**

```bash
git add errors.go errors_test.go
git commit -m "feat: add error type hierarchy"
```

---

## Task 1.4: Create Logger Interface

- [ ] **Step 1: Write logger.go**

```go
package openclaw

import "log"

// Logger interface for customizable logging
type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

// DefaultLogger uses stdlib log
type DefaultLogger struct{}

func (l *DefaultLogger) Debug(msg string, args ...any) { log.Printf("[DEBUG] "+msg, args...) }
func (l *DefaultLogger) Info(msg string, args ...any)  { log.Printf("[INFO] "+msg, args...) }
func (l *DefaultLogger) Warn(msg string, args ...any)  { log.Printf("[WARN] "+msg, args...) }
func (l *DefaultLogger) Error(msg string, args ...any) { log.Printf("[ERROR] "+msg, args...) }
```

- [ ] **Step 2: Commit**

```bash
git add logger.go
git commit -m "feat: add Logger interface"
```

---

## Phase 1 Complete

After this phase, you should have:
- `go.mod` - Go module initialized
- `types.go` - Common types and constants
- `errors.go` - Error type hierarchy
- `logger.go` - Logger interface

All code should compile and tests should pass.
