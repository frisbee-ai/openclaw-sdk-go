// Package connection provides connection management components for OpenClaw SDK.
//
// This package provides:
//   - ConnectionStateMachine: State machine for managing connection lifecycle
//   - PolicyManager: Connection policy configuration
//   - ProtocolNegotiator: Protocol version negotiation
//   - TLS validation: Certificate and configuration validation
package connection

import (
	"time"
)

// PolicyManager manages connection policies.
// It defines rules for reconnection behavior and ping intervals.
type PolicyManager struct {
	maxReconnectAttempts int           // Maximum number of reconnection attempts (0 = infinite)
	pingInterval         time.Duration // Interval between ping messages
}

// NewPolicyManager creates a new policy manager with the specified settings.
func NewPolicyManager(maxReconnectAttempts int, pingInterval time.Duration) *PolicyManager {
	return &PolicyManager{
		maxReconnectAttempts: maxReconnectAttempts,
		pingInterval:         pingInterval,
	}
}

// MaxReconnectAttempts returns the maximum number of reconnect attempts.
// Returns 0 for infinite retries.
func (pm *PolicyManager) MaxReconnectAttempts() int {
	return pm.maxReconnectAttempts
}

// PingInterval returns the configured ping interval.
func (pm *PolicyManager) PingInterval() time.Duration {
	return pm.pingInterval
}

// ShouldReconnect checks if reconnection should be attempted based on the attempt count.
// Returns true if more attempts should be made, false otherwise.
func (pm *PolicyManager) ShouldReconnect(attemptCount int) bool {
	if pm.maxReconnectAttempts == 0 {
		return true // Infinite retries
	}
	return attemptCount < pm.maxReconnectAttempts
}
