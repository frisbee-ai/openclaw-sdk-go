# OpenClaw SDK Go

## What This Is

OpenClaw SDK Go is a feature-complete WebSocket client library for Go, migrated from TypeScript. It provides a robust SDK for connecting to the OpenClaw gateway with connection management, event handling, request/response patterns, and automatic reconnection. End users are Go applications that need to interact with the OpenClaw platform via WebSocket.

## Core Value

**Go developers can integrate the OpenClaw platform in under 10 lines of code** — the SDK handles connection lifecycle, authentication, protocol framing, event dispatch, and reconnection transparently.

## Requirements

### Validated

- ✓ **WebSocket Transport** — gorilla/websocket-based transport with ping/pong, TLS, graceful close
- ✓ **Connection State Machine** — Disconnected → Connecting → Connected → Authenticating → Authenticated with valid transition enforcement
- ✓ **Protocol Layer** — JSON RequestFrame/ResponseFrame/EventFrame serialization, API param/result types for 15 namespaces
- ✓ **Auth Handler** — AuthHandler interface + StaticAuthHandler, CredentialsProvider interface + StaticCredentialsProvider
- ✓ **Event Manager** — Pub/sub event dispatch via buffered channels, handler registration, auto-increment keys
- ✓ **Request Manager** — Request/response correlation by RequestID, pending request map with context channels
- ✓ **Connection Manager** — WebSocket lifecycle, handshake (ConnectParams → HelloOk), server info negotiation
- ✓ **Reconnect Manager** — Fibonacci backoff reconnection with jitter, configurable via ReconnectConfig
- ✓ **Tick Monitor** — Connection health monitoring via periodic tick events
- ✓ **Gap Detector** — Detect missed events via sequence numbering
- ✓ **All 15 API Namespaces** — Chat, Agents, Sessions, Config, Cron, Nodes, Skills, DevicePairing, Browser, Channels, Push, ExecApprovals, System, Secrets, Usage (127 total API methods)
- ✓ **Option Pattern** — Functional options (WithURL, WithAuthHandler, WithCredentialsProvider, WithLogger, etc.)
- ✓ **TLS Validation** — Custom certificate validation via TlsValidator interface
- ✓ **Context Propagation** — All operations respect context.Context cancellation
- ✓ **Graceful Shutdown** — All managers implement Close() with goroutine cleanup via sync.WaitGroup
- ✓ **CI/CD Pipeline** — GitHub Actions: fmt, vet, lint, test with race detection, coverage upload
- ✓ **Pre-commit Hooks** — Local: gofmt, go vet, golangci-lint, go test
- ✓ **Client Struct Organization** — Sub-struct grouping (managers, api, protocol, health) for maintainability
- ✓ **API Semantics** — Clear Close() vs Disconnect() documentation in OpenClawClient interface
- ✓ **Performance** — L1+L2 hot-path benchmarks (JSON marshal/unmarshal, event dispatch, request correlation); custom metrics (bytes_per_frame, channel_overhead_ns, goroutine_count); fuzz test round-trip assertions with 24 corpus files
- ✓ **Release Process** — GoReleaser v2 configured for library mode (mode: github, gomod.proxy: true); v1.0.0 semantic version tag exists; git-cliff changelog automation integrated with release.yml (REL-01, REL-03 satisfied; REL-02 partial — v1.0.1 tag skipped by user)

### Active

- [ ] **API Stability** — Complete fuzz testing, edge case coverage, API contract verification
- [ ] **Developer Experience** — Comprehensive examples, usage guides, migration documentation
- [ ] **Performance** — Benchmark existing hot paths, identify bottlenecks (Validated in Phase 04)
- [ ] **Error Recovery** — Improve error handling granularity, retry policies for specific error types

### Out of Scope

- **Binary distribution** — This is a library SDK, not a CLI application
- **HTTP/REST fallback** — WebSocket-only; no REST API wrapper
- **Connection pooling** — Single connection per client instance; pooling is user responsibility
- **Built-in retry for all errors** — Only reconnect on disconnect; application-level retry is user responsibility

## Context

### Migration Background

This SDK was migrated from `openclaw-sdk-typescript` (TypeScript). The design document at `docs/specs/2026-03-18-typescript-to-go-migration-design.md` details architectural decisions. The migration preserved functional equivalence but not API identity — some APIs were adapted for Go idioms (e.g., option pattern instead of config objects, interfaces instead of classes).

### Current State (2026-03-29)

- All 127 API methods implemented across 15 namespaces
- Full manager-based architecture with 4 managers (event, request, connection, reconnect)
- CI pipeline passing: fmt, vet, lint, test with race detection
- 80%+ test coverage target
- README.md synchronized with implementation
- Client struct organized into 4 sub-struct groups (managers, api, protocol, health)
- Hot-path benchmarks established (protocol, managers) with benchstat CI integration
- Fuzz testing with round-trip assertions and 24 corpus files
- GoReleaser v2 release infrastructure configured and verified (REL-01, REL-03); git-cliff changelog automation operational

### Known Concerns

From codebase analysis (CONCERNS.md):
- README examples need real-world usage validation
- Fuzz tests exist but coverage unknown
- No integration tests against real OpenClaw gateway
- No published version history in CHANGELOG.md beyond initial structure

## Constraints

- **Go 1.21+ runtime** — No breaking changes to Go compatibility
- **No CGO** — Pure Go, no external C dependencies
- **Minimal dependencies** — Only `gorilla/websocket`; stdlib preferred
- **Library distribution** — GoReleaser configured for library mode (no binaries)
- **API compatibility** — Once v1.0.0 released, breaking changes require major version bump

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| Option pattern over config struct | Idiomatic Go, more flexible at call site | ✓ Good |
| Manager-based architecture | Aligns with TypeScript SDK, clear separation of concerns | ✓ Good |
| Channel + Context hybrid | Context for cancellation, channels for event delivery (avoids deadlock) | ✓ Good |
| No built-in retry beyond reconnect | Applications have varying retry needs; SDK provides foundation | ✓ Good |
| Re-export types from subpackages | Single-package import convenience | ✓ Good |
| JSON over binary protocol | OpenClaw gateway uses JSON; performance not critical path | ✓ Good |

---

*Last updated: 2026-03-29 after Phase 05 (release-infrastructure) completion*
