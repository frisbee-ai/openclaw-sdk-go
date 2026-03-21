package managers

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"

	"github.com/frisbee-ai/openclaw-sdk-go/pkg/protocol"
)

func TestRequestManager_SendAndReceive(t *testing.T) {
	ctx := context.Background()
	rm := NewRequestManager(ctx)

	req := protocol.NewRequestFrame("test-123", "test", nil)

	resp := protocol.NewResponseFrameSuccess("test-123", json.RawMessage(`{}`))

	// Send response after a small delay
	go func() {
		time.Sleep(10 * time.Millisecond)
		rm.HandleResponse(resp)
	}()

	got, err := rm.SendRequest(context.Background(), req, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != "test-123" {
		t.Errorf("expected 'test-123', got '%s'", got.ID)
	}

	_ = rm.Close()
}

func TestRequestManager_ContextCancellation(t *testing.T) {
	ctx := context.Background()
	rm := NewRequestManager(ctx)

	req := protocol.NewRequestFrame("test-cancel", "test", nil)

	// Create cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := rm.SendRequest(ctx, req, nil)
	if err == nil {
		t.Error("expected error for cancelled context")
	}

	_ = rm.Close()
}

// TestRequestManager_RaceHandleResponseAndTimeout tests the race condition between
// HandleResponse and cleanup (context timeout). This test uses -race flag to detect.
func TestRequestManager_RaceHandleResponseAndTimeout(t *testing.T) {
	ctx := context.Background()
	rm := NewRequestManager(ctx)
	defer rm.Close()

	req := protocol.NewRequestFrame("race-test", "test", nil)

	// Create context with very short deadline to trigger timeout
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(1*time.Millisecond))
	defer cancel()

	resp := protocol.NewResponseFrameSuccess("race-test", json.RawMessage(`{}`))

	var wg sync.WaitGroup
	wg.Add(2)

	// Concurrently send request (which will timeout) and handle response
	go func() {
		defer wg.Done()
		rm.SendRequest(ctx, req, nil)
	}()

	// Small delay to increase chance of race condition
	go func() {
		time.Sleep(500 * time.Microsecond)
		rm.HandleResponse(resp)
		wg.Done()
	}()

	wg.Wait()
}

// TestRequestManager_ConcurrentHandleResponse tests concurrent HandleResponse calls
// that race with cleanup operations.
func TestRequestManager_ConcurrentHandleResponse(t *testing.T) {
	ctx := context.Background()
	rm := NewRequestManager(ctx)
	defer rm.Close()

	const numPairs = 10

	var wg sync.WaitGroup
	wg.Add(numPairs * 2) // 2 goroutines per pair (SendRequest + HandleResponse)

	for i := 0; i < numPairs; i++ {
		reqID := "concurrent-test"
		req := protocol.NewRequestFrame(reqID, "test", nil)
		resp := protocol.NewResponseFrameSuccess(reqID, json.RawMessage(`{}`))

		// Each pair: one goroutine for SendRequest, one for HandleResponse
		go func(requestID string, response *protocol.ResponseFrame) {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()
			rm.SendRequest(ctx, req, nil)
		}(reqID, resp)

		go func(response *protocol.ResponseFrame) {
			defer wg.Done()
			rm.HandleResponse(response)
		}(resp)
	}

	wg.Wait()
}

// TestRequestManager_AbortRequestRace tests race between AbortRequest and HandleResponse.
func TestRequestManager_AbortRequestRace(t *testing.T) {
	ctx := context.Background()
	rm := NewRequestManager(ctx)
	defer rm.Close()

	req := protocol.NewRequestFrame("abort-race", "test", nil)
	resp := protocol.NewResponseFrameSuccess("abort-race", json.RawMessage(`{}`))

	var wg sync.WaitGroup
	wg.Add(3)

	// Start request that will be aborted
	go func() {
		defer wg.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()
		rm.SendRequest(ctx, req, nil)
	}()

	// Abort the request
	go func() {
		defer wg.Done()
		time.Sleep(10 * time.Millisecond)
		rm.AbortRequest("abort-race")
	}()

	// Handle response concurrently
	go func() {
		defer wg.Done()
		time.Sleep(5 * time.Millisecond)
		rm.HandleResponse(resp)
	}()

	wg.Wait()
}

