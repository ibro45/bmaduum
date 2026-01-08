package cli

import (
	"github.com/spf13/cobra"
)

func newRunCommand(app *App) *cobra.Command {
	return &cobra.Command{
		Use:   "run <story-key>",
		Short: "Run full development cycle",
		Long: `Run the full development cycle for a story:
  1. create-story - Create the story definition
  2. dev-story    - Implement the story
  3. code-review  - Review and fix any issues
  4. git-commit   - Commit and push changes`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			storyKey := args[0]
			ctx := cmd.Context()
			exitCode := app.Runner.RunFullCycle(ctx, storyKey)
			if exitCode != 0 {
				cmd.SilenceUsage = true
				return NewExitError(exitCode)
			}
			return nil
		},
	}
}
