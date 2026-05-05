// Package runner wires together the core driftwatch pipeline: loading
// service definitions, detecting drift, reporting results, and optionally
// notifying on significant findings.
package runner

import (
	"context"
	"fmt"
	"io"

	"driftwatch/internal/drift"
	"driftwatch/internal/notifier"
	"driftwatch/internal/reporter"
	"driftwatch/internal/source"
)

// Runner executes a single drift-check cycle.
type Runner struct {
	loader   *source.Loader
	detector *drift.Detector
	reporter *reporter.Reporter
	notifier *notifier.Notifier
}

// New constructs a Runner from its collaborators.
func New(
	l *source.Loader,
	d *drift.Detector,
	r *reporter.Reporter,
	n *notifier.Notifier,
) *Runner {
	return &Runner{
		loader:   l,
		detector: d,
		reporter: r,
		notifier: n,
	}
}

// Run loads all service definitions, compares each against its live state,
// writes a report, and fires notifications. It returns the first error
// encountered, but always attempts to report whatever results were gathered.
func (r *Runner) Run(ctx context.Context, out io.Writer) error {
	defs, err := r.loader.LoadAll(ctx)
	if err != nil {
		return fmt.Errorf("runner: load definitions: %w", err)
	}

	var results []drift.Result
	for _, def := range defs {
		res, err := r.detector.Compare(ctx, def)
		if err != nil {
			return fmt.Errorf("runner: compare %q: %w", def.Name, err)
		}
		results = append(results, res)
	}

	if err := r.reporter.Write(out, results); err != nil {
		return fmt.Errorf("runner: write report: %w", err)
	}

	if err := r.notifier.Notify(ctx, results); err != nil {
		return fmt.Errorf("runner: notify: %w", err)
	}

	return nil
}
