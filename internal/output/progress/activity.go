// Package progress provides terminal progress display functionality.
package progress

import (
	"fmt"
	"io"
	"strings"
	"time"

	"bmaduum/internal/output/terminal"
)

// ActivityLine handles rendering of the activity line (second-to-last row).
// Format matches Claude Code: · Thinking… (2m 30s · ↓ 13.4k tokens · thought for 4s)
type ActivityLine struct {
	cursor *terminal.Cursor
}

// NewActivityLine creates a new activity line renderer.
func NewActivityLine(out io.Writer) *ActivityLine {
	return &ActivityLine{
		cursor: terminal.NewCursor(out),
	}
}

// Render renders the activity line at height-1.
// Format: ⚡ Thinking… (2m 30s · 47 tools · ↓ 13.4k tokens)
func (a *ActivityLine) Render(height int, state ActivityState) error {
	// Move to activity line (second-to-last row)
	if err := a.cursor.MoveTo(height-1, 1); err != nil {
		return err
	}
	if err := a.cursor.ClearLine(); err != nil {
		return err
	}

	// Build activity content with bmaduum icon
	verb := getThinkingVerb(state.VerbIdx)
	activity := iconBmaduum + " " + verb + "…"

	// Build the parenthesized info section
	var parts []string

	// Timer (elapsed time for this activity)
	if !state.ActivityStart.IsZero() {
		parts = append(parts, formatDurationNatural(time.Since(state.ActivityStart)))
	}

	// Tool count (shows activity level)
	if state.ToolCount > 0 {
		toolText := "tool"
		if state.ToolCount != 1 {
			toolText = "tools"
		}
		parts = append(parts, fmt.Sprintf("%d %s", state.ToolCount, toolText))
	}

	// Token count
	totalTokens := state.InputTokens + state.OutputTokens
	if totalTokens > 0 {
		parts = append(parts, "↓ "+formatTokenCount(totalTokens)+" tokens")
	}

	// Thinking time (only show if we had a response and recorded thinking duration)
	if state.HadFirstResponse && state.ThinkingDuration > 0 {
		parts = append(parts, "thought for "+formatDurationNatural(state.ThinkingDuration))
	}

	// Join parts with middot separator
	if len(parts) > 0 {
		activity += " (" + strings.Join(parts, " · ") + ")"
	}

	// Write activity in orange/red color
	if _, err := a.cursor.WriteString(terminal.FgActivity + terminal.Bold); err != nil {
		return err
	}
	if _, err := a.cursor.WriteString(activity); err != nil {
		return err
	}
	if _, err := a.cursor.WriteString(terminal.ResetAttrs); err != nil {
		return err
	}

	return nil
}

// RenderDone renders the completion message before clearing.
// Format: ✓ Thought for 2m 30s · 47 tools · ↓ 13.4k tokens
func (a *ActivityLine) RenderDone(height int, success bool, duration time.Duration, totalTokens, toolCount, verbIdx int) error {
	if err := a.cursor.MoveTo(height-1, 1); err != nil {
		return err
	}
	if err := a.cursor.ClearLine(); err != nil {
		return err
	}

	var msg string
	if success {
		// Pick a past-tense thinking verb (use current verb index)
		pastTenseVerb := getPastTenseVerb(verbIdx)
		msg = iconSuccess + " " + pastTenseVerb + " for " + formatDurationNatural(duration)
		if toolCount > 0 {
			toolText := "tool"
			if toolCount != 1 {
				toolText = "tools"
			}
			msg += fmt.Sprintf(" · %d %s", toolCount, toolText)
		}
		if totalTokens > 0 {
			msg += " · ↓ " + formatTokenCount(totalTokens) + " tokens"
		}
	} else {
		msg = iconError + " Failed"
		if toolCount > 0 {
			toolText := "tool"
			if toolCount != 1 {
				toolText = "tools"
			}
			msg += fmt.Sprintf(" · %d %s", toolCount, toolText)
		}
		if totalTokens > 0 {
			msg += " · ↓ " + formatTokenCount(totalTokens) + " tokens"
		}
	}

	// Write in activity color
	if _, err := a.cursor.WriteString(terminal.FgActivity + terminal.Bold); err != nil {
		return err
	}
	if _, err := a.cursor.WriteString(msg); err != nil {
		return err
	}
	if _, err := a.cursor.WriteString(terminal.ResetAttrs); err != nil {
		return err
	}

	return nil
}
