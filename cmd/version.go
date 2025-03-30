package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/mkm29/stablemcp/internal/version"
	"github.com/spf13/cobra"
)

// NewVersionCmd returns a new version command
func NewVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version information",
		Long:  `Print detailed version information about the StableMCP server.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Get output format from persistent flags
			outputFormat, _ := cmd.Parent().PersistentFlags().GetString("output")

			// Check the output format
			if outputFormat == "json" {
				// Output as JSON
				info := version.Info()
				jsonData, err := json.MarshalIndent(info, "", "  ")
				if err != nil {
					fmt.Printf("Error: %v\n", err)
					return
				}
				fmt.Println(string(jsonData))
			} else {
				// Human-readable output
				fmt.Printf("StableMCP Version: %s\n", version.Version)
				fmt.Printf("Build Date: %s\n", version.BuildDate)
				fmt.Printf("Git Commit: %s\n", version.GitCommit)
				fmt.Printf("Git Branch: %s\n", version.GitBranch)
			}
		},
	}

	return cmd
}