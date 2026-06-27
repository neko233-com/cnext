package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefault(t *testing.T) {
	cfg := Default()

	if cfg.Package.Name != "my-project" {
		t.Errorf("Default package name = %q, want %q", cfg.Package.Name, "my-project")
	}
	if cfg.Package.Version != "0.1.0" {
		t.Errorf("Default package version = %q, want %q", cfg.Package.Version, "0.1.0")
	}
	if cfg.Package.Edition != "2024" {
		t.Errorf("Default package edition = %q, want %q", cfg.Package.Edition, "2024")
	}

	if cfg.Build.Compiler != "auto" {
		t.Errorf("Default build compiler = %q, want %q", cfg.Build.Compiler, "auto")
	}
	if cfg.Build.STD != "c++20" {
		t.Errorf("Default build std = %q, want %q", cfg.Build.STD, "c++20")
	}
	if cfg.Build.Optimization != "release" {
		t.Errorf("Default build optimization = %q, want %q", cfg.Build.Optimization, "release")
	}
	if !cfg.Build.LTO {
		t.Errorf("Default build lto = %v, want %v", cfg.Build.LTO, true)
	}
	if !cfg.Build.PIC {
		t.Errorf("Default build pic = %v, want %v", cfg.Build.PIC, true)
	}

	if cfg.Build.Target.OS != "auto" {
		t.Errorf("Default target os = %q, want %q", cfg.Build.Target.OS, "auto")
	}
	if cfg.Build.Target.Arch != "auto" {
		t.Errorf("Default target arch = %q, want %q", cfg.Build.Target.Arch, "auto")
	}

	if cfg.Dependencies == nil {
		t.Error("Default dependencies should not be nil")
	}
	if cfg.DevDependencies == nil {
		t.Error("Default dev-dependencies should not be nil")
	}
}

