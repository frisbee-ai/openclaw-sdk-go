package utils

import (
	"context"
	"testing"
	"time"
)

func TestNewTimeoutManager_Negative(t *testing.T) {
	tm := NewTimeoutManager(-1 * time.Second)
	if tm.defaultTimeout != 0 {
		t.Errorf("expected 0, got %v", tm.defaultTimeout)
	}
}

func TestNewTimeoutManager_Zero(t *testing.T) {
	tm := NewTimeoutManager(0)
	if tm.defaultTimeout != 0 {
		t.Errorf("expected 0, got %v", tm.defaultTimeout)
	}
}

func TestWithTimeout_Default(t *testing.T) {
	tm := NewTimeoutManager(10 * time.Second)
	ctx, cancel := tm.WithTimeout(context.Background())
	defer cancel()

	if ctx == nil {
		t.Error("expected non-nil context")
	}
}

func TestWithTimeout_ZeroDefault(t *testing.T) {
	tm := NewTimeoutManager(0) // Zero timeout
	ctx, cancel := tm.WithTimeout(context.Background())
	defer cancel()

	if ctx == nil {
		t.Error("expected non-nil context")
	}
}

func TestWithCustomTimeout_Positive(t *testing.T) {
	tm := NewTimeoutManager(10 * time.Second)
	ctx, cancel := tm.WithCustomTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if ctx == nil {
		t.Error("expected non-nil context")
	}
}

func TestWithCustomTimeout_Negative(t *testing.T) {
	tm := NewTimeoutManager(10 * time.Second)
	ctx, cancel := tm.WithCustomTimeout(context.Background(), -5*time.Second)
	defer cancel()

	// Negative timeout should behave like no timeout (WithCancel)
	if ctx == nil {
		t.Error("expected non-nil context")
	}
}

func TestWithCustomTimeout_Zero(t *testing.T) {
	tm := NewTimeoutManager(10 * time.Second)
	ctx, cancel := tm.WithCustomTimeout(context.Background(), 0)
	defer cancel()

	// Zero timeout should behave like no timeout
	if ctx == nil {
		t.Error("expected non-nil context")
	}
}
