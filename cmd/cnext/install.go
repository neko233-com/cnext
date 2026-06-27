package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Masterminds/semver/v3"
	pkg "github.com/neko233-com/cnext/internal/package"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install all dependencies",
	Long: `Install all dependencies listed in cnext.toml.

This downloads and installs all required packages to the vendor directory.

Example:
  cnext install`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		configPath := "cnext.toml"
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			return fmt.Errorf("cnext.toml not found. Run 'cnext init' first")
		}

		mod, err := pkg.LoadModule(configPath)
		if err != nil {
			return fmt.Errorf("failed to load cnext.toml: %w", err)
		}

		if len(mod.Deps) == 0 && len(mod.DevDeps) == 0 {
			fmt.Println("No dependencies to install")
			return nil
		}

		vendorDir := "vendor"
		if err := os.MkdirAll(vendorDir, 0755); err != nil {
			return fmt.Errorf("failed to create vendor directory: %w", err)
		}

		cache := pkg.NewCache()

		for name, constraint := range mod.Deps {
			fmt.Printf("Installing %s@%s...\n", name, constraint)

			if path, ok := cache.Get(name, constraint); ok {
				fmt.Printf("  Using cached version at %s\n", path)
				if err := extractPackage(path, vendorDir, name); err != nil {
					return fmt.Errorf("failed to extract %s: %w", name, err)
				}
				continue
			}

			fmt.Printf("  Downloading %s@%s...\n", name, constraint)
			path, err := downloadPackage(name, constraint)
			if err != nil {
				return fmt.Errorf("failed to download %s: %w", name, err)
			}

			if err := cache.Put(name, constraint, path); err != nil {
				fmt.Printf("  Warning: failed to cache %s: %v\n", name, err)
			}

			if err := extractPackage(path, vendorDir, name); err != nil {
				return fmt.Errorf("failed to extract %s: %w", name, err)
			}
		}

		for name, constraint := range mod.DevDeps {
			fmt.Printf("Installing dev dependency %s@%s...\n", name, constraint)

			if path, ok := cache.Get(name, constraint); ok {
				fmt.Printf("  Using cached version at %s\n", path)
				if err := extractPackage(path, vendorDir, name); err != nil {
					return fmt.Errorf("failed to extract %s: %w", name, err)
				}
				continue
			}

			fmt.Printf("  Downloading %s@%s...\n", name, constraint)
			path, err := downloadPackage(name, constraint)
			if err != nil {
				return fmt.Errorf("failed to download %s: %w", name, err)
			}

			if err := cache.Put(name, constraint, path); err != nil {
				fmt.Printf("  Warning: failed to cache %s: %v\n", name, err)
			}

			if err := extractPackage(path, vendorDir, name); err != nil {
				return fmt.Errorf("failed to extract %s: %w", name, err)
			}
		}

		fmt.Println("All dependencies installed successfully")
		return nil
	},
}

func downloadPackage(name, constraint string) (string, error) {
	registry := pkg.NewHTTPRegistry("https://registry.cnext.dev")

	versions, err := registry.GetVersions(name)
	if err != nil {
		return "", err
	}

	parsedVersions := make([]*semver.Version, 0, len(versions))
	for _, v := range versions {
		pv, err := semver.NewVersion(v)
		if err != nil {
			continue
		}
		parsedVersions = append(parsedVersions, pv)
	}

	c, err := pkg.ParseConstraint(constraint)
	if err != nil {
		return "", err
	}

	resolved, err := pkg.ResolveVersion(parsedVersions, c)
	if err != nil {
		return "", err
	}

	return registry.Download(name, resolved.String())
}

func extractPackage(archivePath, vendorDir, name string) error {
	destDir := filepath.Join(vendorDir, name)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}

	data, err := os.ReadFile(archivePath)
	if err != nil {
		return err
	}

	mainFile := filepath.Join(destDir, name+".h")
	if err := os.WriteFile(mainFile, data, 0644); err != nil {
		return err
	}

	return nil
}

func parsePackageArg(arg string) (name, constraint string) {
	parts := strings.SplitN(arg, "@", 2)
	if len(parts) == 1 {
		return parts[0], "*"
	}
	return parts[0], parts[1]
}
