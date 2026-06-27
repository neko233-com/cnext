package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	buildTarget string
	buildDebug  bool
	buildJobs   int
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the project",
	Long: `Build the current project.

By default, builds all targets defined in cnext.toml.
Use --target to build a specific target.

Example:
  cnext build --target mylib
  cnext build --debug`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("cnext build - not yet implemented")
		return nil
	},
}

func init() {
	buildCmd.Flags().StringVarP(&buildTarget, "target", "t", "", "Build a specific target")
	buildCmd.Flags().BoolVarP(&buildDebug, "debug", "d", false, "Build with debug symbols")
	buildCmd.Flags().IntVarP(&buildJobs, "jobs", "j", 0, "Number of parallel jobs (default: number of CPUs)")
}
