// Package runner orchestrates a single drift-check cycle: load sources,
// compare against live state, report results, and record history.
package runner

import (
	"context"
	"fmt"
	"log/slog"

	"driftwatch/internal/drift"
	"driftwatch/internal/history"
	"driftwatch/internal/metrics"
	"driftwatch/internal/reporter"
	"driftwatch/internal/retry"
	"driftwatch/internal/source"
)

// Runner executes one full drift-detection cycle.
type Runner struct {
	loader   *source.Loader
	detector *drift.Detector
	reporter *reporter.Reporter
	history  *history.History
	metrics  *metrics.Metrics
	policy   retry.Policy
	log      *slog.Logger
}

// New constructs a Runner with the supplied dependencies.
func New(
	l *source.Loader,
	d *drift.Detector,
	r *reporter.Reporter,
	h *history.History,
	m *metrics.Metrics,
	log *slog.Logger,
) *Runner {
	return &Runner{
		loader:   l,
		detector: d,
		reporter: r,
		history:  h,
		metrics:  m,
		policy:   retry.DefaultPolicy(),
		log:      log,
	}
}

// Run performs a single drift-check cycle, retrying transient load errors
// according to the configured retry policy.
func (r *Runner) Run(ctx context.Context) error {
	var definitions []source.Definition

	loadErr := r.policy.Do(ctx, func() error {
		defs, err := r.loader.LoadAll(ctx)
		if err != nil {
			r.log.Warn("source load failed, will retry", "err", err)
			return err
		}
		definitions = defs
		return nil
	})
	if loadErr != nil {
		r.metrics.RecordRun(nil, loadErr)
		return fmt.Errorf("load sources: %w", loadErr)
	}

	results := r.detector.Compare(definitions)

	if err := r.reporter.Write(ctx, results); err != nil {
		r.log.Warn("reporter write failed", "err", err)
	}

	if err := r.history.Record(results); err != nil {
		r.log.Warn("history record failed", "err", err)
	}

	r.metrics.RecordRun(results, nil)
	return nil
}
