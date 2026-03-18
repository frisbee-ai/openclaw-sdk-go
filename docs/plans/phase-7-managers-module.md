# Phase 7: Managers Module

**Files:**
- Create: `managers/event.go`, `managers/event_test.go`
- Create: `managers/request.go`, `managers/request_test.go`
- Create: `managers/connection.go`
- Create: `managers/reconnect.go`

**Depends on:** Phase 1 (types.go), Phase 5 (connection), Phase 6 (events), Phase 4 (transport)

---

## Task 7.1: Event Manager

- [ ] **Step 1: Create managers directory and event.go**

```bash
mkdir -p managers
```

```go
// managers/event.go
package managers

import (
	"context"
	"sync"
	"time"

	"openclaw-sdk-go"
)

// EventManager manages event subscriptions and dispatching
type EventManager struct {
	events   chan openclaw.Event
	handlers map[openclaw.EventType][]openclaw.EventHandler
	ctx      context.Context
	cancel   context.CancelFunc
	mu       sync.RWMutex
	wg       sync.WaitGroup
}

// NewEventManager creates a new event manager
func NewEventManager(ctx context.Context, bufferSize int) *EventManager {
	ctx, cancel := context.WithCancel(ctx)
	return &EventManager{
		events:   make(chan openclaw.Event, bufferSize),
		handlers: make(map[openclaw.EventType][]openclaw.EventHandler),
		ctx:      ctx,
		cancel:   cancel,
	}
}

// Subscribe adds an event handler
func (em *EventManager) Subscribe(eventType openclaw.EventType, handler openclaw.EventHandler) func() {
	em.mu.Lock()
	defer em.mu.Unlock()
	em.handlers[eventType] = append(em.handlers[eventType], handler)
	return func() { em.Unsubscribe(eventType, handler) }
}

// Unsubscribe removes an event handler
func (em *EventManager) Unsubscribe(eventType openclaw.EventType, handler openclaw.EventHandler) {
	em.mu.Lock()
	defer em.mu.Unlock()
	handlers := em.handlers[eventType]
	for i, h := range handlers {
		if h == handler {
			handlers[i] = handlers[len(handlers)-1]
			em.handlers[eventType] = handlers[:len(handlers)-1]
			return
		}
	}
}

// Events returns the event channel
func (em *EventManager) Events() <-chan openclaw.Event {
	return em.events
}

// Emit emits an event
func (em *EventManager) Emit(event openclaw.Event) {
	select {
	case em.events <- event:
	case <-em.ctx.Done():
	}
}

// Start begins the event dispatch loop
func (em *EventManager) Start() {
	em.wg.Add(1)
	go func() {
		defer em.wg.Done()
		for {
			select {
			case <-em.ctx.Done():
				return
			case event := <-em.events:
				em.dispatch(event)
			}
		}
	}()
}

// dispatch sends event to all registered handlers
func (em *EventManager) dispatch(event openclaw.Event) {
	em.mu.RLock()
	defer em.mu.RUnlock()
	for _, handler := range em.handlers[event.Type] {
		defer func() { recover() }()
		handler(event)
	}
}

// Close gracefully shuts down the event manager
func (em *EventManager) Close() error {
	em.cancel()
	em.wg.Wait()
	close(em.events)
	return nil
}
```

- [ ] **Step 2: Write test**

```go
// managers/event_test.go
package managers

import (
	"context"
	"testing"
	"time"

	"openclaw-sdk-go"
)

func TestEventManager(t *testing.T) {
	ctx := context.Background()
	em := NewEventManager(ctx, 10)

	handlerCalled := false
	em.Subscribe(openclaw.EventConnect, func(e openclaw.Event) {
		handlerCalled = true
	})

	em.Start()
	em.Emit(openclaw.Event{Type: openclaw.EventConnect, Timestamp: time.Now()})
	time.Sleep(10 * time.Millisecond)

	if !handlerCalled {
		t.Error("expected handler to be called")
	}

	em.Close()
}
```

- [ ] **Step 3: Run tests and commit**

Run: `go test -v ./managers/...`
Commit: `git add managers/event.go managers/event_test.go && git commit -m "feat: add event manager"`

---

## Task 7.2: Request Manager

- [ ] **Step 1: Write request.go**

