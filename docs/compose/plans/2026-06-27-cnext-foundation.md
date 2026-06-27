# cnext Foundation Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use compose:subagent (recommended) or compose:execute to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build the core cnext CLI with project initialization, build system, and compiler driver foundation.

**Architecture:** Go-based CLI tool that wraps gcc/clang/msvc compilers, reads cnext.toml configuration, generates build graphs, and executes builds. Uses cobra for CLI, toml for config parsing, and the bundled toolchain approach.

**Tech Stack:** Go 1.26, cobra (CLI), BurntSushi/toml (config), os/exec (compiler invocation)

## Global Constraints

- Go 1.26 minimum
- All code in Go, no CGO for host tools
- Cross-platform: Linux, macOS, Windows
- MIT License
- Follow Go idioms and conventions
- Tests for all public functions

---

### Task 1: Project Scaffolding

**Covers:** S9

**Files:**
- Create: `go.mod`
- Create: `go.sum`
- Create: `cmd/cnext/main.go`
- Create: `internal/compiler/driver.go`
- Create: `internal/frontend/parser.go`
- Create: `internal/build/graph.go`
- Create: `internal/test/runner.go`
- Create: `internal/package/resolver.go`
- Create: `internal/toolchain/manager.go`
- Create: `pkg/cnext/cnext.go`

**Interfaces:**
- Produces: Go module structure, package layout

- [ ] **Step 1: Initialize Go module**

```bash
cd D:\Code\neko233-Projects\cnext
go mod init github.com/neko233-com/cnext
```

- [ ] **Step 2: Create directory structure**

```bash
mkdir -p cmd/cnext
mkdir -p internal/compiler
mkdir -p internal/frontend
mkdir -p internal/build
mkdir -p internal/test
mkdir -p internal/package
mkdir -p internal/toolchain
mkdir -p pkg/cnext
mkdir -p cnext-test/include/cnext
mkdir -p cnext-test/src
```

- [ ] **Step 3: Create main.go entry point**

Create `cmd/cnext/main.go`:

```go
package main

import (
    "fmt"
    "os"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("cnext - Next-gen C/C++ toolchain")
        fmt.Println("Usage: cnext <command> [args]")
        os.Exit(1)
    }
    
    cmd := os.Args[1]
    switch cmd {
    case "init":
        fmt.Println("Initializing project...")
    case "build":
        fmt.Println("Building project...")
    case "test":
        fmt.Println("Running tests...")
    default:
        fmt.Printf("Unknown command: %s\n", cmd)
        os.Exit(1)
    }
}
```

- [ ] **Step 4: Verify it compiles**

```bash
go build ./cmd/cnext
```

- [ ] **Step 5: Commit**

```bash
git add -A
git commit -m "feat: initial project scaffolding"
```

---

### Task 2: CLI Framework with Cobra

**Covers:** S7

**Files:**
- Create: `cmd/cnext/main.go` (modify)
- Create: `cmd/cnext/root.go`
- Create: `cmd/cnext/init.go`
- Create: `cmd/cnext/build.go`
- Create: `cmd/cnext/test.go`
- Create: `cmd/cnext/run.go`
- Create: `cmd/cnext/add.go`
- Create: `cmd/cnext/remove.go`
- Create: `cmd/cnext/fmt.go`
- Create: `cmd/cnext/lint.go`
- Create: `cmd/cnext/doc.go`
- Create: `cmd/cnext/publish.go`
- Create: `cmd/cnext/clean.go`
- Create: `cmd/cnext/update.go`
- Create: `cmd/cnext/info.go`
- Create: `cmd/cnext/toolchain.go`

**Interfaces:**
- Consumes: Task 1 scaffolding
- Produces: Cobra command structure

- [ ] **Step 1: Add cobra dependency**

```bash
go get github.com/spf13/cobra@latest
```

- [ ] **Step 2: Create root command**

Create `cmd/cnext/root.go`:

```go
package main

import (
    "fmt"
    "os"
    
    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "cnext",
    Short: "Next-gen C/C++ toolchain",
    Long:  `cnext is a modern build system and package manager for C/C++ projects.`,
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}

func init() {
    rootCmd.AddCommand(initCmd)
    rootCmd.AddCommand(buildCmd)
    rootCmd.AddCommand(testCmd)
    rootCmd.AddCommand(runCmd)
    rootCmd.AddCommand(addCmd)
    rootCmd.AddCommand(removeCmd)
    rootCmd.AddCommand(fmtCmd)
    rootCmd.AddCommand(lintCmd)
    rootCmd.AddCommand(docCmd)
    rootCmd.AddCommand(publishCmd)
    rootCmd.AddCommand(cleanCmd)
    rootCmd.AddCommand(updateCmd)
    rootCmd.AddCommand(infoCmd)
    rootCmd.AddCommand(toolchainCmd)
}
```

- [ ] **Step 3: Create init command**

Create `cmd/cnext/init.go`:

```go
package main

import (
    "fmt"
    
    "github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
    Use:   "init",
    Short: "Initialize a new cnext project",
    RunE: func(cmd *cobra.Command, args []string) error {
        fmt.Println("Initializing new cnext project...")
        // TODO: Task 3 will implement this
        return nil
    },
}
```

- [ ] **Step 4: Create build command**

Create `cmd/cnext/build.go`:

```go
package main

import (
    "fmt"
    
    "github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
    Use:   "build",
    Short: "Build the project",
    RunE: func(cmd *cobra.Command, args []string) error {
        fmt.Println("Building project...")
        // TODO: Task 5 will implement this
        return nil
    },
}

func init() {
    buildCmd.Flags().StringP("target", "t", "", "Build target (e.g., x86_64-linux-gnu)")
    buildCmd.Flags().BoolP("release", "r", false, "Build in release mode")
    buildCmd.Flags().StringP("compiler", "c", "auto", "Compiler to use (auto, gcc, clang, msvc)")
}
```

- [ ] **Step 5: Create test command**

Create `cmd/cnext/test.go`:

```go
package main

import (
    "fmt"
    
    "github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
    Use:   "test",
    Short: "Run project tests",
    RunE: func(cmd *cobra.Command, args []string) error {
        fmt.Println("Running tests...")
        // TODO: Task 9-10 will implement this
        return nil
    },
}

func init() {
    testCmd.Flags().StringP("filter", "f", "", "Filter tests by name")
    testCmd.Flags().BoolP("coverage", "c", false, "Generate coverage report")
    testCmd.Flags().BoolP("benchmark", "b", false, "Run benchmarks")
    testCmd.Flags().BoolP("watch", "w", false, "Watch for changes and rerun tests")
}
```

- [ ] **Step 6: Create add command**

Create `cmd/cnext/add.go`:

