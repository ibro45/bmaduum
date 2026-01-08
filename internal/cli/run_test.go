package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"bmad-automate/internal/claude"
	"bmad-automate/internal/config"
	"bmad-automate/internal/output"
	"bmad-automate/internal/status"
	"bmad-automate/internal/workflow"
)

func setupRunTestApp(tmpDir string) (*App, *claude.MockExecutor, *bytes.Buffer) {
	cfg := config.DefaultConfig()
	buf := &bytes.Buffer{}
	printer := output.NewPrinterWithWriter(buf)
	mockExecutor := &claude.MockExecutor{
		Events: []claude.Event{
			{Type: claude.EventTypeSystem, SessionStarted: true},
			{Type: claude.EventTypeResult, SessionComplete: true},
		},
		ExitCode: 0,
	}
	runner := workflow.NewRunner(mockExecutor, printer, cfg)
	queue := workflow.NewQueueRunner(runner)
	statusReader := status.NewReader(tmpDir)

	return &App{
		Config:       cfg,
		Executor:     mockExecutor,
		Printer:      printer,
		Runner:       runner,
		Queue:        queue,
		StatusReader: statusReader,
	}, mockExecutor, buf
}

func createSprintStatusFile(t *testing.T, tmpDir string, content string) {
	t.Helper()
	artifactsDir := filepath.Join(tmpDir, "_bmad-output", "implementation-artifacts")
	require.NoError(t, os.MkdirAll(artifactsDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(artifactsDir, "sprint-status.yaml"), []byte(content), 0644))
}

func TestRunCommand_StatusBasedRouting(t *testing.T) {
	tests := []struct {
		name             string
		storyKey         string
		statusYAML       string
		expectedWorkflow string
		expectError      bool
		expectExitCode   int
		expectedOutput   string
	}{
		{
			name:     "backlog status routes to create-story",
			storyKey: "STORY-1",
			statusYAML: `development_status:
  STORY-1: backlog`,
			expectedWorkflow: "create-story",
			expectError:      false,
		},
		{
			name:     "ready-for-dev status routes to dev-story",
			storyKey: "STORY-2",
			statusYAML: `development_status:
  STORY-2: ready-for-dev`,
			expectedWorkflow: "dev-story",
			expectError:      false,
		},
		{
			name:     "in-progress status routes to dev-story",
			storyKey: "STORY-3",
			statusYAML: `development_status:
  STORY-3: in-progress`,
			expectedWorkflow: "dev-story",
			expectError:      false,
		},
		{
			name:     "review status routes to code-review",
			storyKey: "STORY-4",
			statusYAML: `development_status:
  STORY-4: review`,
			expectedWorkflow: "code-review",
			expectError:      false,
		},
		{
			name:     "done status prints completion message",
			storyKey: "STORY-5",
			statusYAML: `development_status:
  STORY-5: done`,
			expectedWorkflow: "",
			expectError:      false,
			expectedOutput:   "", // Output goes to fmt.Printf, not captured
		},
		{
			name:     "story not found returns error",
			storyKey: "STORY-NOT-FOUND",
			statusYAML: `development_status:
  STORY-1: backlog`,
			expectError:    true,
			expectExitCode: 1,
			expectedOutput: "", // Output goes to fmt.Printf, not captured
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			createSprintStatusFile(t, tmpDir, tt.statusYAML)

			app, mockExecutor, _ := setupRunTestApp(tmpDir)
			rootCmd := NewRootCommand(app)

			outBuf := &bytes.Buffer{}
			errBuf := &bytes.Buffer{}
			rootCmd.SetOut(outBuf)
			rootCmd.SetErr(errBuf)
			rootCmd.SetArgs([]string{"run", tt.storyKey})

			err := rootCmd.Execute()

			if tt.expectError {
				require.Error(t, err)
				if tt.expectExitCode > 0 {
					code, ok := IsExitError(err)
					assert.True(t, ok, "error should be an ExitError")
					assert.Equal(t, tt.expectExitCode, code)
				}
			} else {
				assert.NoError(t, err)
			}

			if tt.expectedWorkflow != "" {
				assert.NotEmpty(t, mockExecutor.RecordedPrompts, "prompt should have been executed")
			}

			if tt.expectedOutput != "" {
				assert.Contains(t, outBuf.String()+errBuf.String(), tt.expectedOutput)
			}
		})
	}
}

func TestRunCommand_MissingSprintStatusFile(t *testing.T) {
	tmpDir := t.TempDir()
	// Don't create sprint-status.yaml

	app, _, _ := setupRunTestApp(tmpDir)
	rootCmd := NewRootCommand(app)

	outBuf := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}
	rootCmd.SetOut(outBuf)
	rootCmd.SetErr(errBuf)
	rootCmd.SetArgs([]string{"run", "STORY-1"})

	err := rootCmd.Execute()

	require.Error(t, err)
	code, ok := IsExitError(err)
	assert.True(t, ok, "error should be an ExitError")
	assert.Equal(t, 1, code)
}

func TestRunCommand_WorkflowExecution(t *testing.T) {
	tests := []struct {
		name             string
		storyKey         string
		status           string
		expectedWorkflow string
	}{
		{"backlog executes create-story", "S1", "backlog", "create-story"},
		{"ready-for-dev executes dev-story", "S2", "ready-for-dev", "dev-story"},
		{"in-progress executes dev-story", "S3", "in-progress", "dev-story"},
		{"review executes code-review", "S4", "review", "code-review"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			createSprintStatusFile(t, tmpDir, "development_status:\n  "+tt.storyKey+": "+tt.status)

			app, mockExecutor, _ := setupRunTestApp(tmpDir)
			rootCmd := NewRootCommand(app)

			outBuf := &bytes.Buffer{}
			rootCmd.SetOut(outBuf)
			rootCmd.SetErr(outBuf)
			rootCmd.SetArgs([]string{"run", tt.storyKey})

			err := rootCmd.Execute()
			require.NoError(t, err)

			// The workflow runner should have been called
			assert.NotEmpty(t, mockExecutor.RecordedPrompts)
		})
	}
}

func TestRunCommand_WorkflowFailure(t *testing.T) {
	tmpDir := t.TempDir()
	createSprintStatusFile(t, tmpDir, `development_status:
  STORY-1: backlog`)

	cfg := config.DefaultConfig()
	buf := &bytes.Buffer{}
	printer := output.NewPrinterWithWriter(buf)
	mockExecutor := &claude.MockExecutor{
		Events: []claude.Event{
			{Type: claude.EventTypeSystem, SessionStarted: true},
			{Type: claude.EventTypeResult, SessionComplete: true},
		},
		ExitCode: 1, // Simulate failure
	}
	runner := workflow.NewRunner(mockExecutor, printer, cfg)
	queue := workflow.NewQueueRunner(runner)
	statusReader := status.NewReader(tmpDir)

	app := &App{
		Config:       cfg,
		Executor:     mockExecutor,
		Printer:      printer,
		Runner:       runner,
		Queue:        queue,
		StatusReader: statusReader,
	}

	rootCmd := NewRootCommand(app)
	outBuf := &bytes.Buffer{}
	rootCmd.SetOut(outBuf)
	rootCmd.SetErr(outBuf)
	rootCmd.SetArgs([]string{"run", "STORY-1"})

	err := rootCmd.Execute()
	require.Error(t, err)

	code, ok := IsExitError(err)
	assert.True(t, ok, "error should be an ExitError")
	assert.Equal(t, 1, code)
}
