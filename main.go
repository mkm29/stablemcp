package main

import (
	"fmt"
	"os"

	"github.com/mkm29/stablemcp/cmd"
)

func main() {
	// Create MCPServer instance
	// server := mcp.NewMCPServer()

	rootCmd := cmd.NewRootCmd()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
