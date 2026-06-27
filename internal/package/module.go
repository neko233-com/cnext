package pkg

import (
	"bytes"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type Module struct {
	Name    string            `toml:"name"`
	Version string            `toml:"version"`
	Deps    map[string]string `toml:"dependencies"`
	DevDeps map[string]string `toml:"dev-dependencies"`
}

func LoadModule(path string) (*Module, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read module file: %w", err)
	}

	var cfg struct {
		Package struct {
			Name    string `toml:"name"`
			Version string `toml:"version"`
		} `toml:"package"`
		Dependencies   map[string]string `toml:"dependencies"`
		DevDependencies map[string]string `toml:"dev-dependencies"`
	}

	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse module file: %w", err)
	}

	return &Module{
		Name:    cfg.Package.Name,
		Version: cfg.Package.Version,
		Deps:    cfg.Dependencies,
		DevDeps: cfg.DevDependencies,
	}, nil
}

func SaveModule(mod *Module, path string) error {
	cfg := struct {
		Package struct {
			Name    string `toml:"name"`
			Version string `toml:"version"`
		} `toml:"package"`
		Dependencies   map[string]string `toml:"dependencies"`
		DevDependencies map[string]string `toml:"dev-dependencies"`
	}{
		Dependencies:   mod.Deps,
		DevDependencies: mod.DevDeps,
	}
	cfg.Package.Name = mod.Name
	cfg.Package.Version = mod.Version

	var buf bytes.Buffer
	if err := toml.NewEncoder(&buf).Encode(cfg); err != nil {
		return fmt.Errorf("failed to encode module: %w", err)
	}

	if err := os.WriteFile(path, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write module file: %w", err)
	}

	return nil
}
