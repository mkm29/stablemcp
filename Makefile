# Makefile for StableMCP

# Variables
BINARY_NAME=stablemcp
BUILD_DIR=bin
GO=$(shell command -v go)
GOFLAGS=-ldflags="-s -w"
GOLANGCI_LINT=$(shell command -v golangci-lint)
CONFIG_PATH=configs/.stablemcp.yaml

# Default target
.PHONY: all
all: clean lint test build

# Build the application
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@$(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) main.go

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)

# Run the server
.PHONY: run
run: build
	@echo "Running server..."
	@./$(BUILD_DIR)/$(BINARY_NAME) server

# Run the server with a custom config
.PHONY: run-config
run-config: build
	@echo "Running server with custom config..."
	@./$(BUILD_DIR)/$(BINARY_NAME) server --config $(CONFIG_PATH)

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	@$(GO) test ./... -v

# Run a specific test
.PHONY: test-one
test-one:
	@echo "Running specific test $(TEST)..."
	@$(GO) test ./... -v -run $(TEST)

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	@$(GO) fmt ./...

# Run linter
.PHONY: lint
lint:
	@echo "Running linter..."
	@if [ -z "$(GOLANGCI_LINT)" ]; then \
		echo "golangci-lint is not installed. Installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	@golangci-lint run

# Create a default config file
.PHONY: config
config:
	@echo "Creating default config file..."
	@mkdir -p configs
	@cat > configs/.stablemcp.yaml << EOF
server:
  host: "localhost"
  port: 8080
  tls:
    enabled: false
    cert: ""
    key: ""

logging:
  level: "info"
  format: "json"

debug: false
timeout: "30s"

telemetry:
  metrics:
    enabled: false
    port: 9090
    path: "/metrics"
  tracing:
    enabled: false
    port: 9091
    path: "/traces"

openai:
  apiKey: ""
  model: "gpt-3.5-turbo"
  baseUrl: "https://api.openai.com/v1"

downloadPath: "~/Downloads"
EOF
	@echo "Config file created at configs/.stablemcp.yaml"

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all          - Clean, run linter, run tests, and build"
	@echo "  build        - Build the application"
	@echo "  clean        - Remove build artifacts"
	@echo "  run          - Build and run the server"
	@echo "  run-config   - Build and run the server with a custom config"
	@echo "  test         - Run all tests"
	@echo "  test-one     - Run a specific test (usage: make test-one TEST=TestName)"
	@echo "  fmt          - Format code"
	@echo "  lint         - Run linter"
	@echo "  config       - Create a default config file"
	@echo "  help         - Show this help information"