package baseline

import (
	"github.com/driftwatch/internal/cache"
	"github.com/driftwatch/internal/drift"
)

// Checker decides whether a drift result is suppressed by an approved baseline.
type Checker struct {
	store *Store
}

// NewChecker returns a Checker backed by store.
func NewChecker(store *Store) *Checker {
	return &Checker{store: store}
}

// IsSuppressed returns true when the result's drift matches the stored baseline
// checksum, meaning the operator has already acknowledged this drift state.
func (c *Checker) IsSuppressed(r drift.Result) (bool, error) {
	if !r.Drifted() {
		return false, nil
	}

	e, ok, err := c.store.Get(r.Service)
	if err != nil || !ok {
		return false, err
	}

	current := cache.Checksum(r)
	return e.Checksum == current, nil
}

// Approve records the current drift state as an approved baseline.
func (c *Checker) Approve(r drift.Result, note string) error {
	return c.store.Set(Entry{
		Service:  r.Service,
		Checksum: cache.Checksum(r),
		Note:     note,
	})
}

// Revoke removes the baseline approval for a service.
func (c *Checker) Revoke(service string) error {
	return c.store.Delete(service)
}
