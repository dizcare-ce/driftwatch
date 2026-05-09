package priority

import (
	"sort"

	"github.com/driftwatch/internal/drift"
)

// Ranked pairs a drift result with its computed priority level.
type Ranked struct {
	Result drift.Result
	Level  Level
}

// Sort returns results ordered from highest to lowest priority.
// Results with the same level are sorted alphabetically by service name.
func Sort(results []drift.Result, scorer *Scorer) []Ranked {
	ranked := make([]Ranked, len(results))
	for i, r := range results {
		ranked[i] = Ranked{Result: r, Level: scorer.Score(r)}
	}
	sort.Slice(ranked, func(i, j int) bool {
		if ranked[i].Level != ranked[j].Level {
			return ranked[i].Level > ranked[j].Level
		}
		return ranked[i].Result.Service < ranked[j].Result.Service
	})
	return ranked
}
