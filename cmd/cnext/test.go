package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	testTarget string
	testFilter string
	testVerbose bool
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Run tests",
	Long: `Run all tests in the project.

Tests are discovered from the tests/ directory.
Use --filter to run specific tests.

Example:
  cnext test
  cnext test --filter "test_*"
  cnext test --verbose`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("cnext test - not yet implemented")
		return nil
	},
}

func init() {
	testCmd.Flags().StringVarP(&testTarget, "target", "t", "", "Run tests for a specific target")
	testCmd.Flags().StringVarP(&testFilter, "filter", "f", "", "Filter tests by name pattern")
	testCmd.Flags().BoolVarP(&testVerbose, "verbose", "v", false, "Verbose test output")
}
