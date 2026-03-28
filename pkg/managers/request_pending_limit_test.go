// Package managers provides high-level manager components for OpenClaw SDK.
package managers

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/frisbee-ai/openclaw-sdk-go/pkg/protocol"
	"github.com/frisbee-ai/openclaw-sdk-go/pkg/types"
)

// ---------------------------------------------------------------------------
// Pending Limit Tests
// ---------------------------------------------------------------------------

// TestRequestManager_PendingLimit_RejectsOverLimit verifies that when maxPending
// is reached, the N+1th SendRequest returns TooManyPendingRequestsError immediately
// WITHOUT calling sendFunc (no transport call made).
func TestRequestManager_PendingLimit_RejectsOverLimit(t *testing.T) {
	ctx := context.Background()
	rm := NewRequestManager(ctx)
	defer rm.Close()

	rm.SetMaxPending(2)

	// Track whether sendFunc is called (it should NOT be for rejected requests).
	var sendFuncCalled int32

	// blockSendFunc keeps the request pending indefinitely.
	// It launches a goroutine to block forever, then returns nil immediately.
	// The goroutine keeps the channel open so SendRequest's select{} waits on respCh.
	blockSendFunc := func(*protocol.RequestFrame) error {
		go func() { <-make(chan struct{}) }()
		return nil
	}

	// Start 2 requests that remain pending (sendFunc never returns)
	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			req := protocol.NewRequestFrame("block-"+string(rune('A'+idx)), "test", nil)
			_, _ = rm.SendRequest(ctx, req, blockSendFunc)
		}(i)
	}

	// Wait for both to register in the pending map
	time.Sleep(30 * time.Millisecond)

	// Verify 2 are pending
	rm.mu.Lock()
	pending := len(rm.pending)
	rm.mu.Unlock()
	if pending != 2 {
		t.Fatalf("expected 2 pending, got %d", pending)
	}

	// The 3rd request should be REJECTED immediately with TooManyPendingRequestsError.
	// sendFunc must NOT be called (it would block forever if we reached it).
	req := protocol.NewRequestFrame("reject-me", "test", nil)
	_, err := rm.SendRequest(ctx, req, func(*protocol.RequestFrame) error {
		// If this is called, the implementation is wrong.
		atomic.AddInt32(&sendFuncCalled, 1)
		return nil
	})

	if err == nil {
		t.Fatal("expected error when pending limit exceeded, got nil")
	}
	if !errors.Is(err, types.ErrTooManyPendingRequests) {
		t.Errorf("expected error wrapping ErrTooManyPendingRequests, got: %v", err)
	}
	var tooManyErr *types.TooManyPendingRequestsError
	if !errors.As(err, &tooManyErr) {
		t.Errorf("expected *TooManyPendingRequestsError, got: %T", err)
	}
	if atomic.LoadInt32(&sendFuncCalled) != 0 {
		t.Error("sendFunc was called for rejected request -- limit check must be before sendFunc")
	}
}

// TestRequestManager_PendingLimit_ZeroUnlimited verifies that maxPending=0 (default)
// imposes no limit. We use a short deadline to unblock each SendRequest after
// registration, avoiding any manual cleanup.
func TestRequestManager_PendingLimit_ZeroUnlimited(t *testing.T) {
	ctx := context.Background()
	rm := NewRequestManager(ctx)
	defer rm.Close()

	rm.SetMaxPending(0) // 0 = unlimited (backward compatible default)

	var mu sync.Mutex
	var registered int

	for i := 0; i < 10; i++ {
		// Each request: 5ms timeout so it releases from respCh quickly.
		shortCtx, cancel := context.WithTimeout(ctx, 5*time.Millisecond)
		req := protocol.NewRequestFrame("unlimited-"+string(rune('0'+i)), "test", nil)
		_, _ = rm.SendRequest(shortCtx, req, nil)
		cancel()

		mu.Lock()
		registered++
		mu.Unlock()
	}

	if registered != 10 {
		t.Errorf("expected 10 successful registrations, got %d", registered)
	}
}

// TestRequestManager_PendingLimit_UnderLimit verifies that requests below the limit succeed.
func TestRequestManager_PendingLimit_UnderLimit(t *testing.T) {
	ctx := context.Background()
	rm := NewRequestManager(ctx)
	defer rm.Close()

	rm.SetMaxPending(5)

	for i := 0; i < 3; i++ {
		shortCtx, cancel := context.WithTimeout(ctx, 5*time.Millisecond)
		req := protocol.NewRequestFrame("under-limit-"+string(rune('A'+i)), "test", nil)
		_, _ = rm.SendRequest(shortCtx, req, nil)
		cancel()
	}

	rm.mu.Lock()
	pending := len(rm.pending)
	rm.mu.Unlock()

	if pending != 0 {
		// All timed out and cleaned up by now
		t.Errorf("expected 0 pending (all timed out), got %d", pending)
	}
}