func TestRequestManager_ResolveProgress(t *testing.T) {
	ctx := context.Background()
	rm := NewRequestManager(ctx)
	defer rm.Close()

	var progressCalled bool
	var progressValue any
	req := protocol.NewRequestFrame("progress-test", "test", nil)

	// Register request with progress callback
	rm.mu.Lock()
	rm.pending[req.ID] = &pendingRequest{
		responseCh: make(chan *protocol.ResponseFrame, 1),
		onProgress: func(v any) {
			progressCalled = true
			progressValue = v
		},
	}
	rm.mu.Unlock()

	// Deliver progress update
	rm.ResolveProgress("progress-test", "test-progress-data")

	if !progressCalled {
		t.Error("expected progress callback to be called")
	}
	if progressValue != "test-progress-data" {
		t.Errorf("expected 'test-progress-data', got '%v'", progressValue)
	}
}

func TestRequestManager_ResolveProgress_NotFound(t *testing.T) {
	ctx := context.Background()
	rm := NewRequestManager(ctx)
	defer rm.Close()

	// Should not panic for non-existent request
	rm.ResolveProgress("non-existent", "data")
}

func TestRequestManager_ResolveProgress_NoCallback(t *testing.T) {
	ctx := context.Background()
	rm := NewRequestManager(ctx)
	defer rm.Close()

	req := protocol.NewRequestFrame("no-callback", "test", nil)

	// Register request without progress callback
	rm.mu.Lock()
	rm.pending[req.ID] = &pendingRequest{
		responseCh: make(chan *protocol.ResponseFrame, 1),
		onProgress: nil,
	}
	rm.mu.Unlock()

	// Should not panic
	rm.ResolveProgress("no-callback", "data")
}

func TestRequestManager_Clear(t *testing.T) {
	ctx := context.Background()
	rm := NewRequestManager(ctx)

	req := protocol.NewRequestFrame("clear-test", "test", nil)

	// Register request
	rm.mu.Lock()
	rm.pending[req.ID] = &pendingRequest{
		responseCh: make(chan *protocol.ResponseFrame, 1),
	}
	rm.mu.Unlock()

	// Clear should cancel all pending requests
	rm.Clear()

	rm.mu.Lock()
	if len(rm.pending) != 0 {
		t.Errorf("expected pending to be empty after Clear, got %d", len(rm.pending))
	}
	rm.mu.Unlock()

	_ = rm.Close()
}

func TestRequestManager_Clear_AfterClose(t *testing.T) {
	ctx := context.Background()
	rm := NewRequestManager(ctx)
	_ = rm.Close()

	// Clear after close should be no-op (closed flag is checked)
	rm.Clear()
}

func TestRequestManager_Close(t *testing.T) {
	ctx := context.Background()
	rm := NewRequestManager(ctx)

	req := protocol.NewRequestFrame("close-test", "test", nil)

	// Register request
	rm.mu.Lock()
	rm.pending[req.ID] = &pendingRequest{
		responseCh: make(chan *protocol.ResponseFrame, 1),
	}
	rm.mu.Unlock()

	// Close should cleanup
	err := rm.Close()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	rm.mu.Lock()
	if len(rm.pending) != 0 {
		t.Errorf("expected pending to be empty after Close, got %d", len(rm.pending))
	}
	closed := rm.closed
	rm.mu.Unlock()

	if !closed {
		t.Error("expected closed flag to be set after Close")
	}
}

func TestRequestManager_Close_Idempotent(t *testing.T) {
	ctx := context.Background()
	rm := NewRequestManager(ctx)

	// Multiple closes should be safe
	err1 := rm.Close()
	err2 := rm.Close()

	if err1 != nil {
		t.Errorf("first Close error: %v", err1)
	}
	if err2 != nil {
		t.Errorf("second Close error: %v", err2)
	}
}
