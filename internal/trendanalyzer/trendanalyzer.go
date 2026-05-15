// Package trendanalyzer computes drift trend statistics over a rolling
// window of historical run results, indicating whether drift is increasing,
// decreasing, or stable across observed services.
package trendanalyzer

import (
	"time"

	"driftwatch/internal/drift"
)

// Direction describes the overall direction of drift over the window.
type Direction string

const (
	DirectionIncreasing Direction = "increasing"
	DirectionDecreasing Direction = "decreasing"
	DirectionStable     Direction = "stable"
)

// Report holds the trend analysis output for a window of runs.
type Report struct {
	Window    int       // number of runs analysed
	First     time.Time // timestamp of the earliest run
	Last      time.Time // timestamp of the most recent run
	MinDrifts int
	MaxDrifts int
	AvgDrifts float64
	Direction Direction
}

// Snapshot is a lightweight record of a single run's drift count and time.
type Snapshot struct {
	At    time.Time
	Drift int
}

// Analyzer computes trend reports from ordered run snapshots.
type Analyzer struct {
	maxWindow int
}

// New returns an Analyzer that considers at most maxWindow recent snapshots.
func New(maxWindow int) *Analyzer {
	if maxWindow < 1 {
		maxWindow = 1
	}
	return &Analyzer{maxWindow: maxWindow}
}

// FromResults converts a slice of drift results into a Snapshot.
func FromResults(at time.Time, results []drift.Result) Snapshot {
	count := 0
	for _, r := range results {
		count += len(r.Diffs)
	}
	return Snapshot{At: at, Drift: count}
}

// Analyse computes a trend Report over the provided snapshots.
// Snapshots must be ordered oldest-first. At most maxWindow entries are used.
func (a *Analyzer) Analyse(snapshots []Snapshot) Report {
	if len(snapshots) == 0 {
		return Report{Direction: DirectionStable}
	}
	if len(snapshots) > a.maxWindow {
		snapshots = snapshots[len(snapshots)-a.maxWindow:]
	}

	min, max, sum := snapshots[0].Drift, snapshots[0].Drift, 0
	for _, s := range snapshots {
		sum += s.Drift
		if s.Drift < min {
			min = s.Drift
		}
		if s.Drift > max {
			max = s.Drift
		}
	}

	return Report{
		Window:    len(snapshots),
		First:     snapshots[0].At,
		Last:      snapshots[len(snapshots)-1].At,
		MinDrifts: min,
		MaxDrifts: max,
		AvgDrifts: float64(sum) / float64(len(snapshots)),
		Direction: direction(snapshots),
	}
}

// direction determines trend direction by comparing the first and second halves.
func direction(snapshots []Snapshot) Direction {
	if len(snapshots) < 2 {
		return DirectionStable
	}
	mid := len(snapshots) / 2
	early := avg(snapshots[:mid])
	late := avg(snapshots[mid:])
	switch {
	case late > early:
		return DirectionIncreasing
	case late < early:
		return DirectionDecreasing
	default:
		return DirectionStable
	}
}

func avg(ss []Snapshot) float64 {
	if len(ss) == 0 {
		return 0
	}
	sum := 0
	for _, s := range ss {
		sum += s.Drift
	}
	return float64(sum) / float64(len(ss))
}
