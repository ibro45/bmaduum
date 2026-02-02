package cli

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"bmaduum/internal/claude"
	"bmaduum/internal/config"
	"bmaduum/internal/output"
	"bmaduum/internal/status"
	"bmaduum/internal/workflow"
)

func setupTestApp() *App {
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
	statusReader := status.NewReader("")

	return &App{
		Config:       cfg,
		Executor:     mockExecutor,
		Printer:      printer,
		Runner:       runner,
		StatusReader: statusReader,
	}
}

func TestNewApp(t *testing.T) {
	cfg := config.DefaultConfig()
	app := NewApp(cfg)

	assert.NotNil(t, app)
	assert.NotNil(t, app.Config)
	assert.NotNil(t, app.Executor)
	assert.NotNil(t, app.Printer)
	assert.NotNil(t, app.Runner)
	assert.NotNil(t, app.StatusReader)
	assert.Equal(t, cfg, app.Config)
}

func TestNewRootCommand(t *testing.T) {
	app := setupTestApp()
	rootCmd := NewRootCommand(app)

	assert.NotNil(t, rootCmd)
	assert.Equal(t, "bmaduum", rootCmd.Use)
	assert.Contains(t, rootCmd.Short, "BMAD")
}

func TestNewRootCommand_HasAllSubcommands(t *testing.T) {
	app := setupTestApp()
	rootCmd := NewRootCommand(app)

	expectedCommands := []string{
		"story",
		"epic",
		"workflow",
		"raw",
	}

	commands := rootCmd.Commands()
	commandNames := make([]string, len(commands))
	for i, cmd := range commands {
		commandNames[i] = cmd.Name()
	}

	for _, expected := range expectedCommands {
		assert.Contains(t, commandNames, expected, "missing subcommand: %s", expected)
	}

	// Verify workflow has the expected subcommands
	workflowCmd := findCommand(rootCmd, "workflow")
	require.NotNil(t, workflowCmd, "workflow command not found")

	expectedSubcommands := []string{
		"create-story",
		"dev-story",
		"code-review",
		"git-commit",
	}

	subcommands := workflowCmd.Commands()
	subcommandNames := make([]string, len(subcommands))
	for i, cmd := range subcommands {
		subcommandNames[i] = cmd.Name()
	}

	for _, expected := range expectedSubcommands {
		assert.Contains(t, subcommandNames, expected, "missing workflow subcommand: %s", expected)
	}
}

func TestWorkflowCommand(t *testing.T) {
	app := setupTestApp()
	cmd := newWorkflowCommand(app)

	assert.Equal(t, "workflow <workflow-name> <story-key>", cmd.Use)
	assert.NotEmpty(t, cmd.Short)
	assert.NotEmpty(t, cmd.Long)

	// Test args validation - should require at least 2 args (workflow-name, story-key)
	err := cmd.Args(cmd, []string{})
	assert.Error(t, err)

	err = cmd.Args(cmd, []string{"create-story"})
	assert.Error(t, err)

	err = cmd.Args(cmd, []string{"create-story", "story-1"})
	assert.NoError(t, err)
}

func TestWorkflowSubcommands(t *testing.T) {
	app := setupTestApp()
	workflowCmd := newWorkflowCommand(app)

	subcommands := []struct {
		name string
	}{
		{"create-story"},
		{"dev-story"},
		{"code-review"},
		{"git-commit"},
	}

	for _, tt := range subcommands {
		t.Run(tt.name, func(t *testing.T) {
			cmd := findCommand(workflowCmd, tt.name)
			require.NotNil(t, cmd, "workflow subcommand %s not found", tt.name)

			assert.Equal(t, tt.name+" <story-key>", cmd.Use)
			assert.NotEmpty(t, cmd.Short)

			// Test args validation - should require exactly 1 arg (story-key)
			err := cmd.Args(cmd, []string{})
			assert.Error(t, err)

			err = cmd.Args(cmd, []string{"story-1"})
			assert.NoError(t, err)

			err = cmd.Args(cmd, []string{"story-1", "extra"})
			assert.Error(t, err)
		})
	}
}

