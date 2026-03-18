# Phase 9: Main Client

**Files:**
- Create: `client.go`, `client_test.go`

**Depends on:** All previous phases

---

## Task 9.1: Client with Options

- [ ] **Step 1: Write client.go**

```go
package openclaw

import (
	"context"
	"sync"
	"time"

	"openclaw-sdk-go/managers"
	"openclaw-sdk-go/protocol"
	"openclaw-sdk-go/transport"
)

// ClientConfig holds client configuration
type ClientConfig struct {
	URL              string
	AuthHandler      AuthHandler
	ReconnectEnabled bool
	ReconnectConfig  *managers.ReconnectConfig
	Logger           Logger
	Header           map[string][]string
	TLSConfig        *transport.TLSConfig
	EventBufferSize  int
}

// DefaultClientConfig returns default configuration
func DefaultClientConfig() *ClientConfig {
	return &ClientConfig{
		EventBufferSize: 100,
		Logger:          &DefaultLogger{},
	}
}

// ClientOption is a functional option
type ClientOption func(*ClientConfig) error

// WithURL sets the WebSocket URL
func WithURL(url string) ClientOption {
	return func(c *ClientConfig) error {
		c.URL = url
		return nil
	}
}

// WithAuthHandler sets the auth handler
func WithAuthHandler(handler AuthHandler) ClientOption {
	return func(c *ClientConfig) error {
		c.AuthHandler = handler
		return nil
	}
}

// WithReconnect enables or disables reconnect
func WithReconnect(enabled bool) ClientOption {
	return func(c *ClientConfig) error {
		c.ReconnectEnabled = enabled
		return nil
	}
}

// WithLogger sets the logger
func WithLogger(logger Logger) ClientOption {
	return func(c *ClientConfig) error {
		c.Logger = logger
		return nil
	}
}

// OpenClawClient is the main client interface
type OpenClawClient interface {
	Connect(ctx context.Context) error
	Disconnect() error
	State() ConnectionState
	SendRequest(ctx context.Context, req *protocol.RequestFrame) (*protocol.ResponseFrame, error)
	Events() <-chan Event
	Subscribe(eventType EventType, handler EventHandler) func()
	Close() error
}

// client is the concrete implementation
type client struct {
	config   *ClientConfig
	state    ConnectionState
	managers struct {
		event      *managers.EventManager
		request    *managers.RequestManager
		connection *managers.ConnectionManager
		reconnect  *managers.ReconnectManager
	}
	ctx    context.Context
	cancel context.CancelFunc
	mu     sync.Mutex
}

// NewClient creates a new OpenClaw client
func NewClient(opts ...ClientOption) (OpenClawClient, error) {
	cfg := DefaultClientConfig()
	for _, opt := range opts {
		if err := opt(cfg); err != nil {
			return nil, err
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	c := &client{
		config: cfg,
		ctx:    ctx,
		cancel: cancel,
	}

	// Initialize managers
	c.managers.event = managers.NewEventManager(ctx, cfg.EventBufferSize)
	c.managers.request = managers.NewRequestManager(ctx)
	c.managers.connection = managers.NewConnectionManager(cfg, c.managers.event)

	if cfg.ReconnectEnabled {
		c.managers.reconnect = managers.NewReconnectManager(cfg.ReconnectConfig)
	}

	c.managers.event.Start()

	return c, nil
}

// Connect establishes a connection
func (c *client) Connect(ctx context.Context) error {
	return c.managers.connection.Connect(ctx)
}

// Disconnect closes the connection
func (c *client) Disconnect() error {
	return c.managers.connection.Disconnect()
}

// State returns the current connection state
func (c *client) State() ConnectionState {
	return c.managers.connection.State()
}

// SendRequest sends a request
func (c *client) SendRequest(ctx context.Context, req *protocol.RequestFrame) (*protocol.ResponseFrame, error) {
	return c.managers.request.SendRequest(ctx, req)
}

// Events returns the event channel
func (c *client) Events() <-chan Event {
	return c.managers.event.Events()
}

// Subscribe subscribes to events
func (c *client) Subscribe(eventType EventType, handler EventHandler) func() {
	return c.managers.event.Subscribe(eventType, handler)
}

// Close closes the client
func (c *client) Close() error {
	c.cancel()
	c.managers.event.Close()
	c.managers.request.Close()
	return c.managers.connection.Close()
}
```

- [ ] **Step 2: Write basic test**

```go
// client_test.go
package openclaw

import (
	"context"
	"testing"
)

func TestNewClient(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Fatal("expected client to not be nil")
	}
	defer client.Close()

	if client.State() != StateDisconnected {
		t.Errorf("expected disconnected state, got %s", client.State())
	}
}

func TestClientOptions(t *testing.T) {
	client, err := NewClient(
		WithURL("ws://localhost:8080"),
		WithReconnect(true),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer client.Close()
}
```

- [ ] **Step 3: Run tests and commit**

Run: `go build ./... && go test -v ./...`
Commit: `git add client.go client_test.go && git commit -m "feat: add main client with options"`

---

## Phase 9 Complete

After this phase, you should have:
- `client.go` - Main client with options pattern
- All managers integrated

All code should compile and tests should pass.
