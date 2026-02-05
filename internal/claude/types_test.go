package claude

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEventFromStream_SystemInit(t *testing.T) {
	raw := &StreamEvent{
		Type:    "system",
		Subtype: "init",
	}

	event := NewEventFromStream(raw)

	assert.Equal(t, EventTypeSystem, event.Type)
	assert.Equal(t, "init", event.Subtype)
	assert.True(t, event.SessionStarted)
	assert.False(t, event.SessionComplete)
}

func TestNewEventFromStream_AssistantText(t *testing.T) {
	raw := &StreamEvent{
		Type: "assistant",
		Message: &MessageContent{
			Content: []ContentBlock{
				{
					Type: "text",
					Text: "Hello, I'm Claude!",
				},
			},
		},
	}

	event := NewEventFromStream(raw)

	assert.Equal(t, EventTypeAssistant, event.Type)
	assert.Equal(t, "Hello, I'm Claude!", event.Text)
	assert.True(t, event.IsText())
	assert.False(t, event.IsToolUse())
}

func TestNewEventFromStream_AssistantToolUse(t *testing.T) {
	raw := &StreamEvent{
		Type: "assistant",
		Message: &MessageContent{
			Content: []ContentBlock{
				{
					Type: "tool_use",
					Name: "Bash",
					Input: &ToolInput{
						Command:     "ls -la",
						Description: "List files",
					},
				},
			},
		},
	}

	event := NewEventFromStream(raw)

	assert.Equal(t, EventTypeAssistant, event.Type)
	assert.Equal(t, "Bash", event.ToolName)
	assert.Equal(t, "ls -la", event.ToolCommand)
	assert.Equal(t, "List files", event.ToolDescription)
	assert.True(t, event.IsToolUse())
	assert.False(t, event.IsText())
}

func TestNewEventFromStream_ToolResult(t *testing.T) {
	raw := &StreamEvent{
		Type: "user",
		ToolUseResult: &ToolResult{
			Stdout: "file1.go\nfile2.go",
			Stderr: "",
		},
	}

	event := NewEventFromStream(raw)

	assert.Equal(t, EventTypeUser, event.Type)
	assert.Equal(t, "file1.go\nfile2.go", event.ToolStdout)
	assert.True(t, event.IsToolResult())
}

func TestNewEventFromStream_Result(t *testing.T) {
	raw := &StreamEvent{
		Type: "result",
	}

	event := NewEventFromStream(raw)

	assert.Equal(t, EventTypeResult, event.Type)
	assert.True(t, event.SessionComplete)
}

func TestEvent_IsText(t *testing.T) {
	tests := []struct {
		name     string
		event    Event
		expected bool
	}{
		{
			name:     "text event",
			event:    Event{Type: EventTypeAssistant, Text: "hello"},
			expected: true,
		},
		{
			name:     "empty text",
			event:    Event{Type: EventTypeAssistant, Text: ""},
			expected: false,
		},
		{
			name:     "wrong type",
			event:    Event{Type: EventTypeSystem, Text: "hello"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.event.IsText())
		})
	}
}

func TestEvent_IsToolUse(t *testing.T) {
	tests := []struct {
		name     string
		event    Event
		expected bool
	}{
		{
			name:     "tool use event",
			event:    Event{Type: EventTypeAssistant, ToolName: "Bash"},
			expected: true,
		},
		{
			name:     "empty tool name",
			event:    Event{Type: EventTypeAssistant, ToolName: ""},
			expected: false,
		},
		{
			name:     "wrong type",
			event:    Event{Type: EventTypeSystem, ToolName: "Bash"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.event.IsToolUse())
		})
	}
}

func TestEvent_IsToolResult(t *testing.T) {
	tests := []struct {
		name     string
		event    Event
		expected bool
	}{
		{
			name:     "tool result with stdout",
			event:    Event{Type: EventTypeUser, ToolStdout: "output", HasToolResult: true},
			expected: true,
		},
		{
			name:     "tool result with stderr",
			event:    Event{Type: EventTypeUser, ToolStderr: "error", HasToolResult: true},
			expected: true,
		},
		{
			name:     "tool result with empty content",
			event:    Event{Type: EventTypeUser, HasToolResult: true},
			expected: true,
		},
		{
			name:     "user event without tool result",
			event:    Event{Type: EventTypeUser},
			expected: false,
		},
		{
			name:     "wrong type",
			event:    Event{Type: EventTypeAssistant, ToolStdout: "output"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.event.IsToolResult())
		})
	}
}

func TestContentBlock_UnmarshalJSON_RawInput(t *testing.T) {
	// Test that raw JSON is captured for unknown tools
	jsonData := `{
		"type": "tool_use",
		"name": "UnknownTool",
		"input": {
			"foo": "bar",
			"count": 42,
			"nested": {"key": "value"}
		}
	}`

	var block ContentBlock
	err := json.Unmarshal([]byte(jsonData), &block)
	require.NoError(t, err)

	assert.Equal(t, "tool_use", block.Type)
	assert.Equal(t, "UnknownTool", block.Name)
	assert.NotNil(t, block.InputRaw)

	// Verify raw JSON can be unmarshaled
	var rawParams map[string]interface{}
	err = json.Unmarshal(block.InputRaw, &rawParams)
	require.NoError(t, err)
	assert.Equal(t, "bar", rawParams["foo"])
	assert.Equal(t, float64(42), rawParams["count"])
}