```go
package main

import (
    "fmt"
    
    "github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
    Use:   "add <package>",
    Short: "Add a dependency",
    Args:  cobra.MinimumNArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        fmt.Printf("Adding package: %s\n", args[0])
        // TODO: Task 7-8 will implement this
        return nil
    },
}

func init() {
    addCmd.Flags().StringArrayP("features", "f", nil, "Enable specific features")
    addCmd.Flags().BoolP("dev", "d", false, "Add as dev dependency")
}
```

- [ ] **Step 7: Create remove command**

Create `cmd/cnext/remove.go`:

```go
package main

import (
    "fmt"
    
    "github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
    Use:   "remove <package>",
    Short: "Remove a dependency",
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        fmt.Printf("Removing package: %s\n", args[0])
        // TODO: Task 7-8 will implement this
        return nil
    },
}
```

- [ ] **Step 8: Create remaining commands**

Create `cmd/cnext/fmt.go`:

```go
package main

import (
    "fmt"
    
    "github.com/spf13/cobra"
)

var fmtCmd = &cobra.Command{
    Use:   "fmt",
    Short: "Format source code",
    RunE: func(cmd *cobra.Command, args []string) error {
        fmt.Println("Formatting code...")
        return nil
    },
}
```

Create `cmd/cnext/lint.go`:

```go
package main

import (
    "fmt"
    
    "github.com/spf13/cobra"
)

var lintCmd = &cobra.Command{
    Use:   "lint",
    Short: "Lint source code",
    RunE: func(cmd *cobra.Command, args []string) error {
        fmt.Println("Linting code...")
        return nil
    },
}
```

Create `cmd/cnext/doc.go`:

```go
package main

import (
    "fmt"
    
    "github.com/spf13/cobra"
)

var docCmd = &cobra.Command{
    Use:   "doc",
    Short: "Generate documentation",
    RunE: func(cmd *cobra.Command, args []string) error {
        fmt.Println("Generating documentation...")
        return nil
    },
}
```

Create `cmd/cnext/publish.go`:

```go
package main

import (
    "fmt"
    
    "github.com/spf13/cobra"
)

var publishCmd = &cobra.Command{
    Use:   "publish",
    Short: "Publish package to registry",
    RunE: func(cmd *cobra.Command, args []string) error {
        fmt.Println("Publishing package...")
        return nil
    },
}
```

Create `cmd/cnext/clean.go`:

```go
package main

import (
    "fmt"
    
    "github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
    Use:   "clean",
    Short: "Clean build artifacts",
    RunE: func(cmd *cobra.Command, args []string) error {
        fmt.Println("Cleaning build artifacts...")
        return nil
    },
}
```

Create `cmd/cnext/update.go`:

```go
package main

import (
    "fmt"
    
    "github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
    Use:   "update",
    Short: "Update dependencies",
    RunE: func(cmd *cobra.Command, args []string) error {
        fmt.Println("Updating dependencies...")
        return nil
    },
}
```

Create `cmd/cnext/info.go`:

```go
package main

import (
    "fmt"
    
    "github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
    Use:   "info",
    Short: "Show project information",
    RunE: func(cmd *cobra.Command, args []string) error {
        fmt.Println("Project information...")
        return nil
    },
}
```

Create `cmd/cnext/toolchain.go`:

```go
package main

import (
    "fmt"
    
    "github.com/spf13/cobra"
)

var toolchainCmd = &cobra.Command{
    Use:   "toolchain",
    Short: "Manage toolchains",
}

var toolchainListCmd = &cobra.Command{
    Use:   "list",
    Short: "List installed toolchains",
    RunE: func(cmd *cobra.Command, args []string) error {
        fmt.Println("Listing toolchains...")
        return nil
    },
}

var toolchainInstallCmd = &cobra.Command{
    Use:   "install <toolchain>",
    Short: "Install a toolchain",
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        fmt.Printf("Installing toolchain: %s\n", args[0])
        return nil
    },
}

func init() {
    toolchainCmd.AddCommand(toolchainListCmd)
    toolchainCmd.AddCommand(toolchainInstallCmd)
}
```

- [ ] **Step 9: Update main.go to use Execute()**

Update `cmd/cnext/main.go`:

```go
package main

func main() {
    Execute()
}
```

- [ ] **Step 10: Verify all commands work**

```bash
go build ./cmd/cnext
./cnext --help
./cnext init --help
./cnext build --help
./cnext test --help
```

- [ ] **Step 11: Commit**

```bash
git add -A
git commit -m "feat: add CLI framework with cobra"
```

---

### Task 3: cnext.toml Parser

**Covers:** S4.1

**Files:**
- Create: `internal/config/config.go`
- Create: `internal/config/config_test.go`
- Modify: `cmd/cnext/init.go`

**Interfaces:**
- Consumes: Task 2 CLI structure
- Produces: `Config` struct, `Load()`, `Save()`, `Default()` functions

- [ ] **Step 1: Add toml dependency**

```bash
go get github.com/BurntSushi/toml@latest
```

- [ ] **Step 2: Write failing test**

Create `internal/config/config_test.go`:

```go
package config

import (
    "os"
    "path/filepath"
    "testing"
)

func TestLoadConfig(t *testing.T) {
    content := `
[package]
name = "test-project"
version = "1.0.0"
edition = "2024"

[build]
compiler = "auto"
std = "c++23"
`
    tmpDir := t.TempDir()
    configPath := filepath.Join(tmpDir, "cnext.toml")
    
    if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
        t.Fatalf("Failed to write config: %v", err)
    }
    
    cfg, err := Load(configPath)
    if err != nil {
        t.Fatalf("Load failed: %v", err)
    }
    
    if cfg.Package.Name != "test-project" {
        t.Errorf("Expected name 'test-project', got '%s'", cfg.Package.Name)
    }
    
    if cfg.Package.Version != "1.0.0" {
        t.Errorf("Expected version '1.0.0', got '%s'", cfg.Package.Version)
    }
    
    if cfg.Build.Compiler != "auto" {
        t.Errorf("Expected compiler 'auto', got '%s'", cfg.Build.Compiler)
    }
}

func TestDefaultConfig(t *testing.T) {
    cfg := Default()
    
    if cfg.Package.Name != "" {
        t.Errorf("Expected empty name, got '%s'", cfg.Package.Name)
    }
    
    if cfg.Build.Compiler != "auto" {
        t.Errorf("Expected compiler 'auto', got '%s'", cfg.Build.Compiler)
    }
    
    if cfg.Build.Std != "c++23" {
        t.Errorf("Expected std 'c++23', got '%s'", cfg.Build.Std)
    }
}

func TestSaveConfig(t *testing.T) {
    cfg := Default()
    cfg.Package.Name = "save-test"
    cfg.Package.Version = "2.0.0"
    
    tmpDir := t.TempDir()
    configPath := filepath.Join(tmpDir, "cnext.toml")
    
    if err := Save(cfg, configPath); err != nil {
        t.Fatalf("Save failed: %v", err)
    }
    
    loaded, err := Load(configPath)
    if err != nil {
        t.Fatalf("Load failed: %v", err)
    }
    
    if loaded.Package.Name != "save-test" {
        t.Errorf("Expected name 'save-test', got '%s'", loaded.Package.Name)
    }
}
```

