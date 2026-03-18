// pkg/openclaw/connection/policies_test.go
package connection

import (
	"testing"
	"time"
)

func TestPolicyManager_InfiniteReconnect(t *testing.T) {
	pm := NewPolicyManager(0, 30*time.Second)

	if !pm.ShouldReconnect(0) {
		t.Error("expected ShouldReconnect(0) to return true for infinite retries")
	}
	if !pm.ShouldReconnect(100) {
		t.Error("expected ShouldReconnect(100) to return true for infinite retries")
	}
}

func TestPolicyManager_LimitedReconnect(t *testing.T) {
	pm := NewPolicyManager(3, 30*time.Second)

	if !pm.ShouldReconnect(0) {
		t.Error("expected ShouldReconnect(0) to return true")
	}
	if !pm.ShouldReconnect(2) {
		t.Error("expected ShouldReconnect(2) to return true")
	}
	if pm.ShouldReconnect(3) {
		t.Error("expected ShouldReconnect(3) to return false (attempt 3 == max)")
	}
	if pm.ShouldReconnect(10) {
		t.Error("expected ShouldReconnect(10) to return false")
	}
}

func TestPolicyManager_PingInterval(t *testing.T) {
	pm := NewPolicyManager(0, 30*time.Second)

	interval := pm.PingInterval()
	if interval != 30*time.Second {
		t.Errorf("expected 30s, got %v", interval)
	}
}
