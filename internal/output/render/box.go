// Package render provides output rendering components for the CLI.
package render

import (
	"fmt"
	"strings"

	"github.com/mattn/go-runewidth"
)

// WidthProvider provides terminal width information.
type WidthProvider interface {
	TerminalWidth() int
}

// Box provides rounded box drawing capabilities.
type Box struct {
	widthProvider WidthProvider
	maxWidth      int
	minWidth      int
}

// NewBox creates a new box renderer with the given width provider.
// The box will be constrained to maxWidth and minWidth if specified.
func NewBox(provider WidthProvider, maxWidth, minWidth int) *Box {
	return &Box{
		widthProvider: provider,
		maxWidth:      maxWidth,
		minWidth:      minWidth,
	}
}

// getWidth returns the appropriate box width based on terminal size.
func (b *Box) getWidth() int {
	width := b.widthProvider.TerminalWidth()

	// Apply max width constraint
	if b.maxWidth > 0 && width > b.maxWidth {
		width = b.maxWidth
	}

	// Apply min width constraint
	if b.minWidth > 0 && width < b.minWidth {
		width = b.minWidth
	}

	return width
}

// Top returns the top of a rounded box.
func (b *Box) Top() string {
	width := b.getWidth()
	return "╭" + strings.Repeat("─", width-2) + "╮"
}

// Bottom returns the bottom of a rounded box.
func (b *Box) Bottom() string {
	width := b.getWidth()
	return "╰" + strings.Repeat("─", width-2) + "╯"
}

// Line returns a line inside the box, padded to width.
// Uses display width for accurate padding with Unicode characters.
func (b *Box) Line(content string) string {
	width := b.getWidth()
	contentWidth := runewidth.StringWidth(content)
	padding := width - 4 - contentWidth // 4 = "│ " + " │"
	if padding < 0 {
		// Truncate content if too long
		content = TruncateToWidth(content, width-5)
		padding = 0
	}
	return "│ " + content + strings.Repeat(" ", padding) + " │"
}

// Separator returns a separator line inside the box.
func (b *Box) Separator() string {
	width := b.getWidth()
	return "├" + strings.Repeat("─", width-2) + "┤"
}

// WrapWords wraps text at word boundaries to fit within maxWidth.
// Returns a slice of lines, each fitting within maxWidth display cells.
func WrapWords(text string, maxWidth int) []string {
	if maxWidth <= 0 {
		return []string{text}
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{}
	}

	var lines []string
	var currentLine strings.Builder
	currentWidth := 0

	for i, word := range words {
		wordWidth := runewidth.StringWidth(word)

		// If this is the first word on the line
		if currentWidth == 0 {
			// If word itself is too long, truncate it
			if wordWidth > maxWidth {
				lines = append(lines, TruncateToWidth(word, maxWidth))
				continue
			}
			currentLine.WriteString(word)
			currentWidth = wordWidth
			continue
		}

		// Check if adding this word (with space) would exceed max
		if currentWidth+1+wordWidth > maxWidth {
			// Start a new line
			lines = append(lines, currentLine.String())
			currentLine.Reset()

			// If word itself is too long, truncate it
			if wordWidth > maxWidth {
				lines = append(lines, TruncateToWidth(word, maxWidth))
				currentWidth = 0
				continue
			}

			currentLine.WriteString(word)
			currentWidth = wordWidth
		} else {
			// Add word to current line
			if i > 0 {
				currentLine.WriteString(" ")
				currentWidth++
			}
			currentLine.WriteString(word)
			currentWidth += wordWidth
		}
	}

	// Don't forget the last line
	if currentLine.Len() > 0 {
		lines = append(lines, currentLine.String())
	}

	return lines
}

// WrapWordsInBox wraps content at word boundaries inside the box.
// The first line contains "Label:" and subsequent lines are indented.
func (b *Box) WrapWordsInBox(label, content string) []string {
	width := b.getWidth()
	var lines []string
	innerWidth := width - 4 // "│ " + " │"

	// First line: "Label:"
	labelLine := label + ":"
	lines = append(lines, b.Line(labelLine))

	// Content lines: indented
	indent := "  " // 2 spaces indent for content
	contentWidth := innerWidth - runewidth.StringWidth(indent)

	wrappedContent := WrapWords(content, contentWidth)
	for _, line := range wrappedContent {
		lines = append(lines, b.Line(indent+line))
	}

	return lines
}

// TruncateToWidth truncates a string to fit within maxWidth display cells.
// If the string is too long, it truncates and adds "…" at the end.
func TruncateToWidth(s string, maxWidth int) string {
	if maxWidth <= 0 {
		return ""
	}
	if runewidth.StringWidth(s) <= maxWidth {
		return s
	}
	// Truncate rune by rune until it fits
	runes := []rune(s)
	for i := len(runes) - 1; i >= 0; i-- {
		truncated := string(runes[:i]) + "…"
		if runewidth.StringWidth(truncated) <= maxWidth {
			return truncated
		}
	}
	return ""
}

// TruncateString truncates a string to maxLen, adding "..." if truncated.
// This is a simple byte-length truncation, not display-width aware.
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// BoxTop returns the top of a rounded box with a specific width.
// This is a convenience function for backward compatibility.
func BoxTop(width int) string {
	return "╭" + strings.Repeat("─", width-2) + "╮"
}

// BoxTop is an alias for the Box.Top method - convenience function.
func boxTop(width int) string {
	return BoxTop(width)
}

