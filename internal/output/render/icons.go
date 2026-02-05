// Package render provides output rendering components for the CLI.
package render

// Icons for status indicators in terminal output.
// These are the single source of truth for all icon constants.
const (
	// Status icons
	IconSuccess    = "✓" // Completed successfully
	IconError      = "✗" // Failed
	IconPending    = "○" // Not yet started
	IconInProgress = "●" // Currently running

	// Tool icons
	IconTool   = "⏺" // Tool invocation (filled circle)
	IconOutput = "⎿" // Tool output (right angle bracket)

	// Brand icons
	IconBmaduum = "⚡" // Bmaduum logo (high voltage)

	// Utility icons
	IconClock     = "⏱" // Clock/timer (no variation selector for consistency)
	IconPaused    = "⏸" // Paused rate limit
	IconRetrying  = "↻" // Retrying (simpler arrow)
	IconSeparator = "│" // Vertical separator
)

// Indentation constants matching Claude Code terminal output format.
// Claude Code uses 2-space base indentation for all output.
const (
	// IndentToolUse is the leading indent for tool use lines (⏺)
	// Format: "  ⏺ ToolName(params)"
	IndentToolUse = "  "

	// IndentToolResult is the leading indent for tool result lines (⎿)
	// Format: "    ⎿  content" (4 spaces before bracket, 2 after)
	IndentToolResult = "    "

	// SpaceAfterBracket is the spacing after the ⎿ bracket
	// Claude Code uses 2 spaces after the bracket before content
	SpaceAfterBracket = "  "
)
