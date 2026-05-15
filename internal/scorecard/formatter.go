package scorecard

import (
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"
)

// Format writes the scorecard entries to w in the requested format.
// Supported formats: "text" (default), "json".
func Format(w io.Writer, entries []Entry, format string) error {
	switch format {
	case "json":
		return formatJSON(w, entries)
	default:
		return formatText(w, entries)
	}
}

func formatText(w io.Writer, entries []Entry) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "SERVICE\tSCORE\tGRADE\tDRIFTED\tTOTAL")
	for _, e := range entries {
		fmt.Fprintf(tw, "%s\t%d\t%s\t%d\t%d\n",
			e.Service, e.Score, e.Grade, e.Drifted, e.Total)
	}
	return tw.Flush()
}

func formatJSON(w io.Writer, entries []Entry) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(entries)
}
