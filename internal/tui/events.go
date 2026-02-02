package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// TextEvent indicates new text content has arrived.
type TextEvent struct{ Text string }

// ToolUseEvent indicates a tool is being invoked.
type ToolUseEvent struct {
	Name        string
	Description string
	Command     string
	FilePath    string
}

// ToolResultEvent contains the result of a tool execution.
type ToolResultEvent struct {
	Stdout string
	Stderr string
}

// SessionStartEvent indicates a Claude session has started.
type SessionStartEvent struct{}

// StepStartMsg indicates a new workflow step has started.
type StepStartMsg struct {
	Step     int
	Total    int
	StepName string
	StoryKey string
}

// StepCompleteMsg indicates the current step has completed.
type StepCompleteMsg struct {
	Success  bool
	Duration time.Duration
}

// CompleteMsg indicates the entire TUI workflow is complete.
type CompleteMsg struct {
	ExitCode int
}

// SpinnerUpdateMsg updates the spinner state.
type SpinnerUpdateMsg struct {
	tea.Msg
}

// SessionCompleteEvent indicates a Claude session has completed.
type SessionCompleteEvent struct{}
