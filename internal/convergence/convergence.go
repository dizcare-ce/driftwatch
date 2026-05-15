// Package convergence estimates how long a drifted service has been
// out of sync and projects when it might return to a clean state based
// on historical drift records.
package convergence

import (
	"fmt"
	"math"
	"time"

	"driftwatch/internal/drift"
)

// Estimate holds the convergence analysis for a single service.
type Estimate struct {
	Service       string
	DriftingSince time.Time
	DriftDuration time.Duration
	// ProjectedClean is zero when convergence cannot be projected.
	ProjectedClean time.Time
	Trend          string // "improving", "worsening", "stable"
}

// String returns a human-readable summary of the estimate.
func (e Estimate) String() string {
	if e.ProjectedClean.IsZero() {
		return fmt.Sprintf("%s: drifting for %s (%s, no projection)",
			e.Service, fmtDuration(e.DriftDuration), e.Trend)
	}
	return fmt.Sprintf("%s: drifting for %s (%s, clean ~%s)",
		e.Service, fmtDuration(e.DriftDuration), e.Trend,
		e.ProjectedClean.Format(time.RFC3339))
}

// Analyse computes convergence estimates from an ordered slice of
// historical result sets (oldest first). now is the reference time.
func Analyse(history [][]drift.Result, now time.Time) []Estimate {
	if len(history) == 0 {
		return nil
	}

	// Build per-service diff-count series.
	type series struct {
		firstSeen time.Time
		counts    []int
	}
	data := map[string]*series{}

	for _, snapshot := range history {
		for _, r := range snapshot {
			if _, ok := data[r.Service]; !ok {
				data[r.Service] = &series{firstSeen: now}
			}
			if len(r.Diffs) > 0 {
				data[r.Service].counts = append(data[r.Service].counts, len(r.Diffs))
				if r.CheckedAt.Before(data[r.Service].firstSeen) {
					data[r.Service].firstSeen = r.CheckedAt
				}
			}
		}
	}

	var estimates []Estimate
	for svc, s := range data {
		if len(s.counts) == 0 {
			continue
		}
		trendStr, slope := trendOf(s.counts)
		driftDur := now.Sub(s.firstSeen)
		var projected time.Time
		if slope < 0 && math.Abs(slope) > 0.1 {
			last := float64(s.counts[len(s.counts)-1])
			stepsToZero := last / math.Abs(slope)
			projected = now.Add(time.Duration(stepsToZero) * time.Hour)
		}
		estimates = append(estimates, Estimate{
			Service:        svc,
			DriftingSince:  s.firstSeen,
			DriftDuration:  driftDur,
			ProjectedClean: projected,
			Trend:          trendStr,
		})
	}
	return estimates
}

func trendOf(counts []int) (string, float64) {
	if len(counts) < 2 {
		return "stable", 0
	}
	first := float64(counts[0])
	last := float64(counts[len(counts)-1])
	slope := (last - first) / float64(len(counts)-1)
	switch {
	case slope < -0.1:
		return "improving", slope
	case slope > 0.1:
		return "worsening", slope
	default:
		return "stable", slope
	}
}

func fmtDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm", int(d.Minutes()))
	}
	return fmt.Sprintf("%dh%dm", int(d.Hours()), int(d.Minutes())%60)
}