func TestStoryCommand(t *testing.T) {
	app := setupTestApp()
	cmd := newStoryCommand(app)

	assert.Equal(t, "story <story-key> [story-key...]", cmd.Use)
	assert.NotEmpty(t, cmd.Short)
	assert.NotEmpty(t, cmd.Long)

	// Test args validation - should require at least 1 arg
	err := cmd.Args(cmd, []string{})
	assert.Error(t, err)

	err = cmd.Args(cmd, []string{"story-1"})
	assert.NoError(t, err)
}

func TestRawCommand(t *testing.T) {
	app := setupTestApp()
	cmd := newRawCommand(app)

	assert.Equal(t, "raw <prompt>", cmd.Use)
	assert.NotEmpty(t, cmd.Short)
	assert.NotEmpty(t, cmd.Long)

	// Test args validation - should require at least 1 arg
	err := cmd.Args(cmd, []string{})
	assert.Error(t, err)

	err = cmd.Args(cmd, []string{"hello"})
	assert.NoError(t, err)

	err = cmd.Args(cmd, []string{"hello", "world", "test"})
	assert.NoError(t, err)
}

func TestRootCommand_Help(t *testing.T) {
	app := setupTestApp()
	rootCmd := NewRootCommand(app)

	// Capture help output
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"--help"})

	err := rootCmd.Execute()
	require.NoError(t, err)

	helpOutput := buf.String()
	assert.Contains(t, helpOutput, "BMAD")
	assert.Contains(t, helpOutput, "Available Commands")
}

func TestSubcommand_Help(t *testing.T) {
	app := setupTestApp()
	rootCmd := NewRootCommand(app)

	tests := []struct {
		name    string
		command string
	}{
			{"workflow create-story help", "workflow"},
		{"workflow dev-story help", "workflow"},
		{"workflow code-review help", "workflow"},
		{"workflow git-commit help", "workflow"},
		{"story help", "story"},
		{"epic help", "epic"},
		{"raw help", "raw"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			rootCmd.SetOut(buf)
			rootCmd.SetArgs([]string{tt.command, "--help"})

			err := rootCmd.Execute()
			require.NoError(t, err)

			helpOutput := buf.String()
			assert.NotEmpty(t, helpOutput)
		})
	}
}

// findCommand finds a subcommand by name
func findCommand(root *cobra.Command, name string) *cobra.Command {
	for _, cmd := range root.Commands() {
		if cmd.Name() == name {
			return cmd
		}
	}
	return nil
}

func TestCommandsHaveRunEFunctions(t *testing.T) {
	app := setupTestApp()
	rootCmd := NewRootCommand(app)

	commands := []string{
		"story",
		"epic",
		"raw",
		"workflow",
	}

	for _, cmdName := range commands {
		t.Run(cmdName, func(t *testing.T) {
			cmd := findCommand(rootCmd, cmdName)
			require.NotNil(t, cmd, "command %s not found", cmdName)
			assert.NotNil(t, cmd.RunE, "command %s should have a RunE function", cmdName)
		})
	}

	// Test workflow subcommands
	workflowCmd := findCommand(rootCmd, "workflow")
	require.NotNil(t, workflowCmd, "workflow command not found")

	workflowSubcommands := []string{
		"create-story",
		"dev-story",
		"code-review",
		"git-commit",
	}

	for _, subcmdName := range workflowSubcommands {
		t.Run("workflow."+subcmdName, func(t *testing.T) {
			cmd := findCommand(workflowCmd, subcmdName)
			require.NotNil(t, cmd, "workflow subcommand %s not found", subcmdName)
			assert.NotNil(t, cmd.RunE, "workflow subcommand %s should have a RunE function", subcmdName)
		})
	}
}

// setupFailingTestApp creates an App with a mock executor that returns exit code 1.
func setupFailingTestApp() *App {
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
	statusReader := status.NewReader("")

	return &App{
		Config:       cfg,
		Executor:     mockExecutor,
		Printer:      printer,
		Runner:       runner,
		StatusReader: statusReader,
	}
}

