package diff

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Renderer renders diffs with rich formatting including background colors,
// line numbers, and optional syntax highlighting.
type Renderer struct {
	lineNumWidth  int    // Width of line number column
	showLineNums  bool   // Whether to show line numbers
	syntaxHL      bool   // Whether to apply syntax highlighting
	language      string // Language for syntax highlighting
	showGutter    bool   // Whether to show gutter markers (┃) for changes
	maxContentLen int    // Maximum content length per line (0 = no limit)
	highlighter   HighlightFunc // Optional syntax highlighter function
	// Styles
	summaryStyle      lipgloss.Style
	contextLineStyle  lipgloss.Style
	lineNumStyle      lipgloss.Style
	gutterMarkerStyle lipgloss.Style
}

// HighlightFunc is a function that applies syntax highlighting to code.
// It takes code and a language identifier, and returns the highlighted code.
type HighlightFunc func(code, language string) string

// Option is a functional option for configuring Renderer.
type Option func(*Renderer)

// WithLineNumbers enables or disables line number display.
func WithLineNumbers(show bool) Option {
	return func(r *Renderer) {
		r.showLineNums = show
	}
}

// WithSyntaxHighlight enables syntax highlighting for the given language.
func WithSyntaxHighlight(language string) Option {
	return func(r *Renderer) {
		r.syntaxHL = language != ""
		r.language = language
	}
}

// WithGutter enables or disables gutter markers for changed lines.
func WithGutter(show bool) Option {
	return func(r *Renderer) {
		r.showGutter = show
	}
}

// WithMaxContentLength sets the maximum content length per line.
func WithMaxContentLength(n int) Option {
	return func(r *Renderer) {
		r.maxContentLen = n
	}
}

// WithHighlighter sets a custom syntax highlighter function.
func WithHighlighter(h HighlightFunc) Option {
	return func(r *Renderer) {
		r.highlighter = h
	}
}

// NewRenderer creates a renderer with the given options.
// Default configuration: line numbers enabled, gutter enabled, no syntax highlighting.
func NewRenderer(opts ...Option) *Renderer {
	r := &Renderer{
		lineNumWidth:  4,
		showLineNums:  true,
		syntaxHL:      false,
		showGutter:    true,
		maxContentLen: 0,
	}
	for _, opt := range opts {
		opt(r)
	}
	r.setStyles()
	return r
}

// setStyles configures the lipgloss styles based on terminal capabilities.
func (r *Renderer) setStyles() {
	if !supportsColor() {
		// No colors - use plain styles
		r.summaryStyle = lipgloss.NewStyle()
		r.contextLineStyle = lipgloss.NewStyle()
		r.lineNumStyle = lipgloss.NewStyle()
		r.gutterMarkerStyle = lipgloss.NewStyle()
		return
	}

	// Use rich colors for terminal output
	r.summaryStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#8b949e")).
		Italic(true)

	r.contextLineStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#8b949e"))

	r.lineNumStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6e7681"))

	r.gutterMarkerStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#58a6ff")).
		Bold(true)
}

// supportsColor checks if the terminal supports color output.
func supportsColor() bool {
	// Check if NO_COLOR environment variable is set
	if v := strings.ToLower(getEnv("NO_COLOR", "")); v != "" && v != "0" {
		return false
	}
	// Check if TERM is set to a known value
	term := getEnv("TERM", "")
	if term == "dumb" || term == "" {
		return false
	}
	return true
}

// getEnv gets an environment variable with a fallback value.
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// Render formats a diff for terminal display with rich formatting.
//
// The output includes:
//   - Line numbers in left gutter
//   - Background colors for added/deleted lines
//   - Gutter markers (┃) for changed sections
//   - Optional syntax highlighting
func (r *Renderer) Render(diff *Diff) string {
	if diff == nil {
		return ""
	}

	if !supportsColor() {
		return r.renderPlain(diff)
	}

	var buf strings.Builder

	// Calculate line number width based on max line number
	r.lineNumWidth = r.calculateLineNumWidth(diff)

	// Render each hunk
	for i, hunk := range diff.Hunks {
		if i > 0 {
			// Separator between hunks
			buf.WriteString(r.renderHunkSeparator())
		}
		r.renderHunk(&buf, &hunk)
	}

	return buf.String()
}

// RenderWithSummary returns the diff with a summary header.
func (r *Renderer) RenderWithSummary(diff *Diff) string {
	if diff == nil {
		return ""
	}

	var buf strings.Builder

	// Summary line
	summary := diff.Summary()
	buf.WriteString(r.summaryStyle.Render(summary))
	buf.WriteString("\n")

	// Render the diff content
	buf.WriteString(r.Render(diff))

	return buf.String()
}

// calculateLineNumWidth determines the width needed for line numbers.
func (r *Renderer) calculateLineNumWidth(diff *Diff) int {
	maxLine := 0
	for _, hunk := range diff.Hunks {
		end := hunk.NewStart + hunk.NewCount
		if end > maxLine {
			maxLine = end
		}
		end = hunk.OldStart + hunk.OldCount
		if end > maxLine {
			maxLine = end
		}
	}

	width := 1
	for maxLine >= 10 {
		width++
		maxLine /= 10
	}

	if width < 3 {
		width = 3
	}
	return width
}

// renderHunk renders a single hunk with line numbers and styling.
func (r *Renderer) renderHunk(buf *strings.Builder, hunk *Hunk) {
	for _, line := range hunk.Lines {
		r.renderLine(buf, &line)
	}
}

