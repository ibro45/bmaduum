// Package output provides terminal output formatting using lipgloss styles.
//
// The package provides structured output for CLI operations including session
// lifecycle, step progress, tool usage display, and batch operation summaries.
// All output is styled using the lipgloss library for consistent terminal rendering.
//
// Key types:
//   - [Printer] - Interface for structured terminal output operations
//   - [DefaultPrinter] - Production implementation using lipgloss styles
//   - [StepResult] - Result of a single workflow step execution
//   - [StoryResult] - Result of processing a story in queue/epic operations
//
// Use [NewPrinter] for production output to stdout, or [NewPrinterWithWriter]
// to capture output in tests by providing a custom io.Writer.
package output

import (
	"github.com/charmbracelet/lipgloss"
)

// Claude Code color palette (hex values).
// These colors are inspired by Claude's visual style.
var (
	// Brand colors - Claude's terracotta warmth
	colorBrand      = lipgloss.Color("#C15F3C") // Claude terracotta - primary brand
	colorBrandLight = lipgloss.Color("#D97B5C") // Lighter terracotta
	colorBrandDark  = lipgloss.Color("#A14E30") // Darker terracotta

	// Accent colors
	colorAccent    = lipgloss.Color("#6B4EE6") // Purple accent
	colorAccentDim = lipgloss.Color("#5A3FD4") // Darker purple

	// Functional colors
	colorTool      = lipgloss.Color("#58A6FF") // Blue - tool invocations
	colorOutput    = lipgloss.Color("#8B949E") // Gray - tool output
	colorText      = lipgloss.Color("#E6EDF3") // Off-white - text content
	colorTextMuted = lipgloss.Color("#8B949E") // Gray - secondary text
	colorTextDim   = lipgloss.Color("#6E7681") // Dark gray - subtle text
	colorSuccess   = lipgloss.Color("#3FB950") // Green - success
	colorError     = lipgloss.Color("#F85149") // Red - errors
	colorWarning   = lipgloss.Color("#D29922") // Orange - warnings
	colorInfo      = lipgloss.Color("#58A6FF") // Blue - info
	colorBorder    = lipgloss.Color("#30363D") // Dark border
	colorMuted     = lipgloss.Color("#6E7681") // Gray - secondary info

	// Status bar colors
	colorStatusBg  = lipgloss.Color("#C15F3C") // Terracotta background
	colorStatusFg  = lipgloss.Color("#FFFFFF") // White text on status bar
	colorStatusDim = lipgloss.Color("#F4F3EE") // Cream for secondary text

	// Syntax highlighting colors
	colorComment  = lipgloss.Color("#8B949E") // Gray
	colorKeyword  = lipgloss.Color("#FF7B72") // Red/pink
	colorString   = lipgloss.Color("#A5D6FF") // Light blue
	colorFunction = lipgloss.Color("#D2A8FF") // Purple
	colorNumber   = lipgloss.Color("#79C0FF") // Blue
	colorDiffAdd  = lipgloss.Color("#3FB950") // Green - added lines
	colorDiffDel  = lipgloss.Color("#F85149") // Red - deleted lines
)

// Adaptive colors for light/dark terminal detection.
var (
	adaptiveBrand = lipgloss.AdaptiveColor{
		Light: "#A14E30", // Darker on light backgrounds
		Dark:  "#C15F3C", // Standard on dark backgrounds
	}
	adaptiveText = lipgloss.AdaptiveColor{
		Light: "#24292F", // Dark text on light
		Dark:  "#E6EDF3", // Light text on dark
	}
	adaptiveSuccess = lipgloss.AdaptiveColor{
		Light: "#1A7F37", // Darker green on light
		Dark:  "#3FB950", // Bright green on dark
	}
	adaptiveError = lipgloss.AdaptiveColor{
		Light: "#CF222E", // Darker red on light
		Dark:  "#F85149", // Bright red on dark
	}
)

