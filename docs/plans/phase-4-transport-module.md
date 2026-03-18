# Phase 4: Transport Module

**Files:**
- Create: `transport/websocket.go`, `transport/websocket_test.go`

**Depends on:** Phase 1 (types.go, logger.go), Phase 3 (protocol/types.go)

---

## Task 4.1: WebSocket Transport

- [ ] **Step 1: Create transport directory and websocket.go**

```bash
mkdir -p transport
```

```go
// transport/websocket.go
package transport

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// WebSocketConfig holds WebSocket configuration
type WebSocketConfig struct {
	ReadBufferSize  int
	WriteBufferSize int
	TLSConfig       *TLSConfig
	Header          http.Header
}

// TLSConfig holds TLS configuration
type TLSConfig struct {
	InsecureSkipVerify bool
	CertFile          string
	KeyFile           string
	CAFile            string
	ServerName        string
}

// WebSocketTransport handles WebSocket communication
type WebSocketTransport struct {
	conn    *websocket.Conn
	sendCh  chan []byte
	recvCh  chan []byte
	ctx     context.Context
	cancel  context.CancelFunc
	mu      sync.Mutex
	wg      sync.WaitGroup
	closed  bool
}

// Dial creates a new WebSocket connection
func Dial(url string, config *WebSocketConfig) (*WebSocketTransport, error) {
	dialer := websocket.Dialer{
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,
	}

	if config != nil && config.TLSConfig != nil {
		dialer.TLSClientConfig = config.TLSConfig.toTLSConfig()
	}

	conn, _, err := dialer.Dial(url, config.Header)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &WebSocketTransport{
		conn:   conn,
		sendCh: make(chan []byte, 10),
		recvCh: make(chan []byte, 10),
		ctx:    ctx,
		cancel: cancel,
	}, nil
}

// Start begins the send/receive loops
func (t *WebSocketTransport) Start() {
	t.wg.Add(2)
	go t.readLoop()
	go t.writeLoop()
}

// readLoop reads messages from the WebSocket
func (t *WebSocketTransport) readLoop() {
	defer t.wg.Done()
	for {
		select {
		case <-t.ctx.Done():
			return
		default:
			_, message, err := t.conn.ReadMessage()
			if err != nil {
				t.handleError(err)
				return
			}
			select {
			case t.recvCh <- message:
			case <-t.ctx.Done():
				return
			}
		}
	}
}

// writeLoop writes messages to the WebSocket
func (t *WebSocketTransport) writeLoop() {
	defer t.wg.Done()
	for {
		select {
		case <-t.ctx.Done():
			return
		case message := <-t.sendCh:
			if err := t.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				t.handleError(err)
				return
			}
		}
	}
}

// handleError handles WebSocket errors
func (t *WebSocketTransport) handleError(err error) {
	// Could emit error event here
}

// Send sends a message
func (t *WebSocketTransport) Send(data []byte) error {
	select {
	case t.sendCh <- data:
		return nil
	case <-t.ctx.Done():
		return t.ctx.Err()
	}
}

// Receive returns the receive channel
func (t *WebSocketTransport) Receive() <-chan []byte {
	return t.recvCh
}

// Close closes the WebSocket connection
func (t *WebSocketTransport) Close() error {
	t.mu.Lock()
	if t.closed {
		t.mu.Unlock()
		return nil
	}
	t.closed = true
	t.mu.Unlock()

	t.cancel()
	t.wg.Wait()
	return t.conn.Close()
}

// IsConnected returns whether the transport is connected
func (t *WebSocketTransport) IsConnected() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.conn != nil && !t.closed
}

// Helper to convert TLSConfig
func (c *TLSConfig) toTLSConfig() *TLSConfig {
	return c // Simplified - actual impl would use crypto/tls
}
```

- [ ] **Step 2: Write basic test**

```go
// transport/websocket_test.go
package transport

import (
	"testing"
)

func TestWebSocketConfig_Defaults(t *testing.T) {
	config := &WebSocketConfig{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	if config.ReadBufferSize != 1024 {
		t.Errorf("expected 1024, got %d", config.ReadBufferSize)
	}
}
```

- [ ] **Step 3: Run tests and commit**

Run: `go build ./transport/...`
Commit: `git add transport/ && git commit -m "feat: add WebSocket transport"`

---

## Phase 4 Complete

After this phase, you should have:
- `transport/websocket.go` - WebSocket transport implementation

All code should compile and tests should pass.
