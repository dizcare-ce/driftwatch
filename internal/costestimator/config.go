package costestimator

import (
	"fmt"
	"strings"
)

// Config holds the YAML-decoded configuration for the cost estimator.
type Config struct {
	// DefaultCost is applied to any drifted field not listed in Weights.
	DefaultCost float64 `yaml:"default_cost"`
	// Weights maps field names to their cost multipliers.
	Weights map[string]float64 `yaml:"weights"`
}

// DefaultConfig returns a sensible out-of-the-box configuration.
func DefaultConfig() Config {
	return Config{
		DefaultCost: 1.0,
		Weights: map[string]float64{
			"replicas": 5.0,
			"image":    3.0,
			"env":      2.0,
			"resources": 4.0,
		},
	}
}

// Validate checks that all weight values are non-negative and the default cost
// is positive.
func (c Config) Validate() error {
	if c.DefaultCost <= 0 {
		return fmt.Errorf("costestimator: default_cost must be positive, got %.2f", c.DefaultCost)
	}
	for k, v := range c.Weights {
		if v < 0 {
			return fmt.Errorf("costestimator: weight for field %q must be non-negative, got %.2f", k, v)
		}
	}
	return nil
}

// NewFromConfig constructs an Estimator from a validated Config.
func NewFromConfig(c Config) (*Estimator, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}
	normalised := make(FieldWeight, len(c.Weights))
	for k, v := range c.Weights {
		normalised[strings.ToLower(k)] = v
	}
	return &Estimator{weights: normalised, defaultCost: c.DefaultCost}, nil
}
