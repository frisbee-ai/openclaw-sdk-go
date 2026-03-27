# Testing Patterns

**Analysis Date:** 2026-03-28

## Test Framework

**Runner:**
- Standard Go testing: `testing` package
- No external test framework (no testify, no ginkgo)

**Run Commands:**
```bash
go test ./...                    # Run all tests
go test -cover ./...            # With coverage
go test -race ./...              # Race detection
go test -run TestName ./pkg/...  # Run specific test
go test -v ./pkg/...             # Verbose output
```

**Assertion Library:**
- Standard Go `if` statements and `t.Errorf`/`t.Fatalf`
- No external assertion library

**Pre-commit Hook:**
- `go test ./...` runs all tests before commit

## Test File Organization

**Location:**
- Co-located with source: `pkg/managers/event_test.go` alongside `pkg/managers/event.go`

**Naming:**
- Pattern: `*_test.go`
- Test functions: `TestSubjectName_Scenario` or `TestSubjectName`

**Structure:**
```
pkg/
├── managers/
│   ├── event.go           # Implementation
│   ├── event_test.go      # Unit tests
│   ├── request.go
│   ├── request_test.go
```

## Test Structure

**Basic Pattern:**
```go
func TestSubjectName_ExpectedBehavior(t *testing.T) {
    // Setup
    ctx := context.Background()
    em := NewEventManager(ctx, 10, 20*time.Millisecond)

    // Action
    em.Start()
    em.Emit(types.Event{Type: types.EventConnect, Timestamp: time.Now()})

    // Assertion with synchronization
    timeout := time.After(100 * time.Millisecond)
    done := make(chan struct{})
    go func() {
        // Wait for condition
        close(done)
    }()

    select {
    case <-timeout:
        t.Error("timeout waiting for handler")
    case <-done:
        // Success
    }

    _ = em.Close()
}
```

## Table-Driven Tests

**Preferred for Multiple Test Cases:**
```go
func TestNewAPIError_AuthErrors(t *testing.T) {
    tests := []struct {
        name  string
        shape *ErrorShape
    }{
        {
            name: "AUTH_TOKEN_EXPIRED",
            shape: &ErrorShape{
                Code:    "AUTH_TOKEN_EXPIRED",
                Message: "Token has expired",
            },
        },
        // ... more cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := NewAPIError(tt.shape)
            if !IsAuthError(err) {
                t.Errorf("IsAuthError() = false, want true for code %s", tt.shape.Code)
            }
        })
    }
}
```

## Synchronization Patterns

**Avoid `time.Sleep` for synchronization.** Use proper synchronization:

**WaitGroup for Multiple Operations:**
```go
var wg sync.WaitGroup
wg.Add(2)
go func() {
    defer wg.Done()
    // work
}()
go func() {
    defer wg.Done()
    // work
}()
wg.Wait()
```

**Channel-Based Coordination:**
```go
done := make(chan struct{})
go func() {
    wg.Wait()
    close(done)
}()
select {
case <-time.After(timeout):
    t.Error("timeout")
case <-done:
    // completed
}
```

**Atomic Operations:**
```go
var counter atomic.Int64
counter.Add(1)
if counter.Load() == 0 { /* ... */ }
```

## Mocking

**No External Mocking Library.** Use direct instantiation or test doubles:

**Interface Verification:**
```go
// Verify implementation satisfies interface
var _ Transport = (*WebSocketTransport)(nil)
```

**Test Structs:**
```go
// Direct struct creation for unit testing
transport := &WebSocketTransport{
    sendCh: make(chan []byte, 1),
    recvCh: make(chan []byte, 1),
    errCh:  make(chan error, 1),
    ctx:    ctx,
    cancel: cancel,
}
```

**Fake Servers for Integration Tests:**
```go
// From pkg/transport/websocket_test.go
type testServer struct {
    URL     string
    server  *http.Server
    done    chan struct{}
}

func newTestServer(t *testing.T, handler http.HandlerFunc) *testServer {
    listener, err := net.Listen("tcp", "127.0.0.1:0")
    // ... setup
    return &testServer{...}
}

func (s *testServer) Close() {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    _ = s.server.Shutdown(ctx)
    <-s.done
}
```

## Fixtures and Test Data

**Inline Test Data:**
```go
creds := map[string]string{
    "username": "testuser",
    "password": "testpass",
}
handler, err := NewStaticAuthHandler(creds)
```

