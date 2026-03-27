# Domain Pitfalls: Go WebSocket SDKs

**Project:** OpenClaw SDK Go
**Domain:** WebSocket client library (Go)
**Researched:** 2026-03-28
**Confidence:** HIGH (gorilla/websocket official docs) / MEDIUM (community patterns from WebFetch)

---

## Critical Pitfalls

Mistakes that cause rewrites, security vulnerabilities, or production outages.

### Pitfall 1: Concurrent Write Corruption

**What goes wrong:** The WebSocket connection enters a corrupt state when multiple goroutines call write methods (`WriteMessage`, `WriteJSON`, `NextWriter`) concurrently.

**Root cause:** gorilla/websocket requires exactly one concurrent writer. The underlying protocol frames can interleave and corrupt data.

**Consequences:** Silent data corruption, interleaved messages, connection hangs.

**Prevention:**
- All write methods must be serialized behind a single mutex or channel.
- The SDK's `RequestManager` and `ConnectionManager` must never call write methods concurrently.
- Document that users must not call `SendRequest` concurrently from multiple goroutines without their own synchronization.

**Detection:** Race detector (`go test -race`) catches this. Run with race detection always.

---

### Pitfall 2: Failing to Read the Connection (Deadlock)

**What goes wrong:** If the application stops reading messages, the WebSocket connection eventually blocks because the peer cannot send control messages (pings, close frames).

**Root cause:** The WebSocket protocol requires reading to process `pong`, `close`, and `ping` messages. Without a read loop, the write side fills the TCP buffers and deadlocks.

**Consequences:** Connection hangs, goroutine leaks, no response to server pings.

**Prevention:**
- Always run a read loop in a dedicated goroutine.
- Ensure the read loop handles `CloseError` gracefully and returns.
- The SDK's `transport.readLoop()` must be started before any write and must survive until `Close()`.

**Detection:** Goroutine dumps show blocked send on connection. Test with `timeouts` on read operations.

---

### Pitfall 3: Missing Origin Checking (CSRF Vulnerability)

**What goes wrong:** A malicious website can open a WebSocket connection to your server on behalf of a logged-in user.

**Root cause:** The `Upgrader` defaults to allowing all origins if `CheckOrigin` is nil. Browsers send the `Origin` header automatically.

**Consequences:** Cross-site WebSocket hijacking, session token theft.

**Prevention:**
- Set explicit `CheckOrigin` policy in the `Upgrader`.
- For browser clients, validate `Origin` header matches expected domain.
- The SDK uses `Dial` (client-side), so this is less critical than server-side, but still document that the server must validate origins.

**Detection:** Security audit. Missing `CheckOrigin` in server implementations.

---

### Pitfall 4: Write Deadline Corruption

**What goes wrong:** After a write times out, the WebSocket connection is in a corrupt state. All subsequent writes fail.

**Root cause:** gorilla/websocket documentation explicitly states: "After a write has timed out, the websocket state is corrupt and all future writes will return an error."

**Consequences:** SDK enters unrecoverable error state after any transient network slowdown.

**Prevention:**
- Set write deadlines only when absolutely necessary.
- After a write error, the connection must be closed and reconnected, not reused.
- The SDK's `ReconnectManager` should handle this case by reconnecting.

**Detection:** Logs show `websocket: write timeout` followed by cascade of errors.

---

### Pitfall 5: Ignoring CloseError

**What goes wrong:** Close errors contain important metadata (close code, reason) that is silently discarded.

**Root cause:** `websocket.IsUnexpectedCloseError` is not used to distinguish expected vs unexpected closes.

**Consequences:** Masked disconnects (e.g., server restart vs client ban), lost debugging information.

**Prevention:**
- Always check `websocket.IsUnexpectedCloseError(err, expectedCodes...)`.
- Log close codes for debugging.
- Handle `CloseError` separately from network errors.

**Detection:** `go test -v` shows passing tests but real disconnects are silent in production logs.

---

## Moderate Pitfalls

### Pitfall 6: Channel + Lock Ordering Violation (Deadlock)

**What goes wrong:** Sending to a buffered channel while holding a mutex causes a deadlock if the receiver holds the same mutex.

**Root cause:** Classic deadlock pattern: `lock() -> send to channel -> receiver tries to lock() -> deadlock`.

**Consequences:** All goroutines freeze. Process becomes unresponsive.

**Prevention:**
- **Rule:** Never send to a channel while holding a lock. Release lock BEFORE sending.
- The SDK's `ConnectionStateMachine` follows this: it releases `RLock` before sending to `eventChan`.
- All callbacks from managers to client must not hold locks.

**Detection:** Go race detector does not catch this. Use deadlock detection in tests or careful code review.

---

### Pitfall 7: Unbounded Pending Request Map (Memory Leak)

**What goes wrong:** RequestManager's pending request map grows without bound on long-running connections.

**Root cause:** Each `SendRequest` adds to the map. If requests timeout but the map entry is not removed, or if there are many requests, memory grows indefinitely.

**Consequences:** Memory leak, eventual OOM on long-running clients.

**Prevention:**
- Set a maximum pending request limit. Reject new requests with `ErrTooManyPendingRequests` when limit is reached.
- Always remove entries on timeout (already done in SDK, but must be verified).
- Document the limit for users.

**Detection:** Memory profiler shows growing `map[RequestID]*requestInfo`.

---

### Pitfall 8: Timer Leak in Reconnect Loop

**What goes wrong:** `time.NewTimer` or `time.NewTicker` is not stopped on context cancellation, leaking timers.

**Root cause:** Timers are not cleaned up when the reconnection context is cancelled (e.g., user calls `Close()`).

