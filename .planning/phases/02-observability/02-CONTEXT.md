# Phase 2: Observability - Context

**Gathered:** 2026-03-28
**Status:** Ready for planning

<domain>
## Phase Boundary

Phase 2 delivers **connection health visibility and graceful degradation under load** for the OpenClaw SDK. Four requirements:

- **OBS-01**: `ConnectionMetrics` struct with Latency, LastTickAge, ReconnectCount exposed via `GetMetrics()`
- **OBS-02**: Per-request timeout via `SendRequest(ctx, req, opts...)` with variadic options
- **OBS-03**: Event priority levels; when EventChannel is full, low-priority events drop first
- **OBS-04**: `EventBufferSize` configurable via client option

**OBS-04 is already implemented** — `ClientConfig.EventBufferSize` and `WithEventBufferSize()` exist in `pkg/client.go`. No work needed.

</domain>

<decisions>
## Implementation Decisions

### OBS-01: ConnectionMetrics

- **D-01 (Latency measurement):** Tick-based estimation using time between consecutive ticks. No extra probing needed — uses `TickMonitor`'s existing tick interval. `Latency` = `tickIntervalMs * staleMultiplier` as baseline estimate.
- **D-02 (Metrics API surface):** `GetMetrics() ConnectionMetrics` is added to the `OpenClawClient` interface. All users can call it. `ConnectionMetrics` is a public struct.
- **D-03 (Metrics struct fields):**
  ```go
  type ConnectionMetrics struct {
      Latency        time.Duration // tick-based estimate: tickInterval * staleMultiplier
      LastTickAge    time.Duration // actual time since last tick received
      ReconnectCount int           // total reconnection attempts made
      IsStale        bool         // whether connection is currently stale
  }
  ```
- **D-04 (ReconnectCount source):** `ReconnectManager` needs an exported attempt counter. Add `AttemptCount()` method returning the local `attempt` counter from `run()`. Thread-safe via existing mutex. Counter is read-only snapshot (not reset on success).

### OBS-02: Per-Request Timeout

- **D-05 (SendRequest signature):** Extend to variadic options for backward compatibility.
  ```go
  SendRequest(ctx context.Context, req *protocol.RequestFrame, opts ...RequestOption) (*protocol.ResponseFrame, error)
  ```
  Existing callers with `SendRequest(ctx, req)` continue to work unchanged. No new method name needed.
- **D-06 (RequestOption design):** `RequestOption` is a functional option type.
  ```go
  type RequestOption func(*requestConfig)
  type requestConfig struct {
      timeout time.Duration
      // reserved for future options (progress callback, idempotency key, etc.)
  }
  func WithRequestTimeout(d time.Duration) RequestOption
  ```
- **D-07 (Timeout precedence):** `WithRequestTimeout` wraps the incoming `ctx` with `context.WithTimeout`. If the caller already set a deadline on `ctx`, it is overwritten by `WithRequestTimeout`. This is explicit — caller chooses per-request timeout, not additive.
- **D-08 (Placement):** Options check happens in `client.SendRequest` (not `RequestManager.SendRequest`). The `client` layer is where the options are parsed; `RequestManager.SendRequest` receives the already-wrapped context.

### OBS-03: Event Priority Levels

- **D-09 (Priority levels):** 3 levels — `EventPriorityHigh`, `EventPriorityMedium`, `EventPriorityLow`. Type: `EventPriority` (string or int type).
  ```go
  type EventPriority int
  const (
      EventPriorityLow    EventPriority = 0
      EventPriorityMedium EventPriority = 1
      EventPriorityHigh   EventPriority = 2
  )
  ```
- **D-10 (Event struct change):** Add `Priority EventPriority` field to `types.Event`. Default priority is `EventPriorityMedium` for backward compatibility (existing events without explicit priority get MEDIUM).
  ```go
  type Event struct {
      Priority  EventPriority
      Type      EventType
      Payload   any
      Err       error
      Timestamp time.Time
  }
  ```
- **D-11 (Priority assignment — HIGH):** `EventError`, `EventDisconnect`, `EventStateChange`, `EventGap`
- **D-12 (Priority assignment — MEDIUM):** `EventTick`, `EventResponse`, `EventConnect`
- **D-13 (Priority assignment — LOW):** `EventMessage`, `EventRequest`
- **D-14 (Emit drop behavior):** When `events` channel is full, `EventManager` prioritizes draining HIGH events. Drop order: LOW first, then MEDIUM, then HIGH. HIGH events should never drop unless channel is completely full and all drain attempts failed (panic-level safety net).
- **D-15 (Emit implementation approach):** `Emit()` checks priority before selecting the send path. Simple approach: iterate from LOW to HIGH, trying each priority level's buffered channel. If all channels full, drop based on priority. See `pkg/managers/event.go` for existing `emitTimer` + select pattern — priority can layer on top.

### Claude's Discretion

