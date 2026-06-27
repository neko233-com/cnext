package main

import (
	"fmt"
	"os"

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
Use --all to also remove downloaded dependencies and cmake cache.

Example:
  cnext clean
  cnext clean --all`,
	RunE: func(cmd *cobra.Command, args []string) error {
		dirs := []string{"build"}

		if cleanAll {
			dirs = append(dirs, ".cnext", "vendor")
		}

		for _, dir := range dirs {
			if _, err := os.Stat(dir); os.IsNotExist(err) {
				continue
			}

			fmt.Printf("Removing %s/\n", dir)
			if err := os.RemoveAll(dir); err != nil {
				return fmt.Errorf("failed to remove %s: %w", dir, err)
			}
		}

		fmt.Println("Clean completed successfully")
		return nil
	},
}

func init() {
	cleanCmd.Flags().BoolVarP(&cleanAll, "all", "a", false, "Remove all generated files including dependencies")
}
