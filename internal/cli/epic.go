package cli

import (
	"github.com/spf13/cobra"
)

func newEpicCommand(app *App) *cobra.Command {
	return &cobra.Command{
		Use:   "epic <epic-id>",
		Short: "Run appropriate workflow for all stories in an epic",
		Long: `Run the appropriate workflow for all stories in an epic based on their status in sprint-status.yaml.

Finds all stories matching the pattern {epic-id}-{N}-* where N is numeric,
sorts them by story number, and runs them in order.

Status routing:
  - backlog       → create-story
  - ready-for-dev → dev-story
  - in-progress   → dev-story
  - review        → code-review
  - done          → skipped (story complete)

The epic command stops on the first failure. Done stories are skipped and do not cause failure.

Example:
  bmad-automate epic 6
  # Runs 6-1-*, 6-2-*, 6-3-*, etc. in order`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			epicID := args[0]

			// Get all stories for this epic
			storyKeys, err := app.StatusReader.GetEpicStories(epicID)
			if err != nil {
				cmd.SilenceUsage = true
				return NewExitError(1)
			}

			// Run all stories using the queue runner (which handles status-based routing)
			exitCode := app.Queue.RunQueueWithStatus(ctx, storyKeys, app.StatusReader)
			if exitCode != 0 {
				cmd.SilenceUsage = true
				return NewExitError(exitCode)
			}
			return nil
		},
	}
}