**Consequences:** Goroutine leak, timer goroutines accumulate over reconnect cycles.

**Prevention:**
- Always `defer timer.Stop()` or use `sync.Once` to ensure cleanup.
- Use `runtime.SetFinalizer` sparingly as a backstop.
- The SDK's `ReconnectManager` uses `time.NewTimer` per iteration with `defer stop()` pattern.

**Detection:** `go test -race` with `-timer-check` or looking for timer goroutines.

---

### Pitfall 9: Deprecated API Usage

**What goes wrong:** Using old package-level functions (`websocket.Upgrade()`, `websocket.NewClient()`) instead of `Upgrader` and `Dialer`.

**Root cause:** The package-level functions are deprecated and miss security features like origin checking.

**Consequences:** Security vulnerabilities, behavior differences in edge cases.

**Prevention:**
- Use `websocket.Upgrader` for servers, `websocket.Dialer` for clients.
- The SDK uses `Dialer.Dial()` for outgoing connections, which is correct.
- Lint with `golangci-lint` to catch deprecated API usage.

**Detection:** `go vet` catches some; add explicit nolint comments if deprecated is intentional.

---

### Pitfall 10: Compression Misunderstanding

**What goes wrong:** Enabling compression assumes "context takeover" (sliding window reuse), which gorilla/websocket does not support.

**Root cause:** gorilla/websocket compression is experimental and does not maintain compression state across messages.

**Consequences:** Unexpected memory usage, worse compression ratios than expected.

**Prevention:**
- If compression is enabled, benchmark performance.
- Do not assume compression improves with message correlation.
- The SDK does not enable compression by default.

**Detection:** Memory profiler shows compression buffers not being reused.

---

## Minor Pitfalls

### Pitfall 11: Buffer Size vs Message Size Confusion

**What goes wrong:** Setting `ReadBufferSize` / `WriteBufferSize` thinking it limits max message size.

**Root cause:** Buffer sizes are for I/O efficiency, not message limits. `ReadLimit` is the actual message size limit.

**Consequences:** Large messages are accepted despite small buffers; small buffers cause performance issues.

**Prevention:**
- Use `ReadLimit` to set maximum message size.
- Use buffers sized to 99th percentile message size for efficiency.
- Document this distinction.

**Detection:** Benchmarking shows unexpected memory allocations.

---

### Pitfall 12: InsecureSkipVerify Without Warning

**What goes wrong:** `InsecureSkipVerify: true` allows man-in-the-middle attacks.

**Root cause:** Convenience option used in development but shipped to production.

**Consequences:** Full TLS MITM, credential theft, data tampering.

**Prevention:**
- Default to secure TLS. Never default to skip verify.
- Log a warning when `InsecureSkipVerify` is true.
- CONCERNS.md (line 51) flags this: add warning log.
- Consider removing entirely; force users to provide custom `TlsValidator` for self-signed certs.

**Detection:** Security audit, code review.

---

### Pitfall 13: Context Inheritance from HTTP Request

**What goes wrong:** Using `r.Context()` from an HTTP request for the WebSocket connection lifecycle.

**Root cause:** `http.Hijacker` interface changes the context behavior unexpectedly.

**Consequences:** Connection closes when HTTP request context cancels, not when intended.

**Prevention:**
- Create a fresh context for the WebSocket connection: `context.WithTimeout(context.Background(), ...)` or `context.WithCancel(context.Background())`.
- The SDK correctly creates fresh contexts for its operations.

**Detection:** Connection drops immediately after handshake in some deployments.

---

### Pitfall 14: Silent Event Drops Under Load

**What goes wrong:** When the event channel is full, events are dropped after a 200ms timeout.

**Root cause:** Non-blocking send with fixed 200ms `EventEmitTimeout`. High-throughput scenarios lose events silently.

**Consequences:** Missed events, state inconsistency, hard-to-reproduce bugs in production.

**Prevention:**
- Make `EventBufferSize` configurable.
- Consider overflow handling: slow consumer, backpressure signal, or metrics.
- Document the behavior for users.

**Detection:** Missing events in high-throughput scenarios; monitoring shows event channel full.

---

## Phase-Specific Warnings

| Phase Topic | Likely Pitfall | Mitigation |
|-------------|---------------|------------|
| Connection lifecycle | Write deadline corruption (Pitfall 4) | Reconnect on any write error, never reuse connection after timeout |
| Event dispatch | Channel + lock deadlock (Pitfall 6) | Enforce: lock release before channel send |
| Reconnection | Timer leak (Pitfall 8) | Always defer timer.Stop() in reconnect loop |
| Request/Response | Unbounded map (Pitfall 7) | Add max pending requests limit |
| TLS/Certificates | InsecureSkipVerify without warning (Pitfall 12) | Log warning, prefer custom TlsValidator |
| Read loop | Failing to read (Pitfall 2) | readLoop must run until Close(), never exit early on error |
| API stability testing | Concurrent write corruption (Pitfall 1) | Race detector (`-race`) must pass |

---

## Sources

- [gorilla/websocket official docs](https://pkg.go.dev/github.com/gorilla/websocket) — HIGH confidence
- [gorilla/websocket issues (GitHub)](https://github.com/gorilla/websocket/issues) — HIGH confidence
- [CONCERNS.md (.planning/codebase)](file://.planning/codebase/CONCERNS.md) — HIGH confidence (internal)
- [coder/websocket comparison](https://github.com/nhooyr/websocket) — MEDIUM confidence
- [nhooyr/websocket docs](https://pkg.go.dev/nhooyr.io/websocket) — MEDIUM confidence
