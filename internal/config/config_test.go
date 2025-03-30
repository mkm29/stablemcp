package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	config := NewConfig()
	assert.NotNil(t, config, "NewConfig should return a non-nil config")
}

func TestLoadConfigWithDefaults(t *testing.T) {
	v := viper.New()
	config, err := LoadConfig(v)
	
	assert.NoError(t, err, "LoadConfig should not return an error with default values")
	assert.NotNil(t, config, "LoadConfig should return a non-nil config")
	
	// Check default values
	assert.Equal(t, "localhost", config.Server.Host)
	assert.Equal(t, 8080, config.Server.Port)
	assert.Equal(t, false, config.Server.TLS.Enabled)
	assert.Equal(t, "30s", config.Timeout)
	assert.Equal(t, false, config.Debug)
	assert.Equal(t, "info", config.Logging.Level)
	assert.Equal(t, "json", config.Logging.Format)
}

func TestLoadConfigWithFile(t *testing.T) {
	// Create a temporary config file
	dir, err := os.MkdirTemp("", "config-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)
	
	configPath := filepath.Join(dir, ".stablemcp.yaml")
	
	// Write test config
	configContent := `
server:
  host: "testhost"
  port: 9000
  tls:
    enabled: true
    cert: "cert-path"
    key: "key-path"
logging:
  level: "debug"
  format: "text"
debug: true
timeout: "60s"
downloadPath: "/test/path"
`
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}
	
	// Create viper instance and load config
	v := viper.New()
	v.SetConfigFile(configPath)
	
	// Force read the config file since we're providing a specific path
	err = v.ReadInConfig()
	assert.NoError(t, err, "ReadInConfig should not return an error with valid config file")
	
	config, err := LoadConfig(v)
	assert.NoError(t, err, "LoadConfig should not return an error with valid config file")
	assert.NotNil(t, config, "LoadConfig should return a non-nil config")
	
	// Check values from file
	assert.Equal(t, "testhost", config.Server.Host)
	assert.Equal(t, 9000, config.Server.Port)
	assert.Equal(t, true, config.Server.TLS.Enabled)
	assert.Equal(t, "cert-path", config.Server.TLS.Cert)
	assert.Equal(t, "key-path", config.Server.TLS.Key)
	assert.Equal(t, "debug", config.Logging.Level)
	assert.Equal(t, "text", config.Logging.Format)
	assert.Equal(t, true, config.Debug)
	assert.Equal(t, "60s", config.Timeout)
	assert.Equal(t, "/test/path", config.DownloadPath)
}

func TestLoadConfigWithFlags(t *testing.T) {
	// Create a temporary config file
	dir, err := os.MkdirTemp("", "config-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)
	
	configPath := filepath.Join(dir, ".stablemcp.yaml")
	
	// Write test config
	configContent := `
server:
  host: "confighost"
  port: 8000
logging:
  level: "info"
debug: false
`
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}
	
	// Create viper instance and set flag values
	v := viper.New()
	v.SetConfigFile(configPath)
	
	// Force read the config file first
	err = v.ReadInConfig()
	assert.NoError(t, err, "ReadInConfig should not return an error with valid config file")
	
	// Set values via flags (these should override the config file)
	v.Set("server.host", "flaghost")      // Override host
	v.Set("server.port", 9090)           // Override port
	v.Set("logging.level", "warning")    // Override log level
	v.Set("debug", true)                 // Override debug
	
	config, err := LoadConfig(v)
	assert.NoError(t, err, "LoadConfig should not return an error with valid flags")
	assert.NotNil(t, config, "LoadConfig should return a non-nil config")
	
	// Check that flag values take precedence
	assert.Equal(t, "flaghost", config.Server.Host)
	assert.Equal(t, 9090, config.Server.Port)
	assert.Equal(t, "warning", config.Logging.Level)
	assert.Equal(t, true, config.Debug)
}

func TestLoadConfigWithInvalidFile(t *testing.T) {
	// Create a temporary invalid config file
	dir, err := os.MkdirTemp("", "config-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)
	
	configPath := filepath.Join(dir, ".stablemcp.yaml")
	
	// Write syntactically valid YAML but with invalid data (string where int expected)
	// Viper won't catch this at read time but should fail at unmarshal time
	invalidContent := `
server:
  host: "testhost"
  port: "not-a-number"  # This will cause an error when unmarshaling to int
`
	err = os.WriteFile(configPath, []byte(invalidContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}
	
	// Create viper instance and load config
	v := viper.New()
	v.SetConfigFile(configPath)
	
	// Reading the file should succeed because the YAML is syntactically valid
	err = v.ReadInConfig()
	assert.NoError(t, err, "ReadInConfig should not return an error with valid YAML syntax")
	
	// But unmarshaling should fail because "not-a-number" can't be converted to int
	_, err = LoadConfig(v)
	assert.Error(t, err, "LoadConfig should return an error with semantically invalid config file")
	assert.Contains(t, err.Error(), "error unmarshaling", "Error should be related to unmarshaling")
}

func TestLoadConfigWithNonExistentFile(t *testing.T) {
	// Create viper instance with non-existent file
	v := viper.New()
	v.SetConfigFile("/path/that/does/not/exist.yaml")
	
	// This should use defaults
	config, err := LoadConfig(v)
	assert.NoError(t, err, "LoadConfig should not return an error with non-existent file")
	assert.NotNil(t, config, "LoadConfig should return a non-nil config with defaults")
	
	// Check some default values
	assert.Equal(t, "localhost", config.Server.Host)
	assert.Equal(t, 8080, config.Server.Port)
}