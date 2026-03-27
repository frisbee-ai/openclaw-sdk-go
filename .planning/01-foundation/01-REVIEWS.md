---
phase: 1
reviewers: [codex]
reviewed_at: 2026-03-28T02:20:00Z
plans_reviewed: [01-PLAN.md, 02-PLAN.md, 03-PLAN.md]
---

# Cross-AI Plan Review — Phase 1

## Codex Review

# Plan 01 Review

## Summary
Plan 01 is a reasonable contract-first foundation, but it currently mixes three different concerns with different compatibility profiles: request throttling, reconnect semantics, and new errors. The biggest issue is not implementation complexity, but API shape: adding `MaxRetries` alongside the existing `MaxAttempts` in `pkg/types/types.go:63` risks creating duplicate public semantics before the later waves even begin.

## Strengths
- Introduces shared contracts before wiring behavior, which is the right order for the later waves.
- Keeps reusable definitions in `pkg/types`, matching the current package structure.
- Includes tests at the contract layer, which should catch default-value regressions early.
- Puts the new pending-request and reconnect errors in a central place rather than scattering string literals.

## Concerns
- `[HIGH]` `MaxRetries` overlaps with existing `MaxAttempts` on the public `ReconnectConfig` type in `pkg/types/types.go:63`. Without a deprecation/migration story, Plan 03 will inherit ambiguous precedence rules.
- `[MEDIUM]` Sentinel errors do not match the repo's current typed-error pattern in `pkg/types/errors.go:223`. Plain sentinels alone will be weaker than the existing `ReconnectError` / `RequestError` model.
- `[MEDIUM]` `RequestRateLimiter.Allow() bool` is probably too narrow for FOUND-01. It cannot represent waiting, context cancellation, or retry-after information, so it may force local rejection instead of pacing.
- `[LOW]` The plan does not define zero-value semantics for limiter configuration, which matters for backward compatibility and easy adoption.

## Suggestions
- Pick one reconnect budget field. If `MaxRetries` is the future API, make `MaxAttempts` a deprecated alias or derive one from the other explicitly.
- Make the limiter contract context-aware, for example `Wait(ctx)` or `Allow() error`, so the SDK can pace instead of only reject.
- Add explicit thread-safety expectations for limiter implementations.
- Integrate new conditions with the existing typed error system, not only `var` sentinels.

## Risk Assessment
**MEDIUM**. The implementation itself is not hard, but a weak public contract here will create avoidable confusion and rework in Plans 02 and 03.

---

# Plan 02 Review

## Summary
Plan 02 targets the right code paths, but it underestimates the current concurrency model. The request path is already serialized by the client mutex, and the request manager has channel-lifecycle hazards. Without addressing those first, rate limiting and pending limits may appear to work in tests while still missing real edge cases.

## Strengths
- Places pending-limit enforcement close to the pending map, which is the correct ownership boundary.
- Chooses immediate failure on limit breach, matching FOUND-04.
- Keeps request behavior changes mostly isolated to request/client code.
- Includes tests in both manager and client layers.

## Concerns
- `[HIGH]` Public requests are currently serialized by the client-wide lock in `pkg/client.go:563`. Because `SendRequest` holds `c.mu` while waiting for the response, client-level concurrency tests can easily miss real pending-map pressure and rate-limit behavior.
- `[HIGH]` `RequestManager` has a double-close risk today: `SendRequest` cleanup closes `respCh` in `pkg/managers/request.go:70`, while `Clear` and `Close` also close it in `pkg/managers/request.go:162` and `pkg/managers/request.go:185`. Adding more state around the same map increases race risk.
- `[MEDIUM]` `Default pending limit = 256` is arbitrary and may be too low for batch-heavy users. Hard-coding it without a public option is likely to cause friction.
- `[MEDIUM]` If the limiter remains `Allow() bool`, the SDK will reject bursts locally instead of smoothing them, which does not cleanly satisfy the "prevent server rejection under load" success criterion.
- `[LOW]` The plan does not mention duplicate request IDs or concurrent overflow tests.

## Suggestions
- Reduce the scope of the client mutex in `SendRequest`, or snapshot the needed state under lock and wait outside it.
- Fix request-manager channel ownership before adding limits. One side should own channel close; not both.
- Make the pending limit configurable through `ClientConfig` / `With...` option, with a documented default.
- Add race-detector tests that launch more than the limit concurrently and verify fast failure without panics.
- Decide whether local rate limiting should block, reject, or support both modes.

## Risk Assessment
**HIGH**. The plan is aimed at the right features, but current locking and channel ownership make this area fragile. If implemented narrowly, it may pass tests and still leave concurrency bugs.

---

# Plan 03 Review

## Summary
Plan 03 is the weakest of the three as scoped. The goals are valid, but the listed files do not cover the live TLS connection path, and the reconnect budget is being added on top of a reconnect manager that currently starts retrying immediately after a healthy connect. As written, this plan is unlikely to fully achieve FOUND-02, FOUND-03, or FOUND-05.

