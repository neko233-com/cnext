package test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/neko233-com/cnext/internal/compiler"
	"github.com/neko233-com/cnext/internal/config"
)

type Runner struct {
	cfg      *config.Config
	comp     compiler.Compiler
	filter   string
	verbose  bool
	buildDir string
}

type TestResult struct {
	Name   string
	Passed bool
	Output string
}

type RunResult struct {
	Passed  bool
	Total   int
	Failed  int
	Results []TestResult
}

func New(cfg *config.Config, comp compiler.Compiler) *Runner {
	return &Runner{
		cfg:      cfg,
		comp:     comp,
		buildDir: filepath.Join(".", "build", "test"),
	}
}

func (r *Runner) SetFilter(filter string) {
	r.filter = filter
}

func (r *Runner) SetVerbose(verbose bool) {
	r.verbose = verbose
}

func (r *Runner) DiscoverTests(dirs []string) ([]string, error) {
	var testFiles []string

	for _, dir := range dirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			continue
		}
		if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if info.IsDir() {
				return nil
			}

			ext := filepath.Ext(path)
			if ext == ".cpp" || ext == ".cc" {
				base := filepath.Base(path)
				if strings.HasSuffix(base, "_test"+ext) {
					testFiles = append(testFiles, path)
				}
			}
			return nil
		}); err != nil {
			return nil, fmt.Errorf("failed to walk directory %s: %w", dir, err)
		}
	}

	return testFiles, nil
}

func (r *Runner) CompileTests(testFiles []string) (string, error) {
	if err := os.MkdirAll(r.buildDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create build directory: %w", err)
	}

	output := filepath.Join(r.buildDir, "test_runner")
	if runtime.GOOS == "windows" {
		output += ".exe"
	}

	var sources []string
	sources = append(sources, testFiles...)

	var incDirs []string
	for _, dir := range []string{"include", "."} {
		if _, err := os.Stat(dir); err == nil {
			incDirs = append(incDirs, dir)
		}
	}

	compileOpts := compiler.CompileOptions{
		Sources:      sources,
		Output:       output,
		Std:          r.cfg.Build.STD,
		IncludeDirs:  incDirs,
		Optimization: "debug",
	}

	if err := r.comp.Compile(compileOpts); err != nil {
		return "", fmt.Errorf("failed to compile tests: %w", err)
	}

	return output, nil
}

func (r *Runner) RunTests(binaryPath string) (*RunResult, error) {
	args := []string{}
	if r.filter != "" {
		args = append(args, "--filter", r.filter)
	}
	if r.verbose {
		args = append(args, "--verbose")
	}

	cmd := exec.Command(binaryPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	result := &RunResult{
		Passed: err == nil,
	}

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() != 0 {
				result.Passed = false
				result.Total = 1
				result.Failed = 1
			}
		} else {
			return nil, fmt.Errorf("failed to run tests: %w", err)
		}
	} else {
		result.Total = 1
		result.Failed = 0
	}

	return result, nil
}
