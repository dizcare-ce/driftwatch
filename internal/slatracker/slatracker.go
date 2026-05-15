// Package slatracker measures whether services are meeting drift-free SLA
// targets over a rolling window of historical results.
package slatracker

import (
	"fmt"
	"sync"
	"time"

	"github.com/driftwatch/internal/drift"
)

// Entry records whether a single check run was compliant (no drift).
type Entry struct {
	Service   string
	Timestamp time.Time
	Compliant bool
}

// Report summarises SLA compliance for one service.
type Report struct {
	Service        string
	WindowStart    time.Time
	WindowEnd      time.Time
	TotalChecks    int
	CompliantChecks int
	CompliancePct  float64
	MeetsSLA       bool
}

// Tracker accumulates check results and evaluates SLA compliance.
type Tracker struct {
	mu        sync.Mutex
	window    time.Duration
	threshold float64 // 0–100
	entries   []Entry
}

// New creates a Tracker that evaluates compliance over the given rolling window
// and considers a service compliant when compliancePct >= threshold.
func New(window time.Duration, threshold float64) *Tracker {
	return &Tracker{
		window:    window,
		threshold: threshold,
	}
}

// Record adds a drift result to the tracker.
func (t *Tracker) Record(result drift.Result) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.entries = append(t.entries, Entry{
		Service:   result.Service,
		Timestamp: time.Now(),
		Compliant: !result.Drifted(),
	})
}

// Report returns an SLA report for the named service over the rolling window.
func (t *Tracker) Report(service string) (Report, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-t.window)

	var total, compliant int
	for _, e := range t.entries {
		if e.Service != service || e.Timestamp.Before(cutoff) {
			continue
		}
		total++
		if e.Compliant {
			compliant++
		}
	}

	if total == 0 {
		return Report{}, fmt.Errorf("slatracker: no data for service %q", service)
	}

	pct := float64(compliant) / float64(total) * 100
	return Report{
		Service:         service,
		WindowStart:     cutoff,
		WindowEnd:       now,
		TotalChecks:     total,
		CompliantChecks: compliant,
		CompliancePct:   pct,
		MeetsSLA:        pct >= t.threshold,
	}, nil
}

// Prune removes entries older than the rolling window to bound memory usage.
func (t *Tracker) Prune() {
	t.mu.Lock()
	defer t.mu.Unlock()
	cutoff := time.Now().Add(-t.window)
	kept := t.entries[:0]
	for _, e := range t.entries {
		if !e.Timestamp.Before(cutoff) {
			kept = append(kept, e)
		}
	}
	t.entries = kept
}
