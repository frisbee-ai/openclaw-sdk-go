# Phase 6: Events Module

**Files:**
- Create: `events/tick.go`, `events/tick_test.go`
- Create: `events/gap.go`, `events/gap_test.go`

**Depends on:** Phase 1 (types.go)

---

## Task 6.1: Tick Monitor

- [ ] **Step 1: Create events directory and tick.go**

```bash
mkdir -p events
```

```go
// events/tick.go
package events

import (
	"context"
	"sync"
	"time"
)

// TickMonitor monitors connection heartbeat
type TickMonitor struct {
	interval   time.Duration
	ticker     *time.Ticker
	tickCh     chan time.Time
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
	onTick     func(time.Time)
	onTimeout  func()
	timeout    time.Duration
}

// NewTickMonitor creates a new tick monitor
func NewTickMonitor(interval time.Duration, timeout time.Duration) *TickMonitor {
	ctx, cancel := context.WithCancel(context.Background())
	return &TickMonitor{
		interval: interval,
		ticker:   time.NewTicker(interval),
		tickCh:   make(chan time.Time, 1),
		ctx:      ctx,
		cancel:   cancel,
		timeout:  timeout,
	}
}

// SetOnTick sets the tick callback
func (tm *TickMonitor) SetOnTick(f func(time.Time)) {
	tm.onTick = f
}

// SetOnTimeout sets the timeout callback
func (tm *TickMonitor) SetOnTimeout(f func()) {
	tm.onTimeout = f
}

// Start begins the tick monitoring
func (tm *TickMonitor) Start() {
	tm.wg.Add(1)
	go tm.run()
}

// Stop stops the tick monitoring
func (tm *TickMonitor) Stop() {
	tm.cancel()
	tm.ticker.Stop()
	tm.wg.Wait()
	close(tm.tickCh)
}

func (tm *TickMonitor) run() {
	defer tm.wg.Done()

	var lastTick time.Time
	for {
		select {
		case <-tm.ctx.Done():
			return
		case tick := <-tm.ticker.C:
			lastTick = tick
			select {
			case tm.tickCh <- tick:
			default:
			}
			if tm.onTick != nil {
				tm.onTick(tick)
			}
		case <-time.After(tm.timeout):
			if tm.onTimeout != nil && !lastTick.IsZero() {
				tm.onTimeout()
			}
		}
	}
}

// TickChannel returns the tick channel
func (tm *TickMonitor) TickChannel() <-chan time.Time {
	return tm.tickCh
}
```

- [ ] **Step 2: Write test**

```go
// events/tick_test.go
package events

import (
	"testing"
	"time"
)

func TestTickMonitor(t *testing.T) {
	monitor := NewTickMonitor(50*time.Millisecond, 100*time.Millisecond)
	tickCount := 0

	monitor.SetOnTick(func(t time.Time) {
		tickCount++
	})

	monitor.Start()
	time.Sleep(200 * time.Millisecond)
	monitor.Stop()

	if tickCount == 0 {
		t.Error("expected at least one tick")
	}
}
```

- [ ] **Step 3: Run tests and commit**

Run: `go test -v ./events/...`
Commit: `git add events/tick.go events/tick_test.go && git commit -m "feat: add tick monitor"`

---

## Task 6.2: Gap Detector

- [ ] **Step 1: Write gap.go**

```go
// events/gap.go
package events

import (
	"sync"
)

// GapDetector detects message gaps
type GapDetector struct {
	mu           sync.Mutex
	expectedSeq  uint64
	detectedGaps []Gap
	onGap        func(start, end uint64)
}

// Gap represents a detected gap
type Gap struct {
	Start uint64
	End   uint64
}

// NewGapDetector creates a new gap detector
func NewGapDetector() *GapDetector {
	return &GapDetector{
		expectedSeq:  0,
		detectedGaps: make([]Gap, 0),
	}
}

// SetOnGap sets the gap callback
func (gd *GapDetector) SetOnGap(f func(start, end uint64)) {
	gd.onGap = f
}

// Record records a message sequence number
func (gd *GapDetector) Record(seq uint64) {
	gd.mu.Lock()
	defer gd.mu.Unlock()

	if gd.expectedSeq == 0 {
		gd.expectedSeq = seq + 1
		return
	}

	if seq > gd.expectedSeq {
		gap := Gap{Start: gd.expectedSeq, End: seq - 1}
		gd.detectedGaps = append(gd.detectedGaps, gap)
		if gd.onGap != nil {
			gd.onGap(gap.Start, gap.End)
		}
	}

	gd.expectedSeq = seq + 1
}

// Gaps returns the detected gaps
func (gd *GapDetector) Gaps() []Gap {
	gd.mu.Lock()
	defer gd.mu.Unlock()
	return gd.detectedGaps
}

// Reset resets the gap detector
func (gd *GapDetector) Reset() {
	gd.mu.Lock()
	defer gd.mu.Unlock()
	gd.expectedSeq = 0
	gd.detectedGaps = nil
}
```

- [ ] **Step 2: Write test**

```go
// events/gap_test.go
package events

import (
	"testing"
)

func TestGapDetector(t *testing.T) {
	detector := NewGapDetector()
	gapDetected := false

	detector.SetOnGap(func(start, end uint64) {
		gapDetected = true
		if start != 2 || end != 3 {
			t.Errorf("expected gap 2-3, got %d-%d", start, end)
		}
	})

	detector.Record(1)
	detector.Record(4) // Should detect gap 2-3

	if !gapDetected {
		t.Error("expected gap to be detected")
	}
}
```

- [ ] **Step 3: Run tests and commit**

Run: `go test -v ./events/...`
Commit: `git add events/gap.go events/gap_test.go && git commit -m "feat: add gap detector"`

---

## Phase 6 Complete

After this phase, you should have:
- `events/tick.go` - Tick monitor (heartbeat)
- `events/gap.go` - Gap detector

All code should compile and tests should pass.