## Strengths
- Correctly treats CRL validation as either real security work or an explicit limitation, not a fake implementation.
- Groups reconnect and TLS hardening together, which is sensible from a production-safety standpoint.
- Calls out precedence between old and new retry settings, which is necessary if both fields exist temporarily.

## Concerns
- `[HIGH]` `ClientConfig.TLSConfig` is not wired into the connection path. `ConnectionManager.Connect` calls `transport.Dial(..., nil)` in `pkg/managers/connection.go:73`, so changing only `pkg/connection/tls.go` will not affect real connections.
- `[HIGH]` The insecure warning already happens in transport via direct `stderr` output in `pkg/transport/websocket.go:109`. A `SetLogger` on `TlsValidator` is insufficient unless client, manager, and transport are also rewired.
- `[HIGH]` Reconnect starts immediately after a successful `Connect()` in `pkg/client.go:507`, and `ReconnectManager.Start()` always launches the loop in `pkg/managers/reconnect.go:74`. That means retries can be consumed without any disconnect event.
- `[MEDIUM]` "MaxRetries takes precedence over MaxAttempts" is too vague without exact nil/zero/negative semantics.
- `[MEDIUM]` The repo already has a TODO-style revocation stub. Replacing it with another comment does not materially improve FOUND-03 unless the limitation is exposed in docs or behavior.

## Suggestions
- Expand the plan scope to include `pkg/client.go`, `pkg/managers/connection.go`, and `pkg/transport/websocket.go`.
- Pass `TLSConfig` and `Logger` through the actual dial path, then log the insecure warning through `types.Logger` once per connection attempt.
- Rework reconnect triggering so retries start only after a disconnect/failure signal, not after initial healthy connect.
- Return a typed reconnect error when the budget is exhausted, consistent with existing error patterns.
- Add tests for:
  - no reconnect attempts while the initial connection remains healthy
  - warning is emitted at connection time when `InsecureSkipVerify=true`
  - retry budget stops at the configured count and returns the expected error

## Risk Assessment
**HIGH**. The plan's intent is right, but the implementation scope is incomplete and the current reconnect behavior is already semantically wrong for a retry-budget feature.

---

# Overall

Plan 01 is mostly solid if the public contract is tightened. Plan 02 needs concurrency and lifecycle issues addressed, not just feature wiring. Plan 03 needs scope expansion; otherwise it will not hit the phase success criteria in the live code path.

The main cross-plan risk is hidden overlap in Wave 2. Both later plans touch request/connect behavior, and Plan 03 really needs files currently absent from its scope. I would revise the sequence to make Plan 01 finalize the public config/error model first, then split Wave 2 into:
1. request-path concurrency and limits
2. reconnect/TLS wiring through the real dial path

---

## Consensus Summary

**Note:** Only one reviewer (Codex) was available for this cross-AI review. The following represents Codex's independent assessment.

### Agreed Strengths
- Contract-first approach in Plan 01 (shared contracts before wiring behavior)
- Test coverage at multiple layers (manager and client)
- Centralized error definitions in pkg/types
- Correct treatment of CRL as either real work or explicit limitation

### Key Concerns (Priority Order)

**HIGH Priority — Must address before execution:**
1. **Plan 03 scope gap**: Files `pkg/client.go`, `pkg/managers/connection.go`, `pkg/transport/websocket.go` are missing but required for TLS warning and reconnect budget to work in actual code path
2. **Plan 02 channel ownership risk**: `RequestManager` has double-close hazard that will be exacerbated by adding pending limits
3. **Plan 01 MaxRetries/MaxAttempts overlap**: Two public fields with similar semantics creates API ambiguity without clear precedence rules
4. **Plan 02 client mutex scope**: `SendRequest` holds lock while waiting for response, which masks real concurrency pressure in tests

**MEDIUM Priority — Should address:**
5. Sentinel errors don't match existing typed-error pattern in pkg/types/errors.go
6. `RequestRateLimiter.Allow() bool` too narrow for pacing vs rejection
7. Hard-coded pending limit (256) without public option may cause user friction
8. Reconnect starts immediately after healthy connect, consuming retries without disconnect events

### Divergent Views
N/A — only one reviewer participated

### Recommended Actions

**Before /gsd:execute-phase:**
1. Expand Plan 03 scope to include `pkg/client.go`, `pkg/managers/connection.go`, `pkg/transport/websocket.go`
2. Fix `RequestManager` channel ownership in Plan 02 (one side owns close, not both)
3. Clarify MaxRetries vs MaxAttempts precedence with exact nil/zero/negative semantics in Plan 01
4. Reduce client mutex scope in Plan 02 or snapshot state under lock and wait outside
5. Consider making `RequestRateLimiter` context-aware (`Allow(ctx) error` instead of `Allow() bool`)

**To incorporate feedback:**
```bash
/gsd:plan-phase 1 --reviews
```

---

*Cross-AI review completed: 2026-03-28*
*Reviewer: Codex (1 of 3 CLIs available — Gemini missing, Claude skipped for independence)*
