// Package core provides the fundamental types and interfaces for terminal output.
//
// The core package contains the Printer interface that defines the contract
// for all output operations, along with the core data types used throughout
// the output system.
//
// Key types:
//   - [Printer] - Interface for terminal output operations
//   - [StepResult] - Result of a single workflow step
//   - [StoryResult] - Result of processing a story
//   - [ToolParams] - Parameters for tool invocations
package core

import (
	"encoding/json"
	"time"

	"bmaduum/internal/claude"
)

// StepResult represents the result of a single workflow step execution.
type StepResult struct {
	Name     string
	Duration time.Duration
	Success  bool
}

// StoryResult represents the result of processing a story in queue or epic operations.
type StoryResult struct {
	Key      string
	Success  bool
	Duration time.Duration
	FailedAt string
	Skipped  bool
}

// ToolParams contains parameters for a tool invocation.
type ToolParams struct {
	Name        string
	Description string
	Command     string // Bash
	FilePath    string // Read, Write, Edit
	OldString   string // Edit - text to replace
	NewString   string // Edit - replacement text
	Pattern     string // Glob, Grep
	Query       string // WebSearch
	URL         string // WebFetch
	Path        string // Glob, Grep directory
	Content     string // Write - file content

	// Raw JSON for unknown tools
	InputRaw json.RawMessage

	// Task tool fields
	SubagentType string
	Prompt       string

	// NotebookEdit fields
	NotebookPath string
	CellID       string
	NewSource    string
	EditMode     string
	CellType     string

	// AskUserQuestion fields
	Questions []claude.Question

	// Skill tool fields
	Skill string
	Args  string

	// TodoWrite fields
	Todos []claude.TodoItem
}

// Printer defines the interface for structured terminal output operations.
//
// The Printer interface provides a comprehensive set of methods for all
// output operations in the CLI, including:
//   - Session lifecycle (SessionStart, SessionEnd)
//   - Step lifecycle (StepStart, StepEnd)
//   - Tool output (ToolUse, ToolResult)
//   - Text and formatting (Text, Divider)
//   - Cycle operations (CycleHeader, CycleSummary, CycleFailed)
//   - Queue operations (QueueHeader, QueueStoryStart, QueueSummary)
//   - Command operations (CommandHeader, CommandFooter)
type Printer interface {
	SessionStart()
	SessionEnd(duration time.Duration, success bool)
	StepStart(step, total int, name string)
	StepEnd(duration time.Duration, success bool)
	ToolUse(params ToolParams)
	ToolResult(stdout, stderr string, truncateLines int)
	Text(message string)
	Divider()
	CycleHeader(storyKey string)
	CycleSummary(storyKey string, steps []StepResult, totalDuration time.Duration)
	CycleFailed(storyKey string, failedStep string, duration time.Duration)
	QueueHeader(count int, stories []string)
	QueueStoryStart(index, total int, storyKey string)
	QueueSummary(results []StoryResult, allKeys []string, totalDuration time.Duration)
	CommandHeader(label, prompt string, truncateLength int)
	CommandFooter(duration time.Duration, success bool, exitCode int)
}
