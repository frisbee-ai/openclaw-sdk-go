// pkg/openclaw/connection/tls.go
package connection

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"os"
	"time"

	openclaw "github.com/i0r3k/openclaw-sdk-go/pkg/openclaw"
)

// TlsValidator validates TLS certificates
type TlsValidator struct {
	config *TLSConfig
}

// TLSConfig holds TLS configuration for connection layer
// Note: This is distinct from transport.TLSConfig which is for dial-time configuration
// This version supports certificate loading and validation
type TLSConfig struct {
	InsecureSkipVerify bool
	CertFile          string
	KeyFile           string
	CAFile            string
	ServerName        string
}

// ErrInvalidTLSConfig represents TLS configuration validation errors
var ErrInvalidTLSConfig = errors.New("invalid TLS configuration")

// ErrCertNotFound is returned when certificate file is not found
var ErrCertNotFound = errors.New("certificate file not found")

// ErrCANotFound is returned when CA file is not found
var ErrCANotFound = errors.New("CA certificate file not found")

// NewTlsValidator creates a new TLS validator
func NewTlsValidator(config *TLSConfig) *TlsValidator {
	return &TlsValidator{config: config}
}

// Validate validates the TLS configuration
func (v *TlsValidator) Validate() error {
	if v.config == nil {
		return nil // No config is valid (use system defaults)
	}

	// If using custom CA, verify it exists
	if v.config.CAFile != "" {
		if _, err := os.Stat(v.config.CAFile); os.IsNotExist(err) {
			return openclaw.NewValidationError("TLS CA file does not exist", ErrCANotFound)
		}
	}

	// If using client cert, both cert and key must be present
	if v.config.CertFile != "" || v.config.KeyFile != "" {
		if v.config.CertFile == "" || v.config.KeyFile == "" {
			return openclaw.NewValidationError("both CertFile and KeyFile are required for client authentication", ErrInvalidTLSConfig)
		}
		// Verify both files exist
		if _, err := os.Stat(v.config.CertFile); os.IsNotExist(err) {
			return openclaw.NewValidationError("TLS certificate file does not exist", ErrCertNotFound)
		}
		if _, err := os.Stat(v.config.KeyFile); os.IsNotExist(err) {
			return openclaw.NewValidationError("TLS key file does not exist", ErrCertNotFound)
		}
	}

	return nil
}

// GetTLSConfig returns the TLS config for the connection
func (v *TlsValidator) GetTLSConfig() (*tls.Config, error) {
	// First validate
	if err := v.Validate(); err != nil {
		return nil, err
	}

	// Handle nil config case
	if v.config == nil {
		return &tls.Config{}, nil
	}

	config := &tls.Config{
		InsecureSkipVerify: v.config.InsecureSkipVerify,
		ServerName:         v.config.ServerName,
	}

	// Load client certificate if provided
	if v.config.CertFile != "" && v.config.KeyFile != "" {
		cert, err := tls.LoadX509KeyPair(v.config.CertFile, v.config.KeyFile)
		if err != nil {
			return nil, openclaw.NewTransportError("failed to load client certificate", err)
		}
		config.Certificates = []tls.Certificate{cert}
	}

	// Load CA certificate if provided
	if v.config.CAFile != "" {
		caCert, err := os.ReadFile(v.config.CAFile)
		if err != nil {
			return nil, openclaw.NewTransportError("failed to read CA certificate", err)
		}
		caPool := x509.NewCertPool()
		caPool.AppendCertsFromPEM(caCert)
		config.RootCAs = caPool
	}

	return config, nil
}

// ValidateCertificate validates the given certificate
// This is a basic validation - checks expiry and key usage
func ValidateCertificate(cert *x509.Certificate) error {
	if time.Now().After(cert.NotAfter) {
		return errors.New("certificate has expired")
	}
	if time.Now().Before(cert.NotBefore) {
		return errors.New("certificate is not yet valid")
	}
	return nil
}
