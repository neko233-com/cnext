package main

import (
	"fmt"
	"os"

	pkg "github.com/neko233-com/cnext/internal/package"
	"github.com/spf13/cobra"
)

var (
	addVersion string
	addDev     bool
)

var addCmd = &cobra.Command{
	Use:   "add <package>[@version]",
	Short: "Add a dependency",
	Long: `Add a package dependency to the project.

This downloads and adds the package to cnext.toml.

Example:
  cnext add fmt
  cnext add yaml@2.3.1
  cnext add test-framework --dev`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		configPath := "cnext.toml"
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			return fmt.Errorf("cnext.toml not found. Run 'cnext init' first")
		}

		mod, err := pkg.LoadModule(configPath)
		if err != nil {
			return fmt.Errorf("failed to load cnext.toml: %w", err)
		}

		for _, arg := range args {
			name, constraint := parsePackageArg(arg)

			if addVersion != "" {
				constraint = addVersion
			}

			if mod.Deps == nil {
				mod.Deps = make(map[string]string)
			}
			if mod.DevDeps == nil {
				mod.DevDeps = make(map[string]string)
			}

			if addDev {
				mod.DevDeps[name] = constraint
			} else {
				mod.Deps[name] = constraint
			}

			fmt.Printf("Added %s@%s\n", name, constraint)
		}

		if err := pkg.SaveModule(mod, configPath); err != nil {
			return fmt.Errorf("failed to save cnext.toml: %w", err)
		}

		fmt.Println("Dependencies updated successfully")
		return nil
	},
}

func init() {
	addCmd.Flags().StringVarP(&addVersion, "version", "v", "", "Specific version to add")
	addCmd.Flags().BoolVar(&addDev, "dev", false, "Add as a development dependency")
}
