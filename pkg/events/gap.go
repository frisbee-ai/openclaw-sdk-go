// Package events provides event handling utilities for the OpenClaw SDK.
//
// This package provides:
//   - TickMonitor: Connection heartbeat monitoring with timeout detection
//   - GapDetector: Message gap detection for ordered message streams
package events

import (
	"sync"
)

// GapDetector detects message gaps in ordered message streams.
// It tracks sequence numbers and identifies when expected messages are missing.
type GapDetector struct {
	mu           sync.Mutex              // Mutex for thread-safety
	expectedSeq  uint64                  // Next expected sequence number
	detectedGaps []Gap                   // List of detected gaps
	onGap        func(start, end uint64) // Callback for gap detection
}

// Gap represents a detected gap in the message sequence.
// It indicates messages from Start to End (inclusive) were missing.
type Gap struct {
	Start uint64 // Start of the gap (first missing sequence number)
	End   uint64 // End of the gap (last missing sequence number)
}

// NewGapDetector creates a new gap detector with an empty gap list.
func NewGapDetector() *GapDetector {
	return &GapDetector{
		detectedGaps: make([]Gap, 0),
	}
}

// SetOnGap sets the callback function to be called when a gap is detected.
func (gd *GapDetector) SetOnGap(f func(start, end uint64)) {
	gd.mu.Lock()
	defer gd.mu.Unlock()
	gd.onGap = f
}

// Record records a message sequence number.
// It detects gaps when sequence numbers are skipped.
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

// Gaps returns a copy of detected gaps.
// Thread-safe method that returns a copy to prevent external mutation.
func (gd *GapDetector) Gaps() []Gap {
	gd.mu.Lock()
	defer gd.mu.Unlock()

	// Return a copy to prevent external mutation
	result := make([]Gap, len(gd.detectedGaps))
	copy(result, gd.detectedGaps)
	return result
}

// GapCount returns the number of detected gaps.
func (gd *GapDetector) GapCount() int {
	gd.mu.Lock()
	defer gd.mu.Unlock()
	return len(gd.detectedGaps)
}

// Reset resets the gap detector to its initial state.
// Clears all detected gaps and resets the expected sequence.
func (gd *GapDetector) Reset() {
	gd.mu.Lock()
	defer gd.mu.Unlock()
	gd.expectedSeq = 0
	gd.detectedGaps = nil
	gd.detectedGaps = make([]Gap, 0)
}

// ExpectedSequence returns the next expected sequence number.
func (gd *GapDetector) ExpectedSequence() uint64 {
	gd.mu.Lock()
	defer gd.mu.Unlock()
	return gd.expectedSeq
}
