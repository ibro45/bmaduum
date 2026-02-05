// Package terminal provides low-level terminal control primitives.
//
// This package encapsulates ANSI escape sequences and terminal manipulation
// functions used by the output package. It provides a clean abstraction
// over raw terminal control codes.
package terminal

import (
	"fmt"
	"io"
	"sync"
)

// Cursor provides cursor positioning and control methods.
type Cursor struct {
	out io.Writer
	mu  sync.Mutex
}

// NewCursor creates a new Cursor instance for the given writer.
func NewCursor(out io.Writer) *Cursor {
	return &Cursor{out: out}
}

// Save saves the current cursor position.
func (c *Cursor) Save() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, err := fmt.Fprint(c.out, SaveCursor)
	return err
}

// Restore restores the previously saved cursor position.
func (c *Cursor) Restore() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, err := fmt.Fprint(c.out, RestoreCursor)
	return err
}

// MoveTo moves the cursor to the specified row and column (1-indexed).
func (c *Cursor) MoveTo(row, col int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, err := fmt.Fprintf(c.out, MoveToFormat, row, col)
	return err
}

// ClearLine clears the entire line the cursor is on.
func (c *Cursor) ClearLine() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, err := fmt.Fprint(c.out, ClearLine)
	return err
}

// Up moves the cursor up by the specified number of lines.
func (c *Cursor) Up(lines int) error {
	if lines <= 0 {
		return nil
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// Use the standard CursorUp sequence repeated
	for i := 0; i < lines; i++ {
		if _, err := fmt.Fprint(c.out, CursorUp); err != nil {
			return err
		}
	}
	return nil
}

// Hide hides the cursor.
func (c *Cursor) Hide() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, err := fmt.Fprint(c.out, HideCursor)
	return err
}

// Show shows the cursor.
func (c *Cursor) Show() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, err := fmt.Fprint(c.out, ShowCursor)
	return err
}

// DisableWrap disables line wrapping.
func (c *Cursor) DisableWrap() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, err := fmt.Fprint(c.out, DisableWrap)
	return err
}

// EnableWrap enables line wrapping.
func (c *Cursor) EnableWrap() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, err := fmt.Fprint(c.out, EnableWrap)
	return err
}

// Write writes data directly to the underlying writer.
func (c *Cursor) Write(p []byte) (n int, err error) {
	return c.out.Write(p)
}

// WriteString writes a string directly to the underlying writer.
func (c *Cursor) WriteString(s string) (n int, err error) {
	return io.WriteString(c.out, s)
}
