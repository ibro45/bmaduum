// Package cli provides the command-line interface for bmad-automate.
package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"bmad-automate/internal/claude"
	"bmad-automate/internal/config"
	"bmad-automate/internal/output"
	"bmad-automate/internal/status"
	"bmad-automate/internal/workflow"
)

// WorkflowRunner is the interface for running workflows.
type WorkflowRunner interface {
	RunSingle(ctx context.Context, workflowName, storyKey string) int
	RunRaw(ctx context.Context, prompt string) int
}

// StatusReader is the interface for reading story status.
type StatusReader interface {
	GetStoryStatus(storyKey string) (status.Status, error)
	GetEpicStories(epicID string) ([]string, error)
}

// StatusWriter is the interface for updating story status.
type StatusWriter interface {
	UpdateStatus(storyKey string, newStatus status.Status) error
}

// App holds the application dependencies.
type App struct {
	Config       *config.Config
	Executor     claude.Executor
	Printer      output.Printer
	Runner       WorkflowRunner
	Queue        *workflow.QueueRunner
	StatusReader StatusReader
	StatusWriter StatusWriter
}

// NewApp creates a new application with all dependencies wired up.
func NewApp(cfg *config.Config) *App {
	printer := output.NewPrinter()

	executor := claude.NewExecutor(claude.ExecutorConfig{
		BinaryPath:   cfg.Claude.BinaryPath,
		OutputFormat: cfg.Claude.OutputFormat,
		StderrHandler: func(line string) {
			// Print stderr to stderr
			os.Stderr.WriteString("[stderr] " + line + "\n")
		},
	})

	runner := workflow.NewRunner(executor, printer, cfg)
	queue := workflow.NewQueueRunner(runner)
	statusReader := status.NewReader("")
	statusWriter := status.NewWriter("")

	return &App{
		Config:       cfg,
		Executor:     executor,
		Printer:      printer,
		Runner:       runner,
		Queue:        queue,
		StatusReader: statusReader,
		StatusWriter: statusWriter,
	}
}

// NewRootCommand creates the root cobra command.
func NewRootCommand(app *App) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "bmad-automate",
		Short: "BMAD Automation CLI",
		Long: `BMAD Automation CLI - Automate development workflows with Claude.

This tool orchestrates Claude to run development workflows including
story creation, development, code review, and git operations.`,
	}

	// Add subcommands
	rootCmd.AddCommand(
		newCreateStoryCommand(app),
		newDevStoryCommand(app),
		newCodeReviewCommand(app),
		newGitCommitCommand(app),
		newRunCommand(app),
		newQueueCommand(app),
		newEpicCommand(app),
		newRawCommand(app),
	)

	return rootCmd
}

// ExecuteResult holds the result of running the CLI.
type ExecuteResult struct {
	ExitCode int
	Err      error
}

// RunWithConfig creates the app and executes the root command with the given config.
// This is the testable core of Execute().
func RunWithConfig(cfg *config.Config) ExecuteResult {
	app := NewApp(cfg)
	rootCmd := NewRootCommand(app)

	if err := rootCmd.Execute(); err != nil {
		// Check if it's an ExitError from a command
		if code, ok := IsExitError(err); ok {
			return ExecuteResult{ExitCode: code, Err: err}
		}
		// Other errors (e.g., unknown command) - exit code 1
		return ExecuteResult{ExitCode: 1, Err: err}
	}
	return ExecuteResult{ExitCode: 0, Err: nil}
}

// Run loads config and runs the CLI, returning the result.
// This is the fully testable version of Execute().
func Run() ExecuteResult {
	cfg, err := config.NewLoader().Load()
	if err != nil {
		return ExecuteResult{
			ExitCode: 1,
			Err:      fmt.Errorf("error loading config: %w", err),
		}
	}
	return RunWithConfig(cfg)
}

// Execute runs the CLI application.
// This is the entry point called by main() and handles os.Exit().
func Execute() {
	result := Run()
	if result.ExitCode != 0 {
		os.Exit(result.ExitCode)
	}
}