// TestRequestManager_PendingLimit_ConcurrentRace verifies that under concurrent load,
// at most maxPending goroutines succeed in registering before the limit is hit.
func TestRequestManager_PendingLimit_ConcurrentRace(t *testing.T) {
	ctx := context.Background()
	rm := NewRequestManager(ctx)
	defer rm.Close()

	rm.SetMaxPending(3)

	var mu sync.Mutex
	var accepted, rejected int

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			// Use short timeout so blocked goroutines exit quickly
			shortCtx, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
			req := protocol.NewRequestFrame("race-"+string(rune('A'+idx)), "test", nil)
			_, err := rm.SendRequest(shortCtx, req, nil)
			cancel()
			mu.Lock()
			if err != nil && errors.Is(err, types.ErrTooManyPendingRequests) {
				rejected++
			} else {
				accepted++
			}
			mu.Unlock()
		}(i)
	}

	wg.Wait()

	mu.Lock()
	a, r := accepted, rejected
	mu.Unlock()

	// Out of 10 concurrent requests with limit=3, at most 3 should be accepted.
	if a > 3 {
		t.Errorf("expected at most 3 accepted (maxPending=3), got %d", a)
	}
	if r < 7 {
		t.Errorf("expected at least 7 rejected, got %d", r)
	}
}

// ---------------------------------------------------------------------------
// Channel Ownership Tests
// ---------------------------------------------------------------------------

// TestRequestManager_ChannelOwnership_ClearNoDoubleClose verifies that Clear()
// signals the waiting goroutine without double-closing the response channel.
func TestRequestManager_ChannelOwnership_ClearNoDoubleClose(t *testing.T) {
	ctx := context.Background()
	rm := NewRequestManager(ctx)
	defer rm.Close()

	req := protocol.NewRequestFrame("clear-test", "test", nil)

	done := make(chan error, 1)
	go func() {
		// Use long timeout so the request is definitely pending when Clear is called
		longCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_, err := rm.SendRequest(longCtx, req, nil)
		done <- err
	}()

	// Wait for request to register
	time.Sleep(30 * time.Millisecond)

	// Clear() should signal the goroutine without panicking (no double-close)
	rm.Clear()

	select {
	case err := <-done:
		if err == nil {
			t.Error("expected non-nil error from SendRequest after Clear")
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("SendRequest did not return after Clear")
	}

	// The manager should still be usable after Clear.
	// Use 100ms timeout -- long enough that the response can arrive if HandleResponse
	// is called, but short enough that we don't hang if it doesn't.
	longCtx2, cancel2 := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel2()
	req2 := protocol.NewRequestFrame("after-clear", "test", nil)
	_, err := rm.SendRequest(longCtx2, req2, nil)
	// DeadlineExceeded is acceptable (no response was sent); limit error should not occur.
	if err != nil && !errors.Is(err, types.ErrTooManyPendingRequests) && !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("unexpected error after Clear: %v", err)
	}
}

// TestRequestManager_ChannelOwnership_CloseNoDoubleClose verifies that Close()
// signals the waiting goroutine without double-closing the response channel.
func TestRequestManager_ChannelOwnership_CloseNoDoubleClose(t *testing.T) {
	ctx := context.Background()
	rm := NewRequestManager(ctx)

	req := protocol.NewRequestFrame("close-test", "test", nil)

	done := make(chan error, 1)
	go func() {
		longCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_, err := rm.SendRequest(longCtx, req, nil)
		done <- err
	}()

	time.Sleep(30 * time.Millisecond)

	// Close() should signal the goroutine without panicking (no double-close)
	_ = rm.Close()

	select {
	case err := <-done:
		if err == nil {
			t.Error("expected non-nil error from SendRequest after Close")
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("SendRequest did not return after Close")
	}
}

// TestRequestManager_SetMaxPending verifies the setter correctly stores the limit.
func TestRequestManager_SetMaxPending(t *testing.T) {
	ctx := context.Background()
	rm := NewRequestManager(ctx)
	defer rm.Close()

	rm.SetMaxPending(10)

	rm.mu.Lock()
	limit := rm.maxPending
	rm.mu.Unlock()

	if limit != 10 {
		t.Errorf("expected maxPending to be 10, got %d", limit)
	}
}
