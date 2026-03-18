# Phase 8: Utils Module

**Files:**
- Create: `pkg/openclaw/utils/timeout.go`

**Project Structure:** Go module in root, source files in `pkg/openclaw/` directory

**Depends on:** Phase 1

---

## Task 8.1: Timeout Manager

- [ ] **Step 1: Create utils directory and timeout.go**

```bash
mkdir -p pkg/openclaw/utils
```

```go
// pkg/openclaw/utils/timeout.go
package utils

import (
	"context"
	"time"
)

// TimeoutManager manages timeouts
type TimeoutManager struct {
	defaultTimeout time.Duration
}

// NewTimeoutManager creates a new timeout manager
// If defaultTimeout is negative, it will be set to 0 (no timeout)
func NewTimeoutManager(defaultTimeout time.Duration) *TimeoutManager {
	if defaultTimeout < 0 {
		defaultTimeout = 0
	}
	return &TimeoutManager{defaultTimeout: defaultTimeout}
}

// WithTimeout wraps a context with timeout
func (tm *TimeoutManager) WithTimeout(parent context.Context) (context.Context, context.CancelFunc) {
	if tm.defaultTimeout > 0 {
		return context.WithTimeout(parent, tm.defaultTimeout)
	}
	return context.WithCancel(parent)
}

// WithCustomTimeout wraps a context with a custom timeout
// If timeout is zero or negative, behaves like WithCancel (no timeout)
func (tm *TimeoutManager) WithCustomTimeout(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if timeout <= 0 {
		return context.WithCancel(parent)
	}
	return context.WithTimeout(parent, timeout)
}

// DefaultTimeoutManager is a global default
var DefaultTimeoutManager = NewTimeoutManager(30 * time.Second)
```

- [ ] **Step 2: Write test**

```go
// pkg/openclaw/utils/timeout_test.go
package utils

import (
	"context"
	"testing"
	"time"
)

func TestNewTimeoutManager_Negative(t *testing.T) {
	tm := NewTimeoutManager(-1 * time.Second)
	if tm.defaultTimeout != 0 {
		t.Errorf("expected 0, got %v", tm.defaultTimeout)
	}
}

func TestNewTimeoutManager_Zero(t *testing.T) {
	tm := NewTimeoutManager(0)
	if tm.defaultTimeout != 0 {
		t.Errorf("expected 0, got %v", tm.defaultTimeout)
	}
}

func TestWithTimeout_Default(t *testing.T) {
	tm := NewTimeoutManager(10 * time.Second)
	ctx, cancel := tm.WithTimeout(context.Background())
	defer cancel()

	if ctx == nil {
		t.Error("expected non-nil context")
	}
}

func TestWithCustomTimeout_Negative(t *testing.T) {
	tm := NewTimeoutManager(10 * time.Second)
	ctx, cancel := tm.WithCustomTimeout(context.Background(), -5*time.Second)
	defer cancel()

	// Negative timeout should behave like no timeout (WithCancel)
	if ctx == nil {
		t.Error("expected non-nil context")
	}
}

func TestWithCustomTimeout_Zero(t *testing.T) {
	tm := NewTimeoutManager(10 * time.Second)
	ctx, cancel := tm.WithCustomTimeout(context.Background(), 0)
	defer cancel()

	// Zero timeout should behave like no timeout
	if ctx == nil {
		t.Error("expected non-nil context")
	}
}
```

- [ ] **Step 3: Run tests**

Run: `go test -v ./pkg/openclaw/utils/... -race`

- [ ] **Step 4: Commit**

```bash
git add pkg/openclaw/utils/timeout.go go.mod
git commit -m "feat: add timeout manager"
```

---

## Phase 8 Complete

After this phase, you should have:
- `pkg/openclaw/utils/timeout.go` - Timeout manager

All code should compile.
