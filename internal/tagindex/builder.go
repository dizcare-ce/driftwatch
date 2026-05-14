package tagindex

import "driftwatch/internal/source"

// Builder populates an Index from a slice of source definitions.
type Builder struct {
	idx *Index
}

// NewBuilder returns a Builder backed by idx.
func NewBuilder(idx *Index) *Builder {
	return &Builder{idx: idx}
}

// Build clears the index and re-populates it from defs.
// Labels are read from Definition.Labels; services without labels are
// still registered with an empty label map so they appear in All().
func (b *Builder) Build(defs []source.Definition) {
	// Remove services no longer present.
	existing := b.idx.All()
	current := make(map[string]struct{}, len(defs))
	for _, d := range defs {
		current[d.Name] = struct{}{}
	}
	for _, e := range existing {
		if _, ok := current[e.Service]; !ok {
			b.idx.Remove(e.Service)
		}
	}

	// Upsert each definition.
	for _, d := range defs {
		labels := d.Labels
		if labels == nil {
			labels = map[string]string{}
		}
		b.idx.Add(d.Name, labels)
	}
}
