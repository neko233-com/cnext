package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	updateDryRun bool
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update dependencies",
	Long: `Update all dependencies to their latest compatible versions.

Example:
  cnext update
  cnext update --dry-run`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("cnext update - not yet implemented")
		return nil
	},
}

func init() {
	updateCmd.Flags().BoolVar(&updateDryRun, "dry-run", false, "Preview what would be updated")
}
