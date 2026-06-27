package compiler

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/neko233-com/cnext/internal/config"
)

func TestAutoDetect(t *testing.T) {
	comp, err := AutoDetect()
	if err != nil {
		t.Skipf("No compiler available: %v", err)
	}

	if comp == nil {
		t.Fatal("Expected non-nil compiler")
	}

	name := comp.Name()
	if name == "" {
		t.Error("Expected non-empty compiler name")
	}

	version := comp.Version()
	if version == "" {
		t.Error("Expected non-empty compiler version")
	}

	t.Logf("Detected compiler: %s (%s)", name, version)
}

func TestAutoDetectForTarget(t *testing.T) {
	target := config.Target{
		OS:   "auto",
		Arch: "auto",
	}

	comp, err := AutoDetectForTarget(target)
	if err != nil {
		t.Skipf("No compiler available: %v", err)
	}

	if comp == nil {
		t.Fatal("Expected non-nil compiler")
	}

	t.Logf("Detected compiler for target: %s", comp.Name())
}

func TestGCCName(t *testing.T) {
	gcc := NewGCC("/usr/bin/gcc")
	if gcc.Name() != "gcc" {
		t.Errorf("Expected name 'gcc', got '%s'", gcc.Name())
	}
}

func TestClangName(t *testing.T) {
	clang := NewClang("/usr/bin/clang")
	if clang.Name() != "clang" {
		t.Errorf("Expected name 'clang', got '%s'", clang.Name())
	}
}

func TestMSVCName(t *testing.T) {
	msvc := NewMSVC("cl.exe")
	if msvc.Name() != "msvc" {
		t.Errorf("Expected name 'msvc', got '%s'", msvc.Name())
	}
}

func TestCompileOptions(t *testing.T) {
	opts := CompileOptions{
		Sources:      []string{"main.cpp"},
		Output:       "main.o",
		Std:          "c++20",
		IncludeDirs:  []string{"include"},
		Defines:      []string{"DEBUG"},
		Flags:        []string{"-Wall"},
		Optimization: "debug",
	}

	if len(opts.Sources) != 1 {
		t.Errorf("Expected 1 source, got %d", len(opts.Sources))
	}
	if opts.Output != "main.o" {
		t.Errorf("Expected output 'main.o', got '%s'", opts.Output)
	}
	if opts.Std != "c++20" {
		t.Errorf("Expected std 'c++20', got '%s'", opts.Std)
	}
}

func TestLinkOptions(t *testing.T) {
	opts := LinkOptions{
		Objects:   []string{"main.o", "utils.o"},
		Output:    "main",
		Libraries: []string{"m", "pthread"},
		LibDirs:   []string{"/usr/lib"},
		Flags:     []string{"-static"},
	}

	if len(opts.Objects) != 2 {
		t.Errorf("Expected 2 objects, got %d", len(opts.Objects))
	}
	if opts.Output != "main" {
		t.Errorf("Expected output 'main', got '%s'", opts.Output)
	}
	if len(opts.Libraries) != 2 {
		t.Errorf("Expected 2 libraries, got %d", len(opts.Libraries))
	}
}

func TestCompilerInterface(t *testing.T) {
	var _ Compiler = (*GCC)(nil)
	var _ Compiler = (*Clang)(nil)
	var _ Compiler = (*MSVC)(nil)
}

func TestRunCommandIncludesStderr(t *testing.T) {
	cmd := exec.Command("go", "run", "nonexistent.go")
	err := runCommand(cmd)
	if err == nil {
		t.Fatal("Expected error from runCommand")
	}
	if !strings.Contains(err.Error(), "compiler output:") {
		t.Errorf("Expected error to contain 'compiler output:', got: %v", err)
	}
}
