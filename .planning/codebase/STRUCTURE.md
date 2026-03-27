# Codebase Structure

**Analysis Date:** 2026-03-28

## Directory Layout

```
openclaw-sdk-go/
├── CLAUDE.md               # Project guidance for Claude Code
├── README.md               # Public documentation
├── CHANGELOG.md            # Release history
├── LICENSE                 # Apache 2.0
├── go.mod                  # Go module definition
├── go.sum                  # Dependency checksums
├── coverage.out            # Test coverage output
├── .golangci.yaml          # Linter configuration
├── .pre-commit-config.yaml # Pre-commit hook definitions
├── .goreleaser.yaml        # Release build configuration
├── .github/                # GitHub workflows and templates
├── docs/                   # Specifications and plans
│   ├── specs/             # Architecture and design specs
│   └── plans/             # Phase implementation plans
├── examples/               # Usage examples
├── pkg/                    # Main source code
└── .planning/codebase/    # GSD planning documents
```

## Directory Purposes

**`pkg/` - Main Source Code:**
- Purpose: All production code for the SDK
- Contains: Go packages organized by layer/responsibility

**`pkg/api/` - API Namespace Clients:**
- Purpose: Typed wrappers for server API methods
- Contains: `api.go`, `chat.go`, `agents.go`, `sessions.go`, `config.go`, `cron.go`, `nodes.go`, `skills.go`, `device_pairing.go`, `browser.go`, `channels.go`, `push.go`, `exec_approvals.go`, `system.go`, `secrets.go`, `usage.go`, `shared.go`
- Key files: `pkg/api/api.go` (type aliases re-exported from protocol), `pkg/api/chat.go` (API pattern example)

**`pkg/auth/` - Authentication:**
- Purpose: Auth handler and credentials provider abstractions
- Contains: `handler.go`, `provider.go`
- Key files: `pkg/auth/handler.go` (`AuthHandler` interface, `StaticAuthHandler`)

**`pkg/connection/` - Connection Management:**
- Purpose: State machine, handshake protocol, TLS, policies
- Contains: `state.go`, `connection_types.go`, `protocol.go`, `policies.go`, `tls.go`, `policies_test.go`, `protocol_test.go`, `state_test.go`, `tls_test.go`
- Key files: `pkg/connection/state.go` (`ConnectionStateMachine`), `pkg/connection/connection_types.go` (`ConnectParams`, `HelloOk`, `Snapshot`)

**`pkg/events/` - Event Monitors:**
- Purpose: Tick/heartbeat monitoring and gap detection
- Contains: `tick.go`, `tick_test.go`, `gap.go` (if present)
- Key files: `pkg/events/tick.go` (`TickMonitor`)

**`pkg/managers/` - Manager Components:**
- Purpose: High-level coordination (event, request, connection, reconnect)
- Contains: `event.go`, `request.go`, `connection.go`, `reconnect.go`, `interfaces.go`, `*_test.go`
- Key files: `pkg/managers/event.go` (`EventManager`), `pkg/managers/request.go` (`RequestManager`)

**`pkg/protocol/` - Protocol Definition:**
- Purpose: Wire format frames and API method signatures
- Contains: `types.go`, `validation.go`, `errors.go`, `api_*.go` (16 API modules), `wire_*.go`, `fuzz_test.go`
- Key files: `pkg/protocol/types.go` (`RequestFrame`, `ResponseFrame`, `EventFrame`)

**`pkg/transport/` - WebSocket Transport:**
- Purpose: Low-level WebSocket I/O
- Contains: `websocket.go`, `websocket_test.go`
- Key files: `pkg/transport/websocket.go` (`WebSocketTransport`, `Transport` interface)

**`pkg/types/` - Shared Types:**
- Purpose: Core types (states, events, errors, logging)
- Contains: `types.go`, `errors.go`, `logger.go`, `*_test.go`
- Key files: `pkg/types/types.go` (`ConnectionState`, `EventType`, `Event`)

**`pkg/utils/` - Utilities:**
- Purpose: Timeout management helpers
- Contains: `timeout.go`, `timeout_test.go`

**`pkg/client.go` - SDK Entry Point:**
- Purpose: Client factory, manager coordination, type re-exports
- Re-exports types from subpackages for single-package import convenience

## Key File Locations

