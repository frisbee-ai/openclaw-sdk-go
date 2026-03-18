# Phase 3: Protocol Module

**Files:**
- Create: `protocol/types.go`, `protocol/types_test.go`
- Create: `protocol/validation.go`, `protocol/validation_test.go`

**Depends on:** Phase 1 (types.go, errors.go)

---

## Task 3.1: Protocol Types

- [ ] **Step 1: Create protocol directory and types.go**

```bash
mkdir -p protocol
```

```go
// protocol/types.go
package protocol

import (
	"encoding/json"
	"time"
)

// FrameType represents the type of frame
type FrameType string

const (
	FrameTypeGateway   FrameType = "gateway"
	FrameTypeRequest   FrameType = "request"
	FrameTypeResponse  FrameType = "response"
	FrameTypeEvent     FrameType = "event"
	FrameTypeError     FrameType = "error"
)

// GatewayFrame is the main frame type
type GatewayFrame struct {
	Type      FrameType          `json:"type"`
	Timestamp time.Time          `json:"timestamp"`
	Payload   json.RawMessage    `json:"payload,omitempty"`
}

// RequestFrame represents a request frame
type RequestFrame struct {
	RequestID  string          `json:"requestId"`
	Method     string          `json:"method"`
	Params     json.RawMessage `json:"params,omitempty"`
	Timestamp  time.Time      `json:"timestamp"`
}

// ResponseFrame represents a response frame
type ResponseFrame struct {
	RequestID string          `json:"requestId"`
	Success   bool            `json:"success"`
	Result    json.RawMessage `json:"result,omitempty"`
	Error     *ResponseError `json:"error,omitempty"`
	Timestamp time.Time      `json:"timestamp"`
}

// ResponseError represents an error in a response
type ResponseError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// EventFrame represents an event frame
type EventFrame struct {
	EventType string          `json:"eventType"`
	Data      json.RawMessage `json:"data,omitempty"`
	Timestamp time.Time       `json:"timestamp"`
}
```

- [ ] **Step 2: Write test**

```go
// protocol/types_test.go
package protocol

import (
	"encoding/json"
	"testing"
	"time"
)

func TestGatewayFrameSerialization(t *testing.T) {
	frame := GatewayFrame{
		Type:      FrameTypeGateway,
		Timestamp: time.Now(),
		Payload:   json.RawMessage(`{"key":"value"}`),
	}

	data, err := json.Marshal(frame)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var decoded GatewayFrame
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if decoded.Type != frame.Type {
		t.Errorf("expected %s, got %s", frame.Type, decoded.Type)
	}
}
```

- [ ] **Step 3: Run tests and commit**

Run: `go test -v ./protocol/...`
Commit: `git add protocol/ && git commit -m "feat: add protocol types"`

---

## Task 3.2: Protocol Validation

- [ ] **Step 1: Write validation.go**

```go
// protocol/validation.go
package protocol

import (
	"errors"
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Message
}

// Validator validates protocol frames
type Validator struct{}

// NewValidator creates a new validator
func NewValidator() *Validator {
	return &Validator{}
}

// ValidateGatewayFrame validates a gateway frame
func (v *Validator) ValidateGatewayFrame(frame *GatewayFrame) error {
	if frame == nil {
		return errors.New("frame is nil")
	}
	if frame.Type == "" {
		return &ValidationError{Field: "Type", Message: "is required"}
	}
	return nil
}

// ValidateRequestFrame validates a request frame
func (v *Validator) ValidateRequestFrame(frame *RequestFrame) error {
	if frame == nil {
		return errors.New("frame is nil")
	}
	if frame.RequestID == "" {
		return &ValidationError{Field: "RequestID", Message: "is required"}
	}
	if frame.Method == "" {
		return &ValidationError{Field: "Method", Message: "is required"}
	}
	return nil
}

// ValidateResponseFrame validates a response frame
func (v *Validator) ValidateResponseFrame(frame *ResponseFrame) error {
	if frame == nil {
		return errors.New("frame is nil")
	}
	if frame.RequestID == "" {
		return &ValidationError{Field: "RequestID", Message: "is required"}
	}
	return nil
}
```

- [ ] **Step 2: Write test**

```go
// protocol/validation_test.go
package protocol

import (
	"testing"
)

func TestValidator_ValidateGatewayFrame(t *testing.T) {
	v := NewValidator()
	err := v.ValidateGatewayFrame(nil)
	if err == nil {
		t.Error("expected error for nil frame")
	}

	err = v.ValidateGatewayFrame(&GatewayFrame{Type: FrameTypeGateway})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	err = v.ValidateGatewayFrame(&GatewayFrame{})
	if err == nil {
		t.Error("expected error for empty type")
	}
}

func TestValidator_ValidateRequestFrame(t *testing.T) {
	v := NewValidator()
	err := v.ValidateRequestFrame(&RequestFrame{
		RequestID: "123",
		Method:   "test",
	})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
```

- [ ] **Step 3: Run tests and commit**

Run: `go test -v ./protocol/...`
Commit: `git add protocol/ && git commit -m "feat: add protocol validation"`

---

## Phase 3 Complete

After this phase, you should have:
- `protocol/types.go` - Protocol frame types
- `protocol/validation.go` - Frame validation

All code should compile and tests should pass.
