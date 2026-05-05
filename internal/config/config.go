package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Kind identifies the type of backend being checked.
type Kind string

const (
	KindKubernetes Kind = "kubernetes"
	KindDocker     Kind = "docker"
)

// Config is the top-level driftwatch configuration.
type Config struct {
	Sources   SourcesConfig  `yaml:"sources"`
	Drift     DriftConfig    `yaml:"drift"`
	Reporter  ReporterConfig `yaml:"reporter"`
	Scheduler SchedulerConfig `yaml:"scheduler"`
}

// SourcesConfig describes where service definitions are loaded from.
type SourcesConfig struct {
	Dir string `yaml:"dir"`
}

// DriftConfig controls drift-detection behaviour.
type DriftConfig struct {
	Kind Kind `yaml:"kind"`
}

// ReporterConfig controls output format.
type ReporterConfig struct {
	Format string `yaml:"format"` // "text" or "json"
}

// SchedulerConfig controls how often drift checks run.
type SchedulerConfig struct {
	Interval time.Duration `yaml:"interval"`
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Sources:  SourcesConfig{Dir: "services"},
		Drift:    DriftConfig{Kind: KindKubernetes},
		Reporter: ReporterConfig{Format: "text"},
		Scheduler: SchedulerConfig{Interval: 5 * time.Minute},
	}
}

// Load reads a YAML config file from path, merging over defaults.
func Load(path string) (Config, error) {
	cfg := DefaultConfig()

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return cfg, fmt.Errorf("config file not found: %s", path)
		}
		return cfg, fmt.Errorf("reading config: %w", err)
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("parsing config: %w", err)
	}

	if cfg.Drift.Kind != KindKubernetes && cfg.Drift.Kind != KindDocker {
		return cfg, fmt.Errorf("unsupported drift kind %q: must be 'kubernetes' or 'docker'", cfg.Drift.Kind)
	}

	if cfg.Scheduler.Interval <= 0 {
		cfg.Scheduler.Interval = DefaultConfig().Scheduler.Interval
	}

	return cfg, nil
}
