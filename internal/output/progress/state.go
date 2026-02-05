// Package progress provides terminal progress display functionality.
//
// This package manages a fixed two-line status area at the bottom of the terminal
// with scrolling output above. The status area consists of:
//   - Activity line (second-to-last): Spinner + activity verb + timer + token count
//   - Status bar (last row): Operation + step info + story key + model + total timer
package progress

import "time"

// State holds all the current progress state.
// This struct is designed to be passed atomically for updates.
type State struct {
	// Step tracking
	Step     int
	Total    int
	StepName string
	StoryKey string
	Model    string

	// Operation context (e.g., "Epic 6", "Story 2/3")
	Operation string

	// Current activity
	CurrentTool string

	// Timing
	StartTime     time.Time
	StepStartTime time.Time
	ActivityStart time.Time

	// Thinking time tracking (time before first response)
	ThinkingStart    time.Time
	ThinkingDuration time.Duration
	HadFirstResponse bool

	// Token counts
	InputTokens  int
	OutputTokens int

	// Tool tracking
	ToolCount int // Number of tools used in this session

	// Animation state
	SpinnerIdx int
	VerbIdx    int
}

// ActivityState holds just the activity-related state for rendering.
type ActivityState struct {
	VerbIdx          int
	ActivityStart    time.Time
	CurrentTool      string
	InputTokens      int
	OutputTokens     int
	HadFirstResponse bool
	ThinkingDuration time.Duration
	ToolCount        int
}

// StatusState holds just the status bar state for rendering.
type StatusState struct {
	Step          int
	Total         int
	StepName      string
	StoryKey      string
	Model         string
	Operation     string
	StepStartTime time.Time
	StartTime     time.Time
	Width         int
}

// Result holds completion result data.
type Result struct {
	Success  bool
	Duration time.Duration
}
