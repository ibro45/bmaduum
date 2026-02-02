package main

import "bmaduum/internal/cli"

// Version information, set via ldflags during build
var (
	// Version is the semantic version (e.g., "1.2.3")
	version = "dev"
	// Commit is the git commit hash (short form)
	commit = "unknown"
	// Date is the build timestamp (RFC3339 or similar)
	date = "unknown"
	// BuiltBy identifies the builder (e.g., "goreleaser")
	builtBy = "unknown"
)

func main() {
	// Inject version information into CLI package
	cli.SetVersionInfo(version, commit, date, builtBy)
	cli.Execute()
}
