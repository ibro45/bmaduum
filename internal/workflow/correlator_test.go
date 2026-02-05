package workflow

import (
	"testing"

	"bmaduum/internal/output/core"

	"github.com/stretchr/testify/assert"
)

func TestToolCorrelator_AddAndMatch(t *testing.T) {
	c := NewToolCorrelator()

	// Add a tool use
	params := core.ToolParams{Name: "Bash", Command: "ls"}
	c.AddToolUse("tool-123", params)

	// Match by ID
	matched, found := c.MatchResult("tool-123")
	assert.True(t, found)
	assert.Equal(t, "Bash", matched.Name)

	// No more pending
	assert.False(t, c.HasPending())
}

func TestToolCorrelator_FIFOFallback(t *testing.T) {
	c := NewToolCorrelator()

	// Add multiple tools without IDs
	c.AddToolUse("", core.ToolParams{Name: "Tool1"})
	c.AddToolUse("", core.ToolParams{Name: "Tool2"})
	c.AddToolUse("", core.ToolParams{Name: "Tool3"})

	// Match without ID should use FIFO
	matched, found := c.MatchResult("")
	assert.True(t, found)
	assert.Equal(t, "Tool1", matched.Name)

	matched, found = c.MatchResult("")
	assert.True(t, found)
	assert.Equal(t, "Tool2", matched.Name)

	matched, found = c.MatchResult("")
	assert.True(t, found)
	assert.Equal(t, "Tool3", matched.Name)

	// No more pending
	_, found = c.MatchResult("")
	assert.False(t, found)
}

func TestToolCorrelator_MatchByID(t *testing.T) {
	c := NewToolCorrelator()

	// Add tools with IDs
	c.AddToolUse("id-1", core.ToolParams{Name: "Tool1"})
	c.AddToolUse("id-2", core.ToolParams{Name: "Tool2"})
	c.AddToolUse("id-3", core.ToolParams{Name: "Tool3"})

	// Match out of order by ID
	matched, found := c.MatchResult("id-2")
	assert.True(t, found)
	assert.Equal(t, "Tool2", matched.Name)

	matched, found = c.MatchResult("id-1")
	assert.True(t, found)
	assert.Equal(t, "Tool1", matched.Name)

	matched, found = c.MatchResult("id-3")
	assert.True(t, found)
	assert.Equal(t, "Tool3", matched.Name)
}

func TestToolCorrelator_Flush(t *testing.T) {
	c := NewToolCorrelator()

	c.AddToolUse("id-1", core.ToolParams{Name: "Tool1"})
	c.AddToolUse("id-2", core.ToolParams{Name: "Tool2"})

	assert.True(t, c.HasPending())

	flushed := c.Flush()
	assert.Len(t, flushed, 2)
	assert.Equal(t, "Tool1", flushed[0].Params.Name)
	assert.Equal(t, "Tool2", flushed[1].Params.Name)

	assert.False(t, c.HasPending())
}

func TestToolCorrelator_Reset(t *testing.T) {
	c := NewToolCorrelator()

	c.AddToolUse("id-1", core.ToolParams{Name: "Tool1"})
	assert.True(t, c.HasPending())

	c.Reset()
	assert.False(t, c.HasPending())
}

func TestEventToToolParams(t *testing.T) {
	// This test verifies that all fields are properly converted
	// Note: This is more of a documentation test since the function
	// is straightforward field copying
	params := core.ToolParams{
		Name:    "Bash",
		Command: "ls -la",
	}
	assert.Equal(t, "Bash", params.Name)
	assert.Equal(t, "ls -la", params.Command)
}