// Lipgloss styles for different output elements.
// These styles are used by [DefaultPrinter] methods.
var (
	// headerStyle formats major section headers.
	headerStyle = conditionalStyle(lipgloss.NewStyle().
			Bold(true).
			Foreground(colorBrand))

	// stepHeaderStyle formats step progress headers.
	stepHeaderStyle = conditionalStyle(lipgloss.NewStyle().
			Bold(true).
			Foreground(colorBrand).
			Margin(1, 0))

	// successStyle formats success indicators (checkmarks, completion messages).
	successStyle = conditionalStyle(lipgloss.NewStyle().
			Bold(true).
			Foreground(colorSuccess))

	// errorStyle formats error indicators (X marks, failure messages).
	errorStyle = conditionalStyle(lipgloss.NewStyle().
			Bold(true).
			Foreground(colorError))

	// mutedStyle formats secondary information (stderr, pending items).
	mutedStyle = conditionalStyle(lipgloss.NewStyle().
			Foreground(colorTextMuted))

	// labelStyle formats command labels and highlights.
	labelStyle = conditionalStyle(lipgloss.NewStyle().
			Bold(true).
			Foreground(colorBrand))

	// bulletStyle formats the bullet icon for tool use and assistant text (green).
	bulletStyle = conditionalStyle(lipgloss.NewStyle().
			Foreground(colorSuccess))

	// toolNameStyle formats the tool name (bold white).
	toolNameStyle = conditionalStyle(lipgloss.NewStyle().
			Bold(true))

	// toolParamsStyle formats tool parameters (muted).
	toolParamsStyle = conditionalStyle(lipgloss.NewStyle().
			Foreground(colorTextMuted))

	// toolUseStyle formats tool invocation lines (legacy, kept for compatibility).
	toolUseStyle = conditionalStyle(lipgloss.NewStyle().
			Foreground(colorTool).
			Bold(true))

	// toolOutputStyle formats tool output lines (muted gray).
	toolOutputStyle = conditionalStyle(lipgloss.NewStyle().
			Foreground(colorOutput))

	// textStyle formats plain text content (Claude's output).
	textStyle = conditionalStyle(lipgloss.NewStyle().
			Foreground(colorText))

	// dividerStyle formats visual separator lines.
	dividerStyle = conditionalStyle(lipgloss.NewStyle().
			Foreground(colorBorder))

	// summaryStyle formats completion summary boxes.
	summaryStyle = conditionalStyle(lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(colorBrand).
			Padding(0, 1))

	// queueHeaderStyle formats queue operation headers.
	queueHeaderStyle = conditionalStyle(lipgloss.NewStyle().
				Bold(true).
				Foreground(colorBrand).
				Margin(1, 0))

	// Status bar style - reverse video with brand color
	statusBarStyle = conditionalStyle(lipgloss.NewStyle().
			Background(colorStatusBg).
			Foreground(colorStatusFg).
			Bold(true))

	// Status bar dim text (for separators, secondary info)
	statusBarDimStyle = conditionalStyle(lipgloss.NewStyle().
				Background(colorStatusBg).
				Foreground(colorStatusDim))

	// Styles for basic diff highlighting (used by FormatDiff fallback)
	diffAddStyle = conditionalStyle(lipgloss.NewStyle().
			Foreground(colorDiffAdd))

	diffDelStyle = conditionalStyle(lipgloss.NewStyle().
			Foreground(colorDiffDel))

	diffMetaStyle = conditionalStyle(lipgloss.NewStyle().
			Foreground(colorTextMuted))

	diffHeaderStyle = conditionalStyle(lipgloss.NewStyle().
			Foreground(colorInfo).
			Bold(true))

	// Rich diff rendering styles (used by DiffRenderer)

	// Added line content - dark green background with bright green text
	diffAddedLineStyle = conditionalStyle(lipgloss.NewStyle().
				Background(lipgloss.Color("#0d2818")).
				Foreground(lipgloss.Color("#3fb950")))

	// Deleted line content - dark red background with bright red text
	diffDeletedLineStyle = conditionalStyle(lipgloss.NewStyle().
				Background(lipgloss.Color("#2d0a0a")).
				Foreground(lipgloss.Color("#f85149")))

	// Context line content - no background, muted text
	diffContextLineStyle = conditionalStyle(lipgloss.NewStyle().
				Foreground(colorTextMuted))

	// Line number for context lines - dim
	diffLineNumStyle = conditionalStyle(lipgloss.NewStyle().
				Foreground(colorTextDim))

	// Line number for added lines - green tint
	diffLineNumAddedStyle = conditionalStyle(lipgloss.NewStyle().
				Foreground(lipgloss.Color("#238636")))

	// Line number for deleted lines - red tint
	diffLineNumDeletedStyle = conditionalStyle(lipgloss.NewStyle().
				Foreground(lipgloss.Color("#b62324")))

	// Gutter marker (┃) for changed sections - brand color
	diffGutterMarkerStyle = conditionalStyle(lipgloss.NewStyle().
				Foreground(colorBrand).
				Bold(true))

	// + marker for added lines
	diffAddedMarkerStyle = conditionalStyle(lipgloss.NewStyle().
				Foreground(colorDiffAdd).
				Bold(true))

	// - marker for deleted lines
	diffDeletedMarkerStyle = conditionalStyle(lipgloss.NewStyle().
				Foreground(colorDiffDel).
				Bold(true))

	// Summary line ("Added X lines, removed Y lines")
	diffSummaryStyle = conditionalStyle(lipgloss.NewStyle().
				Foreground(colorTextMuted).
				Italic(true))

	// Question header style for AskUserQuestion tool
	questionHeaderStyle = conditionalStyle(lipgloss.NewStyle().
				Foreground(colorInfo).
				Bold(true))

	// Styles for code syntax highlighting (used by syntax.go)
	codeKeywordStyle = conditionalStyle(lipgloss.NewStyle().
				Foreground(colorKeyword))

	codeStringStyle = conditionalStyle(lipgloss.NewStyle().
			Foreground(colorString))

	codeCommentStyle = conditionalStyle(lipgloss.NewStyle().
				Foreground(colorComment).
				Italic(true))

	codeFunctionStyle = conditionalStyle(lipgloss.NewStyle().
				Foreground(colorFunction))

	codeNumberStyle = conditionalStyle(lipgloss.NewStyle().
			Foreground(colorNumber))
)

