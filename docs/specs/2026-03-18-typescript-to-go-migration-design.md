# OpenClaw SDK Go Migration Design

**Date**: 2026-03-18
**Status**: Revised (v2 - addressing review feedback)
**Source**: openclaw-sdk-typescript (TypeScript)

## 1. Overview

This document describes the design for migrating the OpenClaw SDK from TypeScript to Go. The Go SDK will provide feature-parity with the TypeScript version while following Go idioms and best practices.

### 1.1 Goals

- **Feature Parity**: All TypeScript SDK functionality available in Go
- **Go Idioms**: Follow Go standard library patterns and conventions
- **Dual Support**: Support both server-side and CLI usage scenarios
- **Standard Library First**: Minimize external dependencies

### 1.2 Design Decisions Summary

| Aspect | Decision |
|--------|----------|
| API Style | Go idiomatic (Option pattern, context.Context) |
| Architecture | Context + Channel hybrid |
| Dependencies | Standard library preferred |
| Compatibility | Feature-equivalent, not API-identical |

---

## 2. Project Structure

```
openclaw-sdk-go/
├── go.mod
├── go.sum
├── LICENSE
├── README.md
├── client.go              // Main OpenClawClient
├── errors.go              // Error type hierarchy
├── types.go               // Common type definitions
├── client_test.go
├── errors_test.go
│
├── auth/
│   ├── provider.go        // CredentialsProvider interface
│   ├── provider_test.go
│   └── handler.go         // AuthHandler
│
├── transport/
│   ├── websocket.go        // WebSocketTransport
│   ├── websocket_test.go
│   └── ...
│
├── protocol/
│   ├── types.go           // GatewayFrame, RequestFrame, ResponseFrame
│   ├── types_test.go
│   └── validation.go      // Frame validation
│   └── validation_test.go
│
├── connection/
│   ├── protocol.go        // ProtocolNegotiator
│   ├── state.go           // ConnectionStateMachine
│   ├── policies.go        // PolicyManager
│   ├── tls.go             // TlsValidator
│   └── ...
│
├── events/
│   ├── tick.go            // TickMonitor (heartbeat)
│   ├── tick_test.go
│   ├── gap.go             // GapDetector (message gap detection)
│   └── gap_test.go
│
├── managers/
│   ├── connection.go      // ConnectionManager
│   ├── request.go         // RequestManager
│   ├── event.go           // EventManager
│   ├── reconnect.go       // ReconnectManager
│   └── ...
│
├── utils/
│   └── timeout.go         // TimeoutManager
│   └── timeout_test.go
│
└── examples/
    ├── cmd/
    │   └── main.go        // CLI example
    └── server/
        └── main.go        // Server example
```

---

## 3. Core Design Patterns

### 3.1 Option Pattern

Instead of TypeScript constructor parameters, Go uses functional options:

```go
// Usage: client, err := NewClient(WithURL("ws://..."), WithTimeout(30 * time.Second))

type ClientOption func(*ClientConfig) error

type ClientConfig struct {
    URL              string
    AuthHandler      AuthHandler
    ReconnectEnabled bool
    ReconnectConfig  *ReconnectConfig
    Logger           Logger
    // ... more options
}

func WithURL(url string) ClientOption { ... }
func WithAuthHandler(handler AuthHandler) ClientOption { ... }
func WithReconnect(enabled bool) ClientOption { ... }
func WithLogger(logger Logger) ClientOption { ... }
func WithTimeout(timeout time.Duration) ClientOption { ... }
```

### 3.2 Context + Channel Hybrid Events

The event system uses Go channels for event delivery with context for lifecycle management:

