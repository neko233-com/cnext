package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	runTarget string
	runArgs   []string
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the built executable",
	Long: `Run the built executable.

By default, runs the main executable defined in cnext.toml.
Use --target to run a specific executable.

Example:
  cnext run
  cnext run --target myapp`,
	Args: cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("cnext run - not yet implemented")
		return nil
	},
}

func init() {
	runCmd.Flags().StringVarP(&runTarget, "target", "t", "", "Run a specific executable")
}
