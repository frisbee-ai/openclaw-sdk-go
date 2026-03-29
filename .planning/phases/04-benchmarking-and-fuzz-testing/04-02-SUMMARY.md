---
phase: 04-benchmarking-and-fuzz-testing
plan: "02"
subsystem: ci
tags: [benchmark, ci, benchstat, performance]
dependency_graph:
  requires: []
  provides: [TEST-03]
  affects: [.github/workflows/ci.yml]
tech_stack:
  added: [golang.org/x/perf/cmd/benchstat, actions/upload-artifact@v6, actions/download-artifact@v6]
  patterns: [statistical regression detection, artifact versioning]
key_files:
  created:
    - .planning/phases/04-benchmarking-and-fuzz-testing/04-02-SUMMARY.md
    - benchmark-results/.gitkeep
  modified:
    - .github/workflows/ci.yml
decisions:
  - id: D-09
    text: "Separate benchmark job parallel to test job"
  - id: D-10
    text: "PR branch benchmark vs main branch baseline. Store baseline .bench files as GitHub Actions artifacts"
  - id: D-11
    text: "Regression threshold - block merge when p<0.05 AND >10%"
  - id: D-12
    text: "On PR: run benchmarks → benchstat vs stored baseline → pass/warn/fail; On main: upload .bench artifacts"
  - id: D-13
    text: "Install benchstat via go install golang.org/x/perf/cmd/benchstat@latest"
  - id: D-14
    text: "Use same go-version: ['1.24'] as test job"
metrics:
  duration: "~5 minutes"
  completed: "2026-03-29"
---

# Phase 04 Plan 02: Benchstat CI Integration Summary

## One-liner

Add benchstat CI job for hot-path performance regression detection with p<0.05 AND >10% threshold blocking.

## What Was Built

Added `benchmark` job to `.github/workflows/ci.yml` that:

- Runs hot-path benchmarks (`go test -bench=. -benchmem -run=^$ -count=1 ./pkg/protocol/ ./pkg/managers/`)
- Installs benchstat via `go install golang.org/x/perf/cmd/benchstat@latest`
- On PR: downloads baseline artifact, runs `benchstat baseline.txt bench.txt`, blocks merge if p<0.05 AND regression>10%
- On main push: uploads `benchmark-baseline` artifact for future PR comparisons
- Uses `actions/upload-artifact@v6` and `actions/download-artifact@v6` with 30-day retention

## Acceptance Criteria Verified

| Criterion | Status |
|-----------|--------|
| `.github/workflows/ci.yml` contains `benchmark:` | PASS |
| `.github/workflows/ci.yml` contains `benchstat` | PASS |
| `.github/workflows/ci.yml` contains `benchmark-results` | PASS |
| `.github/workflows/ci.yml` contains `actions/upload-artifact@v6` | PASS |
| `.github/workflows/ci.yml` contains `actions/download-artifact@v6` | PASS |
| Regression threshold: p<0.05 AND >10% blocking logic | PASS |
| benchstat installable and runs | PASS |

## Deviations from Plan

None - plan executed exactly as written.

## Self-Check: PASSED

- `.github/workflows/ci.yml` modified with new `benchmark` job
- `benchmark-results/.gitkeep` created
- benchstat installed and verified working
- Regression threshold logic implemented: `p<0.05 AND >10%` blocks merge
