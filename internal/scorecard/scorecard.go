// Package scorecard aggregates per-service drift metrics into a
// human-readable health score between 0 (perfect) and 100 (fully drifted).
package scorecard

import (
	"fmt"
	"sort"

	"driftwatch/internal/drift"
)

// Entry holds the computed score for a single service.
type Entry struct {
	Service string
	Score   int    // 0–100
	Grade   string // A, B, C, D, F
	Drifted int
	Total   int
}

// Build computes a scorecard from a slice of drift results.
// Results with no fields are skipped.
func Build(results []drift.Result) []Entry {
	entries := make([]Entry, 0, len(results))
	for _, r := range results {
		if r.Total == 0 {
			continue
		}
		score := computeScore(r.Drifted, r.Total)
		entries = append(entries, Entry{
			Service: r.Service,
			Score:   score,
			Grade:   grade(score),
			Drifted: r.Drifted,
			Total:   r.Total,
		})
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Score != entries[j].Score {
			return entries[i].Score > entries[j].Score // worst first
		}
		return entries[i].Service < entries[j].Service
	})
	return entries
}

// computeScore returns 0 for no drift and scales linearly to 100.
func computeScore(drifted, total int) int {
	if total == 0 {
		return 0
	}
	return (drifted * 100) / total
}

func grade(score int) string {
	switch {
	case score == 0:
		return "A"
	case score <= 25:
		return "B"
	case score <= 50:
		return "C"
	case score <= 75:
		return "D"
	default:
		return "F"
	}
}

// Summary returns a one-line overview of the scorecard.
func Summary(entries []Entry) string {
	if len(entries) == 0 {
		return "scorecard: no services evaluated"
	}
	failing := 0
	for _, e := range entries {
		if e.Grade == "F" {
			failing++
		}
	}
	return fmt.Sprintf("scorecard: %d service(s) evaluated, %d failing", len(entries), failing)
}
