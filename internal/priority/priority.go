// Package priority assigns and compares severity levels to drift results,
// allowing operators to triage which services need immediate attention.
package priority

import (
	"fmt"
	"strings"

	"github.com/driftwatch/internal/drift"
)

// Level represents the urgency of a drift result.
type Level int

const (
	LevelLow Level = iota
	LevelMedium
	LevelHigh
	LevelCritical
)

func (l Level) String() string {
	switch l {
	case LevelLow:
		return "low"
	case LevelMedium:
		return "medium"
	case LevelHigh:
		return "high"
	case LevelCritical:
		return "critical"
	default:
		return "unknown"
	}
}

// ParseLevel converts a string to a Level.
func ParseLevel(s string) (Level, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "low":
		return LevelLow, nil
	case "medium":
		return LevelMedium, nil
	case "high":
		return LevelHigh, nil
	case "critical":
		return LevelCritical, nil
	default:
		return LevelLow, fmt.Errorf("priority: unknown level %q", s)
	}
}

// Scorer assigns a priority Level to a drift result based on the number
// and nature of its diffs.
type Scorer struct {
	mediumThreshold  int
	highThreshold    int
	criticalThreshold int
}

// New returns a Scorer with sensible defaults.
func New() *Scorer {
	return &Scorer{
		mediumThreshold:  1,
		highThreshold:    3,
		criticalThreshold: 6,
	}
}

// Score returns the Level for the given result.
func (s *Scorer) Score(r drift.Result) Level {
	n := len(r.Diffs)
	switch {
	case n == 0:
		return LevelLow
	case n < s.highThreshold:
		return LevelMedium
	case n < s.criticalThreshold:
		return LevelHigh
	default:
		return LevelCritical
	}
}

// ScoreAll returns a map of service name to Level for a slice of results.
func (s *Scorer) ScoreAll(results []drift.Result) map[string]Level {
	out := make(map[string]Level, len(results))
	for _, r := range results {
		out[r.Service] = s.Score(r)
	}
	return out
}