// BoxBottom returns the bottom of a rounded box with a specific width.
// This is a convenience function for backward compatibility.
func BoxBottom(width int) string {
	return "╰" + strings.Repeat("─", width-2) + "╯"
}

// BoxLine returns a line inside the box with a specific width.
// This is a convenience function for backward compatibility.
func BoxLine(content string, width int) string {
	contentWidth := runewidth.StringWidth(content)
	padding := width - 4 - contentWidth // 4 = "│ " + " │"
	if padding < 0 {
		// Truncate content if too long
		content = TruncateToWidth(content, width-5)
		padding = 0
	}
	return "│ " + content + strings.Repeat(" ", padding) + " │"
}

// BoxLineWrapWords wraps content at word boundaries inside the box.
// This is a convenience function for backward compatibility.
func BoxLineWrapWords(label, content string, width int) []string {
	var lines []string
	innerWidth := width - 4 // "│ " + " │"

	// First line: "Label:"
	labelLine := label + ":"
	lines = append(lines, BoxLine(labelLine, width))

	// Content lines: indented
	indent := "  " // 2 spaces indent for content
	contentWidth := innerWidth - runewidth.StringWidth(indent)

	wrappedContent := WrapWords(content, contentWidth)
	for _, line := range wrappedContent {
		lines = append(lines, BoxLine(indent+line, width))
	}

	return lines
}

// BoxLineStyled returns a box line with separate styles for border and content.
// borderStyle is applied to the │ characters, contentStyle is applied to the inner content.
func BoxLineStyled(content string, width int, borderStyle, contentStyle func(string) string) string {
	contentWidth := runewidth.StringWidth(content)
	padding := width - 4 - contentWidth // 4 = "│ " + " │"
	if padding < 0 {
		content = TruncateToWidth(content, width-5)
		padding = 0
	}
	return borderStyle("│") + " " + contentStyle(content+strings.Repeat(" ", padding)) + " " + borderStyle("│")
}

// BoxLineWrapWordsStyled wraps content at word boundaries inside the box with separate styles.
// Returns lines with borderStyle applied to │ characters and contentStyle applied to inner content.
func BoxLineWrapWordsStyled(label, content string, width int, borderStyle, contentStyle func(string) string) []string {
	var lines []string
	innerWidth := width - 4 // "│ " + " │"

	// First line: "Label:"
	labelLine := label + ":"
	lines = append(lines, BoxLineStyled(labelLine, width, borderStyle, contentStyle))

	// Content lines: indented
	indent := "  " // 2 spaces indent for content
	contentWidth := innerWidth - runewidth.StringWidth(indent)

	wrappedContent := WrapWords(content, contentWidth)
	for _, line := range wrappedContent {
		lines = append(lines, BoxLineStyled(indent+line, width, borderStyle, contentStyle))
	}

	return lines
}

// StyleProvider applies styling to box lines.
type BoxStyleProvider interface {
	RenderHeader(s string) string
	RenderSuccess(s string) string
	RenderError(s string) string
	RenderMuted(s string) string
	RenderDivider(s string) string
}

// StyledBox renders styled boxes for various output elements.
type StyledBox struct {
	box     *Box
	styles  BoxStyleProvider
	writeln func(string) // Function to write output
}

// NewStyledBox creates a new styled box renderer.
// The writeln function controls where output is written.
func NewStyledBox(widthProvider WidthProvider, styles BoxStyleProvider, maxWidth, minWidth int, writeln func(string)) *StyledBox {
	if writeln == nil {
		// Default to stdout if no writer provided
		writeln = func(s string) { fmt.Println(s) }
	}
	return &StyledBox{
		box:     NewBox(widthProvider, maxWidth, minWidth),
		styles:  styles,
		writeln: writeln,
	}
}

// RenderHeaderBox renders a header-style box.
func (sb *StyledBox) RenderHeaderBox(lines []string) {
	sb.writeln(sb.styles.RenderHeader(sb.box.Top()))
	for _, line := range lines {
		sb.writeln(sb.styles.RenderHeader(sb.box.Line(line)))
	}
	sb.writeln(sb.styles.RenderHeader(sb.box.Bottom()))
}

// RenderSuccessBox renders a success-style box.
func (sb *StyledBox) RenderSuccessBox(lines []string) {
	sb.writeln(sb.styles.RenderSuccess(sb.box.Top()))
	for _, line := range lines {
		sb.writeln(sb.styles.RenderSuccess(sb.box.Line(line)))
	}
	sb.writeln(sb.styles.RenderSuccess(sb.box.Bottom()))
}

// RenderErrorBox renders an error-style box.
func (sb *StyledBox) RenderErrorBox(lines []string) {
	sb.writeln(sb.styles.RenderError(sb.box.Top()))
	for _, line := range lines {
		sb.writeln(sb.styles.RenderError(sb.box.Line(line)))
	}
	sb.writeln(sb.styles.RenderError(sb.box.Bottom()))
}

// RenderMixedBox renders a box with mixed styling (header line, separator, body).
func (sb *StyledBox) RenderMixedBox(header, separator string, body []string, bodyStyle func(string) string) {
	sb.writeln(sb.styles.RenderHeader(sb.box.Top()))
	sb.writeln(sb.styles.RenderHeader(sb.box.Line(header)))
	sb.writeln(sb.styles.RenderHeader(separator))
	for _, line := range body {
		sb.writeln(bodyStyle(sb.box.Line(line)))
	}
	sb.writeln(sb.styles.RenderHeader(sb.box.Bottom()))
}

// GetWidth returns the current box width.
func (sb *StyledBox) GetWidth() int {
	return sb.box.getWidth()
}