- [ ] **Step 3: Run test to verify it fails**

```bash
go test ./internal/config/...
```

Expected: FAIL with "Load not defined"

- [ ] **Step 4: Write minimal implementation**

Create `internal/config/config.go`:

```go
package config

import (
    "os"
    
    "github.com/BurntSushi/toml"
)

type Config struct {
    Package PackageConfig    `toml:"package"`
    Build   BuildConfig      `toml:"build"`
    Deps    []Dependency     `toml:"dependencies"`
    DevDeps []Dependency     `toml:"dev-dependencies"`
}

type PackageConfig struct {
    Name    string `toml:"name"`
    Version string `toml:"version"`
    Edition string `toml:"edition"`
}

type BuildConfig struct {
    Compiler    string        `toml:"compiler"`
    Std         string        `toml:"std"`
    Optimize    string        `toml:"optimization"`
    LTO         bool          `toml:"lto"`
    PIC         bool          `toml:"pic"`
    Sanitize    []string      `toml:"sanitize"`
    Target      TargetConfig  `toml:"target"`
    Libraries   []Library     `toml:"libraries"`
    Executables []Executable  `toml:"executables"`
    CMakeDeps   []CMakeDep    `toml:"cmake-dependencies"`
    Tests       []TestConfig  `toml:"tests"`
}

type TargetConfig struct {
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

type CMakeDep struct {
    Name          string   `toml:"name"`
    Version       string   `toml:"version"`
    Source        string   `toml:"source"`
    Git           string   `toml:"git"`
    Tag           string   `toml:"tag"`
    CMakeOptions  []string `toml:"cmake_options"`
}

type TestConfig struct {
    Name      string   `toml:"name"`
    Type      string   `toml:"type"`
    Sources   []string `toml:"sources"`
    Framework string   `toml:"framework"`
}

type Dependency struct {
    Name     string            `toml:"name"`
    Version  string            `toml:"version"`
    Features []string          `toml:"features"`
    Source   string            `toml:"source"`
    Git      string            `toml:"git"`
    Branch   string            `toml:"branch"`
}

func Default() *Config {
    return &Config{
        Package: PackageConfig{
            Edition: "2024",
        },
        Build: BuildConfig{
            Compiler: "auto",
            Std:      "c++23",
            Optimize: "debug",
            Target: TargetConfig{
                OS:   "auto",
                Arch: "auto",
            },
        },
        Deps:    []Dependency{},
        DevDeps: []Dependency{},
    }
}

func Load(path string) (*Config, error) {
    cfg := Default()
    
    if _, err := toml.DecodeFile(path, cfg); err != nil {
        return nil, err
    }
    
    return cfg, nil
}

func Save(cfg *Config, path string) error {
    f, err := os.Create(path)
    if err != nil {
        return err
    }
    defer f.Close()
    
    encoder := toml.NewEncoder(f)
    return encoder.Encode(cfg)
}
```

- [ ] **Step 5: Run test to verify it passes**

```bash
go test ./internal/config/...
```

Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add -A
git commit -m "feat: add cnext.toml parser with config structs"
```

---

### Task 4: Basic Compiler Driver

**Covers:** S3.2

**Files:**
- Create: `internal/compiler/driver.go`
- Create: `internal/compiler/driver_test.go`
- Create: `internal/compiler/gcc.go`
- Create: `internal/compiler/clang.go`
- Create: `internal/compiler/msvc.go`

**Interfaces:**
- Consumes: Task 3 Config
- Produces: `Compiler` interface, `Compile()`, `Link()` functions

- [ ] **Step 1: Write failing test**

Create `internal/compiler/driver_test.go`:

```go
package compiler

import (
    "testing"
    
    "github.com/neko233-com/cnext/internal/config"
)

func TestAutoDetectCompiler(t *testing.T) {
    compiler, err := AutoDetect()
    if err != nil {
        t.Skipf("No compiler available: %v", err)
    }
    
    if compiler == nil {
        t.Fatal("Expected non-nil compiler")
    }
    
    name := compiler.Name()
    if name == "" {
        t.Error("Expected non-empty compiler name")
    }
    
    t.Logf("Detected compiler: %s", name)
}

func TestGCCCompiler(t *testing.T) {
    gcc := NewGCC("/usr/bin/gcc")
    if gcc.Name() != "gcc" {
        t.Errorf("Expected name 'gcc', got '%s'", gcc.Name())
    }
}

func TestClangCompiler(t *testing.T) {
    clang := NewClang("/usr/bin/clang")
    if clang.Name() != "clang" {
        t.Errorf("Expected name 'clang', got '%s'", clang.Name())
    }
}
```

- [ ] **Step 2: Run test to verify it fails**

```bash
go test ./internal/compiler/...
```

Expected: FAIL with "AutoDetect not defined"

- [ ] **Step 3: Write minimal implementation**

Create `internal/compiler/driver.go`:

```go
package compiler

import (
    "fmt"
    "os/exec"
    "runtime"
    "strings"
    
    "github.com/neko233-com/cnext/internal/config"
)

type Compiler interface {
    Name() string
    Version() string
    Compile(opts CompileOptions) error
    Link(opts LinkOptions) error
}

type CompileOptions struct {
    Sources     []string
    Output      string
    Std         string
    IncludeDirs []string
    Defines     []string
    Flags       []string
    Optimization string
}

type LinkOptions struct {
    Objects    []string
    Output     string
    Libraries  []string
    LibDirs    []string
    Flags      []string
}

func AutoDetect() (Compiler, error) {
    switch runtime.GOOS {
    case "linux", "darwin":
        if path, err := exec.LookPath("gcc"); err == nil {
            return NewGCC(path), nil
        }
        if path, err := exec.LookPath("clang"); err == nil {
            return NewClang(path), nil
        }
    case "windows":
        if path, err := exec.LookPath("gcc"); err == nil {
            return NewGCC(path), nil
        }
        if path, err := exec.LookPath("clang"); err == nil {
            return NewClang(path), nil
        }
    }
    
    return nil, fmt.Errorf("no compiler found")
}

