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
	timeout    time.Duration
	ticker     *time.Ticker
	timer      *time.Timer
	tickCh     chan time.Time
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
	mu         sync.RWMutex
	onTick     func(time.Time)
	onTimeout  func()
	running    bool
	stopped   chan struct{}
}

// NewTickMonitor creates a new tick monitor
// Returns error if interval or timeout is zero or negative
func NewTickMonitor(interval time.Duration, timeout time.Duration) (*TickMonitor, error) {
	if interval <= 0 {
		return nil, &ValidationError{Field: "interval", Message: "must be positive"}
	}
	if timeout <= 0 {
		return nil, &ValidationError{Field: "timeout", Message: "must be positive"}
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &TickMonitor{
		interval: interval,
		timeout:  timeout,
		ticker:   time.NewTicker(interval),
		timer:    time.NewTimer(timeout),
		tickCh:   make(chan time.Time, 1),
		ctx:      ctx,
		cancel:   cancel,
		stopped:  make(chan struct{}),
	}, nil
}

// ValidationError represents validation failure
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Message
}

// SetOnTick sets the tick callback (thread-safe)
func (tm *TickMonitor) SetOnTick(f func(time.Time)) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.onTick = f
}

// SetOnTimeout sets the timeout callback (thread-safe)
func (tm *TickMonitor) SetOnTimeout(f func()) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.onTimeout = f
}

// Start begins the tick monitoring
func (tm *TickMonitor) Start() {
	tm.mu.Lock()
	if tm.running {
		tm.mu.Unlock()
		return
	}
	tm.running = true
	tm.mu.Unlock()

	tm.wg.Add(1)
	go tm.run()
}

// Stop stops the tick monitoring (idempotent)
func (tm *TickMonitor) Stop() {
	tm.mu.Lock()
	if !tm.running {
		tm.mu.Unlock()
		return
	}
	tm.running = false
	tm.mu.Unlock()

	tm.cancel()

	// Stop timer and ticker
	if !tm.timer.Stop() {
		// Drain timer channel if it fired
		select {
		case <-tm.timer.C:
		default:
		}
	}
	tm.ticker.Stop()

	// Wait for goroutine to finish
	tm.wg.Wait()

	// Close channel only after goroutine is done
	close(tm.stopped)
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
			// Call callback with lock
			tm.mu.RLock()
			onTick := tm.onTick
			tm.mu.RUnlock()
			if onTick != nil {
				onTick(tick)
			}
		case <-tm.timer.C:
			// Timer fired - timeout occurred
			tm.mu.RLock()
			onTimeout := tm.onTimeout
			tm.mu.RUnlock()
			if onTimeout != nil && !lastTick.IsZero() {
				onTimeout()
			}
			// Reset timer for next timeout
			if !lastTick.IsZero() {
				tm.timer.Reset(tm.timeout)
			}
		}
	}
}

// TickChannel returns the tick channel
func (tm *TickMonitor) TickChannel() <-chan time.Time {
	return tm.tickCh
}

// IsRunning returns whether the monitor is running
func (tm *TickMonitor) IsRunning() bool {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return tm.running
}
```

- [ ] **Step 2: Write comprehensive test**

```go
// events/tick_test.go
package events

import (
	"sync"
	"testing"
	"time"
)

func TestNewTickMonitor_InvalidInterval(t *testing.T) {
	_, err := NewTickMonitor(0, time.Second)
	if err == nil {
		t.Error("expected error for zero interval")
	}

	_, err = NewTickMonitor(-1, time.Second)
	if err == nil {
		t.Error("expected error for negative interval")
	}
}

func TestNewTickMonitor_InvalidTimeout(t *testing.T) {
	_, err := NewTickMonitor(time.Second, 0)
	if err == nil {
		t.Error("expected error for zero timeout")
	}

	_, err = NewTickMonitor(time.Second, -1)
	if err == nil {
		t.Error("expected error for negative timeout")
	}
}

