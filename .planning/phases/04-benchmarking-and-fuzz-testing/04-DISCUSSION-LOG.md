# Phase 4: Benchmarking and Fuzz Testing - Discussion Log

> **Audit trail only.** Do not use as input to planning, research, or execution agents.
> Decisions are captured in CONTEXT.md — this log preserves the alternatives considered.

**Date:** 2026-03-29
**Phase:** 04-benchmarking-and-fuzz-testing
**Areas discussed:** Benchmark scope, Fuzz test depth, benchstat CI strategy, Performance metrics reporting

---

## Area 1: Benchmark Scope

| Option | Description | Selected |
|--------|-------------|----------|
| L1 only | Pure function layer: protocol.Marshal/Unmarshal | |
| L1 + L2 | Pure function + component layer: L1 + EventManager.Emit(), RequestManager.SendRequest() with mocked channels | ✓ |
| L1 + L2 + L3 | Full end-to-end with mock WebSocket server | |

**User's choice:** L1 + L2
**Notes:** L3 (end-to-end) excluded — requires mock WebSocket server, complexity/value mismatch for Phase 4 benchmark infrastructure. Focus on highest-frequency operations (JSON serialization, channel dispatch) that are stable and repeatable.

---

## Area 2: Fuzz Test Depth

| Option | Description | Selected |
|--------|-------------|----------|
| Keep panic-only tests | Existing fuzz tests only check no panics occur | |
| Add round-trip assertions | JSON marshal → unmarshal → marshal → compare fields | ✓ |

**User's choice:** Add round-trip assertions
**Notes:** Current fuzz tests have a critical gap — data is never used for actual validation. Round-trip assertions ensure JSON encoding/decoding preserves data correctly. Compare ID, Method, Action, Payload fields. Report via t.Errorf (not panic).

---

## Area 3: benchstat CI Strategy

| Option | Description | Selected |
|--------|-------------|----------|
| A — GitHub PR comment | Compare vs main, post comment | |
| B — Main branch baseline artifact | Upload .bench on main merge, download on PR | |
| C — PR vs main, fail block | Run benchstat on PR, block merge if regression | ✓ |
| D — Local only | No CI integration | |

**User's choice:** C — PR vs main, fail block
**Regression threshold:** p < 0.05 AND regression > 10% — both conditions required to block. Warn but pass if only one condition met.
**Notes:** Statistical significance without magnitude = too noisy; magnitude without significance = catching trivial diffs. Both gates required for signal quality.

---

## Area 4: Performance Metrics Reporting

| Option | Description | Selected |
|--------|-------------|----------|
| Built-in only | ns/op, B/op, allocs/op via Go built-in | |
| Built-in + ReportMetric | Add custom metrics: bytes_per_frame, channel_overhead_ns | ✓ |

**User's choice:** Built-in + ReportMetric
**Metrics to report:**
- `ns/op`, `B/op`, `allocs/op` — built-in, standard
- `bytes_per_frame` — JSON frame size (encoder overhead)
- `channel_overhead_ns` — EventManager Emit() channel send overhead
- `goroutine_count` — concurrency safety indicator

**Notes:** Only report metrics that map to SDK core value. No vanity metrics.

---

## Claude's Discretion

All four areas had clear decisions from user. No areas delegated to Claude.

## Deferred Ideas

None — discussion stayed within phase scope.

