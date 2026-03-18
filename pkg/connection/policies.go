// pkg/openclaw/connection/policies.go
package connection

import (
	"time"
)

// PolicyManager manages connection policies
type PolicyManager struct {
	maxReconnectAttempts int
	pingInterval         time.Duration
}

// NewPolicyManager creates a new policy manager
func NewPolicyManager(maxReconnectAttempts int, pingInterval time.Duration) *PolicyManager {
	return &PolicyManager{
		maxReconnectAttempts: maxReconnectAttempts,
		pingInterval:         pingInterval,
	}
}

// MaxReconnectAttempts returns the max reconnect attempts
// Returns 0 for infinite retries
func (pm *PolicyManager) MaxReconnectAttempts() int {
	return pm.maxReconnectAttempts
}

// PingInterval returns the ping interval
func (pm *PolicyManager) PingInterval() time.Duration {
	return pm.pingInterval
}

// ShouldReconnect checks if reconnection should be attempted based on attempt count
func (pm *PolicyManager) ShouldReconnect(attemptCount int) bool {
	if pm.maxReconnectAttempts == 0 {
		return true // Infinite retries
	}
	return attemptCount < pm.maxReconnectAttempts
}
