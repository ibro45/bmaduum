package cli

import (
	"os"
	"path/filepath"
	"testing"

	"bmad-automate/internal/config"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSprintCommand(t *testing.T) {
	cfg := &config.Config{}
	app := NewApp(cfg)
	cmd := newSprintCommand(app)

	assert.Equal(t, "sprint", cmd.Use)
	assert.NotEmpty(t, cmd.Short)
	assert.NotEmpty(t, cmd.Long)
}

func TestNewSprintRebuildCommand(t *testing.T) {
	cfg := &config.Config{}
	app := NewApp(cfg)
	cmd := newSprintRebuildCommand(app)

	assert.Equal(t, "rebuild", cmd.Use)
	assert.NotEmpty(t, cmd.Short)
	assert.NotEmpty(t, cmd.Long)

	// Check dry-run flag exists
	flag := cmd.Flags().Lookup("dry-run")
	require.NotNil(t, flag)
	assert.Equal(t, "bool", flag.Value.Type())
}

func TestSprintRebuildCommand_Integration(t *testing.T) {
	tmpDir := t.TempDir()

	// Create stories directory
	storiesDir := filepath.Join(tmpDir, "_bmad-output", "implementation-artifacts", "stories")
	err := os.MkdirAll(storiesDir, 0755)
	require.NoError(t, err)

	// Create sprint-status.yaml
	statusDir := filepath.Join(tmpDir, "_bmad-output", "implementation-artifacts")
	sprintPath := filepath.Join(statusDir, "sprint-status.yaml")
	err = os.WriteFile(sprintPath, []byte("development_status:\n"), 0644)
	require.NoError(t, err)

	// Create test story files
	storyFiles := map[string]string{
		"3-5-implement-accessibility-features.md": "# Story 3.5: Implement Accessibility Features\n\nStatus: review\n",
		"3-6-add-user-profile.md":                 "# Story 3.6: Add User Profile\n\nStatus: in-progress\n",
		"4-1-refactor-database.md":                "# Story 4.1: Refactor Database\n\nStatus: done\n",
	}

	for filename, content := range storyFiles {
		err := os.WriteFile(filepath.Join(storiesDir, filename), []byte(content), 0644)
		require.NoError(t, err)
	}

	// Change to temp dir for command execution
	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(origDir)

	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	// Create app and run command
	cfg := &config.Config{}
	app := NewApp(cfg)
	cmd := newSprintRebuildCommand(app)

	// Test normal execution
	err = cmd.Execute()
	require.NoError(t, err)

	// Verify sprint-status.yaml was updated
	data, err := os.ReadFile(sprintPath)
	require.NoError(t, err)

	content := string(data)
	assert.Contains(t, content, "3-5-implement-accessibility-features:")
	assert.Contains(t, content, "3-6-add-user-profile:")
	assert.Contains(t, content, "4-1-refactor-database:")
	assert.Contains(t, content, "review")
	assert.Contains(t, content, "in-progress")
	assert.Contains(t, content, "done")
}

func TestSprintRebuildCommand_DryRun(t *testing.T) {
	tmpDir := t.TempDir()

	// Create stories directory
	storiesDir := filepath.Join(tmpDir, "_bmad-output", "implementation-artifacts", "stories")
	err := os.MkdirAll(storiesDir, 0755)
	require.NoError(t, err)

	// Create a test story file
	err = os.WriteFile(
		filepath.Join(storiesDir, "1-1-test-story.md"),
		[]byte("# Story 1.1: Test Story\n\nStatus: backlog\n"),
		0644,
	)
	require.NoError(t, err)

	// Create empty sprint-status.yaml
	statusDir := filepath.Join(tmpDir, "_bmad-output", "implementation-artifacts")
	sprintPath := filepath.Join(statusDir, "sprint-status.yaml")
	err = os.WriteFile(sprintPath, []byte("development_status:\n"), 0644)
	require.NoError(t, err)

	// Change to temp dir
	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(origDir)

	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	// Create app and run command with dry-run
	cfg := &config.Config{}
	app := NewApp(cfg)
	cmd := newSprintRebuildCommand(app)

	// Set dry-run flag
	err = cmd.Flags().Set("dry-run", "true")
	require.NoError(t, err)

	err = cmd.Execute()
	require.NoError(t, err)

	// Verify sprint-status.yaml was NOT modified (still empty)
	data, err := os.ReadFile(sprintPath)
	require.NoError(t, err)
	assert.Equal(t, "development_status:\n", string(data))
}

func TestSprintRebuildCommand_SkipsInvalidFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Create stories directory
	storiesDir := filepath.Join(tmpDir, "_bmad-output", "implementation-artifacts", "stories")
	err := os.MkdirAll(storiesDir, 0755)
	require.NoError(t, err)

	// Create valid story file
	err = os.WriteFile(
		filepath.Join(storiesDir, "1-1-valid-story.md"),
		[]byte("# Story 1.1: Valid Story\n\nStatus: backlog\n"),
		0644,
	)
	require.NoError(t, err)

	// Create story file without status line
	err = os.WriteFile(
		filepath.Join(storiesDir, "1-2-no-status.md"),
		[]byte("# Story 1.2: No Status\n\nJust some content.\n"),
		0644,
	)
	require.NoError(t, err)

	// Create story file with invalid status
	err = os.WriteFile(
		filepath.Join(storiesDir, "1-3-invalid-status.md"),
		[]byte("# Story 1.3: Invalid Status\n\nStatus: not-a-real-status\n"),
		0644,
	)
	require.NoError(t, err)

	// Create non-md file (should be ignored)
	err = os.WriteFile(
		filepath.Join(storiesDir, "README.txt"),
		[]byte("Status: backlog\n"),
		0644,
	)
	require.NoError(t, err)

	// Create sprint-status.yaml
	statusDir := filepath.Join(tmpDir, "_bmad-output", "implementation-artifacts")
	sprintPath := filepath.Join(statusDir, "sprint-status.yaml")
	err = os.MkdirAll(statusDir, 0755)
	require.NoError(t, err)
	err = os.WriteFile(sprintPath, []byte("development_status:\n"), 0644)
	require.NoError(t, err)

	// Change to temp dir
	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(origDir)

	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	// Create app and run command
	cfg := &config.Config{}
	app := NewApp(cfg)
	cmd := newSprintRebuildCommand(app)

	err = cmd.Execute()
	require.NoError(t, err)

	// Verify sprint-status.yaml only contains valid story
	data, err := os.ReadFile(sprintPath)
	require.NoError(t, err)

	content := string(data)
	assert.Contains(t, content, "1-1-valid-story:")
	assert.NotContains(t, content, "1-2-no-status")
	assert.NotContains(t, content, "1-3-invalid-status")
	assert.NotContains(t, content, "README")
}

