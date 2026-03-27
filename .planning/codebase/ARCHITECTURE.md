# Architecture

**Analysis Date:** 2026-03-28

## Pattern Overview

**Overall:** Manager-based WebSocket client with channel-based event dispatch

**Key Characteristics:**
- **Option Pattern**: Client configuration via functional options (`WithURL()`, `WithAuthHandler()`, etc.)
- **Context + Channel Hybrid**: `context.Context` for lifecycle/cancellation, buffered channels for event delivery
- **Manager Coordination**: Four managers (event, request, connection, reconnect) coordinated by a single `client` struct
- **Graceful Shutdown**: All managers implement `Close()` with proper goroutine cleanup via `sync.WaitGroup`

## Layers

**SDK Entry Layer (`pkg/client.go`):**
- Purpose: Main public API, client factory, manager coordination
- Location: `pkg/client.go`
- Contains: `OpenClawClient` interface, `client` struct, `NewClient()`, all API namespace accessors
- Depends on: All subpackages
- Used by: End users

**API Namespace Layer (`pkg/api/`):**
- Purpose: Typed API method wrappers per domain (chat, agents, sessions, etc.)
- Location: `pkg/api/*.go`
- Contains: `ChatAPI`, `AgentsAPI`, `SessionsAPI`, `ConfigAPI`, `CronAPI`, `NodesAPI`, `SkillsAPI`, `DevicePairingAPI`, `BrowserAPI`, `ChannelsAPI`, `PushAPI`, `ExecApprovalsAPI`, `SystemAPI`, `SecretsAPI`, `UsageAPI`
- Depends on: `pkg/protocol` for types, `pkg/api` for `RequestFn`
- Used by: `pkg/client.go`

**Protocol Layer (`pkg/protocol/`):**
- Purpose: Wire protocol frame definitions, validation, and API method signatures
- Location: `pkg/protocol/types.go`, `pkg/protocol/validation.go`, `pkg/protocol/api_*.go`
- Contains: `RequestFrame`, `ResponseFrame`, `EventFrame`, `ErrorShape`, `StateVersion`, API param/result types
- Depends on: Standard library only
- Used by: `pkg/api/`, `pkg/managers/request.go`

**Managers Layer (`pkg/managers/`):**
- Purpose: High-level coordination for events, requests, connections, reconnection
- Location: `pkg/managers/event.go`, `pkg/managers/request.go`, `pkg/managers/connection.go`, `pkg/managers/reconnect.go`
- Contains: `EventManager`, `RequestManager`, `ConnectionManager`, `ReconnectManager`
- Depends on: `pkg/transport`, `pkg/connection`, `pkg/types`, `pkg/events`
- Used by: `pkg/client.go`

**Connection Layer (`pkg/connection/`):**
- Purpose: Connection state machine, handshake types, TLS validation, policy management
- Location: `pkg/connection/state.go`, `pkg/connection/connection_types.go`, `pkg/connection/policies.go`, `pkg/connection/tls.go`
- Contains: `ConnectionStateMachine`, `ConnectParams`, `HelloOk`, `Snapshot`, `Policy`, `TLSConfig`
- Depends on: `pkg/types`
- Used by: `pkg/managers/connection.go`

**Transport Layer (`pkg/transport/`):**
- Purpose: Low-level WebSocket I/O
- Location: `pkg/transport/websocket.go`
- Contains: `WebSocketTransport`, `Transport` interface, `WebSocketConfig`, `TLSConfig`
- Depends on: `github.com/gorilla/websocket`, `pkg/connection` (for TLS config)
- Used by: `pkg/managers/connection.go`

**Events Layer (`pkg/events/`):**
- Purpose: Connection health monitoring and gap detection
- Location: `pkg/events/tick.go`
- Contains: `TickMonitor`, `GapDetector`
- Depends on: Standard library only
- Used by: `pkg/client.go`

**Auth Layer (`pkg/auth/`):**
- Purpose: Authentication handler and credentials provider interfaces
- Location: `pkg/auth/handler.go`, `pkg/auth/provider.go`
- Contains: `AuthHandler`, `CredentialsProvider`, `StaticAuthHandler`, `StaticCredentialsProvider`
- Depends on: Standard library only
- Used by: `pkg/client.go` (via `ClientConfig.CredentialsProvider`)

**Types Layer (`pkg/types/`):**
- Purpose: Shared core types (connection states, event types, errors, logging)
- Location: `pkg/types/types.go`, `pkg/types/errors.go`, `pkg/types/logger.go`
- Contains: `ConnectionState`, `EventType`, `Event`, `EventHandler`, `ReconnectConfig`, error types, `Logger` interface
- Depends on: Standard library only
- Used by: All layers (re-exported from `pkg/client.go`)

**Utils Layer (`pkg/utils/`):**
- Purpose: Utility functions (timeout management)
- Location: `pkg/utils/timeout.go`
- Depends on: Standard library only
- Used by: `pkg/managers/request.go`

