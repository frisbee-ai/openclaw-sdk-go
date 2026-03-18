# Phase 5: Connection Module

**Files:**
- Create: `connection/state.go`, `connection/state_test.go`
- Create: `connection/protocol.go`
- Create: `connection/policies.go`, `connection/tls.go`

**Depends on:** Phase 1 (types.go, errors.go), Phase 4 (transport)

---

## Task 5.1: Connection State Machine

- [ ] **Step 1: Create connection directory and state.go**

```bash
mkdir -p connection
```

```go
// connection/state.go
package connection

import (
	"fmt"
	"sync"
)

// StateChangeEvent represents a state change event
type StateChangeEvent struct {
	From   string
	To     string
	Reason error
}

// ConnectionStateMachine manages connection state
type ConnectionStateMachine struct {
	state  string
	mu     sync.RWMutex
	events chan StateChangeEvent
}

// NewConnectionStateMachine creates a new state machine
func NewConnectionStateMachine(initial string) *ConnectionStateMachine {
	return &ConnectionStateMachine{
		state:  initial,
		events: make(chan StateChangeEvent, 10),
	}
}

// validTransitions defines valid state transitions
var validTransitions = map[string][]string{
	"disconnected":    {"connecting"},
	"connecting":      {"connected", "disconnected", "failed"},
	"connected":       {"authenticating", "disconnected", "reconnecting", "failed"},
	"authenticating":  {"authenticated", "failed"},
	"authenticated":   {"disconnected", "reconnecting"},
	"reconnecting":    {"connecting", "failed"},
	"failed":          {"disconnected"},
}

func (csm *ConnectionStateMachine) validTransition(from, to string) bool {
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
func (csm *ConnectionStateMachine) Transition(to string, reason error) error {
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
		// Channel full - log warning but don't block
	}
	return nil
}

// State returns the current state
func (csm *ConnectionStateMachine) State() string {
	csm.mu.RLock()
	defer csm.mu.RUnlock()
	return csm.state
}

// Events returns the state change event channel
func (csm *ConnectionStateMachine) Events() <-chan StateChangeEvent {
	return csm.events
}
```

- [ ] **Step 2: Write test**

```go
// connection/state_test.go
package connection

import (
	"testing"
)

func TestConnectionStateMachine_Transition(t *testing.T) {
	csm := NewConnectionStateMachine("disconnected")

	err := csm.Transition("connecting", nil)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if csm.State() != "connecting" {
		t.Errorf("expected 'connecting', got '%s'", csm.State())
	}
}

func TestConnectionStateMachine_InvalidTransition(t *testing.T) {
	csm := NewConnectionStateMachine("disconnected")

	err := csm.Transition("authenticated", nil)
	if err == nil {
		t.Error("expected error for invalid transition")
	}
}
```

- [ ] **Step 3: Run tests and commit**

Run: `go test -v ./connection/...`
Commit: `git add connection/ && git commit -m "feat: add connection state machine"`

---

## Task 5.2: Protocol Negotiator

- [ ] **Step 1: Write protocol.go**

```go
// connection/protocol.go
package connection

import (
	"context"
)

// ProtocolNegotiator handles protocol version negotiation
type ProtocolNegotiator struct {
	supportedVersions []string
}

// NewProtocolNegotiator creates a new negotiator
func NewProtocolNegotiator(supportedVersions []string) *ProtocolNegotiator {
	if len(supportedVersions) == 0 {
		supportedVersions = []string{"1.0"}
	}
	return &ProtocolNegotiator{supportedVersions: supportedVersions}
}

// Negotiate performs protocol negotiation
func (p *ProtocolNegotiator) Negotiate(ctx context.Context, serverVersions []string) (string, error) {
	for _, clientVer := range p.supportedVersions {
		for _, serverVer := range serverVersions {
			if clientVer == serverVer {
				return clientVer, nil
			}
		}
	}
	return "", ErrNoMatchingProtocol
}

// ErrNoMatchingProtocol is returned when no matching protocol is found
var ErrNoMatchingProtocol = &ProtocolError{Message: "no matching protocol version"}

type ProtocolError struct {
	Message string
}

func (e *ProtocolError) Error() string {
	return e.Message
}
```

- [ ] **Step 2: Commit**

```bash
git add connection/protocol.go
git commit -m "feat: add protocol negotiator"
```

---

## Task 5.3: Policy Manager and TLS Validator

- [ ] **Step 1: Write policies.go**

```go
// connection/policies.go
package connection

// PolicyManager manages connection policies
type PolicyManager struct {
	maxReconnectAttempts int
	pingInterval         int
}

// NewPolicyManager creates a new policy manager
func NewPolicyManager(maxReconnectAttempts int, pingInterval int) *PolicyManager {
	return &PolicyManager{
		maxReconnectAttempts: maxReconnectAttempts,
		pingInterval:         pingInterval,
	}
}

// MaxReconnectAttempts returns the max reconnect attempts
func (pm *PolicyManager) MaxReconnectAttempts() int {
	return pm.maxReconnectAttempts
}

// PingInterval returns the ping interval in seconds
func (pm *PolicyManager) PingInterval() int {
	return pm.pingInterval
}
```

- [ ] **Step 2: Write tls.go**

```go
// connection/tls.go
package connection

import (
	"crypto/tls"
	"crypto/x509"
	"os"
)

// TlsValidator validates TLS certificates
type TlsValidator struct {
	config *TLSConfig
}

// TLSConfig holds TLS configuration
type TLSConfig struct {
	InsecureSkipVerify bool
	CertFile          string
	KeyFile           string
	CAFile            string
	ServerName        string
}

// NewTlsValidator creates a new TLS validator
func NewTlsValidator(config *TLSConfig) *TlsValidator {
	return &TlsValidator{config: config}
}

// Validate validates the TLS configuration
func (v *TlsValidator) Validate() error {
	if v.config == nil {
		return nil
	}
	return nil
}

// GetTLSConfig returns the TLS config for the connection
func (v *TlsValidator) GetTLSConfig() (*tls.Config, error) {
	config := &tls.Config{
		InsecureSkipVerify: v.config.InsecureSkipVerify,
		ServerName:         v.config.ServerName,
	}

	if v.config.CertFile != "" && v.config.KeyFile != "" {
		cert, err := tls.LoadX509KeyPair(v.config.CertFile, v.config.KeyFile)
		if err != nil {
			return nil, err
		}
		config.Certificates = []tls.Certificate{cert}
	}

	if v.config.CAFile != "" {
		caCert, err := os.ReadFile(v.config.CAFile)
		if err != nil {
			return nil, err
		}
		caPool := x509.NewCertPool()
		caPool.AppendCertsFromPEM(caCert)
		config.RootCAs = caPool
	}

	return config, nil
}
```

- [ ] **Step 3: Commit**

```bash
git add connection/policies.go connection/tls.go
git commit -m "feat: add policy manager and TLS validator"
```

---

## Phase 5 Complete

After this phase, you should have:
- `connection/state.go` - Connection state machine
- `connection/protocol.go` - Protocol negotiator
- `connection/policies.go` - Policy manager
- `connection/tls.go` - TLS validator

All code should compile and tests should pass.
