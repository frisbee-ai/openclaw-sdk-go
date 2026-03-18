// pkg/openclaw/connection/protocol_test.go
package connection

import (
	"context"
	"testing"
	"time"
)

func TestProtocolNegotiator_Negotiate_Match(t *testing.T) {
	negotiator := NewProtocolNegotiator([]string{"1.0", "2.0"})

	version, err := negotiator.Negotiate(context.Background(), []string{"1.0", "1.1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if version != "1.0" {
		t.Errorf("expected '1.0', got '%s'", version)
	}
}

func TestProtocolNegotiator_Negotiate_NoMatch(t *testing.T) {
	negotiator := NewProtocolNegotiator([]string{"1.0", "2.0"})

	_, err := negotiator.Negotiate(context.Background(), []string{"3.0", "4.0"})
	if err == nil {
		t.Error("expected error for no matching version")
	}
}

func TestProtocolNegotiator_Negotiate_ContextCancel(t *testing.T) {
	negotiator := NewProtocolNegotiator([]string{"1.0"})

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err := negotiator.Negotiate(ctx, []string{"1.0"})
	if err == nil {
		t.Error("expected error for cancelled context")
	}
}

func TestProtocolNegotiator_DefaultVersions(t *testing.T) {
	negotiator := NewProtocolNegotiator(nil)

	version, err := negotiator.Negotiate(context.Background(), []string{"1.0"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if version != "1.0" {
		t.Errorf("expected '1.0', got '%s'", version)
	}
}

func TestProtocolNegotiator_Negotiate_Timeout(t *testing.T) {
	negotiator := NewProtocolNegotiator([]string{"1.0"})

	// Create a context that times out immediately
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()

	// Wait for timeout
	time.Sleep(10 * time.Millisecond)

	_, err := negotiator.Negotiate(ctx, []string{"1.0"})
	if err == nil {
		t.Error("expected error for timeout")
	}
}