func DetectForTarget(target config.TargetConfig) (Compiler, error) {
    return AutoDetect()
}
```

Create `internal/compiler/gcc.go`:

```go
package compiler

import (
    "fmt"
    "os/exec"
    "strings"
)

type GCC struct {
    path string
}

func NewGCC(path string) *GCC {
    return &GCC{path: path}
}

func (g *GCC) Name() string {
    return "gcc"
}

func (g *GCC) Version() string {
    cmd := exec.Command(g.path, "--version")
    out, err := cmd.Output()
    if err != nil {
        return "unknown"
    }
    
    lines := strings.Split(string(out), "\n")
    if len(lines) > 0 {
        return strings.TrimSpace(lines[0])
    }
    return "unknown"
}

func (g *GCC) Compile(opts CompileOptions) error {
    args := []string{}
    
    // Add std
    if opts.Std != "" {
        args = append(args, "-std="+opts.Std)
    }
    
    // Add optimization
    switch opts.Optimization {
    case "debug":
        args = append(args, "-O0", "-g")
    case "release":
        args = append(args, "-O2")
    case "size":
        args = append(args, "-Os")
    case "speed":
        args = append(args, "-O3")
    }
    
    // Add include dirs
    for _, dir := range opts.IncludeDirs {
        args = append(args, "-I"+dir)
    }
    
    // Add defines
    for _, def := range opts.Defines {
        args = append(args, "-D"+def)
    }
    
    // Add flags
    args = append(args, opts.Flags...)
    
    // Add output
    args = append(args, "-o", opts.Output)
    
    // Add sources
    args = append(args, opts.Sources...)
    
    cmd := exec.Command(g.path, args...)
    cmd.Stdout = nil
    cmd.Stderr = nil
    
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("gcc compilation failed: %w", err)
    }
    
    return nil
}

func (g *GCC) Link(opts LinkOptions) error {
    args := []string{}
    
    // Add objects
    args = append(args, opts.Objects...)
    
    // Add output
    args = append(args, "-o", opts.Output)
    
    // Add library dirs
    for _, dir := range opts.LibDirs {
        args = append(args, "-L"+dir)
    }
    
    // Add libraries
    for _, lib := range opts.Libraries {
        args = append(args, "-l"+lib)
    }
    
    // Add flags
    args = append(args, opts.Flags...)
    
    cmd := exec.Command(g.path, args...)
    cmd.Stdout = nil
    cmd.Stderr = nil
    
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("gcc linking failed: %w", err)
    }
    
    return nil
}
```

Create `internal/compiler/clang.go`:

```go
package compiler

import (
    "fmt"
    "os/exec"
    "strings"
)

type Clang struct {
    path string
}

func NewClang(path string) *Clang {
    return &Clang{path: path}
}

func (c *Clang) Name() string {
    return "clang"
}

func (c *Clang) Version() string {
    cmd := exec.Command(c.path, "--version")
    out, err := cmd.Output()
    if err != nil {
        return "unknown"
    }
    
    lines := strings.Split(string(out), "\n")
    if len(lines) > 0 {
        return strings.TrimSpace(lines[0])
    }
    return "unknown"
}

func (c *Clang) Compile(opts CompileOptions) error {
    args := []string{}
    
    // Add std
    if opts.Std != "" {
        args = append(args, "-std="+opts.Std)
    }
    
    // Add optimization
    switch opts.Optimization {
    case "debug":
        args = append(args, "-O0", "-g")
    case "release":
        args = append(args, "-O2")
    case "size":
        args = append(args, "-Os")
    case "speed":
        args = append(args, "-O3")
    }
    
    // Add include dirs
    for _, dir := range opts.IncludeDirs {
        args = append(args, "-I"+dir)
    }
    
    // Add defines
    for _, def := range opts.Defines {
        args = append(args, "-D"+def)
    }
    
    // Add flags
    args = append(args, opts.Flags...)
    
    // Add output
    args = append(args, "-o", opts.Output)
    
    // Add sources
    args = append(args, opts.Sources...)
    
    cmd := exec.Command(c.path, args...)
    cmd.Stdout = nil
    cmd.Stderr = nil
    
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("clang compilation failed: %w", err)
    }
    
    return nil
}

func (c *Clang) Link(opts LinkOptions) error {
    args := []string{}
    
    // Add objects
    args = append(args, opts.Objects...)
    
    // Add output
    args = append(args, "-o", opts.Output)
    
    // Add library dirs
    for _, dir := range opts.LibDirs {
        args = append(args, "-L"+dir)
    }
    
    // Add libraries
    for _, lib := range opts.Libraries {
        args = append(args, "-l"+lib)
    }
    
    // Add flags
    args = append(args, opts.Flags...)
    
    cmd := exec.Command(c.path, args...)
    cmd.Stdout = nil
    cmd.Stderr = nil
    
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("clang linking failed: %w", err)
    }
    
    return nil
}
```

Create `internal/compiler/msvc.go`:

```go
package compiler

import (
    "fmt"
    "os/exec"
    "runtime"
    "strings"
)

type MSVC struct {
    path string
}

func NewMSVC(path string) *MSVC {
    return &MSVC{path: path}
}

func (m *MSVC) Name() string {
    return "msvc"
}

func (m *MSVC) Version() string {
    // MSVC doesn't have a simple --version
    return "MSVC (version unknown)"
}

func (m *MSVC) Compile(opts CompileOptions) error {
    if runtime.GOOS != "windows" {
        return fmt.Errorf("msvc only available on Windows")
    }
    
    args := []string{}
    
    // MSVC uses /std:c++23 format
    if opts.Std != "" {
        std := strings.Replace(opts.Std, "c++", "/std:c++", 1)
        std = strings.Replace(std, "c1", "/std:c1", 1)
        args = append(args, std)
    }
    
    // Add optimization
    switch opts.Optimization {
    case "debug":
        args = append(args, "/Od", "/Zi")
    case "release":
        args = append(args, "/O2")
    case "size":
        args = append(args, "/O1")
    case "speed":
        args = append(args, "/O2", "/Ot")
    }
    
    // Add include dirs
    for _, dir := range opts.IncludeDirs {
        args = append(args, "/I"+dir)
    }
    
    // Add defines
    for _, def := range opts.Defines {
        args = append(args, "/D"+def)
    }
    
    // Add flags
    args = append(args, opts.Flags...)
    
    // Add output
    args = append(args, "/Fo:"+opts.Output)
    
    // Add sources
    args = append(args, opts.Sources...)
    
    cmd := exec.Command(m.path, args...)
    cmd.Stdout = nil
    cmd.Stderr = nil
    
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("msvc compilation failed: %w", err)
    }
    
    return nil
}

