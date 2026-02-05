// Package render provides output rendering components for the CLI.
package render

import (
	"fmt"
	"strings"
	"time"
)

// CycleStyleProvider provides styling functions for rendered output.
type CycleStyleProvider interface {
	RenderHeader(s string) string
	RenderSuccess(s string) string
	RenderError(s string) string
	RenderMuted(s string) string
	RenderDivider(s string) string
}

// CycleWidthProvider provides terminal width information.
type CycleWidthProvider interface {
	TerminalWidth() int
}

// OutputWriter provides output writing capability.
type OutputWriter interface {
	Writeln(format string, args ...interface{})
	Divider()
}

// CycleRenderer handles rendering for cycle and queue operations.
type CycleRenderer struct {
	writer OutputWriter
	styles CycleStyleProvider
	width  CycleWidthProvider
	box    *Box
}

// NewCycleRenderer creates a new cycle renderer.
func NewCycleRenderer(writer OutputWriter, styles CycleStyleProvider, width CycleWidthProvider) *CycleRenderer {
	box := NewBox(width, 80, 40)
	return &CycleRenderer{
		writer: writer,
		styles: styles,
		width:  width,
		box:    box,
	}
}

// CycleHeader prints the header for a full cycle run.
func (r *CycleRenderer) CycleHeader(storyKey string) {
	width := r.width.TerminalWidth()
	if width > 80 {
		width = 80
	}

	r.writer.Writeln(r.styles.RenderHeader(BoxTop(width)))
	r.writer.Writeln(r.styles.RenderHeader(BoxLine(IconBmaduum+" BMAD Full Cycle", width)))
	r.writer.Writeln(r.styles.RenderHeader(BoxLine("Story: "+storyKey, width)))
	r.writer.Writeln(r.styles.RenderHeader(BoxLine("Steps: create-story -> dev-story -> code-review -> git-commit", width)))
	r.writer.Writeln(r.styles.RenderHeader(BoxBottom(width)))
}

// CycleSummary prints the summary after a successful cycle.
func (r *CycleRenderer) CycleSummary(storyKey string, steps []StepResult, totalDuration time.Duration) {
	width := r.width.TerminalWidth()
	if width > 70 {
		width = 70
	}

	r.writer.Writeln("")
	r.writer.Writeln(r.styles.RenderSuccess(BoxTop(width)))
	r.writer.Writeln(r.styles.RenderSuccess(BoxLine(IconSuccess+" CYCLE COMPLETE", width)))
	r.writer.Writeln(r.styles.RenderSuccess(BoxLine("Story: "+storyKey, width)))
	r.writer.Writeln(r.styles.RenderSuccess("├" + strings.Repeat("─", width-2) + "┤"))

	for i, step := range steps {
		line := fmt.Sprintf("[%d] %-15s %s", i+1, step.Name, step.Duration.Round(time.Millisecond))
		r.writer.Writeln(r.styles.RenderSuccess(BoxLine(line, width)))
	}

	r.writer.Writeln(r.styles.RenderSuccess("├" + strings.Repeat("─", width-2) + "┤"))
	r.writer.Writeln(r.styles.RenderSuccess(BoxLine(fmt.Sprintf("Total: %s", totalDuration.Round(time.Millisecond)), width)))
	r.writer.Writeln(r.styles.RenderSuccess(BoxBottom(width)))
}

// CycleFailed prints failure information when a cycle fails.
func (r *CycleRenderer) CycleFailed(storyKey string, failedStep string, duration time.Duration) {
	width := r.width.TerminalWidth()
	if width > 60 {
		width = 60
	}

	r.writer.Writeln("")
	r.writer.Writeln(r.styles.RenderError(BoxTop(width)))
	r.writer.Writeln(r.styles.RenderError(BoxLine(IconError+" CYCLE FAILED", width)))
	r.writer.Writeln(r.styles.RenderError(BoxLine("Story: "+storyKey, width)))
	r.writer.Writeln(r.styles.RenderError(BoxLine("Failed at: "+failedStep, width)))
	r.writer.Writeln(r.styles.RenderError(BoxLine("Duration: "+duration.Round(time.Millisecond).String(), width)))
	r.writer.Writeln(r.styles.RenderError(BoxBottom(width)))
}