```go
// managers/request.go
package managers

import (
	"context"
	"sync"
	"time"

	"openclaw-sdk-go/protocol"
)

// RequestManager manages pending requests
type RequestManager struct {
	pending  map[string]chan *protocol.ResponseFrame
	timeouts map[string]context.CancelFunc
	mu       sync.Mutex
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
}

// NewRequestManager creates a new request manager
func NewRequestManager(ctx context.Context) *RequestManager {
	ctx, cancel := context.WithCancel(ctx)
	return &RequestManager{
		pending:  make(map[string]chan *protocol.ResponseFrame),
		timeouts: make(map[string]context.CancelFunc),
		ctx:      ctx,
		cancel:   cancel,
	}
}

// SendRequest sends a request and waits for response
func (rm *RequestManager) SendRequest(ctx context.Context, frame *protocol.RequestFrame) (*protocol.ResponseFrame, error) {
	respCh := make(chan *protocol.ResponseFrame, 1)

	rm.mu.Lock()
	rm.pending[frame.RequestID] = respCh
	rm.mu.Unlock()

	cleanup := func() {
		rm.mu.Lock()
		delete(rm.pending, frame.RequestID)
		if cancel, ok := rm.timeouts[frame.RequestID]; ok {
			cancel()
			delete(rm.timeouts, frame.RequestID)
		}
		rm.mu.Unlock()
		close(respCh)
	}
	defer cleanup()

	// TODO: Send via transport

	select {
	case resp := <-respCh:
		return resp, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// HandleResponse handles an incoming response
func (rm *RequestManager) HandleResponse(frame *protocol.ResponseFrame) {
	rm.mu.Lock()
	ch, ok := rm.pending[frame.RequestID]
	rm.mu.Unlock()

	if ok && ch != nil {
		select {
		case ch <- frame:
		default:
		}
	}
}

// Close cleans up all pending requests
func (rm *RequestManager) Close() error {
	rm.cancel()
	rm.mu.Lock()
	defer rm.mu.Unlock()
	for id, ch := range rm.pending {
		close(ch)
		delete(rm.pending, id)
	}
	for _, cancel := range rm.timeouts {
		cancel()
	}
	return nil
}
```

- [ ] **Step 2: Write test**

```go
// managers/request_test.go
package managers

import (
	"context"
	"testing"
	"time"

	"openclaw-sdk-go/protocol"
)

func TestRequestManager(t *testing.T) {
	ctx := context.Background()
	rm := NewRequestManager(ctx)

	frame := &protocol.RequestFrame{
		RequestID: "test-123",
		Method:   "test",
		Timestamp: time.Now(),
	}

	resp := &protocol.ResponseFrame{
		RequestID: "test-123",
		Success:   true,
		Timestamp: time.Now(),
	}

	go func() {
		time.Sleep(10 * time.Millisecond)
		rm.HandleResponse(resp)
	}()

	got, err := rm.SendRequest(context.Background(), frame)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.RequestID != "test-123" {
		t.Errorf("expected 'test-123', got '%s'", got.RequestID)
	}

	rm.Close()
}
```

- [ ] **Step 3: Run tests and commit**

Run: `go test -v ./managers/...`
Commit: `git add managers/request.go managers/request_test.go && git commit -m "feat: add request manager"`

---

## Task 7.3: Connection Manager

- [ ] **Step 1: Write connection.go**

```go
// managers/connection.go
package managers

import (
	"context"
	"sync"
	"time"

	"openclaw-sdk-go"
	"openclaw-sdk-go/connection"
	"openclaw-sdk-go/transport"
)

// ConnectionManager manages WebSocket connections
type ConnectionManager struct {
	config    *ClientConfig
	state     *connection.ConnectionStateMachine
	transport *transport.WebSocketTransport
	eventMgr  *EventManager
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
	mu        sync.Mutex
}

// NewConnectionManager creates a new connection manager
func NewConnectionManager(config *ClientConfig, eventMgr *EventManager) *ConnectionManager {
	ctx, cancel := context.WithCancel(context.Background())
	return &ConnectionManager{
		config:    config,
		state:     connection.NewConnectionStateMachine("disconnected"),
		eventMgr:  eventMgr,
		ctx:       ctx,
		cancel:    cancel,
	}
}

// Connect establishes a connection
func (cm *ConnectionManager) Connect(ctx context.Context) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if cm.transport != nil {
		return openclaw.NewError(openclaw.ErrCodeConnection, "already connected", nil)
	}

	if err := cm.state.Transition("connecting", nil); err != nil {
		return err
	}

	t, err := transport.Dial(cm.config.URL, &transport.WebSocketConfig{
		Header: cm.config.Header,
	})
	if err != nil {
		cm.state.Transition("failed", err)
		return err
	}

	cm.transport = t
	t.Start()

	if err := cm.state.Transition("connected", nil); err != nil {
		return err
	}

	cm.eventMgr.Emit(openclaw.Event{
		Type:      openclaw.EventConnect,
		Timestamp: time.Now(),
	})

	return nil
}

// Disconnect closes the connection
func (cm *ConnectionManager) Disconnect() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if cm.transport == nil {
		return nil
	}

	err := cm.transport.Close()
	cm.transport = nil
	cm.state.Transition("disconnected", nil)

	cm.eventMgr.Emit(openclaw.Event{
		Type:      openclaw.EventDisconnect,
		Timestamp: time.Now(),
	})

	return err
}

// State returns the current connection state
func (cm *ConnectionManager) State() openclaw.ConnectionState {
	return openclaw.ConnectionState(cm.state.State())
}

// Transport returns the underlying transport
func (cm *ConnectionManager) Transport() *transport.WebSocketTransport {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	return cm.transport
}

// Close closes the connection manager
func (cm *ConnectionManager) Close() error {
	cm.cancel()
	return cm.Disconnect()
}
```

