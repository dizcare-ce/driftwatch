// Package labelfilter selects or excludes services based on key=value labels
// attached to their source definitions.
package labelfilter

import "strings"

// Filter evaluates service labels against a set of required selectors.
type Filter struct {
	selectors map[string]string
}

// New returns a Filter that matches services carrying all of the provided
// selectors. Selectors are expressed as "key=value" strings; malformed
// entries are silently ignored.
func New(selectors []string) *Filter {
	sm := make(map[string]string, len(selectors))
	for _, s := range selectors {
		k, v, ok := strings.Cut(s, "=")
		if !ok || strings.TrimSpace(k) == "" {
			continue
		}
		sm[strings.ToLower(strings.TrimSpace(k))] = strings.TrimSpace(v)
	}
	return &Filter{selectors: sm}
}

// Match reports whether the provided labels satisfy every selector in f.
// If no selectors were configured, all label maps are accepted.
func (f *Filter) Match(labels map[string]string) bool {
	if len(f.selectors) == 0 {
		return true
	}
	for k, want := range f.selectors {
		got, ok := labels[k]
		if !ok {
			return false
		}
		if !strings.EqualFold(got, want) {
			return false
		}
	}
	return true
}

// MatchAll filters names to those whose labels satisfy f, returning the
// surviving names in the same order they were provided.
func (f *Filter) MatchAll(candidates []Entry) []string {
	out := make([]string, 0, len(candidates))
	for _, c := range candidates {
		if f.Match(c.Labels) {
			out = append(out, c.Name)
		}
	}
	return out
}

// Entry pairs a service name with its label map.
type Entry struct {
	Name   string
	Labels map[string]string
}
