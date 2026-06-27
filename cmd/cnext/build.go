package main

import (
	"fmt"
	"os"

	"github.com/neko233-com/cnext/internal/build"
	"github.com/neko233-com/cnext/internal/compiler"
	"github.com/neko233-com/cnext/internal/config"
	"github.com/spf13/cobra"
)

var (
	buildTarget   string
	buildDebug    bool
	buildJobs     int
	buildRelease  bool
	buildCompiler string
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
		cfgPath := "cnext.toml"
		if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
			return fmt.Errorf("cnext.toml not found in current directory")
		}

		cfg, err := config.Load(cfgPath)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		if buildCompiler != "" {
			cfg.Build.Compiler = buildCompiler
		}
		if buildTarget != "" {
			cfg.Build.Target.OS = buildTarget
		}
		if buildRelease {
			cfg.Build.Optimization = "release"
		} else if buildDebug {
			cfg.Build.Optimization = "debug"
		}

		graph, err := build.Generate(cfg)
		if err != nil {
			return fmt.Errorf("failed to generate build graph: %w", err)
		}

		comp, err := compiler.AutoDetect()
		if err != nil {
			return fmt.Errorf("failed to detect compiler: %w", err)
		}

		fmt.Printf("Building %s v%s with %s\n", cfg.Package.Name, cfg.Package.Version, comp.Name())
		fmt.Printf("STD: %s | Optimization: %s\n", cfg.Build.STD, cfg.Build.Optimization)

		if err := graph.Execute(".", comp); err != nil {
			return fmt.Errorf("build failed: %w", err)
		}

		fmt.Println("\nBuild completed successfully!")
		return nil
	},
}

func init() {
	buildCmd.Flags().StringVarP(&buildTarget, "target", "t", "", "Build a specific target")
	buildCmd.Flags().BoolVarP(&buildDebug, "debug", "d", false, "Build with debug symbols")
	buildCmd.Flags().BoolVarP(&buildRelease, "release", "r", false, "Build in release mode")
	buildCmd.Flags().IntVarP(&buildJobs, "jobs", "j", 0, "Number of parallel jobs (default: number of CPUs)")
	buildCmd.Flags().StringVarP(&buildCompiler, "compiler", "c", "", "Compiler to use (gcc, clang, etc.)")
}