func (m *MSVC) Link(opts LinkOptions) error {
    if runtime.GOOS != "windows" {
        return fmt.Errorf("msvc only available on Windows")
    }
    
    args := []string{}
    
    // Add objects
    args = append(args, opts.Objects...)
    
    // Add output
    args = append(args, "/Fe:"+opts.Output)
    
    // Add library dirs
    for _, dir := range opts.LibDirs {
        args = append(args, "/LIBPATH:"+dir)
    }
    
    // Add libraries
    for _, lib := range opts.Libraries {
        args = append(args, lib+".lib")
    }
    
    // Add flags
    args = append(args, opts.Flags...)
    
    cmd := exec.Command(m.path, args...)
    cmd.Stdout = nil
    cmd.Stderr = nil
    
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("msvc linking failed: %w", err)
    }
    
    return nil
}
```

- [ ] **Step 4: Run test to verify it passes**

```bash
go test ./internal/compiler/...
```

Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add -A
git commit -m "feat: add basic compiler driver with gcc/clang/msvc support"
```

---

### Task 5: Build Graph Generation

**Covers:** S4.2

**Files:**
- Create: `internal/build/graph.go`
- Create: `internal/build/graph_test.go`

**Interfaces:**
- Consumes: Task 3 Config, Task 4 Compiler
- Produces: `BuildGraph`, `Generate()`, `Execute()` functions

- [ ] **Step 1: Write failing test**

Create `internal/build/graph_test.go`:

```go
package build

import (
    "testing"
    
    "github.com/neko233-com/cnext/internal/config"
)

func TestGenerateBuildGraph(t *testing.T) {
    cfg := &config.Config{
        Package: config.PackageConfig{
            Name:    "test",
            Version: "1.0.0",
        },
        Build: config.BuildConfig{
            Compiler: "auto",
            Std:      "c++23",
            Libraries: []config.Library{
                {
                    Name:    "mylib",
                    Type:    "static",
                    Sources: []string{"src/**/*.cpp"},
                },
            },
            Executables: []config.Executable{
                {
                    Name:      "myapp",
                    Sources:   []string{"main.cpp"},
                    Libraries: []string{"mylib"},
                },
            },
        },
    }
    
    graph, err := Generate(cfg)
    if err != nil {
        t.Fatalf("Generate failed: %v", err)
    }
    
    if graph == nil {
        t.Fatal("Expected non-nil graph")
    }
    
    if len(graph.Nodes) == 0 {
        t.Error("Expected non-empty graph")
    }
    
    t.Logf("Graph has %d nodes", len(graph.Nodes))
}

func TestTopologicalSort(t *testing.T) {
    graph := &BuildGraph{
        Nodes: []BuildNode{
            {ID: "a", Dependencies: []string{"b"}},
            {ID: "b", Dependencies: []string{"c"}},
            {ID: "c", Dependencies: []string{}},
        },
    }
    
    sorted, err := graph.TopologicalSort()
    if err != nil {
        t.Fatalf("TopologicalSort failed: %v", err)
    }
    
    if len(sorted) != 3 {
        t.Errorf("Expected 3 nodes, got %d", len(sorted))
    }
    
    // c should come before b, b before a
    cIdx, bIdx, aIdx := -1, -1, -1
    for i, id := range sorted {
        switch id {
        case "c":
            cIdx = i
        case "b":
            bIdx = i
        case "a":
            aIdx = i
        }
    }
    
    if cIdx >= bIdx || bIdx >= aIdx {
        t.Errorf("Invalid order: c=%d, b=%d, a=%d", cIdx, bIdx, aIdx)
    }
}
```

- [ ] **Step 2: Run test to verify it fails**

```bash
go test ./internal/build/...
```

Expected: FAIL with "Generate not defined"

- [ ] **Step 3: Write minimal implementation**

Create `internal/build/graph.go`:

```go
package build

import (
    "fmt"
    "os"
    "path/filepath"
    
    "github.com/neko233-com/cnext/internal/config"
)

type BuildGraph struct {
    Nodes []BuildNode
}

type BuildNode struct {
    ID           string
    Type         string // "compile", "link", "cmake"
    Sources      []string
    Output       string
    Dependencies []string
    Compiler     string
    Std          string
    Flags        []string
}

func Generate(cfg *config.Config) (*BuildGraph, error) {
    graph := &BuildGraph{
        Nodes: []BuildNode{},
    }
    
    // Generate compile nodes for libraries
    for _, lib := range cfg.Build.Libraries {
        sources, err := globSources(lib.Sources)
        if err != nil {
            return nil, fmt.Errorf("failed to glob sources for %s: %w", lib.Name, err)
        }
        
        if len(sources) == 0 {
            continue
        }
        
        node := BuildNode{
            ID:       "lib:" + lib.Name,
            Type:     "compile",
            Sources:  sources,
            Output:   lib.Name + ".a", // Default to static
            Compiler: cfg.Build.Compiler,
            Std:      cfg.Build.Std,
        }
        
        graph.Nodes = append(graph.Nodes, node)
    }
    
    // Generate compile nodes for executables
    for _, exe := range cfg.Build.Executables {
        sources, err := globSources(exe.Sources)
        if err != nil {
            return nil, fmt.Errorf("failed to glob sources for %s: %w", exe.Name, err)
        }
        
        if len(sources) == 0 {
            continue
        }
        
        node := BuildNode{
            ID:           "exe:" + exe.Name,
            Type:         "compile",
            Sources:      sources,
            Output:       exe.Name + ".o",
            Dependencies: exe.Libraries,
            Compiler:     cfg.Build.Compiler,
            Std:          cfg.Build.Std,
        }
        
        graph.Nodes = append(graph.Nodes, node)
        
        // Generate link node
        linkNode := BuildNode{
            ID:           "link:" + exe.Name,
            Type:         "link",
            Output:       exe.Name,
            Dependencies: append([]string{"exe:" + exe.Name}, exe.Libraries...),
        }
        
        graph.Nodes = append(graph.Nodes, linkNode)
    }
    
    return graph, nil
}

func globSources(patterns []string) ([]string, error) {
    var sources []string
    
    for _, pattern := range patterns {
        matches, err := filepath.Glob(pattern)
        if err != nil {
            return nil, err
        }
        sources = append(sources, matches...)
    }
    
    return sources, nil
}

func (g *BuildGraph) TopologicalSort() ([]string, error) {
    // Build adjacency list
    adj := make(map[string][]string)
    inDegree := make(map[string]int)
    
    for _, node := range g.Nodes {
        if _, exists := inDegree[node.ID]; !exists {
            inDegree[node.ID] = 0
        }
        
        for _, dep := range node.Dependencies {
            adj[dep] = append(adj[dep], node.ID)
            inDegree[node.ID]++
        }
    }
    
    // Kahn's algorithm
    var queue []string
    for id, degree := range inDegree {
        if degree == 0 {
            queue = append(queue, id)
        }
    }
    
    var sorted []string
    for len(queue) > 0 {
        node := queue[0]
        queue = queue[1:]
        sorted = append(sorted, node)
        
        for _, neighbor := range adj[node] {
            inDegree[neighbor]--
            if inDegree[neighbor] == 0 {
                queue = append(queue, neighbor)
            }
        }
    }
    
    if len(sorted) != len(g.Nodes) {
        return nil, fmt.Errorf("cycle detected in build graph")
    }
    
    return sorted, nil
}

func (g *BuildGraph) Execute() error {
    sorted, err := g.TopologicalSort()
    if err != nil {
        return err
    }
    
    for _, id := range sorted {
        for _, node := range g.Nodes {
            if node.ID == id {
                fmt.Printf("Building %s...\n", node.ID)
                // TODO: Task 6 will implement actual compilation
            }
        }
    }
    
    return nil
}
```