func TestLoad(t *testing.T) {
	content := `
[package]
name = "test-project"
version = "1.0.0"
edition = "2024"

[build]
compiler = "gcc"
std = "c++23"
optimization = "debug"
lto = false
pic = true
sanitize = ["address"]

[build.target]
os = "linux"
arch = "x86_64"

[[build.libraries]]
name = "mylib"
type = "static"
sources = ["src/**/*.cpp"]
include_dirs = ["include"]

[[build.executables]]
name = "myapp"
sources = ["main.cpp"]
libraries = ["mylib"]

[[build.cmake-dependencies]]
name = "opencv"
version = "4.8.0"
source = "registry"
git = "https://github.com/opencv/opencv.git"
tag = "4.8.0"
cmake_options = ["-DBUILD_SHARED_LIBS=ON"]

[[build.tests]]
name = "my-tests"
type = "unit"
sources = ["tests/**/*.cpp"]
framework = "cnext"

[dependencies]
fmt = "10.2.1"
spdlog = "1.13.0"

[dev-dependencies]
gtest = "1.14.0"
benchmark = "1.8.3"
`
	dir := t.TempDir()
	configPath := filepath.Join(dir, "cnext.toml")
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Package.Name != "test-project" {
		t.Errorf("Package.Name = %q, want %q", cfg.Package.Name, "test-project")
	}
	if cfg.Package.Version != "1.0.0" {
		t.Errorf("Package.Version = %q, want %q", cfg.Package.Version, "1.0.0")
	}
	if cfg.Package.Edition != "2024" {
		t.Errorf("Package.Edition = %q, want %q", cfg.Package.Edition, "2024")
	}

	if cfg.Build.Compiler != "gcc" {
		t.Errorf("Build.Compiler = %q, want %q", cfg.Build.Compiler, "gcc")
	}
	if cfg.Build.STD != "c++23" {
		t.Errorf("Build.STD = %q, want %q", cfg.Build.STD, "c++23")
	}
	if cfg.Build.Optimization != "debug" {
		t.Errorf("Build.Optimization = %q, want %q", cfg.Build.Optimization, "debug")
	}
	if cfg.Build.LTO {
		t.Errorf("Build.LTO = %v, want %v", cfg.Build.LTO, false)
	}
	if !cfg.Build.PIC {
		t.Errorf("Build.PIC = %v, want %v", cfg.Build.PIC, true)
	}
	if len(cfg.Build.Sanitize) != 1 || cfg.Build.Sanitize[0] != "address" {
		t.Errorf("Build.Sanitize = %v, want [address]", cfg.Build.Sanitize)
	}

	if cfg.Build.Target.OS != "linux" {
		t.Errorf("Build.Target.OS = %q, want %q", cfg.Build.Target.OS, "linux")
	}
	if cfg.Build.Target.Arch != "x86_64" {
		t.Errorf("Build.Target.Arch = %q, want %q", cfg.Build.Target.Arch, "x86_64")
	}

	if len(cfg.Build.Libraries) != 1 {
		t.Fatalf("Build.Libraries len = %d, want 1", len(cfg.Build.Libraries))
	}
	lib := cfg.Build.Libraries[0]
	if lib.Name != "mylib" {
		t.Errorf("Library.Name = %q, want %q", lib.Name, "mylib")
	}
	if lib.Type != "static" {
		t.Errorf("Library.Type = %q, want %q", lib.Type, "static")
	}

	if len(cfg.Build.Executables) != 1 {
		t.Fatalf("Build.Executables len = %d, want 1", len(cfg.Build.Executables))
	}
	exec := cfg.Build.Executables[0]
	if exec.Name != "myapp" {
		t.Errorf("Executable.Name = %q, want %q", exec.Name, "myapp")
	}

	if len(cfg.Build.CMakeDependencies) != 1 {
		t.Fatalf("Build.CMakeDependencies len = %d, want 1", len(cfg.Build.CMakeDependencies))
	}
	cmake := cfg.Build.CMakeDependencies[0]
	if cmake.Name != "opencv" {
		t.Errorf("CMakeDependency.Name = %q, want %q", cmake.Name, "opencv")
	}
	if cmake.Source != "registry" {
		t.Errorf("CMakeDependency.Source = %q, want %q", cmake.Source, "registry")
	}

	if len(cfg.Build.Tests) != 1 {
		t.Fatalf("Build.Tests len = %d, want 1", len(cfg.Build.Tests))
	}
	test := cfg.Build.Tests[0]
	if test.Name != "my-tests" {
		t.Errorf("Test.Name = %q, want %q", test.Name, "my-tests")
	}
	if test.Framework != "cnext" {
		t.Errorf("Test.Framework = %q, want %q", test.Framework, "cnext")
	}

	if cfg.Dependencies["fmt"] != "10.2.1" {
		t.Errorf("Dependencies[fmt] = %q, want %q", cfg.Dependencies["fmt"], "10.2.1")
	}
	if cfg.Dependencies["spdlog"] != "1.13.0" {
		t.Errorf("Dependencies[spdlog] = %q, want %q", cfg.Dependencies["spdlog"], "1.13.0")
	}

	if cfg.DevDependencies["gtest"] != "1.14.0" {
		t.Errorf("DevDependencies[gtest] = %q, want %q", cfg.DevDependencies["gtest"], "1.14.0")
	}
	if cfg.DevDependencies["benchmark"] != "1.8.3" {
		t.Errorf("DevDependencies[benchmark] = %q, want %q", cfg.DevDependencies["benchmark"], "1.8.3")
	}
}

func TestLoadFileNotFound(t *testing.T) {
	_, err := Load("/nonexistent/cnext.toml")
	if err == nil {
		t.Error("Load() should return error for nonexistent file")
	}
}

