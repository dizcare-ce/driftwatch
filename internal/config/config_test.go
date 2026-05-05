package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/driftwatch/driftwatch/internal/config"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "driftwatch-*.yaml")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_Defaults(t *testing.T) {
	path := writeTempConfig(t, "sources:\n  - name: svc\n    path: /tmp/svc.yaml\n    kind: raw\n")
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.PollInterval != 30*time.Second {
		t.Errorf("expected default poll_interval 30s, got %s", cfg.PollInterval)
	}
	if cfg.LogLevel != "info" {
		t.Errorf("expected default log_level info, got %s", cfg.LogLevel)
	}
}

func TestLoad_Override(t *testing.T) {
	path := writeTempConfig(t, "poll_interval: 10s\nlog_level: debug\nsources:\n  - name: api\n    path: /etc/api.yaml\n    kind: kubernetes\n")
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.PollInterval != 10*time.Second {
		t.Errorf("expected 10s, got %s", cfg.PollInterval)
	}
	if cfg.LogLevel != "debug" {
		t.Errorf("expected debug, got %s", cfg.LogLevel)
	}
	if len(cfg.Sources) != 1 || cfg.Sources[0].Kind != "kubernetes" {
		t.Errorf("unexpected sources: %+v", cfg.Sources)
	}
}

func TestLoad_InvalidKind(t *testing.T) {
	path := writeTempConfig(t, "sources:\n  - name: svc\n    path: /tmp/x.yaml\n    kind: unknown\n")
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected validation error for unknown kind")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := config.Load("/nonexistent/path.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_PollIntervalTooShort(t *testing.T) {
	path := writeTempConfig(t, "poll_interval: 500ms\nsources:\n  - name: s\n    path: /p\n    kind: raw\n")
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected validation error for short poll_interval")
	}
}
