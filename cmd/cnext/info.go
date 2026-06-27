package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show project information",
	Long: `Display information about the current project.

Shows project name, version, dependencies, and build configuration.

Example:
  cnext info`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("cnext info - not yet implemented")
		return nil
	},
}
