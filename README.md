# StableMCP

A Model Context Protocol (MCP) server for generating images using Stable Diffusion.

## Features

- MCP-compliant API endpoint for image generation
- Integration with Stable Diffusion image generation models
- Support for various image parameters (size, style, prompt)
- API key authentication (optional)
- Configurable image size and quality settings
- Rate limiting and request validation

## Project Structure

```
.
├── api           # API definitions, routes, and documentation
├── bin           # Build artifacts
├── cmd           # Application entrypoints
├── configs       # Configuration files
├── examples      # Example usages
├── internal      # Private application code
│   ├── config    # Configuration handling
│   ├── models    # Data models
│   └── utils     # Utility functions
├── pkg           # Public packages
│   ├── auth      # Authentication/authorization
│   ├── handlers  # Request handlers
│   ├── mcp       # MCP protocol implementation
│   └── stablediffusion # Stable Diffusion client
└── scripts       # Utility scripts
```

## Prerequisites

- Go 1.22 or higher
- A running Stable Diffusion API (local or remote)

## Getting Started

```bash
# Clone the repository
git clone https://github.com/yourusername/stablemcp.git
cd stablemcp

# Build the server
go build -o bin/stablemcp ./main.go

# Run with default config
./bin/stablemcp server

# Run with a custom config
./bin/stablemcp server --config configs/custom.yaml
```

## Configuration

The application uses [Viper](https://github.com/spf13/viper) for configuration management. Configuration values can be provided via:

1. Custom config file specified with the `--config` flag (highest priority)
2. Configuration files named `.stablemcp.yaml` or `.stablemcp.json` in standard locations (checked in this order):
   - `./configs/.stablemcp.yaml` (in the current directory)
   - `$HOME/.config/.stablemcp.yaml` (in the user's home directory)
   - `/etc/.stablemcp.yaml` (system-wide)
3. Default values (lowest priority)

### Configuration Options

You can customize the application by setting the following options in your YAML configuration file:

```yaml
server:
  host: "localhost"                    # Server host address (default: "localhost")
  port: 8080                           # Server port (default: 8080)
  tls:
    enabled: true                      # Enable/disable TLS (default: false)
    # There are no default values for the following options, so these must be set if TLS is enabled
    cert: "/path/to/cert.pem"          # TLS certificate path (default: "")
    key: "/path/to/key.pem"            # TLS key path (default: "")

logging:
  level: "info"                        # Log level: debug, info, warn, error (default: "info")
  format: "json"                       # Log format: json, text (default: "json")

debug: false                           # Enable debug mode (default: false)
timeout: "30s"                         # Request timeout (default: "30s")

telemetry:
  metrics:
    enabled: false                     # Enable metrics collection (default: false)
    port: 9090                         # Metrics server port (default: 9090)
    path: "/metrics"                   # Metrics endpoint path (default: "/metrics")
  tracing:
    enabled: false                     # Enable distributed tracing (default: false)
    port: 9091                         # Tracing server port (default: 9091)
    path: "/traces"                    # Tracing endpoint path (default: "/traces")

# OpenAI configuration
openai:
  apiKey: "your-openai-api-key"        # OpenAI API key for API calls (default: "")
  model: "chatgpt-4o"                  # Model to use (default: "chatgpt-4o")
  baseUrl: "https://api.openai.com/v1" # Base URL for API calls

# download path for generated images
downloadPath: "/path/to/downloads"     # Path where generated images will be saved (default: "~/Downloads")
```

Any values not specified in your configuration file will use the defaults shown above.

### Using a Custom Configuration File

You can specify a custom configuration file when running the server:

```bash
./bin/stablemcp server --config path/to/your/config.yaml
```

Or create one of these standard configuration files:

```bash
# In your project directory
touch configs/.stablemcp.yaml

# In your home directory
touch ~/.config/.stablemcp.yaml

# System-wide (requires sudo)
sudo touch /etc/.stablemcp.yaml
```

## API Usage

### Generate an Image

```bash
curl -X POST http://localhost:8080/v1/generate \
  -H "Content-Type: application/json" \
  -d '{
    "prompt": "a photo of a cat in space",
    "width": 512,
    "height": 512,
    "num_inference_steps": 50
  }'
```

## Development

```bash
# Run tests
go test ./...

# Format code
go fmt ./...

# Run linter
golangci-lint run
```

## License

[MIT License](LICENSE)