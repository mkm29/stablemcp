package mcp

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestJSONRPCRequestMarshaling tests marshaling and unmarshaling of JSONRPCRequest
func TestJSONRPCRequestMarshaling(t *testing.T) {
	// Create a test request
	req := JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  "initialize",
		ID:      "1",
		Params:  map[string]interface{}{"foo": "bar"},
	}

	// Marshal to JSON
	data, err := json.Marshal(req)
	require.NoError(t, err, "Marshaling should not error")

	// Verify the JSON structure
	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	require.NoError(t, err, "Unmarshaling should not error")

	// Check fields
	assert.Equal(t, "2.0", result["jsonrpc"], "JSONRPC field mismatch")
	assert.Equal(t, "initialize", result["method"], "Method field mismatch")
	assert.Equal(t, "1", result["id"], "ID field mismatch")

	// Check the params field
	params, ok := result["params"].(map[string]interface{})
	assert.True(t, ok, "Params should be a map")
	assert.Equal(t, "bar", params["foo"], "Params field mismatch")

	// Test unmarshaling back to a struct
	var newReq JSONRPCRequest
	err = json.Unmarshal(data, &newReq)
	require.NoError(t, err, "Unmarshaling back to struct should not error")
	
	// Verify the unmarshaled struct
	assert.Equal(t, req.JSONRPC, newReq.JSONRPC)
	assert.Equal(t, req.Method, newReq.Method)
	assert.Equal(t, req.ID, newReq.ID)
	
	// Verify params (needs type assertion because it's 'any')
	newParams, ok := newReq.Params.(map[string]interface{})
	assert.True(t, ok, "Unmarshaled params should be a map")
	assert.Equal(t, "bar", newParams["foo"], "Unmarshaled params value mismatch")
}

// TestJSONRPCResponseMarshaling tests marshaling and unmarshaling of JSONRPCResponse
func TestJSONRPCResponseMarshaling(t *testing.T) {
	// Create a test response with a result
	result := map[string]interface{}{"success": true}
	res := JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      "1",
		Result:  result,
	}

	// Marshal to JSON
	data, err := json.Marshal(res)
	require.NoError(t, err, "Marshaling should not error")

	// Verify the JSON structure
	var jsonResult map[string]interface{}
	err = json.Unmarshal(data, &jsonResult)
	require.NoError(t, err, "Unmarshaling should not error")

	// Check fields
	assert.Equal(t, "2.0", jsonResult["jsonrpc"], "JSONRPC field mismatch")
	assert.Equal(t, "1", jsonResult["id"], "ID field mismatch")
	
	// Check result field
	resultMap, ok := jsonResult["result"].(map[string]interface{})
	assert.True(t, ok, "Result should be a map")
	assert.Equal(t, true, resultMap["success"], "Result field mismatch")
	assert.NotContains(t, jsonResult, "error", "Error should be omitted when not present")

	// Test with error instead of result
	errorRes := JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      "1",
		Error: &JSONRPCError{
			Code:    -32600,
			Message: "Invalid Request",
			Data:    "Additional error data",
		},
	}

	// Marshal to JSON
	errorData, err := json.Marshal(errorRes)
	require.NoError(t, err, "Marshaling should not error")

	// Verify the JSON structure
	var errorJsonResult map[string]interface{}
	err = json.Unmarshal(errorData, &errorJsonResult)
	require.NoError(t, err, "Unmarshaling should not error")

	// Check error field
	errorMap, ok := errorJsonResult["error"].(map[string]interface{})
	assert.True(t, ok, "Error should be a map")
	assert.Equal(t, float64(-32600), errorMap["code"], "Error code mismatch")
	assert.Equal(t, "Invalid Request", errorMap["message"], "Error message mismatch")
	assert.Equal(t, "Additional error data", errorMap["data"], "Error data mismatch")
	assert.NotContains(t, errorJsonResult, "result", "Result should be omitted when error is present")
}

