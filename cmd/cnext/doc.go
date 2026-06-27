package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	docFormat string
	docOutput string
)

var docCmd = &cobra.Command{
	Use:   "doc",
	Short: "Generate documentation",
	Long: `Generate documentation for the project.

Uses Doxygen or other tools to generate API documentation from source comments.

Example:
  cnext doc
  cnext doc --format html
  cnext doc --output docs/`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("cnext doc - not yet implemented")
		return nil
	},
}

func init() {
	docCmd.Flags().StringVarP(&docFormat, "format", "f", "html", "Documentation format (html, latex, man)")
	docCmd.Flags().StringVarP(&docOutput, "output", "o", "docs", "Output directory")
}