func TestTickMonitor_StartStop(t *testing.T) {
	monitor, err := NewTickMonitor(50*time.Millisecond, 100*time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	monitor.Start()
	time.Sleep(20 * time.Millisecond)

	if !monitor.IsRunning() {
		t.Error("expected monitor to be running")
	}

	monitor.Stop()

	if monitor.IsRunning() {
		t.Error("expected monitor to be stopped")
	}
}

func TestTickMonitor_Stop_Idempotent(t *testing.T) {
	monitor, err := NewTickMonitor(50*time.Millisecond, 100*time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	monitor.Start()
	time.Sleep(10 * time.Millisecond)

	// Should not panic
	monitor.Stop()
	monitor.Stop()
}

func TestTickMonitor_TickCallback(t *testing.T) {
	monitor, err := NewTickMonitor(20*time.Millisecond, 100*time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var tickCount int
	var mu sync.Mutex

	monitor.SetOnTick(func(t time.Time) {
		mu.Lock()
		tickCount++
		mu.Unlock()
	})

	monitor.Start()
	time.Sleep(60 * time.Millisecond)
	monitor.Stop()

	mu.Lock()
	if tickCount == 0 {
		t.Error("expected at least one tick")
	}
	mu.Unlock()
}

func TestTickMonitor_TimeoutCallback(t *testing.T) {
	monitor, err := NewTickMonitor(100*time.Millisecond, 50*time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	timeoutCalled := false
	monitor.SetOnTimeout(func() {
		timeoutCalled = true
	})

	monitor.Start()
	time.Sleep(120 * time.Millisecond)
	monitor.Stop()

	if !timeoutCalled {
		t.Error("expected timeout callback to be called")
	}
}

func TestTickMonitor_TickChannel(t *testing.T) {
	monitor, err := NewTickMonitor(20*time.Millisecond, 100*time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	monitor.Start()

	select {
	case <-monitor.TickChannel():
		// Expected
	case <-time.After(100 * time.Millisecond):
		t.Error("timeout waiting for tick")
	}

	monitor.Stop()
}

func TestTickMonitor_ConcurrentCallbacks(t *testing.T) {
	monitor, err := NewTickMonitor(10*time.Millisecond, 100*time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	// Set callbacks from different goroutines
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			monitor.SetOnTick(func(t time.Time) {})
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			monitor.SetOnTimeout(func() {})
		}
	}()

	monitor.Start()
	time.Sleep(30 * time.Millisecond)
	monitor.Stop()

	wg.Wait()
}
```

- [ ] **Step 3: Run tests and commit**

Run: `go test -v ./events/... -race`
Commit: `git add events/tick.go events/tick_test.go && git commit -m "feat: add tick monitor with thread-safe callbacks"`

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
		detectedGaps: make([]Gap, 0),
	}
}

// SetOnGap sets the gap callback (thread-safe)
func (gd *GapDetector) SetOnGap(f func(start, end uint64)) {
	gd.mu.Lock()
	defer gd.mu.Unlock()
	gd.onGap = f
}

// Record records a message sequence number
func (gd *GapDetector) Record(seq uint64) {
	gd.mu.Lock()
	defer gd.mu.Unlock()

	// Handle first message
	if gd.expectedSeq == 0 {
		gd.expectedSeq = seq + 1
		return
	}

	// Detect gap (but not for duplicate sequences)
	if seq > gd.expectedSeq {
		gap := Gap{Start: gd.expectedSeq, End: seq - 1}
		gd.detectedGaps = append(gd.detectedGaps, gap)

		// Call callback OUTSIDE the lock to prevent deadlock
		onGap := gd.onGap
		if onGap != nil {
			gd.mu.Unlock()
			onGap(gap.Start, gap.End)
			gd.mu.Lock()
		}
	}

	// Advance expected (handle duplicates by not going backwards)
	if seq+1 > gd.expectedSeq {
		gd.expectedSeq = seq + 1
	}
}

// Gaps returns a copy of detected gaps (thread-safe)
func (gd *GapDetector) Gaps() []Gap {
	gd.mu.Lock()
	defer gd.mu.Unlock()

	// Return a copy to prevent external mutation
	result := make([]Gap, len(gd.detectedGaps))
	copy(result, gd.detectedGaps)
	return result
}

// GapCount returns the number of detected gaps
func (gd *GapDetector) GapCount() int {
	gd.mu.Lock()
	defer gd.mu.Unlock()
	return len(gd.detectedGaps)
}

// Reset resets the gap detector (thread-safe)
func (gd *GapDetector) Reset() {
	gd.mu.Lock()
	defer gd.mu.Unlock()
	gd.expectedSeq = 0
	gd.detectedGaps = nil
	gd.detectedGaps = make([]Gap, 0)
}

// ExpectedSequence returns the next expected sequence number
func (gd *GapDetector) ExpectedSequence() uint64 {
	gd.mu.Lock()
	defer gd.mu.Unlock()
	return gd.expectedSeq
}
```

- [ ] **Step 2: Write comprehensive test**

```go
// events/gap_test.go
package events

import (
	"sync"
	"testing"
)

func TestGapDetector_Basic(t *testing.T) {
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

func TestGapDetector_NoGap(t *testing.T) {
	detector := NewGapDetector()

	detector.Record(1)
	detector.Record(2)
	detector.Record(3)

	gaps := detector.Gaps()
	if len(gaps) != 0 {
		t.Errorf("expected no gaps, got %d", len(gaps))
	}
}

func TestGapDetector_DuplicateSequence(t *testing.T) {
	detector := NewGapDetector()

	detector.Record(1)
	detector.Record(1) // Duplicate
	detector.Record(2)

	gaps := detector.Gaps()
	if len(gaps) != 0 {
		t.Errorf("expected no gaps, got %d", len(gaps))
	}
}

func TestGapDetector_GapsReturnsCopy(t *testing.T) {
	detector := NewGapDetector()

	detector.Record(1)
	detector.Record(4)

	gaps1 := detector.Gaps()
	gaps1[0].Start = 999 // Modify copy

	gaps2 := detector.Gaps()
	if gaps2[0].Start == 999 {
		t.Error("Gaps() should return a copy, not original")
	}
}

func TestGapDetector_GapCount(t *testing.T) {
	detector := NewGapDetector()

	detector.Record(1)
	detector.Record(4)
	detector.Record(10)

	if detector.GapCount() != 2 {
		t.Errorf("expected 2 gaps, got %d", detector.GapCount())
	}
}

func TestGapDetector_Reset(t *testing.T) {
	detector := NewGapDetector()

	detector.Record(1)
	detector.Record(4)

	detector.Reset()

	if detector.GapCount() != 0 {
		t.Error("expected 0 gaps after reset")
	}
	if detector.ExpectedSequence() != 0 {
		t.Error("expected 0 expectedSeq after reset")
	}
}

func TestGapDetector_Concurrent(t *testing.T) {
	detector := NewGapDetector()

	var wg sync.WaitGroup
	wg.Add(2)

	// Concurrent record from multiple goroutines
	go func() {
		defer wg.Done()
		for i := uint64(0); i < 1000; i++ {
			detector.Record(i * 2) // Even numbers
		}
	}()

	go func() {
		defer wg.Done()
		for i := uint64(0); i < 1000; i++ {
			detector.Record(i*2 + 1) // Odd numbers
		}
	}()

	wg.Wait()

	// Should have gaps when interleaved
	gaps := detector.Gaps()
	_ = gaps // Verify no panic
}

func TestGapDetector_ConcurrentCallbacks(t *testing.T) {
	detector := NewGapDetector()

	var wg sync.WaitGroup
	wg.Add(2)

	// Set callback while recording
	detector.SetOnGap(func(start, end uint64) {})

	go func() {
		defer wg.Done()
		for i := uint64(0); i < 100; i++ {
			detector.Record(i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			detector.SetOnGap(func(start, end uint64) {})
		}
	}()

	wg.Wait()
}
```

- [ ] **Step 3: Run tests and commit**

Run: `go test -v ./events/... -race`
Commit: `git add events/gap.go events/gap_test.go && git commit -m "feat: add gap detector with thread-safe operations"`

---

## Phase 6 Complete

After this phase, you should have:
- `events/tick.go` - Tick monitor (heartbeat) with thread-safe callbacks
- `events/tick_test.go` - Comprehensive tick monitor tests
- `events/gap.go` - Gap detector with thread-safe operations
- `events/gap_test.go` - Comprehensive gap detector tests

All code should compile and tests should pass.

Key fixes from review:
1. Fixed race condition in Stop() - proper synchronization with done channel
2. Fixed memory leak - replaced time.After() with time.Timer.Reset()
3. Fixed callback thread-safety - added mutex protection
4. Fixed Gaps() returning internal slice - now returns copy
5. Fixed gap callback deadlock - callback called outside lock
6. Added input validation for NewTickMonitor
7. Added comprehensive test coverage including concurrent tests
8. Added Idempotent Stop() - can be called multiple times safely
9. Added IsRunning() method
