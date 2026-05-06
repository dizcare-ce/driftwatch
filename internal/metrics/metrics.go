// Package metrics provides lightweight in-process counters for drift
// detection runs, allowing operators to observe daemon behaviour without
// an external metrics backend.
package metrics

import (
	"fmt"
	"io"
	"sync"
	"sync/atomic"
	"time"
)

// Counters holds cumulative statistics for a single daemon lifetime.
type Counters struct {
	Runs        atomic.Int64
	Drifted     atomic.Int64
	Clean       atomic.Int64
	Errors      atomic.Int64
	lastRunAt   atomic.Value // stores time.Time
	mu          sync.Mutex
}

// New returns a zeroed Counters instance ready for use.
func New() *Counters {
	return &Counters{}
}

// RecordRun increments the run counter, records the timestamp, and
// increments either the drifted or clean counter based on hasDrift.
func (c *Counters) RecordRun(hasDrift bool, err error) {
	c.Runs.Add(1)
	c.lastRunAt.Store(time.Now())

	switch {
	case err != nil:
		c.Errors.Add(1)
	case hasDrift:
		c.Drifted.Add(1)
	default:
		c.Clean.Add(1)
	}
}

// LastRunAt returns the time of the most recent run, and false if no
// run has been recorded yet.
func (c *Counters) LastRunAt() (time.Time, bool) {
	v := c.lastRunAt.Load()
	if v == nil {
		return time.Time{}, false
	}
	return v.(time.Time), true
}

// Write formats the current counters as plain text and writes them to w.
func (c *Counters) Write(w io.Writer) error {
	last := "never"
	if t, ok := c.LastRunAt(); ok {
		last = t.UTC().Format(time.RFC3339)
	}

	_, err := fmt.Fprintf(w,
		"runs=%d drifted=%d clean=%d errors=%d last_run=%s\n",
		c.Runs.Load(),
		c.Drifted.Load(),
		c.Clean.Load(),
		c.Errors.Load(),
		last,
	)
	return err
}
