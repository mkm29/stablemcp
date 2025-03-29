package config

import "github.com/spf13/viper"

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
	config := &Config{
		Server: Server{
			Port: 8080,
			Host: "localhost",
			TLS: TLSConfig{
				Enabled: false,
			},
		},
		Timeout: "30s",
		Debug:   false,
		Telemetry: Telemetry{
			Metrics: MetricsConfig{
				Enabled: false,
			},
			Tracing: MetricsConfig{
				Enabled: false,
			},
		},
		Logging: LoggingConfig{
			Level:  "info",
			Format: "json",
		},
		OpenAI: OpenAIConfig{
			Model:   "gpt-3.5-turbo",
			BaseURL: "https://api.openai.com/v1",
		},
		DownloadPath: "~/Downloads",
	}
	return config
}

// LoadConfig loads the configuration from a file.
func LoadConfig(v *viper.Viper) (*Config, error) {
	v.SetConfigName(".stablemcp")
	v.SetConfigType("yaml") // Default type, will be overridden if .stablemcp.json is found

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

	// Check if config file is specified via flag (highest priority)
	if configFile := v.GetString("config"); configFile != "" {
		v.SetConfigFile(configFile)
	}

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	// Unmarshal the config into the Config struct
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
