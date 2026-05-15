// Package shadowmode provides a shadow-run facility that executes drift
// detection without emitting notifications or triggering remediations.
// Results are recorded for comparison against live-mode outputs, making
// it safe to trial new rules or config changes in production.
package shadowmode

import (
	"context"
	"fmt"
	"sync"
	"time"

	"driftwatch/internal/drift"
)

// Record holds the outcome of a single shadow run.
type Record struct {
	RunAt   time.Time
	Results []drift.Result
	Err     error
}

// Runner executes a drift-check function in shadow mode and stores
// the last N records for later inspection.
type Runner struct {
	mu      sync.Mutex
	records []Record
	cap     int
}

// New returns a Runner that retains at most capacity records.
func New(capacity int) *Runner {
	if capacity <= 0 {
		capacity = 10
	}
	return &Runner{cap: capacity}
}

// Run invokes fn inside a shadow context. Panics are recovered and
// surfaced as errors so the host process is never destabilised.
func (r *Runner) Run(ctx context.Context, fn func(ctx context.Context) ([]drift.Result, error)) {
	rec := Record{RunAt: time.Now()}

	func() {
		defer func() {
			if p := recover(); p != nil {
				rec.Err = fmt.Errorf("shadow run panic: %v", p)
			}
		}()
		rec.Results, rec.Err = fn(ctx)
	}()

	r.mu.Lock()
	defer r.mu.Unlock()
	r.records = append(r.records, rec)
	if len(r.records) > r.cap {
		r.records = r.records[len(r.records)-r.cap:]
	}
}

// Last returns up to n most-recent records, oldest first.
func (r *Runner) Last(n int) []Record {
	r.mu.Lock()
	defer r.mu.Unlock()
	if n <= 0 || n > len(r.records) {
		n = len(r.records)
	}
	out := make([]Record, n)
	copy(out, r.records[len(r.records)-n:])
	return out
}

// Reset discards all stored records.
func (r *Runner) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.records = nil
}
