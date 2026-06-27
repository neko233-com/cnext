package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	cleanAll bool
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean build artifacts",
	Long: `Remove build artifacts and temporary files.

By default, removes the build/ directory.
Use --all to also remove downloaded dependencies.

Example:
  cnext clean
  cnext clean --all`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("cnext clean - not yet implemented")
		return nil
	},
}

func init() {
	cleanCmd.Flags().BoolVarP(&cleanAll, "all", "a", false, "Remove all generated files including dependencies")
}
