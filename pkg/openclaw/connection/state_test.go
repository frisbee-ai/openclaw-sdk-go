package connection

import (
	"testing"
	"time"

	openclaw "github.com/i0r3k/openclaw-sdk-go/pkg/openclaw"
)

func TestConnectionStateMachine_Transition(t *testing.T) {
	csm := NewConnectionStateMachine(openclaw.StateDisconnected)

	err := csm.Transition(openclaw.StateConnecting, nil)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if csm.State() != openclaw.StateConnecting {
		t.Errorf("expected 'connecting', got '%s'", csm.State())
	}
}

func TestConnectionStateMachine_InvalidTransition(t *testing.T) {
	csm := NewConnectionStateMachine(openclaw.StateDisconnected)

	err := csm.Transition(openclaw.StateAuthenticated, nil)
	if err == nil {
		t.Error("expected error for invalid transition")
	}
}

func TestConnectionStateMachine_StateChangeEvent(t *testing.T) {
	csm := NewConnectionStateMachine(openclaw.StateDisconnected)

	err := csm.Transition(openclaw.StateConnecting, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	select {
	case event := <-csm.Events():
		if event.From != openclaw.StateDisconnected {
			t.Errorf("expected from 'disconnected', got '%s'", event.From)
		}
		if event.To != openclaw.StateConnecting {
			t.Errorf("expected to 'connecting', got '%s'", event.To)
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for state change event")
	}
}
