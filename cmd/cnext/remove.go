package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove [package]",
	Short: "Remove a dependency",
	Long: `Remove a package dependency from the project.

This removes the package from cnext.toml and cleans up downloaded files.

Example:
  cnext remove fmt`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("cnext remove %s - not yet implemented\n", args[0])
		return nil
	},
}