```go
// EventManager - Pub/Sub with channels
// Thread-safe event manager with buffered channel to prevent deadlocks
type EventManager struct {
    events   chan Event  // Buffered channel to prevent blocking
    handlers map[EventType][]EventHandler
    ctx      context.Context
    cancel   context.CancelFunc
    mu       sync.RWMutex
    wg       sync.WaitGroup
}

// NewEventManager creates a new EventManager with buffered events channel
func NewEventManager(ctx context.Context, bufferSize int) *EventManager {
    ctx, cancel := context.WithCancel(ctx)
    return &EventManager{
        events:   make(chan Event, bufferSize),  // Buffered to prevent deadlocks
        handlers: make(map[EventType][]EventHandler),
        ctx:      ctx,
        cancel:   cancel,
    }
}

// Close gracefully shuts down the event manager
func (em *EventManager) Close() error {
    em.cancel()
    em.wg.Wait()
    close(em.events)
    return nil
}

// Event types
type EventType string

const (
    EventConnect        EventType = "connect"
    EventDisconnect     EventType = "disconnect"
    EventError          EventType = "error"
    EventMessage        EventType = "message"
    EventRequest        EventType = "request"
    EventResponse       EventType = "response"
    EventTick           EventType = "tick"
    EventGap            EventType = "gap"
    EventStateChange    EventType = "stateChange"
)

// Event structure
type Event struct {
    Type    EventType
    Payload interface{}
    Err     error  // For Error events
    Timestamp time.Time
}

// EventHandler type
type EventHandler func(Event)

// Subscribe returns unsubscribe function
func (em *EventManager) Subscribe(eventType EventType, handler EventHandler) func() {
    em.mu.Lock()
    defer em.mu.Unlock()
    em.handlers[eventType] = append(em.handlers[eventType], handler)
    return func() { em.Unsubscribe(eventType, handler) }
}

// Unsubscribe removes a handler from the event type
func (em *EventManager) Unsubscribe(eventType EventType, handler EventHandler) {
    em.mu.Lock()
    defer em.mu.Unlock()
    handlers := em.handlers[eventType]
    for i, h := range handlers {
        if h == handler {
            // Remove handler by swapping with last and trimming
            handlers[i] = handlers[len(handlers)-1]
            em.handlers[eventType] = handlers[:len(handlers)-1]
            return
        }
    }
}

// Events returns readonly event channel
func (em *EventManager) Events() <-chan Event {
    return em.events
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
func (em *EventManager) dispatch(event Event) {
    em.mu.RLock()
    defer em.mu.RUnlock()
    for _, handler := range em.handlers[event.Type] {
        // Handle panics from handlers
        defer func() { recover() }()
        handler(event)
    }
}
```

### 3.3 Request/Response with Context

Requests use context for timeout/cancellation. Thread-safe with proper cleanup:

```go
type RequestManager struct {
    pending  map[string]chan *ResponseFrame
    timeouts map[string]context.CancelFunc
    mu       sync.Mutex
    ctx      context.Context
    cancel   context.CancelFunc
}

// NewRequestManager creates a new RequestManager
func NewRequestManager(ctx context.Context) *RequestManager {
    ctx, cancel := context.WithCancel(ctx)
    return &RequestManager{
        pending: make(map[string]chan *ResponseFrame),
        timeouts: make(map[string]context.CancelFunc),
        ctx:    ctx,
        cancel: cancel,
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

func (rm *RequestManager) SendRequest(ctx context.Context, frame *RequestFrame) (*ResponseFrame, error) {
    respCh := make(chan *ResponseFrame, 1)

    // Register request BEFORE releasing lock
    rm.mu.Lock()
    rm.pending[frame.RequestID] = respCh
    rm.mu.Unlock()

    // Ensure cleanup happens - cancel timeout on success, remove from pending
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

    // Send frame via transport
    // ...

    // Wait for response or context cancellation
    select {
    case resp := <-respCh:
        return resp, nil
    case <-ctx.Done():
        return nil, ctx.Err()
    }
}
```

### 3.4 Connection State Machine

Thread-safe state machine with buffered event channel to prevent deadlocks:

```go
type ConnectionState string

const (
    StateDisconnected      ConnectionState = "disconnected"
    StateConnecting        ConnectionState = "connecting"
    StateConnected         ConnectionState = "connected"
    StateAuthenticating    ConnectionState = "authenticating"
    StateAuthenticated     ConnectionState = "authenticated"
    StateReconnecting      ConnectionState = "reconnecting"
    StateFailed            ConnectionState = "failed"
)

type ConnectionStateMachine struct {
    state   ConnectionState
    mu      sync.RWMutex
    events  chan StateChangeEvent  // Buffered to prevent blocking
}

type StateChangeEvent struct {
    From    ConnectionState
    To      ConnectionState
    Reason  error
}

// NewConnectionStateMachine creates a new state machine
func NewConnectionStateMachine(initial ConnectionState) *ConnectionStateMachine {
    return &ConnectionStateMachine{
        state:  initial,
        events: make(chan StateChangeEvent, 10),  // Buffered
    }
}

// Transition atomically changes state and sends event
// IMPORTANT: Must NOT send to channel while holding lock to prevent deadlock
func (csm *ConnectionStateMachine) Transition(to ConnectionState, reason error) error {
    csm.mu.Lock()
    from := csm.state
    if !csm.validTransition(from, to) {
        csm.mu.Unlock()
        return fmt.Errorf("invalid state transition from %s to %s", from, to)
    }
    csm.state = to
    csm.mu.Unlock()  // Release lock BEFORE sending to channel

    // Send event outside lock to prevent deadlock
    select {
    case csm.events <- StateChangeEvent{From: from, To: to, Reason: reason}:
    default:
        // Channel full - log warning but don't block
    }
    return nil
}
```

