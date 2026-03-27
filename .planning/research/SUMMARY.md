# Project Research Summary

**Project:** OpenClaw SDK Go
**Domain:** WebSocket Client SDK Library (Go)
**Researched:** 2026-03-28
**Confidence:** MEDIUM

## Executive Summary

OpenClaw SDK Go is a production-grade WebSocket client library providing feature-complete connection management, request/response correlation, event dispatch, and automatic reconnection for the OpenClaw gateway. The existing implementation has a solid foundation (Go 1.21+, gorilla/websocket v1.5.3, proper layered architecture) but requires targeted hardening before v1 release: client-side rate limiting, retry budgets, and TLS CRL validation are missing. The architecture follows established patterns (channel-based event dispatch, pending-request map correlation, state machine for connection lifecycle) that are well-suited to the domain. The primary risks are concurrency-related pitfalls around gorilla/websocket's single-writer requirement and channel+lock ordering violations.

## Key Findings

### Recommended Stack

The core stack is already established and working. The research identified gaps in tooling depth rather than core technology changes.

**Core technologies:**
- **Go 1.21+** (runtime) / **Go 1.24** (CI) -- Language runtime; already established with strong backward-compatibility guarantees for v1.x
- **gorilla/websocket v1.5.3** -- WebSocket client; de facto standard for Go WebSocket, already in use
- **Go modules** -- Dependency management; already in use with go.mod/go.sum committed
- **GoReleaser v2** -- Release automation; partially configured, needs completion for library distribution

**Tooling gaps to address:**
- Fuzz tests exist but are shallow (panic checks only, no correctness assertions) -- need round-trip validation and corpus files
- No benchmark files exist -- need hot-path benchmarks for transport write/read, event dispatch, request correlation
- No version tags exist -- need `git tag v0.1.0` progression to `v1.0.0`
- GoReleaser config incomplete -- needs `blobs: true` for GitHub Releases hosting and `gomod.proxy: true` for verified builds

### Expected Features

**Must have (table stakes) -- mostly done, 3 gaps remain:**
- Connection lifecycle, request/response correlation, event handling, auth, typed errors, context propagation, graceful shutdown -- all DONE
- Reconnection with backoff, heartbeat/ping-pong, gap detection, protocol negotiation, server policy enforcement -- all DONE
- **Client-side rate limiting** -- MISSING (CONCERNS.md security risk: no configurable limits to prevent server rejection under load)
- **Retry budget** -- MISSING (CONCERNS.md critical: MaxRetries=0 means unlimited retries; needs sensible default like 10)
- **TLS CRL checking** -- STUB (CONCERNS.md tech debt: `CheckCertificateRevocation` not implemented)

**Should have (competitive, add in v1.x):**
- Connection health metrics -- latency, last tick age, reconnect count for observability
- Per-request timeout -- different requests may need different timeouts
- Graceful degradation -- priority levels for events when channel is full
- Integration tests against real OpenClaw gateway

**Defer (v2+):**
- Middleware/interceptors -- request/response hooks for logging, metrics, auth refresh
- OpenTelemetry tracing -- span propagation for WebSocket frames
- Circuit breaker -- prevent cascade failures when gateway is unhealthy
- Request deduplication -- idempotency key support

### Architecture Approach

The SDK uses a layered architecture with clear component boundaries: SDK entry (client.go), manager coordination layer (EventManager, RequestManager, ConnectionManager, ReconnectManager), protocol/connection layer (state machine, policies, protocol negotiation), and transport layer (gorilla/websocket wrapper). Data flows through connection lifecycle (Dial -> handshake -> authenticated), request/response (SendRequest -> correlation -> response), and event dispatch (server frame -> readLoop -> EventManager.Emit -> fan-out to handlers).

**Major components:**
1. **Transport** -- Low-level WebSocket I/O via gorilla/websocket; `readLoop`/`writeLoop` with channel-based communication
2. **Managers** -- Four managers with single responsibilities; EventManager (pub/sub with buffered channels), RequestManager (RequestID correlation), ConnectionManager (lifecycle + state transitions), ReconnectManager (Fibonacci backoff)
3. **Connection** -- State machine (7 states with enforced transitions), PolicyManager (server capability caching), ProtocolNegotiator (version negotiation), TlsValidator

