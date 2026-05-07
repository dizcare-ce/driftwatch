// Package ratelimit provides a simple token-bucket rate limiter for
// controlling how frequently drift notifications are emitted per service.
package ratelimit

import (
	"sync"
	"time"
)

// Limiter tracks per-service notification rate limits.
type Limiter struct {
	mu       sync.Mutex
	buckets  map[string]time.Time
	cooldown time.Duration
}

// New returns a Limiter that enforces the given cooldown between
// successive notifications for the same service.
func New(cooldown time.Duration) *Limiter {
	return &Limiter{
		buckets:  make(map[string]time.Time),
		cooldown: cooldown,
	}
}

// Allow reports whether a notification for the named service is
// permitted at the current time. If allowed, the service's last-seen
// timestamp is updated.
func (l *Limiter) Allow(service string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	last, seen := l.buckets[service]
	if seen && now.Sub(last) < l.cooldown {
		return false
	}
	l.buckets[service] = now
	return true
}

// Reset clears the recorded timestamp for a service, allowing the next
// notification to pass through immediately.
func (l *Limiter) Reset(service string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.buckets, service)
}

// ResetAll clears all recorded timestamps.
func (l *Limiter) ResetAll() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.buckets = make(map[string]time.Time)
}
