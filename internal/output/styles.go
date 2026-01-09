// Package output provides terminal output formatting using lipgloss styles.
//
// The package provides structured output for CLI operations including session
// lifecycle, step progress, tool usage display, and batch operation summaries.
// All output is styled using the lipgloss library for consistent terminal rendering.
//
// Key types:
//   - [Printer] - Interface for structured terminal output operations
//   - [DefaultPrinter] - Production implementation using lipgloss styles
//   - [StepResult] - Result of a single workflow step execution
//   - [StoryResult] - Result of processing a story in queue/epic operations
//
// Use [NewPrinter] for production output to stdout, or [NewPrinterWithWriter]
// to capture output in tests by providing a custom io.Writer.
package output

import (
	"github.com/charmbracelet/lipgloss"
)

// Color palette for terminal output styling.
// These are internal colors used by the style definitions below.
var (
	colorPrimary   = lipgloss.Color("39")  // Bright blue - headers, borders
	colorSuccess   = lipgloss.Color("42")  // Green - success indicators
	colorError     = lipgloss.Color("196") // Red - error indicators
	colorWarning   = lipgloss.Color("214") // Orange - tool names
	colorMuted     = lipgloss.Color("245") // Gray - secondary info, dividers
	colorHighlight = lipgloss.Color("177") // Purple - labels, queue headers
)

// Lipgloss styles for different output elements.
// These styles are used by [DefaultPrinter] methods.
var (
	// headerStyle formats major section headers with double borders.
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorPrimary).
			BorderStyle(lipgloss.DoubleBorder()).
			BorderForeground(colorPrimary).
			Padding(0, 1)

	// stepHeaderStyle formats step progress headers with rounded borders.
	stepHeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorPrimary).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(colorMuted).
			Padding(0, 1)

	// successStyle formats success indicators (checkmarks, completion messages).
	successStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorSuccess)

	// errorStyle formats error indicators (X marks, failure messages).
	errorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorError)

	// mutedStyle formats secondary information (stderr, pending items).
	mutedStyle = lipgloss.NewStyle().
			Foreground(colorMuted)

	// labelStyle formats command labels and highlights.
	labelStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorHighlight)

	// toolNameStyle formats Claude tool names in output.
	toolNameStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorWarning)

	// dividerStyle formats visual separator lines.
	dividerStyle = lipgloss.NewStyle().
			Foreground(colorMuted)

	// summaryStyle formats completion summary boxes.
	summaryStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.DoubleBorder()).
			BorderForeground(colorPrimary).
			Padding(0, 1)

	// queueHeaderStyle formats queue operation headers.
	queueHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(colorHighlight).
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(colorHighlight).
				Padding(0, 1)
)

// Icons for status indicators in terminal output.
// These provide visual feedback for operation states.
const (
	iconSuccess    = "✓"  // Completed successfully
	iconError      = "✗"  // Failed
	iconPending    = "○"  // Not yet started
	iconInProgress = "●"  // Currently running
	iconTool       = "┌─" // Tool block start
	iconToolEnd    = "└─" // Tool block end
	iconToolLine   = "│"  // Tool block continuation
)
