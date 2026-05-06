// Package filter provides service name filtering based on include/exclude glob patterns.
package filter

import (
	"path"
	"strings"
)

// Filter decides which service names should be processed based on
// configured include and exclude patterns.
type Filter struct {
	include []string
	exclude []string
}

// New returns a Filter. An empty includes list means "match everything".
func New(include, exclude []string) *Filter {
	return &Filter{
		include: include,
		exclude: exclude,
	}
}

// Allow reports whether the given service name passes the filter.
// A name is allowed when:
//  1. It matches at least one include pattern (or no include patterns are set), AND
//  2. It does not match any exclude pattern.
func (f *Filter) Allow(name string) bool {
	if matchesAny(name, f.exclude) {
		return false
	}
	if len(f.include) == 0 {
		return true
	}
	return matchesAny(name, f.include)
}

// Apply filters a slice of service names, returning only those that pass.
func (f *Filter) Apply(names []string) []string {
	out := make([]string, 0, len(names))
	for _, n := range names {
		if f.Allow(n) {
			out = append(out, n)
		}
	}
	return out
}

// matchesAny reports whether name matches any of the provided glob patterns.
// Matching is case-insensitive.
func matchesAny(name string, patterns []string) bool {
	lower := strings.ToLower(name)
	for _, p := range patterns {
		matched, err := path.Match(strings.ToLower(p), lower)
		if err == nil && matched {
			return true
		}
	}
	return false
}
