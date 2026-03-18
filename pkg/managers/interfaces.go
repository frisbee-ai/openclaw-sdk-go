package managers

import (
	"context"

	"github.com/i0r3k/openclaw-sdk-go/pkg/protocol"
	"github.com/i0r3k/openclaw-sdk-go/pkg/transport"
	"github.com/i0r3k/openclaw-sdk-go/pkg/types"
)

// EventEmitter is the interface for event emission
type EventEmitter interface {
	Emit(event types.Event)
	Events() <-chan types.Event
}

// EventManagerInterface defines the interface for event management
type EventManagerInterface interface {
	Subscribe(eventType types.EventType, handler types.EventHandler) func()
	Unsubscribe(eventType types.EventType, handler types.EventHandler)
	Events() <-chan types.Event
	Emit(event types.Event)
	Start()
	Close() error
}

// RequestManagerInterface defines the interface for request management
type RequestManagerInterface interface {
	SendRequest(ctx context.Context, req *protocol.RequestFrame, sendFunc func(*protocol.RequestFrame) error) (*protocol.ResponseFrame, error)
	HandleResponse(frame *protocol.ResponseFrame)
	Close() error
}

// ConnectionManagerInterface defines the interface for connection management
type ConnectionManagerInterface interface {
	Connect(ctx context.Context) error
	Disconnect() error
	State() types.ConnectionState
	Transport() transport.Transport
	Close() error
}

// ReconnectManagerInterface defines the interface for reconnection management
type ReconnectManagerInterface interface {
	SetOnReconnect(f func() error)
	SetOnReconnectFailed(f func(err error))
	Start()
	Stop()
}
