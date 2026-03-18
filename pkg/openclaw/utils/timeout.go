package utils

import (
	"context"
	"time"
)

// TimeoutManager manages timeouts
type TimeoutManager struct {
	defaultTimeout time.Duration
}

// NewTimeoutManager creates a new timeout manager
// If defaultTimeout is negative, it will be set to 0 (no timeout)
func NewTimeoutManager(defaultTimeout time.Duration) *TimeoutManager {
	if defaultTimeout < 0 {
		defaultTimeout = 0
	}
	return &TimeoutManager{defaultTimeout: defaultTimeout}
}

// WithTimeout wraps a context with timeout
func (tm *TimeoutManager) WithTimeout(parent context.Context) (context.Context, context.CancelFunc) {
	if tm.defaultTimeout > 0 {
		return context.WithTimeout(parent, tm.defaultTimeout)
	}
	return context.WithCancel(parent)
}

// WithCustomTimeout wraps a context with a custom timeout
// If timeout is zero or negative, behaves like WithCancel (no timeout)
func (tm *TimeoutManager) WithCustomTimeout(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if timeout <= 0 {
		return context.WithCancel(parent)
	}
	return context.WithTimeout(parent, timeout)
}

// DefaultTimeoutManager is a global default
var DefaultTimeoutManager = NewTimeoutManager(30 * time.Second)
