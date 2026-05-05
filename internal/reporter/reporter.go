package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/driftwatch/internal/drift"
)

// Format controls the output format of the report.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Report holds the results of a drift detection run.
type Report struct {
	Timestamp time.Time     `json:"timestamp"`
	Results   []drift.Result `json:"results"`
	DriftCount int           `json:"drift_count"`
}

// Reporter writes drift reports to an output destination.
type Reporter struct {
	format Format
	out    io.Writer
}

// New creates a Reporter with the given format. If out is nil, os.Stdout is used.
func New(format Format, out io.Writer) *Reporter {
	if out == nil {
		out = os.Stdout
	}
	return &Reporter{format: format, out: out}
}

// Write renders the report to the configured output.
func (r *Reporter) Write(results []drift.Result) error {
	report := Report{
		Timestamp:  time.Now().UTC(),
		Results:    results,
		DriftCount: countDrifted(results),
	}

	switch r.format {
	case FormatJSON:
		return r.writeJSON(report)
	default:
		return r.writeText(report)
	}
}

func (r *Reporter) writeJSON(report Report) error {
	enc := json.NewEncoder(r.out)
	enc.SetIndent("", "  ")
	return enc.Encode(report)
}

func (r *Reporter) writeText(report Report) error {
	fmt.Fprintf(r.out, "Drift Report — %s\n", report.Timestamp.Format(time.RFC3339))
	fmt.Fprintf(r.out, "Services checked: %d | Drifted: %d\n\n", len(report.Results), report.DriftCount)

	for _, res := range report.Results {
		status := "OK"
		if res.Drifted {
			status = "DRIFTED"
		}
		fmt.Fprintf(r.out, "[%s] %s\n", status, res.Service)
		for _, d := range res.Diffs {
			fmt.Fprintf(r.out, "  - %s: expected=%q actual=%q\n", d.Field, d.Expected, d.Actual)
		}
	}
	return nil
}

func countDrifted(results []drift.Result) int {
	n := 0
	for _, r := range results {
		if r.Drifted {
			n++
		}
	}
	return n
}
