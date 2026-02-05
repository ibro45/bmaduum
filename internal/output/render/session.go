// Package render provides output rendering components for the CLI.
package render

import (
	"fmt"
	"io"
	"strings"
	"time"
)

// SessionStyleProvider provides styling functions for session rendering.
type SessionStyleProvider interface {
	RenderHeader(s string) string
	RenderSuccess(s string) string
	RenderError(s string) string
	RenderMuted(s string) string
	RenderDivider(s string) string
	RenderBullet(s string) string
	RenderText(s string) string
}

// SessionWidthProvider provides terminal width information.
type SessionWidthProvider interface {
	TerminalWidth() int
}

// SessionRenderer handles rendering for session, command, and text output.
type SessionRenderer struct {
	writer         io.Writer
	styles         SessionStyleProvider
	width          SessionWidthProvider
	box            *Box
	renderMarkdown func(message string) string
}

// NewSessionRenderer creates a new session renderer.
// The markdownRender function should convert markdown to styled terminal output.
func NewSessionRenderer(writer io.Writer, styles SessionStyleProvider, width SessionWidthProvider, markdownRender func(string) string) *SessionRenderer {
	box := NewBox(width, 100, 40)
	return &SessionRenderer{
		writer:         writer,
		styles:         styles,
		width:          width,
		box:            box,
		renderMarkdown: markdownRender,
	}
}

// writeln writes a formatted line to the output.
func (r *SessionRenderer) Writeln(format string, args ...interface{}) {
	fmt.Fprintf(r.writer, format+"\n", args...)
}

// SessionStart prints session start indicator.
func (r *SessionRenderer) SessionStart() {
	r.Writeln("%s Session started", IconInProgress)
}

// SessionEnd prints session end with status.
func (r *SessionRenderer) SessionEnd(duration time.Duration, success bool) {
	if success {
		r.Writeln("%s Session complete (%s)", r.styles.RenderSuccess(IconSuccess), duration.Round(time.Millisecond))
	} else {
		r.Writeln("%s Session failed (%s)", r.styles.RenderError(IconError), duration.Round(time.Millisecond))
	}
}

// CommandHeader prints a nice box with command information.
func (r *SessionRenderer) CommandHeader(label, prompt string, truncateLength int) {
	width := r.width.TerminalWidth()
	if width > 100 {
		width = 100
	}
	if width < 40 {
		width = 40
	}

	// Parse label to extract step info
	// Label format is typically "workflow-name: story-key" or just "workflow-name"
	parts := strings.SplitN(label, ": ", 2)
	workflowName := parts[0]
	storyKey := ""
	if len(parts) > 1 {
		storyKey = parts[1]
	}

	// Build the box (minimal spacing - output follows directly)
	r.Writeln(r.styles.RenderHeader(BoxTop(width)))

	// Title line
	title := IconBmaduum + " " + workflowName
	if storyKey != "" {
		title += " | " + storyKey
	}
	r.Writeln(r.styles.RenderHeader(BoxLine(title, width)))

	// Separator
	r.Writeln(r.styles.RenderHeader("├" + strings.Repeat("─", width-2) + "┤"))

	// Command (wrapped at word boundaries) - border stays header color, content is muted
	commandLines := BoxLineWrapWordsStyled("Command", prompt, width, r.styles.RenderHeader, r.styles.RenderMuted)
	for _, line := range commandLines {
		r.Writeln(line)
	}

	// Bottom - output follows directly after
	r.Writeln(r.styles.RenderHeader(BoxBottom(width)))
}

// CommandFooter prints the footer after a command completes.
func (r *SessionRenderer) CommandFooter(duration time.Duration, success bool, exitCode int) {
	if success {
		r.Writeln("%s Complete (%s)", r.styles.RenderSuccess(IconSuccess), duration.Round(time.Millisecond))
	} else {
		r.Writeln("%s Failed (exit %d, %s)", r.styles.RenderError(IconError), exitCode, duration.Round(time.Millisecond))
	}
}

// Text prints a text message from Claude.
// Format: "  ● text" with 2-space base indent and bullet, matching Claude Code style.
// Markdown is rendered with proper formatting (bold, code, headers, etc.)
func (r *SessionRenderer) Text(message string) {
	if message == "" {
		return
	}

	// Render markdown first (glamour handles word wrapping)
	rendered := r.renderMarkdown(message)

	bullet := r.styles.RenderBullet(IconTool)
	// Print Claude's text with bullet on first line
	// Uses IndentToolUse (2 spaces) for base indent matching Claude Code
	lines := strings.Split(rendered, "\n")
	for i, line := range lines {
		if line == "" {
			fmt.Fprintln(r.writer)
		} else if i == 0 {
			// First line gets the bullet with 2-space base indent
			fmt.Fprintf(r.writer, "%s%s %s\n", IndentToolUse, bullet, line)
		} else {
			// Continuation lines are indented to align with text after bullet (4 spaces)
			fmt.Fprintf(r.writer, "%s  %s\n", IndentToolUse, line)
		}
	}
}

// Divider prints a visual divider (thin line).
func (r *SessionRenderer) Divider() {
	width := r.width.TerminalWidth()
	if width > 80 {
		width = 80
	}
	r.Writeln(r.styles.RenderDivider(strings.Repeat("─", width)))
}
