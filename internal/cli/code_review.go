package cli

import (
	"github.com/spf13/cobra"
)

func newCodeReviewCommand(app *App) *cobra.Command {
	return &cobra.Command{
		Use:   "code-review <story-key>",
		Short: "Run code-review workflow",
		Long:  `Run the code-review workflow for the specified story key.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			storyKey := args[0]
			ctx := cmd.Context()
			exitCode := app.Runner.RunSingle(ctx, "code-review", storyKey)
			if exitCode != 0 {
				cmd.SilenceUsage = true
				return NewExitError(exitCode)
			}
			return nil
		},
	}
}