// QueueHeader prints the header for a queue run.
func (r *CycleRenderer) QueueHeader(count int, stories []string) {
	width := r.width.TerminalWidth()
	if width > 80 {
		width = 80
	}

	r.writer.Writeln(r.styles.RenderHeader(BoxTop(width)))
	r.writer.Writeln(r.styles.RenderHeader(BoxLine(fmt.Sprintf("%s BMAD Queue: %d stories", IconBmaduum, count), width)))

	// Show stories (wrap at word boundaries)
	storiesStr := strings.Join(stories, ", ")
	storyLines := BoxLineWrapWords("Stories", storiesStr, width)
	for _, line := range storyLines {
		r.writer.Writeln(r.styles.RenderMuted(line))
	}

	r.writer.Writeln(r.styles.RenderHeader(BoxBottom(width)))
}

// QueueStoryStart prints the header for starting a story in a queue.
func (r *CycleRenderer) QueueStoryStart(index, total int, storyKey string) {
	r.writer.Writeln("")
	r.writer.Divider()
	r.writer.Writeln(r.styles.RenderHeader(fmt.Sprintf("  Queue [%d/%d]: %s", index, total, storyKey)))
	r.writer.Divider()
}

// QueueSummary prints the summary after a queue completes or fails.
func (r *CycleRenderer) QueueSummary(results []StoryResult, allKeys []string, totalDuration time.Duration) {
	completed := 0
	failed := 0
	skipped := 0
	for _, r := range results {
		if r.Skipped {
			skipped++
		} else if r.Success {
			completed++
		} else {
			failed++
		}
	}
	remaining := len(allKeys) - len(results)

	width := r.width.TerminalWidth()
	if width > 70 {
		width = 70
	}

	r.writer.Writeln("")

	// Header
	if failed == 0 && remaining == 0 {
		r.writer.Writeln(r.styles.RenderSuccess(BoxTop(width)))
		r.writer.Writeln(r.styles.RenderSuccess(BoxLine(IconSuccess+" QUEUE COMPLETE", width)))
	} else {
		r.writer.Writeln(r.styles.RenderError(BoxTop(width)))
		r.writer.Writeln(r.styles.RenderError(BoxLine(IconError+" QUEUE STOPPED", width)))
	}

	// Summary line
	summaryLine := fmt.Sprintf("Completed: %d | Skipped: %d | Failed: %d | Remaining: %d",
		completed, skipped, failed, remaining)
	r.writer.Writeln(r.styles.RenderMuted(BoxLine(summaryLine, width)))

	// Separator
	if failed == 0 && remaining == 0 {
		r.writer.Writeln(r.styles.RenderSuccess("├" + strings.Repeat("─", width-2) + "┤"))
	} else {
		r.writer.Writeln(r.styles.RenderError("├" + strings.Repeat("─", width-2) + "┤"))
	}

	// Results
	for _, result := range results {
		var status, suffix string
		if result.Skipped {
			status = r.styles.RenderMuted("↷")
			suffix = "(done)"
		} else if result.Success {
			status = r.styles.RenderSuccess(IconSuccess)
			suffix = result.Duration.Round(time.Second).String()
		} else {
			status = r.styles.RenderError(IconError)
			suffix = result.Duration.Round(time.Second).String()
		}
		line := fmt.Sprintf("%s %-30s %s", status, result.Key, suffix)
		r.writer.Writeln(BoxLine(line, width))
	}

	// Remaining
	if remaining > 0 {
		for i := len(results); i < len(allKeys); i++ {
			line := fmt.Sprintf("%s %-30s (pending)", r.styles.RenderMuted(IconPending), allKeys[i])
			r.writer.Writeln(BoxLine(line, width))
		}
	}

	// Footer
	if failed == 0 && remaining == 0 {
		r.writer.Writeln(r.styles.RenderSuccess("├" + strings.Repeat("─", width-2) + "┤"))
		r.writer.Writeln(r.styles.RenderSuccess(BoxLine("Total: "+totalDuration.Round(time.Second).String(), width)))
		r.writer.Writeln(r.styles.RenderSuccess(BoxBottom(width)))
	} else {
		r.writer.Writeln(r.styles.RenderError("├" + strings.Repeat("─", width-2) + "┤"))
		r.writer.Writeln(r.styles.RenderError(BoxLine("Total: "+totalDuration.Round(time.Second).String(), width)))
		r.writer.Writeln(r.styles.RenderError(BoxBottom(width)))
	}
}
