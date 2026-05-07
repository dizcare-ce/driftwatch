// Package audit provides a structured audit log of drift-check runs,
// recording when checks occurred, how many services were evaluated, and
// whether any drift was detected.
package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Entry represents a single audit log record.
type Entry struct {
	Timestamp   time.Time `json:"timestamp"`
	ServicesRun int       `json:"services_run"`
	DriftCount  int       `json:"drift_count"`
	ErrorCount  int       `json:"error_count"`
	DurationMs  int64     `json:"duration_ms"`
}

// Log appends audit entries to a newline-delimited JSON file.
type Log struct {
	path string
}

// New creates a Log that writes to dir/audit.jsonl, creating dir if needed.
func New(dir string) (*Log, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("audit: create directory: %w", err)
	}
	return &Log{path: filepath.Join(dir, "audit.jsonl")}, nil
}

// Record appends e to the audit log file.
func (l *Log) Record(e Entry) error {
	f, err := os.OpenFile(l.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("audit: open log: %w", err)
	}
	defer f.Close()

	line, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("audit: marshal entry: %w", err)
	}
	_, err = fmt.Fprintf(f, "%s\n", line)
	return err
}

// ReadAll returns all entries from the audit log in append order.
// Returns an empty slice if the file does not yet exist.
func (l *Log) ReadAll() ([]Entry, error) {
	data, err := os.ReadFile(l.path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("audit: read log: %w", err)
	}

	var entries []Entry
	dec := json.NewDecoder(
		// wrap bytes in a reader via a small helper
		newBytesReader(data),
	)
	for dec.More() {
		var e Entry
		if err := dec.Decode(&e); err != nil {
			return nil, fmt.Errorf("audit: decode entry: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, nil
}

// newBytesReader wraps a byte slice in a *bytesReader that satisfies io.Reader.
func newBytesReader(b []byte) *bytesReader { return &bytesReader{b: b} }

type bytesReader struct {
	b   []byte
	pos int
}

func (r *bytesReader) Read(p []byte) (int, error) {
	if r.pos >= len(r.b) {
		return 0, fmt.Errorf("EOF")
	}
	n := copy(p, r.b[r.pos:])
	r.pos += n
	return n, nil
}
