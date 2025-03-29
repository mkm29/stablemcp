package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/mkm29/stablemcp/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	cfg *config.Config
	v   *viper.Viper
)

// NewRootCmd creates and returns the root command
func NewRootCmd() *cobra.Command {
	// Create the root command
	rootCmd := &cobra.Command{
		Use:   "stablemcp",
		Short: "StableMCP: A Model Context Protocol server for generating images",
		Long: `StableMCP is a server that implements the Model Context Protocol (MCP)
for generating images using Stable Diffusion models.

It provides a standardized API for generating images with various parameters
and supports features like configurable quality settings and authentication.`,
		Run: func(cmd *cobra.Command, args []string) {
			// If no subcommand is provided, print help
			cmd.Help()
		},
	}

	// Add version flag to the root command
	rootCmd.Flags().BoolP("version", "v", false, "Print the version number")
	
	// Add the server command as a subcommand
	rootCmd.AddCommand(NewServerCmd())
	
	return rootCmd
}

func NewServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start the stablemcp server",
		Long:  `Start the stablemcp server with the specified configuration.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Initialize viper
			v = viper.New()

			// Get the config flag
			configPath, _ := cmd.Flags().GetString("config")
			if configPath != "" {
				// Set the config flag in viper
				v.Set("config", configPath)
				// Check if the file exists
				if _, err := os.Stat(configPath); os.IsNotExist(err) {
					return fmt.Errorf("config file not found: %s", configPath)
				}
			}

			// Bind all flags to viper
			bindFlags(cmd, v)

			// Load the configuration
			var err error
			cfg, err = config.LoadConfig(v)
			if err != nil {
				return fmt.Errorf("error loading config: %w", err)
			}

			if cfg.Debug {
				log.Printf("Loaded configuration: %+v\n", cfg)
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			if cfg.Debug {
				log.Println("Running stablemcp Server")
			}
			if err := runServer(); err != nil {
				log.Fatalf("Server error: %v", err)
			}
		},
	}

	// Add flags
	cmd.Flags().String("config", "", "Path to the configuration file")
	
	// Add server-related flags
	cmd.Flags().Int("port", 8080, "Server port")
	cmd.Flags().String("host", "localhost", "Server host")
	cmd.Flags().Bool("debug", false, "Enable debug mode")
	
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

// Bind each cobra flag to its associated viper configuration (config file and environment variable)
func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		// Determine the naming convention of the flags when represented in the config file
		configName := f.Name

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && v.IsSet(configName) {
			val := v.Get(configName)
			if err := cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val)); err != nil {
				log.Fatalf("Unable to set flag '%s' from config: %v", f.Name, err)
			}
		}
	})
}
