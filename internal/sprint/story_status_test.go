package sprint

import (
	"os"
	"path/filepath"
	"testing"

	"bmad-automate/internal/status"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestDir creates a temporary directory with story files and sprint-status.yaml
func setupTestDir(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()

	// Create stories directory
	storiesDir := filepath.Join(tmpDir, "_bmad-output", "implementation-artifacts", "stories")
	err := os.MkdirAll(storiesDir, 0755)
	require.NoError(t, err)

	// Create sprint-status.yaml
	statusDir := filepath.Join(tmpDir, "_bmad-output", "implementation-artifacts")
	sprintStatus := `development_status:
  3-5-implement-accessibility-features: backlog
  3-6-add-user-profile: ready-for-dev
  3-7-fix-login-bug: in-progress
  4-1-refactor-database: review
  4-2-update-docs: done
  epic-3-retrospective: backlog
`
	err = os.WriteFile(filepath.Join(statusDir, "sprint-status.yaml"), []byte(sprintStatus), 0644)
	require.NoError(t, err)

	return tmpDir
}

// createStoryFile creates a story file with the given content
func createStoryFile(t *testing.T, storyDir, filename, content string) {
	t.Helper()
	path := filepath.Join(storyDir, filename)
	err := os.WriteFile(path, []byte(content), 0644)
	require.NoError(t, err)
}

func TestNewStoryStatusManager(t *testing.T) {
	t.Run("uses default paths when empty", func(t *testing.T) {
		m := NewStoryStatusManager("", "")

		assert.Equal(t, DefaultStoryDir, m.storyDir)
		assert.Equal(t, status.DefaultStatusPath, m.sprintPath)
	})

	t.Run("uses provided paths", func(t *testing.T) {
		m := NewStoryStatusManager("/custom/stories", "/custom/sprint.yaml")

		assert.Equal(t, "/custom/stories", m.storyDir)
		assert.Equal(t, "/custom/sprint.yaml", m.sprintPath)
	})
}

func TestStoryStatusManager_GetStoryStatus(t *testing.T) {
	tmpDir := setupTestDir(t)
	storiesDir := filepath.Join(tmpDir, "_bmad-output", "implementation-artifacts", "stories")

	// Create test story files
	createStoryFile(t, storiesDir, "3-5-implement-accessibility-features.md", `# Story 3.5: Implement Accessibility Features

Status: review

## Description
This story implements accessibility features.
`)

	createStoryFile(t, storiesDir, "3-6-add-user-profile.md", `# Story 3.6: Add User Profile

Status: in-progress
`)

	createStoryFile(t, storiesDir, "3-7-fix-login-bug.md", `# Story 3.7: Fix Login Bug
Status: done
`)

	createStoryFile(t, storiesDir, "epic-3-retrospective.md", `# Epic 3 Retrospective

Status: backlog
`)

	manager := NewStoryStatusManager(storiesDir, filepath.Join(tmpDir, status.DefaultStatusPath))

	t.Run("reads status from story file", func(t *testing.T) {
		st, err := manager.GetStoryStatus("3-5-implement-accessibility-features")
		require.NoError(t, err)
		assert.Equal(t, status.StatusReview, st)
	})

	t.Run("reads different statuses correctly", func(t *testing.T) {
		st, err := manager.GetStoryStatus("3-6-add-user-profile")
		require.NoError(t, err)
		assert.Equal(t, status.StatusInProgress, st)
	})

	t.Run("finds status without blank line after title", func(t *testing.T) {
		st, err := manager.GetStoryStatus("3-7-fix-login-bug")
		require.NoError(t, err)
		assert.Equal(t, status.StatusDone, st)
	})

	t.Run("reads epic retrospective story", func(t *testing.T) {
		st, err := manager.GetStoryStatus("epic-3-retrospective")
		require.NoError(t, err)
		assert.Equal(t, status.StatusBacklog, st)
	})

	t.Run("returns error for non-existent story", func(t *testing.T) {
		_, err := manager.GetStoryStatus("9-9-non-existent")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get status")
	})
}

func TestStoryStatusManager_GetStoryStatus_InvalidFiles(t *testing.T) {
	tmpDir := setupTestDir(t)
	storiesDir := filepath.Join(tmpDir, "_bmad-output", "implementation-artifacts", "stories")

	manager := NewStoryStatusManager(storiesDir, filepath.Join(tmpDir, status.DefaultStatusPath))

	t.Run("returns error when status line is missing", func(t *testing.T) {
		createStoryFile(t, storiesDir, "no-status.md", `# Story Without Status

This file has no status line.
`)

		_, err := manager.GetStoryStatus("no-status")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "status line not found")
	})

	t.Run("returns error for empty file", func(t *testing.T) {
		createStoryFile(t, storiesDir, "empty.md", ``)

		_, err := manager.GetStoryStatus("empty")
		assert.Error(t, err)
	})
}

