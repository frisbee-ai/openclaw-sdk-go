# External Integrations

**Analysis Date:** 2026-03-28

## APIs & External Services

**WebSocket Gateway:**
- OpenClaw Gateway Server - Primary service this SDK connects to
  - Protocol: WebSocket (wss:// for TLS)
  - SDK handles: connection lifecycle, authentication, request/response correlation
  - Users provide gateway URL via `WithURL()` option
  - Example URL: `wss://gateway.example.com/ws`

**No third-party API integrations** - This is a WebSocket client SDK, not a server

## Data Storage

**None** - SDK is a client library with no persistent storage

## Authentication & Identity

**Custom Auth Handlers:**
- `auth.CredentialsProvider` interface - Pluggable credential source
  - `pkg/auth/provider.go` - `StaticCredentialsProvider` implementation
  - Users implement interface for custom credential retrieval
- `auth.AuthHandler` interface - Custom authentication logic
  - `pkg/auth/handler.go` - Handler implementation
  - Deprecated option in favor of CredentialsProvider

**Auth Parameters:**
- Passed via `ClientConfig.Auth` (`*connection.ConnectParamsAuth`)
- Or via `ClientConfig.CredentialsProvider`
- Connection params include: auth credentials, device pairing credentials

## Network Security

**TLS/SSL Support:**
- `transport.TLSConfig` (`pkg/transport/websocket.go`)
  - `InsecureSkipVerify` - Skip certificate verification (insecure, testing only)
  - `CertFile` / `KeyFile` - Client certificate authentication
  - `CAFile` - Custom CA certificate
  - `ServerName` - SNI server name
- `connection.TLSConfig` (`pkg/connection/tls.go`)
  - `NewTlsValidator()` - Loads and validates TLS certificates
  - `GetTLSConfig()` - Returns `crypto/tls.Config`

## Monitoring & Observability

**Logging:**
- Custom `Logger` interface (`pkg/types/logger.go`)
- Built-in implementations:
  - `DefaultLogger` - Logs to stderr with timestamps
  - `NopLogger` - No-op logger
  - `WithContext` / `FromContext` - Context-aware logging
- Set via `WithLogger()` option

**No external observability integrations:**
- No error tracking (Sentry, etc.)
- No metrics (Prometheus, etc.)
- No distributed tracing

## CI/CD & Deployment

**GitHub Actions CI:**
- Runs on: push to main, pull requests
- Coverage upload: Codecov

**Release:**
- GoReleaser - Publishes to Go module proxy
- Mode: replace existing artifacts
- Triggers on git tags

## Protocol Layer

**Request/Response:**
- JSON-serialized frames over WebSocket
- `RequestFrame` (`pkg/protocol/types.go`) - Outbound requests
- `ResponseFrame` (`pkg/protocol/types.go`) - Inbound responses
- `GatewayFrame` - Server-initiated messages

**No external message brokers or queues**

## Environment Configuration

**No required environment variables** - SDK is configured programmatically

**User provides configuration via:**
```go
client, err := openclaw.NewClient(
    openclaw.WithURL("wss://gateway.example.com/ws"),
    openclaw.WithClientID("my-client"),
    openclaw.WithAuthHandler(handler),
)
```

## Webhooks & Callbacks

**SDK-Side Callbacks (user provides):**
- `EventHandler` - Subscribe to event types
- `TickMonitorConfig.OnStale` / `OnRecovered` - Connection health callbacks
- `GapDetectorConfig.OnGap` - Gap detection callback (interface exists, callback not wired)

**No inbound webhooks** - This is a client SDK

**No outbound webhooks** - Server pushes events via WebSocket

---

*Integration audit: 2026-03-28*
