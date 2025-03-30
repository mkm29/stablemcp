package helpers

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestPrettyJSON tests the PrettyJSON function
func TestPrettyJSON(t *testing.T) {
	testCases := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:  "Simple Object",
			input: map[string]string{"key": "value"},
			expected: `{
  "key": "value"
}`,
		},
		{
			name:  "Nested Object",
			input: map[string]interface{}{"key": map[string]int{"nested": 42}},
			expected: `{
  "key": {
    "nested": 42
  }
}`,
		},
		{
			name:  "Array",
			input: []string{"one", "two", "three"},
			expected: `[
  "one",
  "two",
  "three"
]`,
		},
		{
			name:  "Mixed Types",
			input: map[string]interface{}{"string": "value", "number": 42, "bool": true, "null": nil},
			expected: `{
  "bool": true,
  "null": null,
  "number": 42,
  "string": "value"
}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := PrettyJSON(tc.input)
			
			// Normalize line endings for platform independence
			expectedNormalized := strings.ReplaceAll(tc.expected, "\r\n", "\n")
			resultNormalized := strings.ReplaceAll(result, "\r\n", "\n")
			
			assert.Equal(t, expectedNormalized, resultNormalized, "PrettyJSON output should match expected format")
		})
	}
}

// TestPrettyJSONError tests the error handling in PrettyJSON
func TestPrettyJSONError(t *testing.T) {
	// Create a value that cannot be marshaled to JSON
	badValue := make(chan int) // Channels cannot be marshaled to JSON
	
	result := PrettyJSON(badValue)
	assert.Contains(t, result, "Error:", "PrettyJSON should return an error message for unmarshalable values")
	assert.Contains(t, result, "json", "PrettyJSON error should mention JSON")
}

// TestPrettyJSONWithStruct tests PrettyJSON with a custom struct
func TestPrettyJSONWithStruct(t *testing.T) {
	type TestStruct struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}
	
	testStruct := TestStruct{
		Name:  "test",
		Value: 42,
	}
	
	result := PrettyJSON(testStruct)
	expected := `{
  "name": "test",
  "value": 42
}`
	
	// Normalize line endings for platform independence
	expectedNormalized := strings.ReplaceAll(expected, "\r\n", "\n")
	resultNormalized := strings.ReplaceAll(result, "\r\n", "\n")
	
	assert.Equal(t, expectedNormalized, resultNormalized, "PrettyJSON should correctly format struct")
}

// TestPrettyJSONWithJSONTags tests PrettyJSON respects JSON struct tags
func TestPrettyJSONWithJSONTags(t *testing.T) {
	type TaggedStruct struct {
		Name      string `json:"customName"`
		Value     int    `json:"customValue"`
		Omitted   string `json:"-"`
		OmitEmpty string `json:"omitMe,omitempty"`
	}
	
	testStruct := TaggedStruct{
		Name:      "test",
		Value:     42,
		Omitted:   "this should not appear in JSON",
		OmitEmpty: "",
	}
	
	result := PrettyJSON(testStruct)
	expected := `{
  "customName": "test",
  "customValue": 42
}`
	
	// Normalize line endings for platform independence
	expectedNormalized := strings.ReplaceAll(expected, "\r\n", "\n")
	resultNormalized := strings.ReplaceAll(result, "\r\n", "\n")
	
	assert.Equal(t, expectedNormalized, resultNormalized, "PrettyJSON should respect JSON struct tags")
	assert.NotContains(t, result, "Omitted", "Fields with json:\"-\" should be omitted")
	assert.NotContains(t, result, "omitMe", "Fields with omitempty and empty values should be omitted")
	
	// Now test with a non-empty OmitEmpty field
	testStruct.OmitEmpty = "now I appear"
	result = PrettyJSON(testStruct)
	assert.Contains(t, result, "omitMe", "Fields with omitempty and non-empty values should be included")
	assert.Contains(t, result, "now I appear", "The value should be included")
}