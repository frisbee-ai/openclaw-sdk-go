# Phase 2: Authentication Module

**Files:**
- Create: `auth/provider.go`, `auth/provider_test.go`
- Create: `auth/handler.go`

**Depends on:** Phase 1 (types.go)

---

## Task 2.1: CredentialsProvider Interface

- [ ] **Step 1: Create auth directory and provider.go**

```bash
mkdir -p auth
```

```go
// auth/provider.go
package auth

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
func NewStaticCredentialsProvider(credentials map[string]string) *StaticCredentialsProvider {
	return &StaticCredentialsProvider{credentials: credentials}
}

func (p *StaticCredentialsProvider) GetCredentials() (map[string]string, error) {
	return p.credentials, nil
}
```

- [ ] **Step 2: Write test**

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
	provider := NewStaticCredentialsProvider(creds)
	got, err := provider.GetCredentials()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["username"] != "testuser" {
		t.Errorf("expected 'testuser', got '%s'", got["username"])
	}
}
```

- [ ] **Step 3: Run tests and commit**

Run: `go test -v ./auth/...`
Commit: `git add auth/ && git commit -m "feat: add CredentialsProvider interface"`

---

## Task 2.2: AuthHandler

- [ ] **Step 1: Write handler.go**

```go
// auth/handler.go
package auth

import (
	"context"
)

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
func NewStaticAuthHandler(credentials map[string]string) *StaticAuthHandler {
	return &StaticAuthHandler{credentials: credentials}
}

func (h *StaticAuthHandler) Authenticate(ctx context.Context) (CredentialsProvider, error) {
	return NewStaticCredentialsProvider(h.credentials), nil
}
```

- [ ] **Step 2: Commit**

```bash
git add auth/handler.go
git commit -m "feat: add AuthHandler interface"
```

---

## Phase 2 Complete

After this phase, you should have:
- `auth/provider.go` - CredentialsProvider interface
- `auth/handler.go` - AuthHandler interface

All code should compile and tests should pass.