**Protocol Frames:**
```go
req := protocol.NewRequestFrame("test-123", "test", nil)
resp := protocol.NewResponseFrameSuccess("test-123", json.RawMessage(`{}`))
```

## Race Detection

**Run with `-race` flag:**
```bash
go test -race ./...
```

**Tests Using Race Flag Explicitly:**
```go
// TestRequestManager_RaceHandleResponseAndTimeout tests the race condition between
// HandleResponse and cleanup (context timeout). This test uses -race flag to detect.
func TestRequestManager_RaceHandleResponseAndTimeout(t *testing.T) { ... }
```

## Error Path Testing

**Test Both Success and Failure:**
```go
func TestClientConnectWithoutURL(t *testing.T) {
    client, err := NewClient()
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    defer func() { _ = client.Close() }()

    ctx := context.Background()
    err = client.Connect(ctx)
    if err == nil {
        t.Error("expected error when connecting without URL")
    }
}
```

**Table-Driven for Error Types:**
```go
func TestIsRetryable(t *testing.T) {
    tests := []struct {
        name    string
        err     error
        want    bool
    }{
        {
            name:    "retryable AuthError",
            err:     NewAuthError("AUTH_TOKEN_EXPIRED", "Token expired", true, nil),
            want:    true,
        },
        // ... more cases
    }
}
```

## Coverage

**Current Status:** Target appears to be 80%+ based on project guidelines

**View Coverage:**
```bash
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

**Integration Test Coverage:**
```go
// From pkg/integration_test.go - tests client creation, options, subscriptions
func TestClient_NewClient(t *testing.T) { ... }
func TestClient_Options(t *testing.T) { ... }
func TestClient_Subscribe(t *testing.T) { ... }
```

## Common Patterns

**Async Testing:**
```go
func TestEventManager_Subscribe(t *testing.T) {
    em.Start()
    em.Emit(types.Event{Type: types.EventConnect, Timestamp: time.Now()})

    timeout := time.After(100 * time.Millisecond)
    done := make(chan struct{})
    go func() {
        mu.Lock()
        for !handlerCalled {
            mu.Unlock()
            time.Sleep(10 * time.Millisecond)
            mu.Lock()
        }
        mu.Unlock()
        close(done)
    }()

    select {
    case <-timeout:
        t.Error("timeout waiting for handler")
    case <-done:
        // Handler was called
    }
}
```

**Concurrent Testing:**
```go
func TestEventManager_ConcurrentSubscribeUnsubscribe(t *testing.T) {
    var wg sync.WaitGroup
    em.Start()

    for i := 0; i < 50; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            handler := func(e types.Event) {}
            unsubscribe := em.Subscribe(types.EventConnect, handler)
            time.Sleep(time.Microsecond)
            unsubscribe()
        }()
    }

    done := make(chan struct{})
    go func() {
        wg.Wait()
        close(done)
    }()

    select {
    case <-time.After(2 * time.Second):
        t.Error("timeout in concurrent subscribe/unsubscribe")
    case <-done:
    }
}
```

**Idempotency Testing:**
```go
func TestEventManager_CloseIdempotent(t *testing.T) {
    em.Start()

    // Close multiple times - should be idempotent
    err1 := em.Close()
    err2 := em.Close()
    err3 := em.Close()

    if err1 != nil {
        t.Errorf("first Close error: %v", err1)
    }
    if err2 != nil {
        t.Errorf("second Close error: %v", err2)
    }
    if err3 != nil {
        t.Errorf("third Close error: %v", err3)
    }
}
```

**Context Cancellation:**
```go
func TestRequestManager_ContextCancellation(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background())
    cancel() // Cancel immediately

    _, err := rm.SendRequest(ctx, req, nil)
    if err == nil {
        t.Error("expected error for cancelled context")
    }
}
```

## Test Categories

**Unit Tests:**
- Test individual components in isolation
- Mock dependencies via direct instantiation
- Focus: `pkg/types/`, `pkg/managers/`, `pkg/protocol/`

**Integration Tests:**
- Test client creation and options: `pkg/integration_test.go`
- Test client-level functionality without real server

**Transport Tests:**
- Use fake WebSocket server for E2E tests
- Located in `pkg/transport/websocket_test.go`

---

*Testing analysis: 2026-03-28*
