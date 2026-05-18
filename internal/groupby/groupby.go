// Package groupby provides utilities for grouping drift results by
// arbitrary string keys derived from each result, such as owner, label,
// or environment tag.
package groupby

import (
	"fmt"
	"sort"

	"driftwatch/internal/drift"
)

// KeyFunc extracts a grouping key from a drift result.
// Returning an empty string places the result in the "(ungrouped)" bucket.
type KeyFunc func(r drift.Result) string

// Group holds a named collection of drift results.
type Group struct {
	Key     string
	Results []drift.Result
}

// Grouper partitions drift results using a caller-supplied KeyFunc.
type Grouper struct {
	keyFn    KeyFunc
	fallback string
}

// New returns a Grouper that uses keyFn to assign results to buckets.
// If keyFn returns an empty string the result is placed under fallback.
// If fallback is empty it defaults to "(ungrouped)".
func New(keyFn KeyFunc, fallback string) (*Grouper, error) {
	if keyFn == nil {
		return nil, fmt.Errorf("groupby: keyFn must not be nil")
	}
	if fallback == "" {
		fallback = "(ungrouped)"
	}
	return &Grouper{keyFn: keyFn, fallback: fallback}, nil
}

// Apply partitions results into groups and returns them sorted by key.
// An empty input returns nil.
func (g *Grouper) Apply(results []drift.Result) []Group {
	if len(results) == 0 {
		return nil
	}

	buckets := make(map[string][]drift.Result)
	for _, r := range results {
		k := g.keyFn(r)
		if k == "" {
			k = g.fallback
		}
		buckets[k] = append(buckets[k], r)
	}

	groups := make([]Group, 0, len(buckets))
	for k, rs := range buckets {
		groups = append(groups, Group{Key: k, Results: rs})
	}
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Key < groups[j].Key
	})
	return groups
}

// ByService is a convenience KeyFunc that groups by service name.
func ByService(r drift.Result) string { return r.Service }
