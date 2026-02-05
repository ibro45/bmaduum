// Package terminal provides low-level terminal control primitives.
//
// This package encapsulates ANSI escape sequences and terminal manipulation
// functions used by the output package. It provides a clean abstraction
// over raw terminal control codes.
package terminal

// ANSI escape sequences for terminal control.
const (
	// DECSTBM - Set Top and Bottom Margins (scrolling region)
	SetScrollRegionFormat = "\x1b[%d;%dr"
	ResetScrollRegion     = "\x1b[r"

	// Cursor positioning
	SaveCursor    = "\x1b7"       // DEC save (more portable than CSI s)
	RestoreCursor = "\x1b8"       // DEC restore
	MoveToFormat  = "\x1b[%d;%dH" // Move to row, column
	ClearLine     = "\x1b[2K"     // Clear entire line
	CursorUp      = "\x1b[1A"     // Move cursor up one line
	HideCursor    = "\x1b[?25l"   // Hide cursor during updates
	ShowCursor    = "\x1b[?25h"   // Show cursor
	DisableWrap   = "\x1b[?7l"    // Disable line wrap
	EnableWrap    = "\x1b[?7h"    // Enable line wrap

	// Colors for status bar
	ResetAttrs   = "\x1b[0m"                // Reset all attributes
	BgTerracotta = "\x1b[48;2;193;95;60m"   // #C15F3C background
	FgWhite      = "\x1b[38;2;255;255;255m" // White foreground
	FgActivity   = "\x1b[38;2;255;107;107m" // #FF6B6B orange/red for activity
	Bold         = "\x1b[1m"                // Bold
)