func TestCommandExecution_Success(t *testing.T) {
	// Note: "story" and "epic" commands excluded - they require sprint-status.yaml and are tested in their respective test files
	tests := []struct {
		name    string
		command string
		args    []string
	}{
		{"workflow create-story", "workflow", []string{"create-story", "TEST-123"}},
		{"workflow dev-story", "workflow", []string{"dev-story", "TEST-123"}},
		{"workflow code-review", "workflow", []string{"code-review", "TEST-123"}},
		{"workflow git-commit", "workflow", []string{"git-commit", "TEST-123"}},
		{"raw", "raw", []string{"hello", "world"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := setupTestApp()
			rootCmd := NewRootCommand(app)

			buf := &bytes.Buffer{}
			rootCmd.SetOut(buf)
			rootCmd.SetErr(buf)
			rootCmd.SetArgs(append([]string{tt.command}, tt.args...))

			err := rootCmd.Execute()
			assert.NoError(t, err)
		})
	}
}

func TestCommandExecution_Failure(t *testing.T) {
	// Note: "story" and "epic" commands excluded - they require sprint-status.yaml and are tested in their respective test files
	tests := []struct {
		name    string
		command string
		args    []string
	}{
		{"workflow create-story", "workflow", []string{"create-story", "TEST-123"}},
		{"workflow dev-story", "workflow", []string{"dev-story", "TEST-123"}},
		{"workflow code-review", "workflow", []string{"code-review", "TEST-123"}},
		{"workflow git-commit", "workflow", []string{"git-commit", "TEST-123"}},
		{"raw", "raw", []string{"hello"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := setupFailingTestApp()
			rootCmd := NewRootCommand(app)

			buf := &bytes.Buffer{}
			rootCmd.SetOut(buf)
			rootCmd.SetErr(buf)
			rootCmd.SetArgs(append([]string{tt.command}, tt.args...))

			err := rootCmd.Execute()
			require.Error(t, err)

			code, ok := IsExitError(err)
			assert.True(t, ok, "error should be an ExitError")
			assert.Equal(t, 1, code)
		})
	}
}

func TestCommandExecution_InvalidArgs(t *testing.T) {
	tests := []struct {
		name    string
		command string
		args    []string
	}{
		{"workflow create-story no args", "workflow", []string{"create-story"}},
		{"workflow dev-story no args", "workflow", []string{"dev-story"}},
		{"workflow code-review no args", "workflow", []string{"code-review"}},
		{"workflow git-commit no args", "workflow", []string{"git-commit"}},
		{"story no args", "story", []string{}},
		{"epic no args", "epic", []string{}},
		{"raw no args", "raw", []string{}},
		{"workflow too few args", "workflow", []string{}},
		{"workflow create-story too many args", "workflow", []string{"create-story", "a", "b"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := setupTestApp()
			rootCmd := NewRootCommand(app)

			buf := &bytes.Buffer{}
			rootCmd.SetOut(buf)
			rootCmd.SetErr(buf)
			rootCmd.SetArgs(append([]string{tt.command}, tt.args...))

			err := rootCmd.Execute()
			require.Error(t, err)
		})
	}
}

func TestRunWithConfig_Success(t *testing.T) {
	cfg := config.DefaultConfig()

	// With no args, rootCmd.Execute() shows help and returns nil
	result := RunWithConfig(cfg)

	assert.Equal(t, 0, result.ExitCode)
	assert.NoError(t, result.Err)
}

func TestExecuteResult(t *testing.T) {
	t.Run("zero exit code", func(t *testing.T) {
		result := ExecuteResult{ExitCode: 0, Err: nil}
		assert.Equal(t, 0, result.ExitCode)
		assert.NoError(t, result.Err)
	})

	t.Run("non-zero exit code with error", func(t *testing.T) {
		result := ExecuteResult{ExitCode: 1, Err: NewExitError(1)}
		assert.Equal(t, 1, result.ExitCode)
		assert.Error(t, result.Err)
	})
}
