package config

import (
	"bytes"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type DependencyValue struct {
	Version  string   `toml:"version"`
	Features []string `toml:"features"`
}

func (d *DependencyValue) UnmarshalTOML(v any) error {
	switch val := v.(type) {
	case string:
		d.Version = val
		d.Features = nil
		return nil
	case map[string]any:
		if ver, ok := val["version"].(string); ok {
			d.Version = ver
		}
		if feats, ok := val["features"].([]any); ok {
			d.Features = make([]string, 0, len(feats))
			for _, f := range feats {
				if s, ok := f.(string); ok {
					d.Features = append(d.Features, s)
				}
			}
		}
		return nil
	default:
		return fmt.Errorf("cannot unmarshal dependency value from %T", v)
	}
}

func (d DependencyValue) MarshalTOML() ([]byte, error) {
	if len(d.Features) == 0 {
		return []byte(fmt.Sprintf("%q", d.Version)), nil
	}
	var buf bytes.Buffer
	buf.WriteString("{ ")
	buf.WriteString(fmt.Sprintf("version = %q, ", d.Version))
	buf.WriteString("features = [")
	for i, f := range d.Features {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(fmt.Sprintf("%q", f))
	}
	buf.WriteString("] }")
	return buf.Bytes(), nil
}

func (d DependencyValue) VersionString() string {
	return d.Version
}

type Config struct {
	Package        Package                      `toml:"package"`
	Build          Build                        `toml:"build"`
	Dependencies   map[string]DependencyValue   `toml:"dependencies"`
	DevDependencies map[string]DependencyValue   `toml:"dev-dependencies"`
}

type Package struct {
	Name    string `toml:"name"`
	Version string `toml:"version"`
	Edition string `toml:"edition"`
}

type Build struct {
	Compiler          string             `toml:"compiler"`
	STD               string             `toml:"std"`
	Optimization      string             `toml:"optimization"`
	LTO               bool               `toml:"lto"`
	PIC               bool               `toml:"pic"`
	Sanitize          []string           `toml:"sanitize"`
	Target            Target             `toml:"target"`
	Libraries         []Library          `toml:"libraries"`
	Executables       []Executable       `toml:"executables"`
	CMakeDependencies []CMakeDependency  `toml:"cmake-dependencies"`
	Tests             []TestTarget       `toml:"tests"`
}

type Target struct {
	OS   string `toml:"os"`
	Arch string `toml:"arch"`
}

type Library struct {
	Name        string   `toml:"name"`
	Type        string   `toml:"type"`
	Sources     []string `toml:"sources"`
	IncludeDirs []string `toml:"include_dirs"`
}

type Executable struct {
	Name       string   `toml:"name"`
	Sources    []string `toml:"sources"`
	Libraries  []string `toml:"libraries"`
}

type CMakeDependency struct {
	Name         string   `toml:"name"`
	Version      string   `toml:"version"`
	Source       string   `toml:"source"`
	Git          string   `toml:"git"`
	Tag          string   `toml:"tag"`
	CMakeOptions []string `toml:"cmake_options"`
}

type TestTarget struct {
	Name      string   `toml:"name"`
	Type      string   `toml:"type"`
	Sources   []string `toml:"sources"`
	Framework string   `toml:"framework"`
}

func Default() *Config {
	return &Config{
		Package: Package{
			Name:    "my-project",
			Version: "0.1.0",
			Edition: "2024",
		},
		Build: Build{
			Compiler:     "auto",
			STD:          "c++20",
			Optimization: "release",
			LTO:          true,
			PIC:          true,
			Sanitize:     []string{},
			Target: Target{
				OS:   "auto",
				Arch: "auto",
			},
			Libraries:         []Library{},
			Executables:       []Executable{},
			CMakeDependencies: []CMakeDependency{},
			Tests:             []TestTarget{},
		},
		Dependencies:   map[string]DependencyValue{},
		DevDependencies: map[string]DependencyValue{},
	}
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &cfg, nil
}

func Save(path string, cfg *Config) error {
	var buf bytes.Buffer
	if err := toml.NewEncoder(&buf).Encode(cfg); err != nil {
		return fmt.Errorf("failed to encode config: %w", err)
	}

	if err := os.WriteFile(path, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
