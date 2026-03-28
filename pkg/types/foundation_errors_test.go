package types

import (
	"errors"
	"testing"
)

// TestErrTooManyPendingRequests_SentinelIs verifies sentinel error works with errors.Is.
func TestErrTooManyPendingRequests_SentinelIs(t *testing.T) {
	if !errors.Is(ErrTooManyPendingRequests, ErrTooManyPendingRequests) {
		t.Error("errors.Is(ErrTooManyPendingRequests, ErrTooManyPendingRequests) = false, want true")
	}
}

// TestErrTooManyPendingRequests_TypedIs verifies typed error wraps sentinel for errors.Is.
func TestErrTooManyPendingRequests_TypedIs(t *testing.T) {
	err := NewTooManyPendingRequestsError(256)
	if !errors.Is(err, ErrTooManyPendingRequests) {
		t.Error("errors.Is(NewTooManyPendingRequestsError(256), ErrTooManyPendingRequests) = false, want true")
	}
}

// TestErrTooManyPendingRequests_OpenClawError verifies typed error implements OpenClawError.
func TestErrTooManyPendingRequests_OpenClawError(t *testing.T) {
	err := NewTooManyPendingRequestsError(256)

	var oe OpenClawError = err
	if oe == nil {
		t.Fatal("TooManyPendingRequestsError does not implement OpenClawError")
	}

	// Verify all interface methods are accessible
	_ = oe.Code()
	_ = oe.Retryable()
	_ = oe.Details()
	_ = oe.Unwrap()
	_ = oe.Error()
}

// TestErrTooManyPendingRequests_Retryable verifies TooManyPendingRequestsError is retryable.
func TestErrTooManyPendingRequests_Retryable(t *testing.T) {
	err := NewTooManyPendingRequestsError(256)
	if !err.Retryable() {
		t.Error("TooManyPendingRequestsError.Retryable() = false, want true (transient condition)")
	}
}

// TestErrTooManyPendingRequests_Code verifies error code.
func TestErrTooManyPendingRequests_Code(t *testing.T) {
	err := NewTooManyPendingRequestsError(256)
	if err.Code() != "TOO_MANY_PENDING_REQUESTS" {
		t.Errorf("Code() = %s, want TOO_MANY_PENDING_REQUESTS", err.Code())
	}
}

// TestErrTooManyPendingRequests_Details verifies details contains limit.
func TestErrTooManyPendingRequests_Details(t *testing.T) {
	err := NewTooManyPendingRequestsError(256)
	details := err.Details()
	if details == nil {
		t.Fatal("Details() = nil, want map with limit")
	}
	m, ok := details.(map[string]int)
	if !ok {
		t.Fatalf("Details() type = %T, want map[string]int", details)
	}
	if m["limit"] != 256 {
		t.Errorf("Details()[\"limit\"] = %d, want 256", m["limit"])
	}
}

// TestErrTooManyPendingRequests_ErrorMessage verifies error message format.
func TestErrTooManyPendingRequests_ErrorMessage(t *testing.T) {
	err := NewTooManyPendingRequestsError(256)
	if err.Error() != "pending request limit (256) exceeded" {
		t.Errorf("Error() = %s, want 'pending request limit (256) exceeded'", err.Error())
	}
}

// TestErrMaxRetriesExceeded_SentinelIs verifies sentinel error works with errors.Is.
func TestErrMaxRetriesExceeded_SentinelIs(t *testing.T) {
	if !errors.Is(ErrMaxRetriesExceeded, ErrMaxRetriesExceeded) {
		t.Error("errors.Is(ErrMaxRetriesExceeded, ErrMaxRetriesExceeded) = false, want true")
	}
}

// TestErrMaxRetriesExceeded_TypedIs verifies typed error wraps sentinel for errors.Is.
func TestErrMaxRetriesExceeded_TypedIs(t *testing.T) {
	err := NewMaxRetriesExceededError(10)
	if !errors.Is(err, ErrMaxRetriesExceeded) {
		t.Error("errors.Is(NewMaxRetriesExceededError(10), ErrMaxRetriesExceeded) = false, want true")
	}
}

// TestErrMaxRetriesExceeded_OpenClawError verifies typed error implements OpenClawError.
func TestErrMaxRetriesExceeded_OpenClawError(t *testing.T) {
	err := NewMaxRetriesExceededError(10)

	var oe OpenClawError = err
	if oe == nil {
		t.Fatal("MaxRetriesExceededError does not implement OpenClawError")
	}

	// Verify all interface methods are accessible
	_ = oe.Code()
	_ = oe.Retryable()
	_ = oe.Details()
	_ = oe.Unwrap()
	_ = oe.Error()
}

// TestErrMaxRetriesExceeded_NotRetryable verifies MaxRetriesExceededError is NOT retryable.
func TestErrMaxRetriesExceeded_NotRetryable(t *testing.T) {
	err := NewMaxRetriesExceededError(10)
	if err.Retryable() {
		t.Error("MaxRetriesExceededError.Retryable() = true, want false (terminal condition)")
	}
}

// TestErrMaxRetriesExceeded_Code verifies error code.
func TestErrMaxRetriesExceeded_Code(t *testing.T) {
	err := NewMaxRetriesExceededError(10)
	if err.Code() != "MAX_RETRIES_EXCEEDED" {
		t.Errorf("Code() = %s, want MAX_RETRIES_EXCEEDED", err.Code())
	}
}

// TestErrMaxRetriesExceeded_Details verifies details contains max_retries.
func TestErrMaxRetriesExceeded_Details(t *testing.T) {
	err := NewMaxRetriesExceededError(10)
	details := err.Details()
	if details == nil {
		t.Fatal("Details() = nil, want map with max_retries")
	}
	m, ok := details.(map[string]int)
	if !ok {
		t.Fatalf("Details() type = %T, want map[string]int", details)
	}
	if m["max_retries"] != 10 {
		t.Errorf("Details()[\"max_retries\"] = %d, want 10", m["max_retries"])
	}
}

// TestErrMaxRetriesExceeded_ErrorMessage verifies error message format.
func TestErrMaxRetriesExceeded_ErrorMessage(t *testing.T) {
	err := NewMaxRetriesExceededError(10)
	if err.Error() != "max retries exceeded: 10 attempts" {
		t.Errorf("Error() = %s, want 'max retries exceeded: 10 attempts'", err.Error())
	}
}
