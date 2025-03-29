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
go build -o bin/stablemcp cmd/server.go

# Run with default config
./bin/stablemcp

# Run with a custom config
./bin/stablemcp --config configs/custom.yaml
```

## Configuration

Edit `configs/default.yaml` to configure:

- Server host and port
- Stable Diffusion API endpoint
- Authentication settings
- Logging options

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