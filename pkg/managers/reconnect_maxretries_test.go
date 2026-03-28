package managers

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/frisbee-ai/openclaw-sdk-go/pkg/types"
)

// TestReconnectManager_MaxRetries_StopsAtLimit verifies MaxRetries enforcement (FOUND-02)
func TestReconnectManager_MaxRetries_StopsAtLimit(t *testing.T) {
	config := DefaultReconnectConfig()
	config.MaxAttempts = 10 // Should be ignored
	config.MaxRetries = 3   // Should take precedence
	config.InitialDelay = 5 * time.Millisecond

	rm := NewReconnectManager(config)

	var mu sync.Mutex
	attemptCount := 0

	rm.SetOnReconnect(func() error {
		mu.Lock()
		attemptCount++
		mu.Unlock()
		return &testError{msg: "connection failed"}
	})

	rm.SetOnReconnectFailed(func(err error) {
		// Will be called on each failure
	})

	rm.Start()
	time.Sleep(80 * time.Millisecond) // Wait for all attempts
	rm.Stop()

	mu.Lock()
	if attemptCount != 3 {
		t.Errorf("expected 3 attempts, got %d", attemptCount)
	}
	mu.Unlock()
}

// TestReconnectManager_MaxRetries_FallbackToMaxAttempts verifies backward compatibility (FOUND-02)
func TestReconnectManager_MaxRetries_FallbackToMaxAttempts(t *testing.T) {
	config := DefaultReconnectConfig()
	config.MaxAttempts = 4
	config.MaxRetries = 0 // Should fall back to MaxAttempts
	config.InitialDelay = 5 * time.Millisecond

	rm := NewReconnectManager(config)

	var mu sync.Mutex
	attemptCount := 0

	rm.SetOnReconnect(func() error {
		mu.Lock()
		attemptCount++
		mu.Unlock()
		return &testError{msg: "connection failed"}
	})

	rm.SetOnReconnectFailed(func(err error) {
		// Will be called on each failure
	})

	rm.Start()
	time.Sleep(100 * time.Millisecond) // Wait for all attempts
	rm.Stop()

	mu.Lock()
	if attemptCount != 4 {
		t.Errorf("expected 4 attempts (MaxAttempts fallback), got %d", attemptCount)
	}
	mu.Unlock()
}

// TestReconnectManager_MaxRetries_TakesPrecedence verifies MaxRetries > MaxAttempts (FOUND-02)
func TestReconnectManager_MaxRetries_TakesPrecedence(t *testing.T) {
	config := DefaultReconnectConfig()
	config.MaxAttempts = 10 // Should be ignored
	config.MaxRetries = 2   // Should take precedence
	config.InitialDelay = 5 * time.Millisecond

	rm := NewReconnectManager(config)

	var mu sync.Mutex
	attemptCount := 0

	rm.SetOnReconnect(func() error {
		mu.Lock()
		attemptCount++
		mu.Unlock()
		return &testError{msg: "connection failed"}
	})

	rm.Start()
	time.Sleep(60 * time.Millisecond) // Wait for all attempts
	rm.Stop()

	mu.Lock()
	if attemptCount != 2 {
		t.Errorf("expected 2 attempts (MaxRetries precedence), got %d", attemptCount)
	}
	mu.Unlock()
}

// TestReconnectManager_MaxRetries_NegativeTreatedAsZero verifies negative handling (FOUND-02)
func TestReconnectManager_MaxRetries_NegativeTreatedAsZero(t *testing.T) {
	config := DefaultReconnectConfig()
	config.MaxAttempts = 2
	config.MaxRetries = -5 // Should be treated as 0, fall back to MaxAttempts
	config.InitialDelay = 5 * time.Millisecond

	rm := NewReconnectManager(config)

	var mu sync.Mutex
	attemptCount := 0

	rm.SetOnReconnect(func() error {
		mu.Lock()
		attemptCount++
		mu.Unlock()
		return &testError{msg: "connection failed"}
	})

	rm.Start()
	time.Sleep(60 * time.Millisecond) // Wait for all attempts
	rm.Stop()

	mu.Lock()
	if attemptCount != 2 {
		t.Errorf("expected 2 attempts (MaxRetries=-5 treated as 0, falls back to MaxAttempts=2), got %d", attemptCount)
	}
	mu.Unlock()
}

// TestReconnectManager_MaxRetries_CallsFailedWithTypedError verifies typed error (FOUND-02)
func TestReconnectManager_MaxRetries_CallsFailedWithTypedError(t *testing.T) {
	config := DefaultReconnectConfig()
	config.MaxAttempts = 5 // Should be ignored
	config.MaxRetries = 2  // Takes precedence
	config.InitialDelay = 5 * time.Millisecond

	rm := NewReconnectManager(config)

	var mu sync.Mutex
	var lastErr error

	rm.SetOnReconnect(func() error {
		return &testError{msg: "connection failed"}
	})

	rm.SetOnReconnectFailed(func(err error) {
		mu.Lock()
		lastErr = err
		mu.Unlock()
	})

	rm.Start()
	time.Sleep(60 * time.Millisecond) // Wait for all attempts
	rm.Stop()

	mu.Lock()
	defer mu.Unlock()
	if lastErr == nil {
		t.Fatal("expected onReconnectFailed to be called")
	}
	if !errors.Is(lastErr, types.ErrMaxRetriesExceeded) {
		t.Errorf("expected error wrapping ErrMaxRetriesExceeded, got %v", lastErr)
	}
}

