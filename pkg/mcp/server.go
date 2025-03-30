package mcp

import (
	"encoding/json"
	"log"
	"os"

	"github.com/mkm29/stablemcp/internal/helpers"
	"github.com/mkm29/stablemcp/internal/version"
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

// NewMCPServer creates a new MCP server with the specified name
// and using the version from the version package
func NewMCPServer(name string) *MCPServer {
	return &MCPServer{
		Name:    name,
		Version: version.Version,
		Capabilities: Capabilities{
			Tools: make(map[string]any),
		},
		Decoder: json.NewDecoder(os.Stdin),
		Encoder: json.NewEncoder(os.Stdout),
	}
}

// NewMCPServerWithVersion creates a new MCP server with the specified name and version
func NewMCPServerWithVersion(name, customVersion string) *MCPServer {
	return &MCPServer{
		Name:    name,
		Version: customVersion,
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

// HandleRequest handles incoming JSON-RPC requests
func (s *MCPServer) handleRequest() {
	// Read the request from the decoder
	var req JSONRPCRequest
	if err := s.Decoder.Decode(&req); err != nil {
		log.Println("Error decoding request:", err)
		return
	}

	// Log the request
	log.Println("Received request:", helpers.PrettyJSON(req))

	// Handle the request based on the method
	if req.JSONRPC != "2.0" {
		log.Println("Invalid JSON-RPC version")
		return
	}

	var res any
	var sendResponse bool
	switch req.Method {
	case "initialize":
		sendResponse = true
		res = JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: InitializeResult{
				ProtocolVersion: "1.0",
				ServerInfo: ServerInfo{
					Name:    s.Name,
					Version: s.Version,
				},
				Capabilities: s.Capabilities,
			},
		}
	case "notifications/initialize":
		log.Printf("Server initialized with capabilities: %s\n", helpers.PrettyJSON(s.Capabilities))
		return
	default:
		log.Println("Unknown method:", req.Method)
		return
	}

	if sendResponse {
		// Log the response
		log.Println("Sending response:", helpers.PrettyJSON(res))
		// Send the response
		if err := s.Encoder.Encode(res); err != nil {
			log.Println("Error encoding response:", err)
		}
	}
}

// Handler function
func (s *MCPServer) Handle() {
	// Handle incoming requests
	for {
		s.handleRequest()
	}
}