**Anti-patterns to fix:**
- `client` struct is too large (direct fields for 15 API namespaces, protocolNegotiator, policyManager, tickMonitor, gapDetector) -- group into logical sub-structs
- Close/Disconnect confusion in ConnectionManager -- both methods exist with unclear semantics; needs single method with clear intent

### Critical Pitfalls

1. **Concurrent Write Corruption** -- gorilla/websocket requires exactly one concurrent writer; all writes must be serialized behind a mutex. Prevention: `go test -race` must pass; document that users must not call `SendRequest` concurrently without their own synchronization.
2. **Failing to Read the Connection (Deadlock)** -- WebSocket protocol requires reading for pongs/closes; without read loop the write side fills TCP buffers and deadlocks. Prevention: `readLoop` must run until `Close()`, never exit early on error.
3. **Write Deadline Corruption** -- After a write times out, the WebSocket connection is corrupt and all future writes fail. Prevention: set write deadlines only when necessary; after write error, close and reconnect, never reuse.
4. **Channel + Lock Ordering Violation (Deadlock)** -- Sending to a buffered channel while holding a mutex causes deadlock if the receiver tries to acquire the same lock. Prevention: **Rule: never send to a channel while holding a lock. Release lock BEFORE sending.**
5. **Timer Leak in Reconnect Loop** -- `time.NewTimer` not stopped on context cancellation leaks goroutines. Prevention: always `defer timer.Stop()` in reconnect loop.

## Implications for Roadmap

Based on research, the recommended phase structure follows the build order from ARCHITECTURE.md (dependency order) combined with the P1 feature gaps from FEATURES.md.

### Phase 1: Foundation Hardening
**Rationale:** Must fix concurrency safety and memory management before adding features. These are the root causes of production outages.
**Delivers:** Production-safe core SDK with rate limiting, retry budgets, and TLS CRL validation.
**Addresses:**
- Rate limiting: Add `RequestRateLimiter` interface with `WithRateLimit()` option (FEATURES.md P1)
- Retry budget: Add `MaxRetries` to `ReconnectConfig`; replace 0 (unlimited) with default 10 (FEATURES.md P1)
- TLS CRL: Implement actual `CheckCertificateRevocation` or mark stub with comment (FEATURES.md P1)
**Avoids:**
- Pitfall 1 (concurrent write corruption) -- ensure all write paths are serialized
- Pitfall 7 (unbounded pending request map) -- add max pending requests limit with `ErrTooManyPendingRequests`
- Pitfall 12 (InsecureSkipVerify without warning) -- add warning log when TLS skip is used
**Research flag:** None needed -- patterns are well-documented in gorilla/websocket docs and CONCERNS.md.

### Phase 2: Observability
**Rationale:** Production deployments require metrics and health signals. Adds value without changing core behavior.
**Delivers:** Connection health metrics, structured logging adapters, Prometheus-compatible endpoint.
**Addresses:**
- Connection health metrics: `ConnectionMetrics` struct with Latency, LastTickAge, ReconnectCount; expose via `GetMetrics()` (FEATURES.md P2)
- Graceful degradation: priority levels for events; when EventChannel is full, drop low-priority first (FEATURES.md P2)
**Avoids:**
- Pitfall 14 (silent event drops under load) -- configurable `EventBufferSize` and overflow handling
**Research flag:** Phase 2 -- confirm metrics format with OpenClaw gateway team; integration tests against real gateway needed.

### Phase 3: Client Struct Refactor
**Rationale:** The `client` struct has grown too large with direct fields for all sub-components. Grouping into logical sub-structs improves testability and maintainability without changing behavior.
**Delivers:** Refactored client struct with `core`, `protocol`, `health`, and `api` sub-structs.
**Addresses:**
- Anti-pattern: God Client Struct (ARCHITECTURE.md)
- Anti-pattern: Close/Disconnect confusion (ARCHITECTURE.md) -- single method with clear semantics
**Avoids:**
- Pitfall 6 (channel+lock deadlock) -- refactor ensures lock-release-before-channel-send pattern is consistent
**Research flag:** None needed -- refactoring of existing code, no new research needed.

### Phase 4: Benchmarking and Fuzz Testing
**Rationale:** Library maturity requires performance validation and regression detection. This is tooling, not feature work.
**Delivers:** Benchmark files for hot paths, CI-integrated fuzzing with correctness assertions, `benchstat` for regression detection.
**Addresses:**
- STACK.md: shallow fuzz tests need corpus files and round-trip validation
- STACK.md: no benchmark files exist
**Avoids:**
- Pitfall 11 (buffer size confusion) -- benchmarks will validate sizing decisions
**Research flag:** None needed -- Go stdlib `testing.F` and `testing.B` are well-documented.

