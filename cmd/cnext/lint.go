package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	lintFix    bool
	lintOutput string
)

var lintCmd = &cobra.Command{
	Use:   "lint",
	Short: "Lint source code",
	Long: `Run static analysis and linting on source files.

Uses clang-tidy and other tools to check for common issues.

Example:
  cnext lint
  cnext lint --fix`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("cnext lint - not yet implemented")
		return nil
	},
}

func init() {
	lintCmd.Flags().BoolVarP(&lintFix, "fix", "f", false, "Automatically fix lint issues")
	lintCmd.Flags().StringVarP(&lintOutput, "output", "o", "", "Output format (text, json, checkstyle)")
}
