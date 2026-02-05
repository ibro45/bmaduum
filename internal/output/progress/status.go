// Package progress provides terminal progress display functionality.
package progress

import (
	"fmt"
	"io"
	"strings"
	"time"

	"bmaduum/internal/output/terminal"
)

// StatusBar handles rendering of the status bar (bottom row).
type StatusBar struct {
	cursor *terminal.Cursor
}

// NewStatusBar creates a new status bar renderer.
func NewStatusBar(out io.Writer) *StatusBar {
	return &StatusBar{
		cursor: terminal.NewCursor(out),
	}
}

// Render renders the status bar at the bottom row.
func (s *StatusBar) Render(width int, state StatusState) error {
	// Move to status bar line (bottom row)
	if err := s.cursor.MoveTo(width, 1); err != nil {
		return err
	}

	// Disable wrap, clear line, set colors
	if err := s.cursor.DisableWrap(); err != nil {
		return err
	}
	if _, err := s.cursor.WriteString(terminal.BgTerracotta + terminal.FgWhite + terminal.Bold); err != nil {
		return err
	}
	if err := s.cursor.ClearLine(); err != nil {
		return err
	}

	// Build and write content
	content := s.buildLine(state)
	if _, err := s.cursor.WriteString(content); err != nil {
		return err
	}

	// Pad to full width for complete background (use display width)
	contentWidth := displayWidth(content)
	if contentWidth < width {
		if _, err := s.cursor.WriteString(strings.Repeat(" ", width-contentWidth)); err != nil {
			return err
		}
	}

	// Reset
	if _, err := s.cursor.WriteString(terminal.ResetAttrs); err != nil {
		return err
	}
	if err := s.cursor.EnableWrap(); err != nil {
		return err
	}

	return nil
}

// buildLine constructs the status bar content with clear visual hierarchy.
// Format: ▸ Epic 6 │ Story 3/8 · 6-3-api │ Step 2/5 dev-story │ opus-4-5 │ 12:34
func (s *StatusBar) buildLine(state StatusState) string {
	available := state.Width - 2 // Leave some margin
	sep := " │ "                 // Box drawing separator for cleaner look

	// Total timer (always shown, goes at end)
	totalTimer := ""
	if !state.StartTime.IsZero() {
		totalTimer = formatDuration(time.Since(state.StartTime))
	}

	// Model (shortened if too long)
	model := shortenModel(state.Model)

	// Step info with timer
	var stepInfo string
	if state.Total > 0 {
		stepInfo = fmt.Sprintf("Step %d/%d %s", state.Step, state.Total, state.StepName)
	} else if state.StepName != "" {
		stepInfo = state.StepName
	}

	// Story key
	storyKey := state.StoryKey

	// Operation (e.g., "Epic 6", "Story 2/3")
	operation := state.Operation

	// Build parts based on available width
	// Priority order for narrow terminals: timer > operation > step > story > model
	var parts []string
	usedWidth := 0

	// Leader icon
	leader := "▸ "
	usedWidth += displayWidth(leader)

	// Always reserve space for timer
	timerWidth := displayWidth(totalTimer) + displayWidth(sep)
	available -= timerWidth

	// Add operation (highest priority context)
	if operation != "" {
		needed := displayWidth(operation) + displayWidth(sep)
		if usedWidth+needed <= available {
			parts = append(parts, operation)
			usedWidth += needed
		}
	}

	// Add story key (combine with operation if epic format detected)
	if storyKey != "" {
		needed := displayWidth(storyKey) + displayWidth(sep)
		if usedWidth+needed <= available {
			parts = append(parts, storyKey)
			usedWidth += needed
		}
	}

	// Add step info
	if stepInfo != "" {
		needed := displayWidth(stepInfo) + displayWidth(sep)
		if usedWidth+needed <= available {
			parts = append(parts, stepInfo)
			usedWidth += needed
		} else {
			// Try shorter version without "Step" prefix
			shortStep := ""
			if state.Total > 0 {
				shortStep = fmt.Sprintf("%d/%d %s", state.Step, state.Total, state.StepName)
			}
			if shortStep != "" {
				needed = displayWidth(shortStep) + displayWidth(sep)
				if usedWidth+needed <= available {
					parts = append(parts, shortStep)
					usedWidth += needed
				}
			}
		}
	}

	// Add model
	if model != "" {
		needed := displayWidth(model) + displayWidth(sep)
		if usedWidth+needed <= available {
			parts = append(parts, model)
		}
	}

	// Join parts and add timer at end
	content := leader + strings.Join(parts, sep)
	if totalTimer != "" {
		if len(parts) > 0 {
			content += sep + totalTimer
		} else {
			content += totalTimer
		}
	}

	return content
}

// shortenModel returns a shorter model name for display.
// "claude-opus-4-5-20251101" → "opus-4-5"
// "claude-sonnet-4-5" → "sonnet-4-5"
func shortenModel(model string) string {
	if model == "" {
		return ""
	}

	// Remove "claude-" prefix
	short := strings.TrimPrefix(model, "claude-")

	// Remove date suffix (e.g., "-20251101")
	if len(short) > 9 && short[len(short)-9] == '-' {
		// Check if the suffix is all digits
		suffix := short[len(short)-8:]
		allDigits := true
		for _, c := range suffix {
			if c < '0' || c > '9' {
				allDigits = false
				break
			}
		}
		if allDigits {
			short = short[:len(short)-9]
		}
	}

	return short
}
