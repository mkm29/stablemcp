package version

// Build-time variables
var (
	// Version is the current version of StableMCP
	Version = "0.1.1"
	// BuildDate is the date when StableMCP was built
	BuildDate = "unknown"
	// GitCommit is the git commit hash of the build
	GitCommit = "unknown"
	// GitBranch is the git branch of the build
	GitBranch = "unknown"
)

// Info returns a map with version information
func Info() map[string]string {
	return map[string]string{
		"version":    Version,
		"buildDate":  BuildDate,
		"gitCommit":  GitCommit,
		"gitBranch":  GitBranch,
	}
}