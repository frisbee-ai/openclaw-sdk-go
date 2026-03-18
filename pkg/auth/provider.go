// Package auth provides authentication types and handlers for OpenClaw SDK.
//
// This package provides:
//   - CredentialsProvider: Interface for credential sources
//   - StaticCredentialsProvider: Simple static credential implementation
//   - AuthHandler: Interface for authentication logic
//   - StaticAuthHandler: Simple static authentication implementation
package auth

import "errors"

// CredentialsProvider provides credentials for authentication.
// Implement this interface to provide custom credential sources.
type CredentialsProvider interface {
	// GetCredentials returns credentials map
	GetCredentials() (map[string]string, error)
}

// StaticCredentialsProvider provides static credentials.
// It stores credentials in memory and returns them on request.
type StaticCredentialsProvider struct {
	credentials map[string]string
}

// NewStaticCredentialsProvider creates a new static credentials provider.
// Returns error if credentials is nil or empty.
func NewStaticCredentialsProvider(credentials map[string]string) (*StaticCredentialsProvider, error) {
	if credentials == nil {
		return nil, errors.New("credentials cannot be nil")
	}
	if len(credentials) == 0 {
		return nil, errors.New("credentials cannot be empty")
	}
	return &StaticCredentialsProvider{credentials: credentials}, nil
}

// GetCredentials returns the stored credentials map.
func (p *StaticCredentialsProvider) GetCredentials() (map[string]string, error) {
	return p.credentials, nil
}
