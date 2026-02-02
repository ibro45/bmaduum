package events

import (
	"time"

	"bmaduum/internal/claude"
)

// Message types for Bubble Tea TUI updates.

// ClaudeEventMsg wraps a Claude event for TUI consumption.
type ClaudeEventMsg struct {
	Event claude.Event
}

// TextContentMsg indicates new text content has arrived.
type TextContentMsg struct {
	Text string
}

// ToolUseMsg indicates a tool is being invoked.
type ToolUseMsg struct {
	Name        string
	Description string
	Command     string
	FilePath    string
}

// ToolResultMsg contains the result of a tool execution.
type ToolResultMsg struct {
	Stdout string
	Stderr string
}

// StepStartMsg indicates a new workflow step has started.
type StepStartMsg struct {
	Step      int
	Total     int
	StepName  string
	StoryKey  string
}

// StepCompleteMsg indicates the current step has completed.
type StepCompleteMsg struct {
	Success  bool
	Duration time.Duration
}

// SessionStartMsg indicates a Claude session has started.
type SessionStartMsg struct {
	SessionID string
}

// SessionCompleteMsg indicates a Claude session has completed.
type SessionCompleteMsg struct {
	Success bool
}

// CompleteMsg indicates the entire TUI workflow is complete.
type CompleteMsg struct {
	ExitCode int
}

// ThinkingMsg toggles the thinking spinner state.
type ThinkingMsg struct {
	Thinking bool
	Status   string
}

// ResizeMsg indicates the terminal has been resized.
type ResizeMsg struct {
	Width  int
	Height int
}

// TypewriterTickMsg is sent for each character animation frame.
type TypewriterTickMsg struct {
	Time time.Time
}

// SpinnerTickMsg is sent for each spinner animation frame.
type SpinnerTickMsg struct {
	Time time.Time
}
