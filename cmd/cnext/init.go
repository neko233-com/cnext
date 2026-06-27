package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [name]",
	Short: "Initialize a new cnext project",
	Long: `Initialize a new cnext project in the current directory.

This creates the basic project structure including:
  - cnext.toml configuration file
  - src/ directory for source files
  - tests/ directory for test files
  - include/ directory for public headers

Example:
  cnext init my-project`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("cnext init - not yet implemented")
		return nil
	},
}
