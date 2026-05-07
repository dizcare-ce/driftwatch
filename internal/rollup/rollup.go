// Package rollup aggregates drift results across multiple services into
// a single summary report, grouping findings by severity and service name.
package rollup

import (
	"fmt"
	"sort"
	"strings"

	"github.com/driftwatch/internal/drift"
)

// Summary holds aggregated drift information across all checked services.
type Summary struct {
	Total    int
	Drifted  int
	Clean    int
	Services []ServiceSummary
}

// ServiceSummary holds the drift result for a single service.
type ServiceSummary struct {
	Name    string
	Drifted bool
	Diffs   int
	Details []string
}

// Aggregate combines a slice of drift results into a single Summary.
func Aggregate(results []drift.Result) Summary {
	s := Summary{
		Total: len(results),
	}

	sorted := make([]drift.Result, len(results))
	copy(sorted, results)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Service < sorted[j].Service
	})

	for _, r := range sorted {
		ss := ServiceSummary{
			Name:    r.Service,
			Drifted: r.Drifted(),
			Diffs:   len(r.Diffs),
		}
		for _, d := range r.Diffs {
			ss.Details = append(ss.Details, fmt.Sprintf("%s: %v → %v", d.Field, d.Got, d.Want))
		}
		if ss.Drifted {
			s.Drifted++
		} else {
			s.Clean++
		}
		s.Services = append(s.Services, ss)
	}
	return s
}

// String returns a human-readable representation of the summary.
func (s Summary) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Drift summary: %d total, %d drifted, %d clean\n", s.Total, s.Drifted, s.Clean)
	for _, svc := range s.Services {
		status := "clean"
		if svc.Drifted {
			status = "DRIFTED"
		}
		fmt.Fprintf(&sb, "  [%s] %s (%d diff(s))\n", status, svc.Name, svc.Diffs)
		for _, detail := range svc.Details {
			fmt.Fprintf(&sb, "      %s\n", detail)
		}
	}
	return sb.String()
}
