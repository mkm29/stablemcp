package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

func NewServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start the stablemcp server",
		Long:  `Start the stablemcp server with the specified configuration.`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// do any necessary setup before running the server
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("Running stablemcp Server")
			if err := runServer(); err != nil {
				log.Fatalf("Server error: %v", err)
			}
		},
	}

	// Add server-related flags
	cmd.Flags().Int("port", 8080, "Server port")
	cmd.Flags().String("host", "localhost", "Server host")

	return cmd
}

func runServer() error {
	// Here you would typically set up the server, logging, etc.
	// For now, we'll just print information about the configuration
	fmt.Printf("Starting server on %s:%d\n", cfg.Server.Host, cfg.Server.Port)
	fmt.Printf("Debug mode: %v\n", cfg.Debug)
	fmt.Printf("Log level: %s\n", cfg.Logging.Level)

	// Simulate server running
	select {}
}
