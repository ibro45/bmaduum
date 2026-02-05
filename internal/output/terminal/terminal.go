// Package terminal provides low-level terminal control primitives.
//
// This package encapsulates ANSI escape sequences and terminal manipulation
// functions used by the output package. It provides a clean abstraction
// over raw terminal control codes.
package terminal

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"

	"golang.org/x/term"
)

// Terminal represents a terminal output device with low-level control methods.
type Terminal struct {
	out    io.Writer
	fd     int
	mu     sync.Mutex
	width  int
	height int
}

// New creates a new Terminal instance for the given writer.
// If the writer is an *os.File, it will be used for TTY detection
// and terminal size queries.
func New(out io.Writer) *Terminal {
	t := &Terminal{
		out: out,
		fd:  -1,
	}

	if f, ok := out.(*os.File); ok {
		t.fd = int(f.Fd())
		if IsTTY(f) {
			t.width, t.height = GetSize()
		}
	}

	return t
}

// IsTTY returns true if the file descriptor refers to a terminal.
//
// This function uses golang.org/x/term to detect if the given file
// is connected to a terminal device. Returns false for pipes, redirected
// output, and non-interactive environments.
func IsTTY(f *os.File) bool {
	return term.IsTerminal(int(f.Fd()))
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
	if !IsTTY(os.Stdout) {
		return false
	}

	// Respect NO_COLOR convention (https://no-color.org/)
	if os.Getenv("NO_COLOR") != "" {
		return false
	}

	// Check for dumb terminal
	termEnv := os.Getenv("TERM")
	if termEnv == "dumb" {
		return false
	}

	return true
}

// TerminalWidth returns the width of the terminal in columns.
//
// Returns the number of columns (characters) that fit in the terminal.
// If the width cannot be determined (not a TTY or error), returns
// a sensible default of 80 columns.
func TerminalWidth() int {
	if !IsTTY(os.Stdout) {
		return 80
	}

	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 80
	}

	return width
}

// GetSize returns the dimensions of the terminal.
//
// Returns width and height in columns and rows. If the size cannot
// be determined, returns 80x24 as a default.
func GetSize() (width, height int) {
	if !IsTTY(os.Stdout) {
		return 80, 24
	}

	w, h, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 80, 24
	}

	return w, h
}

// IsWindows returns true if running on Windows.
//
// Uses runtime.GOOS to detect the operating system. Windows terminals
// may have different color and formatting capabilities.
func IsWindows() bool {
	return runtime.GOOS == "windows"
}

// SetScrollRegion sets the scrolling region of the terminal.
// This reserves lines outside the region for status displays.
func (t *Terminal) SetScrollRegion(top, bottom int) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	_, err := fmt.Fprintf(t.out, SetScrollRegionFormat, top, bottom)
	return err
}

// ResetScrollRegion resets the scrolling region to the full terminal.
func (t *Terminal) ResetScrollRegion() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	_, err := fmt.Fprint(t.out, ResetScrollRegion)
	return err
}

// Clear clears the entire current line.
func (t *Terminal) Clear() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	_, err := fmt.Fprint(t.out, ClearLine)
	return err
}

// Size returns the current terminal dimensions.
func (t *Terminal) Size() (width, height int) {
	t.mu.Lock()
	defer t.mu.Unlock()

	return t.width, t.height
}

// UpdateSize refreshes the cached terminal size.
func (t *Terminal) UpdateSize() (width, height int) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.width, t.height = GetSize()
	return t.width, t.height
}

// FileDescriptor returns the file descriptor for the terminal, or -1 if not available.
func (t *Terminal) FileDescriptor() int {
	return t.fd
}

// Writer returns the underlying io.Writer.
func (t *Terminal) Writer() io.Writer {
	return t.out
}
