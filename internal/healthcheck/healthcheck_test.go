package healthcheck_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/driftwatch/driftwatch/internal/healthcheck"
	"github.com/driftwatch/driftwatch/internal/metrics"
)

func newMetrics(t *testing.T) *metrics.Metrics {
	t.Helper()
	m, err := metrics.New()
	if err != nil {
		t.Fatalf("metrics.New: %v", err)
	}
	return m
}

func TestHandler_NoRunYet(t *testing.T) {
	m := newMetrics(t)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)

	healthcheck.Handler(m)(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var s healthcheck.Status
	if err := json.NewDecoder(rec.Body).Decode(&s); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if s.Message != "no run completed yet" {
		t.Errorf("unexpected message: %q", s.Message)
	}
}

func TestHandler_CleanRun(t *testing.T) {
	m := newMetrics(t)
	m.RecordRun(3, 0, nil)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	healthcheck.Handler(m)(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var s healthcheck.Status
	_ = json.NewDecoder(rec.Body).Decode(&s)
	if !s.OK {
		t.Error("expected ok=true")
	}
	if s.Total != 3 || s.Drifted != 0 {
		t.Errorf("counts mismatch: total=%d drifted=%d", s.Total, s.Drifted)
	}
	if s.Message != "all services in sync" {
		t.Errorf("unexpected message: %q", s.Message)
	}
}

func TestHandler_DriftedRun(t *testing.T) {
	m := newMetrics(t)
	m.RecordRun(4, 2, nil)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	healthcheck.Handler(m)(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d", rec.Code)
	}
	var s healthcheck.Status
	_ = json.NewDecoder(rec.Body).Decode(&s)
	if s.OK {
		t.Error("expected ok=false")
	}
	if s.Drifted != 2 {
		t.Errorf("expected drifted=2, got %d", s.Drifted)
	}
}

func TestHandler_ContentType(t *testing.T) {
	m := newMetrics(t)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	healthcheck.Handler(m)(rec, req)

	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected application/json, got %q", ct)
	}
}

func TestHandler_LastRunAt_Populated(t *testing.T) {
	m := newMetrics(t)
	before := time.Now()
	m.RecordRun(1, 0, nil)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	healthcheck.Handler(m)(rec, req)

	var s healthcheck.Status
	_ = json.NewDecoder(rec.Body).Decode(&s)
	if s.LastRunAt.Before(before) {
		t.Error("last_run_at should be after test start")
	}
}
