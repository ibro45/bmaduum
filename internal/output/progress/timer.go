// Package progress provides terminal progress display functionality.
package progress

import (
	"fmt"
	"time"

	"github.com/mattn/go-runewidth"
)

// formatDuration formats a duration as MM:SS or H:MM:SS.
func formatDuration(d time.Duration) string {
	if d < time.Second {
		return "0:00"
	}

	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60

	if d < time.Hour {
		return fmt.Sprintf("%d:%02d", minutes, seconds)
	}

	hours := int(d.Hours())
	minutes = int(d.Minutes()) % 60
	return fmt.Sprintf("%d:%02d:%02d", hours, minutes, seconds)
}

// formatDurationNatural formats a duration in Claude Code's natural style.
// Examples: "4s", "2m 30s", "1h 5m"
func formatDurationNatural(d time.Duration) string {
	if d < time.Second {
		return "0s"
	}

	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if hours > 0 {
		if minutes > 0 {
			return fmt.Sprintf("%dh %dm", hours, minutes)
		}
		return fmt.Sprintf("%dh", hours)
	}

	if minutes > 0 {
		if seconds > 0 {
			return fmt.Sprintf("%dm %ds", minutes, seconds)
		}
		return fmt.Sprintf("%dm", minutes)
	}

	return fmt.Sprintf("%ds", seconds)
}

// formatTokenCount formats a token count for display (e.g., 5700 -> "5.7k").
func formatTokenCount(count int) string {
	if count < 1000 {
		return fmt.Sprintf("%d", count)
	}
	if count < 1000000 {
		return fmt.Sprintf("%.1fk", float64(count)/1000)
	}
	return fmt.Sprintf("%.2fM", float64(count)/1000000)
}

// displayWidth returns the visual width of a string in terminal cells,
// accounting for wide characters like emoji and CJK.
func displayWidth(s string) int {
	return runewidth.StringWidth(s)
}

// truncateToDisplayWidth truncates a string to fit within maxWidth display cells.
func truncateToDisplayWidth(s string, maxWidth int) string {
	if maxWidth <= 0 {
		return ""
	}
	if displayWidth(s) <= maxWidth {
		return s
	}
	// Truncate rune by rune until it fits
	runes := []rune(s)
	for i := len(runes) - 1; i >= 0; i-- {
		truncated := string(runes[:i]) + "…"
		if displayWidth(truncated) <= maxWidth {
			return truncated
		}
	}
	return ""
}
