// pkg/openclaw/connection/tls_test.go
package connection

import (
	"testing"
)

func TestTlsValidator_Validate_NilConfig(t *testing.T) {
	v := NewTlsValidator(nil)

	err := v.Validate()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestTlsValidator_Validate_MissingCAFile(t *testing.T) {
	v := NewTlsValidator(&TLSConfig{
		CAFile: "/nonexistent/ca.pem",
	})

	err := v.Validate()
	if err == nil {
		t.Error("expected error for missing CA file")
	}
}

func TestTlsValidator_Validate_IncompleteClientCert(t *testing.T) {
	v := NewTlsValidator(&TLSConfig{
		CertFile: "/path/to/cert.pem",
		// KeyFile missing
	})

	err := v.Validate()
	if err == nil {
		t.Error("expected error for incomplete client cert")
	}
}

func TestTlsValidator_Validate_ValidConfig(t *testing.T) {
	// Create temp files for testing
	// In real tests, use temp files

	v := NewTlsValidator(&TLSConfig{
		InsecureSkipVerify: true,
		ServerName:         "example.com",
	})

	err := v.Validate()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestTlsValidator_GetTLSConfig_Insecure(t *testing.T) {
	v := NewTlsValidator(&TLSConfig{
		InsecureSkipVerify: true,
		ServerName:        "example.com",
	})

	config, err := v.GetTLSConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if config.InsecureSkipVerify != true {
		t.Error("InsecureSkipVerify not set correctly")
	}
	if config.ServerName != "example.com" {
		t.Error("ServerName not set correctly")
	}
}

func TestTlsValidator_GetTLSConfig_NoConfig(t *testing.T) {
	v := NewTlsValidator(nil)

	config, err := v.GetTLSConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if config == nil {
		t.Fatal("expected non-nil config")
	}
}
