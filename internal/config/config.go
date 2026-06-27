package config

import (
	"bytes"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Package        Package                      `toml:"package"`
	Build          Build                        `toml:"build"`
	Dependencies   map[string]string            `toml:"dependencies"`
	DevDependencies map[string]string            `toml:"dev-dependencies"`
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
		Dependencies:   map[string]string{},
		DevDependencies: map[string]string{},
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
