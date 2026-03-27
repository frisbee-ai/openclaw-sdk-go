# Feature Research

**Domain:** WebSocket Client SDK (Go)
**Researched:** 2026-03-28
**Confidence:** MEDIUM

**Note:** WebSearch tool was unavailable during research. Findings are based on codebase analysis, established WebSocket SDK patterns, and the CONCERNS.md audit. Some market observations are drawn from training data and should be verified against current library documentation.

## Feature Landscape

### Table Stakes (Users Expect These)

Features users assume exist in a mature WebSocket SDK. Missing these = product feels incomplete or production-unready.

| Feature | Why Expected | Complexity | Status | Notes |
|---------|--------------|------------|--------|-------|
| Connection lifecycle management | Users need to connect, disconnect, reconnect | LOW | DONE | State machine with valid transitions |
| Request/response correlation | WebSocket is async; responses must match requests | LOW | DONE | RequestID-based correlation in RequestManager |
| Event subscription/handling | Push-based communication requires pub/sub | LOW | DONE | EventManager with buffered channels |
| Authentication | Gateways require credentials | LOW | DONE | AuthHandler, CredentialsProvider interfaces |
| Typed errors | Debugging requires categorizable failures | LOW | DONE | 8 error types with code, retryable, details |
| Context propagation | Lifecycle control, cancellation | LOW | DONE | All operations accept context.Context |
| Graceful shutdown | Clean resource release | LOW | DONE | Close() with WaitGroup cleanup |
| TLS support | Production requires encrypted connections | MEDIUM | PARTIAL | Custom validator exists; CRL checking is stub |
| Reconnection with backoff | Networks are unreliable | MEDIUM | DONE | Fibonacci backoff with jitter |
| Heartbeat/ping-pong | Detect dead connections | MEDIUM | DONE | TickMonitor with stale detection |
| Gap detection | Ordered message streams can have gaps | MEDIUM | DONE | GapDetector with sequence tracking |
| Protocol version negotiation | API evolution without breaking clients | MEDIUM | DONE | ProtocolNegotiator |
| Server policy enforcement | Respect server-side limits | MEDIUM | DONE | PolicyManager with MaxPayload, rate limits |
| **Client-side rate limiting** | Prevent server rejection from overload | MEDIUM | **MISSING** | CONCERNS.md flagged; no configurable limits |
| **Retry budget for reconnection** | Prevent infinite retry loops | MEDIUM | **MISSING** | MaxAttempts=0 means unlimited; CONCERNS.md flagged |
| **Connection health metrics** | Observability beyond state changes | MEDIUM | **MISSING** | CONCERNS.md flagged; no latency/quality tracking |

### Differentiators (Competitive Advantage)

Features that set the product apart. Not required, but valuable for a Go ecosystem SDK.

| Feature | Value Proposition | Complexity | Notes |
|---------|-------------------|------------|-------|
| **Middleware/interceptors** | Allow users to log, metrics, modify requests without touching core | MEDIUM | RequestFn wrapper chain; few Go SDKs have this |
| **Request deduplication** | Idempotency keys prevent duplicate sends on retry | MEDIUM | gorilla/websocket has no built-in dedup |
| **Per-request timeout** | Different requests may need different timeouts | LOW | Currently only context-level timeout; RequestManager applies no internal deadline |
| **Streaming responses** | Server-push, chunked data delivery | HIGH | Not in OpenClaw protocol spec but future-proofing |
| **OpenTelemetry integration** | Distributed tracing for WebSocket operations | MEDIUM | Growing expectation in cloud-native Go libraries |
| **Connection resilience patterns** | Circuit breaker, bulkhead isolation | HIGH | Complex to implement correctly; valuable for robustness |
| **Structured logging adapter** | Integrate with log/slog/zap | LOW | Logger interface exists; could add adapters |
| **Metrics endpoint** | Prometheus-compatible metrics for dashboards | MEDIUM | CONCERNS.md flagged as missing observability |

### Anti-Features (Commonly Requested, Often Problematic)

Features that seem good but create problems or violate project constraints.

| Feature | Why Requested | Why Problematic | Alternative |
|---------|---------------|-----------------|-------------|
| **Binary protocol support** | Performance gains over JSON | OpenClaw gateway uses JSON; would fork the protocol | Stay with JSON; document perf expectations |
| **HTTP/REST fallback** | Fall back when WebSocket fails | Violates WebSocket-only design; adds massive complexity | Document WebSocket reliability; let users implement HTTP if needed |
| **Connection pooling** | Handle high throughput | Single connection per client is architectural constraint; pooling is user responsibility | Document multi-client pattern if high throughput needed |
| **Built-in retry for all errors** | Automatic recovery from transient failures | Varies wildly by error type and business logic; SDK cannot know retry semantics | Keep retry scoped to reconnect; let users implement application retry |
| **Connection multiplexing (channels)** | Multiple logical streams over one connection | Gorilla doesn't support it; would need different WS library | Single connection is simpler and explicit |
| **Auto-reconnection for ALL operations** | Seamlessly recover mid-flight requests | Request state is lost; responses can't be recovered | Document that in-flight requests fail on disconnect; user implements idempotent retry |

## Feature Dependencies

```
[ReconnectManager]
    └──requires──> [ConnectionManager]

[GapDetector]
    └──requires──> [EventManager]

[TickMonitor]
    └──requires──> [EventManager] ──emits──> [Tick Events]
    └──monitors──> [ConnectionManager]

[PolicyManager]
    └──enforces──> [RequestManager] (MaxPayload validation)

[ProtocolNegotiator]
    └──configures──> [ConnectionManager]

[Middleware/Interceptors] ──wraps──> [RequestManager]
```

