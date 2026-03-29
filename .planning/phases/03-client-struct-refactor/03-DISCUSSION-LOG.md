# Phase 3: Client Struct Refactor - Discussion Log

> **Audit trail only.** Do not use as input to planning, research, or execution agents.
> Decisions are captured in CONTEXT.md — this log preserves the alternatives considered.

**Date:** 2026-03-29
**Phase:** 03-client-struct-refactor
**Areas discussed:** API sub-struct, remaining fields grouping, Close vs Disconnect semantics

---

## API Sub-struct

### Area: API namespace fields grouping

| Option | Description | Selected |
|--------|-------------|----------|
| Nested `api struct { ... }` | Group 15 API fields into nested sub-struct, accessor methods unchanged (delegate to c.api.*) | ✓ |
| Keep flat fields | All 15 API fields as top-level struct fields | |

**User's choice:** Nested `api struct { ... }`
**Notes:** Consistent with existing `managers` sub-struct pattern. Backward compatible (accessor methods unchanged). Research confirmed no reflection/encoding/mutex issues with nested sub-structs in Go.

---

## Remaining Fields Grouping

### Area: protocol and health field organization

| Option | Description | Selected |
|--------|-------------|----------|
| 4 sub-structs (managers + api + protocol + health) | Group by concern; core fields remain flat | ✓ |
| Only `health` sub-struct | Group tick/gap fields only; protocol fields remain flat | |
| Keep all flat | No additional grouping beyond existing managers | |

**User's choice:** 4 sub-structs (managers + api + protocol + health)
**Notes:** protocol group = {protocolNegotiator, policyManager, serverInfo, snapshot}. health group = {tickMonitor, gapDetector, tickHandlerUnsub}. core flat = {config, ctx, cancel, mu, requestFn}.

---

## Close vs Disconnect Semantics

### Area: Disconnect() documentation and existence

| Option | Description | Selected |
|--------|-------------|----------|
| Document clarification (recommended) | Keep both methods; update interface docs to clearly distinguish: Close=full shutdown, Disconnect=disconnect without shutdown | ✓ |
| Merge into one method | Remove Disconnect(); only Close() exists | |

**User's choice:** Document clarification
**Notes:** Keep both. Update OpenClawClient interface docs. Disconnect() comment should say "disconnects without shutting down client; call Connect() to reconnect."

---

## Claude's Discretion

Delegated to researcher/planner:
- Exact initialization order in NewClient for the 4 sub-structs
- Whether to use embedded types vs non-embedded sub-structs
- Test strategy for the refactor (verify all fields accessible after restructuring)

---

## Deferred Ideas

None — all scope items addressed in discussion.