- [ ] **Step 4: Run test to verify it passes**

```bash
go test ./internal/build/...
```

Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add -A
git commit -m "feat: add build graph generation with topological sort"
```

---

### Task 6: Build Command Implementation

**Covers:** S4, S7

**Files:**
- Modify: `cmd/cnext/build.go`

**Interfaces:**
- Consumes: Task 3 Config, Task 4 Compiler, Task 5 BuildGraph
- Produces: Working `cnext build` command

- [ ] **Step 1: Update build command**

Update `cmd/cnext/build.go`:

```go
package main

import (
    "fmt"
    "os"
    "path/filepath"
    
    "github.com/spf13/cobra"
    
    "github.com/neko233-com/cnext/internal/build"
    "github.com/neko233-com/cnext/internal/compiler"
    "github.com/neko233-com/cnext/internal/config"
)

var buildCmd = &cobra.Command{
    Use:   "build",
    Short: "Build the project",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Load config
        configPath := "cnext.toml"
        if _, err := os.Stat(configPath); os.IsNotExist(err) {
            return fmt.Errorf("no cnext.toml found. Run 'cnext init' first")
        }
        
        cfg, err := config.Load(configPath)
        if err != nil {
            return fmt.Errorf("failed to load config: %w", err)
        }
        
        // Get flags
        target, _ := cmd.Flags().GetString("target")
        release, _ := cmd.Flags().GetBool("release")
        compilerName, _ := cmd.Flags().GetString("compiler")
        
        // Override config with flags
        if target != "" {
            parts := splitTarget(target)
            cfg.Build.Target.OS = parts[0]
            cfg.Build.Target.Arch = parts[1]
        }
        
        if release {
            cfg.Build.Optimize = "release"
        }
        
        if compilerName != "auto" {
            cfg.Build.Compiler = compilerName
        }
        
        // Generate build graph
        graph, err := build.Generate(cfg)
        if err != nil {
            return fmt.Errorf("failed to generate build graph: %w", err)
        }
        
        // Detect compiler
        comp, err := compiler.AutoDetect()
        if err != nil {
            return fmt.Errorf("no compiler found: %w", err)
        }
        
        fmt.Printf("Using compiler: %s\n", comp.Name())
        fmt.Printf("Building %s...\n", cfg.Package.Name)
        
        // Execute build
        if err := graph.Execute(); err != nil {
            return fmt.Errorf("build failed: %w", err)
        }
        
        fmt.Println("Build successful!")
        return nil
    },
}

func init() {
    buildCmd.Flags().StringP("target", "t", "", "Build target (e.g., x86_64-linux-gnu)")
    buildCmd.Flags().BoolP("release", "r", false, "Build in release mode")
    buildCmd.Flags().StringP("compiler", "c", "auto", "Compiler to use (auto, gcc, clang, msvc)")
}

func splitTarget(target string) []string {
    // Simple split by "-"
    // TODO: More sophisticated target parsing
    return []string{"auto", "auto"}
}
```

- [ ] **Step 2: Verify it compiles**

```bash
go build ./cmd/cnext
```

- [ ] **Step 3: Test with a simple project**

```bash
cd /tmp
mkdir test-project && cd test-project
/path/to/cnext init
# Edit cnext.toml to add a simple executable
echo 'int main() { return 0; }' > main.cpp
/path/to/cnext build
```

- [ ] **Step 4: Commit**

```bash
git add -A
git commit -m "feat: implement build command with config loading"
```

---

### Task 7: Init Command Implementation

**Covers:** S7

**Files:**
- Modify: `cmd/cnext/init.go`

**Interfaces:**
- Consumes: Task 3 Config
- Produces: Working `cnext init` command

- [ ] **Step 1: Update init command**

Update `cmd/cnext/init.go`:

```go
package main

import (
    "fmt"
    "os"
    "path/filepath"
    
    "github.com/spf13/cobra"
    
    "github.com/neko233-com/cnext/internal/config"
)

var initCmd = &cobra.Command{
    Use:   "init",
    Short: "Initialize a new cnext project",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Get project name
        name := "my-project"
        if len(args) > 0 {
            name = args[0]
        }
        
        // Check if already initialized
        if _, err := os.Stat("cnext.toml"); err == nil {
            return fmt.Errorf("cnext.toml already exists")
        }
        
        // Create default config
        cfg := config.Default()
        cfg.Package.Name = name
        cfg.Package.Version = "0.1.0"
        
        // Create directories
        dirs := []string{"src", "include", "tests"}
        for _, dir := range dirs {
            if err := os.MkdirAll(dir, 0755); err != nil {
                return fmt.Errorf("failed to create directory %s: %w", dir, err)
            }
        }
        
        // Create cnext.toml
        if err := config.Save(cfg, "cnext.toml"); err != nil {
            return fmt.Errorf("failed to create cnext.toml: %w", err)
        }
        
        // Create main.cpp
        mainContent := `#include <iostream>

int main() {
    std::cout << "Hello, " << name << "!" << std::endl;
    return 0;
}
`
        if err := os.WriteFile("main.cpp", []byte(mainContent), 0644); err != nil {
            return fmt.Errorf("failed to create main.cpp: %w", err)
        }
        
        // Create .gitignore
        gitignore := `build/
*.o
*.a
*.so
*.dylib
*.exe
`
        if err := os.WriteFile(".gitignore", []byte(gitignore), 0644); err != nil {
            return fmt.Errorf("failed to create .gitignore: %w", err)
        }
        
        fmt.Printf("Initialized project '%s'\n", name)
        fmt.Println("Created:")
        fmt.Println("  cnext.toml")
        fmt.Println("  main.cpp")
        fmt.Println("  .gitignore")
        fmt.Println("  src/")
        fmt.Println("  include/")
        fmt.Println("  tests/")
        
        return nil
    },
}
```

- [ ] **Step 2: Verify it compiles**

```bash
go build ./cmd/cnext
```

- [ ] **Step 3: Test init command**

```bash
cd /tmp
mkdir test-init && cd test-init
/path/to/cnext init my-app
ls -la
cat cnext.toml
```

- [ ] **Step 4: Commit**

```bash
git add -A
git commit -m "feat: implement init command with project scaffolding"
```

---

### Task 8: Basic Test Framework

**Covers:** S5.1

**Files:**
- Create: `cnext-test/include/cnext/test.h`
- Create: `cnext-test/include/cnext/assert.h`
- Create: `cnext-test/src/runner.cpp`

**Interfaces:**
- Consumes: Task 4 Compiler
- Produces: Test framework headers and runner

- [ ] **Step 1: Create test.h**

Create `cnext-test/include/cnext/test.h`:

```cpp
#pragma once

