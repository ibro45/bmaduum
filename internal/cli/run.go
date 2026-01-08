package cli

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"bmad-automate/internal/router"
)

func newRunCommand(app *App) *cobra.Command {
	return &cobra.Command{
		Use:   "run <story-key>",
		Short: "Run the appropriate workflow based on story status",
		Long: `Run the appropriate workflow for a story based on its status in sprint-status.yaml:
  - backlog       → create-story
  - ready-for-dev → dev-story
  - in-progress   → dev-story
  - review        → code-review
  - done          → no action (story complete)`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			storyKey := args[0]
			ctx := cmd.Context()

			// Get story status from sprint-status.yaml
			status, err := app.StatusReader.GetStoryStatus(storyKey)
			if err != nil {
				cmd.SilenceUsage = true
				fmt.Printf("Error: %v\n", err)
				return NewExitError(1)
			}

			// Route to appropriate workflow
			workflowName, err := router.GetWorkflow(status)
			if err != nil {
				cmd.SilenceUsage = true
				if errors.Is(err, router.ErrStoryComplete) {
					fmt.Printf("Story %s is already complete, no action needed\n", storyKey)
					return nil
				}
				if errors.Is(err, router.ErrUnknownStatus) {
					fmt.Printf("Error: unknown status value: %s\n", status)
					return NewExitError(1)
				}
				fmt.Printf("Error: %v\n", err)
				return NewExitError(1)
			}

			// Run the appropriate workflow
			exitCode := app.Runner.RunSingle(ctx, workflowName, storyKey)
			if exitCode != 0 {
				cmd.SilenceUsage = true
				return NewExitError(exitCode)
			}
			return nil
		},
	}
}
