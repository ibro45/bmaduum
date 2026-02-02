package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version information - defaults are overridden by ldflags during build
var (
	// Version is the current semantic version (set via ldflags)
	Version = "dev"
	// Commit is the git commit hash (set via ldflags)
	Commit = "unknown"
	// Date is the build date (set via ldflags)
	Date = "unknown"
	// BuiltBy is the builder identifier (set via ldflags)
	BuiltBy = "unknown"
)

// SetVersionInfo sets the version information from build-time ldflags.
// This is called by main() before Execute().
func SetVersionInfo(version, commit, date, builtBy string) {
	if version != "" {
		Version = version
	}
	if commit != "" {
		Commit = commit
	}
	if date != "" {
		Date = date
	}
	if builtBy != "" {
		BuiltBy = builtBy
	}
}

func newVersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Display version information",
		Long:  `Display the version, release date, and other build information for bmaduum.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Print version
			cmd.Printf("bmaduum version %s\n", Version)

			// Print additional build info if available
			if Commit != "unknown" {
				cmd.Printf("commit: %s\n", Commit)
			}
			if Date != "unknown" {
				cmd.Printf("built at: %s\n", Date)
			}
			if BuiltBy != "unknown" {
				cmd.Printf("built by: %s\n", BuiltBy)
			}
		},
	}

	return cmd
}

// GetVersion returns the version string.
// Useful for other packages that need to display the version.
func GetVersion() string {
	return Version
}

// GetVersionInfo returns all version information as a map.
// Useful for debugging and testing.
func GetVersionInfo() map[string]string {
	return map[string]string{
		"version": Version,
		"commit":  Commit,
		"date":    Date,
		"builtBy": BuiltBy,
	}
}

// FormatVersion formats version information as a single-line string.
// Useful for logging and user-facing messages.
func FormatVersion() string {
	if Commit != "unknown" {
		return fmt.Sprintf("%s (%s)", Version, Commit[:7])
	}
	return Version
}