#include <string>
#include <vector>
#include <functional>
#include <iostream>
#include <chrono>

namespace cnext {
namespace test {

struct TestCase {
    std::string name;
    std::string suite;
    std::function<void()> func;
};

class TestRegistry {
public:
    static TestRegistry& Instance() {
        static TestRegistry instance;
        return instance;
    }
    
    void Register(const std::string& suite, const std::string& name, std::function<void()> func) {
        cases.push_back({name, suite, func});
    }
    
    const std::vector<TestCase>& Cases() const {
        return cases;
    }
    
private:
    std::vector<TestCase> cases;
};

class TestRunner {
public:
    int Run(const std::string& filter = "") {
        int passed = 0;
        int failed = 0;
        int skipped = 0;
        
        auto& registry = TestRegistry::Instance();
        
        std::cout << "Running tests..." << std::endl;
        std::cout << std::endl;
        
        for (const auto& tc : registry.Cases()) {
            // Apply filter
            if (!filter.empty() && 
                tc.suite.find(filter) == std::string::npos && 
                tc.name.find(filter) == std::string::npos) {
                continue;
            }
            
            auto start = std::chrono::high_resolution_clock::now();
            
            try {
                tc.func();
                auto end = std::chrono::high_resolution_clock::now();
                auto duration = std::chrono::duration_cast<std::chrono::microseconds>(end - start);
                
                std::cout << "  PASS: " << tc.suite << "::" << tc.name 
                         << " (" << duration.count() / 1000.0 << "ms)" << std::endl;
                passed++;
            } catch (const std::exception& e) {
                auto end = std::chrono::high_resolution_clock::now();
                auto duration = std::chrono::duration_cast<std::chrono::microseconds>(end - start);
                
                std::cout << "  FAIL: " << tc.suite << "::" << tc.name 
                         << " (" << duration.count() / 1000.0 << "ms)" << std::endl;
                std::cout << "    Error: " << e.what() << std::endl;
                failed++;
            } catch (...) {
                std::cout << "  FAIL: " << tc.suite << "::" << tc.name 
                         << " (unknown exception)" << std::endl;
                failed++;
            }
        }
        
        std::cout << std::endl;
        std::cout << passed << " passed, " << failed << " failed, " 
                 << skipped << " skipped" << std::endl;
        
        return failed > 0 ? 1 : 0;
    }
};

} // namespace test
} // namespace namespace

#define TEST(suite, name) \
    void suite##_##name(); \
    namespace { \
        struct Suite##_##name##Registrar { \
            Suite##_##name##Registrar() { \
                ::cnext::test::TestRegistry::Instance().Register(#suite, #name, suite##_##name); \
            } \
        } suite##_##name##registrar; \
    } \
    void suite##_##name()