### Dependency Notes

- **ReconnectManager requires ConnectionManager:** Reconnect calls `connection.Reconnect()` - the dependency is implicit via callback
- **GapDetector and TickMonitor both require EventManager:** They subscribe to EventTick and emit EventGap events
- **PolicyManager enforces RequestManager:** Payload size validation happens in `SendRequest()` before queuing
- **Middleware would wrap RequestManager:** A `WithRequestInterceptor` option could add middleware chain without changing existing code

## MVP Definition

### Launch With (v1.0)

Minimum viable product for production use. Based on CONCERNS.md critical gaps.

- [x] Connection lifecycle (done)
- [x] Request/response correlation (done)
- [x] Event handling (done)
- [x] Auth (done)
- [x] Reconnect with backoff (done)
- [x] TickMonitor (done)
- [x] **Rate limiting** — Add `RequestRateLimiter` interface with `WithRateLimit()` option; prevents server rejection under load (CONCERNS.md security risk)
- [x] **Retry budget** — Add `MaxRetries` to `ReconnectConfig`; replace 0 (unlimited) with sensible default like 10; document behavior change (CONCERNS.md critical)
- [x] **TLS CRL checking stub** — Implement actual `CheckCertificateRevocation` using `cert.CRLDistributionPoints` or mark as intentionally not implemented with comment (CONCERNS.md tech debt)

### Add After Validation (v1.x)

Features to add once core v1 is stable and users have provided feedback.

- [ ] **Connection health metrics** — Add `ConnectionMetrics` struct with Latency, LastTickAge, ReconnectCount; expose via `GetMetrics()` on client (CONCERNS.md missing critical feature)
- [ ] **Per-request timeout** — Allow `SendRequest(ctx, req, timeout)` or `WithRequestTimeout()` option per-call; useful for long-running operations
- [ ] **Graceful degradation** — Add priority levels to events; when EventChannel is full, drop low-priority first (CONCERNS.md scaling limit)
- [ ] **Integration tests** — Tests against real OpenClaw gateway (CONCERNS.md test gap)

### Future Consideration (v2+)

Features to defer until product-market fit is established.

- [ ] **Middleware/interceptors** — Request/response hooks for logging, metrics, auth token refresh
- [ ] **OpenTelemetry tracing** — Span propagation for WebSocket frames
- [ ] **Circuit breaker** — Prevent cascade failures when gateway is unhealthy
- [ ] **Request deduplication** — Idempotency key support for at-least-once delivery

## Feature Prioritization Matrix

| Feature | User Value | Implementation Cost | Priority | Source |
|---------|------------|---------------------|----------|--------|
| Rate limiting | HIGH | MEDIUM | P1 | CONCERNS.md (security risk) |
| Retry budget | HIGH | LOW | P1 | CONCERNS.md (critical gap) |
| TLS CRL stub fix | MEDIUM | MEDIUM | P1 | CONCERNS.md (tech debt) |
| Connection health metrics | HIGH | MEDIUM | P2 | CONCERNS.md (observability) |
| Per-request timeout | MEDIUM | LOW | P2 | DX improvement |
| Graceful degradation | MEDIUM | MEDIUM | P2 | CONCERNS.md (scaling) |
| Middleware/interceptors | MEDIUM | MEDIUM | P3 | DX improvement |
| OpenTelemetry | MEDIUM | HIGH | P3 | Future ecosystem |
| Circuit breaker | MEDIUM | HIGH | P3 | Future resilience |
| Request deduplication | LOW | MEDIUM | P3 | Edge case handling |

**Priority key:**
- P1: Must have for v1 production release
- P2: Should have, add in v1.x
- P3: Nice to have, future consideration

## Competitor Feature Analysis

| Feature | gorilla/websocket | gobwas/ws | nhooyr/websocket | Our Approach |
|---------|------------------|-----------|------------------|--------------|
| Connection management | Basic dial/close | Basic | Basic | Enhanced state machine with 7 states |
| Reconnection | None (user implements) | None | None | Built-in Fibonacci backoff with jitter |
| Heartbeat | Ping/Pong helpers only | None | None | TickMonitor with stale detection |
| Error types | Single error | Single error | Single error | 8 typed errors with codes and retryable |
| Rate limiting | None | None | None | **MISSING - add in v1** |
| Retry budget | None | None | None | **MISSING - add in v1** |
| Request correlation | None (user implements) | None | None | RequestID-based correlation done |
| Event system | None | None | None | EventManager with pub/sub done |
| Metrics/observability | None | None | None | **MISSING - add in v1.x** |
| TLS custom validation | Yes | Yes | Yes | Partial (CRL stub) |
| Protocol negotiation | None | None | None | ProtocolNegotiator done |
| Middleware | None | Pipeline support | None | **MISSING - future** |

**Analysis:** The SDK is more feature-complete than raw gorilla/websocket but lacks production hardening (rate limiting, retry budgets, observability) that enterprise users expect. The TypeScript migration brought the API surface; Go-specific hardening is still needed.

## Sources

- Codebase analysis: `/Users/linyang/workspace/my-projects/openclaw-sdk-go/pkg/` (client.go, managers/, types/, events/)
- CONCERNS.md: Security risks, performance bottlenecks, missing critical features
- PROJECT.md: Out-of-scope items, current state
- ARCHITECTURE.md: System structure, data flows
- gorilla/websocket documentation: https://pkg.go.dev/github.com/gorilla/websocket
- gobs/ws documentation: https://pkg.go.dev/github.com/gobwas/ws
- WebSocket RFC 6455: Protocol specification reference

---
*Feature research for: WebSocket Client SDK (Go)*
*Researched: 2026-03-28*
