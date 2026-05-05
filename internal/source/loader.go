package source

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// ServiceDefinition represents a service's expected configuration state
// as defined in the source (e.g., a Git repo or local directory).
type ServiceDefinition struct {
	Name        string            `yaml:"name"`
	Version     string            `yaml:"version"`
	Environment string            `yaml:"environment"`
	Config      map[string]string `yaml:"config"`
}

// Loader reads service definitions from a source directory.
type Loader struct {
	BasePath string
}

// NewLoader creates a Loader rooted at basePath.
func NewLoader(basePath string) *Loader {
	return &Loader{BasePath: basePath}
}

// LoadAll reads every *.yaml file in BasePath and returns the parsed definitions.
func (l *Loader) LoadAll() ([]ServiceDefinition, error) {
	pattern := filepath.Join(l.BasePath, "*.yaml")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("source: glob %q: %w", pattern, err)
	}

	var defs []ServiceDefinition
	for _, path := range matches {
		def, err := l.loadFile(path)
		if err != nil {
			return nil, err
		}
		defs = append(defs, def)
	}
	return defs, nil
}

// LoadOne reads and parses a single service definition file.
func (l *Loader) LoadOne(name string) (ServiceDefinition, error) {
	path := filepath.Join(l.BasePath, name+".yaml")
	return l.loadFile(path)
}

func (l *Loader) loadFile(path string) (ServiceDefinition, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return ServiceDefinition{}, fmt.Errorf("source: read %q: %w", path, err)
	}

	var def ServiceDefinition
	if err := yaml.Unmarshal(data, &def); err != nil {
		return ServiceDefinition{}, fmt.Errorf("source: parse %q: %w", path, err)
	}

	if def.Name == "" {
		return ServiceDefinition{}, fmt.Errorf("source: %q: missing required field 'name'", path)
	}
	return def, nil
}