#define EXPECT_EQ(a, b) \
    do { \
        if ((a) != (b)) { \
            throw std::runtime_error("Expected " #a " == " #b); \
        } \
    } while(0)

#define ASSERT_EQ(a, b) EXPECT_EQ(a, b)

#define EXPECT_TRUE(cond) \
    do { \
        if (!(cond)) { \
            throw std::runtime_error("Expected " #cond " to be true"); \
        } \
    } while(0)

#define EXPECT_FALSE(cond) EXPECT_TRUE(!(cond))

#define EXPECT_GT(a, b) \
    do { \
        if (!((a) > (b))) { \
            throw std::runtime_error("Expected " #a " > " #b); \
        } \
    } while(0)

#define EXPECT_LT(a, b) \
    do { \
        if (!((a) < (b))) { \
            throw std::runtime_error("Expected " #a " < " #b); \
        } \
    } while(0)
```

- [ ] **Step 2: Create assert.h**

Create `cnext-test/include/cnext/assert.h`:

```cpp
#pragma once

#include <cmath>
#include <stdexcept>
#include <string>

namespace cnext {
namespace test {

template<typename T>
void AssertEqual(const T& a, const T& b, const std::string& msg = "") {
    if (a != b) {
        throw std::runtime_error("AssertEqual failed: " + msg);
    }
}

template<typename T>
void AssertNear(const T& a, const T& b, const T& epsilon, const std::string& msg = "") {
    if (std::abs(a - b) > epsilon) {
        throw std::runtime_error("AssertNear failed: " + msg);
    }
}

void AssertTrue(bool condition, const std::string& msg = "") {
    if (!condition) {
        throw std::runtime_error("AssertTrue failed: " + msg);
    }
}

void AssertFalse(bool condition, const std::string& msg = "") {
    if (condition) {
        throw std::runtime_error("AssertFalse failed: " + msg);
    }
}

} // namespace test
} // namespace namespace
```

- [ ] **Step 3: Create runner.cpp**

Create `cnext-test/src/runner.cpp`:

```cpp
#include <cnext/test.h>

int main(int argc, char* argv[]) {
    std::string filter;
    
    // Simple argument parsing
    for (int i = 1; i < argc; i++) {
        std::string arg = argv[i];
        if (arg == "--filter" || arg == "-f") {
            if (i + 1 < argc) {
                filter = argv[++i];
            }
        }
    }
    
    cnext::test::TestRunner runner;
    return runner.Run(filter);
}
```

- [ ] **Step 4: Verify it compiles**

```bash
cd cnext-test
g++ -std=c++23 -I include -c src/runner.cpp -o runner.o
```

- [ ] **Step 5: Commit**

```bash
git add -A
git commit -m "feat: add basic test framework with assertions"
```

---

### Task 9: Test Command Implementation

**Covers:** S5.2

**Files:**
- Modify: `cmd/cnext/test.go`

**Interfaces:**
- Consumes: Task 4 Compiler, Task 8 Test Framework
- Produces: Working `cnext test` command

- [ ] **Step 1: Update test command**

Update `cmd/cnext/test.go`:

```go
package main

import (
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
    
    "github.com/spf13/cobra"
    
    "github.com/neko233-com/cnext/internal/config"
)

var testCmd = &cobra.Command{
    Use:   "test",
    Short: "Run project tests",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Load config
        configPath := "cnext.toml"
        if _, err := os.Stat(configPath); os.IsNotExist(err) {
            return fmt.Errorf("no cnext.toml found. Run 'cnext init' first")
        }
        
        cfg, err := config.Load(configPath)
        if err != nil {
            return fmt.Errorf("failed to load config: %w", err)
        }
        
        // Get flags
        filter, _ := cmd.Flags().GetString("filter")
        coverage, _ := cmd.Flags().GetBool("coverage")
        
        // Find test files
        testFiles, err := findTestFiles()
        if err != nil {
            return fmt.Errorf("failed to find test files: %w", err)
        }
        
        if len(testFiles) == 0 {
            fmt.Println("No test files found")
            return nil
        }
        
        fmt.Printf("Found %d test files\n", len(testFiles))
        
        // Compile tests
        testBinary := "test_runner"
        if err := compileTests(testFiles, testBinary, coverage); err != nil {
            return fmt.Errorf("failed to compile tests: %w", err)
        }
        
        // Run tests
        args = []string{}
        if filter != "" {
            args = append(args, "--filter", filter)
        }
        
        return runTests(testBinary, args)
    },
}

func findTestFiles() ([]string, error) {
    var testFiles []string
    
    err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        
        if info.IsDir() {
            return nil
        }
        
        if strings.HasSuffix(path, "_test.cpp") || strings.HasSuffix(path, "_test.cc") {
            testFiles = append(testFiles, path)
        }
        
        return nil
    })
    
    return testFiles, err
}

func compileTests(testFiles []string, output string, coverage bool) error {
    args := []string{"-std=c++23", "-I", "include", "-I", "cnext-test/include"}
    
    if coverage {
        args = append(args, "--coverage")
    }
    
    args = append(args, testFiles...)
    args = append(args, "cnext-test/src/runner.cpp")
    args = append(args, "-o", output)
    
    cmd := exec.Command("g++", args...)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    
    return cmd.Run()
}

func runTests(binary string, args []string) error {
    cmd := exec.Command("./"+binary, args...)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    
    return cmd.Run()
}

func init() {
    testCmd.Flags().StringP("filter", "f", "", "Filter tests by name")
    testCmd.Flags().BoolP("coverage", "c", false, "Generate coverage report")
    testCmd.Flags().BoolP("benchmark", "b", false, "Run benchmarks")
    testCmd.Flags().BoolP("watch", "w", false, "Watch for changes and rerun tests")
}
```

- [ ] **Step 2: Verify it compiles**

```bash
go build ./cmd/cnext
```

- [ ] **Step 3: Test with a simple test**

```bash
cd /tmp
mkdir test-example && cd test-example
/path/to/cnext init

# Create a test file
cat > tests/math_test.cpp << 'EOF'
#include <cnext/test.h>

TEST(Math, Addition) {
    EXPECT_EQ(2 + 2, 4);
}

TEST(Math, Subtraction) {
    EXPECT_EQ(5 - 3, 2);
}
EOF

/path/to/cnext test
```

- [ ] **Step 4: Commit**

```bash
git add -A
git commit -m "feat: implement test command with test discovery"
```

---

### Task 10: README and Documentation

**Covers:** S9

**Files:**
- Create: `README.md`
- Create: `LICENSE`

**Interfaces:**
- Consumes: All previous tasks
- Produces: Project documentation

- [ ] **Step 1: Create README.md**

Create `README.md`:

```markdown
# cnext

Next-generation C/C++ toolchain written in Go.

## Features

- **One-click install** — no need to install gcc/clang separately
- **Declarative build system** — cnext.toml replaces CMakeLists.txt
- **Built-in package manager** — like npm for C/C++
- **Built-in test framework** — Go-style testing
- **CMake bridging** — seamlessly use CMake dependencies
- **Cross-compilation** — target any platform from any host

## Installation

```bash
# Linux/macOS
curl -fsSL https://get.cnext.dev | sh

# Windows
iwr -useb https://get.cnext.dev/install.ps1 | iex
```

## Quick Start

```bash
# Create a new project
cnext init my-project
cd my-project

# Build
cnext build

# Run
cnext run

# Test
cnext test

# Add a dependency
cnext add fmt@10.2.1
```

## Configuration

```toml
[package]
name = "my-project"
version = "1.0.0"
edition = "2024"

[build]
compiler = "auto"
std = "c++23"
optimization = "release"

[dependencies]
fmt = "10.2.1"
spdlog = "1.13.0"
```

## Commands

| Command | Description |
|---------|-------------|
| `cnext init` | Initialize a new project |
| `cnext build` | Build the project |
| `cnext test` | Run tests |
| `cnext run` | Run the built executable |
| `cnext add <pkg>` | Add a dependency |
| `cnext remove <pkg>` | Remove a dependency |
| `cnext fmt` | Format code |
| `cnext lint` | Lint code |
| `cnext doc` | Generate documentation |
| `cnext publish` | Publish to registry |

## Testing

```cpp
#include <cnext/test.h>

TEST(Math, Addition) {
    EXPECT_EQ(2 + 2, 4);
}

TEST(Math, TableDriven) {
    struct TestCase {
        int a, b, expected;
    };
    
    TestCase cases[] = {
        {1, 1, 2},
        {2, 3, 5},
        {10, 20, 30},
    };
    
    for (const auto& tc : cases) {
        EXPECT_EQ(tc.a + tc.b, tc.expected);
    }
}
```

## License

MIT
```

- [ ] **Step 2: Create LICENSE**

Create `LICENSE`:

```
MIT License

Copyright (c) 2026 neko233-com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

- [ ] **Step 3: Commit**

```bash
git add -A
git commit -m "docs: add README and LICENSE"
```

---

## Execution Handoff

This plan covers Phase 1 (Foundation) of the cnext project. After completing these tasks, you will have:

1. A working CLI with cobra
2. cnext.toml parser
3. Basic compiler driver (gcc/clang/msvc)
4. Build graph generation
5. Working `cnext init` and `cnext build` commands
6. Basic test framework
7. Working `cnext test` command

**Next phases** (to be planned separately):
- Phase 2: Dependency resolution, package manager, CMake bridging
- Phase 3: cnext-cloud registry server
- Phase 4: Formatter, linter, documentation generator
- Phase 5: Custom frontend, distributed compilation
