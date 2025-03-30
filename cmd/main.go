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
	// Initialize viper
	v = viper.New()

	// Create the root command
	rootCmd := &cobra.Command{
		Use:   "stablemcp",
		Short: "StableMCP: A Model Context Protocol server for generating images",
		Long: `StableMCP is a server that implements the Model Context Protocol (MCP)
for generating images using Stable Diffusion models.

It provides a standardized API for generating images with various parameters
and supports features like configurable quality settings and authentication.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// bind the global flags to viper
			bindFlags(cmd, v)
			// Get the config path
			configPath, _ := v.Get("config").(string)
			if configPath != "" {
				// Check if the file exists
				if _, err := os.Stat(configPath); os.IsNotExist(err) {
					return fmt.Errorf("config file not found: %s", configPath)
				}
			}

			// Initialize and load the configuration
			cfg = config.NewConfig()
			var err error
			cfg, err = config.LoadConfig(v)
			if err != nil {
				return fmt.Errorf("error loading config: %w", err)
			}
			// if debug print out config
			if cfg.Debug {
				log.Printf("Loaded configuration: %+v\n", cfg)
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			// If no subcommand is provided, print help
			if err := cmd.Help(); err != nil {
				log.Printf("Error displaying help: %v", err)
			}
		},
	}

	// Add global flags
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Enable debug mode (default: false)")
	rootCmd.PersistentFlags().StringP("log-level", "l", "info", "Set the logging level (debug, info, warn, error)")
	rootCmd.PersistentFlags().StringP("output", "o", "json", "Output format (json, text)")
	rootCmd.PersistentFlags().StringP("config", "c", "", "Path to the configuration file")

	// Add the server command as a subcommand
	rootCmd.AddCommand(NewServerCmd())

	// Add the version command as a subcommand
	rootCmd.AddCommand(NewVersionCmd())

	return rootCmd
}

// Bind each cobra flag to its associated viper configuration (config file and environment variable)
func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	// Bind local flags
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

		// Bind the flag to viper
		if err := v.BindPFlag(configName, f); err != nil {
			log.Fatalf("Unable to bind flag '%s' to viper: %v", f.Name, err)
		}
	})

	// Also check persistent flags if they exist
	if cmd.PersistentFlags() != nil {
		cmd.PersistentFlags().VisitAll(func(f *pflag.Flag) {
			configName := f.Name
			if !f.Changed && v.IsSet(configName) {
				val := v.Get(configName)
				if err := cmd.PersistentFlags().Set(f.Name, fmt.Sprintf("%v", val)); err != nil {
					log.Fatalf("Unable to set persistent flag '%s' from config: %v", f.Name, err)
				}
			}

			// Bind the persistent flag to viper
			if err := v.BindPFlag(configName, f); err != nil {
				log.Fatalf("Unable to bind persistent flag '%s' to viper: %v", f.Name, err)
			}
		})
	}
}
