// Package driftbudget enforces a maximum allowable number of drifted fields
// across a set of results within a rolling time window. When the budget is
// exhausted the caller is expected to escalate or halt remediation attempts.
package driftbudget

import (
	"fmt"
	"sync"
	"time"

	"driftwatch/internal/drift"
)

// Budget tracks drift consumption against a configurable limit.
type Budget struct {
	mu      sync.Mutex
	limit   int
	window  time.Duration
	entries []entry
}

type entry struct {
	recordedAt time.Time
	count      int
}

// New returns a Budget that allows at most limit drifted fields within window.
func New(limit int, window time.Duration) (*Budget, error) {
	if limit <= 0 {
		return nil, fmt.Errorf("driftbudget: limit must be positive, got %d", limit)
	}
	if window <= 0 {
		return nil, fmt.Errorf("driftbudget: window must be positive, got %s", window)
	}
	return &Budget{limit: limit, window: window}, nil
}

// Record adds the drifted field count from results to the budget ledger.
// It returns an error if the budget would be exceeded after recording.
func (b *Budget) Record(results []drift.Result) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
	b.evict(now)

	count := 0
	for _, r := range results {
		count += len(r.Diffs)
	}

	current := b.total()
	if current+count > b.limit {
		return fmt.Errorf("driftbudget: budget exceeded — %d drifted fields recorded in window, limit is %d",
			current+count, b.limit)
	}

	if count > 0 {
		b.entries = append(b.entries, entry{recordedAt: now, count: count})
	}
	return nil
}

// Remaining returns the number of drifted fields still permitted in the
// current window.
func (b *Budget) Remaining() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.evict(time.Now())
	remaining := b.limit - b.total()
	if remaining < 0 {
		return 0
	}
	return remaining
}

// Reset clears all ledger entries, restoring the full budget.
func (b *Budget) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.entries = nil
}

func (b *Budget) evict(now time.Time) {
	cutoff := now.Add(-b.window)
	filtered := b.entries[:0]
	for _, e := range b.entries {
		if e.recordedAt.After(cutoff) {
			filtered = append(filtered, e)
		}
	}
	b.entries = filtered
}

func (b *Budget) total() int {
	n := 0
	for _, e := range b.entries {
		n += e.count
	}
	return n
}