// renderLine renders a single diff line with appropriate styling.
func (r *Renderer) renderLine(buf *strings.Builder, line *Line) {
	// Build gutter marker (outside the background block)
	var gutter string
	if r.showGutter {
		switch line.Type {
		case LineTypeAdded, LineTypeDeleted:
			gutter = r.gutterMarkerStyle.Render("┃") + " "
		default:
			gutter = "  " // Two spaces to align with "┃ "
		}
	}

	// Build line number (plain, will be wrapped with background)
	var lineNum string
	if r.showLineNums {
		switch line.Type {
		case LineTypeAdded:
			lineNum = fmt.Sprintf("%*d ", r.lineNumWidth, line.NewLineNum)
		case LineTypeDeleted:
			lineNum = fmt.Sprintf("%*d ", r.lineNumWidth, line.OldLineNum)
		case LineTypeContext:
			lineNum = fmt.Sprintf("%*d ", r.lineNumWidth, line.NewLineNum)
			lineNum = r.lineNumStyle.Render(lineNum)
		}
	}

	// Build change marker (plain, will be wrapped with background)
	var marker string
	switch line.Type {
	case LineTypeAdded:
		marker = "+"
	case LineTypeDeleted:
		marker = "-"
	case LineTypeContext:
		marker = " "
	}

	// Apply syntax highlighting to content if enabled
	content := line.Content
	if r.maxContentLen > 0 {
		runes := []rune(content)
		if len(runes) > r.maxContentLen {
			content = string(runes[:r.maxContentLen-3]) + "..."
		}
	}
	if r.syntaxHL && r.language != "" {
		content = r.highlightCode(content, r.language)
	}

	// Build and style the full line content (line number + marker + content)
	// For added/deleted lines, wrap everything in background color
	var styledLine string
	switch line.Type {
	case LineTypeAdded:
		lineContent := lineNum + marker + " " + content
		if r.syntaxHL && r.language != "" {
			// Syntax HL provides foreground colors, just add background
			styledLine = applyDiffStyle(lineContent, "\x1b[48;2;13;40;24m", "")
		} else {
			// No syntax HL, apply both bg (#0d2818) and fg (#3fb950)
			styledLine = applyDiffStyle(lineContent, "\x1b[48;2;13;40;24m", "\x1b[38;2;63;185;80m")
		}
	case LineTypeDeleted:
		lineContent := lineNum + marker + " " + content
		if r.syntaxHL && r.language != "" {
			// Syntax HL provides foreground colors, just add background
			styledLine = applyDiffStyle(lineContent, "\x1b[48;2;45;10;10m", "")
		} else {
			// No syntax HL, apply both bg (#2d0a0a) and fg (#f85149)
			styledLine = applyDiffStyle(lineContent, "\x1b[48;2;45;10;10m", "\x1b[38;2;248;81;73m")
		}
	case LineTypeContext:
		if r.syntaxHL && r.language != "" {
			// Syntax HL: apply muted foreground that survives ANSI resets
			styledLine = lineNum + marker + " " + applyDiffStyle(content, "", "\x1b[38;2;139;148;158m")
		} else {
			styledLine = lineNum + marker + " " + r.contextLineStyle.Render(content)
		}
	}

	buf.WriteString(gutter)
	buf.WriteString(styledLine)
	buf.WriteString("\n")
}

// renderHunkSeparator renders the separator between hunks.
func (r *Renderer) renderHunkSeparator() string {
	if r.showGutter {
		return "  " + r.lineNumStyle.Render(strings.Repeat(" ", r.lineNumWidth)) + "  " +
			r.contextLineStyle.Render("...") + "\n"
	}
	return r.lineNumStyle.Render(strings.Repeat(" ", r.lineNumWidth)) + "  " +
		r.contextLineStyle.Render("...") + "\n"
}

// renderPlain renders the diff without colors for non-TTY output.
func (r *Renderer) renderPlain(diff *Diff) string {
	var buf strings.Builder

	for _, hunk := range diff.Hunks {
		for _, line := range hunk.Lines {
			var marker string
			switch line.Type {
			case LineTypeAdded:
				marker = "+"
			case LineTypeDeleted:
				marker = "-"
			case LineTypeContext:
				marker = " "
			}

			if r.showLineNums {
				var lineNum int
				if line.Type == LineTypeDeleted {
					lineNum = line.OldLineNum
				} else {
					lineNum = line.NewLineNum
				}
				buf.WriteString(fmt.Sprintf("%*d %s %s\n", r.lineNumWidth, lineNum, marker, line.Content))
			} else {
				buf.WriteString(fmt.Sprintf("%s %s\n", marker, line.Content))
			}
		}
	}

	return buf.String()
}

// applyDiffStyle applies background and optional foreground colors to content.
// It injects the style codes at the start and after each reset sequence
// to ensure the styling persists through syntax highlighting color changes.
// If fgCode is empty, only the background is applied (for syntax-highlighted content).
func applyDiffStyle(content, bgCode, fgCode string) string {
	const ansiReset = "\x1b[0m"
	styleCode := bgCode + fgCode
	result := styleCode + strings.ReplaceAll(content, ansiReset, ansiReset+styleCode) + ansiReset
	return result
}

// highlightCode applies syntax highlighting to code using the configured highlighter.
func (r *Renderer) highlightCode(code, language string) string {
	if r.highlighter == nil {
		return code
	}
	return r.highlighter(code, language)
}
