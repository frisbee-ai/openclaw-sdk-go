# Phase 8: Utils Module

**Files:**
- Create: `utils/timeout.go`

**Depends on:** Phase 1

---

## Task 8.1: Timeout Manager

- [ ] **Step 1: Create utils directory and timeout.go**

```bash
mkdir -p utils
```

```go
// utils/timeout.go
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
func NewTimeoutManager(defaultTimeout time.Duration) *TimeoutManager {
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
func (tm *TimeoutManager) WithCustomTimeout(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(parent, timeout)
}

// DefaultTimeoutManager is a global default
var DefaultTimeoutManager = NewTimeoutManager(30 * time.Second)
```

- [ ] **Step 2: Commit**

```bash
git add utils/timeout.go
git commit -m "feat: add timeout manager"
```

---

## Phase 8 Complete

After this phase, you should have:
- `utils/timeout.go` - Timeout manager

All code should compile.