The following are delegated to researcher/planner judgment:
- Exact internal channel structure for priority-based event dispatch (separate per-priority channels vs. single channel with priority tagging)
- How `ReconnectManager.AttemptCount()` is implemented (expose attempt as atomic int or mutex-protected int)
- Whether `GetMetrics()` should return a copy (defensive) or direct reference (for zero-copy observability)
- Whether `TickMonitor` needs a `GetTickInterval()` method to support tick-based latency estimation

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Phase 1 Decisions (must be consistent)
- `.planning/01-foundation/CONTEXT.md` — Rate limiter placement, channel ownership fix, mutex scope patterns
- `.planning/01-foundation/01-PLAN.md` — Foundational patterns that apply to all phases
- `.planning/01-foundation/02-PLAN.md` — Rate limiting and pending limit implementation
- `.planning/01-foundation/03-PLAN.md` — TLS and reconnect implementation

### Project Requirements
- `.planning/ROADMAP.md` §Phase 2 — Phase 2 goals, success criteria, OBS-01 through OBS-04
- `.planning/REQUIREMENTS.md` §OBS-01 through OBS-04 — Acceptance criteria for each observability requirement
- `.planning/STATE.md` — Prior decisions and accumulated context

### Codebase Architecture
- `pkg/client.go` — Main client, `OpenClawClient` interface, `ClientConfig`, option functions
- `pkg/managers/event.go` — `EventManager` with existing `Emit()` + `emitTimer` + select drop pattern
- `pkg/managers/request.go` — `RequestManager` with existing `RequestOptions` struct (has `Timeout any`)
- `pkg/managers/reconnect.go` — `ReconnectManager` (needs `AttemptCount()` method)
- `pkg/events/tick.go` — `TickMonitor` with `GetTimeSinceLastTick()`, `IsStale()`, `GetStaleDuration()`
- `pkg/types/types.go` — `Event` struct (needs `Priority` field), `EventType` constants

### No external specs
No external specs — requirements fully captured in decisions above.

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- **`TickMonitor.GetTimeSinceLastTick()`** (`pkg/events/tick.go:164`): Already returns `int64` milliseconds since last tick — wrap in `ConnectionMetrics.LastTickAge`
- **`TickMonitor.IsStale()`** (`pkg/events/tick.go:139`): Already exists — expose in `ConnectionMetrics.IsStale`
- **`EventManager.emitTimer`** (`pkg/managers/event.go:33`): Reusable timer pattern for bounded-wait sends
- **`RequestOptions` struct** (`pkg/managers/request.go:19`): Already exists with `Timeout any` field — extend with `RequestOption` functional option pattern
- **`types.Logger`** (`pkg/types/logger.go`): `EventManager` already uses it — use same logger for metrics-related logging

### Established Patterns
- **Functional options**: `WithURL()`, `WithRateLimit()`, `WithMaxPending()` all use `ClientOption func(*ClientConfig) error`. `RequestOption` follows the same pattern.
- **Mutex scope discipline**: Lock held for minimum time; release before expensive operations. Applies to `GetMetrics()` — snapshot under lock, return copy.
- **No new goroutines in getters**: `TickMonitor` uses `sync.RWMutex` with `RLock` for read-only state access.
- **Error types**: `NewRequestError()` already exists — per-request timeout errors should use this.

### Integration Points
- **`OpenClawClient` interface** (`pkg/client.go:327`): `GetMetrics() ConnectionMetrics` added here. All 15 API accessor methods listed after it — new method goes before them.
- **`client` struct** (`pkg/client.go:369`): Already has `managers`, `tickMonitor`, `gapDetector`. Needs `GetMetrics()` method aggregating data from `managers.reconnect` + `tickMonitor`.
- **`EventManager`**: `events` channel (`make(chan types.Event, bufferSize)`) needs priority logic. `EventManager.NewEventManager()` constructor already takes `bufferSize` — no config change needed.
- **`Event.Type` field** (`pkg/types/types.go:50`): Add `Priority EventPriority` field to existing struct.

</code_context>

<specifics>
## Specific Ideas

- **Latency = tick-based estimate**: `Latency = time.Duration(tm.tickIntervalMs) * time.Duration(tm.staleMultiplier) * time.Millisecond`. This is the expected tick interval as a proxy for RTT — not actual RTT measurement.
- **Drop behavior must be deterministic**: When buffer is full, LOW events must drop before MEDIUM, MEDIUM before HIGH. The implementation should iterate from lowest to highest priority.
- **ReconnectCount is a snapshot**: `ReconnectManager.AttemptCount()` returns current count at call time. It does not reset when reconnect succeeds (that's the nature of a counter).
- **Variadic options for SendRequest**: `client.SendRequest(ctx, req)` with 0 args works exactly as today. Only new callers use `WithRequestTimeout`.

</specifics>

<deferred>
## Deferred Ideas

### Reviewed Todos (not folded)
None — no pending todos matched Phase 2 scope.

### Out-of-Scope (discussed but not this phase)
- **OBS-01 Latency — actual RTT via ping/pong**: Deferred to Phase 4 (Benchmarking). Tick-based estimate is good enough for v1.0 observability.
- **OBS-03 — more than 3 priority levels**: Phase 2 uses 3 levels. 5-level scheme deferred to future if real use cases emerge.
- **Metrics export via Prometheus/OpenTelemetry**: Phase 5 REL mentions this. Internal `GetMetrics()` struct is the foundation.

</deferred>

---

*Phase: 02-observability*
*Context gathered: 2026-03-28*
