package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/neko233-com/cnext/internal/config"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [name]",
	Short: "Initialize a new cnext project",
	Long: `Initialize a new cnext project in the current directory.

This creates the basic project structure including:
  - cnext.toml configuration file
  - src/ directory for source files
  - tests/ directory for test files
  - include/ directory for public headers

Example:
  cnext init my-project`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := "my-project"
		if len(args) > 0 {
			name = args[0]
		}

		configPath := "cnext.toml"
		if _, err := os.Stat(configPath); err == nil {
			return fmt.Errorf("cnext.toml already exists in the current directory")
		}

		dirs := []string{"src", "include", "tests"}
		for _, dir := range dirs {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", dir, err)
			}
		}

		cfg := config.Default()
		cfg.Package.Name = name
		cfg.Build.Executables = []config.Executable{
			{
				Name:    name,
				Sources: []string{"src/**/*.cpp"},
			},
		}
		if err := config.Save(configPath, cfg); err != nil {
			return fmt.Errorf("failed to create cnext.toml: %w", err)
		}

		mainCpp := filepath.Join("src", "main.cpp")
		mainContent := `#include <iostream>

int main() {
    std::cout << "Hello, World!" << std::endl;
    return 0;
}
`
		if err := os.WriteFile(mainCpp, []byte(mainContent), 0644); err != nil {
			return fmt.Errorf("failed to create main.cpp: %w", err)
		}

		gitignore := `.gitignore
build/
*.o
*.obj
*.exe
*.dll
*.so
*.dylib
.cache/
`
		if err := os.WriteFile(".gitignore", []byte(gitignore), 0644); err != nil {
			return fmt.Errorf("failed to create .gitignore: %w", err)
		}

		fmt.Printf("Created project '%s' with the following structure:\n", name)
		fmt.Println("  cnext.toml")
		fmt.Println("  src/main.cpp")
		fmt.Println("  include/")
		fmt.Println("  tests/")
		fmt.Println("  .gitignore")
		return nil
	},
}
