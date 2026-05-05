package source_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/driftwatch/internal/source"
)

func writeServiceFile(t *testing.T, dir, name, content string) {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writeServiceFile: %v", err)
	}
}

func TestLoadAll_ReturnsDefinitions(t *testing.T) {
	dir := t.TempDir()
	writeServiceFile(t, dir, "api.yaml", "name: api\nversion: \"1.2.0\"\nenvironment: production\nconfig:\n  port: \"8080\"\n")
	writeServiceFile(t, dir, "worker.yaml", "name: worker\nversion: \"0.9.1\"\nenvironment: staging\nconfig:\n  concurrency: \"4\"\n")

	l := source.NewLoader(dir)
	defs, err := l.LoadAll()
	if err != nil {
		t.Fatalf("LoadAll: unexpected error: %v", err)
	}
	if len(defs) != 2 {
		t.Fatalf("expected 2 definitions, got %d", len(defs))
	}
}

func TestLoadOne_Success(t *testing.T) {
	dir := t.TempDir()
	writeServiceFile(t, dir, "api.yaml", "name: api\nversion: \"2.0.0\"\nenvironment: production\nconfig:\n  port: \"443\"\n")

	l := source.NewLoader(dir)
	def, err := l.LoadOne("api")
	if err != nil {
		t.Fatalf("LoadOne: unexpected error: %v", err)
	}
	if def.Name != "api" {
		t.Errorf("expected name 'api', got %q", def.Name)
	}
	if def.Version != "2.0.0" {
		t.Errorf("expected version '2.0.0', got %q", def.Version)
	}
	if def.Config["port"] != "443" {
		t.Errorf("expected config port '443', got %q", def.Config["port"])
	}
}

func TestLoadOne_MissingFile(t *testing.T) {
	l := source.NewLoader(t.TempDir())
	_, err := l.LoadOne("nonexistent")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoadOne_MissingName(t *testing.T) {
	dir := t.TempDir()
	writeServiceFile(t, dir, "bad.yaml", "version: \"1.0.0\"\nenvironment: dev\n")

	l := source.NewLoader(dir)
	_, err := l.LoadOne("bad")
	if err == nil {
		t.Fatal("expected validation error for missing name, got nil")
	}
}

func TestLoadAll_EmptyDir(t *testing.T) {
	l := source.NewLoader(t.TempDir())
	defs, err := l.LoadAll()
	if err != nil {
		t.Fatalf("LoadAll on empty dir: unexpected error: %v", err)
	}
	if len(defs) != 0 {
		t.Errorf("expected 0 definitions, got %d", len(defs))
	}
}
