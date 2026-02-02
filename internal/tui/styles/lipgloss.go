package styles

import (
	"github.com/charmbracelet/lipgloss"
)

// Styles provides pre-defined lipgloss styles for the TUI.
type Styles struct {
	// Header styles
	Header       lipgloss.Style
	HeaderText   lipgloss.Style
	HeaderMuted  lipgloss.Style
	HeaderActive lipgloss.Style

	// Content styles
	Text       lipgloss.Style
	TextMuted  lipgloss.Style
	TextDim    lipgloss.Style
	ToolUse    lipgloss.Style
	ToolOutput lipgloss.Style

	// Symbol styles
	ToolIcon   lipgloss.Style
	OutputIcon lipgloss.Style
	Success    lipgloss.Style
	Error      lipgloss.Style
	Warning    lipgloss.Style
	Info       lipgloss.Style

	// Section styles
	SectionDivider lipgloss.Style
}

// DefaultStyles returns the default Claude Code-inspired styles.
func DefaultStyles() Styles {
	return Styles{
		// Header with purple background
		Header: lipgloss.NewStyle().
			Background(lipgloss.Color(ClaudeColors.Primary)).
			Foreground(lipgloss.Color("#FFFFFF")).
			Bold(true).
			Padding(0, 1),

		HeaderText: lipgloss.NewStyle().
			Background(lipgloss.Color(ClaudeColors.Primary)).
			Foreground(lipgloss.Color("#FFFFFF")),

		HeaderMuted: lipgloss.NewStyle().
			Background(lipgloss.Color(ClaudeColors.Primary)).
			Foreground(lipgloss.Color("#CCCCCC")),

		HeaderActive: lipgloss.NewStyle().
			Background(lipgloss.Color("#FFFFFF")).
			Foreground(lipgloss.Color(ClaudeColors.Primary)).
			Bold(true).
			Padding(0, 1),

		// Content
		Text: lipgloss.NewStyle().
			Foreground(lipgloss.Color(ClaudeColors.Text)),

		TextMuted: lipgloss.NewStyle().
			Foreground(lipgloss.Color(ClaudeColors.TextMuted)),

		TextDim: lipgloss.NewStyle().
			Foreground(lipgloss.Color(ClaudeColors.TextDim)),

		// Tool use display
		ToolUse: lipgloss.NewStyle().
			Foreground(lipgloss.Color(ClaudeColors.Text)),

		ToolOutput: lipgloss.NewStyle().
			Foreground(lipgloss.Color(ClaudeColors.TextMuted)),

		// Symbols
		ToolIcon: lipgloss.NewStyle().
			Foreground(lipgloss.Color(ClaudeColors.ToolIcon)).
			Bold(true),

		OutputIcon: lipgloss.NewStyle().
			Foreground(lipgloss.Color(ClaudeColors.OutputIcon)),

		Success: lipgloss.NewStyle().
			Foreground(lipgloss.Color(ClaudeColors.Success)).
			Bold(true),

		Error: lipgloss.NewStyle().
			Foreground(lipgloss.Color(ClaudeColors.Error)).
			Bold(true),

		Warning: lipgloss.NewStyle().
			Foreground(lipgloss.Color(ClaudeColors.Warning)),

		Info: lipgloss.NewStyle().
			Foreground(lipgloss.Color(ClaudeColors.Info)),

		// Section divider
		SectionDivider: lipgloss.NewStyle().
			Foreground(lipgloss.Color(ClaudeColors.Border)).
			MarginTop(1).
			MarginBottom(1),
	}
}

// Default returns the default styles instance for convenience.
var Default = DefaultStyles()
