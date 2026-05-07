// Package debounce provides a mechanism to suppress rapid repeated triggers
// for the same service, ensuring that only one action fires after a quiet period.
package debounce

import (
	"sync"
	"time"
)

// Debouncer delays execution of a callback until a service has been quiet
// for at least the configured wait duration.
type Debouncer struct {
	wait    time.Duration
	mu      sync.Mutex
	timers  map[string]*time.Timer
	callback func(service string)
}

// New creates a Debouncer that waits for the given duration of inactivity
// before invoking callback for a given service key.
func New(wait time.Duration, callback func(service string)) *Debouncer {
	return &Debouncer{
		wait:     wait,
		timers:   make(map[string]*time.Timer),
		callback: callback,
	}
}

// Trigger schedules callback(service) to fire after the wait duration.
// If Trigger is called again for the same service before the timer fires,
// the timer is reset, effectively debouncing rapid successive calls.
func (d *Debouncer) Trigger(service string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if t, ok := d.timers[service]; ok {
		t.Reset(d.wait)
		return
	}

	d.timers[service] = time.AfterFunc(d.wait, func() {
		d.mu.Lock()
		delete(d.timers, service)
		d.mu.Unlock()
		d.callback(service)
	})
}

// Cancel cancels any pending timer for the given service.
// It is a no-op if no timer is pending.
func (d *Debouncer) Cancel(service string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if t, ok := d.timers[service]; ok {
		t.Stop()
		delete(d.timers, service)
	}
}

// Pending returns true if a timer is currently active for the given service.
func (d *Debouncer) Pending(service string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	_, ok := d.timers[service]
	return ok
}