// Icons for status indicators in terminal output.
// Re-exported from render package for convenience and backward compatibility.
// These are lowercase aliases to the exported Icon* constants in the render package.
const (
	iconSuccess    = "✓" // Completed successfully
	iconError      = "✗" // Failed
	iconPending    = "○" // Not yet started
	iconInProgress = "●" // Currently running
	iconTool       = "⏺" // Tool invocation (filled circle)
	iconOutput     = "⎿" // Tool output (right angle bracket)
	iconBmaduum    = "⚡" // Bmaduum logo (high voltage)
	iconClock      = "⏱" // Clock/timer (no variation selector for consistency)
	iconPaused     = "⏸" // Paused rate limit
	iconRetrying   = "↻" // Retrying (simpler arrow)
	iconSeparator  = "│" // Vertical separator
)

// Spinner frames for animated thinking indicator.
// Claude Code uses a simple static dot, no animation.
var spinnerFrames = []string{"·"}

// Thinking verbs for Claude Code-style status messages.
// Randomly selected during "thinking" states for whimsy.
var thinkingVerbs = []string{
	"Thinking",
	"Pondering",
	"Analyzing",
	"Processing",
	"Cogitating",
	"Contemplating",
	"Reasoning",
	"Evaluating",
	"Computing",
	"Deliberating",
	"Mulling",
	"Considering",
	"Musing",
	"Ruminating",
	"Synthesizing",
	"Clauding", // Easter egg
}

// Past-tense versions of thinking verbs for completion messages.
// Must match the order of thinkingVerbs for correct pairing.
var pastTenseVerbs = []string{
	"Thought",
	"Pondered",
	"Analyzed",
	"Processed",
	"Cogitated",
	"Contemplated",
	"Reasoned",
	"Evaluated",
	"Computed",
	"Deliberated",
	"Mulled",
	"Considered",
	"Mused",
	"Ruminated",
	"Synthesized",
	"Clauded", // Easter egg past tense
}

// conditionalStyle returns the style if colors are supported, or a plain style otherwise.
//
// This function ensures that styling is only applied when the terminal supports it.
// When output is piped to a file or the terminal doesn't support colors, returns
// an unstyled lipgloss.Style to avoid ANSI escape codes in the output.
func conditionalStyle(style lipgloss.Style) lipgloss.Style {
	if SupportsColor() {
		return style
	}
	return lipgloss.NewStyle()
}
