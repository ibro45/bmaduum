package workflow

import (
	"context"
	"errors"
	"fmt"
	"time"

	"bmad-automate/internal/output"
	"bmad-automate/internal/router"
	"bmad-automate/internal/status"
)

// QueueRunner processes multiple stories in sequence.
type QueueRunner struct {
	runner *Runner
}

// NewQueueRunner creates a new queue runner.
func NewQueueRunner(runner *Runner) *QueueRunner {
	return &QueueRunner{runner: runner}
}

// RunQueueWithStatus executes the appropriate workflow for each story based on status.
// Done stories are skipped. It stops on the first failure.
func (q *QueueRunner) RunQueueWithStatus(ctx context.Context, storyKeys []string, statusReader *status.Reader) int {
	queueStart := time.Now()
	results := make([]output.StoryResult, 0, len(storyKeys))

	q.runner.printer.QueueHeader(len(storyKeys), storyKeys)

	for i, storyKey := range storyKeys {
		q.runner.printer.QueueStoryStart(i+1, len(storyKeys), storyKey)

		storyStart := time.Now()

		// Get story status
		storyStatus, err := statusReader.GetStoryStatus(storyKey)
		if err != nil {
			fmt.Printf("  Error: %v\n", err)
			result := output.StoryResult{
				Key:      storyKey,
				Success:  false,
				Duration: time.Since(storyStart),
				FailedAt: "status",
			}
			results = append(results, result)
			q.runner.printer.QueueSummary(results, storyKeys, time.Since(queueStart))
			return 1
		}

		// Route to appropriate workflow
		workflowName, err := router.GetWorkflow(storyStatus)
		if err != nil {
			if errors.Is(err, router.ErrStoryComplete) {
				// Done stories are skipped, not failures
				fmt.Printf("  ↷ Skipped (already done)\n")
				result := output.StoryResult{
					Key:      storyKey,
					Success:  true,
					Duration: time.Since(storyStart),
					Skipped:  true,
				}
				results = append(results, result)
				fmt.Println() // Add spacing between stories
				continue
			}
			if errors.Is(err, router.ErrUnknownStatus) {
				fmt.Printf("  Error: unknown status value: %s\n", storyStatus)
			} else {
				fmt.Printf("  Error: %v\n", err)
			}
			result := output.StoryResult{
				Key:      storyKey,
				Success:  false,
				Duration: time.Since(storyStart),
				FailedAt: "routing",
			}
			results = append(results, result)
			q.runner.printer.QueueSummary(results, storyKeys, time.Since(queueStart))
			return 1
		}

		// Run the workflow
		exitCode := q.runner.RunSingle(ctx, workflowName, storyKey)
		duration := time.Since(storyStart)

		result := output.StoryResult{
			Key:      storyKey,
			Success:  exitCode == 0,
			Duration: duration,
		}

		if exitCode != 0 {
			result.FailedAt = workflowName
			results = append(results, result)
			q.runner.printer.QueueSummary(results, storyKeys, time.Since(queueStart))
			return exitCode
		}

		results = append(results, result)
		fmt.Println() // Add spacing between stories
	}

	q.runner.printer.QueueSummary(results, storyKeys, time.Since(queueStart))
	return 0
}
