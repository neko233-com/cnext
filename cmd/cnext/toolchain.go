package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/neko233-com/cnext/internal/toolchain"
	"github.com/spf13/cobra"
)

var toolchainCmd = &cobra.Command{
	Use:   "toolchain",
	Short: "Manage toolchains",
	Long: `Manage C/C++ toolchains.

List available toolchains, install new ones, or switch between them.

Examples:
  cnext toolchain list
  cnext toolchain install gcc@14.2.0
  cnext toolchain install clang
  cnext toolchain remove gcc@14.2.0`,
}

var toolchainListCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed toolchains",
	Long:  `List all installed toolchains.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		mgr := toolchain.NewManager()
		installed := mgr.List()

		if len(installed) == 0 {
			fmt.Println("No toolchains installed.")
			fmt.Println("\nInstall a toolchain with: cnext toolchain install <name>[@version]")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tVERSION\tPLATFORM\tPATH")
		fmt.Fprintln(w, "----\t-------\t--------\t----")
		for _, tc := range installed {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", tc.Name, tc.Version, tc.Platform, tc.Path)
		}
		w.Flush()

		return nil
	},
}

var toolchainInstallCmd = &cobra.Command{
	Use:   "install [name@version]",
	Short: "Install a toolchain",
	Long: `Install a new C/C++ toolchain.

If version is not specified, the latest available version is installed.

Supported toolchains:
  - gcc (various versions)
  - clang (various versions)

Examples:
  cnext toolchain install gcc@14.2.0
  cnext toolchain install clang
  cnext toolchain install gcc@14.2.0 --platform linux-arm64`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name, version := toolchain.ParseToolchainArg(args[0])
		platform, _ := cmd.Flags().GetString("platform")
		if platform == "" {
			platform = toolchain.CurrentPlatform()
		}

		var tc *toolchain.ToolchainInfo
		var err error

		if version != "" {
			tc, err = toolchain.FindInRegistry(name, version, platform)
		} else {
			tc, err = toolchain.FindLatest(name, platform)
		}

		if err != nil {
			return err
		}

		mgr := toolchain.NewManager()
		if mgr.IsInstalled(tc.Name, tc.Version, tc.Platform) {
			fmt.Printf("Toolchain %s@%s (%s) is already installed.\n", tc.Name, tc.Version, tc.Platform)
			return nil
		}

		fmt.Printf("Installing %s@%s for %s...\n", tc.Name, tc.Version, tc.Platform)
		if err := mgr.Install(*tc); err != nil {
			return fmt.Errorf("failed to install toolchain: %w", err)
		}

		fmt.Printf("Successfully installed %s@%s (%s)\n", tc.Name, tc.Version, tc.Platform)
		return nil
	},
}

var toolchainRemoveCmd = &cobra.Command{
	Use:   "remove [name@version]",
	Short: "Remove a toolchain",
	Long: `Remove an installed toolchain.

Examples:
  cnext toolchain remove gcc@14.2.0
  cnext toolchain remove clang`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name, version := toolchain.ParseToolchainArg(args[0])
		platform, _ := cmd.Flags().GetString("platform")
		if platform == "" {
			platform = toolchain.CurrentPlatform()
		}

		if version == "" {
			latest, err := toolchain.FindLatest(name, platform)
			if err != nil {
				return err
			}
			version = latest.Version
		}

		mgr := toolchain.NewManager()
		if !mgr.IsInstalled(name, version, platform) {
			return fmt.Errorf("toolchain %s@%s (%s) is not installed", name, version, platform)
		}

		fmt.Printf("Removing %s@%s (%s)...\n", name, version, platform)
		if err := mgr.Remove(name, version, platform); err != nil {
			return fmt.Errorf("failed to remove toolchain: %w", err)
		}

		fmt.Printf("Successfully removed %s@%s (%s)\n", name, version, platform)
		return nil
	},
}

func init() {
	toolchainInstallCmd.Flags().StringP("platform", "p", "", "Target platform (default: current)")
	toolchainRemoveCmd.Flags().StringP("platform", "p", "", "Target platform (default: current)")

	toolchainCmd.AddCommand(toolchainListCmd)
	toolchainCmd.AddCommand(toolchainInstallCmd)
	toolchainCmd.AddCommand(toolchainRemoveCmd)
}
