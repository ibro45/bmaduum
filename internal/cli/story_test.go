package cli

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"bmaduum/internal/config"
	"bmaduum/internal/output"
	"bmaduum/internal/status"
)

// TestStoryCommand_SingleStory tests that story command works for a single story
func TestStoryCommand_SingleStory(t *testing.T) {
	tests := []struct {
		name              string
		storyKey          string
		statusYAML        string
		expectedWorkflows []string
		expectedStatuses  []StatusUpdate
		expectError       bool
		failOnWorkflow    string
	}{
		{
			name:     "single backlog story runs full lifecycle",
			storyKey: "STORY-1",
			statusYAML: `development_status:
  STORY-1: backlog`,
			expectedWorkflows: []string{
				"create-story", "dev-story", "code-review", "git-commit",
			},
			expectedStatuses: []StatusUpdate{
				{StoryKey: "STORY-1", NewStatus: status.StatusReadyForDev},
				{StoryKey: "STORY-1", NewStatus: status.StatusReview},
				{StoryKey: "STORY-1", NewStatus: status.StatusDone},
				{StoryKey: "STORY-1", NewStatus: status.StatusDone},
			},
			expectError: false,
		},
		{
			name:     "story at review runs only remaining workflows",
			storyKey: "STORY-1",
			statusYAML: `development_status:
  STORY-1: review`,
			expectedWorkflows: []string{
				"code-review", "git-commit",
			},
			expectedStatuses: []StatusUpdate{
				{StoryKey: "STORY-1", NewStatus: status.StatusDone},
				{StoryKey: "STORY-1", NewStatus: status.StatusDone},
			},
			expectError: false,
		},
		{
			name:     "done story is skipped",
			storyKey: "STORY-1",
			statusYAML: `development_status:
  STORY-1: done`,
			expectedWorkflows: nil,
			expectedStatuses:  nil,
			expectError:       false,
		},
		{
			name:     "story failure stops execution",
			storyKey: "STORY-1",
			statusYAML: `development_status:
  STORY-1: backlog`,
			failOnWorkflow: "dev-story",
			expectedWorkflows: []string{
				"create-story", "dev-story",
			},
			expectedStatuses: []StatusUpdate{
				{StoryKey: "STORY-1", NewStatus: status.StatusReadyForDev},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			createSprintStatusFile(t, tmpDir, tt.statusYAML)

			mockRunner := &MockWorkflowRunner{
				FailOnWorkflow: tt.failOnWorkflow,
			}
			mockWriter := &MockStatusWriter{}
			statusReader := status.NewReader(tmpDir)
			buf := &bytes.Buffer{}
			printer := output.NewPrinterWithWriter(buf)

			app := &App{
				Config:       config.DefaultConfig(),
				StatusReader: statusReader,
				StatusWriter: mockWriter,
				Runner:       mockRunner,
				Printer:      printer,
			}

			rootCmd := NewRootCommand(app)
			outBuf := &bytes.Buffer{}
			rootCmd.SetOut(outBuf)
			rootCmd.SetErr(outBuf)
			rootCmd.SetArgs([]string{"story", tt.storyKey})

			err := rootCmd.Execute()

			if tt.expectError {
				require.Error(t, err)
				code, ok := IsExitError(err)
				assert.True(t, ok, "error should be an ExitError")
				assert.Equal(t, 1, code)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedWorkflows, mockRunner.ExecutedWorkflows,
				"workflows should be executed in lifecycle order")

			if tt.expectedStatuses != nil {
				require.Len(t, mockWriter.Updates, len(tt.expectedStatuses))
				for i, expected := range tt.expectedStatuses {
					assert.Equal(t, expected.StoryKey, mockWriter.Updates[i].StoryKey)
					assert.Equal(t, expected.NewStatus, mockWriter.Updates[i].NewStatus)
				}
			}
		})
	}
}

