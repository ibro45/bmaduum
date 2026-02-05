package workflow

import (
	"bmaduum/internal/claude"
	"bmaduum/internal/output/core"
)

// PendingTool represents a buffered tool use waiting for its result.
type PendingTool struct {
	ID     string          // Tool use ID for correlation
	Params core.ToolParams // Parameters to pass to ToolUse
}

// ToolCorrelator buffers tool uses and correlates them with their results.
// This enables printing tool+result pairs together rather than separately.
type ToolCorrelator struct {
	pending []PendingTool
}

// NewToolCorrelator creates a new tool correlator.
func NewToolCorrelator() *ToolCorrelator {
	return &ToolCorrelator{
		pending: make([]PendingTool, 0),
	}
}

// Reset clears all pending tools. Call this at the start of each execution.
func (c *ToolCorrelator) Reset() {
	c.pending = make([]PendingTool, 0)
}

// AddToolUse buffers a tool use event for later correlation with its result.
func (c *ToolCorrelator) AddToolUse(id string, params core.ToolParams) {
	c.pending = append(c.pending, PendingTool{
		ID:     id,
		Params: params,
	})
}

// MatchResult finds and removes the pending tool that matches the given result.
// If toolUseID is provided, matches by ID. Otherwise uses FIFO ordering.
// Returns the matched tool params and true if found, or empty params and false if not.
func (c *ToolCorrelator) MatchResult(toolUseID string) (core.ToolParams, bool) {
	if len(c.pending) == 0 {
		return core.ToolParams{}, false
	}

	// Try to match by ID first
	if toolUseID != "" {
		for i, tool := range c.pending {
			if tool.ID == toolUseID {
				// Remove this tool from pending
				c.pending = append(c.pending[:i], c.pending[i+1:]...)
				return tool.Params, true
			}
		}
	}

	// Fall back to FIFO ordering (pop first)
	tool := c.pending[0]
	c.pending = c.pending[1:]
	return tool.Params, true
}

// Flush returns all pending tools and clears the buffer.
// Call this when you need to print buffered tools without waiting for results
// (e.g., when text arrives or session ends).
func (c *ToolCorrelator) Flush() []PendingTool {
	tools := c.pending
	c.pending = make([]PendingTool, 0)
	return tools
}

// HasPending returns true if there are buffered tools waiting for results.
func (c *ToolCorrelator) HasPending() bool {
	return len(c.pending) > 0
}

// EventToToolParams converts a tool use event to core.ToolParams.
func EventToToolParams(event claude.Event) core.ToolParams {
	return core.ToolParams{
		Name:        event.ToolName,
		Description: event.ToolDescription,
		Command:     event.ToolCommand,
		FilePath:    event.ToolFilePath,
		OldString:   event.ToolOldString,
		NewString:   event.ToolNewString,
		Pattern:     event.ToolPattern,
		Query:       event.ToolQuery,
		URL:         event.ToolURL,
		Path:        event.ToolPath,
		Content:     event.ToolContent,
		InputRaw:    event.ToolInputRaw,
		// Task tool fields
		SubagentType: event.ToolSubagentType,
		Prompt:       event.ToolPrompt,
		// NotebookEdit fields
		NotebookPath: event.ToolNotebookPath,
		CellID:       event.ToolCellID,
		NewSource:    event.ToolNewSource,
		EditMode:     event.ToolEditMode,
		CellType:     event.ToolCellType,
		// AskUserQuestion fields
		Questions: event.ToolQuestions,
		// Skill tool fields
		Skill: event.ToolSkill,
		Args:  event.ToolArgs,
		// TodoWrite fields
		Todos: event.ToolTodos,
	}
}
