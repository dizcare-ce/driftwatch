package drift

// Diff describes a single field-level discrepancy between expected and actual state.
type Diff struct {
	Field    string `json:"field"`
	Expected string `json:"expected"`
	Actual   string `json:"actual"`
}

// Result holds the outcome of comparing a single service against its definition.
type Result struct {
	Service string `json:"service"`
	Drifted bool   `json:"drifted"`
	Diffs   []Diff `json:"diffs,omitempty"`
}

// Summary returns a human-readable one-line summary of the result.
func (r Result) Summary() string {
	if !r.Drifted {
		return r.Service + ": no drift detected"
	}
	return r.Service + ": drift detected in " + pluralise(len(r.Diffs), "field")
}

func pluralise(n int, word string) string {
	if n == 1 {
		return "1 " + word
	}
	return fmt.Sprintf("%d %ss", n, word)
}
