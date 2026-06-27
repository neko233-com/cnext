package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	addVersion string
	addDev     bool
)

var addCmd = &cobra.Command{
	Use:   "add [package]",
	Short: "Add a dependency",
	Long: `Add a package dependency to the project.

This downloads and adds the package to cnext.toml.

Example:
  cnext add fmt
  cnext add yaml --version 2.3.1
  cnext add test-framework --dev`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("cnext add %s - not yet implemented\n", args[0])
		return nil
	},
}

func init() {
	addCmd.Flags().StringVarP(&addVersion, "version", "v", "", "Specific version to add")
	addCmd.Flags().BoolVar(&addDev, "dev", false, "Add as a development dependency")
}
