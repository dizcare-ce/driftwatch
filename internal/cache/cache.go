// Package cache provides a simple in-memory result cache for drift detection
// runs, avoiding redundant comparisons when service definitions have not changed.
package cache

import (
	"sync"
	"time"

	"github.com/driftwatch/internal/drift"
)

// Entry holds a cached drift result alongside metadata.
type Entry struct {
	Result    drift.Result
	CachedAt  time.Time
	Checksum  string
}

// Cache stores the most recent drift result per service name.
type Cache struct {
	mu      sync.RWMutex
	entries map[string]Entry
	ttl     time.Duration
}

// New returns a Cache that expires entries after ttl.
func New(ttl time.Duration) *Cache {
	return &Cache{
		entries: make(map[string]Entry),
		ttl:     ttl,
	}
}

// Set stores a result for the given service, keyed by checksum.
func (c *Cache) Set(service, checksum string, result drift.Result) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[service] = Entry{
		Result:   result,
		CachedAt: time.Now(),
		Checksum: checksum,
	}
}

// Get retrieves a cached result. It returns (entry, true) when a valid,
// non-expired entry exists whose checksum matches; otherwise (zero, false).
func (c *Cache) Get(service, checksum string) (Entry, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	e, ok := c.entries[service]
	if !ok {
		return Entry{}, false
	}
	if e.Checksum != checksum {
		return Entry{}, false
	}
	if c.ttl > 0 && time.Since(e.CachedAt) > c.ttl {
		return Entry{}, false
	}
	return e, true
}

// Invalidate removes the cached entry for a service.
func (c *Cache) Invalidate(service string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, service)
}

// Len returns the number of entries currently held in the cache.
func (c *Cache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.entries)
}
