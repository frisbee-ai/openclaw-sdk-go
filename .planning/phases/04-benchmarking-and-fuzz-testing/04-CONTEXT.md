# Phase 4: Benchmarking and Fuzz Testing - Context

**Gathered:** 2026-03-29
**Status:** Ready for planning

<domain>
## Phase Boundary

Phase 4 delivers **performance validation and regression detection for hot paths** for the OpenClaw SDK.

Three requirements:
- **TEST-01**: Hot-path benchmarks — `*_bench_test.go` files for transport write/read, event dispatch, request correlation using `b.Loop()`
- **TEST-02**: Fuzz test depth — round-trip correctness assertions in existing fuzz tests; corpus files in `testdata/fuzz/`
- **TEST-03**: Benchmark CI integration — `benchstat` in CI for regression detection on hot paths

</domain>

<decisions>
## Implementation Decisions

### TEST-01: Benchmark Scope (L1 + L2)

- **D-01 (Scope — L1 pure function layer):** Benchmark `protocol.Marshal()` and `protocol.Unmarshal()` JSON serialization — highest frequency operations, no external dependencies, stable and repeatable.
- **D-02 (Scope — L2 component layer):** Benchmark `EventManager.Emit()` and `RequestManager.SendRequest()` with mocked transport channels — covers manager hot-path logic without network IO.
- **D-03 (Excluded — L3 end-to-end):** NOT benchmarking full `client.SendRequest()` with mock WebSocket server — complexity/value mismatch. End-to-end transport tests belong in integration tests, not Phase 4 benchmark scope.
- **D-04 (File location):** `pkg/protocol/bench_test.go` for L1, `pkg/managers/event_bench_test.go` and `pkg/managers/request_bench_test.go` for L2. Follow existing co-location pattern.

### TEST-02: Fuzz Test Depth

- **D-05 (Round-trip assertions required):** Existing fuzz tests in `pkg/protocol/fuzz_test.go` currently only check that parsing does not panic. Add actual JSON round-trip correctness assertions:
  - `json.Unmarshal` → `json.Marshal` → `json.Unmarshal` → compare key fields
  - Compare `ID`, `Method`, `Action`, `Payload` fields between original and round-tripped frame
  - Report `t.Errorf` on mismatch (not panic — fuzz test should continue)
- **D-06 (Corpus files):** Create `testdata/fuzz/` directory with structured corpus files for each frame type. Each corpus file contains one valid example that exercises a specific code path.
- **D-07 (Corpus file format):** One JSON blob per file, named by the frame type and variant (e.g., `request_minimal.json`, `response_with_error.json`, `event_state_change.json`).
- **D-08 (Fuzz target coverage):** Ensure each `Fuzz*` function has corpus seeds covering: empty input, valid minimal, valid with payload, invalid JSON, large input (1MB+), special characters, null bytes.

### TEST-03: Benchstat CI Integration

- **D-09 (CI integration approach):** Add a separate `benchmark` job to `ci.yml` that runs on PR and block merge if regression detected.
- **D-10 (Comparison baseline):** PR branch benchmark vs main branch baseline. Store baseline `.bench` files as GitHub Actions artifacts (upload on main merge, download on PR).
- **D-11 (Regression threshold):** Block merge (fail CI) when BOTH conditions are met:
  - p-value < 0.05 (statistically significant difference)
  - AND regression > 10% (magnitude threshold)
  - Otherwise: warn but pass (avoid noise from CI environment variance)
- **D-12 (CI pipeline):**
  ```
  On PR: run benchmarks → benchstat vs stored baseline → pass/warn/fail
  On main merge: run benchmarks → upload .bench artifacts
  ```
- **D-13 (benchstat tool):** Install via `go install golang.org/x/perf/cmd/benchstat@latest` in CI before running.
- **D-14 (Go version):** Use same `go-version: ['1.24']` matrix as existing test job.

### TEST-04: Performance Metrics Reporting

- **D-15 (Built-in metrics):** All benchmarks report `ns/op`, `B/op`, `allocs/op` via Go's built-in benchmark framework (default when using `b.Loop()`).
- **D-16 (Custom ReportMetric):** Use `b.ReportMetric()` for SDK-specific metrics:
  - `bytes_per_frame` — JSON frame size in bytes (encoder overhead)
  - `channel_overhead_ns` — EventManager Emit() channel send overhead
  - `goroutine_count` — Goroutines active during benchmark (concurrency safety indicator)
- **D-17 (Reporter gate):** Only report metrics that map to SDK core value (fast serialization, low memory, predictable latency). Do NOT report vanity metrics.

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Phase Decisions (must be consistent)
- `.planning/phases/01-foundation/CONTEXT.md` — Phase 1 patterns: `testing` stdlib only, no external mock libraries
- `.planning/phases/02-observability/02-CONTEXT.md` — Phase 2 patterns: EventManager priority channels, ConnectionMetrics snapshot
- `.planning/phases/03-client-struct-refactor/03-CONTEXT.md` — Phase 3 patterns: client sub-struct layout, functional option pattern

