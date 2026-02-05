package output

import (
	"strings"

	"github.com/charmbracelet/glamour"
)

// MarkdownConfig contains configuration for markdown rendering.
type MarkdownConfig struct {
	Enabled  bool
	Style    string
	WordWrap int
	Emoji    bool
}

// DefaultMarkdownConfig returns sensible defaults for markdown rendering.
func DefaultMarkdownConfig() MarkdownConfig {
	return MarkdownConfig{
		Enabled:  true,
		Style:    "dark",
		WordWrap: 100,
		Emoji:    true,
	}
}

// MarkdownRenderer renders markdown for terminal display.
type MarkdownRenderer struct {
	renderer *glamour.TermRenderer
	enabled  bool
}

// NewMarkdownRenderer creates a renderer with default configuration.
// Uses dark theme by default to avoid auto-detection delays.
func NewMarkdownRenderer() *MarkdownRenderer {
	return NewMarkdownRendererWithConfig(DefaultMarkdownConfig())
}

// NewMarkdownRendererWithConfig creates a renderer with custom configuration.
func NewMarkdownRendererWithConfig(cfg MarkdownConfig) *MarkdownRenderer {
	if !cfg.Enabled {
		return &MarkdownRenderer{enabled: false}
	}

	// Check if terminal supports color
	if !SupportsColor() {
		return &MarkdownRenderer{enabled: false}
	}

	opts := []glamour.TermRendererOption{
		glamour.WithWordWrap(cfg.WordWrap),
	}

	if cfg.Emoji {
		opts = append(opts, glamour.WithEmoji())
	}

	// Use specified style (avoid "auto" for performance)
	style := cfg.Style
	if style == "" || style == "auto" {
		style = "dark"
	}
	opts = append(opts, glamour.WithStandardStyle(style))

	r, err := glamour.NewTermRenderer(opts...)
	if err != nil {
		return &MarkdownRenderer{enabled: false}
	}
	return &MarkdownRenderer{renderer: r, enabled: true}
}

// Render converts markdown to styled terminal output.
func (m *MarkdownRenderer) Render(markdown string) string {
	if !m.enabled || markdown == "" || !looksLikeMarkdown(markdown) {
		return markdown
	}

	out, err := m.renderer.Render(markdown)
	if err != nil {
		return markdown // Fallback to raw
	}

	// Glamour adds trailing newlines; trim them for cleaner integration
	return strings.TrimSuffix(out, "\n")
}

// looksLikeMarkdown checks if text contains markdown syntax.
// Used to skip rendering overhead for plain text messages.
func looksLikeMarkdown(s string) bool {
	// Quick check for common markdown characters
	if !strings.ContainsAny(s, "*_`#[]()|->.0123456789") {
		return false
	}

	// More specific patterns to avoid false positives
	// Check for actual markdown constructs
	patterns := []string{
		"**",   // Bold
		"__",   // Bold alt
		"*",    // Italic (single asterisk with content)
		"_",    // Italic alt
		"```",  // Code block
		"`",    // Inline code
		"# ",   // Header
		"## ",  // H2
		"### ", // H3
		"- ",   // List item
		"* ",   // List item alt
		"[",    // Link start
		"![",   // Image
		"| ",   // Table
		"---",  // HR
		"***",  // HR alt
		"> ",   // Blockquote
	}

	for _, p := range patterns {
		if strings.Contains(s, p) {
			return true
		}
	}

	// Check for ordered list at start of string or after newline
	if strings.HasPrefix(s, "1. ") || strings.HasPrefix(s, "2. ") || strings.HasPrefix(s, "3. ") {
		return true
	}
	if strings.Contains(s, "\n1. ") || strings.Contains(s, "\n2. ") || strings.Contains(s, "\n3. ") {
		return true
	}

	return false
}
