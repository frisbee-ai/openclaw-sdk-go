// Package protocol provides protocol frame types and utilities for OpenClaw SDK.
//
// This package provides validation for protocol frames:
//   - Validator: Validates GatewayFrame, RequestFrame, ResponseFrame, EventFrame
//   - ValidationError: Structured validation errors with field and message
package protocol

import (
	"errors"
	"strings"
)

// ValidationError represents a validation error with field name and message.
// Used by Validator to provide structured error information.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Message
}

// Validator validates protocol frames.
// Provides methods to validate each frame type according to protocol rules.
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
	if !frame.Type.IsValid() {
		return &ValidationError{Field: "Type", Message: "is not a valid frame type"}
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
	// Validate method format (namespace.method or namespace.sub.method)
	parts := strings.Split(frame.Method, ".")
	if len(parts) < 2 {
		return &ValidationError{Field: "Method", Message: "must be in format 'namespace.method'"}
	}
	for _, part := range parts {
		if part == "" {
			return &ValidationError{Field: "Method", Message: "must be in format 'namespace.method'"}
		}
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
	// Success and Error are mutually exclusive
	if frame.Success && frame.Error != nil {
		return &ValidationError{Field: "Error", Message: "must be nil when Success is true"}
	}
	if !frame.Success && frame.Error == nil {
		return &ValidationError{Field: "Error", Message: "is required when Success is false"}
	}
	return nil
}

// ValidateEventFrame validates an event frame
func (v *Validator) ValidateEventFrame(frame *EventFrame) error {
	if frame == nil {
		return errors.New("frame is nil")
	}
	if frame.EventType == "" {
		return &ValidationError{Field: "EventType", Message: "is required"}
	}
	return nil
}
