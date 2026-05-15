// Package anomaly detects statistical anomalies in drift counts across
// successive runs, flagging services whose drift frequency deviates
// significantly from their recent baseline.
package anomaly

import (
	"fmt"
	"math"
	"sort"

	"driftwatch/internal/drift"
)

// Anomaly describes a service whose drift count is unusually high relative
// to its recent history.
type Anomaly struct {
	Service  string
	Current  int
	Mean     float64
	StdDev   float64
	ZScore   float64
}

// String returns a human-readable summary of the anomaly.
func (a Anomaly) String() string {
	return fmt.Sprintf("%s: current=%d mean=%.2f stddev=%.2f z=%.2f",
		a.Service, a.Current, a.Mean, a.StdDev, a.ZScore)
}

// Detector identifies services with anomalous drift counts.
type Detector struct {
	// Threshold is the minimum absolute z-score required to flag an anomaly.
	// Defaults to 2.0 (two standard deviations).
	Threshold float64
	// MinSamples is the minimum number of historical snapshots required
	// before anomaly detection is applied for a service.
	MinSamples int
}

// New returns a Detector with sensible defaults.
func New() *Detector {
	return &Detector{Threshold: 2.0, MinSamples: 3}
}

// Detect analyses the current batch of results against historical snapshots.
// history is a slice of past run results, ordered oldest-first.
// current is the most recent run.
// Returns a slice of Anomaly values, sorted by descending z-score magnitude.
func (d *Detector) Detect(history [][]drift.Result, current []drift.Result) []Anomaly {
	if len(history) < d.MinSamples {
		return nil
	}

	// Build per-service historical drift counts.
	counts := make(map[string][]float64)
	for _, run := range history {
		for _, r := range run {
			counts[r.Service] = append(counts[r.Service], float64(len(r.Diffs)))
		}
	}

	var anomalies []Anomaly
	for _, r := range current {
		hist, ok := counts[r.Service]
		if !ok || len(hist) < d.MinSamples {
			continue
		}
		mean, stddev := stats(hist)
		if stddev == 0 {
			continue
		}
		z := (float64(len(r.Diffs)) - mean) / stddev
		if math.Abs(z) >= d.Threshold {
			anomalies = append(anomalies, Anomaly{
				Service: r.Service,
				Current: len(r.Diffs),
				Mean:    mean,
				StdDev:  stddev,
				ZScore:  z,
			})
		}
	}

	sort.Slice(anomalies, func(i, j int) bool {
		return math.Abs(anomalies[i].ZScore) > math.Abs(anomalies[j].ZScore)
	})
	return anomalies
}

// stats returns the mean and population standard deviation of xs.
func stats(xs []float64) (mean, stddev float64) {
	for _, v := range xs {
		mean += v
	}
	mean /= float64(len(xs))
	var variance float64
	for _, v := range xs {
		d := v - mean
		variance += d * d
	}
	stddev = math.Sqrt(variance / float64(len(xs)))
	return
}
