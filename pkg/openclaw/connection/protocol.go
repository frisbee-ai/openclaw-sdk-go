// pkg/openclaw/connection/protocol.go
package connection

import (
	"context"
	"errors"
	"time"

	openclaw "github.com/i0r3k/openclaw-sdk-go/pkg/openclaw"
)

// ProtocolNegotiator handles protocol version negotiation
type ProtocolNegotiator struct {
	supportedVersions []string
	defaultTimeout    time.Duration
}

// NewProtocolNegotiator creates a new negotiator
func NewProtocolNegotiator(supportedVersions []string) *ProtocolNegotiator {
	if len(supportedVersions) == 0 {
		supportedVersions = []string{"1.0"}
	}
	return &ProtocolNegotiator{
		supportedVersions: supportedVersions,
		defaultTimeout:    5 * time.Second,
	}
}

// Negotiate performs protocol negotiation with context support
func (p *ProtocolNegotiator) Negotiate(ctx context.Context, serverVersions []string) (string, error) {
	// Create a timeout if context doesn't have one
	ctx, cancel := context.WithTimeout(ctx, p.defaultTimeout)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return "", openclaw.NewProtocolError("protocol negotiation timeout", ctx.Err())
		default:
			// Check for matching versions
			for _, clientVer := range p.supportedVersions {
				for _, serverVer := range serverVersions {
					if clientVer == serverVer {
						return clientVer, nil
					}
				}
			}
			// No match found
			return "", openclaw.NewProtocolError("no matching protocol version", nil)
		}
	}
}

// ErrNoMatchingProtocol is a sentinel error for protocol negotiation failures
// Use errors.Is() to check for this specific error
var ErrNoMatchingProtocol = errors.New("no matching protocol version")
