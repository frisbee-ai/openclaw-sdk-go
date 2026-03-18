# Phase 2: Authentication Module

**Project Structure:** Go module in root, source files in `pkg/openclaw/` directory

**Files:**
- Create: `pkg/openclaw/auth/provider.go`, `pkg/openclaw/auth/provider_test.go`
- Create: `pkg/openclaw/auth/handler.go`, `pkg/openclaw/auth/handler_test.go`

**Depends on:** Phase 1 (pkg/openclaw/types.go, pkg/openclaw/errors.go)

---

## Task 2.1: CredentialsProvider Interface

- [ ] **Step 1: Create auth directory and provider.go**

```bash
mkdir -p pkg/openclaw/auth
```

```go
// auth/provider.go
package auth

import "errors"

// CredentialsProvider provides credentials for authentication
type CredentialsProvider interface {
	// GetCredentials returns credentials map
	GetCredentials() (map[string]string, error)
}

// StaticCredentialsProvider provides static credentials
type StaticCredentialsProvider struct {
	credentials map[string]string
}

// NewStaticCredentialsProvider creates a new static credentials provider
// Returns error if credentials is nil or empty
func NewStaticCredentialsProvider(credentials map[string]string) (*StaticCredentialsProvider, error) {
	if credentials == nil {
		return nil, errors.New("credentials cannot be nil")
	}
	if len(credentials) == 0 {
		return nil, errors.New("credentials cannot be empty")
	}
	return &StaticCredentialsProvider{credentials: credentials}, nil
}

func (p *StaticCredentialsProvider) GetCredentials() (map[string]string, error) {
	return p.credentials, nil
}
```

- [ ] **Step 2: Write comprehensive tests**

```go
// auth/provider_test.go
package auth

import (
	"testing"
)

func TestStaticCredentialsProvider(t *testing.T) {
	creds := map[string]string{
		"username": "testuser",
		"password": "testpass",
	}
	provider, err := NewStaticCredentialsProvider(creds)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := provider.GetCredentials()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["username"] != "testuser" {
		t.Errorf("expected 'testuser', got '%s'", got["username"])
	}
}

func TestStaticCredentialsProvider_Nil(t *testing.T) {
	_, err := NewStaticCredentialsProvider(nil)
	if err == nil {
		t.Error("expected error for nil credentials")
	}
}

func TestStaticCredentialsProvider_Empty(t *testing.T) {
	_, err := NewStaticCredentialsProvider(map[string]string{})
	if err == nil {
		t.Error("expected error for empty credentials")
	}
}

// Compile-time check: StaticCredentialsProvider implements CredentialsProvider
var _ CredentialsProvider = (*StaticCredentialsProvider)(nil)
```

- [ ] **Step 3: Run tests and commit**

Run: `go test -v ./pkg/openclaw/auth/...`
Commit: `git add pkg/openclaw/auth/ && git commit -m "feat: add CredentialsProvider interface with validation"`

---

## Task 2.2: AuthHandler

- [ ] **Step 1: Write handler.go**

```go
// auth/handler.go
package auth

import (
	"context"
	"errors"
)

// ErrNoCredentials is returned when no credentials are provided
var ErrNoCredentials = errors.New("no credentials provided")

// AuthHandler handles authentication
type AuthHandler interface {
	// Authenticate performs authentication and returns credentials
	Authenticate(ctx context.Context) (CredentialsProvider, error)
}

// StaticAuthHandler is a simple auth handler that returns static credentials
type StaticAuthHandler struct {
	credentials map[string]string
}

// NewStaticAuthHandler creates a new static auth handler
// Returns error if credentials is nil or empty
func NewStaticAuthHandler(credentials map[string]string) (*StaticAuthHandler, error) {
	if credentials == nil {
		return nil, ErrNoCredentials
	}
	if len(credentials) == 0 {
		return nil, ErrNoCredentials
	}
	return &StaticAuthHandler{credentials: credentials}, nil
}

func (h *StaticAuthHandler) Authenticate(ctx context.Context) (CredentialsProvider, error) {
	// Check for context cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	return NewStaticCredentialsProvider(h.credentials)
}
```

- [ ] **Step 2: Write handler tests**

```go
// auth/handler_test.go
package auth

import (
	"context"
	"testing"
	"time"
)

func TestStaticAuthHandler(t *testing.T) {
	creds := map[string]string{
		"username": "testuser",
		"password": "testpass",
	}
	handler, err := NewStaticAuthHandler(creds)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	provider, err := handler.Authenticate(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := provider.GetCredentials()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["username"] != "testuser" {
		t.Errorf("expected 'testuser', got '%s'", got["username"])
	}
}

func TestStaticAuthHandler_NilCredentials(t *testing.T) {
	_, err := NewStaticAuthHandler(nil)
	if err == nil {
		t.Error("expected error for nil credentials")
	}
}

func TestStaticAuthHandler_EmptyCredentials(t *testing.T) {
	_, err := NewStaticAuthHandler(map[string]string{})
	if err == nil {
		t.Error("expected error for empty credentials")
	}
}

func TestStaticAuthHandler_ContextCancellation(t *testing.T) {
	creds := map[string]string{
		"username": "testuser",
		"password": "testpass",
	}
	handler, err := NewStaticAuthHandler(creds)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err = handler.Authenticate(ctx)
	if err == nil {
		t.Error("expected error for cancelled context")
	}
}

func TestStaticAuthHandler_ContextTimeout(t *testing.T) {
	creds := map[string]string{
		"username": "testuser",
		"password": "testpass",
	}
	handler, err := NewStaticAuthHandler(creds)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	// Wait for context to expire
	time.Sleep(100 * time.Millisecond)

	_, err = handler.Authenticate(ctx)
	if err == nil {
		t.Error("expected error for timed out context")
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expected DeadlineExceeded error, got %v", err)
	}
}

// Compile-time check: StaticAuthHandler implements AuthHandler
var _ AuthHandler = (*StaticAuthHandler)(nil)
```

- [ ] **Step 3: Run tests and commit**

Run: `go test -v ./pkg/openclaw/auth/...`
Commit: `git add pkg/openclaw/auth/handler.go pkg/openclaw/auth/handler_test.go && git commit -m "feat: add AuthHandler interface with validation"`

---

## Phase 2 Complete

After this phase, you should have:
- `pkg/openclaw/auth/provider.go` - CredentialsProvider interface with validation
- `pkg/openclaw/auth/provider_test.go` - Comprehensive provider tests
- `pkg/openclaw/auth/handler.go` - AuthHandler interface with context support
- `pkg/openclaw/auth/handler_test.go` - Comprehensive handler tests

All code should compile and tests should pass.