func TestReadStoryStatusLine(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name        string
		content     string
		wantStatus  string
		wantErr     bool
	}{
		{
			name:        "valid status",
			content:     "# Story\n\nStatus: review\n",
			wantStatus:  "review",
			wantErr:     false,
		},
		{
			name:        "lowercase status",
			content:     "# Story\n\nstatus: done\n",
			wantStatus:  "done",
			wantErr:     false,
		},
		{
			name:        "no status line",
			content:     "# Story\n\nSome content.\n",
			wantStatus:  "",
			wantErr:     true,
		},
		{
			name:        "empty file",
			content:     "",
			wantStatus:  "",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(tmpDir, tt.name+".md")
			err := os.WriteFile(path, []byte(tt.content), 0644)
			require.NoError(t, err)

			got, err := readStoryStatusLine(path)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantStatus, got)
			}
		})
	}
}

func TestSprintCommand_AddedToRoot(t *testing.T) {
	cfg := &config.Config{}
	app := NewApp(cfg)
	rootCmd := NewRootCommand(app)

	// Find sprint command
	var sprintCmd *cobra.Command
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "sprint" {
			sprintCmd = cmd
			break
		}
	}

	require.NotNil(t, sprintCmd, "sprint command should be added to root")
	assert.Equal(t, "sprint", sprintCmd.Use)

	// Check rebuild subcommand exists
	var rebuildCmd *cobra.Command
	for _, cmd := range sprintCmd.Commands() {
		if cmd.Name() == "rebuild" {
			rebuildCmd = cmd
			break
		}
	}

	require.NotNil(t, rebuildCmd, "rebuild subcommand should exist")
}
