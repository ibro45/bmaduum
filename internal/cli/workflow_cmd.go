package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

func newWorkflowCommand(app *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "workflow <workflow-name> <story-key>",
		Short: "Run individual BMAD workflow steps (advanced)",
		Long: `Run individual BMAD workflow steps directly.

These are the same workflow commands that are used in BMAD-METHOD and are
automatically executed by the 'story' and 'epic' commands. Use these commands
to run individual workflow steps outside of the full lifecycle automation.

Available workflows:
  - create-story: Create a story definition from backlog
  - dev-story: Implement a story (ready-for-dev or in-progress)
  - code-review: Review code changes (review status)
  - git-commit: Commit and push changes after review

Most users should use 'story' or 'epic' commands instead, which automatically
run the appropriate workflows based on story status.

Use these individual workflow commands when:
  - A workflow fails and you want to retry just that step
  - You need to run a step out of the normal sequence
  - You're testing or developing workflow prompts

Example:
  bmaduum workflow create-story PROJ-123
  bmaduum workflow dev-story PROJ-123
  bmaduum workflow code-review PROJ-123
  bmaduum workflow git-commit PROJ-123`,
		Args: cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			workflowName := args[0]
			storyKey := args[1]

			// Get the subcommand that matches the workflow name
			for _, subcmd := range cmd.Commands() {
				if subcmd.Name() == workflowName {
					// Call the subcommand directly with remaining args
					_ = storyKey // avoid unused variable warning
					return subcmd.RunE(subcmd, args[1:])
				}
			}

			return fmt.Errorf("unknown workflow: %s (valid workflows: create-story, dev-story, code-review, git-commit)", workflowName)
		},
	}

	// Add subcommands for each workflow
	cmd.AddCommand(
		newCreateStoryWorkflowCommand(app),
		newDevStoryWorkflowCommand(app),
		newCodeReviewWorkflowCommand(app),
		newGitCommitWorkflowCommand(app),
	)

	return cmd
}

// newCreateStoryWorkflowCommand creates the create-story workflow subcommand
func newCreateStoryWorkflowCommand(app *App) *cobra.Command {
	var autoRetry bool

	cmd := &cobra.Command{
		Use:   "create-story <story-key>",
		Short: "Create a story definition from backlog",
		Long: `Create a story definition from backlog status.

This workflow generates a detailed story definition including acceptance criteria,
technical requirements, and implementation plan.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			storyKey := args[0]

			return executeWorkflowWithRetry(ctx, cmd, app, "create-story", storyKey, autoRetry)
		},
	}

	cmd.Flags().BoolVar(&autoRetry, "auto-retry", false, "Automatically retry on rate limit errors")
	return cmd
}

// newDevStoryWorkflowCommand creates the dev-story workflow subcommand
func newDevStoryWorkflowCommand(app *App) *cobra.Command {
	var autoRetry bool

	cmd := &cobra.Command{
		Use:   "dev-story <story-key>",
		Short: "Implement a story",
		Long: `Implement a story through development.

This workflow handles stories with status ready-for-dev or in-progress,
running the actual implementation work.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			storyKey := args[0]

			return executeWorkflowWithRetry(ctx, cmd, app, "dev-story", storyKey, autoRetry)
		},
	}

	cmd.Flags().BoolVar(&autoRetry, "auto-retry", false, "Automatically retry on rate limit errors")
	return cmd
}

// newCodeReviewWorkflowCommand creates the code-review workflow subcommand
func newCodeReviewWorkflowCommand(app *App) *cobra.Command {
	var autoRetry bool

	cmd := &cobra.Command{
		Use:   "code-review <story-key>",
		Short: "Review code changes",
		Long: `Review code changes for a story.

This workflow reviews the implementation and generates feedback.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			storyKey := args[0]

			return executeWorkflowWithRetry(ctx, cmd, app, "code-review", storyKey, autoRetry)
		},
	}

	cmd.Flags().BoolVar(&autoRetry, "auto-retry", false, "Automatically retry on rate limit errors")
	return cmd
}

// newGitCommitWorkflowCommand creates the git-commit workflow subcommand
func newGitCommitWorkflowCommand(app *App) *cobra.Command {
	var autoRetry bool

	cmd := &cobra.Command{
		Use:   "git-commit <story-key>",
		Short: "Commit and push changes",
		Long: `Commit and push changes for a story.

This workflow commits the changes and pushes them to the remote repository.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			storyKey := args[0]

			return executeWorkflowWithRetry(ctx, cmd, app, "git-commit", storyKey, autoRetry)
		},
	}

	cmd.Flags().BoolVar(&autoRetry, "auto-retry", false, "Automatically retry on rate limit errors")
	return cmd
}

// executeWorkflowWithRetry executes a single workflow with optional retry logic
func executeWorkflowWithRetry(ctx context.Context, cmd *cobra.Command, app *App, workflowName, storyKey string, autoRetry bool) error {
	exitCode := app.Runner.RunSingle(ctx, workflowName, storyKey)
	if exitCode != 0 {
		cmd.SilenceUsage = true
		return NewExitError(exitCode)
	}
	return nil
}
