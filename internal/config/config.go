package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the top-level driftwatch daemon configuration.
type Config struct {
	PollInterval time.Duration `yaml:"poll_interval"`
	LogLevel     string        `yaml:"log_level"`
	Sources      []Source      `yaml:"sources"`
}

// Source defines a single service definition source to watch.
type Source struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
	Kind string `yaml:"kind"` // e.g. "kubernetes", "docker-compose", "raw"
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		PollInterval: 30 * time.Second,
		LogLevel:     "info",
		Sources:      []Source{},
	}
}

// Load reads a YAML config file from the given path and merges it
// over the defaults.
func Load(path string) (*Config, error) {
	cfg := DefaultConfig()

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("config: open %q: %w", path, err)
	}
	defer f.Close()

	dec := yaml.NewDecoder(f)
	dec.KnownFields(true)
	if err := dec.Decode(cfg); err != nil {
		return nil, fmt.Errorf("config: decode %q: %w", path, err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config: validation: %w", err)
	}

	return cfg, nil
}

// validate performs basic sanity checks on the loaded configuration.
func (c *Config) validate() error {
	if c.PollInterval < time.Second {
		return fmt.Errorf("poll_interval must be at least 1s, got %s", c.PollInterval)
	}
	for i, s := range c.Sources {
		if s.Name == "" {
			return fmt.Errorf("sources[%d]: name is required", i)
		}
		if s.Path == "" {
			return fmt.Errorf("sources[%d]: path is required", i)
		}
		switch s.Kind {
		case "kubernetes", "docker-compose", "raw":
		default:
			return fmt.Errorf("sources[%d]: unsupported kind %q", i, s.Kind)
		}
	}
	return nil
}
