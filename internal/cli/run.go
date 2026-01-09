package cli

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"bmad-automate/internal/lifecycle"
	"bmad-automate/internal/router"
)

func newRunCommand(app *App) *cobra.Command {
	return &cobra.Command{
		Use:   "run <story-key>",
		Short: "Run the full story lifecycle to completion",
		Long: `Run the complete lifecycle for a story from its current status to done.

The command executes all remaining workflows based on the story's current status:
  - backlog       → create-story → dev-story → code-review → git-commit → done
  - ready-for-dev → dev-story → code-review → git-commit → done
  - in-progress   → dev-story → code-review → git-commit → done
  - review        → code-review → git-commit → done
  - done          → no action (story already complete)

Status is updated in sprint-status.yaml after each successful workflow.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			storyKey := args[0]
			ctx := cmd.Context()

			// Create lifecycle executor with app dependencies
			executor := lifecycle.NewExecutor(app.Runner, app.StatusReader, app.StatusWriter)

			// Execute the full lifecycle
			err := executor.Execute(ctx, storyKey)
			if err != nil {
				cmd.SilenceUsage = true
				if errors.Is(err, router.ErrStoryComplete) {
					fmt.Printf("Story %s is already complete, no action needed\n", storyKey)
					return nil
				}
				fmt.Printf("Error: %v\n", err)
				return NewExitError(1)
			}

			return nil
		},
	}
}
