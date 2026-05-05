package notifier

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/driftwatch/internal/drift"
)

// Level represents the severity threshold for notifications.
type Level string

const (
	LevelAll   Level = "all"
	LevelDrift Level = "drift"
	LevelNone  Level = "none"
)

// Notifier decides whether to surface drift results based on a threshold level.
type Notifier struct {
	level  Level
	writer io.Writer
}

// New creates a Notifier writing to w at the given level.
// If w is nil, os.Stderr is used. If level is empty, LevelDrift is used.
func New(level Level, w io.Writer) *Notifier {
	if w == nil {
		w = os.Stderr
	}
	if level == "" {
		level = LevelDrift
	}
	return &Notifier{level: level, writer: w}
}

// Notify writes a notification for each result that meets the threshold.
// Returns the number of results that triggered a notification.
func (n *Notifier) Notify(results []drift.Result) (int, error) {
	if n.level == LevelNone {
		return 0, nil
	}

	count := 0
	var errs []string

	for _, r := range results {
		if n.level == LevelDrift && !r.HasDrift() {
			continue
		}
		if _, err := fmt.Fprintln(n.writer, r.Summary()); err != nil {
			errs = append(errs, err.Error())
			continue
		}
		count++
	}

	if len(errs) > 0 {
		return count, fmt.Errorf("notifier: write errors: %s", strings.Join(errs, "; "))
	}
	return count, nil
}

// ParseLevel converts a string to a Level, returning an error for unknown values.
func ParseLevel(s string) (Level, error) {
	switch Level(strings.ToLower(s)) {
	case LevelAll:
		return LevelAll, nil
	case LevelDrift:
		return LevelDrift, nil
	case LevelNone:
		return LevelNone, nil
	default:
		return "", fmt.Errorf("notifier: unknown level %q (want all|drift|none)", s)
	}
}
