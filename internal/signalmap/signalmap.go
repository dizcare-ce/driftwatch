// Package signalmap maps drift results to named severity signals that can be
// consumed by external alerting pipelines. Each signal carries a name, level,
// and the service that triggered it.
package signalmap

import (
	"fmt"
	"sync"

	"driftwatch/internal/drift"
	"driftwatch/internal/priority"
)

// Signal represents a named alert signal derived from a drift result.
type Signal struct {
	Service string
	Name    string
	Level   priority.Level
	Message string
}

// Mapper converts drift results into signals using a configurable name template.
type Mapper struct {
	mu      sync.Mutex
	prefix  string
	scorer  *priority.Scorer
}

// New returns a Mapper that prefixes every signal name with prefix.
func New(prefix string, scorer *priority.Scorer) *Mapper {
	return &Mapper{prefix: prefix, scorer: scorer}
}

// Map converts a slice of drift results into signals, skipping clean results
// whose level falls below minLevel.
func (m *Mapper) Map(results []drift.Result, minLevel priority.Level) []Signal {
	m.mu.Lock()
	defer m.mu.Unlock()

	var signals []Signal
	for _, r := range results {
		lvl := m.scorer.Score(r)
		if lvl < minLevel {
			continue
		}
		sig := Signal{
			Service: r.Service,
			Name:    fmt.Sprintf("%s.%s", m.prefix, r.Service),
			Level:   lvl,
			Message: r.Summary(),
		}
		signals = append(signals, sig)
	}
	return signals
}
