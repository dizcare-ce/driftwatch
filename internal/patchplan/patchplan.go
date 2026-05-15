// Package patchplan generates a human-readable remediation plan from drift
// results, suggesting the minimal set of changes needed to bring each service
// back in line with its source definition.
package patchplan

import (
	"fmt"
	"sort"
	"strings"

	"driftwatch/internal/drift"
)

// Action describes a single remediation step for one field.
type Action struct {
	// Field is the dotted path of the drifted key.
	Field string
	// Expected is the value defined in the source.
	Expected any
	// Actual is the value observed in the deployed service.
	Actual any
	// Description is a short human-readable instruction.
	Description string
}

// Plan holds all remediation actions for a single service.
type Plan struct {
	Service string
	Actions []Action
}

// Build converts a slice of drift results into a slice of Plans.
// Results with no diffs produce no Plan entry.
func Build(results []drift.Result) []Plan {
	var plans []Plan
	for _, r := range results {
		if len(r.Diffs) == 0 {
			continue
		}
		p := Plan{Service: r.Service}
		for _, d := range r.Diffs {
			p.Actions = append(p.Actions, Action{
				Field:       d.Field,
				Expected:    d.Expected,
				Actual:      d.Actual,
				Description: describe(d),
			})
		}
		sort.Slice(p.Actions, func(i, j int) bool {
			return p.Actions[i].Field < p.Actions[j].Field
		})
		plans = append(plans, p)
	}
	sort.Slice(plans, func(i, j int) bool {
		return plans[i].Service < plans[j].Service
	})
	return plans
}

// describe returns a short instruction for a single diff.
func describe(d drift.Diff) string {
	field := d.Field
	switch {
	case d.Expected == nil:
		return fmt.Sprintf("remove field %q (unexpected value: %v)", field, d.Actual)
	case d.Actual == nil:
		return fmt.Sprintf("set field %q to %v (currently absent)", field, d.Expected)
	default:
		return fmt.Sprintf("update field %q from %v to %v", field, d.Actual, d.Expected)
	}
}

// Format renders all plans as a multi-line string suitable for display.
func Format(plans []Plan) string {
	if len(plans) == 0 {
		return "No remediation required — all services are in sync."
	}
	var sb strings.Builder
	for _, p := range plans {
		fmt.Fprintf(&sb, "Service: %s\n", p.Service)
		for _, a := range p.Actions {
			fmt.Fprintf(&sb, "  • %s\n", a.Description)
		}
	}
	return strings.TrimRight(sb.String(), "\n")
}
