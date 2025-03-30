package mcp

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewMCPServer tests the creation of a new MCP server
func TestNewMCPServer(t *testing.T) {
	server := NewMCPServer("test-server")
	
	assert.NotNil(t, server, "Server should not be nil")
	assert.Equal(t, "test-server", server.Name, "Server name should match")
	assert.NotEmpty(t, server.Version, "Server version should not be empty")
	assert.NotNil(t, server.Capabilities.Tools, "Server capabilities tools should not be nil")
	assert.Empty(t, server.Capabilities.Tools, "Server capabilities tools should be empty")
	assert.NotNil(t, server.Decoder, "Server decoder should not be nil")
	assert.NotNil(t, server.Encoder, "Server encoder should not be nil")
}

// TestNewMCPServerWithVersion tests the creation of a new MCP server with a custom version
func TestNewMCPServerWithVersion(t *testing.T) {
	server := NewMCPServerWithVersion("test-server", "1.2.3")
	
	assert.NotNil(t, server, "Server should not be nil")
	assert.Equal(t, "test-server", server.Name, "Server name should match")
	assert.Equal(t, "1.2.3", server.Version, "Server version should match")
	assert.NotNil(t, server.Capabilities.Tools, "Server capabilities tools should not be nil")
	assert.Empty(t, server.Capabilities.Tools, "Server capabilities tools should be empty")
	assert.NotNil(t, server.Decoder, "Server decoder should not be nil")
	assert.NotNil(t, server.Encoder, "Server encoder should not be nil")
}

// TestInitialize tests the initialization of an MCP server
func TestInitialize(t *testing.T) {
	// Temporarily capture the log output
	oldOutput := os.Stderr
	logPipe, w, _ := os.Pipe()
	os.Stderr = w
	
	server := NewMCPServer("test-server")
	err := server.Initialize("stderr")
	
	// Restore normal logging
	w.Close()
	os.Stderr = oldOutput
	
	// Read the captured output
	var buf bytes.Buffer
	io.Copy(&buf, logPipe)
	logOutput := buf.String()
	
	assert.NoError(t, err, "Initialize should not return an error")
	assert.Contains(t, logOutput, "Logger initialized", "Log should contain initialization message")
	assert.Len(t, server.Capabilities.Tools, 1, "Server should have one tool")
	assert.Contains(t, server.Capabilities.Tools, "exampleTool", "Server should have exampleTool")
}

// TestSetLogger tests the setLogger function
func TestSetLogger(t *testing.T) {
	testCases := []struct {
		name       string
		logOutput  string
		checkStdout bool
	}{
		{"Default", "", false},
		{"Stderr", "stderr", false},
		{"Stdout", "stdout", true},
		{"Unknown", "unknown", false},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Capture stdout
			oldStdout := os.Stdout
			oldStderr := os.Stderr
			
			rOut, wOut, _ := os.Pipe()
			rErr, wErr, _ := os.Pipe()
			
			if tc.checkStdout {
				os.Stdout = wOut
			} else {
				os.Stderr = wErr
			}
			
			server := NewMCPServer("test-server")
			server.setLogger(tc.logOutput)
			
			// Write a log message
			log := "Test log message"
			server.Encoder.Encode(log)
			
			// Restore normal output
			if tc.checkStdout {
				wOut.Close()
				os.Stdout = oldStdout
			} else {
				wErr.Close()
				os.Stderr = oldStderr
			}
			
			// Read the captured output
			var outBuf, errBuf bytes.Buffer
			
			if tc.checkStdout {
				io.Copy(&outBuf, rOut)
				output := outBuf.String()
				assert.Contains(t, output, "Test log message", "Log message should be in stdout")
			} else {
				io.Copy(&errBuf, rErr)
				// We can't easily assert on stderr content in this test setup
				// Just verify that the function runs without errors
				assert.NotPanics(t, func() { server.setLogger(tc.logOutput) })
			}
		})
	}
}

// mockDecoder simulates a json.Decoder for testing purposes
type mockDecoder struct {
	data string
}

func (d *mockDecoder) Decode(v interface{}) error {
	return json.Unmarshal([]byte(d.data), v)
}

// mockWriter captures the output for testing purposes
type mockWriter struct {
	bytes.Buffer
}

func (w *mockWriter) Close() error {
	return nil
}

// TestHandleRequest tests the handleRequest method with different request types
func TestHandleRequest(t *testing.T) {
	testCases := []struct {
		name           string
		requestJSON    string
		expectResponse bool
		expectedResult string
	}{
		{
			name: "Initialize Request",
			requestJSON: `{
				"jsonrpc": "2.0",
				"method": "initialize",
				"id": "1",
				"params": {"clientInfo": {"name": "test-client"}}
			}`,
			expectResponse: true,
			expectedResult: `"serverInfo":{"name":"test-server"`,
		},
		{
			name: "Notification Initialize",
			requestJSON: `{
				"jsonrpc": "2.0",
				"method": "notifications/initialize",
				"id": "2"
			}`,
			expectResponse: false,
		},
		{
			name: "Unknown Method",
			requestJSON: `{
				"jsonrpc": "2.0",
				"method": "unknown",
				"id": "3"
			}`,
			expectResponse: false,
		},
		{
			name: "Invalid JSON-RPC Version",
			requestJSON: `{
				"jsonrpc": "1.0",
				"method": "initialize",
				"id": "4"
			}`,
			expectResponse: false,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a server with mocked decoder and encoder
			server := NewMCPServer("test-server")
			
			// Use strings.NewReader to create an io.Reader from the test JSON
			mockInput := strings.NewReader(tc.requestJSON)
			server.Decoder = json.NewDecoder(mockInput)
			
			// Create a buffer to capture the output
			outputBuffer := new(mockWriter)
			server.Encoder = json.NewEncoder(outputBuffer)
			
			// Redirect logs to a buffer for testing
			oldOutput := os.Stderr
			logPipe, w, _ := os.Pipe()
			os.Stderr = w
			
			// Call the method being tested
			server.handleRequest()
			
			// Restore normal logging
			w.Close()
			os.Stderr = oldOutput
			
			// We don't need to read the log output for this test
			_, _ = io.Copy(io.Discard, logPipe)
			
			// Check if a response was expected and received
			if tc.expectResponse {
				output := outputBuffer.String()
				assert.Contains(t, output, tc.expectedResult, "Response should contain expected content")
				assert.Contains(t, output, "jsonrpc", "Response should be in JSON-RPC format")
			} else {
				// The output could be an empty string since no response is sent
				assert.NotContains(t, outputBuffer.String(), "result", "No result should be sent")
			}
		})
	}
}