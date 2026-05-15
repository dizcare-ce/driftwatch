// Package ownermap maps services to their declared owners (team, squad, or
// individual) and provides lookup helpers used by the reporter and notifier
// when attributing drift to a responsible party.
package ownermap

import (
	"fmt"
	"strings"
	"sync"
)

// Owner holds ownership metadata for a single service.
type Owner struct {
	Service string
	Team    string
	Contact string // email or Slack handle
}

// Map stores service → owner associations.
type Map struct {
	mu      sync.RWMutex
	entries map[string]Owner
}

// New returns an empty Map.
func New() *Map {
	return &Map{entries: make(map[string]Owner)}
}

// Set registers or replaces the owner for the given service name.
// Service names are normalised to lowercase before storage.
func (m *Map) Set(service, team, contact string) error {
	service = strings.ToLower(strings.TrimSpace(service))
	if service == "" {
		return fmt.Errorf("ownermap: service name must not be empty")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.entries[service] = Owner{Service: service, Team: team, Contact: contact}
	return nil
}

// Get returns the Owner for the given service and true when found.
// The lookup is case-insensitive.
func (m *Map) Get(service string) (Owner, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	o, ok := m.entries[strings.ToLower(service)]
	return o, ok
}

// Remove deletes the ownership record for the given service.
// It is a no-op when the service is not registered.
func (m *Map) Remove(service string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.entries, strings.ToLower(service))
}

// All returns a snapshot of every registered Owner, sorted by service name.
func (m *Map) All() []Owner {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]Owner, 0, len(m.entries))
	for _, o := range m.entries {
		out = append(out, o)
	}
	sortOwners(out)
	return out
}

func sortOwners(owners []Owner) {
	for i := 1; i < len(owners); i++ {
		for j := i; j > 0 && owners[j].Service < owners[j-1].Service; j-- {
			owners[j], owners[j-1] = owners[j-1], owners[j]
		}
	}
}
