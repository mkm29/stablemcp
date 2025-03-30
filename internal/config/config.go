package config

import (
	"fmt"
	
	"github.com/spf13/viper"
)

type TLSConfig struct {
	Enabled bool   `json:"enabled" yaml:"enabled"`
	Cert    string `json:"cert" yaml:"cert"`
	Key     string `json:"key" yaml:"key"`
}

type Server struct {
	Port int `json:"port" yaml:"port"`
	// Host is the host address for the server.
	Host string `json:"host" yaml:"host"`
	// TLS configuration
	TLS TLSConfig `json:"tls" yaml:"tls"`
}

type MetricsConfig struct {
	Enabled bool   `json:"enabled" yaml:"enabled"`
	Port    int    `json:"port" yaml:"port"`
	Path    string `json:"path" yaml:"path"`
}

type Telemetry struct {
	Metrics MetricsConfig `json:"metrics" yaml:"metrics"`
	Tracing MetricsConfig `json:"tracing" yaml:"tracing"`
}

type LoggingConfig struct {
	Level  string `json:"level" yaml:"level"`
	Format string `json:"format" yaml:"format"`
}

type OpenAIConfig struct {
	APIKey  string `json:"apiKey" yaml:"apiKey"`
	Model   string `json:"model" yaml:"model"`
	BaseURL string `json:"baseUrl" yaml:"baseUrl"`
}

type Config struct {
	// Server configuration
	Server Server `json:"server" yaml:"server"`
	// OpenAI configuration
	OpenAI OpenAIConfig `json:"openai" yaml:"openai"`
	// Download path for generated images
	DownloadPath string `json:"downloadPath" yaml:"downloadPath"`
	// Timeout is the timeout duration for requests.
	Timeout string `json:"timeout" yaml:"timeout"`
	// Debug is a flag to enable or disable debug mode.
	Debug bool `json:"debug" yaml:"debug"`
	// Telemetry configuration
	Telemetry Telemetry `json:"telemetry" yaml:"telemetry"`
	// Logging configuration
	Logging LoggingConfig `json:"logging" yaml:"logging"`
}

// NewConfig creates a new Config with default values.
func NewConfig() *Config {
	config := &Config{}
	return config
}

// LoadConfig loads the configuration from a file and command line flags.
func LoadConfig(v *viper.Viper) (*Config, error) {
	// Set config name and type
	v.SetConfigName(".stablemcp")
	v.SetConfigType("yaml") // Default type, will be overridden if .stablemcp.json is found

	// Set environment variable prefix
	v.SetEnvPrefix("STABLEMCP")
	v.AutomaticEnv() // Read environment variables

	// Search for config in multiple locations (in order of priority)
	v.AddConfigPath("./configs")
	v.AddConfigPath("$HOME/.config")
	v.AddConfigPath("/etc")

	// Set default values to match NewConfig()
	// Server defaults
	v.SetDefault("server.host", "localhost")
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.tls.enabled", false)

	// Timeout default
	v.SetDefault("timeout", "30s")

	// Debug mode default
	v.SetDefault("debug", false)

	// Telemetry defaults
	v.SetDefault("telemetry.metrics.enabled", false)
	v.SetDefault("telemetry.metrics.port", 9090)
	v.SetDefault("telemetry.metrics.path", "/metrics")
	v.SetDefault("telemetry.tracing.enabled", false)
	v.SetDefault("telemetry.tracing.port", 9091)
	v.SetDefault("telemetry.tracing.path", "/traces")

	// Logging defaults
	v.SetDefault("logging.level", "info")
	v.SetDefault("logging.format", "json")

	// Aliases for command line flags (kebab-case) to config keys (dot notation)
	v.RegisterAlias("log-level", "logging.level")
	v.RegisterAlias("port", "server.port")
	v.RegisterAlias("host", "server.host")

	// Check if config file is specified via flag (highest priority)
	if configFile := v.GetString("config"); configFile != "" {
		v.SetConfigFile(configFile)
	}

	// Read config file - don't return error if file not found
	if err := v.ReadInConfig(); err != nil {
		// Only return error if it's not a "config file not found" error
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		// If config file not found, we'll use defaults and command line flags
	}

	// Unmarshal the config into the Config struct
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Handle specific flag-to-config mappings that couldn't be done automatically
	// and ensure command line flags take precedence over config file
	
	// Log level flag handling
	if logLevel := v.GetString("log-level"); logLevel != "" {
		config.Logging.Level = logLevel
	}
	
	// Debug flag handling
	if v.IsSet("debug") {
		config.Debug = v.GetBool("debug")
	}
	
	// Timeout flag handling
	if timeout := v.GetString("timeout"); timeout != "" {
		config.Timeout = timeout
	}
	
	// Server host flag handling
	if host := v.GetString("host"); host != "" {
		config.Server.Host = host
	}
	
	// Server port flag handling
	if v.IsSet("port") {
		config.Server.Port = v.GetInt("port")
	}

	return &config, nil
}
