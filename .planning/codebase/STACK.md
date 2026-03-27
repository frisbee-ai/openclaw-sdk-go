# Technology Stack

**Analysis Date:** 2026-03-28

## Languages

**Primary:**
- Go 1.24 - Core SDK implementation

## Runtime

**Environment:**
- Go standard runtime (no external runtime dependencies)
- Pure Go - no CGO required

**Package Manager:**
- Go modules (go.mod/go.sum)
- Lockfile: present

## Frameworks

**Core:**
- No framework - pure library SDK
- WebSocket: `github.com/gorilla/websocket v1.5.3` - WebSocket client implementation

**Testing:**
- Go testing package (built-in `testing.T`)
- Fuzz testing via `testing.F` in `pkg/protocol/fuzz_test.go`

**Build/Dev:**
- GoReleaser v2 - Release automation for library distribution
- pre-commit hooks - Local hooks for format/vet/lint/test

## Key Dependencies

**Critical:**
- `github.com/gorilla/websocket v1.5.3` - WebSocket protocol implementation
  - Used in: `pkg/transport/websocket.go`
  - Purpose: All WebSocket connectivity

**No other external dependencies** - SDK uses only Go standard library + gorilla/websocket

## Configuration

**Environment:**
- No runtime environment configuration
- SDK is configured programmatically via option pattern

**Build:**
- `.golangci.yaml` - Linter configuration (excludes test files from errcheck/unused)
- `.goreleaser.yaml` - Library release configuration (no binary builds)
- `.pre-commit-config.yaml` - Local hooks: gofmt, go vet, golangci-lint, go test

## Platform Requirements

**Development:**
- Go 1.24+
- golangci-lint (for linting)
- pre-commit (optional, for local hooks)

**Production:**
- Go 1.21+ (runtime compatibility)
- No platform-specific requirements

## Tooling

| Tool | Version/Config | Purpose |
|------|---------------|---------|
| gofmt | built-in | Code formatting |
| go vet | built-in | Static analysis |
| golangci-lint | latest via GitHub Action | Fast linting with auto-fix |
| GoReleaser | v2 | Library release automation |
| codecov | v6 | Coverage reporting |

## CI/CD Pipeline

**GitHub Actions (`.github/workflows/ci.yml`):**
1. Checkout code
2. Set up Go 1.24
3. Download/verify dependencies
4. Run go fmt
5. Run go vet
6. Run golangci-lint
7. Run tests with race detection + coverage
8. Upload coverage to Codecov

**Go Version Matrix:**
- `['1.24']` - Only Go 1.24 tested in CI

---

*Stack analysis: 2026-03-28*
