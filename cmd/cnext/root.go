package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cnext",
	Short: "A build tool and package manager for C/C++ projects",
	Long: `cnext is a modern build tool and package manager for C/C++ projects.

It provides a simple and consistent interface for managing C/C++ projects,
including building, testing, dependency management, and more.

Get started by creating a new project:

  cnext init my-project`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(testCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(fmtCmd)
	rootCmd.AddCommand(lintCmd)
	rootCmd.AddCommand(docCmd)
	rootCmd.AddCommand(publishCmd)
	rootCmd.AddCommand(cleanCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(infoCmd)
	rootCmd.AddCommand(toolchainCmd)
}
