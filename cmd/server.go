package cmd

import (
	"fmt"
	"log"

	"github.com/mkm29/stablemcp/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	cfg *config.Config
	v   *viper.Viper
)

func NewServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start the stablemcp server",
		Long:  `Start the stablemcp server with the specified configuration.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if cfg == nil {
				c, err := initializeConfig(cmd)
				if err != nil {
					return err
				}
				cfg = c
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			if cfg.EnableDebugMode {
				log.Println("Running sigpilot Server")
			}
			runServer()
		},
	}

	cmd.Flags().StringP("config", "c", "", "Path to the configuration file")
	return cmd
}

func runServer() error {
	// Here you would typically load the configuration, set up logging, etc.
	// For this example, we'll just print a message.
	println("Server is starting...")

	// Simulate server running
	select {}

	return nil
}

func initializeConfig(cmd *cobra.Command) (*config.Config, error) {
	// Initialize config
	v = viper.New()
	c, err := config.LoadConfig(v)
	if err != nil {
		return nil, err
	}
	cfg = c
	if c.EnableDebugMode {
		log.Printf("Config: %+v\n", cfg)
	}

	bindFlags(cmd, v)

	return c, nil
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
				log.Fatalf("unable to set flag '%s' from config: %v", f.Name, err)
			}
		}
	})
}

func Execute() {
	if err := NewServerCmd().Execute(); err != nil {
		panic(err)
	}
}
