package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/neko233-com/cnext/internal/compiler"
	"github.com/neko233-com/cnext/internal/config"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show project information",
	Long: `Display information about the current project.

Shows project name, version, dependencies, and build configuration.

Example:
  cnext info`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfgPath := "cnext.toml"
		if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
			return fmt.Errorf("cnext.toml not found in current directory")
		}

		cfg, err := config.Load(cfgPath)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

		fmt.Fprintln(w, "Package:")
		fmt.Fprintf(w, "  Name:\t%s\n", cfg.Package.Name)
		fmt.Fprintf(w, "  Version:\t%s\n", cfg.Package.Version)
		fmt.Fprintf(w, "  Edition:\t%s\n", cfg.Package.Edition)
		fmt.Fprintln(w)

		fmt.Fprintln(w, "Build:")
		fmt.Fprintf(w, "  Compiler:\t%s\n", cfg.Build.Compiler)
		fmt.Fprintf(w, "  STD:\t%s\n", cfg.Build.STD)
		fmt.Fprintf(w, "  Optimization:\t%s\n", cfg.Build.Optimization)
		fmt.Fprintf(w, "  LTO:\t%v\n", cfg.Build.LTO)
		fmt.Fprintf(w, "  PIC:\t%v\n", cfg.Build.PIC)
		fmt.Fprintf(w, "  Target OS:\t%s\n", cfg.Build.Target.OS)
		fmt.Fprintf(w, "  Target Arch:\t%s\n", cfg.Build.Target.Arch)
		fmt.Fprintln(w)

		if len(cfg.Build.Libraries) > 0 {
			fmt.Fprintln(w, "Libraries:")
			for _, lib := range cfg.Build.Libraries {
				fmt.Fprintf(w, "  - %s (%s)\n", lib.Name, lib.Type)
			}
			fmt.Fprintln(w)
		}

		if len(cfg.Build.Executables) > 0 {
			fmt.Fprintln(w, "Executables:")
			for _, exe := range cfg.Build.Executables {
				fmt.Fprintf(w, "  - %s\n", exe.Name)
			}
			fmt.Fprintln(w)
		}

		if len(cfg.Dependencies) > 0 {
			fmt.Fprintln(w, "Dependencies:")
			for name, dep := range cfg.Dependencies {
				fmt.Fprintf(w, "  - %s@%s\n", name, dep.Version)
			}
			fmt.Fprintln(w)
		}

		if len(cfg.DevDependencies) > 0 {
			fmt.Fprintln(w, "Dev Dependencies:")
			for name, dep := range cfg.DevDependencies {
				fmt.Fprintf(w, "  - %s@%s\n", name, dep.Version)
			}
			fmt.Fprintln(w)
		}

		comp, err := compiler.AutoDetect()
		if err == nil {
			fmt.Fprintln(w, "Detected Compiler:")
			fmt.Fprintf(w, "  Name:\t%s\n", comp.Name())
			fmt.Fprintf(w, "  Version:\t%s\n", comp.Version())
		}

		w.Flush()
		return nil
	},
}