### Phase 5: Tooling and Release Infrastructure
**Rationale:** Library distribution requires proper release automation and semantic versioning.
**Delivers:** GoReleaser library mode completion, git tags v0.1.0 progression, CHANGELOG automation with git-cliff.
**Addresses:**
- STACK.md: GoReleaser needs `blobs: true` and `gomod.proxy: true`
- STACK.md: no version tags exist
**Research flag:** None needed -- GoReleaser documentation is comprehensive.

### Phase Ordering Rationale

1. **Foundation Hardening first** -- Concurrency safety and resource limits are prerequisites for production use. Adding features (Phase 2) on top of unsafe code is wasteful.
2. **Observability second** -- Metrics and health signals are expected in production libraries but are additive changes.
3. **Client Struct Refactor third** -- Refactoring the client struct should happen before or alongside observability (both are additive improvements). This refactor is low-risk if done after Phase 1 hardening.
4. **Benchmarking/Fuzz fourth** -- These are quality-assurance tools, not features. They should be in place before heavy feature work begins.
5. **Release Infrastructure last** -- Versioning and distribution tooling is the final step before v1.0 release.

### Research Flags

Phases likely needing deeper research during planning:
- **Phase 2 (Observability):** Integration tests against real OpenClaw gateway needed -- `gsd:research-phase` recommended for API compatibility verification
- **Phase 2 (Observability):** Prometheus metrics format should be confirmed with gateway team before committing to labels

Phases with standard patterns (skip research-phase):
- **Phase 1:** gorilla/websocket well-documented; CONCERNS.md provides clear requirements
- **Phase 3:** Refactoring existing code, no new external dependencies
- **Phase 4:** Go stdlib `testing.F` and `testing.B` are fully documented
- **Phase 5:** GoReleaser official docs are comprehensive

## Confidence Assessment

| Area | Confidence | Notes |
|------|------------|-------|
| Stack | MEDIUM | Core stack established and working; gaps in tooling are well-understood but some sources (GoReleaser web docs) were incomplete |
| Features | MEDIUM | Based on codebase analysis and CONCERNS.md; WebSearch unavailable during research, some competitive analysis from training data |
| Architecture | MEDIUM-HIGH | Based on deep codebase analysis; patterns are standard and well-documented |
| Pitfalls | HIGH | gorilla/websocket official docs and GitHub issues provide high-confidence prevention strategies |

**Overall confidence:** MEDIUM

The SDK is functionally complete for core use cases. The research clearly identifies what remains (3 P1 gaps, several P2 gaps). Confidence is MEDIUM rather than HIGH because WebSearch was unavailable for competitive landscape validation and some tooling recommendations are based on incomplete extraction.

### Gaps to Address

- **Competitive landscape:** WebSearch was unavailable; some competitor feature comparisons (nhooyr/websocket, gobwas/ws) drawn from training data and should be verified against current library documentation
- **Integration testing:** No tests against real OpenClaw gateway; Phase 2 needs this before metrics format is finalized
- **TLS CRL:** The stub implementation needs either actual `cert.CRLDistributionPoints` checking or explicit decision that this is out-of-scope for v1

## Sources

### Primary (HIGH confidence)
- gorilla/websocket official docs (pkg.go.dev) -- Pitfalls, concurrency requirements, write deadline behavior
- gorilla/websocket GitHub issues -- Concurrent write corruption, timer leaks
- Current codebase analysis (2026-03-28) -- Architecture, feature status, build order
- CONCERNS.md -- Critical security risks, performance bottlenecks, missing features

### Secondary (MEDIUM confidence)
- Go stdlib `testing.F` / `testing.B` documentation -- Fuzz and benchmark tooling
- Go module versioning and semantic import versioning (go.dev/blog) -- Semantic versioning rules
- GoReleaser official documentation -- Library distribution patterns (web extraction incomplete)

### Tertiary (LOW confidence)
- coder/websocket and nhooyr/websocket competitive comparison -- Training data; verify against current docs

---

*Research completed: 2026-03-28*
*Ready for roadmap: yes*
