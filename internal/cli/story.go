package cli

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"bmaduum/internal/lifecycle"
	"bmaduum/internal/router"
	"bmaduum/internal/tui"
)

func newStoryCommand(app *App) *cobra.Command {
	var dryRun bool
	var autoRetry bool
	var useTUI bool

	cmd := &cobra.Command{
		Use:   "story <story-key> [story-key...]",
		Short: "Run the full story lifecycle to completion",
		Long: `Run the complete lifecycle for one or more stories from their current status to done.

Each story is run to completion before moving to the next.

For each story, executes all remaining workflows based on its current status:
  - backlog       → create-story → dev-story → code-review → git-commit → done
  - ready-for-dev → dev-story → code-review → git-commit → done
  - in-progress   → dev-story → code-review → git-commit → done
  - review        → code-review → git-commit → done
  - done          → skipped (story already complete)

The command stops on the first failure. Done stories are skipped and do not cause failure.
Status is updated in sprint-status.yaml after each successful workflow.

Use --dry-run to preview workflows without executing them.
Use --auto-retry to automatically retry on rate limit errors.
Use --tui to enable the interactive TUI mode.

Examples:
  bmaduum story 6-1
  bmaduum story 6-1 6-2 6-3
  bmaduum story --tui 6-1`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			storyKeys := args

			// Create lifecycle executor with app dependencies
			executor := lifecycle.NewExecutor(app.Runner, app.StatusReader, app.StatusWriter)

			// Handle dry-run mode
			if dryRun {
				return runStoryDryRun(cmd, app, executor, storyKeys)
			}

			// Handle TUI mode (only for single story)
			if useTUI {
				if len(storyKeys) > 1 {
					return fmt.Errorf("TUI mode only supports single story execution")
				}
				return executeStoryTUI(ctx, app, executor, storyKeys[0])
			}

			// Execute full lifecycle for each story in order
			for i, storyKey := range storyKeys {
				// Show story progress for multiple stories
				if len(storyKeys) > 1 {
					fmt.Printf("─── Story %d of %d: %s\n", i+1, len(storyKeys), storyKey)
				}

				err := executeWithRetry(ctx, executor, storyKey, autoRetry, 10, func(stepIndex, totalSteps int, workflow string) {
					app.Printer.StepStart(stepIndex, totalSteps, workflow)
				})
				if err != nil {
					cmd.SilenceUsage = true
					if errors.Is(err, router.ErrStoryComplete) {
						fmt.Printf("Story %s is already complete, skipping\n", storyKey)
						continue
					}
					fmt.Printf("Error running lifecycle for story %s: %v\n", storyKey, err)
					return NewExitError(1)
				}

				// Show completion message
				if len(storyKeys) > 1 {
					fmt.Printf("Story %s completed successfully\n\n", storyKey)
				}
			}

			if len(storyKeys) > 1 {
				fmt.Printf("All %d stories processed\n", len(storyKeys))
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview workflows without executing them")
	cmd.Flags().BoolVar(&autoRetry, "auto-retry", false, "Automatically retry on rate limit errors")
	cmd.Flags().BoolVar(&useTUI, "tui", false, "Enable interactive TUI mode (single story only)")

	return cmd
}

func runStoryDryRun(cmd *cobra.Command, app *App, executor *lifecycle.Executor, storyKeys []string) error {
	// Single story dry-run - simpler output
	if len(storyKeys) == 1 {
		storyKey := storyKeys[0]
		steps, err := executor.GetSteps(storyKey)
		if err != nil {
			cmd.SilenceUsage = true
			if errors.Is(err, router.ErrStoryComplete) {
				fmt.Printf("Story is already complete, no workflows to run\n")
				return nil
			}
			fmt.Printf("Error: %v\n", err)
			return NewExitError(1)
		}

		fmt.Printf("Dry run for story %s:\n", storyKey)
		for i, step := range steps {
			modelInfo := ""
			model := app.Config.GetModel(step.Workflow)
			if model != "" {
				modelInfo = fmt.Sprintf(" (%s)", model)
			}
			fmt.Printf("  %d. %s%s → %s\n", i+1, step.Workflow, modelInfo, step.NextStatus)
		}
		return nil
	}

	// Multiple stories dry-run - detailed output
	fmt.Printf("Dry run for %d stories:\n", len(storyKeys))

	totalWorkflows := 0
	storiesWithWork := 0
	storiesComplete := 0

	for _, storyKey := range storyKeys {
		fmt.Println()
		fmt.Printf("Story %s:\n", storyKey)

		steps, err := executor.GetSteps(storyKey)
		if err != nil {
			if errors.Is(err, router.ErrStoryComplete) {
				fmt.Printf("  (already complete)\n")
				storiesComplete++
				continue
			}
			cmd.SilenceUsage = true
			fmt.Printf("  Error: %v\n", err)
			return NewExitError(1)
		}

		for i, step := range steps {
			modelInfo := ""
			model := app.Config.GetModel(step.Workflow)
			if model != "" {
				modelInfo = fmt.Sprintf(" (%s)", model)
			}
			fmt.Printf("  %d. %s%s → %s\n", i+1, step.Workflow, modelInfo, step.NextStatus)
		}
		totalWorkflows += len(steps)
		storiesWithWork++
	}

	fmt.Println()
	if storiesComplete > 0 {
		fmt.Printf("Total: %d workflows across %d stories (%d already complete)\n", totalWorkflows, storiesWithWork, storiesComplete)
	} else {
		fmt.Printf("Total: %d workflows across %d stories\n", totalWorkflows, storiesWithWork)
	}

	return nil
}

// executeStoryTUI runs a story lifecycle using the TUI interface.
func executeStoryTUI(ctx context.Context, app *App, executor *lifecycle.Executor, storyKey string) error {
	// Get the steps for this story
	steps, err := executor.GetSteps(storyKey)
	if err != nil {
		if errors.Is(err, router.ErrStoryComplete) {
			fmt.Printf("Story %s is already complete\n", storyKey)
			return nil
		}
		return err
	}

	// Convert lifecycle steps to TUI steps
	tuiSteps := make([]tui.StepInfo, len(steps))
	for i, step := range steps {
		tuiSteps[i] = tui.StepInfo{
			Name:       step.Workflow,
			StoryKey:   storyKey,
			NextStatus: string(step.NextStatus),
		}
	}

	// Create TUI runner and execute
	tuiRunner := tui.NewRunner(app.Executor, app.Config)
	exitCode := tuiRunner.RunMultiStep(ctx, tuiSteps, storyKey)

	if exitCode != 0 {
		return NewExitError(exitCode)
	}

	return nil
}