**Entry Points:**
- `pkg/client.go`: Main `NewClient()` factory, `OpenClawClient` interface, all API accessors

**Configuration:**
- `pkg/client.go`: `ClientConfig` struct, `ClientOption` functional options (`WithURL`, `WithAuthHandler`, etc.)
- `pkg/connection/connection_types.go`: `ConnectParams` for handshake

**Core Logic:**
- `pkg/managers/connection.go`: WebSocket lifecycle, handshake, state transitions
- `pkg/managers/event.go`: Pub/sub event dispatch
- `pkg/managers/request.go`: Request/response correlation
- `pkg/managers/reconnect.go`: Fibonacci backoff reconnection
- `pkg/transport/websocket.go`: Low-level WebSocket send/receive loops
- `pkg/connection/state.go`: State machine with valid transitions map

**Testing:**
- `pkg/*/*_test.go`: Co-located test files following Go convention
- `pkg/integration_test.go`: Integration tests

## Naming Conventions

**Files:**
- Go source: `lowercase.go` (e.g., `client.go`, `event.go`)
- Test files: `*_test.go` (e.g., `event_test.go`)
- API modules: `api_*.go` (e.g., `api_chat.go`, `api_agent.go`)

**Packages:**
- Short, lowercase: `transport`, `managers`, `protocol`, `connection`
- No `utils` or `helpers` - functionality grouped by domain

**Types:**
- PascalCase: `ConnectionStateMachine`, `EventManager`, `WebSocketTransport`
- Interfaces often noun or noun+er: `Transport`, `AuthHandler`, `CredentialsProvider`
- Error types: `<Domain>Error` (e.g., `OpenClawError`, `ConnectionError`)

**Functions:**
- CamelCase: `NewClient`, `Connect`, `SendRequest`, `Subscribe`
- Constructor: `New<StructName>` (e.g., `NewEventManager`, `NewTickMonitor`)
- Option: `With<Option>` (e.g., `WithURL`, `WithLogger`)

**Variables:**
- CamelCase: `eventMgr`, `requestFn`, `tickMonitor`
- Short for receivers: `c` (client), `tm` (tick monitor), `em` (event manager)
- `mu` for mutex, `wg` for WaitGroup

**Constants:**
- CamelCase for state/event constants: `StateConnected`, `EventTick`
- All caps for exported error codes: `ProtocolErrFrameTooLarge`

## Where to Add New Code

**New API Namespace:**
- Implementation: `pkg/protocol/api_<name>.go` (add types) + `pkg/api/<name>.go` (add client)
- Tests: `pkg/protocol/<name>_test.go` (if any) + `pkg/api/<name>_test.go` (if any)
- Register accessor in `pkg/client.go` if namespace access pattern needed

**New Protocol Type:**
- Add to `pkg/protocol/types.go` or new `pkg/protocol/api_<domain>.go`
- Add validation in `pkg/protocol/validation.go` if needed
- Tests in `pkg/protocol/*_test.go`

**New Event Type:**
- Add constant to `pkg/types/types.go` (`Event<Type> EventType = "<name>"`)
- Emit via `EventManager.Emit()` in appropriate manager
- Subscribe via `client.Subscribe(Event<Type>, handler)`

**New Manager:**
- Add to `pkg/managers/` as `<name>.go`
- Implement `Close()` with goroutine cleanup
- Initialize in `NewClient()` and coordinate from `pkg/client.go`

**New Error Type:**
- Add to `pkg/types/errors.go`
- Add error code constant to appropriate type
- Add constructor `New<Error>Error` function

**Utilities:**
- Timeout helpers: `pkg/utils/timeout.go`
- No `helpers.go` or `utils.go` files - organize by purpose

## Special Directories

**`.github/`:**
- Purpose: GitHub Actions workflows and PR templates
- Generated: No
- Committed: Yes

**`docs/`:**
- Purpose: Design documents and migration plans
- Contains: `docs/specs/`, `docs/plans/`
- Generated: No
- Committed: Yes

**`examples/`:**
- Purpose: Usage examples
- Generated: No
- Committed: Yes

**`.planning/codebase/` (GSD planning):**
- Purpose: Architecture and structure analysis documents
- Generated: Yes (by GSD codebase mapper)
- Committed: No (local planning only)

---

*Structure analysis: 2026-03-28*
