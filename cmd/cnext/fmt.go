package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	fmtCheck   bool
	fmtExclude []string
)

var fmtCmd = &cobra.Command{
	Use:   "fmt",
	Short: "Format source code",
	Long: `Format all source files in the project.

Uses clang-format for C/C++ files with project-specific configuration.

Example:
  cnext fmt
  cnext fmt --check`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("cnext fmt - not yet implemented")
		return nil
	},
}

func init() {
	fmtCmd.Flags().BoolVar(&fmtCheck, "check", false, "Check if formatting is needed without modifying files")
	fmtCmd.Flags().StringArrayVar(&fmtExclude, "exclude", nil, "Exclude directories or patterns")
}