func TestStoryStatusManager_UpdateStatus(t *testing.T) {
	tmpDir := setupTestDir(t)
	storiesDir := filepath.Join(tmpDir, "_bmad-output", "implementation-artifacts", "stories")

	// Create test story file
	createStoryFile(t, storiesDir, "3-5-implement-accessibility-features.md", `# Story 3.5: Implement Accessibility Features

Status: backlog

## Description
This story implements accessibility features.
`)

	sprintPath := filepath.Join(tmpDir, status.DefaultStatusPath)
	manager := NewStoryStatusManager(storiesDir, sprintPath)

	t.Run("updates story file status", func(t *testing.T) {
		err := manager.UpdateStatus("3-5-implement-accessibility-features", status.StatusReadyForDev)
		require.NoError(t, err)

		// Verify story file was updated
		st, err := manager.GetStoryStatus("3-5-implement-accessibility-features")
		require.NoError(t, err)
		assert.Equal(t, status.StatusReadyForDev, st)
	})

	t.Run("returns error for invalid status", func(t *testing.T) {
		err := manager.UpdateStatus("3-5-implement-accessibility-features", "invalid-status")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid status")
	})

	t.Run("returns error for non-existent story", func(t *testing.T) {
		err := manager.UpdateStatus("9-9-non-existent", status.StatusDone)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to update story file")
	})
}

func TestStoryStatusManager_GetEpicStories(t *testing.T) {
	tmpDir := setupTestDir(t)
	storiesDir := filepath.Join(tmpDir, "_bmad-output", "implementation-artifacts", "stories")

	// Create test story files for epic 3
	createStoryFile(t, storiesDir, "3-1-first-story.md", `# Story 3.1: First Story

Status: backlog
`)
	createStoryFile(t, storiesDir, "3-5-implement-accessibility-features.md", `# Story 3.5: Implement Accessibility Features

Status: review
`)
	createStoryFile(t, storiesDir, "3-10-tenth-story.md", `# Story 3.10: Tenth Story

Status: in-progress
`)
	createStoryFile(t, storiesDir, "3-2-second-story.md", `# Story 3.2: Second Story

Status: done
`)

	// Create test story files for epic 4
	createStoryFile(t, storiesDir, "4-1-refactor-database.md", `# Story 4.1: Refactor Database

Status: review
`)

	// Create an epic retrospective (no numeric story number)
	createStoryFile(t, storiesDir, "epic-3-retrospective.md", `# Epic 3 Retrospective

Status: backlog
`)

	manager := NewStoryStatusManager(storiesDir, filepath.Join(tmpDir, status.DefaultStatusPath))

	t.Run("returns stories sorted by number", func(t *testing.T) {
		stories, err := manager.GetEpicStories("3")
		require.NoError(t, err)

		// Should be sorted: 3-1, 3-2, 3-5, 3-10 (numeric sort, not string sort)
		assert.Equal(t, []string{
			"3-1-first-story",
			"3-2-second-story",
			"3-5-implement-accessibility-features",
			"3-10-tenth-story",
		}, stories)
	})

	t.Run("returns stories for different epic", func(t *testing.T) {
		stories, err := manager.GetEpicStories("4")
		require.NoError(t, err)

		assert.Equal(t, []string{"4-1-refactor-database"}, stories)
	})

	t.Run("excludes non-numeric story keys", func(t *testing.T) {
		stories, err := manager.GetEpicStories("3")
		require.NoError(t, err)

		// epic-3-retrospective should not be included
		for _, s := range stories {
			assert.NotEqual(t, "epic-3-retrospective", s)
		}
	})

	t.Run("returns error for epic with no stories", func(t *testing.T) {
		_, err := manager.GetEpicStories("99")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no stories found")
	})
}

