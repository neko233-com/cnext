package main

import (
	"fmt"
	"os"

	pkg "github.com/neko233-com/cnext/internal/package"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove <package>",
	Short: "Remove a dependency",
	Long: `Remove a package dependency from the project.

This removes the package from cnext.toml and cleans up downloaded files.

Example:
  cnext remove fmt`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		configPath := "cnext.toml"
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			return fmt.Errorf("cnext.toml not found. Run 'cnext init' first")
		}

		mod, err := pkg.LoadModule(configPath)
		if err != nil {
			return fmt.Errorf("failed to load cnext.toml: %w", err)
		}

		name := args[0]

		_, inDeps := mod.Deps[name]
		_, inDevDeps := mod.DevDeps[name]

		if !inDeps && !inDevDeps {
			return fmt.Errorf("package %s not found in dependencies", name)
		}

		if inDeps {
			delete(mod.Deps, name)
			fmt.Printf("Removed %s from dependencies\n", name)
		}

		if inDevDeps {
			delete(mod.DevDeps, name)
			fmt.Printf("Removed %s from dev-dependencies\n", name)
		}

		if err := pkg.SaveModule(mod, configPath); err != nil {
			return fmt.Errorf("failed to save cnext.toml: %w", err)
		}

		vendorDir := "vendor"
		pkgDir := fmt.Sprintf("%s/%s", vendorDir, name)
		if _, err := os.Stat(pkgDir); err == nil {
			if err := os.RemoveAll(pkgDir); err != nil {
				fmt.Printf("Warning: failed to remove vendor directory for %s: %v\n", name, err)
			} else {
				fmt.Printf("Removed vendor directory for %s\n", name)
			}
		}

		fmt.Println("Package removed successfully")
		return nil
	},
}