func TestContentBlock_UnmarshalJSON_KnownFields(t *testing.T) {
	// Test that known fields are still parsed into ToolInput
	jsonData := `{
		"type": "tool_use",
		"name": "Bash",
		"input": {
			"command": "ls -la",
			"description": "List files"
		}
	}`

	var block ContentBlock
	err := json.Unmarshal([]byte(jsonData), &block)
	require.NoError(t, err)

	assert.Equal(t, "Bash", block.Name)
	assert.NotNil(t, block.Input)
	assert.Equal(t, "ls -la", block.Input.Command)
	assert.Equal(t, "List files", block.Input.Description)
}

func TestNewEventFromStream_TaskTool(t *testing.T) {
	raw := &StreamEvent{
		Type: "assistant",
		Message: &MessageContent{
			Content: []ContentBlock{
				{
					Type: "tool_use",
					Name: "Task",
					Input: &ToolInput{
						SubagentType: "Explore",
						Prompt:       "Find all configuration files",
					},
				},
			},
		},
	}

	event := NewEventFromStream(raw)

	assert.Equal(t, "Task", event.ToolName)
	assert.Equal(t, "Explore", event.ToolSubagentType)
	assert.Equal(t, "Find all configuration files", event.ToolPrompt)
}

func TestNewEventFromStream_WriteTool(t *testing.T) {
	raw := &StreamEvent{
		Type: "assistant",
		Message: &MessageContent{
			Content: []ContentBlock{
				{
					Type: "tool_use",
					Name: "Write",
					Input: &ToolInput{
						FilePath: "/path/to/file.go",
						Content:  "package main\n\nfunc main() {}",
					},
				},
			},
		},
	}

	event := NewEventFromStream(raw)

	assert.Equal(t, "Write", event.ToolName)
	assert.Equal(t, "/path/to/file.go", event.ToolFilePath)
	assert.Equal(t, "package main\n\nfunc main() {}", event.ToolContent)
}

func TestNewEventFromStream_NotebookEditTool(t *testing.T) {
	raw := &StreamEvent{
		Type: "assistant",
		Message: &MessageContent{
			Content: []ContentBlock{
				{
					Type: "tool_use",
					Name: "NotebookEdit",
					Input: &ToolInput{
						NotebookPath: "/path/to/notebook.ipynb",
						CellID:       "cell-123",
						NewSource:    "print('hello')",
						EditMode:     "replace",
						CellType:     "code",
					},
				},
			},
		},
	}

	event := NewEventFromStream(raw)

	assert.Equal(t, "NotebookEdit", event.ToolName)
	assert.Equal(t, "/path/to/notebook.ipynb", event.ToolNotebookPath)
	assert.Equal(t, "cell-123", event.ToolCellID)
	assert.Equal(t, "print('hello')", event.ToolNewSource)
	assert.Equal(t, "replace", event.ToolEditMode)
	assert.Equal(t, "code", event.ToolCellType)
}

func TestNewEventFromStream_AskUserQuestionTool(t *testing.T) {
	raw := &StreamEvent{
		Type: "assistant",
		Message: &MessageContent{
			Content: []ContentBlock{
				{
					Type: "tool_use",
					Name: "AskUserQuestion",
					Input: &ToolInput{
						Questions: []Question{
							{
								Question: "Which library?",
								Header:   "Library",
								Options: []QuestionOption{
									{Label: "React", Description: "Facebook's library"},
									{Label: "Vue", Description: "Progressive framework"},
								},
							},
						},
					},
				},
			},
		},
	}

	event := NewEventFromStream(raw)

	assert.Equal(t, "AskUserQuestion", event.ToolName)
	require.Len(t, event.ToolQuestions, 1)
	assert.Equal(t, "Which library?", event.ToolQuestions[0].Question)
	assert.Equal(t, "Library", event.ToolQuestions[0].Header)
	require.Len(t, event.ToolQuestions[0].Options, 2)
	assert.Equal(t, "React", event.ToolQuestions[0].Options[0].Label)
}

func TestNewEventFromStream_SkillTool(t *testing.T) {
	raw := &StreamEvent{
		Type: "assistant",
		Message: &MessageContent{
			Content: []ContentBlock{
				{
					Type: "tool_use",
					Name: "Skill",
					Input: &ToolInput{
						Skill: "commit",
						Args:  "-m 'Fix bug'",
					},
				},
			},
		},
	}

	event := NewEventFromStream(raw)

	assert.Equal(t, "Skill", event.ToolName)
	assert.Equal(t, "commit", event.ToolSkill)
	assert.Equal(t, "-m 'Fix bug'", event.ToolArgs)
}