// TestReconnectManager_BothZeroInfinite verifies both zero means infinite (FOUND-02)
func TestReconnectManager_BothZeroInfinite(t *testing.T) {
	config := DefaultReconnectConfig()
	config.MaxAttempts = 0
	config.MaxRetries = 0
	config.InitialDelay = 5 * time.Millisecond

	rm := NewReconnectManager(config)

	var mu sync.Mutex
	attemptCount := 0

	rm.SetOnReconnect(func() error {
		mu.Lock()
		attemptCount++
		mu.Unlock()
		return &testError{msg: "connection failed"}
	})

	rm.Start()
	time.Sleep(80 * time.Millisecond) // Run for limited time (enough for 3+ attempts with 5ms delay)

	// Verify multiple attempts occurred (infinite means it kept going)
	mu.Lock()
	if attemptCount < 3 {
		t.Errorf("expected at least 3 attempts (infinite), got %d", attemptCount)
	}
	mu.Unlock()

	rm.Stop()
}

// TestReconnectManager_MaxRetries_ImmediateSuccessDoesNotCallFailed verifies success case (FOUND-02)
func TestReconnectManager_MaxRetries_ImmediateSuccessDoesNotCallFailed(t *testing.T) {
	config := DefaultReconnectConfig()
	config.MaxAttempts = 5
	config.MaxRetries = 5
	config.InitialDelay = 5 * time.Millisecond

	rm := NewReconnectManager(config)

	var mu sync.Mutex
	failedCount := 0

	rm.SetOnReconnect(func() error {
		return nil // Immediate success - no failures
	})

	rm.SetOnReconnectFailed(func(err error) {
		mu.Lock()
		failedCount++
		mu.Unlock()
	})

	rm.Start()
	time.Sleep(30 * time.Millisecond)
	rm.Stop()

	mu.Lock()
	if failedCount != 0 {
		t.Errorf("expected 0 failed callbacks (immediate success), got %d", failedCount)
	}
	mu.Unlock()
}

func TestReconnectManager_AttemptCount_StartsAtZero(t *testing.T) {
	config := DefaultReconnectConfig()
	config.MaxRetries = 5
	config.InitialDelay = 5 * time.Millisecond
	rm := NewReconnectManager(config)
	if rm.AttemptCount() != 0 {
		t.Errorf("expected 0 before Start, got %d", rm.AttemptCount())
	}
	rm.Stop()
}

func TestReconnectManager_AttemptCount_IncrementsOnAttempts(t *testing.T) {
	config := DefaultReconnectConfig()
	config.MaxRetries = 3
	config.InitialDelay = 5 * time.Millisecond
	rm := NewReconnectManager(config)

	rm.SetOnReconnect(func() error {
		return &testError{msg: "connection failed"}
	})
	rm.SetOnReconnectFailed(func(err error) {})

	rm.Start()
	time.Sleep(50 * time.Millisecond) // Wait for attempts
	rm.Stop()

	if rm.AttemptCount() == 0 {
		t.Error("expected AttemptCount > 0 after failed attempts")
	}
}

func TestReconnectManager_AttemptCount_ThreadSafe(t *testing.T) {
	config := DefaultReconnectConfig()
	config.MaxRetries = 100
	config.InitialDelay = 1 * time.Millisecond
	rm := NewReconnectManager(config)

	rm.SetOnReconnect(func() error {
		time.Sleep(1 * time.Millisecond)
		return &testError{msg: "connection failed"}
	})
	rm.SetOnReconnectFailed(func(err error) {})

	rm.Start()
	time.Sleep(50 * time.Millisecond)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = rm.AttemptCount()
		}()
	}
	wg.Wait()
	rm.Stop()
}

func TestReconnectManager_AttemptCount_NoResetOnSuccess(t *testing.T) {
	config := DefaultReconnectConfig()
	config.MaxRetries = 5
	config.InitialDelay = 5 * time.Millisecond
	rm := NewReconnectManager(config)

	var callCount int
	rm.SetOnReconnect(func() error {
		callCount++
		if callCount == 2 {
			return nil // Success on 2nd attempt
		}
		return &testError{msg: "connection failed"}
	})
	rm.SetOnReconnectFailed(func(err error) {})

	rm.Start()
	time.Sleep(50 * time.Millisecond)
	rm.Stop()

	// After 1 failed + 1 success = 2 total attempts (goroutine exits on success)
	if rm.AttemptCount() != 2 {
		t.Errorf("expected 2 attempts, got %d", rm.AttemptCount())
	}
}

func TestReconnectManager_MaxRetries_ZeroDelay(t *testing.T) {
	config := DefaultReconnectConfig()
	config.MaxAttempts = 3
	config.MaxRetries = 0   // Falls back to MaxAttempts
	config.InitialDelay = 0 // Zero delay
	config.MaxDelay = 0

	rm := NewReconnectManager(config)

	var mu sync.Mutex
	attemptCount := 0

	rm.SetOnReconnect(func() error {
		mu.Lock()
		attemptCount++
		mu.Unlock()
		return &testError{msg: "connection failed"}
	})

	rm.Start()
	time.Sleep(30 * time.Millisecond)
	rm.Stop()

	mu.Lock()
	if attemptCount != 3 {
		t.Errorf("expected 3 attempts with zero delay, got %d", attemptCount)
	}
	mu.Unlock()
}
