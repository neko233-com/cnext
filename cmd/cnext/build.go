package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/neko233-com/cnext/internal/build"
	"github.com/neko233-com/cnext/internal/compiler"
	"github.com/neko233-com/cnext/internal/config"
	"github.com/neko233-com/cnext/internal/cpu"
	"github.com/spf13/cobra"
)

var (
	buildTarget    string
	buildDebug     bool
	buildJobs      int
	buildRelease   bool
	buildCompiler  string
	buildIncremental bool
	buildFull      bool
	buildNoCache   bool
	buildLTO       bool
	buildNative    bool
	buildVerbose   bool
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the project",
	Long: `Build the current project.

By default, builds all targets defined in cnext.toml.
Use --target to build a specific target.

Examples:
  cnext build                     # Incremental build with cache
  cnext build --full              # Full rebuild (ignore cache)
  cnext build --jobs 8            # Parallel build with 8 jobs
  cnext build --lto               # Enable Link-Time Optimization
  cnext build --native            # Optimize for current CPU
  cnext build --debug             # Debug symbols
  cnext build --compiler clang    # Use specific compiler`,
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

		// Detect CPU features
		cpuFeatures := cpu.Detect()

		// Build options
		opts := build.BuildOptions{
			Jobs:        buildJobs,
			Incremental: !buildFull,
			Full:        buildFull,
			Cache:       !buildNoCache,
			LTO:         buildLTO,
			Native:      buildNative,
			Verbose:     buildVerbose,
		}

		fmt.Printf("Building %s v%s with %s\n", cfg.Package.Name, cfg.Package.Version, comp.Name())
		fmt.Printf("STD: %s | Optimization: %s | Jobs: %d\n", cfg.Build.STD, cfg.Build.Optimization, opts.Jobs)
		fmt.Printf("CPU: %s\n", cpuFeatures.String())

		start := time.Now()

		result, err := graph.ExecuteWithOpts(".", comp, opts)
		if err != nil {
			return fmt.Errorf("build failed: %w", err)
		}

		elapsed := time.Since(start)

		fmt.Printf("\nBuild completed in %s\n", elapsed.Round(time.Millisecond))
		fmt.Printf("  Compiled: %d | Cached: %d | Linked: %d | Archived: %d\n",
			result.Compiled, result.Cached, result.Linked, result.Archived)

		return nil
	},
}

func init() {
	numCPU := runtime.NumCPU()
	defaultJobs := numCPU - 1
	if defaultJobs < 1 {
		defaultJobs = 1
	}

	buildCmd.Flags().StringVarP(&buildTarget, "target", "t", "", "Build a specific target")
	buildCmd.Flags().BoolVarP(&buildDebug, "debug", "d", false, "Build with debug symbols")
	buildCmd.Flags().BoolVarP(&buildRelease, "release", "r", false, "Build in release mode")
	buildCmd.Flags().IntVarP(&buildJobs, "jobs", "j", defaultJobs, fmt.Sprintf("Parallel jobs (default: %d)", defaultJobs))
	buildCmd.Flags().StringVarP(&buildCompiler, "compiler", "c", "", "Compiler to use (gcc, clang, etc.)")
	buildCmd.Flags().BoolVar(&buildIncremental, "incremental", true, "Enable incremental builds (default: true)")
	buildCmd.Flags().BoolVar(&buildFull, "full", false, "Full rebuild (ignore cache)")
	buildCmd.Flags().BoolVar(&buildNoCache, "no-cache", false, "Disable build cache")
	buildCmd.Flags().BoolVar(&buildLTO, "lto", false, "Enable Link-Time Optimization")
	buildCmd.Flags().BoolVar(&buildNative, "native", false, "Optimize for current CPU architecture")
	buildCmd.Flags().BoolVarP(&buildVerbose, "verbose", "v", false, "Verbose output")
}