### Project Requirements
- `.planning/ROADMAP.md` §Phase 4 — Phase 4 goal, success criteria (TEST-01, TEST-02, TEST-03)
- `.planning/REQUIREMENTS.md` §TEST-01, TEST-02, TEST-03 — Acceptance criteria for each test requirement
- `.planning/STATE.md` — Prior decisions and accumulated context

### Codebase Architecture
- `pkg/protocol/fuzz_test.go` — Existing fuzz tests needing round-trip assertions (TEST-02)
- `pkg/protocol/types.go` — `RequestFrame`, `ResponseFrame`, `EventFrame` types for round-trip comparison (TEST-02)
- `pkg/protocol/validation.go` — Existing validation functions
- `pkg/managers/event.go` — `EventManager.Emit()` hot path for L2 benchmark (TEST-01)
- `pkg/managers/request.go` — `RequestManager.SendRequest()` hot path for L2 benchmark (TEST-01)
- `pkg/transport/websocket.go` — WebSocket send/receive paths (L1 benchmark, TEST-01)
- `.github/workflows/ci.yml` — Existing CI pipeline needing benchmark job (TEST-03)

### Testing Patterns
- `.planning/codebase/TESTING.md` — Established testing conventions: stdlib `testing`, table-driven, no external mocks, sync patterns

### No external specs beyond Go stdlib
No external specs — all requirements captured in decisions above.

</canonical_refs>

<codebase_context>
## Existing Code Insights

### Reusable Assets
- **`pkg/protocol/fuzz_test.go`**: Already exists with 7 fuzz targets — needs round-trip assertions added to each
- **`protocol.NewRequestFrame()`/`NewResponseFrameSuccess()`**: Factory functions for creating test frames for corpus
- **`EventManager`** (`pkg/managers/event.go:19`): Already has priority channels (HIGH/MED/LOW) — benchmark `Emit()` dispatch loop
- **`RequestManager`** (`pkg/managers/request.go:30`): Has `pending` map + mutex — benchmark `SendRequest`/`HandleResponse` correlation
- **`WebSocketTransport`** (`pkg/transport/websocket.go:65`): Has `sendCh`/`recvCh` buffered channels — benchmark channel send/receive

### Established Patterns
- **`b.Loop()` pattern** (Go 1.24+): Use `for _, n := range []int{1, 100, 1000, 1000000} { b.N = n; b.ResetTimer(); ... }` for scaling benchmarks
- **Pure function benchmarks**: `json.Marshal`/`json.Unmarshal` are pure — no setup needed beyond `b.ReportAllocsPerOp()`
- **Coroutine safety**: Benchmarks should not introduce race conditions — use `b.RunParallel()` for concurrent benchmarks
- **Standard test file naming**: `*_bench_test.go` co-located with source files

### Integration Points
- **`ci.yml`**: Add new `benchmark` job parallel to `test` job. Requires `benchstat` tool installation step.
- **`testdata/fuzz/`**: New directory at repo root for corpus files. Each JSON file seeds one fuzz target.
- **`pkg/protocol/fuzz_test.go`**: Each `Fuzz*` function needs: corpus seeds → round-trip assertions → `t.Errorf` on mismatch

</codebase_context>

<specifics>
## Specific Ideas

- **Round-trip comparison fields**: Only compare fields that round-trip through JSON (ID, Method, Action, Payload). Internal-only fields (pointers, unexported fields) are excluded from comparison.
- **Corpus file naming convention**: `{FrameType}_{variant}.json`, e.g., `request_minimal.json`, `request_large_payload.json`, `response_success.json`, `response_error.json`, `event_connect.json`, `event_error.json`.
- **Regression threshold rationale**: p<0.05 AND >10% — pure statistical significance without p-value is too noisy in CI environments; pure percentage without p-value catches trivial differences. Both gates required.
- **Artifact storage**: Use `actions/upload-artifact@v6` and `actions/download-artifact@v6` with path `benchmark-results/*.bench`.

</specifics>

<deferred>
## Deferred Ideas

### Reviewed Todos (not folded)
None — no pending todos matched Phase 4 scope.

### Out-of-Scope (discussed but not this phase)
- **TEST-01 — L3 end-to-end benchmarks**: Requires mock WebSocket server — complexity too high for Phase 4 benchmark infrastructure. Deferred to future integration test work.
- **Prometheus/OpenTelemetry export**: Phase 5 REL mentions this. Internal benchmark infrastructure is the foundation.
- **Continuous benchmark profiling**: pprof integration for CPU/memory profiling on benchmark regressions — future Phase 4 extension.

</deferred>

---

*Phase: 04-benchmarking-and-fuzz-testing*
*Context gathered: 2026-03-29*
