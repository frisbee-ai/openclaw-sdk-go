# Architecture Research

**Domain:** WebSocket Client SDK Library (Go)
**Project:** OpenClaw SDK Go
**Researched:** 2026-03-28
**Confidence:** MEDIUM-HIGH

## Standard Architecture

### System Overview

Mature WebSocket SDKs follow a layered architecture pattern with clear component boundaries:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           SDK Entry Layer (pkg/client.go)                     │
│  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐           │
│  │ ChatAPI │  │AgentsAPI│  │Sessions │  │  ...    │  │UsageAPI │           │
│  └────┬────┘  └────┬────┘  └────┬────┘  └────┬────┘  └────┬────┘           │
├───────┴────────────┴────────────┴────────────┴────────────┴─────────────────┤
│                         Manager Coordination Layer                            │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐            │
│  │  EventManager   │  │  RequestManager  │  │ConnectionManager│            │
│  │  (pub/sub)      │  │  (req/resp corr) │  │ (state+transport│            │
│  └────────┬────────┘  └────────┬────────┘  └────────┬────────┘            │
├───────────┴────────────────────┴────────────────────┴───────────────────────┤
│                        Protocol / Connection Layer                           │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐            │
│  │ ConnectionState │  │   PolicyManager  │  │ProtocolNegotiator│            │
│  │   Machine       │  │   (server caps)   │  │  (version neg)  │            │
│  └────────┬────────┘  └────────┬────────┘  └────────┬────────┘            │
├───────────┴────────────────────┴────────────────────┴───────────────────────┤
│                           Transport Layer (pkg/transport/)                     │
│  ┌──────────────────────────────────────────────────────────────────────┐  │
│  │                    WebSocketTransport                                  │  │
│  │  ┌──────────┐   ┌──────────┐   ┌──────────┐   ┌──────────┐           │  │
│  │  │ writeLoop│   │ readLoop │   │  sendCh  │   │  recvCh  │           │  │
│  │  └──────────┘   └──────────┘   └──────────┘   └──────────┘           │  │
│  └──────────────────────────────────────────────────────────────────────┘  │
│  ┌──────────────────────────────────────────────────────────────────────┐  │
│  │                    gorilla/websocket (underlying library)              │  │
│  └──────────────────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Component Responsibilities

| Component | Responsibility | Typical Implementation | OpenClaw Status |
|-----------|----------------|------------------------|-----------------|
| SDK Entry | Factory, client coordination, type re-exports | Singleton factory, manager wiring | OK - but `client` struct is too large |
| API Namespaces | Typed method wrappers per domain | Each namespace = separate file | Good - 15 namespaces in pkg/api/ |
| EventManager | Pub/sub dispatch, handler registration | Buffered channel + handler map | Good - proper close order |
| RequestManager | Request/response correlation by ID | Pending request map with channels | Good - sync.Cond notification |
| ConnectionManager | WebSocket lifecycle, state transitions | State machine + transport wrapper | Good - but Close/Disconnect confusion |
| ReconnectManager | Backoff reconnection logic | Fibonacci/exponential backoff | OK - loosely coupled via callbacks |
| Transport | Low-level WebSocket I/O | Channel-based read/write loops | Good - interface allows mocking |
| StateMachine | Valid connection state transitions | Map of valid transitions | Good - enforced at compile time |
| PolicyManager | Server capability caching | Cached struct from handshake | Good - prevents client-side policy bypass |
| ProtocolNegotiator | Version negotiation | Range intersection | Good - extensible for future versions |

## Current Project Structure

```
openclaw-sdk-go/
├── pkg/
│   ├── client.go           # Entry point: NewClient(), OpenClawClient interface
│   ├── api/                # 15 API namespace wrappers
│   │   ├── chat.go, agents.go, sessions.go, ...
│   │   └── api.go          # RequestFn type alias from protocol
│   ├── managers/           # Manager implementations
│   │   ├── event.go        # EventManager (pub/sub)
│   │   ├── request.go      # RequestManager (correlation)
│   │   ├── connection.go   # ConnectionManager (lifecycle)
│   │   ├── reconnect.go    # ReconnectManager (backoff)
│   │   └── interfaces.go   # Manager interface definitions
│   ├── transport/          # WebSocket I/O abstraction
│   │   └── websocket.go    # WebSocketTransport + Transport interface
│   ├── connection/         # Connection state and policies
│   │   ├── state.go        # ConnectionStateMachine
│   │   ├── connection_types.go  # ConnectParams, HelloOk, Snapshot
│   │   ├── policies.go     # PolicyManager
│   │   ├── protocol.go     # ProtocolNegotiator
│   │   └── tls.go          # TlsValidator
│   ├── protocol/           # Wire format and API signatures
│   │   ├── types.go        # RequestFrame, ResponseFrame, EventFrame
│   │   ├── validation.go   # Frame validation helpers
│   │   └── api_*.go        # API method signatures (127 methods)
│   ├── events/             # Health monitoring
│   │   └── tick.go         # TickMonitor, GapDetector
│   ├── auth/               # Authentication interfaces
│   │   ├── handler.go      # AuthHandler interface
│   │   └── provider.go     # CredentialsProvider interface
│   ├── types/              # Shared types (errors, logging, events)
│   │   ├── types.go        # ConnectionState, EventType, Event
│   │   ├── errors.go       # Error types with codes
│   │   └── logger.go       # Logger interface
│   └── utils/
│       └── timeout.go      # TimeoutManager
```

### Structure Rationale

- **`pkg/client.go`**: Re-exports types for single-package import convenience. Coordinates all managers but is getting large (370+ lines).
- **`pkg/api/`**: Domain-separated wrappers keep API methods organized. Each namespace is independent.
- **`pkg/managers/`**: Four managers with single responsibilities. Interface definitions in `interfaces.go` allow clean testing with mocks.
- **`pkg/transport/`**: Minimal interface (`Transport`) abstracts gorilla/websocket. Enables unit testing without real connections.
- **`pkg/connection/`**: Connection-specific logic (state machine, policies, TLS). Separated from transport to isolate WebSocket details.
- **`pkg/protocol/`**: Wire format and API types. No dependencies on other SDK layers - pure data definitions.

## Architectural Patterns

### Pattern 1: Manager Coordination via Interface

**What:** Each manager implements an interface defined in `pkg/managers/interfaces.go`. The `client` struct holds manager instances but not interfaces (except for `EventEmitter`).

**When to use:** When you need to test managers in isolation or swap implementations.

**Trade-offs:**
- Pro: Clear contracts, testable with mocks
- Con: Interface definitions add indirection

**Example from current code:**
```go
// pkg/managers/interfaces.go
type ConnectionManagerInterface interface {
    Connect(ctx context.Context) error
    Disconnect() error
    State() types.ConnectionState
    Transport() transport.Transport
    Close() error
}

// pkg/client.go uses concrete types, not interfaces
type client struct {
    managers struct {
        connection *managers.ConnectionManager  // concrete, not interface
    }
}
```

**Recommendation:** The current concrete-type approach is fine for internal use. Exposing interfaces is valuable for external extension points.

### Pattern 2: Channel-Based Event Dispatch

**What:** Events flow through a buffered channel to a dispatch goroutine that fans out to handlers.

**When to use:** When you need to decouple event producers from handlers and prevent deadlocks.

**Trade-offs:**
- Pro: Non-blocking event emission, goroutine-per-handler prevents slow handlers blocking others
- Con: Channel backpressure can drop events (mitigated with emitTimeout)

**Current implementation (good):**
```go
// pkg/managers/event.go - EventManager
func (em *EventManager) Emit(event types.Event) {
    // ... timeout logic ...
    select {
    case em.events <- event:
    case <-em.ctx.Done():
    case <-em.emitTimer.C:
        em.logger.Warn("event channel full, dropping event", "type", event.Type)
    }
}

func (em *EventManager) dispatch(event types.Event) {
    em.mu.RLock()
    // Copy handlers while holding read lock
    handlers := make([]types.EventHandler, 0, len(handlerMap))
    for _, handler := range handlerMap {
        handlers = append(handlers, handler)
    }
    em.mu.RUnlock()
    // Call handlers without holding lock
    for _, handler := range handlers {
        go handler(event)  // fan-out
    }
}
```

### Pattern 3: Request/Response Correlation via Pending Map

**What:** Outbound requests are stored in a map keyed by RequestID. Inbound responses look up the pending request and signal via channel or context.

**When to use:** For request/response protocols over asynchronous transports.

**Trade-offs:**
- Pro: Natural fit for WebSocket's async message pattern
- Con: Memory growth if responses never arrive (mitigated by context timeout and Close cleanup)

**Current implementation (good):**
```go
// pkg/managers/request.go
type pendingRequest struct {
    frame   *protocol.RequestFrame
    respCh  chan *protocol.ResponseFrame
    ctx     context.Context
    cancel  context.CancelFunc
}

func (rm *RequestManager) SendRequest(ctx context.Context, req *protocol.RequestFrame, sendFunc func(*protocol.RequestFrame) error) (*protocol.ResponseFrame, error) {
    // ... create pendingRequest, store in map ...
    // ... send via sendFunc ...
    // ... wait on respCh or context cancellation ...
}
```

### Pattern 4: State Machine for Connection Lifecycle

**What:** Connection state transitions are validated against a map of allowed transitions.

**When to use:** When state transitions must follow a strict protocol (e.g., must authenticate before sending requests).

**Trade-offs:**
- Pro: Prevents invalid state sequences, clear error messages
- Con: Adds complexity for simple use cases

**Current implementation (good):**
```go
// pkg/connection/state.go
validTransitions := map[ConnectionState][]ConnectionState{
    StateDisconnected:   {StateConnecting},
    StateConnecting:     {StateConnected, StateFailed},
    StateConnected:      {StateAuthenticating, StateDisconnected},
    StateAuthenticating: {StateAuthenticated, StateFailed},
    StateAuthenticated:  {StateDisconnected, StateReconnecting},
    StateReconnecting:   {StateConnecting, StateFailed},
    StateFailed:         {StateDisconnected},
}
```

## Data Flow

### Connection Lifecycle

```
client.Connect(ctx)
    |
    v
ConnectionManager.ConnectWithParams(ctx, connectParams)
    |
    +---> transport.Dial() --> WebSocket connection established
    |
    +---> performHandshake() --> Send ConnectParams, wait for HelloOk
    |
    v
processServerInfo()
    |
    +---> ProtocolNegotiator.Negotiate() --> version agreement
    +---> PolicyManager.SetPolicies() --> cache server capabilities
    +---> Initialize TickMonitor (if configured)
    +---> Initialize GapDetector (if configured)
    |
    v
ConnectionManager transitions: Disconnected -> Connecting -> Connected -> Authenticating -> Authenticated
```

### Request/Response Flow

```
client.Chat().Send(ctx, params)
    |
    v
ChatAPI.Send() calls RequestFn closure
    |
    +---> Generate RequestID
    +---> Marshal to RequestFrame JSON
    +---> client.SendRequest(ctx, req)
    |
    v
RequestManager.SendRequest()
    |
    +---> Validate payload size against PolicyManager
    +---> Store pending request in map with ResponseFrame channel
    +---> Transport.Send() --> writes to WebSocket
    |
    v
[Response arrives via WebSocket readLoop]
    |
    v
ConnectionManager dispatches: is ResponseFrame? --> RequestManager.HandleResponse()
    |
    +---> Lookup pending request by RequestID
    +---> Send response to respCh channel
    |
    v
RequestManager.SendRequest() returns response to ChatAPI
```

### Event Flow

```
[Server sends EventFrame via WebSocket]
    |
    v
WebSocketTransport.readLoop() receives message
    |
    v
ConnectionManager receives from transport.Receive() channel
    |
    +---> If EventFrame: EventManager.Emit(event)
    +---> If ResponseFrame: RequestManager.HandleResponse()
    |
    v
EventManager.dispatch() fans out to handlers in goroutines
```

## Component Dependency Graph

```
pkg/client.go (entry)
    |
    +---> pkg/managers/event.go (EventManager)
    |         |
    |         +---> pkg/types/*.go
    |
    +---> pkg/managers/request.go (RequestManager)
    |         |
    |         +---> pkg/protocol/types.go
    |         +---> pkg/utils/timeout.go
    |
    +---> pkg/managers/connection.go (ConnectionManager)
    |         |
    |         +---> pkg/transport/websocket.go (Transport interface)
    |         +---> pkg/connection/state.go
    |         +---> pkg/connection/connection_types.go
    |         +---> pkg/managers/event.go (for emitting connect/disconnect events)
    |
    +---> pkg/managers/reconnect.go (ReconnectManager)
    |         |
    |         +---> pkg/managers/connection.go (via callback)
    |
    +---> pkg/connection/policies.go (PolicyManager)
    +---> pkg/connection/protocol.go (ProtocolNegotiator)
    +---> pkg/events/tick.go (TickMonitor, GapDetector)
    |
    +---> pkg/api/*.go (15 API namespaces)
              |
              +---> pkg/protocol/*.go (types + api_*)
```

## Recommended Build Order

For understanding or testing the SDK, follow this dependency order:

1. **pkg/types/** -- No dependencies. Shared constants and interfaces.
2. **pkg/protocol/** -- No dependencies. Wire format definitions.
3. **pkg/auth/** -- No dependencies. Authentication interfaces.
4. **pkg/utils/** -- No dependencies. Utility functions.
5. **pkg/connection/state.go** -- Depends on types. State machine logic.
6. **pkg/connection/** (other files) -- Depends on state, types.
7. **pkg/transport/websocket.go** -- Depends on connection (for TLS). Implements Transport interface.
8. **pkg/managers/event.go** -- Depends on types. Event dispatch.
9. **pkg/managers/request.go** -- Depends on protocol, utils.
10. **pkg/managers/connection.go** -- Depends on transport, connection, event.
11. **pkg/managers/reconnect.go** -- Depends on nothing new.
12. **pkg/events/tick.go** -- Depends on types.
13. **pkg/api/*.go** -- Depends on protocol.
14. **pkg/client.go** -- Wires everything together.

## Scaling Considerations

| Scale | Architecture Adjustments |
|-------|-------------------------|
| 0-100 clients | Current architecture fine. Single connection per client. |
| 100-1000 clients | Connection pooling at application level (SDK doesn't manage). Consider connection-per-request patterns. |
| 1000-10000 clients | SDK unchanged. User manages connection fan-out. |
| 10k+ clients | Consider shard-based routing. SDK could add connection multiplexing. |

### Scaling Priorities for OpenClaw SDK

1. **First bottleneck:** Event channel backpressure -- events may drop under high load. Consider larger default buffer or overflow to disk.
2. **Second bottleneck:** Pending request map memory -- long-running requests could accumulate. Consider LRU eviction.
3. **Third bottleneck:** JSON marshaling overhead -- hot paths serialize/deserialize repeatedly. Consider `fastjson` or `sonic` for 2-3x speedup.

## Anti-Patterns

### Anti-Pattern 1: God Client Struct

**What people do:** Stuff everything into a single `client` struct -- all managers, all state, all helpers.

**Why it's wrong:** Makes the code hard to test in isolation, increases cognitive load, creates circular dependency risk.

**Do this instead:** Use composition with clear boundaries. The current `managers` struct is a good start but `client` still has direct fields for `protocolNegotiator`, `policyManager`, `tickMonitor`, `gapDetector`, and all 15 API namespaces.

**Current code showing the problem:**
```go
type client struct {
    config   *ClientConfig
    managers struct { ... }           // good: grouped
    protocolNegotiator *connection.ProtocolNegotiator  // direct field
    policyManager      *connection.PolicyManager       // direct field
    tickMonitor        *events.TickMonitor            // direct field
    gapDetector        *events.GapDetector            // direct field
    // ... 15 API namespaces as direct fields
}
```

**Recommendation:** Group these into logical sub-structs:
```go
type client struct {
    config *ClientConfig
    core struct {
        event      *managers.EventManager
        request    *managers.RequestManager
        connection *managers.ConnectionManager
        reconnect  *managers.ReconnectManager
    }
    protocol struct {
        negotiator *connection.ProtocolNegotiator
        policies   *connection.PolicyManager
    }
    health struct {
        tick *events.TickMonitor
        gap  *events.GapDetector
    }
    api struct {
        chat   *api.ChatAPI
        agents *api.AgentsAPI
        // ...
    }
}
```

### Anti-Pattern 2: Close/Disconnect Confusion

**What people do:** Have both `Close()` and `Disconnect()` methods with unclear semantics.

**Why it's wrong:** Users don't know which to call. Confusion leads to resource leaks or double-cleanup.

**Current situation:** `ConnectionManager` has both `Close()` and `Disconnect()`:
- `Close()` delegates to `Disconnect()`
- `Disconnect()` closes transport and transitions to `StateDisconnected`

**Do this instead:** Single method with clear semantics:
- `Disconnect()` = close connection but keep manager alive for reconnect
- `Close()` = full shutdown, manager cannot be reused

Or remove `Close()` entirely and have `Disconnect()` be the only cleanup method.

### Anti-Pattern 3: Channel + Lock Violation

**What people do:** Send to a channel while holding a mutex lock.

**Why it's wrong:** Can cause死锁 (deadlock) if receiver tries to acquire the same lock.

**Current code (correct):**
```go
func (cm *ConnectionManager) Disconnect() error {
    cm.mu.Lock()
    defer cm.mu.Unlock()
    // ... lock held only for state check and transport access ...
    err := cm.transport.Close()  // transport.Close() doesn't need cm's lock
    // ... state transition without holding lock in emit ...
    if cm.eventMgr != nil {
        // Note: Emit doesn't acquire cm's lock, so no deadlock here
        cm.eventMgr.Emit(...)
    }
}
```

**Rule:** Release lock before sending to channels that might callback into the locked struct.

### Anti-Pattern 4: Context Cancellation Without Timeout

**What people do:** Pass `context.Background()` or unbounded context to `Connect()`.

**Why it's wrong:** DNS lookup, TCP connection, and WebSocket handshake can hang indefinitely.

**Current code:** `Connect(ctx context.Context)` is passed through. No internal timeout.

**Do this instead:** Document that callers should provide a timeout context:
```go
// NewClient should document:
// "ctx is used for client lifecycle. For Connect, provide a context with timeout."
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
```

Or add internal timeouts in `Connect()` that wrap the passed context.

## Integration Points

### External: OpenClaw Gateway

| Aspect | Integration Pattern | Notes |
|--------|---------------------|-------|
| Protocol | WebSocket + JSON frames | Stable - version negotiated at handshake |
| Authentication | ConnectParams + HelloOk exchange | AuthHandler abstracts credential provision |
| Heartbeat | Tick events (server-sent) | Client monitors via TickMonitor |
| Error handling | ErrorFrame with code + message | Retryable flag for transient errors |

### Internal Boundaries

| Boundary | Communication | Notes |
|----------|---------------|-------|
| client -> managers | Direct method calls | Mutex-protected for thread safety |
| managers -> transport | Transport interface | Allows mock in tests |
| transport -> WebSocket | gorilla/websocket API | Only external dependency |
| connection -> event | EventManager.Emit() | Non-blocking channel send |

## Sources

- Current codebase analysis (2026-03-28)
- gorilla/websocket package documentation
- Go concurrency patterns (channel + mutex hybrid)
- Standard state machine patterns for connection protocols

---

*Architecture research for: WebSocket Client SDK*
*Researched: 2026-03-28*
