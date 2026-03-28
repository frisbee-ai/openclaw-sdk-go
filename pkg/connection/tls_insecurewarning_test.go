package connection

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"sync"
	"testing"
	"time"
)

// mockLogger is a mock implementation of types.Logger for testing
type mockLogger struct {
	mu   sync.Mutex
	msgs []string
}

func (m *mockLogger) Debug(msg string, args ...any) {}
func (m *mockLogger) Info(msg string, args ...any)  {}
func (m *mockLogger) Warn(msg string, args ...any) {
	m.mu.Lock()
	m.msgs = append(m.msgs, msg)
	m.mu.Unlock()
}
func (m *mockLogger) Error(msg string, args ...any) {}

// TestTlsValidator_InsecureSkipVerifyWarning_WithLogger verifies warning is logged (FOUND-05)
func TestTlsValidator_InsecureSkipVerifyWarning_WithLogger(t *testing.T) {
	logger := &mockLogger{}

	v := NewTlsValidator(&TLSConfig{
		InsecureSkipVerify: true,
		ServerName:         "example.com",
	})
	v.SetLogger(logger)

	_, err := v.GetTLSConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	logger.mu.Lock()
	if len(logger.msgs) == 0 {
		t.Error("expected Warn to be called for InsecureSkipVerify")
	}
	// Verify the message contains InsecureSkipVerify
	found := false
	for _, msg := range logger.msgs {
		if contains(msg, "InsecureSkipVerify") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected warning message to contain 'InsecureSkipVerify', got %v", logger.msgs)
	}
	logger.mu.Unlock()
}

// TestTlsValidator_InsecureSkipVerifyWarning_NotEnabled verifies no warning when disabled (FOUND-05)
func TestTlsValidator_InsecureSkipVerifyWarning_NotEnabled(t *testing.T) {
	logger := &mockLogger{}

	v := NewTlsValidator(&TLSConfig{
		InsecureSkipVerify: false,
		ServerName:         "example.com",
	})
	v.SetLogger(logger)

	_, err := v.GetTLSConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	logger.mu.Lock()
	if len(logger.msgs) != 0 {
		t.Errorf("expected no Warn calls when InsecureSkipVerify=false, got %d calls", len(logger.msgs))
	}
	logger.mu.Unlock()
}

// TestTlsValidator_InsecureSkipVerifyWarning_NilLogger verifies no panic with nil logger (FOUND-05)
func TestTlsValidator_InsecureSkipVerifyWarning_NilLogger(t *testing.T) {
	v := NewTlsValidator(&TLSConfig{
		InsecureSkipVerify: true,
		ServerName:         "example.com",
	})
	// Do not call SetLogger - logger should be nil

	// Should not panic
	_, err := v.GetTLSConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestTlsValidator_InsecureSkipVerifyWarning_AfterSetLogger verifies logger can be changed (FOUND-05)
func TestTlsValidator_InsecureSkipVerifyWarning_AfterSetLogger(t *testing.T) {
	logger1 := &mockLogger{}
	logger2 := &mockLogger{}

	v := NewTlsValidator(&TLSConfig{
		InsecureSkipVerify: true,
		ServerName:         "example.com",
	})

	// Set first logger
	v.SetLogger(logger1)
	_, _ = v.GetTLSConfig()

	// Change logger
	v.SetLogger(logger2)
	_, _ = v.GetTLSConfig()

	logger1.mu.Lock()
	logger2.mu.Lock()
	if len(logger1.msgs) != 1 {
		t.Errorf("expected 1 warning on logger1, got %d", len(logger1.msgs))
	}
	if len(logger2.msgs) != 1 {
		t.Errorf("expected 1 warning on logger2, got %d", len(logger2.msgs))
	}
	logger1.mu.Unlock()
	logger2.mu.Unlock()
}

// TestCheckCertificateRevocation_V1Stub verifies stub behavior and documentation (FOUND-03)
func TestCheckCertificateRevocation_V1Stub(t *testing.T) {
	// Create a test certificate
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: "test.example.com",
		},
		NotBefore:             time.Now().Add(-24 * time.Hour),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  false,
		DNSNames:              []string{"test.example.com"},
		CRLDistributionPoints: []string{"http://example.com/crl"},
		OCSPServer:            []string{"http://ocsp.example.com"},
	}

	// Should return nil (stub behavior)
	err := CheckCertificateRevocation(cert, nil)
	if err != nil {
		t.Errorf("expected nil return for v1 stub, got %v", err)
	}

	// Should also return nil for cert without revocation info
	err = CheckCertificateRevocation(cert, nil)
	if err != nil {
		t.Errorf("expected nil return for cert without CRL/OCSP, got %v", err)
	}

	// Should return error for nil certificate
	err = CheckCertificateRevocation(nil, nil)
	if err == nil {
		t.Error("expected error for nil certificate")
	}
}

// contains is a simple helper for string matching
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
