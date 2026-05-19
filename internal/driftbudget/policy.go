package driftbudget

import (
	"fmt"
	"time"
)

// Policy describes the configuration for a Budget.
type Policy struct {
	// Limit is the maximum number of drifted fields allowed within Window.
	Limit int `yaml:"limit"`
	// Window is the rolling duration over which drift is counted.
	Window time.Duration `yaml:"window"`
}

// DefaultPolicy returns a sensible starting policy.
func DefaultPolicy() Policy {
	return Policy{
		Limit:  50,
		Window: 1 * time.Hour,
	}
}

// Validate checks that the policy fields are usable.
func (p Policy) Validate() error {
	if p.Limit <= 0 {
		return fmt.Errorf("driftbudget: policy limit must be positive")
	}
	if p.Window <= 0 {
		return fmt.Errorf("driftbudget: policy window must be positive")
	}
	return nil
}

// NewFromPolicy constructs a Budget from a Policy.
func NewFromPolicy(p Policy) (*Budget, error) {
	if err := p.Validate(); err != nil {
		return nil, err
	}
	return New(p.Limit, p.Window)
}
