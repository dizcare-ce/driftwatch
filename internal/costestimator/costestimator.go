// Package costestimator estimates the operational cost of drift by assigning
// a weighted score to each drifted field based on configurable field weights.
// This allows operators to prioritise remediation of high-impact drift.
package costestimator

import (
	"fmt"
	"sort"
	"strings"

	"driftwatch/internal/drift"
)

// FieldWeight maps a field name (case-insensitive) to a cost multiplier.
type FieldWeight map[string]float64

// Estimate holds the cost breakdown for a single service.
type Estimate struct {
	Service    string
	TotalCost  float64
	Breakdown  []FieldCost
}

// FieldCost records the cost contribution of a single drifted field.
type FieldCost struct {
	Field  string
	Weight float64
	Cost   float64
}

// Estimator computes drift cost estimates.
type Estimator struct {
	weights     FieldWeight
	defaultCost float64
}

// New creates an Estimator with the supplied field weights and a default cost
// applied to any field not explicitly listed.
func New(weights FieldWeight, defaultCost float64) *Estimator {
	normalised := make(FieldWeight, len(weights))
	for k, v := range weights {
		normalised[strings.ToLower(k)] = v
	}
	return &Estimator{weights: normalised, defaultCost: defaultCost}
}

// Compute returns cost estimates for each result that contains drift.
// Results with no drift are omitted.
func (e *Estimator) Compute(results []drift.Result) []Estimate {
	var out []Estimate
	for _, r := range results {
		if len(r.Diffs) == 0 {
			continue
		}
		est := e.estimateOne(r)
		out = append(out, est)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].TotalCost != out[j].TotalCost {
			return out[i].TotalCost > out[j].TotalCost
		}
		return out[i].Service < out[j].Service
	})
	return out
}

func (e *Estimator) estimateOne(r drift.Result) Estimate {
	est := Estimate{Service: r.Service}
	for _, d := range r.Diffs {
		key := strings.ToLower(d.Field)
		w, ok := e.weights[key]
		if !ok {
			w = e.defaultCost
		}
		fc := FieldCost{Field: d.Field, Weight: w, Cost: w}
		est.Breakdown = append(est.Breakdown, fc)
		est.TotalCost += fc.Cost
	}
	return est
}

// Format returns a human-readable summary of all estimates.
func Format(estimates []Estimate) string {
	if len(estimates) == 0 {
		return "no drift cost estimated\n"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-30s %8s\n", "SERVICE", "COST"))
	sb.WriteString(strings.Repeat("-", 40) + "\n")
	for _, est := range estimates {
		sb.WriteString(fmt.Sprintf("%-30s %8.2f\n", est.Service, est.TotalCost))
	}
	return sb.String()
}
