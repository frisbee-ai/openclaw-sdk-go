# Coding Conventions

**Analysis Date:** 2026-03-28

## Naming Patterns

**Files:**
- Go source files: lowercase with underscores (`websocket_test.go`, `event_manager.go`)
- Test files: `*_test.go` suffix co-located with source
- Package name matches directory name

**Functions:**
- Exported functions: `PascalCase` (`NewClient`, `SendRequest`)
- Unexported functions: `camelCase` (`buildConnectParams`, `processServerInfo`)
- Test functions: `PascalCase` with descriptive names (`TestEventManager_Subscribe`)

**Variables:**
- Local variables: `camelCase` (`connectParams`, `authHandler`)
- Constants: `SCREAMING_SNAKE_CASE` for actual constants (`StateDisconnected`, `EventConnect`)
- Error variables: `err` prefix (`err`, `errCh`)
- Channel variables: `Ch` suffix or descriptive (`done`, `errCh`, `recvCh`)

**Types:**
- Structs: `PascalCase` (`ClientConfig`, `EventManager`)
- Interfaces: `PascalCase` with `er` suffix where idiomatic (`OpenClawClient`, `Transport`)
- Type aliases: preserve pattern of underlying type

## Code Style

**Formatting:**
- Tool: `gofmt` (runs automatically via pre-commit hook)
- Config: `.golangci.yaml` with `golangci-lint`
- Pre-commit hook: `gofmt -w -l .`

**Linting:**
- Tool: `golangci-lint run --fix --timeout=5m`
- Exclusions: Test files (`_test.go`) exclude `errcheck` and `unused` linters
- Config: `.golangci.yaml`

**Indentation:**
- Standard Go formatting (gofmt handles this)

**Line Length:**
- No explicit limit; gofmt handles wrapping

## Import Organization

**Order:**
1. Standard library packages (`context`, `sync`, `time`, `errors`)
2. External packages (`github.com/gorilla/websocket`)
3. Internal packages (`github.com/frisbee-ai/openclaw-sdk-go/pkg/...`)

**Example:**
```go
import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/frisbee-ai/openclaw-sdk-go/pkg/api"
	"github.com/frisbee-ai/openclaw-sdk-go/pkg/auth"
	"github.com/gorilla/websocket"
)
```

## Error Handling

**Pattern:** Custom error types with error codes

Error hierarchy in `pkg/types/errors.go`:
- `BaseError` - core error with code, message, retryable flag
- `OpenClawError` interface - `Code()`, `Message()`, `Retryable()`, `Unwrap()`, `Details()`
- Specific error types: `AuthError`, `ConnectionError`, `ProtocolError`, `RequestError`, `GatewayError`, `ReconnectError`, `TimeoutError`, `CancelledError`, `AbortError`
- `NewAPIError(*ErrorShape)` - factory that routes error codes to correct type

**Error Creation:**
```go
// From pkg/types/errors.go
NewConnectionError("CONNECTION_CLOSED", "Connection closed", false, nil)
NewAuthError("AUTH_TOKEN_EXPIRED", "Token expired", true, nil)
NewRequestError("METHOD_NOT_FOUND", "Method not found", false, nil)
```

**Error Checking:**
```go
// Type assertions
if authErr, ok := err.(*AuthError); ok {
    // handle auth error
}

// Interface checks
if IsRetryable(err) { /* retry */ }
if IsGatewayError(err) { /* handle */ }

// Standard errors.Is/As
if errors.Is(err, context.Canceled) { /* */ }
```

**Never silently ignore errors:**
```go
// WRONG
_ = someFunction()

// CORRECT
if err := someFunction(); err != nil {
    return err
}
```

## Logging

**Logger interface:** `types.Logger`
- `Debug(msg string, keyvals ...any)`
- `Info(msg string, keyvals ...any)`
- `Warn(msg string, keyval ...any)`
- `Error(msg string, keyvals ...any)`

**Implementations:**
- `DefaultLogger` - structured logging to stdout
- `NopLogger` - no-op implementation
- `WithContext`/`FromContext` for context-scoped loggers

**Usage:**
```go
c.config.Logger.Error("protocol negotiation failed", "error", err)
```

## Comments

**When to Comment:**
- Exported functions and types: package-level doc comments
- Non-obvious logic: inline comments explaining WHY
- Bug workarounds: comment explaining the issue

**Style:** Standard Go doc comments (`// FunctionName does...`)

**Example from `pkg/client.go`:**
```go
// ClientConfig holds client configuration for creating an OpenClaw client.
// It contains all the settings needed to connect to a WebSocket server.
type ClientConfig struct { ... }

// Connect establishes a WebSocket connection to the server.
// Thread-safe method that validates URL and initiates connection.
func (c *client) Connect(ctx context.Context) error { ... }
```

## Function Design

**Size:** Small, focused functions (<50 lines typical)

**Parameters:**
- Context as first parameter when used: `func(ctx context.Context, ...)`
- Options pattern for configuration: `func WithURL(url string) ClientOption`
- Error return: `error` as last return value

**Return Values:**
- Named returns only when clarity benefits
- Interface return types for flexibility

## Module Design

**Exports:**
- Single package re-export pattern in `pkg/client.go` for convenience
- Types re-exported from subpackages: `type ConnectionState = types.ConnectionState`

**Barrel Files:**
- `pkg/client.go` acts as main entry point re-exporting all public types

**Internal Packages:**
- `pkg/` contains all implementation packages
- API subpackages: `pkg/api/` for API method groups

## Thread-Safety Patterns

**Critical Rule:** Never send to a channel while holding a lock

**Mutex Usage:**
```go
// Lock before accessing shared state
c.mu.Lock()
defer c.mu.Unlock()

// RWMutex for read-heavy workloads
em.mu.RLock()
defer em.mu.RUnlock()
```

**State Machine:**
- `sync.RWMutex` for connection state
- Release lock BEFORE sending to channels

**Buffered Channels:**
- All event channels are buffered to prevent deadlocks
- Default buffer size: 100

## Context Usage

**Pattern:** Context for cancellation and timeouts
```go
ctx, cancel := context.WithCancel(context.Background())
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
```

**Never store context in structs** - pass as first parameter

## Graceful Shutdown

**Pattern:** All managers implement `Close()` with goroutine cleanup
```go
func (em *EventManager) Close() error {
    em.closedMu.Lock()
    defer em.closedMu.Unlock()
    if em.closed {
        return nil // Idempotent
    }
    em.closed = true
    em.cancel()
    em.wg.Wait()
    return nil
}
```

## Pre-commit Hooks

Location: `.pre-commit-config.yaml`

**Hooks:**
1. `go fmt` - `gofmt -w -l .`
2. `go vet` - `go vet ./...`
3. `golangci-lint` - `golangci-lint run --fix --timeout=5m`
4. `go test` - `go test ./...`

---

*Convention analysis: 2026-03-28*
