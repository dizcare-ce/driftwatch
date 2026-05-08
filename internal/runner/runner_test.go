package runner_test

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"driftwatch/internal/drift"
	"driftwatch/internal/runner"
)

// stubLoader satisfies the interface used by Runner via a thin wrapper.
type stubLoader struct {
	defs []drift.Definition
	err  error
}

func (s *stubLoader) LoadAll(_ context.Context) ([]drift.Definition, error) {
	return s.defs, s.err
}

type stubDetector struct {
	result drift.Result
	err    error
}

func (s *stubDetector) Compare(_ context.Context, def drift.Definition) (drift.Result, error) {
	s.result.Definition = def
	return s.result, s.err
}

type stubReporter struct{ err error }

func (s *stubReporter) Write(_ interface{}, _ []drift.Result) error { return s.err }

type stubNotifier struct{ err error }

func (s *stubNotifier) Notify(_ context.Context, _ []drift.Result) error { return s.err }

func TestRun_NoDrift_ReturnsNil(t *testing.T) {
	def := drift.Definition{Name: "api"}
	l := &stubLoader{defs: []drift.Definition{def}}
	d := &stubDetector{result: drift.Result{}}
	rep := &stubReporter{}
	not := &stubNotifier{}

	r := runner.New(l, d, rep, not)
	if err := r.Run(context.Background(), &bytes.Buffer{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRun_LoaderError_Propagates(t *testing.T) {
	l := &stubLoader{err: errors.New("disk read failure")}
	r := runner.New(l, &stubDetector{}, &stubReporter{}, &stubNotifier{})

	err := r.Run(context.Background(), &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestRun_DetectorError_Propagates(t *testing.T) {
	l := &stubLoader{defs: []drift.Definition{{Name: "svc"}}}
	d := &stubDetector{err: errors.New("compare failed")}
	r := runner.New(l, d, &stubReporter{}, &stubNotifier{})

	err := r.Run(context.Background(), &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestRun_ReporterError_Propagates(t *testing.T) {
	l := &stubLoader{defs: []drift.Definition{{Name: "svc"}}}
	d := &stubDetector{}
	rep := &stubReporter{err: errors.New("write failed")}
	r := runner.New(l, d, rep, &stubNotifier{})

	err := r.Run(context.Background(), &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestRun_NotifierError_Propagates(t *testing.T) {
	l := &stubLoader{defs: []drift.Definition{{Name: "svc"}}}
	d := &stubDetector{}
	rep := &stubReporter{}
	not := &stubNotifier{err: errors.New("notify failed")}
	r := runner.New(l, d, rep, not)

	err := r.Run(context.Background(), &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
