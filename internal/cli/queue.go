package cli

import (
	"github.com/spf13/cobra"
)

func newQueueCommand(app *App) *cobra.Command {
	return &cobra.Command{
		Use:   "queue <story-key> [story-key...]",
		Short: "Run appropriate workflow for multiple stories based on status",
		Long: `Run the appropriate workflow for multiple stories based on their status in sprint-status.yaml:
  - backlog       → create-story
  - ready-for-dev → dev-story
  - in-progress   → dev-story
  - review        → code-review
  - done          → skipped (story complete)

The queue stops on the first failure. Done stories are skipped and do not cause failure.

Example:
  bmad-automate queue 6-5 6-6 6-7 6-8`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			exitCode := app.Queue.RunQueueWithStatus(ctx, args, app.StatusReader)
			if exitCode != 0 {
				cmd.SilenceUsage = true
				return NewExitError(exitCode)
			}
			return nil
		},
	}
}
