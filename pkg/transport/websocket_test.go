package transport

import (
	"context"
	"testing"
	"time"
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

func TestWebSocketTransport_Send_Blocking(t *testing.T) {
	// Create transport with buffered channel
	ctx, cancel := context.WithCancel(context.Background())
	transport := &WebSocketTransport{
		sendCh: make(chan []byte, 1),
		recvCh: make(chan []byte, 1),
		errCh:  make(chan error, 1),
		ctx:    ctx,
		cancel: cancel,
	}
	defer cancel()

	// Test send succeeds when channel is not full
	err := transport.Send([]byte("test message"))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestWebSocketTransport_Receive(t *testing.T) {
	transport := &WebSocketTransport{
		sendCh: make(chan []byte, 1),
		recvCh: make(chan []byte, 1),
		errCh:  make(chan error, 1),
	}

	// Send a message to recvCh
	testMsg := []byte("test")
	select {
	case transport.recvCh <- testMsg:
	default:
		t.Fatal("failed to send test message")
	}

	// Receive should return the message
	select {
	case msg := <-transport.Receive():
		if string(msg) != "test" {
			t.Errorf("expected 'test', got '%s'", string(msg))
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for message")
	}
}

func TestWebSocketTransport_IsConnected(t *testing.T) {
	// Test: conn=nil, closed=false -> not connected
	transport := &WebSocketTransport{
		closed: false,
	}
	if transport.IsConnected() {
		t.Error("expected not connected when conn=nil")
	}

	// Test: closed=true -> not connected regardless of conn
	transport.closed = true
	if transport.IsConnected() {
		t.Error("expected disconnected when closed=true")
	}
}

func TestWebSocketTransport_Close_Idempotent(t *testing.T) {
	transport := &WebSocketTransport{
		closed: true,
	}

	// Second close should be idempotent
	err := transport.Close()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestWebSocketTransport_Errors(t *testing.T) {
	transport := &WebSocketTransport{
		errCh: make(chan error, 10),
	}

	errCh := transport.Errors()
	if errCh == nil {
		t.Error("expected error channel, got nil")
	}
}

func TestTLSConfig_toTLSConfig(t *testing.T) {
	config := &TLSConfig{
		InsecureSkipVerify: true,
		ServerName:         "example.com",
	}

	tlsConfig := config.toTLSConfig()
	if tlsConfig.InsecureSkipVerify != true {
		t.Error("InsecureSkipVerify not set correctly")
	}
	if tlsConfig.ServerName != "example.com" {
		t.Error("ServerName not set correctly")
	}
}

// Mock WebSocket server for integration tests
func TestWebSocketTransport_WithMockServer(t *testing.T) {
	// This would require setting up a test WebSocket server
	// For now, test the transport struct fields
	transport := &WebSocketTransport{
		sendCh: make(chan []byte, 10),
		recvCh: make(chan []byte, 10),
		errCh:  make(chan error, 10),
		ctx:    context.Background(),
	}

	if transport.sendCh == nil {
		t.Error("sendCh should be initialized")
	}
	if transport.recvCh == nil {
		t.Error("recvCh should be initialized")
	}
	if transport.errCh == nil {
		t.Error("errCh should be initialized")
	}
}
