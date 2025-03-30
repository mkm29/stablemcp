package version

import (
	"testing"
)

func TestInfo(t *testing.T) {
	// Save the original values
	origVersion := Version
	origBuildDate := BuildDate
	origGitCommit := GitCommit
	origGitBranch := GitBranch

	// Restore the original values after the test
	defer func() {
		Version = origVersion
		BuildDate = origBuildDate
		GitCommit = origGitCommit
		GitBranch = origGitBranch
	}()

	// Set test values
	Version = "1.2.3"
	BuildDate = "2025-04-01T12:00:00Z"
	GitCommit = "abcdef123456"
	GitBranch = "main"

	// Call the function
	info := Info()

	// Assert the results
	expected := map[string]string{
		"version":   "1.2.3",
		"buildDate": "2025-04-01T12:00:00Z",
		"gitCommit": "abcdef123456",
		"gitBranch": "main",
	}

	// Check if all keys and values match
	for k, v := range expected {
		if info[k] != v {
			t.Errorf("Expected %s to be %s, but got %s", k, v, info[k])
		}
	}

	// Check if map sizes match
	if len(info) != len(expected) {
		t.Errorf("Expected map with %d entries, but got %d", len(expected), len(info))
	}
}