// TestInitializeResultMarshaling tests marshaling and unmarshaling of InitializeResult
func TestInitializeResultMarshaling(t *testing.T) {
	// Create test capabilities with tools
	tools := map[string]interface{}{
		"tool1": map[string]interface{}{"version": "1.0"},
		"tool2": map[string]interface{}{"version": "2.0"},
	}
	
	// Create a test initialize result
	initResult := InitializeResult{
		ProtocolVersion: "1.0",
		ServerInfo: ServerInfo{
			Name:    "test-server",
			Version: "1.2.3",
		},
		Capabilities: Capabilities{
			Tools: tools,
		},
	}

	// Marshal to JSON
	data, err := json.Marshal(initResult)
	require.NoError(t, err, "Marshaling should not error")

	// Verify the JSON structure
	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	require.NoError(t, err, "Unmarshaling should not error")

	// Check fields
	assert.Equal(t, "1.0", result["protocolVersion"], "ProtocolVersion field mismatch")
	
	// Check serverInfo
	serverInfo, ok := result["serverInfo"].(map[string]interface{})
	assert.True(t, ok, "ServerInfo should be a map")
	assert.Equal(t, "test-server", serverInfo["name"], "ServerInfo name mismatch")
	assert.Equal(t, "1.2.3", serverInfo["version"], "ServerInfo version mismatch")
	
	// Check capabilities and tools
	capabilities, ok := result["capabilities"].(map[string]interface{})
	assert.True(t, ok, "Capabilities should be a map")
	
	toolsMap, ok := capabilities["tools"].(map[string]interface{})
	assert.True(t, ok, "Tools should be a map")
	
	tool1, ok := toolsMap["tool1"].(map[string]interface{})
	assert.True(t, ok, "Tool1 should be a map")
	assert.Equal(t, "1.0", tool1["version"], "Tool1 version mismatch")
	
	tool2, ok := toolsMap["tool2"].(map[string]interface{})
	assert.True(t, ok, "Tool2 should be a map")
	assert.Equal(t, "2.0", tool2["version"], "Tool2 version mismatch")
}

// TestServerInfoMarshaling tests marshaling and unmarshaling of ServerInfo
func TestServerInfoMarshaling(t *testing.T) {
	// Create a test server info
	serverInfo := ServerInfo{
		Name:    "test-server",
		Version: "1.2.3",
	}

	// Marshal to JSON
	data, err := json.Marshal(serverInfo)
	require.NoError(t, err, "Marshaling should not error")

	// Verify the JSON structure
	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	require.NoError(t, err, "Unmarshaling should not error")

	// Check fields
	assert.Equal(t, "test-server", result["name"], "Name field mismatch")
	assert.Equal(t, "1.2.3", result["version"], "Version field mismatch")
}

// TestCapabilitiesMarshaling tests marshaling and unmarshaling of Capabilities
func TestCapabilitiesMarshaling(t *testing.T) {
	// Create a test capabilities
	tools := map[string]interface{}{
		"tool1": map[string]interface{}{"version": "1.0"},
	}
	
	capabilities := Capabilities{
		Tools: tools,
	}

	// Marshal to JSON
	data, err := json.Marshal(capabilities)
	require.NoError(t, err, "Marshaling should not error")

	// Verify the JSON structure
	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	require.NoError(t, err, "Unmarshaling should not error")

	// Check tools
	toolsMap, ok := result["tools"].(map[string]interface{})
	assert.True(t, ok, "Tools should be a map")
	
	tool1, ok := toolsMap["tool1"].(map[string]interface{})
	assert.True(t, ok, "Tool1 should be a map")
	assert.Equal(t, "1.0", tool1["version"], "Tool1 version mismatch")

	// Test empty tools - since it has the omitempty tag, 
	// an empty map won't be included in the JSON output
	emptyCapabilities := Capabilities{
		Tools: make(map[string]interface{}),
	}
	
	emptyData, err := json.Marshal(emptyCapabilities)
	require.NoError(t, err, "Marshaling should not error")
	
	// The JSON should be an empty object since tools has omitempty and is empty
	assert.Equal(t, "{}", string(emptyData), "Empty capabilities should be serialized as an empty JSON object")
	
	// Test nil tools - also shouldn't be in the output
	nilCapabilities := Capabilities{}
	nilData, err := json.Marshal(nilCapabilities)
	require.NoError(t, err, "Marshaling should not error")
	
	// The JSON should also be an empty object
	assert.Equal(t, "{}", string(nilData), "Nil capabilities should be serialized as an empty JSON object")
}