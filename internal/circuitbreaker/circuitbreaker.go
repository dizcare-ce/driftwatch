// Package circuitbreaker provides a simple circuit breaker to prevent
// repeated execution of failing operations. Once the failure threshold is
// reached the circuit opens and all subsequent calls are rejected until the
// reset timeout elapses.
package circuitbreaker

import (
	"errors"
	"sync"
	"time"
)

// ErrOpen is returned when the circuit is open and the call is rejected.
var ErrOpen = errors.New("circuit breaker is open")

// State represents the current state of the circuit breaker.
type State int

const (
	StateClosed State = iota
	StateOpen
)

// Breaker is a circuit breaker that tracks consecutive failures for a named
// service and opens the circuit after a configurable threshold.
type Breaker struct {
	mu        sync.Mutex
	failures  map[string]int
	openUntil map[string]time.Time

	threshold int
	timeout   time.Duration
}

// New creates a Breaker that opens after threshold consecutive failures and
// attempts to reset after timeout.
func New(threshold int, timeout time.Duration) *Breaker {
	return &Breaker{
		failures:  make(map[string]int),
		openUntil: make(map[string]time.Time),
		threshold: threshold,
		timeout:   timeout,
	}
}

// Allow reports whether a call for the named service should be allowed.
// If the circuit is open and the reset timeout has elapsed the circuit is
// moved back to half-open (closed with zero failures) automatically.
func (b *Breaker) Allow(service string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	if until, ok := b.openUntil[service]; ok {
		if time.Now().Before(until) {
			return false
		}
		// reset to half-open
		delete(b.openUntil, service)
		b.failures[service] = 0
	}
	return true
}

// RecordSuccess resets the failure counter for the named service.
func (b *Breaker) RecordSuccess(service string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures[service] = 0
	delete(b.openUntil, service)
}

// RecordFailure increments the failure counter for the named service and opens
// the circuit if the threshold is reached.
func (b *Breaker) RecordFailure(service string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures[service]++
	if b.failures[service] >= b.threshold {
		b.openUntil[service] = time.Now().Add(b.timeout)
	}
}

// StateOf returns the current State for the named service.
func (b *Breaker) StateOf(service string) State {
	b.mu.Lock()
	defer b.mu.Unlock()
	if until, ok := b.openUntil[service]; ok && time.Now().Before(until) {
		return StateOpen
	}
	return StateClosed
}
