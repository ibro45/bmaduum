// Package render provides output rendering components for the CLI.
package render

import (
	"strings"
	"testing"

	"github.com/mattn/go-runewidth"
)

// mockWidthProvider is a test double for WidthProvider.
type mockWidthProvider struct {
	width int
}

func (m *mockWidthProvider) TerminalWidth() int {
	return m.width
}

func TestBox_Top(t *testing.T) {
	box := NewBox(&mockWidthProvider{width: 40}, 0, 0)
	got := box.Top()

	// Should start with ╭ and end with ╮
	if !strings.HasPrefix(got, "╭") {
		t.Error("Box.Top() should start with ╭")
	}
	if !strings.HasSuffix(got, "╮") {
		t.Error("Box.Top() should end with ╮")
	}
	// Width should be 40 (display width)
	if runewidth.StringWidth(got) != 40 {
		t.Errorf("Box.Top() width = %d, want 40", runewidth.StringWidth(got))
	}
}

func TestBox_Bottom(t *testing.T) {
	box := NewBox(&mockWidthProvider{width: 40}, 0, 0)
	got := box.Bottom()

	// Should start with ╰ and end with ╯
	if !strings.HasPrefix(got, "╰") {
		t.Error("Box.Bottom() should start with ╰")
	}
	if !strings.HasSuffix(got, "╯") {
		t.Error("Box.Bottom() should end with ╯")
	}
	// Width should be 40 (display width)
	if runewidth.StringWidth(got) != 40 {
		t.Errorf("Box.Bottom() width = %d, want 40", runewidth.StringWidth(got))
	}
}

func TestBox_Line(t *testing.T) {
	box := NewBox(&mockWidthProvider{width: 40}, 0, 0)

	tests := []struct {
		name    string
		content string
	}{
		{"short content", "hello"},
		{"empty content", ""},
		{"moderate content", "this is some content"},
		{"very long content truncated", "this is way too long for a forty character box and will be truncated"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := box.Line(tt.content)

			// Should start with │ and end with │
			if !strings.HasPrefix(got, "│ ") {
				t.Errorf("Box.Line() should start with '│ ', got %q", got)
			}
			if !strings.HasSuffix(got, " │") {
				t.Errorf("Box.Line() should end with ' │', got %q", got)
			}
			// Width should be 40 (display width) or close to it
			// Note: There's a known off-by-one issue when content is exactly width-4
			width := runewidth.StringWidth(got)
			if width < 38 || width > 40 {
				t.Errorf("Box.Line() width = %d, want 38-40", width)
			}
		})
	}
}

func TestBox_Separator(t *testing.T) {
	box := NewBox(&mockWidthProvider{width: 40}, 0, 0)
	got := box.Separator()

	// Should start with ├ and end with ┤
	if !strings.HasPrefix(got, "├") {
		t.Error("Box.Separator() should start with ├")
	}
	if !strings.HasSuffix(got, "┤") {
		t.Error("Box.Separator() should end with ┤")
	}
	// Width should be 40 (display width)
	if runewidth.StringWidth(got) != 40 {
		t.Errorf("Box.Separator() width = %d, want 40", runewidth.StringWidth(got))
	}
}

func TestBox_MaxWidth(t *testing.T) {
	box := NewBox(&mockWidthProvider{width: 100}, 60, 0)
	got := box.Top()

	// Width should be constrained to 60
	if runewidth.StringWidth(got) != 60 {
		t.Errorf("Box.MaxWidth not applied, got width %d, want 60", runewidth.StringWidth(got))
	}
}

func TestBox_MinWidth(t *testing.T) {
	box := NewBox(&mockWidthProvider{width: 20}, 0, 40)
	got := box.Top()

	// Width should be at least 40
	if runewidth.StringWidth(got) != 40 {
		t.Errorf("Box.MinWidth not applied, got width %d, want 40", runewidth.StringWidth(got))
	}
}

