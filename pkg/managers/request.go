// Package managers provides high-level manager components for OpenClaw SDK.
//
// This package provides:
//   - EventManager: Pub/sub event management
//   - RequestManager: Pending request correlation
//   - ConnectionManager: WebSocket connection lifecycle
//   - ReconnectManager: Automatic reconnection with Fibonacci backoff
package managers

import (
	"context"
	"sync"

	"github.com/frisbee-ai/openclaw-sdk-go/pkg/protocol"
	"github.com/frisbee-ai/openclaw-sdk-go/pkg/types"
)

// RequestOptions contains options for a pending request.
type RequestOptions struct {
	Timeout    any       // Timeout for the request (time.Duration)
	OnProgress func(any) // Progress callback for intermediate updates
}

// pendingRequest holds state for an in-flight request.
type pendingRequest struct {
	responseCh chan *protocol.ResponseFrame
	onProgress func(any)
}

// RequestManager manages pending requests and their responses.
// It correlates outgoing requests with incoming responses using request IDs.
type RequestManager struct {
	pending    map[string]*pendingRequest    // Map of request ID to pending request
	timeouts   map[string]context.CancelFunc // Map of request ID to timeout cancel function
	mu         sync.Mutex                    // Mutex for thread-safe access
	ctx        context.Context               // Context for lifecycle management
	cancel     context.CancelFunc            // Cancel function for context
	closed     bool                          // Flag indicating manager is closed
	maxPending int                           // Max concurrent pending requests; 0 = unlimited (FOUND-04)
}

// NewRequestManager creates a new request manager.
func NewRequestManager(ctx context.Context) *RequestManager {
	ctx, cancel := context.WithCancel(ctx)
	return &RequestManager{
		pending:    make(map[string]*pendingRequest),
		timeouts:   make(map[string]context.CancelFunc),
		ctx:        ctx,
		cancel:     cancel,
		maxPending: 0, // 0 = unlimited (backward compatible)
	}
}

// SetMaxPending sets the maximum number of concurrent pending requests.
// A value of 0 (default) means unlimited. When the limit is reached,
// SendRequest returns a TooManyPendingRequestsError. FOUND-04.
func (rm *RequestManager) SetMaxPending(max int) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	rm.maxPending = max
}

// SendRequest sends a request and waits for a response.
// It registers the request ID, sends the request via sendFunc, and waits for response.
// Returns the response or an error if the request times out, is cancelled, or the pending limit is exceeded.
// Thread-safe: safe to call from multiple goroutines concurrently.
func (rm *RequestManager) SendRequest(ctx context.Context, req *protocol.RequestFrame, sendFunc func(*protocol.RequestFrame) error) (*protocol.ResponseFrame, error) {
	respCh := make(chan *protocol.ResponseFrame, 1)

	rm.mu.Lock()
	// Check pending request limit (FOUND-04)
	if rm.maxPending > 0 && len(rm.pending) >= rm.maxPending {
		rm.mu.Unlock()
		return nil, types.NewTooManyPendingRequestsError(rm.maxPending)
	}
	rm.pending[req.ID] = &pendingRequest{
		responseCh: respCh,
	}

	// Set up timeout cancellation if context has deadline
	if deadline, ok := ctx.Deadline(); ok {
		timeoutCtx, cancel := context.WithDeadline(ctx, deadline)
		rm.timeouts[req.ID] = cancel
		ctx = timeoutCtx
	}
	rm.mu.Unlock()

	cleanup := func() {
		rm.mu.Lock()
		delete(rm.pending, req.ID)
		if cancel, ok := rm.timeouts[req.ID]; ok {
			cancel()
			delete(rm.timeouts, req.ID)
		}
		// Close channel while holding lock to prevent race with HandleResponse.
		// Channel is owned exclusively by this function -- Clear() and Close() signal
		// via non-blocking send instead of closing, so double-close is impossible.
		close(respCh)
		rm.mu.Unlock()
	}
	defer cleanup()

	// Send the request via transport
	if sendFunc != nil {
		if err := sendFunc(req); err != nil {
			return nil, err
		}
	}

	select {
	case resp := <-respCh:
		// resp is nil when Clear() or Close() unblocks us by sending nil on the channel.
		// In that case return ctx.Err() (context was cancelled).
		if resp == nil {
			return nil, ctx.Err()
		}
		return resp, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// HandleResponse handles an incoming response frame.
// It correlates the response with the pending request using RequestID.
func (rm *RequestManager) HandleResponse(frame *protocol.ResponseFrame) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	req, ok := rm.pending[frame.ID]
	if !ok || req.responseCh == nil {
		return
	}

	select {
	case req.responseCh <- frame:
	default:
	}
}

// ResolveProgress delivers a progress update to the pending request.
func (rm *RequestManager) ResolveProgress(requestID string, payload any) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	req, ok := rm.pending[requestID]
	if !ok || req.onProgress == nil {
		return
	}

	req.onProgress(payload)
}

// AbortRequest aborts a pending request by ID.
func (rm *RequestManager) AbortRequest(requestID string) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	req, ok := rm.pending[requestID]
	if !ok {
		return
	}

	// Send a cancelled response
	cancelledResp := &protocol.ResponseFrame{
		Type:  protocol.FrameTypeResponse,
		ID:    requestID,
		Ok:    false,
		Error: &protocol.ErrorShape{Code: "REQUEST_CANCELLED", Message: "Request cancelled"},
	}

	select {
	case req.responseCh <- cancelledResp:
	default:
	}
}

// Clear cancels all pending requests.
// It signals each waiting goroutine by sending nil on its respCh (non-blocking),
// then removes the entry from the pending map. The cleanup() in each SendRequest
// goroutine will close the respCh after it returns.
// This prevents double-close: Clear/Close never close channels, only SendRequest cleanup does.
func (rm *RequestManager) Clear() {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	if rm.closed {
		return
	}

	for id, req := range rm.pending {
		// Signal the waiting goroutine with nil response (non-blocking send)
		select {
		case req.responseCh <- nil:
		default:
		}
		delete(rm.pending, id)
	}
	for _, cancel := range rm.timeouts {
		cancel()
	}
}

// Close cleans up all pending requests.
// It cancels the context, signals all waiting goroutines, removes pending entries,
// and clears timeout functions. Thread-safe and idempotent.
func (rm *RequestManager) Close() error {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	if rm.closed {
		return nil
	}
	rm.closed = true

	// Cancel context to interrupt any long-running operations
	rm.cancel()

	for id, req := range rm.pending {
		// Signal the waiting goroutine with nil response (non-blocking send).
		// Cleanup() in SendRequest will close respCh after this function returns.
		select {
		case req.responseCh <- nil:
		default:
		}
		delete(rm.pending, id)
	}
	for _, cancel := range rm.timeouts {
		cancel()
	}
	return nil
}