### 3.5 WebSocket Transport

Using `gorilla/websocket` (the industry standard Go WebSocket library):

```go
import "github.com/gorilla/websocket"

type WebSocketTransport struct {
    conn    *websocket.Conn
    sendCh  chan []byte
    recvCh  chan []byte
    ctx     context.Context
    cancel  context.CancelFunc
    mu      sync.Mutex
    wg      sync.WaitGroup
}

type WebSocketConfig struct {
    ReadBufferSize  int
    WriteBufferSize int
    TLSConfig       *tls.Config
}

// Dial creates a new WebSocket connection
func Dial(url string, header http.Header, config *WebSocketConfig) (*WebSocketTransport, error) {
    dialer := websocket.Dialer{
        ReadBufferSize:  config.ReadBufferSize,
        WriteBufferSize: config.WriteBufferSize,
        TLSClientConfig: config.TLSConfig,
    }

    conn, _, err := dialer.Dial(url, header)
    if err != nil {
        return nil, err
    }

    ctx, cancel := context.WithCancel(context.Background())
    return &WebSocketTransport{
        conn:   conn,
        sendCh: make(chan []byte, 10),
        recvCh: make(chan []byte, 10),
        ctx:    ctx,
        cancel: cancel,
    }, nil
}

// Close gracefully closes the WebSocket connection
func (t *WebSocketTransport) Close() error {
    t.cancel()
    t.wg.Wait()
    return t.conn.Close()
}
```

---

## 4. Error Type Hierarchy

```go
type OpenClawError interface {
    error
    Code() ErrorCode
    Unwrap() error
}

type ErrorCode string

const (
    ErrCodeConnection    ErrorCode = "CONNECTION_ERROR"
    ErrCodeAuth          ErrorCode = "AUTH_ERROR"
    ErrCodeTimeout       ErrorCode = "TIMEOUT"
    ErrCodeProtocol      ErrorCode = "PROTOCOL_ERROR"
    ErrCodeValidation    ErrorCode = "VALIDATION_ERROR"
    ErrCodeTransport     ErrorCode = "TRANSPORT_ERROR"
    ErrCodeUnknown       ErrorCode = "UNKNOWN"
)

// Concrete error types
type ConnectionError struct {
    code    ErrorCode
    message string
    err     error
}

func (e *ConnectionError) Error() string { return e.message }
func (e *ConnectionError) Code() ErrorCode { return e.code }
func (e *ConnectionError) Unwrap() error { return e.err }

type AuthError struct { /* ... */ }
type TimeoutError struct { /* ... */ }
type ProtocolError struct { /* ... */ }
type ValidationError struct { /* ... */ }
type TransportError struct { /* ... */ }

// Error factory
func NewError(code ErrorCode, message string, err error) OpenClawError { ... }
```

---

## 5. Main Client API

### 5.1 Logger Interface

```go
// Logger interface for customizable logging
// Uses Go 1.21+ slog interface for compatibility
type Logger interface {
    Debug(msg string, args ...any)
    Info(msg string, args ...any)
    Warn(msg string, args ...any)
    Error(msg string, args ...any)
}

// DefaultLogger uses stdlib log
type DefaultLogger struct{}

func (l *DefaultLogger) Debug(msg string, args ...any) { log.Printf("[DEBUG] "+msg, args...) }
func (l *DefaultLogger) Info(msg string, args ...any)  { log.Printf("[INFO] "+msg, args...) }
func (l *DefaultLogger) Warn(msg string, args ...any)  { log.Printf("[WARN] "+msg, args...) }
func (l *DefaultLogger) Error(msg string, args ...any) { log.Printf("[ERROR] "+msg, args...) }
```

