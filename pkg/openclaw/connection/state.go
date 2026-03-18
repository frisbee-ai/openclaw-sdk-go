// Package connection provides connection management components
package connection

import (
	"fmt"
	"sync"

	openclaw "github.com/i0r3k/openclaw-sdk-go/pkg/openclaw"
)

// StateChangeEvent represents a state change event
type StateChangeEvent struct {
	From   openclaw.ConnectionState
	To     openclaw.ConnectionState
	Reason error
}

// ConnectionStateMachine manages connection state
type ConnectionStateMachine struct {
	state  openclaw.ConnectionState
	mu     sync.RWMutex
	events chan StateChangeEvent
	ctx    interface{} // context.Context - added for future cancellation support
}

// NewConnectionStateMachine creates a new state machine
func NewConnectionStateMachine(initial openclaw.ConnectionState) *ConnectionStateMachine {
	return &ConnectionStateMachine{
		state:  initial,
		events: make(chan StateChangeEvent, 10),
	}
}

// validTransitions defines valid state transitions using typed constants
var validTransitions = map[openclaw.ConnectionState][]openclaw.ConnectionState{
	openclaw.StateDisconnected:     {openclaw.StateConnecting},
	openclaw.StateConnecting:       {openclaw.StateConnected, openclaw.StateDisconnected, openclaw.StateFailed},
	openclaw.StateConnected:        {openclaw.StateAuthenticating, openclaw.StateDisconnected, openclaw.StateReconnecting, openclaw.StateFailed},
	openclaw.StateAuthenticating:   {openclaw.StateAuthenticated, openclaw.StateFailed},
	openclaw.StateAuthenticated:    {openclaw.StateDisconnected, openclaw.StateReconnecting},
	openclaw.StateReconnecting:    {openclaw.StateConnecting, openclaw.StateFailed},
	openclaw.StateFailed:           {openclaw.StateDisconnected},
}

func (csm *ConnectionStateMachine) validTransition(from, to openclaw.ConnectionState) bool {
	allowed, ok := validTransitions[from]
	if !ok {
		return false
	}
	for _, s := range allowed {
		if s == to {
			return true
		}
	}
	return false
}

// Transition changes the state
func (csm *ConnectionStateMachine) Transition(to openclaw.ConnectionState, reason error) error {
	csm.mu.Lock()
	from := csm.state
	if !csm.validTransition(from, to) {
		csm.mu.Unlock()
		return fmt.Errorf("invalid state transition from %s to %s", from, to)
	}
	csm.state = to
	csm.mu.Unlock()

	select {
	case csm.events <- StateChangeEvent{From: from, To: to, Reason: reason}:
	default:
		// Channel full - return error so caller knows event was dropped
		return fmt.Errorf("state change event dropped: %s -> %s", from, to)
	}
	return nil
}

// State returns the current state
func (csm *ConnectionStateMachine) State() openclaw.ConnectionState {
	csm.mu.RLock()
	defer csm.mu.RUnlock()
	return csm.state
}

// Events returns the state change event channel
func (csm *ConnectionStateMachine) Events() <-chan StateChangeEvent {
	return csm.events
}
