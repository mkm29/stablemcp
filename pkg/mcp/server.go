package mcp

import (
	"encoding/json"
	"log"
	"os"
)

type MCPServer struct {
	// The server's name
	Name string
	// The server's version
	Version string
	// The server's capabilities
	Capabilities Capabilities
	// JSON Decoder
	Decoder *json.Decoder
	// JSON Encoder
	Encoder *json.Encoder
}

func NewMCPServer(name, version string) *MCPServer {
	return &MCPServer{
		Name:    name,
		Version: version,
		Capabilities: Capabilities{
			Tools: make(map[string]any),
		},
		Decoder: json.NewDecoder(os.Stdin),
		Encoder: json.NewEncoder(os.Stdout),
	}
}

func (s *MCPServer) setLogger(output string) {
	// Set up logging here
	// For example, you can use logrus or any other logging library
	// to log the requests and responses.
	// By default, we will log to stderr.

	if output == "" {
		output = "stderr"
	}
	switch output {
	case "stdout":
		log.SetOutput(os.Stdout)
	case "stderr":
		log.SetOutput(os.Stderr)
	default:
		log.SetOutput(os.Stderr)
	}
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix("MCPServer: ")
	log.Println("Logger initialized")
}

func (s *MCPServer) Initialize(logOutput string) error {
	// Initialize the server
	// setup logging
	s.setLogger(logOutput)

	// Initialize the server's capabilities
	// For example, you can set up the server's capabilities here.
	s.Capabilities.Tools["exampleTool"] = map[string]any{
		"version": "1.0",
	}
	return nil
}
