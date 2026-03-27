# Stack Research

**Domain:** Go WebSocket SDK Library Maturation
**Researched:** 2026-03-28
**Confidence:** MEDIUM

*Note: This research targets maturing an existing Go SDK. It does not re-examine the core stack (gorilla/websocket, Go runtime) which is already established and working.*

## Executive Summary

The project has a solid foundation: Go 1.24, gorilla/websocket v1.5.3, and a basic GoReleaser configuration. The gaps for "maturation" are in tooling depth -- fuzz tests exist but are shallow, no benchmarks exist, the goreleaser config is incomplete for semantic versioning, and no version tags exist yet. All findings are MEDIUM confidence based on official Go documentation and established Go ecosystem patterns.

## Recommended Stack

### Core Technologies

| Technology | Version | Purpose | Why Recommended |
|------------|---------|---------|-----------------|
| Go | 1.21+ (runtime), 1.24 (CI) | Language runtime | Already established; Go provides strong backward-compatibility guarantees for v1.x |
| gorilla/websocket | v1.5.3 | WebSocket client | Already in use; the de facto standard for Go WebSocket |
| Go modules | built-in | Dependency management | Already in use; go.mod/go.sum committed |
| GoReleaser | v2 | Release automation | Already configured; needs completion for proper library distribution |

### Fuzz Testing

| Tool | Version | Purpose | Why Recommended |
|------|---------|---------|-----------------|
| `testing.F` | built-in (Go 1.18+) | Fuzz testing framework | Standard library, no external dependency |
| `go test -fuzz` | built-in | Running fuzz targets | Native Go fuzzing, corpus management |
| golang.org/x/perf/cmd/benchstat | latest | Benchmark statistical comparison | Standard for comparing benchmark results |

**Current state:** Fuzz tests exist in `pkg/protocol/fuzz_test.go` but are shallow -- they only check for panics, not correctness. They do not actually validate that parsing produces correct results.

**What to add:**
- Fuzz corpus files in `testdata/fuzz/<Name>/` for persistent regression inputs
- Correctness assertions (not just panic checks) -- e.g., round-trip marshal/unmarshal
- `FuzzJSONMarshal` and `FuzzJSONUnmarshal` for the JSON serialization layer
- Integration with CI: `go test -fuzz=FuzzXxx -fuzztime=60s` (time-boxed, not infinite)

### Performance Benchmarking

| Tool | Version | Purpose | Why Recommended |
|------|---------|---------|-----------------|
| `testing.B` + `b.Loop()` | built-in (Go 1.24+) | Benchmark framework | Native Go benchmarking with new loop syntax |
| `benchstat` | golang.org/x/perf | Statistical benchmark comparison | Required for meaningful A/B comparisons |
| `pprof` | built-in | CPU/memory profiling | Standard Go profiling tool |

**What to add:**
- Benchmark files: `*_bench_test.go` per package
- Hot path benchmarks: transport write/read, event dispatch, request correlation
- Use `b.Loop()` (Go 1.24+) instead of `for range b.N`
- `b.ReportAllocs()` for allocation-sensitive paths
- `benchstat` in CI for regression detection

### Semantic Versioning

| Practice | Implementation | Notes |
|----------|----------------|-------|
| Start at v0 | `git tag v0.1.0` | No stability guarantees in v0 |
| First stable at v1 | `git tag v1.0.0` | Backward compatibility guarantees begin |
| Breaking changes | Major version bump + module path suffix (v2, v3) | Use `/v2` suffix for module path at v2+ |

**Go semantic import versioning rules:**
- v0.x.x: No guarantees
- v1.x.x: Backward compatible, no breaking changes to public API
- v2.x.x+: Incompatible changes require new module path ending in `/vN`

### Library Distribution

| Component | Current | Recommendation |
|-----------|---------|----------------|
| GoReleaser config | Partial | Add `blobs: true` for GitHub Releases hosting |
| Git tags | None | Add `v0.1.0`, progress to `v1.0.0` |
| Release notes | CHANGELOG.md exists | Configure git-cliff with GoReleaser integration |
| Module proxy | Not configured | Add `gomod.proxy: true` for verified builds |

**GoReleaser library mode configuration needed:**
```yaml
# Add to .goreleaser.yaml
archives:
  - format: zip
    # Library projects may want to skip binaries entirely

blobs:
  - provider: github
    # Enables GitHub Releases as blob host
```

## Installation

```bash
# Benchmark statistical comparison
go install golang.org/x/perf/cmd/benchstat@latest

# GoReleaser (if not already installed)
go install github.com/goreleaser/goreleaser@latest

# No additional dependencies needed -- using stdlib + existing deps
```

## Alternatives Considered

| Recommended | Alternative | When to Use Alternative |
|-------------|-------------|-------------------------|
| `testing.F` built-in fuzzing | `github.com/dvyukov/go-fuzz` | Only if needing older Go version (< 1.18) |
| GoReleaser | Manual git tagging + `go mod publish` | Only for very simple single-person projects |
| `benchstat` | Spreadsheet manual comparison | Only for one-off benchmarks, not recurring |
| `b.Loop()` (Go 1.24+) | `for range b.N` | Only if supporting Go < 1.24 |

## What NOT to Use

| Avoid | Why | Use Instead |
|-------|-----|-------------|
| `github.com/stretchr/testify` | Adds external dep for simple assertions | Use stdlib `require` + `assert` patterns or built-in |
| Custom versioning schemes | Breaks Go module proxy expectations | Follow semantic versioning with git tags |
| `gopkg.in` | Pre-modules era; not needed | Use `github.com/frisbee-ai/openclaw-sdk-go` with semver tags |
| Continuous fuzzing in CI (long-running) | Expensive, may timeout | Time-boxed fuzzing: `-fuzztime=60s` |

## Version Compatibility

| Package | Compatible With | Notes |
|---------|-----------------|-------|
| Go 1.24 (CI) | Go 1.21+ (runtime) | Project constraint; CI tests 1.24, runtime supports 1.21+ |
| gorilla/websocket v1.5.3 | Go 1.21+ | No special constraints |
| GoReleaser v2 | Go 1.18+ | Required for `gomod` directive support |

## Sources

- `testing.F` / fuzz testing: https://pkg.go.dev/testing -- HIGH confidence (official stdlib docs)
- `testing.B` / benchmarking with `b.Loop()`: https://pkg.go.dev/testing -- HIGH confidence (official stdlib docs)
- Go module versioning and semantic import versioning: https://go.dev/blog/using-go-modules -- HIGH confidence (official Go blog)
- Go release workflow: https://go.dev/doc/modules/release-workflow -- HIGH confidence (official docs)
- GoReleaser library mode: https://goreleaser.com -- MEDIUM confidence (official, but WebFetch blocked for full extraction)
- Current project files: `.planning/codebase/STACK.md`, `.goreleaser.yaml`, `pkg/protocol/fuzz_test.go` -- HIGH confidence (project-specific context)

---

*Stack research for: Go WebSocket SDK library maturation*
*Researched: 2026-03-28*