// TestStoryCommand_MultipleStories tests that story command processes multiple stories
func TestStoryCommand_MultipleStories(t *testing.T) {
	tests := []struct {
		name              string
		storyKeys         []string
		statusYAML        string
		expectedWorkflows []string
		expectedStatuses  []StatusUpdate
		expectError       bool
	}{
		{
			name:      "multiple stories run in sequence",
			storyKeys: []string{"STORY-1", "STORY-2"},
			statusYAML: `development_status:
  STORY-1: backlog
  STORY-2: backlog`,
			expectedWorkflows: []string{
				"create-story", "dev-story", "code-review", "git-commit",
				"create-story", "dev-story", "code-review", "git-commit",
			},
			expectedStatuses: []StatusUpdate{
				{StoryKey: "STORY-1", NewStatus: status.StatusReadyForDev},
				{StoryKey: "STORY-1", NewStatus: status.StatusReview},
				{StoryKey: "STORY-1", NewStatus: status.StatusDone},
				{StoryKey: "STORY-1", NewStatus: status.StatusDone},
				{StoryKey: "STORY-2", NewStatus: status.StatusReadyForDev},
				{StoryKey: "STORY-2", NewStatus: status.StatusReview},
				{StoryKey: "STORY-2", NewStatus: status.StatusDone},
				{StoryKey: "STORY-2", NewStatus: status.StatusDone},
			},
			expectError: false,
		},
		{
			name:      "mixed statuses run appropriate workflows",
			storyKeys: []string{"STORY-1", "STORY-2", "STORY-3"},
			statusYAML: `development_status:
  STORY-1: backlog
  STORY-2: review
  STORY-3: done`,
			expectedWorkflows: []string{
				"create-story", "dev-story", "code-review", "git-commit",
				"code-review", "git-commit",
				// STORY-3 is done, skipped
			},
			expectedStatuses: []StatusUpdate{
				{StoryKey: "STORY-1", NewStatus: status.StatusReadyForDev},
				{StoryKey: "STORY-1", NewStatus: status.StatusReview},
				{StoryKey: "STORY-1", NewStatus: status.StatusDone},
				{StoryKey: "STORY-1", NewStatus: status.StatusDone},
				{StoryKey: "STORY-2", NewStatus: status.StatusDone},
				{StoryKey: "STORY-2", NewStatus: status.StatusDone},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			createSprintStatusFile(t, tmpDir, tt.statusYAML)

			mockRunner := &MockWorkflowRunner{}
			mockWriter := &MockStatusWriter{}
			statusReader := status.NewReader(tmpDir)
			buf := &bytes.Buffer{}
			printer := output.NewPrinterWithWriter(buf)

			app := &App{
				Config:       config.DefaultConfig(),
				StatusReader: statusReader,
				StatusWriter: mockWriter,
				Runner:       mockRunner,
				Printer:      printer,
			}

			rootCmd := NewRootCommand(app)
			outBuf := &bytes.Buffer{}
			rootCmd.SetOut(outBuf)
			rootCmd.SetErr(outBuf)
			rootCmd.SetArgs(append([]string{"story"}, tt.storyKeys...))

			err := rootCmd.Execute()

			if tt.expectError {
				require.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedWorkflows, mockRunner.ExecutedWorkflows)

			if tt.expectedStatuses != nil {
				require.Len(t, mockWriter.Updates, len(tt.expectedStatuses))
				for i, expected := range tt.expectedStatuses {
					assert.Equal(t, expected.StoryKey, mockWriter.Updates[i].StoryKey)
					assert.Equal(t, expected.NewStatus, mockWriter.Updates[i].NewStatus)
				}
			}
		})
	}
}

// TestStoryCommand_DryRun tests dry-run mode for story command
func TestStoryCommand_DryRun(t *testing.T) {
	tests := []struct {
		name       string
		storyKeys  []string
		statusYAML string
	}{
		{
			name:      "single story dry run",
			storyKeys: []string{"STORY-1"},
			statusYAML: `development_status:
  STORY-1: backlog`,
		},
		{
			name:      "multiple stories dry run",
			storyKeys: []string{"STORY-1", "STORY-2"},
			statusYAML: `development_status:
  STORY-1: backlog
  STORY-2: review`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			createSprintStatusFile(t, tmpDir, tt.statusYAML)

			mockRunner := &MockWorkflowRunner{}
			mockWriter := &MockStatusWriter{}
			statusReader := status.NewReader(tmpDir)
			buf := &bytes.Buffer{}
			printer := output.NewPrinterWithWriter(buf)

			app := &App{
				Config:       config.DefaultConfig(),
				StatusReader: statusReader,
				StatusWriter: mockWriter,
				Runner:       mockRunner,
				Printer:      printer,
			}

			rootCmd := NewRootCommand(app)
			outBuf := &bytes.Buffer{}
			rootCmd.SetOut(outBuf)
			rootCmd.SetErr(outBuf)
			rootCmd.SetArgs(append([]string{"story", "--dry-run"}, tt.storyKeys...))

			err := rootCmd.Execute()
			assert.NoError(t, err)

			// No workflows should have been executed in dry-run
			assert.Empty(t, mockRunner.ExecutedWorkflows)
			assert.Empty(t, mockWriter.Updates)
		})
	}
}