### 5.2 TLS Configuration

```go
// TLSConfig holds TLS configuration options
type TLSConfig struct {
    InsecureSkipVerify bool          // Skip certificate verification (dev only)
    CertFile          string        // Client certificate file
    KeyFile           string        // Client key file
    CAFile            string        // CA certificate file
    ServerName        string        // Server name for SNI
}
```

### 5.3 Reconnect Configuration

```go
// ReconnectConfig holds reconnection settings
type ReconnectConfig struct {
    MaxAttempts     int           // Maximum reconnection attempts (0 = infinite)
    InitialDelay    time.Duration // Initial delay (default: 1s)
    MaxDelay        time.Duration // Maximum delay (default: 60s)
    BackoffMultiplier float64    // Fibonacci backoff multiplier
}

// DefaultReconnectConfig returns sensible defaults
func DefaultReconnectConfig() *ReconnectConfig {
    return &ReconnectConfig{
        MaxAttempts:      0,
        InitialDelay:    1 * time.Second,
        MaxDelay:        60 * time.Second,
        BackoffMultiplier: 1.618, // Golden ratio
    }
}
```

### 5.4 Client Interface and Constructor

```go
// OpenClawClient interface defines the public API
type OpenClawClient interface {
    // Connection management
    Connect(ctx context.Context) error
    Disconnect() error
    State() ConnectionState

    // Request/Response
    SendRequest(ctx context.Context, req *RequestFrame) (*ResponseFrame, error)

    // Event subscription
    Events() <-chan Event
    Subscribe(eventType EventType, handler EventHandler) func()

    // Lifecycle
    Close() error
}

// openClawClient is the concrete implementation
type openClawClient struct {
    config   *ClientConfig
    state    *ConnectionStateMachine
    // ... managers
}

// NewClient returns the interface for dependency injection and testing
func NewClient(opts ...ClientOption) (OpenClawClient, error) {
    cfg := DefaultClientConfig()
    for _, opt := range opts {
        if err := opt(cfg); err != nil {
            return nil, err
        }
    }
    // ... validation and construction
    return &openClawClient{config: cfg, ...}, nil
}
```

---

## 6. Module Mapping

| TypeScript Module | Go Module | Notes |
|-------------------|-----------|-------|
| `OpenClawClient` | `client.go` | Option pattern |
| `ConnectionManager` | `managers/connection.go` | Goroutine + channel |
| `RequestManager` | `managers/request.go` | Context support |
| `EventManager` | `managers/event.go` | Channel pub/sub |
| `ReconnectManager` | `managers/reconnect.go` | Fibonacci backoff |
| `WebSocketTransport` | `transport/websocket.go` | gorilla/websocket |
| `ProtocolNegotiator` | `connection/protocol.go` | - |
| `ConnectionStateMachine` | `connection/state.go` | - |
| `PolicyManager` | `connection/policies.go` | - |
| `TlsValidator` | `connection/tls.go` | - |
| `TickMonitor` | `events/tick.go` | Ticker channel |
| `GapDetector` | `events/gap.go` | - |
| `CredentialsProvider` | `auth/provider.go` | Interface |
| `AuthHandler` | `auth/handler.go` | - |
| `TimeoutManager` | `utils/timeout.go` | - |
| Error classes | `errors.go` | - |

---

## 7. Testing Strategy

- **Unit Tests**: Each module has corresponding `*_test.go` files
- **Integration Tests**: WebSocket connection and protocol tests
- **Examples**: CLI and server examples as end-to-end validation
- **Coverage Target**: 80%+ test coverage

---

## 8. Dependencies

Minimal external dependencies:

- `github.com/gorilla/websocket` - Industry standard WebSocket library
- `golang.org/x/net` - Network utilities (if needed, e.g., for additional protocols)

> Note: Go's standard library `net/http` does not include WebSocket support. The `gorilla/websocket` package is the de facto standard for Go WebSocket implementations.

---

## 9. Backward Compatibility Notes

While the Go SDK follows Go idioms rather than TypeScript API, it maintains **functional equivalence**:

- All features present in TypeScript SDK
- Same protocol wire format
- Same authentication flow
- Same reconnection behavior (fibonacci backoff)
- Same event types and semantics

Users migrating from TypeScript will find equivalent functionality with Go-idiomatic APIs.
