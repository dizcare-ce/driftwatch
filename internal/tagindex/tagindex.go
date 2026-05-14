// Package tagindex maintains an in-memory index of service definitions
// keyed by their labels, enabling fast lookup by tag or label selector.
package tagindex

import (
	"sync"
)

// Entry holds the service name and its associated labels.
type Entry struct {
	Service string
	Labels  map[string]string
}

// Index maps label key-value pairs to the set of services that carry them.
type Index struct {
	mu      sync.RWMutex
	byLabel map[string]map[string][]string // key -> value -> []service
	entries map[string]Entry
}

// New returns an empty Index.
func New() *Index {
	return &Index{
		byLabel: make(map[string]map[string][]string),
		entries: make(map[string]Entry),
	}
}

// Add registers a service with the given labels in the index.
// Calling Add for an existing service replaces its previous labels.
func (idx *Index) Add(service string, labels map[string]string) {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	// Remove stale entries for this service.
	if prev, ok := idx.entries[service]; ok {
		idx.removeFromLabel(service, prev.Labels)
	}

	copy := make(map[string]string, len(labels))
	for k, v := range labels {
		copy[k] = v
	}
	idx.entries[service] = Entry{Service: service, Labels: copy}

	for k, v := range copy {
		if idx.byLabel[k] == nil {
			idx.byLabel[k] = make(map[string][]string)
		}
		idx.byLabel[k][v] = append(idx.byLabel[k][v], service)
	}
}

// Remove deletes a service from the index.
func (idx *Index) Remove(service string) {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	if prev, ok := idx.entries[service]; ok {
		idx.removeFromLabel(service, prev.Labels)
		delete(idx.entries, service)
	}
}

// Lookup returns all service names that have the given label key and value.
func (idx *Index) Lookup(key, value string) []string {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	if vals, ok := idx.byLabel[key]; ok {
		if svcs, ok := vals[value]; ok {
			out := make([]string, len(svcs))
			copy(out, svcs)
			return out
		}
	}
	return nil
}

// All returns a snapshot of every registered entry.
func (idx *Index) All() []Entry {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	out := make([]Entry, 0, len(idx.entries))
	for _, e := range idx.entries {
		out = append(out, e)
	}
	return out
}

// removeFromLabel removes a service name from every inverted-index bucket
// associated with the provided labels. Must be called with idx.mu held.
func (idx *Index) removeFromLabel(service string, labels map[string]string) {
	for k, v := range labels {
		svcs := idx.byLabel[k][v]
		filtered := svcs[:0]
		for _, s := range svcs {
			if s != service {
				filtered = append(filtered, s)
			}
		}
		if len(filtered) == 0 {
			delete(idx.byLabel[k], v)
		} else {
			idx.byLabel[k][v] = filtered
		}
	}
}