## Data Flow

**Outbound Request Flow:**
1. User calls `client.Chat().Send(ctx, params)` (or any API method)
2. `ChatAPI.Send()` calls `RequestFn` closure in `client`
3. `newRequestFn()` creates `protocol.RequestFrame` with generated ID
4. `client.SendRequest()` acquires mutex, validates against `PolicyManager`
5. `RequestManager.SendRequest()` stores pending request in map, sends via `Transport.Send()`
6. `WebSocketTransport.writeLoop()` writes JSON to WebSocket

**Inbound Response Flow:**
1. `WebSocketTransport.readLoop()` reads message from WebSocket
2. `ConnectionManager` dispatches to `EventManager` or handles as response
3. `RequestManager` correlates response by RequestID, delivers to waiting goroutine
4. User's `RequestFn` returns parsed payload to API method

**Event Flow:**
1. `WebSocketTransport.readLoop()` receives message, detects `EventFrame` type
2. `ConnectionManager` emits event via `EventManager.Emit()`
3. `EventManager.dispatch()` delivers to subscribed handlers in goroutines

**Connection Lifecycle Flow:**
1. `client.Connect(ctx)` acquires mutex, calls `ConnectionManager.ConnectWithParams()`
2. `ConnectionManager.Connect()` dials WebSocket, transitions to `StateConnected`
3. `performHandshake()` sends `ConnectParams`, waits for `HelloOk`
4. `processServerInfo()` stores server info, policies, negotiates protocol
5. State transitions: `Disconnected` -> `Connecting` -> `Connected` -> `Authenticating` -> `Authenticated`

## Key Abstractions

**Transport Interface (`pkg/transport/websocket.go`):**
- Purpose: Abstract WebSocket operations
- Examples: `WebSocketTransport` implementation
- Pattern: Interface with `Send()`, `Receive()`, `Errors()`, `Close()`, `IsConnected()`

**State Machine (`pkg/connection/state.go`):**
- Purpose: Enforce valid connection state transitions
- Examples: `ConnectionStateMachine`
- Pattern: `sync.RWMutex` protecting state, buffered channel for state change events, `validTransitions` map

**Request Manager (`pkg/managers/request.go`):**
- Purpose: Correlate request/response by RequestID
- Examples: Pending request map with context channels
- Pattern: Map of RequestID -> pendingRequest with `sync.Cond` for notification

**Event Manager (`pkg/managers/event.go`):**
- Purpose: Pub/sub event dispatch
- Examples: Handler map by EventType
- Pattern: Buffered channel for events, handler registration with auto-increment keys, atomic ID generation

## Entry Points

**Public SDK Entry (`pkg/client.go`):**
- Location: `pkg/client.go`
- Triggers: User instantiates via `NewClient(opts...)`, then calls `Connect(ctx)`, `SendRequest(ctx, req)`, `Subscribe(eventType, handler)`
- Responsibilities: Client lifecycle, manager coordination, API namespace provision, re-export types for convenience

**WebSocket Transport (`pkg/transport/websocket.go`):**
- Location: `pkg/transport/websocket.go`
- Triggers: Called by `ConnectionManager`
- Responsibilities: WebSocket dial, read/write goroutines, ping/pong, close handling

**API Namespaces (`pkg/api/*.go`):**
- Location: `pkg/api/chat.go`, `pkg/api/agents.go`, etc.
- Triggers: User calls `client.<Namespace>().<Method>(ctx, params)`
- Responsibilities: Method-specific request construction, response parsing

## Error Handling

**Strategy:** Error types with codes, wrapped errors using `errors.Join`, context propagation

**Patterns:**
- Custom error types with `ErrorCode` string constants (`pkg/types/errors.go`)
- Error creation via `New<Type>Error(code, message, retryable, cause)` constructors
- Errors propagated through context and channels, not thrown
- `RequestManager.SendRequest()` returns protocol errors as Go errors
- `EventManager.Emit()` drops events after timeout without failing

## Cross-Cutting Concerns

**Logging:** `Logger` interface in `pkg/types/logger.go` with `DefaultLogger`, `NopLogger` implementations. Injected via `ClientConfig.Logger` or context via `WithContext`/`FromContext`.

**Validation:** `pkg/protocol/validation.go` provides frame validation helpers. Payload size validated against `PolicyManager.GetMaxPayload()` before send.

**Authentication:** `AuthHandler` interface + `CredentialsProvider` interface. `StaticAuthHandler` and `StaticCredentialsProvider` provided as defaults. Passed via `WithAuthHandler()` or `WithCredentialsProvider()` options.

**TLS:** `pkg/connection/tls.go` provides `TlsValidator` for certificate loading. `pkg/transport/websocket.go` converts to `crypto/tls.Config`.

---

*Architecture analysis: 2026-03-28*
