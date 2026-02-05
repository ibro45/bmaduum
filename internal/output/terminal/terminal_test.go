// Package terminal provides low-level terminal control primitives.
package terminal

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestTerminal_New(t *testing.T) {
	t.Run("with bytes.Buffer", func(t *testing.T) {
		buf := &bytes.Buffer{}
		term := New(buf)

		if term == nil {
			t.Fatal("New() returned nil")
		}

		if term.FileDescriptor() != -1 {
			t.Errorf("expected fd -1 for non-file writer, got %d", term.FileDescriptor())
		}

		w, h := term.Size()
		if w != 0 || h != 0 {
			t.Errorf("expected size 0x0 for non-TTY, got %dx%d", w, h)
		}
	})

	t.Run("with os.File (non-TTY)", func(t *testing.T) {
		f, err := os.CreateTemp("", "terminal-test")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(f.Name())
		defer f.Close()

		term := New(f)

		if term.FileDescriptor() == -1 {
			t.Error("expected valid fd for os.File")
		}
	})
}

func TestTerminal_SetScrollRegion(t *testing.T) {
	buf := &bytes.Buffer{}
	term := New(buf)

	err := term.SetScrollRegion(1, 24)
	if err != nil {
		t.Fatalf("SetScrollRegion() error = %v", err)
	}

	expected := "\x1b[1;24r"
	if got := buf.String(); got != expected {
		t.Errorf("SetScrollRegion() = %q, want %q", got, expected)
	}
}

func TestTerminal_ResetScrollRegion(t *testing.T) {
	buf := &bytes.Buffer{}
	term := New(buf)

	err := term.ResetScrollRegion()
	if err != nil {
		t.Fatalf("ResetScrollRegion() error = %v", err)
	}

	expected := "\x1b[r"
	if got := buf.String(); got != expected {
		t.Errorf("ResetScrollRegion() = %q, want %q", got, expected)
	}
}

func TestTerminal_Clear(t *testing.T) {
	buf := &bytes.Buffer{}
	term := New(buf)

	err := term.Clear()
	if err != nil {
		t.Fatalf("Clear() error = %v", err)
	}

	expected := "\x1b[2K"
	if got := buf.String(); got != expected {
		t.Errorf("Clear() = %q, want %q", got, expected)
	}
}

func TestCursor_MoveTo(t *testing.T) {
	buf := &bytes.Buffer{}
	cursor := NewCursor(buf)

	err := cursor.MoveTo(5, 10)
	if err != nil {
		t.Fatalf("MoveTo() error = %v", err)
	}

	expected := "\x1b[5;10H"
	if got := buf.String(); got != expected {
		t.Errorf("MoveTo() = %q, want %q", got, expected)
	}
}

func TestCursor_SaveRestore(t *testing.T) {
	buf := &bytes.Buffer{}
	cursor := NewCursor(buf)

	if err := cursor.Save(); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	if err := cursor.Restore(); err != nil {
		t.Fatalf("Restore() error = %v", err)
	}

	expected := "\x1b7\x1b8"
	if got := buf.String(); got != expected {
		t.Errorf("Save+Restore = %q, want %q", got, expected)
	}
}

func TestCursor_ClearLine(t *testing.T) {
	buf := &bytes.Buffer{}
	cursor := NewCursor(buf)

	err := cursor.ClearLine()
	if err != nil {
		t.Fatalf("ClearLine() error = %v", err)
	}

	expected := "\x1b[2K"
	if got := buf.String(); got != expected {
		t.Errorf("ClearLine() = %q, want %q", got, expected)
	}
}

func TestCursor_Up(t *testing.T) {
	t.Run("single line", func(t *testing.T) {
		buf := &bytes.Buffer{}
		cursor := NewCursor(buf)

		err := cursor.Up(1)
		if err != nil {
			t.Fatalf("Up(1) error = %v", err)
		}

		expected := "\x1b[1A"
		if got := buf.String(); got != expected {
			t.Errorf("Up(1) = %q, want %q", got, expected)
		}
	})

	t.Run("multiple lines", func(t *testing.T) {
		buf := &bytes.Buffer{}
		cursor := NewCursor(buf)

		err := cursor.Up(3)
		if err != nil {
			t.Fatalf("Up(3) error = %v", err)
		}

		expected := strings.Repeat("\x1b[1A", 3)
		if got := buf.String(); got != expected {
			t.Errorf("Up(3) = %q, want %q", got, expected)
		}
	})

	t.Run("zero lines", func(t *testing.T) {
		buf := &bytes.Buffer{}
		cursor := NewCursor(buf)

		err := cursor.Up(0)
		if err != nil {
			t.Fatalf("Up(0) error = %v", err)
		}

		if got := buf.String(); got != "" {
			t.Errorf("Up(0) = %q, want empty string", got)
		}
	})
}

func TestCursor_HideShow(t *testing.T) {
	buf := &bytes.Buffer{}
	cursor := NewCursor(buf)

	if err := cursor.Hide(); err != nil {
		t.Fatalf("Hide() error = %v", err)
	}

	if err := cursor.Show(); err != nil {
		t.Fatalf("Show() error = %v", err)
	}

	expected := "\x1b[?25l\x1b[?25h"
	if got := buf.String(); got != expected {
		t.Errorf("Hide+Show = %q, want %q", got, expected)
	}
}

func TestCursor_WrapControl(t *testing.T) {
	buf := &bytes.Buffer{}
	cursor := NewCursor(buf)

	if err := cursor.DisableWrap(); err != nil {
		t.Fatalf("DisableWrap() error = %v", err)
	}

	if err := cursor.EnableWrap(); err != nil {
		t.Fatalf("EnableWrap() error = %v", err)
	}

	expected := "\x1b[?7l\x1b[?7h"
	if got := buf.String(); got != expected {
		t.Errorf("DisableWrap+EnableWrap = %q, want %q", got, expected)
	}
}

func TestCursor_Write(t *testing.T) {
	buf := &bytes.Buffer{}
	cursor := NewCursor(buf)

	n, err := cursor.WriteString("test")
	if err != nil {
		t.Fatalf("WriteString() error = %v", err)
	}

	if n != 4 {
		t.Errorf("WriteString() returned %d bytes, want 4", n)
	}

	if got := buf.String(); got != "test" {
		t.Errorf("WriteString() = %q, want %q", got, "test")
	}
}

func TestTerminalWidth_DefaultOnNonTTY(t *testing.T) {
	// When run in a non-TTY context (like CI), this should return default
	width := TerminalWidth()
	if width <= 0 {
		t.Errorf("TerminalWidth() = %d, want > 0", width)
	}
}

func TestGetSize_DefaultOnNonTTY(t *testing.T) {
	w, h := GetSize()
	if w <= 0 || h <= 0 {
		t.Errorf("GetSize() = %dx%d, want positive dimensions", w, h)
	}
}

func TestIsWindows(t *testing.T) {
	// This just verifies the function runs without panic
	result := IsWindows()
	_ = result // result depends on the test platform
}

func TestIsTTY(t *testing.T) {
	f, err := os.CreateTemp("", "tty-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	defer f.Close()

	// Temp file is not a TTY
	if IsTTY(f) {
		t.Error("temp file should not be a TTY")
	}
}

func TestSupportsColor(t *testing.T) {
	// Just verify it runs without panic
	result := SupportsColor()
	_ = result
}
