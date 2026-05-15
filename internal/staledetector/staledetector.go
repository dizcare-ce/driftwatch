// Package staledetector identifies services whose drift results have not been
// refreshed within a configurable staleness window. A result is considered
// stale when its timestamp is older than the configured max age, which may
// indicate a scanner failure or a silently-dropped job.
package staledetector

import (
	"time"

	"driftwatch/internal/drift"
)

// Entry describes a single stale service.
type Entry struct {
	Service  string
	LastSeen time.Time
	Age      time.Duration
}

// Detector checks a set of results for staleness.
type Detector struct {
	maxAge time.Duration
	now    func() time.Time
}

// New returns a Detector that flags results older than maxAge.
func New(maxAge time.Duration) *Detector {
	return &Detector{maxAge: maxAge, now: time.Now}
}

// Detect returns an Entry for every result whose timestamp is older than the
// configured max age. Results with a zero timestamp are always considered stale.
func (d *Detector) Detect(results []drift.Result) []Entry {
	var stale []Entry
	now := d.now()
	for _, r := range results {
		age := now.Sub(r.CheckedAt)
		if r.CheckedAt.IsZero() || age > d.maxAge {
			stale = append(stale, Entry{
				Service:  r.Service,
				LastSeen: r.CheckedAt,
				Age:      age,
			})
		}
	}
	return stale
}

// IsStale reports whether a single result is stale.
func (d *Detector) IsStale(r drift.Result) bool {
	if r.CheckedAt.IsZero() {
		return true
	}
	return d.now().Sub(r.CheckedAt) > d.maxAge
}
