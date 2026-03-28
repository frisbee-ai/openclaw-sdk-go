---
phase: 02
slug: observability
status: draft
nyquist_compliant: false
wave_0_complete: false
created: 2026-03-28
---

# Phase 02 — Validation Strategy

> Per-phase validation contract for feedback sampling during execution.

---

## Test Infrastructure

| Property | Value |
|----------|-------|
| **Framework** | Go built-in `testing.T` |
| **Config file** | None — standard Go test layout |
| **Quick run command** | `go test ./pkg/... -run "Test(OBS|Event|Priority|Metrics|Timeout|Reconnect)" -v -count=1` |
| **Full suite command** | `go test ./pkg/... -race -count=1` |
| **Estimated runtime** | ~5 seconds |

---

## Sampling Rate

- **After every task commit:** `go test ./pkg/... -run "<task tests>" -v -count=1`
- **After every plan wave:** `go test ./pkg/... -race -count=1`
- **Before `/gsd:verify-work`:** Full suite must be green
- **Max feedback latency:** 10 seconds

---

## Per-Task Verification Map

| Task ID | Plan | Wave | Requirement | Test Type | Automated Command | File | Status |
|---------|------|------|-------------|-----------|-------------------|------|--------|
| 02-01-01 | OBS-01 | 1 | OBS-01 | unit | `go test -run TestReconnectManager_AttemptCount ./pkg/managers/ -v -race` | `pkg/managers/reconnect_maxretries_test.go` | pending |
| 02-01-02 | OBS-01 | 1 | OBS-01 | unit | `go test -run TestTickMonitor_GetTickIntervalMs ./pkg/events/ -v` | `pkg/events/tick_test.go` | pending |
| 02-01-03 | OBS-01 | 1 | OBS-01 | unit | `go test -run TestGetMetrics_ ./pkg/ -v` | `pkg/client_test.go` | pending |
| 02-01-04 | OBS-01 | 1 | OBS-01 | unit | `go test -run TestEvent_DefaultPriority ./pkg/ -v` | `pkg/types/types_test.go` | pending |
| 02-02-01 | OBS-02 | 2 | OBS-02 | unit | `go test -run TestSendRequest_WithRequestTimeout ./pkg/ -v` | `pkg/client_test.go` | pending |
| 02-02-02 | OBS-02 | 2 | OBS-02 | unit | `go test -run TestSendRequest_NoOptions ./pkg/ -v` | existing | pending |
| 02-02-03 | OBS-02 | 2 | OBS-02 | unit | `go test -run TestSendRequest_TimeoutOverwritesCtxDeadline ./pkg/ -v` | `pkg/client_test.go` | pending |
| 02-03-01 | OBS-03 | 2 | OBS-03 | unit | `go test -run TestEventManager_PriorityHighNeverDrops ./pkg/managers/ -v -race` | `pkg/managers/event_priority_test.go` | pending |
| 02-03-02 | OBS-03 | 2 | OBS-03 | unit | `go test -run TestEventManager_PriorityDropOrder ./pkg/managers/ -v -race` | `pkg/managers/event_priority_test.go` | pending |
| 02-03-03 | OBS-03 | 2 | OBS-03 | unit | `go test -run TestEventManager_PriorityAssignment ./pkg/managers/ -v` | `pkg/managers/event_priority_test.go` | pending |
| 02-03-04 | OBS-04 | 2 | OBS-04 | unit | `go test -run TestWithEventBufferSize ./pkg/ -v` | existing (verify) | pending |

*Status: pending · green · red · flaky*

---

## Wave 0 Requirements

- [ ] `pkg/client_test.go` — stubs for `TestGetMetrics_*` and `TestSendRequest_WithRequestTimeout`
- [ ] `pkg/managers/event_priority_test.go` — stubs for `TestEventManager_Priority*` tests
- [ ] `pkg/events/tick_test.go` — stub for `TestTickMonitor_GetTickIntervalMs`
- [ ] `pkg/types/types_test.go` — stub for `TestEvent_DefaultPriority`
- [ ] `pkg/managers/reconnect_maxretries_test.go` — already exists, add `TestReconnectManager_AttemptCount`

*If none: "Existing infrastructure covers all phase requirements."*

---

## Manual-Only Verifications

| Behavior | Requirement | Why Manual | Test Instructions |
|----------|-------------|------------|-------------------|
| Tick-based latency accuracy | OBS-01 | Requires simulated tick stream | Create test with fake tick times, verify `GetMetrics().Latency` matches expected formula |
| Priority drop under memory pressure | OBS-03 | Hard to reproduce deterministically | Verify drop behavior via test that fills all priority channels |

*If none: "All phase behaviors have automated verification."*

---

## Validation Sign-Off

- [ ] All tasks have automated verify or Wave 0 dependencies
- [ ] Sampling continuity: no 3 consecutive tasks without automated verify
- [ ] Wave 0 covers all MISSING references
- [ ] No watch-mode flags
- [ ] Feedback latency < 10s
- [ ] `nyquist_compliant: true` set in frontmatter

**Approval:** pending
