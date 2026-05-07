// Package throttle provides a token-bucket style throttle that limits
// how many drift-check runs can be triggered within a rolling window.
package throttle

import (
	"sync"
	"time"
)

// Throttle tracks run counts within a sliding window and rejects calls
// that exceed the configured maximum.
type Throttle struct {
	mu       sync.Mutex
	window   time.Duration
	maxRuns  int
	timestamps []time.Time
	now      func() time.Time
}

// New returns a Throttle that allows at most maxRuns calls within window.
func New(window time.Duration, maxRuns int) *Throttle {
	return &Throttle{
		window:  window,
		maxRuns: maxRuns,
		now:     time.Now,
	}
}

// Allow returns true if the call is permitted under the current quota.
// It prunes stale timestamps and records the new call when allowed.
func (t *Throttle) Allow() bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	cutoff := t.now().Add(-t.window)
	active := t.timestamps[:0]
	for _, ts := range t.timestamps {
		if ts.After(cutoff) {
			active = append(active, ts)
		}
	}
	t.timestamps = active

	if len(t.timestamps) >= t.maxRuns {
		return false
	}

	t.timestamps = append(t.timestamps, t.now())
	return true
}

// Remaining returns the number of additional calls permitted in the
// current window without pruning the timestamp list.
func (t *Throttle) Remaining() int {
	t.mu.Lock()
	defer t.mu.Unlock()

	cutoff := t.now().Add(-t.window)
	count := 0
	for _, ts := range t.timestamps {
		if ts.After(cutoff) {
			count++
		}
	}
	rem := t.maxRuns - count
	if rem < 0 {
		return 0
	}
	return rem
}

// Reset clears all recorded timestamps, restoring the full quota.
func (t *Throttle) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.timestamps = nil
}
