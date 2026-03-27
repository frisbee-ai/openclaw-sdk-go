# Codebase Concerns

**Analysis Date:** 2026-03-28

## Tech Debt

**TLS Certificate Revocation Not Implemented:**
- Issue: `CheckCertificateRevocation` function in `pkg/connection/tls.go` (line 254) is a stub that returns nil without performing actual CRL/OCSP checking
- Files: `pkg/connection/tls.go:258-280`
- Impact: Revoked certificates may not be detected during TLS handshake
- Fix approach: Implement actual CRL fetching/parsing or OCSP checking using cert.CRLDistributionPoints and cert.OCSPServer fields

**Large File Sizes:**
- Issue: Several files exceed recommended size limits (200-400 lines typical, 800 max)
- Files:
  - `pkg/client.go`: 858 lines
  - `pkg/types/errors.go`: 404 lines
  - `pkg/api/system.go`: 388 lines
  - `pkg/transport/websocket.go`: 327 lines
  - `pkg/connection/tls.go`: 319 lines
- Impact: Harder to test, review, and maintain
- Fix approach: Extract logical sub-modules (e.g., split client.go by manager initialization, errors.go by error category)

**TLSConfig Duplication:**
- Issue: Two separate TLSConfig structs exist in different packages with overlapping purpose
- Files: `pkg/transport/websocket.go:42-48` and `pkg/connection/tls.go:46-54`
- Impact: Confusion about which to use, potential for bugs when configuring TLS
- Fix approach: Consolidate into single TLSConfig in connection package, have transport reference it

**Large Client Struct:**
- Issue: The `client` struct in `pkg/client.go` (line 336) has 30+ fields including 15 API namespaces
- Files: `pkg/client.go:336-373`
- Impact: High coupling, hard to test individual components in isolation
- Fix approach: Use dependency injection with interfaces to decouple managers and APIs

**Deprecated AuthHandler Still Supported:**
- Issue: ClientOption `WithAuthHandler` is marked deprecated in favor of CredentialsProvider but both exist
- Files: `pkg/client.go:149` and `pkg/auth/handler.go`
- Impact: Users may use deprecated pattern, migration path unclear
- Fix approach: Add deprecation warning, document migration in README

## Known Bugs

**No bugs currently identified:**
- No active bug reports in codebase
- All tests pass

## Security Considerations

**InsecureSkipVerify Option:**
- Risk: `InsecureSkipVerify: true` in TLSConfig allows man-in-the-middle attacks
- Files: `pkg/transport/websocket.go:43`, `pkg/connection/tls.go:49`
- Current mitigation: Requires explicit opt-in, not the default
- Recommendations: Add warning log when InsecureSkipVerify is true, consider removing entirely for production use

**No Rate Limiting:**
- Risk: No built-in rate limiting on requests or connections
- Files: `pkg/managers/request.go`, `pkg/transport/websocket.go`
- Impact: Server may reject clients that send too many requests
- Recommendations: Add client-side rate limiting with configurable thresholds

**Credential Handling:**
- Risk: Auth credentials stored in memory
- Files: `pkg/auth/provider.go`, `pkg/auth/handler.go`
- Current mitigation: Credentials cleared on Close()
- Recommendations: Document that credentials remain in memory during operation; consider zeroing memory after use

## Performance Bottlenecks

**Buffered Channel Blocking Risk:**
- Problem: EventEmitTimeout of 200ms (default) may drop events under load
- Files: `pkg/managers/event.go`
- Cause: Fixed-size event channel with non-blocking send semantics
- Improvement path: Make channel size configurable, or use overflow handling

**Reconnect Manager Memory:**
- Problem: Timers created in reconnect loop (using time.NewTimer) could accumulate if not properly stopped
- Files: `pkg/managers/reconnect.go:93`
- Cause: Timer created per iteration, stopped on context cancellation but memory growth if leak occurs
- Improvement path: Already properly stopped - this is a caution rather than bug

**Large Pending Request Map:**
- Problem: RequestManager uses map without size limit for pending requests
- Files: `pkg/managers/request.go`
- Cause: Each in-flight request adds entry to map
- Improvement path: Add max pending requests limit with rejection/failure for overflow

## Fragile Areas

**Protocol Version Negotiation:**
- Files: `pkg/connection/protocol.go`
- Why fragile: Protocol version mismatch between client/server causes silent failures or unclear errors
- Safe modification: Always add new version checks before removing old ones; maintain backward compatibility
- Test coverage: Covered in protocol tests but edge cases may exist

**State Machine Transitions:**
- Files: `pkg/connection/state.go`
- Why fragile: Invalid state transitions could leave connection in inconsistent state
- Safe modification: Review all transition paths; ensure error is returned for invalid transitions
- Test coverage: 81.5% for connection package

**WebSocket Read Loop:**
- Files: `pkg/transport/websocket.go`
- Why fragile: Read errors may cause connection to close without proper cleanup
- Safe modification: Ensure all goroutine paths call proper cleanup; verify closed channel behavior
- Test coverage: 84.6% for transport package

## Scaling Limits

**Event Channel Buffer:**
- Current capacity: 100 events (configurable)
- Limit: Events dropped after timeout when channel full
- Scaling path: Increase EventBufferSize or implement overflow queue

**Concurrent Connection Limit:**
- Current capacity: Single WebSocket connection per client
- Limit: No multi-connection support
- Scaling path: Create connection pool or multi-client support if needed

**Pending Request Limit:**
- Current capacity: Unbounded map in RequestManager
- Limit: Memory growth with long-running connections and many requests
- Scaling path: Add configurable max pending requests

## Dependencies at Risk

**gorilla/websocket:**
- Risk: Single external dependency; any breaking change affects entire SDK
- Impact: WebSocket functionality completely dependent on this package
- Mitigation: Package is well-maintained and widely used
- Alternative: gobwas/ws if gorilla becomes unmaintained

## Missing Critical Features

**Graceful Degradation Under Load:**
- Problem: No built-in mechanism to prioritize critical events over non-critical ones
- Blocks: High-throughput scenarios where events may be dropped

**Connection Health Metrics:**
- Problem: No built-in metrics for connection quality, latency, or reliability
- Blocks: Observability beyond basic state transitions

**Retry Budget:**
- Problem: ReconnectManager has unlimited retries by default
- Blocks: Preventing infinite retry loops in pathological network conditions

## Test Coverage Gaps

**API Package Below Target:**
- What's not tested: Business logic in `pkg/api/*` modules (agents, channels, nodes, browser, config, etc.)
- Files: `pkg/api/*.go` (19 files)
- Risk: 49.2% coverage means approximately half of API code untested; bugs could go unnoticed
- Priority: High

**Examples Not Tested:**
- What's not tested: `examples/cmd/main.go` and `examples/server/main.go`
- Files: `examples/cmd/main.go`, `examples/server/main.go`
- Risk: Example code could break without detection; users rely on examples for onboarding
- Priority: Medium

**Transport WebSocket Tests:**
- What's not tested: Full WebSocket lifecycle under various network conditions
- Files: `pkg/transport/websocket_test.go`
- Risk: Edge cases in connection handling could fail in production
- Priority: Medium

**Request Manager Concurrency:**
- What's not tested: Race conditions in concurrent request handling
- Files: `pkg/managers/request_test.go`
- Risk: Data races could cause incorrect request/response correlation
- Priority: High

---

*Concerns audit: 2026-03-28*
