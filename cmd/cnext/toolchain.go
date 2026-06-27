package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var toolchainCmd = &cobra.Command{
	Use:   "toolchain",
	Short: "Manage toolchains",
	Long: `Manage C/C++ toolchains.

List available toolchains, install new ones, or switch between them.

Examples:
  cnext toolchain list
  cnext toolchain install gcc-13
  cnext toolchain use clang`,
}

var toolchainListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available toolchains",
	Long:  `List all installed and available toolchains.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("cnext toolchain list - not yet implemented")
		return nil
	},
}

var toolchainInstallCmd = &cobra.Command{
	Use:   "install [toolchain]",
	Short: "Install a toolchain",
	Long: `Install a new C/C++ toolchain.

Supported toolchains:
  - gcc (various versions)
  - clang (various versions)
  - msvc (Microsoft Visual C++)

Example:
  cnext toolchain install gcc-13
  cnext toolchain install clang-16`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("cnext toolchain install %s - not yet implemented\n", args[0])
		return nil
	},
}

var toolchainUseCmd = &cobra.Command{
	Use:   "use [toolchain]",
	Short: "Set active toolchain",
	Long: `Set the active toolchain for the current project.

Example:
  cnext toolchain use gcc-13`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("cnext toolchain use %s - not yet implemented\n", args[0])
		return nil
	},
}

func init() {
	toolchainCmd.AddCommand(toolchainListCmd)
	toolchainCmd.AddCommand(toolchainInstallCmd)
	toolchainCmd.AddCommand(toolchainUseCmd)
}
