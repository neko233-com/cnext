package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/neko233-com/cnext/internal/build"
	"github.com/neko233-com/cnext/internal/compiler"
	"github.com/neko233-com/cnext/internal/config"
	"github.com/spf13/cobra"
)

var (
	runTarget string
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the built executable",
	Long: `Run the built executable.

By default, runs the first executable defined in cnext.toml.
Use --target to run a specific executable.

Example:
  cnext run
  cnext run --target myapp
  cnext run -- arg1 arg2`,
	Args: cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfgPath := "cnext.toml"
		if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
			return fmt.Errorf("cnext.toml not found in current directory")
		}

		cfg, err := config.Load(cfgPath)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		if len(cfg.Build.Executables) == 0 {
			return fmt.Errorf("no executables defined in cnext.toml")
		}

		targetName := runTarget
		if targetName == "" {
			targetName = cfg.Build.Executables[0].Name
		}

		found := false
		for _, exe := range cfg.Build.Executables {
			if exe.Name == targetName {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("executable %q not found in cnext.toml", targetName)
		}

		binaryPath := targetName
		if runtime.GOOS == "windows" {
			binaryPath += ".exe"
		}

		if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
			fmt.Printf("Binary not found, building first...\n")

			comp, err := compiler.AutoDetect()
			if err != nil {
				return fmt.Errorf("failed to detect compiler: %w", err)
			}

			graph, err := build.Generate(cfg)
			if err != nil {
				return fmt.Errorf("failed to generate build graph: %w", err)
			}

			if err := graph.Execute(".", comp); err != nil {
				return fmt.Errorf("build failed: %w", err)
			}
		}

		binaryPath, err = filepath.Abs(binaryPath)
		if err != nil {
			return fmt.Errorf("failed to resolve binary path: %w", err)
		}

		fmt.Printf("Running %s...\n\n", targetName)

		execArgs := append([]string{binaryPath}, args...)
		process := exec.Command(execArgs[0], execArgs[1:]...)
		process.Stdout = os.Stdout
		process.Stderr = os.Stderr
		process.Stdin = os.Stdin

		if err := process.Run(); err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				os.Exit(exitErr.ExitCode())
			}
			return fmt.Errorf("failed to run executable: %w", err)
		}

		return nil
	},
}

func init() {
	runCmd.Flags().StringVarP(&runTarget, "target", "t", "", "Run a specific executable")
}
