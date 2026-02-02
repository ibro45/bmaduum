package cli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"bmad-automate/internal/sprint"
	"bmad-automate/internal/status"
)

// newSprintCommand creates the sprint command with subcommands.
//
// The sprint command provides operations for managing sprint status,
// including rebuilding the sprint-status.yaml cache from story files.
func newSprintCommand(app *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sprint",
		Short: "Manage sprint status",
		Long: `Manage sprint status and cache operations.

The sprint command provides subcommands for working with the sprint-status.yaml
file, which serves as a cache of story statuses. The source of truth is always
the story files in _bmad-output/implementation-artifacts/stories/.`,
	}

	cmd.AddCommand(newSprintRebuildCommand(app))

	return cmd
}

// newSprintRebuildCommand creates the rebuild subcommand.
//
// The rebuild command regenerates sprint-status.yaml by reading all story
// files and extracting their current status. This is useful for recovery
// when the cache becomes out of sync with the story files.
func newSprintRebuildCommand(app *App) *cobra.Command {
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "rebuild",
		Short: "Rebuild sprint-status.yaml from story files",
		Long: `Rebuild the sprint-status.yaml cache by reading all story files.

This command scans all story files in _bmad-output/implementation-artifacts/stories/,
extracts their current status, and regenerates the sprint-status.yaml file.

Use this command to recover when sprint-status.yaml becomes out of sync with
the story files (which are the source of truth).

Use --dry-run to preview changes without writing the file.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			_ = ctx // May be used in future for cancellation

			storyDir := sprint.DefaultStoryDir
			sprintPath := status.DefaultStatusPath

			// Get the base path from the status manager if available
			// For now, use default paths

			// Read all story files
			entries, err := os.ReadDir(storyDir)
			if err != nil {
				cmd.SilenceUsage = true
				return fmt.Errorf("failed to read story directory: %w", err)
			}

			developmentStatus := make(map[string]string)
			var skipped []string

			for _, entry := range entries {
				if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
					continue
				}

				storyKey := strings.TrimSuffix(entry.Name(), ".md")
				storyPath := filepath.Join(storyDir, entry.Name())

				st, err := readStoryStatusLine(storyPath)
				if err != nil {
					skipped = append(skipped, fmt.Sprintf("%s (%v)", storyKey, err))
					continue
				}

				if !status.Status(st).IsValid() {
					skipped = append(skipped, fmt.Sprintf("%s (invalid status: %s)", storyKey, st))
					continue
				}

				developmentStatus[storyKey] = st
			}

			// Build the YAML structure
			sprintStatus := map[string]interface{}{
				"development_status": developmentStatus,
			}

			if dryRun {
				fmt.Printf("Dry run - would write %d statuses to %s:\n", len(developmentStatus), sprintPath)
				for key, st := range developmentStatus {
					fmt.Printf("  %s: %s\n", key, st)
				}
				if len(skipped) > 0 {
					fmt.Printf("\nSkipped %d files:\n", len(skipped))
					for _, s := range skipped {
						fmt.Printf("  - %s\n", s)
					}
				}
				return nil
			}

			// Marshal to YAML
			data, err := yaml.Marshal(sprintStatus)
			if err != nil {
				cmd.SilenceUsage = true
				return fmt.Errorf("failed to marshal sprint status: %w", err)
			}

			// Write atomically
			tmpPath := sprintPath + ".tmp"
			if err := os.WriteFile(tmpPath, data, 0644); err != nil {
				cmd.SilenceUsage = true
				return fmt.Errorf("failed to write sprint status: %w", err)
			}

			if err := os.Rename(tmpPath, sprintPath); err != nil {
				os.Remove(tmpPath)
				cmd.SilenceUsage = true
				return fmt.Errorf("failed to rename sprint status file: %w", err)
			}

			fmt.Printf("Rebuilt %s with %d story statuses\n", sprintPath, len(developmentStatus))
			if len(skipped) > 0 {
				fmt.Printf("Skipped %d files\n", len(skipped))
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview changes without writing the file")

	return cmd
}

// readStoryStatusLine reads the status from a story file.
// Similar to StoryStatusManager.readStoryStatus but standalone to avoid import cycle.
func readStoryStatusLine(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(strings.ToLower(trimmed), "status:") {
			parts := strings.SplitN(trimmed, ":", 2)
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1]), nil
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", fmt.Errorf("status line not found")
}