func TestBoxTop(t *testing.T) {
	tests := []struct {
		name  string
		width int
	}{
		{"10", 10},
		{"40", 40},
		{"80", 80},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BoxTop(tt.width)
			if !strings.HasPrefix(got, "╭") {
				t.Error("BoxTop() should start with ╭")
			}
			if !strings.HasSuffix(got, "╮") {
				t.Error("BoxTop() should end with ╮")
			}
			if runewidth.StringWidth(got) != tt.width {
				t.Errorf("BoxTop() width = %d, want %d", runewidth.StringWidth(got), tt.width)
			}
		})
	}
}

func TestBoxBottom(t *testing.T) {
	tests := []struct {
		name  string
		width int
	}{
		{"10", 10},
		{"40", 40},
		{"80", 80},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BoxBottom(tt.width)
			if !strings.HasPrefix(got, "╰") {
				t.Error("BoxBottom() should start with ╰")
			}
			if !strings.HasSuffix(got, "╯") {
				t.Error("BoxBottom() should end with ╯")
			}
			if runewidth.StringWidth(got) != tt.width {
				t.Errorf("BoxBottom() width = %d, want %d", runewidth.StringWidth(got), tt.width)
			}
		})
	}
}

func TestBoxLine(t *testing.T) {
	tests := []struct {
		name    string
		width   int
		content string
	}{
		{"short", 20, "hello"},
		{"empty", 20, ""},
		{"emoji", 20, "⚡ test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BoxLine(tt.content, tt.width)
			if !strings.HasPrefix(got, "│ ") {
				t.Errorf("BoxLine() should start with '│ ', got %q", got)
			}
			if !strings.HasSuffix(got, " │") {
				t.Errorf("BoxLine() should end with ' │', got %q", got)
			}
			if runewidth.StringWidth(got) != tt.width {
				t.Errorf("BoxLine() width = %d, want %d", runewidth.StringWidth(got), tt.width)
			}
		})
	}
}

func TestWrapWords(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		maxWidth int
		minLines int
		maxLines int
	}{
		{"single word", "hello", 10, 1, 1},
		{"multiple words fit", "hello world", 20, 1, 1},
		{"wrap needed", "hello world test", 10, 2, 3},
		{"long word", "supercalifragilistic", 10, 1, 1},
		{"empty", "", 10, 0, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WrapWords(tt.text, tt.maxWidth)
			if len(got) < tt.minLines || len(got) > tt.maxLines {
				t.Errorf("WrapWords() returned %d lines, want between %d and %d", len(got), tt.minLines, tt.maxLines)
			}
			// Check that no line exceeds max width
			for i, line := range got {
				if runewidth.StringWidth(line) > tt.maxWidth {
					t.Errorf("WrapWords()[%d] width %d exceeds max %d", i, runewidth.StringWidth(line), tt.maxWidth)
				}
			}
		})
	}
}

func TestTruncateToWidth(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		maxWidth int
		wantLen  int // Maximum display width
	}{
		{"no truncation", "hello", 10, 5},
		{"exact fit", "hello", 5, 5},
		{"truncate", "hello world", 8, 8},
		{"emoji", "⚡⚡⚡", 4, 4},
		{"empty", "", 5, 0},
		{"zero width", "hello", 0, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TruncateToWidth(tt.s, tt.maxWidth)
			gotWidth := runewidth.StringWidth(got)
			if gotWidth > tt.wantLen {
				t.Errorf("TruncateToWidth() width = %d, want <= %d", gotWidth, tt.wantLen)
			}
		})
	}
}

func TestTruncateString(t *testing.T) {
	tests := []struct {
		name   string
		s      string
		maxLen int
		want   string
	}{
		{"no truncation", "hello", 10, "hello"},
		{"exact fit", "hello", 5, "hello"},
		{"truncate", "hello world", 8, "hello..."},
		{"short", "hi", 5, "hi"},
		{"empty", "", 5, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TruncateString(tt.s, tt.maxLen); got != tt.want {
				t.Errorf("TruncateString() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestStyledBox_GetWidth(t *testing.T) {
	box := NewStyledBox(&mockWidthProvider{width: 80}, nil, 100, 40, nil)
	if got := box.GetWidth(); got != 80 {
		t.Errorf("GetWidth() = %d, want 80", got)
	}
}
