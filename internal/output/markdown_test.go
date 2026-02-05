package output

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMarkdownRenderer(t *testing.T) {
	r := NewMarkdownRenderer()
	assert.NotNil(t, r)
}

func TestNewMarkdownRendererWithConfig_Disabled(t *testing.T) {
	cfg := MarkdownConfig{Enabled: false}
	r := NewMarkdownRendererWithConfig(cfg)
	assert.NotNil(t, r)
	assert.False(t, r.enabled)
}

func TestNewMarkdownRendererWithConfig_Enabled(t *testing.T) {
	cfg := DefaultMarkdownConfig()
	r := NewMarkdownRendererWithConfig(cfg)
	assert.NotNil(t, r)
	// enabled depends on terminal support
}

func TestMarkdownRenderer_Render_Empty(t *testing.T) {
	r := NewMarkdownRenderer()
	result := r.Render("")
	assert.Equal(t, "", result)
}

func TestMarkdownRenderer_Render_PlainText(t *testing.T) {
	r := NewMarkdownRenderer()
	// Plain text without markdown should pass through unchanged
	result := r.Render("Hello world")
	assert.Equal(t, "Hello world", result)
}

func TestMarkdownRenderer_Render_Disabled(t *testing.T) {
	r := &MarkdownRenderer{enabled: false}
	result := r.Render("**bold**")
	assert.Equal(t, "**bold**", result)
}

func TestMarkdownRenderer_Render_Bold(t *testing.T) {
	cfg := DefaultMarkdownConfig()
	r := NewMarkdownRendererWithConfig(cfg)
	if !r.enabled {
		t.Skip("markdown rendering not enabled (no color support)")
	}

	result := r.Render("**bold text**")
	// Should contain ANSI bold codes or the rendered text
	assert.NotEqual(t, "**bold text**", result)
	assert.Contains(t, result, "bold text")
}

func TestMarkdownRenderer_Render_CodeBlock(t *testing.T) {
	cfg := DefaultMarkdownConfig()
	r := NewMarkdownRendererWithConfig(cfg)
	if !r.enabled {
		t.Skip("markdown rendering not enabled (no color support)")
	}

	result := r.Render("```go\nfunc main() {}\n```")
	// Should contain the code, possibly with syntax highlighting
	assert.Contains(t, result, "func")
	assert.Contains(t, result, "main")
}

func TestMarkdownRenderer_Render_InlineCode(t *testing.T) {
	cfg := DefaultMarkdownConfig()
	r := NewMarkdownRendererWithConfig(cfg)
	if !r.enabled {
		t.Skip("markdown rendering not enabled (no color support)")
	}

	result := r.Render("Use `fmt.Println` for output")
	assert.Contains(t, result, "fmt.Println")
}

func TestMarkdownRenderer_Render_Header(t *testing.T) {
	cfg := DefaultMarkdownConfig()
	r := NewMarkdownRendererWithConfig(cfg)
	if !r.enabled {
		t.Skip("markdown rendering not enabled (no color support)")
	}

	result := r.Render("# Title\n\nSome content")
	assert.Contains(t, result, "Title")
	assert.Contains(t, result, "content")
}

func TestMarkdownRenderer_Render_List(t *testing.T) {
	cfg := DefaultMarkdownConfig()
	r := NewMarkdownRendererWithConfig(cfg)
	if !r.enabled {
		t.Skip("markdown rendering not enabled (no color support)")
	}

	result := r.Render("- Item 1\n- Item 2\n- Item 3")
	assert.Contains(t, result, "Item 1")
	assert.Contains(t, result, "Item 2")
	assert.Contains(t, result, "Item 3")
}

func TestLooksLikeMarkdown(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		// Plain text - should be false
		{"Hello world", false},
		{"Simple text without formatting", false},
		{"Numbers: 123, 456", false},

		// Bold
		{"**bold text**", true},
		{"__also bold__", true},

		// Italic
		{"*italic*", true},
		{"_also italic_", true},

		// Code
		{"`inline code`", true},
		{"```\ncode block\n```", true},

		// Headers
		{"# Header", true},
		{"## Subheader", true},
		{"### H3", true},

		// Lists
		{"- list item", true},
		{"* another item", true},
		{"1. numbered", true},

		// Links
		{"[link](url)", true},
		{"![image](url)", true},

		// Tables
		{"| col1 | col2 |", true},

		// Horizontal rules
		{"---", true},
		{"***", true},

		// Blockquotes
		{"> quoted text", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := looksLikeMarkdown(tt.input)
			assert.Equal(t, tt.expected, result, "looksLikeMarkdown(%q)", tt.input)
		})
	}
}

func TestDefaultMarkdownConfig(t *testing.T) {
	cfg := DefaultMarkdownConfig()
	assert.True(t, cfg.Enabled)
	assert.Equal(t, "dark", cfg.Style)
	assert.Equal(t, 100, cfg.WordWrap)
	assert.True(t, cfg.Emoji)
}

func TestMarkdownRenderer_Render_NoTrailingNewlines(t *testing.T) {
	cfg := DefaultMarkdownConfig()
	r := NewMarkdownRendererWithConfig(cfg)
	if !r.enabled {
		t.Skip("markdown rendering not enabled (no color support)")
	}

	result := r.Render("**test**")
	// Should not end with newline
	assert.False(t, strings.HasSuffix(result, "\n"), "should not have trailing newline")
}
