// Package dependency tracks inter-service dependency relationships and
// surfaces drift results that may be caused by upstream service changes.
package dependency

import (
	"fmt"
	"sync"
)

// Graph holds directed dependency edges between services.
// An edge A → B means service A depends on service B.
type Graph struct {
	mu    sync.RWMutex
	edges map[string][]string // dependant → dependencies
}

// New returns an empty dependency Graph.
func New() *Graph {
	return &Graph{
		edges: make(map[string][]string),
	}
}

// Add registers that service `from` depends on service `to`.
// Duplicate edges are silently ignored.
func (g *Graph) Add(from, to string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	for _, existing := range g.edges[from] {
		if existing == to {
			return
		}
	}
	g.edges[from] = append(g.edges[from], to)
}

// Remove deletes all dependency edges originating from the given service.
func (g *Graph) Remove(service string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	delete(g.edges, service)
}

// Dependants returns all services that declare a dependency on `service`.
func (g *Graph) Dependants(service string) []string {
	g.mu.RLock()
	defer g.mu.RUnlock()

	var out []string
	for from, deps := range g.edges {
		for _, d := range deps {
			if d == service {
				out = append(out, from)
				break
			}
		}
	}
	return out
}

// Dependencies returns the direct dependencies of `service`.
func (g *Graph) Dependencies(service string) []string {
	g.mu.RLock()
	defer g.mu.RUnlock()

	copy := make([]string, len(g.edges[service]))
	for i, v := range g.edges[service] {
		copy[i] = v
	}
	return copy
}

// Validate checks for obvious misconfigurations such as self-referential
// edges and returns a combined error when any are found.
func (g *Graph) Validate() error {
	g.mu.RLock()
	defer g.mu.RUnlock()

	for from, deps := range g.edges {
		for _, to := range deps {
			if from == to {
				return fmt.Errorf("dependency: self-referential edge for service %q", from)
			}
		}
	}
	return nil
}
