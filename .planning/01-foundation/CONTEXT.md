# Phase 1 Context: Foundation Hardening

**Phase:** 1 of 5
**Requirements:** FOUND-01, FOUND-02, FOUND-03, FOUND-04, FOUND-05
**Status:** Ready to plan

## Decisions from Research

### Already Decided

- **Rate limiting location**: `RequestManager.SendRequest()` — pre-send check before serialize
- **Rate limiting interface pattern**: `RequestRateLimiter` interface with `Allow()` method + `WithRateLimit()` option
- **Retry budget field**: `MaxRetries` added to `ReconnectConfig` (distinct from `MaxAttempts` for clarity)
- **Pending request limit location**: `RequestManager.SendRequest()` — check map size before adding
- **Pending limit error**: `ErrTooManyPendingRequests` returned immediately when limit reached
- **InsecureSkipVerify warning location**: `TlsValidator.GetTLSConfig()` — log warning after config built
- **TLS CRL approach**: Stub with explicit comment (not implementing actual CRL fetching for v1)

## Gray Areas (Pending Discussion)

### GA-1: Rate Limiter Interface Design

**Question**: What should the `RequestRateLimiter` interface look like?

Two options:

**Option A — Simple token bucket**:
```go
type RequestRateLimiter interface {
    Allow() bool  // returns false when rate limited
}
```
- Pro: Simple, matches standard Go idiom
- Con: No feedback on when to retry

**Option B — Check with retry info**:
```go
type RequestRateLimiter interface {
    Allow() (allowed bool, retryAfter time.Duration)
}
```
- Pro: Provides retry-after feedback
- Con: More complex, not all limiters support retry-after

**Recommendation**: Option A — simpler, `RequestManager.SendRequest` can use `context.WithTimeout` for retry-after

### GA-2: Rate Limiter Placement in Client Chain

**Question**: Where in the call chain should rate limiting be checked?

**Option A — RequestManager (current plan)**:
- Rate limiting happens at lowest level
- All API calls subject to same limiter
- Con: Can't have different limits per namespace

**Option B — Client.SendRequest (higher level)**:
- Rate limiting at client level
- Easier to add per-client limit
- Pro: Aligns with `WithRateLimit()` option being client-level

**Recommendation**: Option B — `client.SendRequest()` checks limiter before calling `RequestManager.SendRequest`

### GA-3: Default Rate Limit Value

**Question**: What should the default max pending requests limit be?

| Value | Rationale |
|-------|-----------|
| 100 | Conservative — matches EventBufferSize |
| 256 | Moderate — common WebSocket pipeline size |
| 1024 | Aggressive — high-throughput scenarios |

**Recommendation**: 256 — balances memory usage against throughput needs

### GA-4: Default MaxRetries Value

**Question**: `MaxRetries=10` as required, but what should the SDK default be?

Current: `MaxAttempts: 0` (infinite)
Required: `MaxRetries: 10` (FOUND-02)

**Option A — Change DefaultReconnectConfig default to 10**:
- Breaking change if anyone relies on infinite retries
- But: research shows this is a P1 gap, so production users would want this

**Option B — Keep MaxAttempts=0, add MaxRetries=10**:
- Maintain backward compatibility
- `MaxRetries` field in `ReconnectConfig` with 10 default
- Existing `MaxAttempts` field deprecated but functional

**Recommendation**: Option B — backward compatible, explicit MaxRetries field

### GA-5: InsecureSkipVerify Warning Log Level

**Question**: What log level should the InsecureSkipVerify warning use?

| Level | Rationale |
|-------|-----------|
| WARN | Standard warning — user should be concerned |
| ERROR | Stronger signal — this is a security issue |
| INFO | Less alarming — some users legitimately use it for testing |

**Recommendation**: WARN — appropriate severity for a security-sensitive configuration

## Implementation Plan Summary

```
pkg/
├── types/
│   └── errors.go              # Add ErrTooManyPendingRequests
├── managers/
│   ├── request.go            # Add rate limiter, pending limit
│   └── reconnect.go          # Add MaxRetries support
├── connection/
│   └── tls.go               # Add InsecureSkipVerify warning
└── client.go                 # Add rate limiter to SendRequest path
```

## Files to Modify

| File | Changes |
|------|---------|
| `pkg/types/errors.go` | Add `ErrTooManyPendingRequests` |
| `pkg/types/types.go` | Add `RequestRateLimiter` interface, `WithRateLimit()` option |
| `pkg/managers/request.go` | Add rate limiter field, check in SendRequest, check pending map size |
| `pkg/managers/reconnect.go` | Support MaxRetries, return `ErrMaxRetriesExceeded` |
| `pkg/connection/tls.go` | Add warning log in GetTLSConfig when InsecureSkipVerify=true |
| `pkg/client.go` | Add rate limiter to ClientConfig and SendRequest chain |

## Success Criteria

1. Rate limiting: SendRequest returns error immediately when limiter denies
2. Retry budget: After MaxRetries=10, reconnect stops and returns `ErrMaxRetriesExceeded`
3. TLS CRL: Function exists with stub implementation and clear documentation
4. Pending limit: SendRequest returns `ErrTooManyPendingRequests` when limit reached
5. InsecureSkipVerify warning: Warning logged at connection time when enabled

---

*Context created: 2026-03-28*
*Next: Discuss gray areas above, then /gsd:plan-phase 1*