func TestStoryStatusManager_GetAllEpics(t *testing.T) {
	tmpDir := setupTestDir(t)
	storiesDir := filepath.Join(tmpDir, "_bmad-output", "implementation-artifacts", "stories")

	// Create test story files for multiple epics
	createStoryFile(t, storiesDir, "10-some-story.md", `# Story 10.1: Some Story

Status: backlog
`)
	createStoryFile(t, storiesDir, "2-another-story.md", `# Story 2.1: Another Story

Status: in-progress
`)
	createStoryFile(t, storiesDir, "1-first-epic-story.md", `# Story 1.1: First Epic Story

Status: done
`)

	manager := NewStoryStatusManager(storiesDir, filepath.Join(tmpDir, status.DefaultStatusPath))

	t.Run("returns epics sorted numerically", func(t *testing.T) {
		epics, err := manager.GetAllEpics()
		require.NoError(t, err)

		// Should be sorted: 1, 2, 10 (numeric sort)
		assert.Equal(t, []string{"1", "2", "10"}, epics)
	})
}

func TestStoryStatusManager_GetAllEpics_Empty(t *testing.T) {
	tmpDir := t.TempDir()
	storiesDir := filepath.Join(tmpDir, "stories")
	err := os.MkdirAll(storiesDir, 0755)
	require.NoError(t, err)

	manager := NewStoryStatusManager(storiesDir, filepath.Join(tmpDir, "sprint-status.yaml"))

	t.Run("returns empty slice when no stories", func(t *testing.T) {
		epics, err := manager.GetAllEpics()
		require.NoError(t, err)
		assert.Empty(t, epics)
	})
}

func TestParseStoryNumber(t *testing.T) {
	tests := []struct {
		input    string
		want     int
		wantErr  bool
	}{
		{"1", 1, false},
		{"10", 10, false},
		{"99", 99, false},
		{"abc", 0, true},
		{"", 0, true},
		{"1a", 1, false}, // parses as 1, ignores rest
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := parseStoryNumber(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestFindStoryFile(t *testing.T) {
	manager := &StoryStatusManager{storyDir: "/some/path/stories"}

	tests := []struct {
		storyKey string
		want     string
	}{
		{"3-5-implement-accessibility-features", "/some/path/stories/3-5-implement-accessibility-features.md"},
		{"epic-3-retrospective", "/some/path/stories/epic-3-retrospective.md"},
		{"simple-story", "/some/path/stories/simple-story.md"},
	}

	for _, tt := range tests {
		t.Run(tt.storyKey, func(t *testing.T) {
			got := manager.findStoryFile(tt.storyKey)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestReadStoryStatus_CaseInsensitive(t *testing.T) {
	tmpDir := t.TempDir()
	manager := &StoryStatusManager{storyDir: tmpDir}

	tests := []struct {
		name     string
		content  string
		expected status.Status
	}{
		{
			name: "lowercase status",
			content: `# Story

status: review
`,
			expected: status.StatusReview,
		},
		{
			name: "uppercase STATUS",
			content: `# Story

STATUS: done
`,
			expected: status.StatusDone,
		},
		{
			name: "mixed case Status",
			content: `# Story

Status: backlog
`,
			expected: status.StatusBacklog,
		},
		{
			name: "indented status",
			content: `# Story

  Status: in-progress
`,
			expected: status.StatusInProgress,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(tmpDir, tt.name+".md")
			err := os.WriteFile(path, []byte(tt.content), 0644)
			require.NoError(t, err)

			got, err := manager.readStoryStatus(path)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, got)
		})
	}
}
