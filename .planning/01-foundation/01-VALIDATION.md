---
phase: 1
slug: foundation-hardening
status: draft
nyquist_compliant: false
wave_0_complete: false
created: 2026-03-28
---

# Phase 1 — Validation Strategy

> Per-phase validation contract for feedback sampling during execution.

---

## Test Infrastructure

| Property | Value |
|----------|-------|
| **Framework** | Go testing package (built-in) |
| **Config file** | none — uses `go test` |
| **Quick run command** | `go test ./pkg/... -run TestFoundations -v` |
| **Full suite command** | `go test ./pkg/... -race -cover` |
| **Estimated runtime** | ~10 seconds |

---

## Sampling Rate

- **After every task commit:** Run `go test ./pkg/... -run TestFoundations -v`
- **After every plan wave:** Run `go test ./pkg/... -race -cover`
- **Before `/gsd:verify-work`:** Full suite must be green
- **Max feedback latency:** 15 seconds

---

## Per-Task Verification Map

| Task ID | Plan | Wave | Requirement | Test Type | Automated Command | File Exists | Status |
|---------|------|------|-------------|-----------|-------------------|-------------|--------|
| 01-01-01 | 01 | 1 | FOUND-01 | unit | `go test ./pkg/types -run TestRateLimiter` | ✅ exists | ⬜ pending |
| 01-01-02 | 01 | 1 | FOUND-02 | unit | `go test ./pkg/managers -run TestMaxRetries` | ✅ exists | ⬜ pending |
| 01-01-03 | 01 | 1 | FOUND-03 | unit | `go test ./pkg/connection -run TestCRLStub` | ✅ exists | ⬜ pending |
| 01-01-04 | 01 | 1 | FOUND-04 | unit | `go test ./pkg/managers -run TestPendingLimit` | ✅ exists | ⬜ pending |
| 01-01-05 | 01 | 1 | FOUND-05 | unit | `go test ./pkg/connection -run TestInsecureWarning` | ✅ exists | ⬜ pending |

*Status: ⬜ pending · ✅ green · ❌ red · ⚠️ flaky*

---

## Wave 0 Requirements

- [ ] `pkg/types/rate_limiter.go` — `RequestRateLimiter` interface, `TokenBucketLimiter` struct
- [ ] `pkg/types/errors.go` — `ErrTooManyPendingRequests`, `ErrMaxRetriesExceeded`
- [ ] `pkg/managers/request_test.go` — existing test infrastructure
- [ ] `pkg/managers/reconnect_test.go` — existing test infrastructure
- [ ] `pkg/connection/tls_test.go` — existing test infrastructure

*Existing infrastructure covers all phase requirements. All packages have test files.*

---

## Manual-Only Verifications

| Behavior | Requirement | Why Manual | Test Instructions |
|----------|-------------|------------|-------------------|
| TLS CRL stub comment clarity | FOUND-03 | Human judgment | Review `tls.go` comment for clarity on limitation |
| Log output visibility in production | FOUND-05 | Environment-dependent | Check logs when TLS connects with InsecureSkipVerify=true |

*Manual verifications are minimal — most behaviors have automated tests.*

---

## Validation Sign-Off

- [ ] All tasks have `<automated>` verify or Wave 0 dependencies
- [ ] Sampling continuity: no 3 consecutive tasks without automated verify
- [ ] Wave 0 covers all MISSING references
- [ ] No watch-mode flags
- [ ] Feedback latency < 15s
- [ ] `nyquist_compliant: true` set in frontmatter

**Approval:** pending

---

*Phase: 01-foundation-hardening*
*Validation strategy: 2026-03-28*
