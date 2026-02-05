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
	"os"

	"bmaduum/internal/output/terminal"
)

// Re-export terminal functions for backward compatibility.

// IsTTY returns true if the file descriptor refers to a terminal.
//
// This function uses golang.org/x/term to detect if the given file
// is connected to a terminal device. Returns false for pipes, redirected
// output, and non-interactive environments.
func IsTTY(f *os.File) bool {
	return terminal.IsTTY(f)
}

// SupportsColor returns true if the current terminal supports color output.
//
// Color support is disabled when:
//   - Output is not a TTY (piped to file or another process)
//   - NO_COLOR environment variable is set (per no-color.org)
//   - TERM is set to "dumb" (indicating basic terminal capabilities)
//
// Returns true otherwise, indicating ANSI color codes can be used.
func SupportsColor() bool {
	return terminal.SupportsColor()
}

// TerminalWidth returns the width of the terminal in columns.
//
// Returns the number of columns (characters) that fit in the terminal.
// If the width cannot be determined (not a TTY or error), returns
// a sensible default of 80 columns.
func TerminalWidth() int {
	return terminal.TerminalWidth()
}

// TerminalSize returns the dimensions of the terminal.
//
// Returns width and height in columns and rows. If the size cannot
// be determined, returns 80x24 as a default.
func TerminalSize() (width, height int) {
	return terminal.GetSize()
}

// IsWindows returns true if running on Windows.
//
// Uses runtime.GOOS to detect the operating system. Windows terminals
// may have different color and formatting capabilities.
func IsWindows() bool {
	return terminal.IsWindows()
}
