package cli

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVersionCommand(t *testing.T) {
	// Save original values and restore after test
	origVersion := Version
	origCommit := Commit
	defer func() {
		Version = origVersion
		Commit = origCommit
	}()

	// Set test values
	Version = "test-1.2.3"
	Commit = "abc123def"

	cmd := newVersionCommand()

	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	err := cmd.Execute()
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "bmaduum version test-1.2.3")
	assert.Contains(t, output, "commit: abc123def")
}

func TestVersionCommand_Defaults(t *testing.T) {
	// Save original values and restore after test
	origVersion := Version
	origCommit := Commit
	origDate := Date
	origBuiltBy := BuiltBy
	defer func() {
		Version = origVersion
		Commit = origCommit
		Date = origDate
		BuiltBy = origBuiltBy
	}()

	// Set to default values (simulating a dev build)
	Version = "dev"
	Commit = "unknown"
	Date = "unknown"
	BuiltBy = "unknown"

	cmd := newVersionCommand()

	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	err := cmd.Execute()
	require.NoError(t, err)

	output := buf.String()
	// Should only show version line for dev build
	assert.Contains(t, output, "bmaduum version dev")
	// Should not show commit/date for unknown values
	assert.NotContains(t, output, "commit: unknown")
	assert.NotContains(t, output, "built at: unknown")
}

func TestVersionCommand_NoArgs(t *testing.T) {
	cmd := newVersionCommand()

	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	// Version command should work with no args
	err := cmd.Execute()
	assert.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "bmaduum version")
}

func TestSetVersionInfo(t *testing.T) {
	// Save original values and restore after test
	origVersion := Version
	origCommit := Commit
	origDate := Date
	origBuiltBy := BuiltBy
	defer func() {
		Version = origVersion
		Commit = origCommit
		Date = origDate
		BuiltBy = origBuiltBy
	}()

	// Reset to defaults
	Version = "dev"
	Commit = "unknown"
	Date = "unknown"
	BuiltBy = "unknown"

	// Call SetVersionInfo
	SetVersionInfo("1.2.3", "commit456", "2024-01-15", "goreleaser")

	// Verify values were set
	assert.Equal(t, "1.2.3", Version)
	assert.Equal(t, "commit456", Commit)
	assert.Equal(t, "2024-01-15", Date)
	assert.Equal(t, "goreleaser", BuiltBy)
}

func TestGetVersion(t *testing.T) {
	// Save and restore
	origVersion := Version
	defer func() { Version = origVersion }()

	Version = "v2.0.0"
	assert.Equal(t, "v2.0.0", GetVersion())
}

func TestGetVersionInfo(t *testing.T) {
	// Save and restore
	origVersion := Version
	origCommit := Commit
	origDate := Date
	origBuiltBy := BuiltBy
	defer func() {
		Version = origVersion
		Commit = origCommit
		Date = origDate
		BuiltBy = origBuiltBy
	}()

	Version = "v3.0.0"
	Commit = "abc123"
	Date = "2024-01-01"
	BuiltBy = "user"

	info := GetVersionInfo()
	assert.Equal(t, "v3.0.0", info["version"])
	assert.Equal(t, "abc123", info["commit"])
	assert.Equal(t, "2024-01-01", info["date"])
	assert.Equal(t, "user", info["builtBy"])
}

func TestFormatVersion(t *testing.T) {
	// Save and restore
	origVersion := Version
	origCommit := Commit
	defer func() {
		Version = origVersion
		Commit = origCommit
	}()

	t.Run("with commit", func(t *testing.T) {
		Version = "v1.0.0"
		Commit = "abcdef123456789"
		assert.Equal(t, "v1.0.0 (abcdef1)", FormatVersion())
	})

	t.Run("without commit", func(t *testing.T) {
		Version = "v2.0.0"
		Commit = "unknown"
		assert.Equal(t, "v2.0.0", FormatVersion())
	})
}
