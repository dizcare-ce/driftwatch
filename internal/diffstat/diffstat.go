// Package diffstat provides aggregated statistics over a collection of drift
// results, giving operators a quick numerical summary of the current state of
// the fleet.
package diffstat

import "github.com/driftwatch/driftwatch/internal/drift"

// Stats holds aggregated counts derived from a slice of drift results.
type Stats struct {
	Total    int
	Clean    int
	Drifted  int
	Diffs    int
	AvgDiffs float64
}

// Compute walks results and returns a populated Stats value.
func Compute(results []drift.Result) Stats {
	s := Stats{
		Total: len(results),
	}

	for _, r := range results {
		nd := len(r.Diffs)
		if nd == 0 {
			s.Clean++
		} else {
			s.Drifted++
			s.Diffs += nd
		}
	}

	if s.Drifted > 0 {
		s.AvgDiffs = float64(s.Diffs) / float64(s.Drifted)
	}

	return s
}

// DriftRate returns the fraction of services that are drifted, in the range
// [0.0, 1.0]. Returns 0 when there are no results.
func DriftRate(s Stats) float64 {
	if s.Total == 0 {
		return 0
	}
	return float64(s.Drifted) / float64(s.Total)
}
