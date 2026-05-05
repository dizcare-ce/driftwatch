// Package watcher monitors source definition files for changes on disk
// and triggers a drift check when a modification is detected.
package watcher

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"
)

// FileWatcher polls a directory for file modifications and calls onChange
// whenever a watched file's modification time changes.
type FileWatcher struct {
	dir      string
	interval time.Duration
	onChange func(path string)
	states   map[string]time.Time
}

// New creates a FileWatcher that watches dir at the given poll interval.
// onChange is invoked with the changed file path.
func New(dir string, interval time.Duration, onChange func(path string)) *FileWatcher {
	return &FileWatcher{
		dir:      dir,
		interval: interval,
		onChange: onChange,
		states:   make(map[string]time.Time),
	}
}

// Run starts the polling loop and blocks until ctx is cancelled.
func (fw *FileWatcher) Run(ctx context.Context) error {
	if err := fw.snapshot(); err != nil {
		return err
	}
	ticker := time.NewTicker(fw.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			fw.check()
		}
	}
}

func (fw *FileWatcher) snapshot() error {
	return filepath.WalkDir(fw.dir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		fw.states[path] = info.ModTime()
		return nil
	})
}

func (fw *FileWatcher) check() {
	err := filepath.WalkDir(fw.dir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		prev, known := fw.states[path]
		if !known || info.ModTime().After(prev) {
			fw.states[path] = info.ModTime()
			fw.onChange(path)
		}
		return nil
	})
	if err != nil {
		log.Printf("watcher: scan error: %v", err)
	}
}
