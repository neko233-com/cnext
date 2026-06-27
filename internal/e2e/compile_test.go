package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

var cnextBin string

func findCnextBinary(t *testing.T) string {
	t.Helper()
	if cnextBin != "" {
		return cnextBin
	}

	bin := "cnext.exe"
	if runtime.GOOS != "windows" {
		bin = "cnext"
	}

	// Check multiple locations
	candidates := []string{}

	// 1. Current working directory
	if wd, err := os.Getwd(); err == nil {
		candidates = append(candidates,
			filepath.Join(wd, bin),
			filepath.Join(wd, bin+".exe"),
		)
	}

	// 2. Project root (relative to test file)
	_, filename, _, _ := runtime.Caller(0)
	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(filename)))
	candidates = append(candidates,
		filepath.Join(projectRoot, bin),
		filepath.Join(projectRoot, bin+".exe"),
	)

	// 3. Two levels up from test file (for CI)
	twoUp := filepath.Dir(filepath.Dir(filename))
	candidates = append(candidates,
		filepath.Join(twoUp, bin),
		filepath.Join(twoUp, bin+".exe"),
	)

	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			p, _ := filepath.Abs(c)
			cnextBin = p
			return cnextBin
		}
	}

	t.Fatal("cnext binary not found - run 'go build -o cnext.exe ./cmd/cnext' first")
	return ""
}

func runCnext(t *testing.T, dir string, args ...string) (string, error) {
	t.Helper()
	cnextPath := findCnextBinary(t)

	cmd := exec.Command(cnextPath, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func TestEndToEndCompile(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping e2e test in short mode")
	}

	tmpDir := t.TempDir()

	// cnext init creates files in the current directory, not in a subdirectory
	out, err := runCnext(t, tmpDir, "init", "test-project")
	if err != nil {
		t.Fatalf("cnext init failed: %v\n%s", err, out)
	}
	t.Logf("init output: %s", out)

	if _, err := os.Stat(filepath.Join(tmpDir, "cnext.toml")); os.IsNotExist(err) {
		t.Fatal("cnext.toml not created")
	}

	mainCpp := filepath.Join(tmpDir, "src", "main.cpp")
	if _, err := os.Stat(mainCpp); os.IsNotExist(err) {
		t.Fatal("src/main.cpp not created")
	}

	out, err = runCnext(t, tmpDir, "build")
	if err != nil {
		t.Fatalf("cnext build failed: %v\n%s", err, out)
	}
	t.Logf("build output: %s", out)

	binaryName := "test-project"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}
	if _, err := os.Stat(filepath.Join(tmpDir, binaryName)); os.IsNotExist(err) {
		t.Fatalf("binary %s not found after build", binaryName)
	}

	out, err = runCnext(t, tmpDir, "run")
	if err != nil {
		t.Fatalf("cnext run failed: %v\n%s", err, out)
	}
	t.Logf("run output: %s", out)
	if !strings.Contains(out, "Hello, World!") {
		t.Errorf("expected 'Hello, World!' in output, got: %s", out)
	}

	out, err = runCnext(t, tmpDir, "clean")
	if err != nil {
		t.Fatalf("cnext clean failed: %v\n%s", err, out)
	}
	if _, err := os.Stat(filepath.Join(tmpDir, "build")); !os.IsNotExist(err) {
		t.Error("build directory not removed after clean")
	}
}

func TestEndToEndMultiFileBuild(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping e2e test in short mode")
	}

	tmpDir := t.TempDir()

	out, err := runCnext(t, tmpDir, "init", "multi-file")
	if err != nil {
		t.Fatalf("cnext init failed: %v\n%s", err, out)
	}

	includeDir := filepath.Join(tmpDir, "include")
	os.MkdirAll(includeDir, 0755)

	helloHeader := `#pragma once
#include <string>
std::string GetHello();
`
	os.WriteFile(filepath.Join(includeDir, "hello.h"), []byte(helloHeader), 0644)

	helloSrc := `#include "hello.h"
std::string GetHello() { return "Hello from library!"; }
`
	helloSrcDir := filepath.Join(tmpDir, "src")
	os.WriteFile(filepath.Join(helloSrcDir, "hello.cpp"), []byte(helloSrc), 0644)

	mainCpp := `#include <iostream>
#include "hello.h"
int main() {
    std::cout << GetHello() << std::endl;
    return 0;
}
`
	os.WriteFile(filepath.Join(helloSrcDir, "main.cpp"), []byte(mainCpp), 0644)

	toml := `[package]
name = "multi-file"
version = "0.1.0"
edition = "2024"

[build]
compiler = "auto"
std = "c++20"
optimization = "debug"

[[build.executables]]
name = "multi-file"
sources = ["src/main.cpp"]
libraries = ["mylib"]

[[build.libraries]]
name = "mylib"
type = "static"
sources = ["src/hello.cpp"]
include_dirs = ["include"]
`
	os.WriteFile(filepath.Join(tmpDir, "cnext.toml"), []byte(toml), 0644)

	out, err = runCnext(t, tmpDir, "build")
	if err != nil {
		t.Fatalf("cnext build failed: %v\n%s", err, out)
	}
	t.Logf("build output: %s", out)

	out, err = runCnext(t, tmpDir, "run")
	if err != nil {
		t.Fatalf("cnext run failed: %v\n%s", err, out)
	}
	t.Logf("run output: %s", out)
	if !strings.Contains(out, "Hello from library!") {
		t.Errorf("expected 'Hello from library!' in output, got: %s", out)
	}
}

func TestEndToEndInfo(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping e2e test in short mode")
	}

	tmpDir := t.TempDir()

	out, err := runCnext(t, tmpDir, "init", "info-project")
	if err != nil {
		t.Fatalf("cnext init failed: %v\n%s", err, out)
	}

	out, err = runCnext(t, tmpDir, "info")
	if err != nil {
		t.Fatalf("cnext info failed: %v\n%s", err, out)
	}
	t.Logf("info output: %s", out)

	if !strings.Contains(out, "info-project") {
		t.Errorf("expected project name in info output, got: %s", out)
	}
}
