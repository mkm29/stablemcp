package config

import "github.com/spf13/viper"

type Config struct {
	// ServerPort is the port on which the server will listen.
	ServerPort int `json:"server_port" yaml:"server_port"`
	// LogLevel is the logging level for the application.
	LogLevel string `json:"log_level" yaml:"log_level"`
	// Timeout is the timeout duration for requests.
	Timeout string `json:"timeout" yaml:"timeout"`
	// EnableDebugMode is a flag to enable or disable debug mode.
	EnableDebugMode bool `json:"enable_debug_mode" yaml:"enable_debug_mode"`
	// EnableMetrics is a flag to enable or disable metrics collection.
	EnableMetrics bool `json:"enable_metrics" yaml:"enable_metrics"`
	// MetricsPort is the port for the metrics endpoint.
	MetricsPort int `json:"metrics_port" yaml:"metrics_port"`
	// EnableTracing is a flag to enable or disable tracing.
	EnableTracing bool `json:"enable_tracing" yaml:"enable_tracing"`
	// TracingEndpoint is the endpoint for the tracing service.
	TracingEndpoint string `json:"tracing_endpoint" yaml:"tracing_endpoint"`
}

// NewConfig creates a new Config with default values.
func NewConfig() *Config {
	return &Config{
		ServerPort:      8080,
		LogLevel:        "info",
		Timeout:         "30s",
		EnableDebugMode: false,
		EnableMetrics:   false,
		EnableTracing:   false,
	}
}

func (c *Config) IsDebug() bool {
	return c.EnableDebugMode
}

// LoadConfig loads the configuration from a file.
func LoadConfig(v *viper.Viper) (*Config, error) {
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	
	// Search for config in multiple locations (in order of priority)
	v.AddConfigPath("./configs")
	v.AddConfigPath("$HOME/.config/stablemcp")
	v.AddConfigPath("/etc/stablemcp")
	
	// Set default values
	v.SetDefault("server_port", 8080)
	v.SetDefault("log_level", "info")
	v.SetDefault("timeout", "30s")
	v.SetDefault("enable_debug_mode", false)
	v.SetDefault("enable_metrics", false)
	v.SetDefault("enable_tracing", false)

	// Check if config file is specified via flag (highest priority)
	if configFile := v.GetString("config"); configFile != "" {
		v.SetConfigFile(configFile)
	}

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}
	
	// Unmarshal the config into the Config struct
	// This will automatically handle the conversion from JSON/YAML to the struct fields
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
