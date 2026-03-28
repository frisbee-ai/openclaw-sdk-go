package types

import (
	"testing"
	"time"
)

func TestConnectionState(t *testing.T) {
	states := []ConnectionState{
		StateDisconnected,
		StateConnecting,
		StateConnected,
		StateAuthenticating,
		StateAuthenticated,
		StateReconnecting,
		StateFailed,
	}

	for _, s := range states {
		if s == "" {
			t.Error("state should not be empty")
		}
	}
}

func TestEventType(t *testing.T) {
	types := []EventType{
		EventConnect,
		EventDisconnect,
		EventError,
		EventMessage,
		EventRequest,
		EventResponse,
		EventTick,
		EventGap,
		EventStateChange,
	}

	for _, et := range types {
		if et == "" {
			t.Error("event type should not be empty")
		}
	}
}

func TestConnectionMetrics_StructHasFields(t *testing.T) {
	m := ConnectionMetrics{
		Latency:        5 * time.Second,
		LastTickAge:    1 * time.Second,
		ReconnectCount: 3,
		IsStale:        true,
	}
	if m.Latency != 5*time.Second {
		t.Errorf("expected Latency=5s, got %v", m.Latency)
	}
	if m.LastTickAge != 1*time.Second {
		t.Errorf("expected LastTickAge=1s, got %v", m.LastTickAge)
	}
	if m.ReconnectCount != 3 {
		t.Errorf("expected ReconnectCount=3, got %d", m.ReconnectCount)
	}
	if !m.IsStale {
		t.Error("expected IsStale=true")
	}
}

func TestEventPriority_Ordering(t *testing.T) {
	if EventPriorityLow >= EventPriorityMedium {
		t.Error("expected EventPriorityLow < EventPriorityMedium")
	}
	if EventPriorityMedium >= EventPriorityHigh {
		t.Error("expected EventPriorityMedium < EventPriorityHigh")
	}
	if EventPriorityLow >= EventPriorityHigh {
		t.Error("expected EventPriorityLow < EventPriorityHigh")
	}
}

func TestEventPriority_Values(t *testing.T) {
	if EventPriorityLow != 0 {
		t.Errorf("expected EventPriorityLow=0, got %d", EventPriorityLow)
	}
	if EventPriorityMedium != 1 {
		t.Errorf("expected EventPriorityMedium=1, got %d", EventPriorityMedium)
	}
	if EventPriorityHigh != 2 {
		t.Errorf("expected EventPriorityHigh=2, got %d", EventPriorityHigh)
	}
}

func TestEventPriority_Compare(t *testing.T) {
	low := EventPriorityLow
	med := EventPriorityMedium
	high := EventPriorityHigh

	if !(low < med) {
		t.Error("expected low < med")
	}
	if !(med < high) {
		t.Error("expected med < high")
	}
	if low > med {
		t.Error("expected low <= med")
	}
}

func TestDefaultReconnectConfig(t *testing.T) {
	cfg := DefaultReconnectConfig()

	if cfg.MaxAttempts != 0 {
		t.Errorf("expected MaxAttempts=0 (infinite), got %d", cfg.MaxAttempts)
	}
	if cfg.MaxRetries != 10 {
		t.Errorf("expected MaxRetries=10, got %d", cfg.MaxRetries)
	}
	if cfg.InitialDelay != 1*time.Second {
		t.Errorf("expected InitialDelay=1s, got %v", cfg.InitialDelay)
	}
	if cfg.MaxDelay != 60*time.Second {
		t.Errorf("expected MaxDelay=60s, got %v", cfg.MaxDelay)
	}
	if cfg.BackoffMultiplier != 1.618 {
		t.Errorf("expected BackoffMultiplier=1.618, got %f", cfg.BackoffMultiplier)
	}
}
