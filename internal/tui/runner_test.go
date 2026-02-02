package tui

import (
	"testing"

	"bmaduum/internal/claude"
	"bmaduum/internal/config"

	"github.com/stretchr/testify/assert"
)

func TestNewRunner(t *testing.T) {
	cfg := &config.Config{}
	executor := &claude.MockExecutor{}

	runner := NewRunner(executor, cfg)

	assert.NotNil(t, runner)
	assert.Equal(t, executor, runner.executor)
	assert.Equal(t, cfg, runner.config)
}

func TestStepInfo(t *testing.T) {
	step := StepInfo{
		Name:       "dev-story",
		StoryKey:   "PROJ-123",
		NextStatus: "review",
	}

	assert.Equal(t, "dev-story", step.Name)
	assert.Equal(t, "PROJ-123", step.StoryKey)
	assert.Equal(t, "review", step.NextStatus)
}

func TestClaudeEventAdapter(t *testing.T) {
	tests := []struct {
		name     string
		event    claude.Event
		expected interface{}
	}{
		{
			name: "session start",
			event: claude.Event{
				SessionStarted: true,
			},
			expected: SessionStartEvent{},
		},
		{
			name: "text content",
			event: claude.Event{
				Type: claude.EventTypeAssistant,
				Text: "Hello world",
				Raw: &claude.StreamEvent{
					Type: "assistant",
					Message: &claude.MessageContent{
						Content: []claude.ContentBlock{
							{Type: "text", Text: "Hello world"},
						},
					},
				},
			},
			expected: TextEvent{Text: "Hello world"},
		},
		{
			name: "tool use",
			event: claude.Event{
				Type:        claude.EventTypeAssistant,
				ToolName:    "Bash",
				ToolCommand: "ls -la",
				Raw: &claude.StreamEvent{
					Type: "assistant",
					Message: &claude.MessageContent{
						Content: []claude.ContentBlock{
							{Type: "tool_use", Name: "Bash"},
						},
					},
				},
			},
			expected: ToolUseEvent{
				Name:    "Bash",
				Command: "ls -la",
			},
		},
		{
			name: "tool result",
			event: claude.Event{
				Type:       claude.EventTypeUser,
				ToolStdout: "output",
				ToolStderr: "error",
				Raw: &claude.StreamEvent{
					Type: "user",
					ToolUseResult: &claude.ToolResult{
						Stdout: "output",
						Stderr: "error",
					},
				},
			},
			expected: ToolResultEvent{
				Stdout: "output",
				Stderr: "error",
			},
		},
		{
			name: "session complete",
			event: claude.Event{
				Type:            claude.EventTypeResult,
				SessionComplete: true,
			},
			expected: SessionCompleteEvent{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ClaudeEventAdapter(tt.event)
			assert.IsType(t, tt.expected, result)
		})
	}
}

func TestClaudeEventAdapter_Unknown(t *testing.T) {
	// Event with no recognized fields should return nil
	event := claude.Event{}
	result := ClaudeEventAdapter(event)
	assert.Nil(t, result)
}
