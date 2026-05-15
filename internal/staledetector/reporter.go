package staledetector

import (
	"fmt"
	"io"
	"text/tabwriter"
	"time"
)

// Report writes a human-readable summary of stale entries to w.
// If there are no stale entries the output notes that all results are current.
func Report(w io.Writer, entries []Entry) error {
	if len(entries) == 0 {
		_, err := fmt.Fprintln(w, "stale-detector: all results are current")
		return err
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "STALE SERVICES (%d)\n", len(entries))
	fmt.Fprintln(tw, "SERVICE\tLAST SEEN\tAGE")
	for _, e := range entries {
		lastSeen := "never"
		if !e.LastSeen.IsZero() {
			lastSeen = e.LastSeen.UTC().Format(time.RFC3339)
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\n", e.Service, lastSeen, fmtDuration(e.Age))
	}
	return tw.Flush()
}

func fmtDuration(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	if h > 0 {
		return fmt.Sprintf("%dh%02dm%02ds", h, m, s)
	}
	if m > 0 {
		return fmt.Sprintf("%dm%02ds", m, s)
	}
	return fmt.Sprintf("%ds", s)
}
