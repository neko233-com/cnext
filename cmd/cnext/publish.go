package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	publishRegistry string
	publishDryRun   bool
)

var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publish to registry",
	Long: `Publish the project package to a registry.

This builds and uploads the package to the specified registry.

Example:
  cnext publish
  cnext publish --registry my-registry
  cnext publish --dry-run`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("cnext publish - not yet implemented")
		return nil
	},
}

func init() {
	publishCmd.Flags().StringVarP(&publishRegistry, "registry", "r", "", "Registry to publish to")
	publishCmd.Flags().BoolVar(&publishDryRun, "dry-run", false, "Preview what would be published")
}
