package main

import (
	"fmt"
	"os"

	"github.com/neko233-com/cnext/internal/compiler"
	"github.com/neko233-com/cnext/internal/config"
	"github.com/neko233-com/cnext/internal/test"
	"github.com/spf13/cobra"
)

var (
	testTarget    string
	testFilter    string
	testVerbose   bool
	testCoverage  bool
	testBenchmark bool
	testWatch     bool
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Run tests",
	Long: `Run all tests in the project.

Tests are discovered from tests/ directory and project sources.
Use --filter to run specific tests.

Example:
  cnext test
  cnext test --filter "test_*"
  cnext test --verbose`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfgPath := "cnext.toml"
		if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
			return fmt.Errorf("cnext.toml not found in current directory")
		}

		cfg, err := config.Load(cfgPath)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		comp, err := compiler.AutoDetect()
		if err != nil {
			return fmt.Errorf("failed to detect compiler: %w", err)
		}

		fmt.Printf("Using compiler: %s\n", comp.Name())

		runner := test.New(cfg, comp)
		runner.SetFilter(testFilter)
		runner.SetVerbose(testVerbose)

		dirs := []string{"tests", "test"}
		var testFiles []string

		if len(cfg.Build.Tests) > 0 {
			for _, tgt := range cfg.Build.Tests {
				if tgt.Sources != nil {
					testFiles = append(testFiles, tgt.Sources...)
				}
			}
		}

		if len(testFiles) == 0 {
			testFiles, err = runner.DiscoverTests(dirs)
			if err != nil {
				return fmt.Errorf("failed to discover tests: %w", err)
			}
		}

		if len(testFiles) == 0 {
			fmt.Println("No test files found")
			return nil
		}

		fmt.Printf("Found %d test file(s)\n", len(testFiles))
		for _, f := range testFiles {
			fmt.Printf("  - %s\n", f)
		}

		binaryPath, err := runner.CompileTests(testFiles)
		if err != nil {
			return fmt.Errorf("failed to compile tests: %w", err)
		}

		fmt.Printf("Test binary: %s\n", binaryPath)

		result, err := runner.RunTests(binaryPath)
		if err != nil {
			return fmt.Errorf("failed to run tests: %w", err)
		}

		fmt.Printf("\nResults: %d passed, %d failed, %d total\n",
			result.Total-result.Failed, result.Failed, result.Total)

		if !result.Passed {
			return fmt.Errorf("tests failed")
		}

		return nil
	},
}

func init() {
	testCmd.Flags().StringVarP(&testTarget, "target", "t", "", "Run tests for a specific target")
	testCmd.Flags().StringVarP(&testFilter, "filter", "f", "", "Filter tests by name pattern")
	testCmd.Flags().BoolVarP(&testVerbose, "verbose", "v", false, "Verbose test output")
	testCmd.Flags().BoolVarP(&testCoverage, "coverage", "c", false, "Generate coverage report")
	testCmd.Flags().BoolVarP(&testBenchmark, "benchmark", "b", false, "Run benchmarks")
	testCmd.Flags().BoolVarP(&testWatch, "watch", "w", false, "Watch for changes and re-run tests")
}