- [ ] **Step 2: Commit**

```bash
git add managers/connection.go
git commit -m "feat: add connection manager"
```

---

## Task 7.4: Reconnect Manager

- [ ] **Step 1: Write reconnect.go**

```go
// managers/reconnect.go
package managers

import (
	"context"
	"math"
	"sync"
	"time"
)

// ReconnectConfig holds reconnection configuration
type ReconnectConfig struct {
	MaxAttempts      int
	InitialDelay     time.Duration
	MaxDelay         time.Duration
	BackoffMultiplier float64
}

// DefaultReconnectConfig returns default configuration
func DefaultReconnectConfig() *ReconnectConfig {
	return &ReconnectConfig{
		MaxAttempts:      0,
		InitialDelay:     1 * time.Second,
		MaxDelay:         60 * time.Second,
		BackoffMultiplier: 1.618,
	}
}

// ReconnectManager handles automatic reconnection
type ReconnectManager struct {
	config             *ReconnectConfig
	attempts          int
	mu                sync.Mutex
	ctx               context.Context
	cancel            context.CancelFunc
	wg                sync.WaitGroup
	onReconnect       func() error
	onReconnectFailed func(err error)
}

// NewReconnectManager creates a new reconnect manager
func NewReconnectManager(config *ReconnectConfig) *ReconnectManager {
	if config == nil {
		config = DefaultReconnectConfig()
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &ReconnectManager{
		config: config,
		ctx:    ctx,
		cancel: cancel,
	}
}

// SetOnReconnect sets the reconnect callback
func (rm *ReconnectManager) SetOnReconnect(f func() error) {
	rm.onReconnect = f
}

// SetOnReconnectFailed sets the reconnect failed callback
func (rm *ReconnectManager) SetOnReconnectFailed(f func(err error)) {
	rm.onReconnectFailed = f
}

// Start begins the reconnection loop
func (rm *ReconnectManager) Start() {
	rm.wg.Add(1)
	go rm.run()
}

func (rm *ReconnectManager) run() {
	defer rm.wg.Done()

	delay := rm.config.InitialDelay
	attempt := 0

	for {
		attempt++
		select {
		case <-rm.ctx.Done():
			return
		case <-time.After(delay):
			if rm.onReconnect != nil {
				err := rm.onReconnect()
				if err == nil {
					return
				}
				if rm.onReconnectFailed != nil {
					rm.onReconnectFailed(err)
				}
			}

			if rm.config.MaxAttempts > 0 && attempt >= rm.config.MaxAttempts {
				return
			}

			delay = time.Duration(float64(delay) * rm.config.BackoffMultiplier)
			if delay > rm.config.MaxDelay {
				delay = rm.config.MaxDelay
			}
		}
	}
}

// Stop stops the reconnection attempts
func (rm *ReconnectManager) Stop() {
	rm.cancel()
	rm.wg.Wait()
}

// Reset resets the attempt counter
func (rm *ReconnectManager) Reset() {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	rm.attempts = 0
}
```

- [ ] **Step 2: Commit**

```bash
git add managers/reconnect.go
git commit -m "feat: add reconnect manager"
```

---

## Phase 7 Complete

After this phase, you should have:
- `managers/event.go` - Event manager
- `managers/request.go` - Request manager
- `managers/connection.go` - Connection manager
- `managers/reconnect.go` - Reconnect manager

All code should compile and tests should pass.
