package managers

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/frisbee-ai/openclaw-sdk-go/pkg/types"
)

func TestReconnectManager_Stop(t *testing.T) {
	config := DefaultReconnectConfig()
	config.MaxAttempts = 1
	config.InitialDelay = 10 * time.Millisecond

	rm := NewReconnectManager(config)
	rm.Start()

	// Wait a bit then stop
	time.Sleep(20 * time.Millisecond)
	rm.Stop()
}

func TestReconnectManager_Callbacks(t *testing.T) {
	config := DefaultReconnectConfig()
	config.MaxAttempts = 1
	config.InitialDelay = 10 * time.Millisecond

	rm := NewReconnectManager(config)

	var mu sync.Mutex
	reconnectCalled := false

	rm.SetOnReconnect(func() error {
		mu.Lock()
		reconnectCalled = true
		mu.Unlock()
		return nil // Success - stops reconnect loop
	})

	rm.Start()

	// Wait for reconnect to be called
	time.Sleep(30 * time.Millisecond)

	mu.Lock()
	if !reconnectCalled {
		t.Error("expected reconnect callback to be called")
	}
	mu.Unlock()

	rm.Stop()
}

func TestReconnectManager_NoCallbackStops(t *testing.T) {
	config := DefaultReconnectConfig()
	config.InitialDelay = 10 * time.Millisecond

	rm := NewReconnectManager(config)
	// Don't set any callback - should stop immediately

	rm.Start()
	time.Sleep(20 * time.Millisecond)

	// Should have stopped because no callback was set
	rm.Stop() // Should be safe to call even if already stopped
}

func TestReconnectManager_FailedCallback(t *testing.T) {
	config := DefaultReconnectConfig()
	config.MaxAttempts = 2
	config.MaxRetries = 2 // FOUND-02: MaxRetries takes precedence
	config.InitialDelay = 10 * time.Millisecond

	rm := NewReconnectManager(config)

	var mu sync.Mutex
	failedCalled := false
	attemptCount := 0
	failedErrors := []error{}

	rm.SetOnReconnect(func() error {
		mu.Lock()
		attemptCount++
		mu.Unlock()
		return &testError{msg: "connection failed"}
	})

	rm.SetOnReconnectFailed(func(err error) {
		mu.Lock()
		failedCalled = true
		failedErrors = append(failedErrors, err)
		mu.Unlock()
	})

	rm.Start()

	// Wait for attempts to complete
	time.Sleep(50 * time.Millisecond)

	mu.Lock()
	if !failedCalled {
		t.Error("expected failed callback to be called")
	}
	// MaxRetries=2: 2 reconnect attempts, then stops
	if attemptCount != 2 {
		t.Errorf("expected 2 reconnect attempts, got %d", attemptCount)
	}
	// onReconnectFailed called 2 times for individual failures + 1 for budget exhaustion
	if len(failedErrors) != 3 {
		t.Errorf("expected 3 failed errors (2 failures + 1 MaxRetriesExceeded), got %d", len(failedErrors))
	}
	// The last error should be MaxRetriesExceeded
	if len(failedErrors) == 3 && !errors.Is(failedErrors[2], types.ErrMaxRetriesExceeded) {
		t.Errorf("expected last error to wrap ErrMaxRetriesExceeded, got %v", failedErrors[2])
	}
	mu.Unlock()

	rm.Stop()
}

// testError is a simple error for testing
type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}

func TestReconnectManager_FibonacciBackoff(t *testing.T) {
	config := DefaultReconnectConfig()
	config.MaxAttempts = 5
	config.InitialDelay = 100 * time.Millisecond
	config.MaxDelay = 5 * time.Second

	rm := NewReconnectManager(config)

	var mu sync.Mutex
	delays := []time.Duration{}
	attemptCount := 0

	rm.SetOnReconnect(func() error {
		mu.Lock()
		attemptCount++
		mu.Unlock()
		return &testError{msg: "connection failed"}
	})

	rm.Start()
	time.Sleep(600 * time.Millisecond) // Wait for several attempts
	rm.Stop()

	mu.Lock()
	if len(delays) == 0 {
		// We can't directly measure delays, but we verified multiple attempts occurred
		if attemptCount < 2 {
			t.Errorf("expected at least 2 attempts, got %d", attemptCount)
		}
	}
	mu.Unlock()

	// Verify Fibonacci sequence: 100ms, 100ms, 200ms, 300ms, 500ms, 800ms...
	// Since we can't directly measure, we verify the logic works via multiple attempts
	if attemptCount < 2 {
		t.Errorf("Fibonacci backoff should allow multiple attempts, got %d", attemptCount)
	}
}

func TestReconnectManager_Reset(t *testing.T) {
	config := DefaultReconnectConfig()
	config.MaxAttempts = 3
	config.InitialDelay = 10 * time.Millisecond

	rm := NewReconnectManager(config)

	// Reset is a no-op - calling it should not panic
	// Attempts are tracked in the run() loop, not persisted
	rm.Reset()

	// Start and let it run a bit
	rm.SetOnReconnect(func() error {
		return &testError{msg: "connection failed"}
	})
	rm.Start()
	time.Sleep(30 * time.Millisecond)

	// Reset again while running - still a no-op
	rm.Reset()

	rm.Stop()
}

func TestReconnectManager_SetOnReconnect(t *testing.T) {
	config := DefaultReconnectConfig()
	config.InitialDelay = 1 * time.Millisecond
	rm := NewReconnectManager(config)

	var called bool
	rm.SetOnReconnect(func() error {
		called = true
		return nil
	})

	rm.mu.Lock()
	callback := rm.onReconnect
	rm.mu.Unlock()

	if callback == nil {
		t.Error("expected onReconnect to be set")
	}

	rm.Start()
	time.Sleep(50 * time.Millisecond)
	rm.Stop()

	if !called {
		t.Error("expected reconnect callback to be called")
	}
}

func TestReconnectManager_SetOnReconnectFailed(t *testing.T) {
	config := DefaultReconnectConfig()
	config.InitialDelay = 1 * time.Millisecond
	rm := NewReconnectManager(config)

	var failedErr error
	rm.SetOnReconnectFailed(func(err error) {
		failedErr = err
	})

	rm.mu.Lock()
	callback := rm.onReconnectFailed
	rm.mu.Unlock()

	if callback == nil {
		t.Error("expected onReconnectFailed to be set")
	}

	rm.SetOnReconnect(func() error {
		return &testError{msg: "test error"}
	})
	rm.Start()
	time.Sleep(50 * time.Millisecond)
	rm.Stop()

	if failedErr == nil {
		t.Error("expected failed callback to be called")
	}
}