func TestSave(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "cnext.toml")

	cfg := Default()
	cfg.Package.Name = "saved-project"
	cfg.Package.Version = "2.0.0"
	cfg.Build.Compiler = "clang"
	cfg.Build.STD = "c++23"
	cfg.Dependencies["fmt"] = "10.2.1"
	cfg.DevDependencies["gtest"] = "1.14.0"

	if err := Save(configPath, cfg); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	loaded, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load() after Save() error = %v", err)
	}

	if loaded.Package.Name != "saved-project" {
		t.Errorf("Loaded Package.Name = %q, want %q", loaded.Package.Name, "saved-project")
	}
	if loaded.Package.Version != "2.0.0" {
		t.Errorf("Loaded Package.Version = %q, want %q", loaded.Package.Version, "2.0.0")
	}
	if loaded.Build.Compiler != "clang" {
		t.Errorf("Loaded Build.Compiler = %q, want %q", loaded.Build.Compiler, "clang")
	}
	if loaded.Build.STD != "c++23" {
		t.Errorf("Loaded Build.STD = %q, want %q", loaded.Build.STD, "c++23")
	}
	if loaded.Dependencies["fmt"] != "10.2.1" {
		t.Errorf("Loaded Dependencies[fmt] = %q, want %q", loaded.Dependencies["fmt"], "10.2.1")
	}
	if loaded.DevDependencies["gtest"] != "1.14.0" {
		t.Errorf("Loaded DevDependencies[gtest] = %q, want %q", loaded.DevDependencies["gtest"], "1.14.0")
	}
}

func TestSaveRoundTrip(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "cnext.toml")

	cfg := Default()
	cfg.Build.Libraries = []Library{
		{
			Name:       "mylib",
			Type:       "static",
			Sources:    []string{"src/**/*.cpp"},
			IncludeDirs: []string{"include"},
		},
	}
	cfg.Build.Executables = []Executable{
		{
			Name:       "myapp",
			Sources:    []string{"main.cpp"},
			Libraries:  []string{"mylib"},
		},
	}
	cfg.Build.CMakeDependencies = []CMakeDependency{
		{
			Name:    "opencv",
			Version: "4.8.0",
			Source:  "registry",
			Git:     "https://github.com/opencv/opencv.git",
			Tag:     "4.8.0",
			CMakeOptions: []string{"-DBUILD_SHARED_LIBS=ON"},
		},
	}
	cfg.Build.Tests = []TestTarget{
		{
			Name:      "my-tests",
			Type:      "unit",
			Sources:   []string{"tests/**/*.cpp"},
			Framework: "cnext",
		},
	}

	if err := Save(configPath, cfg); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	loaded, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if len(loaded.Build.Libraries) != 1 {
		t.Fatalf("Libraries len = %d, want 1", len(loaded.Build.Libraries))
	}
	if loaded.Build.Libraries[0].Name != "mylib" {
		t.Errorf("Library.Name = %q, want %q", loaded.Build.Libraries[0].Name, "mylib")
	}

	if len(loaded.Build.Executables) != 1 {
		t.Fatalf("Executables len = %d, want 1", len(loaded.Build.Executables))
	}
	if loaded.Build.Executables[0].Name != "myapp" {
		t.Errorf("Executable.Name = %q, want %q", loaded.Build.Executables[0].Name, "myapp")
	}

	if len(loaded.Build.CMakeDependencies) != 1 {
		t.Fatalf("CMakeDependencies len = %d, want 1", len(loaded.Build.CMakeDependencies))
	}
	if loaded.Build.CMakeDependencies[0].Name != "opencv" {
		t.Errorf("CMakeDependency.Name = %q, want %q", loaded.Build.CMakeDependencies[0].Name, "opencv")
	}

	if len(loaded.Build.Tests) != 1 {
		t.Fatalf("Tests len = %d, want 1", len(loaded.Build.Tests))
	}
	if loaded.Build.Tests[0].Name != "my-tests" {
		t.Errorf("Test.Name = %q, want %q", loaded.Build.Tests[0].Name, "my-tests")
	}
}

func TestLoadInvalidToml(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "cnext.toml")

	invalidContent := `this is not valid toml`
	if err := os.WriteFile(configPath, []byte(invalidContent), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	_, err := Load(configPath)
	if err == nil {
		t.Error("Load() should return error for invalid TOML")
	}
}

func TestLoadPartialConfig(t *testing.T) {
	content := `
[package]
name = "minimal-project"
`
	dir := t.TempDir()
	configPath := filepath.Join(dir, "cnext.toml")
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Package.Name != "minimal-project" {
		t.Errorf("Package.Name = %q, want %q", cfg.Package.Name, "minimal-project")
	}
	if cfg.Build.Compiler != "" {
		t.Errorf("Build.Compiler should be empty for partial config, got %q", cfg.Build.Compiler)
	}
}
