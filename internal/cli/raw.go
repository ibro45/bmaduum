package cli

import (
	"strings"

	"github.com/spf13/cobra"
)

func newRawCommand(app *App) *cobra.Command {
	return &cobra.Command{
		Use:   "raw <prompt>",
		Short: "Run an arbitrary prompt",
		Long: `Run an arbitrary prompt directly with Claude.
Useful for testing or one-off commands.

Example:
  bmad-automate raw "List all Go files in the project"`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			prompt := strings.Join(args, " ")
			ctx := cmd.Context()
			exitCode := app.Runner.RunRaw(ctx, prompt)
			if exitCode != 0 {
				cmd.SilenceUsage = true
				return NewExitError(exitCode)
			}
			return nil
		},
	}
}
