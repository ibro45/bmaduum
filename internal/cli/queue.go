package cli

import (
	"github.com/spf13/cobra"
)

func newQueueCommand(app *App) *cobra.Command {
	return &cobra.Command{
		Use:   "queue <story-key> [story-key...]",
		Short: "Run full cycle on multiple stories",
		Long: `Run the full development cycle on multiple stories in sequence.
The queue stops on the first failure.

Example:
  bmad-automate queue 6-5 6-6 6-7 6-8`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			exitCode := app.Queue.RunQueue(ctx, args)
			if exitCode != 0 {
				cmd.SilenceUsage = true
				return NewExitError(exitCode)
			}
			return nil
		},
	}
}
