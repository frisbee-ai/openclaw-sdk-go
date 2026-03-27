# OpenClaw SDK Go

[![OpenClaw SDK](https://img.shields.io/badge/OpenClaw-SDK-orange?logo=github)](https://openclaw.ai)
[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Codecov](https://codecov.io/gh/frisbee-ai/openclaw-sdk-go/branch/main/graph/badge.svg)](https://codecov.io/gh/frisbee-ai/openclaw-sdk-go)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

> Feature-complete WebSocket SDK for Go with automatic reconnection, event handling, and request/response correlation.

OpenClaw SDK Go is a Go implementation migrated from the TypeScript version, providing a fully-featured WebSocket client with connection management, event handling, request/response patterns, and automatic reconnection.

## Features

- **Connection Management** - Automatic connection state management with disconnect handling
- **Event System** - Publish/subscribe pattern for event handling with backpressure timeout
- **Request/Response** - Automatic request-response correlation with timeout support
- **Auto-Reconnect** - Intelligent reconnection with Fibonacci backoff
- **TLS Support** - Configurable TLS connection options with certificate validation
- **Thread-Safe** - All public APIs are concurrency-safe with proper lock management
- **Context Support** - Full `context.Context` integration for cancellation and timeouts
- **Extensible Logging** - Built-in logger interface with custom implementation support
- **Protocol Validation** - Inbound payload size validation against server policy
- **Tick Monitoring** - Heartbeat monitoring with stale event detection
- **Gap Detection** - Event sequence gap detection for reliable message delivery
- **Multiple API Namespaces** - Complete coverage of OpenClaw protocol APIs across 16 modules: Agents, Browser, Channels, Chat, Config, Cron, DevicePairing, ExecApprovals, Nodes, Push, Secrets, Sessions, Skills, System, and Usage

## Installation

```bash
go get github.com/frisbee-ai/openclaw-sdk-go
```

## Quick Start

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "log"

    openclaw "github.com/frisbee-ai/openclaw-sdk-go/pkg"
    "github.com/frisbee-ai/openclaw-sdk-go/pkg/protocol"
    "github.com/frisbee-ai/openclaw-sdk-go/pkg/types"
)

func main() {
    // Create client
    client, err := openclaw.NewClient(
        openclaw.WithURL("ws://localhost:8080/ws"),
        openclaw.WithReconnect(true),
    )
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // Subscribe to events
    client.Subscribe(types.EventConnect, func(e types.Event) {
        fmt.Println("Connected!")
    })

    // Connect to server
    ctx := context.Background()
    if err := client.Connect(ctx); err != nil {
        log.Fatal(err)
    }

    // Send request
    resp, err := client.SendRequest(ctx, &protocol.RequestFrame{
        RequestID: "req-001",
        Action:    "ping",
        Payload:   map[string]interface{}{"message": "hello"},
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Response: %v\n", resp)
}
```

### Enable Auto-Reconnect

```go
client, err := openclaw.NewClient(
    openclaw.WithURL("wss://api.example.com/ws"),
    openclaw.WithReconnect(true),
    openclaw.WithReconnectConfig(&types.ReconnectConfig{
        MaxAttempts: 10,                // Maximum reconnection attempts
        InitialDelay: 2 * time.Second,  // Initial delay
        MaxDelay:     60 * time.Second, // Maximum delay
    }),
)
```

### TLS Configuration

```go
client, err := openclaw.NewClient(
    openclaw.WithURL("wss://secure.example.com/ws"),
    openclaw.WithTLSConfig(&transport.TLSConfig{
        InsecureSkipVerify: false,
        CertFile:          "/path/to/client.crt",
        KeyFile:           "/path/to/client.key",
        CAFile:            "/path/to/ca.crt",
    }),
)
```

### Custom Logger

```go
type MyLogger struct{}

func (l *MyLogger) Debug(msg string, args ...any) {
    log.Printf("[DEBUG] %s %v", msg, args)
}

func (l *MyLogger) Info(msg string, args ...any) {
    log.Printf("[INFO] %s %v", msg, args)
}

func (l *MyLogger) Warn(msg string, args ...any) {
    log.Printf("[WARN] %s %v", msg, args)
}

func (l *MyLogger) Error(msg string, args ...any) {
    log.Printf("[ERROR] %s %v", msg, args)
}

client, err := openclaw.NewClient(
    openclaw.WithURL("ws://localhost:8080/ws"),
    openclaw.WithLogger(&MyLogger{}),
)
```

### Authentication

The SDK supports multiple authentication methods:

#### 1. Using AuthHandler (Dynamic Authentication)

```go
import "github.com/frisbee-ai/openclaw-sdk-go/pkg/auth"

handler, _ := auth.NewStaticAuthHandler(map[string]string{
    "token": "your-auth-token",
})

client, err := openclaw.NewClient(
    openclaw.WithURL("wss://api.example.com/ws"),
    openclaw.WithClientID("my-client"),
    openclaw.WithAuthHandler(handler),
)
```

#### 2. Using CredentialsProvider (Advanced)

For dynamic credential refresh or custom credential sources:

```go
type MyCredentialsProvider struct{}

func (p *MyCredentialsProvider) GetCredentials() (map[string]string, error) {
    // Fetch credentials from your source (database, API, vault, etc.)
    return map[string]string{
        "token": fmt.Sprintf("Bearer %s", getAccessToken()),
    }, nil
}

provider := &MyCredentialsProvider{}

// Create a custom AuthHandler that uses your provider
handler := &DynamicAuthHandler{provider: provider}

client, err := openclaw.NewClient(
    openclaw.WithURL("wss://api.example.com/ws"),
    openclaw.WithClientID("my-client"),
    openclaw.WithAuthHandler(handler),
)
```

### Device Pairing

For device-to-device authentication:

```go
import "github.com/frisbee-ai/openclaw-sdk-go/pkg/connection"

device := &connection.ConnectParamsDevice{
    ID:        "device-123",
    PublicKey: "base64-encoded-public-key",
    Signature: "device-signature",
    SignedAt:  time.Now().Unix(),
    Nonce:     "random-nonce",
}

// Note: Device pairing is configured via ClientConfig
// See Advanced Configuration section for details
```

### Tick Monitor (Heartbeat Monitoring)

Configure automatic heartbeat monitoring to detect stale connections:

```go
client, err := openclaw.NewClient(
    openclaw.WithURL("wss://api.example.com/ws"),
    openclaw.WithClientID("my-client"),
    openclaw.WithTickMonitor(&openclaw.TickMonitorConfig{
        TickIntervalMs:  30000,           // Check every 30 seconds
        StaleMultiplier: 2,               // Connection stale after 2 missed ticks
        OnStale: func() {
            log.Println("⚠️  Connection is stale!")
        },
        OnRecovered: func() {
            log.Println("✅ Connection recovered!")
        },
    }),
)

// After connecting, you can check tick monitor status
tickMonitor := client.GetTickMonitor()
if tickMonitor != nil && tickMonitor.IsStale() {
    log.Printf("Stale for %d ms", tickMonitor.GetStaleDuration())
}
```

### Gap Detector (Message Sequence Tracking)

Configure gap detection for reliable message delivery:

```go
client, err := openclaw.NewClient(
    openclaw.WithURL("wss://api.example.com/ws"),
    openclaw.WithClientID("my-client"),
    openclaw.WithGapDetector(&openclaw.GapDetectorConfig{
        RecoveryMode:     "reconnect",    // or "snapshot", "skip"
        SnapshotEndpoint: "/api/snapshot", // For snapshot recovery mode
        MaxGaps:          100,
        OnGap: func(gaps []openclaw.GapInfo) {
            log.Printf("❌ Detected %d message gaps", len(gaps))
            for _, gap := range gaps {
                log.Printf("   Expected %d, got %d", gap.Expected, gap.Received)
            }
        },
    }),
)

// Check gap detector status
gapDetector := client.GetGapDetector()
if gapDetector != nil && gapDetector.HasGap() {
    log.Printf("Detected %d gaps", gapDetector.GapCount())
}
```

### Server Information Access

After connecting, you can access server information, snapshots, and policies:

```go
// Connect first
err := client.Connect(ctx)
if err != nil {
    log.Fatal(err)
}

// Get server information
serverInfo := client.GetServerInfo()
if serverInfo != nil {
    log.Printf("Connected to: %s (conn: %s)",
        serverInfo.Server.Version, serverInfo.Server.ConnID)
    log.Printf("Protocol: %d", serverInfo.Protocol)
}

// Get server snapshot (contains agents, nodes, health, etc.)
snapshot := client.GetSnapshot()
if snapshot != nil {
    log.Printf("Server uptime: %d ms", snapshot.UptimeMs)
    log.Printf("State version: %d", snapshot.StateVersion)
    log.Printf("Agents: %v", snapshot.Agents)
}

// Get server policy (limits and constraints)
policy := client.GetPolicy()
if policy != nil {
    log.Printf("Max payload: %d bytes", policy.MaxPayload)
    log.Printf("Tick interval: %d ms", policy.TickIntervalMs)
}
```

## API Documentation

### Client Options

| Option | Type | Description |
|--------|------|-------------|
| `WithURL(url string)` | string | WebSocket server URL |
| `WithClientID(id string)` | string | Client identifier (required for server connection) |
| `WithAuthHandler(handler)` | AuthHandler | Authentication handler |
| `WithCredentialsProvider(p CredentialsProvider)` | CredentialsProvider | Credentials provider interface |
| `WithReconnect(enabled bool)` | bool | Enable auto-reconnect |
| `WithReconnectConfig(cfg)` | *ReconnectConfig | Reconnect configuration |
| `WithLogger(logger)` | Logger | Custom logger |
| `WithHeader(header)` | map[string][]string | Custom HTTP headers |
| `WithTLSConfig(cfg)` | *TLSConfig | TLS configuration |
| `WithEventBufferSize(n)` | int | Event buffer size |
| `WithEventEmitTimeout(t)` | time.Duration | Timeout for Emit when channel is full (default 200ms) |
| `WithTickMonitor(cfg)` | *TickMonitorConfig | Heartbeat monitoring configuration |
| `WithGapDetector(cfg)` | *GapDetectorConfig | Message gap detection configuration |

### Connection States

```go
const (
    StateDisconnected   ConnectionState = "disconnected"
    StateConnecting     ConnectionState = "connecting"
    StateConnected      ConnectionState = "connected"
    StateAuthenticating ConnectionState = "authenticating"
    StateAuthenticated  ConnectionState = "authenticated"
    StateReconnecting   ConnectionState = "reconnecting"
    StateFailed         ConnectionState = "failed"
)
```

### Event Types

```go
const (
    EventConnect     EventType = "connect"
    EventDisconnect  EventType = "disconnect"
    EventError       EventType = "error"
    EventMessage     EventType = "message"
    EventRequest     EventType = "request"
    EventResponse    EventType = "response"
    EventTick        EventType = "tick"
    EventGap         EventType = "gap"
    EventStateChange EventType = "stateChange"
)
```

### API Modules

The SDK provides 16 API modules accessible via client methods:

| Module | Accessor Method | Description |
|--------|-----------------|-------------|
| Agents | `client.Agents()` | Agent management and operations |
| Browser | `client.Browser()` | Browser automation control |
| Channels | `client.Channels()` | Channel management |
| Chat | `client.Chat()` | Chat and messaging operations |
| Config | `client.Config()` | Configuration management |
| Cron | `client.Cron()` | Scheduled task management |
| DevicePairing | `client.DevicePairing()` | Device pairing operations |
| ExecApprovals | `client.ExecApprovals()` | Execution approval workflow |
| Nodes | `client.Nodes()` | Node management |
| Push | `client.Push()` | Push notification services |
| Secrets | `client.Secrets()` | Secret management |
| Sessions | `client.Sessions()` | Session management |
| Skills | `client.Skills()` | Skill operations |
| System | `client.System()` | System-level operations (includes TTS and Wizard) |
| Usage | `client.Usage()` | Usage tracking and reporting |

#### API Usage Example

```go
// After connecting, access API modules directly
agentsList, err := client.Agents().List(ctx)
if err != nil {
    log.Fatal(err)
}

// Send a chat message
chatResponse, err := client.Chat().Create(ctx, &ChatCreateParams{
    AgentID: "agent-123",
    Message: "Hello, world!",
})

// System operations include TTS and Wizard functionality
ttsResult, err := client.System().TextToSpeech(ctx, &TTSParams{
    Text: "Hello from OpenClaw Go SDK!",
    Voice: "default",
})
```

## Project Structure

```
openclaw-sdk-go/
├── pkg/
│   ├── client.go          # Main client API (Option pattern configuration)
│   ├── api/               # API namespace modules (16 modules: agents, browser, channels, chat, config, cron, device_pairing, exec_approvals, nodes, push, secrets, sessions, skills, shared, system, usage)
│   ├── types/             # Shared types (ConnectionState, Event, errors, Logger)
│   ├── auth/              # Authentication (CredentialsProvider, AuthHandler)
│   ├── transport/         # WebSocket transport layer (gorilla/websocket)
│   ├── protocol/          # Protocol frames (RequestFrame, ResponseFrame, validation)
│   ├── connection/        # Connection state machine, policies, TLS validator
│   ├── events/            # Tick monitor, Gap detector
│   ├── managers/          # High-level managers (event, request, connection, reconnect)
│   └── utils/             # Timeout manager
└── examples/
    ├── cmd/               # CLI example
    └── server/            # Echo server example
```

## Design Patterns

### Option Pattern

Client configuration uses functional options for flexible, readable construction:

```go
client, err := openclaw.NewClient(
    openclaw.WithURL("ws://..."),
    openclaw.WithTimeout(30*time.Second),
    openclaw.WithLogger(myLogger),
)
```

### Context + Channel Hybrid

- `context.Context` for lifecycle management and cancellation
- Buffered channels for event delivery (prevents deadlocks)
- **Critical Rule**: Never send to a channel while holding a lock

### Graceful Shutdown

All managers implement `Close()` with proper goroutine cleanup using `sync.WaitGroup`.

## Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for specific package
go test ./pkg/transport

# Run single test
go test -run TestWebSocketTransportDial ./pkg/transport

# Run tests with race detector
go test -race ./...
```

### Code Quality

```bash
# Format code
gofmt -w -l .

# Static analysis
go vet ./...

# Lint with golangci-lint
golangci-lint run

# Run all pre-commit hooks manually
pre-commit run --all-files
```

### Examples

Run built-in examples:

```bash
# Start echo server
go run examples/server/main.go

# Run client in another terminal
go run examples/cmd/main.go
```

## Dependencies

This project minimizes external dependencies:

- `github.com/gorilla/websocket` - Industry standard WebSocket library for Go

> Note: Go's standard library `net/http` does not include WebSocket support. The `gorilla/websocket` package is the de facto standard for Go WebSocket implementations.

## Migration Notes

This is a Go implementation migrated from `openclaw-sdk-typescript`. While the Go SDK follows Go idioms rather than the TypeScript API, it maintains **functional equivalence**:

- All features from TypeScript SDK are available in Go
- Same protocol wire format
- Same authentication flow
- Same reconnection behavior (Fibonacci backoff)
- Same event types and semantics

Users migrating from TypeScript will find equivalent functionality with Go-idiomatic APIs.

## License

Apache License 2.0 - see [LICENSE](LICENSE) for details

## Resources

- [Design Document](docs/specs/2026-03-18-typescript-to-go-migration-design.md) - Architecture decisions
- [Implementation Plans](docs/plans/) - Phased implementation plans
- [GoDoc](https://pkg.go.dev/github.com/frisbee-ai/openclaw-sdk-go) - API reference

---

Copyright © 2026 @frisbee-ai
