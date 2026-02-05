package output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"testing"
	"time"

	"bmaduum/internal/claude"
	"bmaduum/internal/output/core"
	"bmaduum/internal/output/diff"

	"github.com/stretchr/testify/assert"
)

func TestNewPrinter(t *testing.T) {
	p := NewPrinter()
	assert.NotNil(t, p)
}

func TestNewPrinterWithWriter(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinterWithWriter(&buf)
	assert.NotNil(t, p)
}

func TestDefaultPrinter_SessionStart(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinterWithWriter(&buf)

	p.SessionStart()

	output := buf.String()
	assert.Contains(t, output, "Session started")
}

func TestDefaultPrinter_SessionEnd(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinterWithWriter(&buf)

	p.SessionEnd(5*time.Second, true)

	output := buf.String()
	assert.Contains(t, output, "Session complete")
}

func TestDefaultPrinter_StepStart(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinterWithWriter(&buf)

	p.StepStart(1, 4, "create-story")

	// StepStart now outputs nothing - step info is in CommandHeader and ProgressLine
	output := buf.String()
	assert.Empty(t, output)
}

func TestDefaultPrinter_ToolUse(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinterWithWriter(&buf)

	p.ToolUse(core.ToolParams{Name: "Bash", Description: "List files", Command: "ls -la"})

	output := buf.String()
	assert.Contains(t, output, "⏺")
	assert.Contains(t, output, "Bash")
	assert.Contains(t, output, "(ls -la)")
}

func TestDefaultPrinter_ToolUse_WithFilePath(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinterWithWriter(&buf)

	p.ToolUse(core.ToolParams{Name: "Read", FilePath: "/path/to/file.go"})

	output := buf.String()
	assert.Contains(t, output, "Read")
	assert.Contains(t, output, "/path/to/file.go")
}

func TestDefaultPrinter_ToolUse_WithPattern(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinterWithWriter(&buf)

	p.ToolUse(core.ToolParams{Name: "Glob", Pattern: "**/*.go"})

	output := buf.String()
	assert.Contains(t, output, "Glob")
	assert.Contains(t, output, "**/*.go")
}

func TestDefaultPrinter_ToolUse_WithPatternAndPath(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinterWithWriter(&buf)

	p.ToolUse(core.ToolParams{Name: "Grep", Pattern: "func main", Path: "/src"})

	output := buf.String()
	assert.Contains(t, output, "Grep")
	assert.Contains(t, output, "func main")
	assert.Contains(t, output, "/src")
}

func TestDefaultPrinter_ToolResult(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinterWithWriter(&buf)

	p.ToolResult("file1.go\nfile2.go", "", 20)

	output := buf.String()
	assert.Contains(t, output, "file1.go")
	assert.Contains(t, output, "file2.go")
}

func TestDefaultPrinter_ToolResult_WithStderr(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinterWithWriter(&buf)

	p.ToolResult("", "error message", 20)

	output := buf.String()
	assert.Contains(t, output, "✗") // Error icon
	assert.Contains(t, output, "error message")
}

func TestDefaultPrinter_Text(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinterWithWriter(&buf)

	p.Text("Hello from Claude!")

	output := buf.String()
	// Claude Code style: no "Claude:" prefix
	assert.Contains(t, output, "Hello from Claude!")
}

func TestDefaultPrinter_Text_Empty(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinterWithWriter(&buf)

	p.Text("")

	output := buf.String()
	assert.Empty(t, output)
}

func TestDefaultPrinter_CommandHeader(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinterWithWriter(&buf)

	p.CommandHeader("create-story: test-123", "Long prompt here", 20)

	output := buf.String()
	// New box format shows workflow name and story key separately
	assert.Contains(t, output, "create-story")
	assert.Contains(t, output, "test-123")
	assert.Contains(t, output, "Command")
	assert.Contains(t, output, "Long prompt here")
	// Should have box characters
	assert.Contains(t, output, "╭")
	assert.Contains(t, output, "╰")
}

func TestDefaultPrinter_CommandFooter_Success(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinterWithWriter(&buf)

	p.CommandFooter(5*time.Second, true, 0)

	output := buf.String()
	assert.Contains(t, output, "Complete")
	assert.Contains(t, output, "✓")
}

func TestDefaultPrinter_CommandFooter_Failure(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinterWithWriter(&buf)

	p.CommandFooter(5*time.Second, false, 1)

	output := buf.String()
	assert.Contains(t, output, "Failed")
	assert.Contains(t, output, "exit 1")
	assert.Contains(t, output, "✗")
}

func TestDefaultPrinter_CycleHeader(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinterWithWriter(&buf)

	p.CycleHeader("test-story")

	output := buf.String()
	assert.Contains(t, output, "BMAD Full Cycle")
	assert.Contains(t, output, "test-story")
}

func TestDefaultPrinter_CycleSummary(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinterWithWriter(&buf)

	steps := []core.StepResult{
		{Name: "create-story", Duration: 10 * time.Second, Success: true},
		{Name: "dev-story", Duration: 30 * time.Second, Success: true},
	}

	p.CycleSummary("test-story", steps, 40*time.Second)

	output := buf.String()
	assert.Contains(t, output, "CYCLE COMPLETE")
	assert.Contains(t, output, "test-story")
	assert.Contains(t, output, "create-story")
	assert.Contains(t, output, "dev-story")
}

func TestDefaultPrinter_CycleFailed(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinterWithWriter(&buf)

	p.CycleFailed("test-story", "dev-story", 15*time.Second)

	output := buf.String()
	assert.Contains(t, output, "CYCLE FAILED")
	assert.Contains(t, output, "test-story")
	assert.Contains(t, output, "dev-story")
}

func TestDefaultPrinter_QueueHeader(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinterWithWriter(&buf)

	p.QueueHeader(3, []string{"story-1", "story-2", "story-3"})

	output := buf.String()
	assert.Contains(t, output, "BMAD Queue")
	assert.Contains(t, output, "3 stories")
}

func TestDefaultPrinter_QueueStoryStart(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinterWithWriter(&buf)

	p.QueueStoryStart(2, 5, "story-key")

	output := buf.String()
	assert.Contains(t, output, "Queue")
	assert.Contains(t, output, "[2/5]")
	assert.Contains(t, output, "story-key")
}

func TestDefaultPrinter_QueueSummary_Success(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinterWithWriter(&buf)

	results := []core.StoryResult{
		{Key: "story-1", Success: true, Duration: 10 * time.Second},
		{Key: "story-2", Success: true, Duration: 20 * time.Second},
	}

	p.QueueSummary(results, []string{"story-1", "story-2"}, 30*time.Second)

	output := buf.String()
	assert.Contains(t, output, "QUEUE COMPLETE")
	assert.Contains(t, output, "Completed: 2")
}

func TestDefaultPrinter_QueueSummary_WithFailure(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinterWithWriter(&buf)

	results := []core.StoryResult{
		{Key: "story-1", Success: true, Duration: 10 * time.Second},
		{Key: "story-2", Success: false, Duration: 5 * time.Second, FailedAt: "dev-story"},
	}

	p.QueueSummary(results, []string{"story-1", "story-2", "story-3"}, 15*time.Second)

	output := buf.String()
	assert.Contains(t, output, "QUEUE STOPPED")
	assert.Contains(t, output, "Failed: 1")
	assert.Contains(t, output, "Remaining: 1")
	assert.Contains(t, output, "(pending)")
}

func TestTruncateString(t *testing.T) {
	tests := []struct {
		input    string
		maxLen   int
		expected string
	}{
		{"short", 10, "short"},
		{"exactly10!", 10, "exactly10!"},
		{"this is a long string", 10, "this is..."},
		{"", 10, ""},
	}

	for _, tt := range tests {
		result := truncateString(tt.input, tt.maxLen)
		assert.Equal(t, tt.expected, result)
	}
}

func TestTruncateOutput(t *testing.T) {
	// Create 30 lines
	lines := make([]string, 30)
	for i := range lines {
		lines[i] = "line"
	}
	input := strings.Join(lines, "\n")

	result := truncateOutput(input, 10)

	// New format matches Claude Code: "… +N lines"
	assert.Contains(t, result, "… +20 lines")
}

func TestTruncateOutput_NoTruncation(t *testing.T) {
	input := "line1\nline2\nline3"
	result := truncateOutput(input, 10)

	assert.Equal(t, input, result)
}

func TestTruncateOutput_ZeroMaxLines(t *testing.T) {
	input := "line1\nline2\nline3"
	result := truncateOutput(input, 0)

	assert.Equal(t, input, result)
}

func TestBoxLineWrapWords(t *testing.T) {
	// Test that content wraps properly
	// boxLineWrapWords puts label on first line, content on subsequent lines
	lines := boxLineWrapWords("Label", "Short content", 50)
	assert.Len(t, lines, 2) // Label line + content line
	assert.Contains(t, lines[0], "Label:")
	assert.Contains(t, lines[1], "Short content")

	// Test with very long content
	longContent := strings.Repeat("word ", 30)
	lines = boxLineWrapWords("Command", longContent, 50)
	assert.Greater(t, len(lines), 2, "Long content should wrap to multiple lines")
}

func TestDefaultPrinter_ToolUse_Write(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinterWithWriter(&buf)

	p.ToolUse(core.ToolParams{
		Name:     "Write",
		FilePath: "/path/to/file.go",
		Content:  "package main\n\nfunc main() {}",
	})

	output := buf.String()
	assert.Contains(t, output, "Write")
	assert.Contains(t, output, "/path/to/file.go")
	// Should show diff with added lines
	assert.Contains(t, output, "+")
}

func TestDefaultPrinter_ToolUse_Task(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinterWithWriter(&buf)

	p.ToolUse(core.ToolParams{
		Name:         "Task",
		SubagentType: "Explore",
		Prompt:       "Find all configuration files in the codebase",
	})

	output := buf.String()
	assert.Contains(t, output, "Task")
	assert.Contains(t, output, "Explore")
	assert.Contains(t, output, "Find all configuration")
}

func TestDefaultPrinter_ToolUse_Task_LongPrompt(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinterWithWriter(&buf)

	longPrompt := strings.Repeat("word ", 50) // Very long prompt
	p.ToolUse(core.ToolParams{
		Name:         "Task",
		SubagentType: "Plan",
		Prompt:       longPrompt,
	})

	output := buf.String()
	assert.Contains(t, output, "Task")
	assert.Contains(t, output, "Plan")
	assert.Contains(t, output, "...")
}

func TestDefaultPrinter_ToolUse_NotebookEdit(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinterWithWriter(&buf)

	p.ToolUse(core.ToolParams{
		Name:         "NotebookEdit",
		NotebookPath: "/path/to/notebook.ipynb",
		CellID:       "cell-123",
		NewSource:    "print('hello')",
		EditMode:     "replace",
		CellType:     "code",
	})

	output := buf.String()
	assert.Contains(t, output, "NotebookEdit")
	assert.Contains(t, output, "notebook.ipynb")
	assert.Contains(t, output, "cell-123")
	assert.Contains(t, output, "replace")
}

func TestDefaultPrinter_ToolUse_AskUserQuestion(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinterWithWriter(&buf)

	p.ToolUse(core.ToolParams{
		Name: "AskUserQuestion",
		Questions: []claude.Question{
			{
				Question: "Which library should we use?",
				Header:   "Library",
				Options: []claude.QuestionOption{
					{Label: "React", Description: "Facebook's library"},
					{Label: "Vue", Description: "Progressive framework"},
				},
			},
		},
	})

	output := buf.String()
	assert.Contains(t, output, "AskUserQuestion")
	assert.Contains(t, output, "1 question")
	assert.Contains(t, output, "Library")
	assert.Contains(t, output, "Which library")
	assert.Contains(t, output, "React")
	assert.Contains(t, output, "Vue")
}

func TestDefaultPrinter_ToolUse_AskUserQuestion_Multiple(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinterWithWriter(&buf)

	p.ToolUse(core.ToolParams{
		Name: "AskUserQuestion",
		Questions: []claude.Question{
			{Question: "First question?"},
			{Question: "Second question?"},
		},
	})

	output := buf.String()
	assert.Contains(t, output, "2 questions")
}

func TestDefaultPrinter_ToolUse_Skill(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinterWithWriter(&buf)

	p.ToolUse(core.ToolParams{
		Name:  "Skill",
		Skill: "commit",
		Args:  "-m 'Fix bug'",
	})

	output := buf.String()
	assert.Contains(t, output, "Skill")
	assert.Contains(t, output, "commit")
	assert.Contains(t, output, "-m 'Fix bug'")
}

func TestDefaultPrinter_ToolUse_TodoWrite(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinterWithWriter(&buf)

	p.ToolUse(core.ToolParams{
		Name: "TodoWrite",
		Todos: []claude.TodoItem{
			{ID: "1", Content: "First task", Status: "completed"},
			{ID: "2", Content: "Second task", Status: "in_progress", ActiveForm: "Working on second task"},
			{ID: "3", Content: "Third task", Status: "pending"},
		},
	})

	output := buf.String()
	assert.Contains(t, output, "TodoWrite")
	assert.Contains(t, output, "3 todos")
	assert.Contains(t, output, "✓") // completed
	assert.Contains(t, output, "●") // in_progress
	assert.Contains(t, output, "○") // pending
	assert.Contains(t, output, "First task")
	assert.Contains(t, output, "Working on second task") // activeForm used for in_progress
	assert.Contains(t, output, "Third task")
}

func TestDefaultPrinter_ToolUse_TodoWrite_Empty(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinterWithWriter(&buf)

	p.ToolUse(core.ToolParams{
		Name:  "TodoWrite",
		Todos: []claude.TodoItem{},
	})

	output := buf.String()
	assert.Contains(t, output, "TodoWrite")
	assert.NotContains(t, output, "todos")
}

func TestDefaultPrinter_ToolUse_UnknownTool(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinterWithWriter(&buf)

	rawJSON := json.RawMessage(`{"foo":"bar","count":42}`)
	p.ToolUse(core.ToolParams{
		Name:     "FutureTool",
		InputRaw: rawJSON,
	})

	output := buf.String()
	assert.Contains(t, output, "FutureTool")
	assert.Contains(t, output, "count=42")
	assert.Contains(t, output, "foo=bar")
}

func TestFormatUnknownToolParams(t *testing.T) {
	tests := []struct {
		name     string
		raw      json.RawMessage
		contains []string
	}{
		{
			name:     "simple string",
			raw:      json.RawMessage(`{"name":"test"}`),
			contains: []string{"name=test"},
		},
		{
			name:     "number",
			raw:      json.RawMessage(`{"count":42}`),
			contains: []string{"count=42"},
		},
		{
			name:     "boolean",
			raw:      json.RawMessage(`{"enabled":true}`),
			contains: []string{"enabled=true"},
		},
		{
			name:     "multiple fields",
			raw:      json.RawMessage(`{"a":"1","b":"2"}`),
			contains: []string{"a=1", "b=2"},
		},
		{
			name:     "long string truncated",
			raw:      json.RawMessage(`{"long":"` + strings.Repeat("x", 50) + `"}`),
			contains: []string{"..."},
		},
		{
			name:     "empty array",
			raw:      json.RawMessage(`{"items":[]}`),
			contains: []string{"items=[]"},
		},
		{
			name:     "array with items",
			raw:      json.RawMessage(`{"items":["a","b","c"]}`),
			contains: []string{"items=[3 items]"},
		},
		{
			name:     "nested object",
			raw:      json.RawMessage(`{"config":{"a":1,"b":2}}`),
			contains: []string{"config={2 keys}"},
		},
		{
			name:     "null value",
			raw:      json.RawMessage(`{"value":null}`),
			contains: []string{"value=null"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatUnknownToolParams(tt.raw)
			for _, expected := range tt.contains {
				assert.Contains(t, result, expected)
			}
		})
	}
}

func TestCreateWriteDiff(t *testing.T) {
	content := "line1\nline2\nline3"
	d := createWriteDiff(content)

	assert.Equal(t, 3, d.Added)
	assert.Equal(t, 0, d.Deleted)
	assert.Len(t, d.Hunks, 1)
	assert.Len(t, d.Hunks[0].Lines, 3)

	// All lines should be additions
	for _, line := range d.Hunks[0].Lines {
		assert.Equal(t, diff.LineTypeAdded, line.Type)
	}
}

// Test helper functions (internal to package for testing)

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func truncateOutput(output string, maxLines int) string {
	if maxLines <= 0 {
		return output
	}

	lines := strings.Split(output, "\n")
	if len(lines) <= maxLines {
		return output
	}

	// Show first maxLines, then indicate how many more
	first := strings.Join(lines[:maxLines], "\n")
	remaining := len(lines) - maxLines

	return fmt.Sprintf("%s\n… +%d lines", first, remaining)
}

func boxLineWrapWords(label, content string, width int) []string {
	var lines []string
	innerWidth := width - 4 // "│ " + " │"

	// First line: "Label:"
	labelLine := label + ":"
	lines = append(lines, boxLine(labelLine, width))

	// Content lines: indented
	indent := "  " // 2 spaces indent for content
	contentWidth := innerWidth - len(indent)

	wrappedContent := wrapWords(content, contentWidth)
	for _, line := range wrappedContent {
		lines = append(lines, boxLine(indent+line, width))
	}

	return lines
}

func boxLine(content string, width int) string {
	if len(content) >= width-4 {
		content = truncateString(content, width-5)
	}
	padding := width - 4 - len(content)
	return "│ " + content + strings.Repeat(" ", padding) + " │"
}

func wrapWords(text string, maxWidth int) []string {
	if maxWidth <= 0 {
		return []string{text}
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{}
	}

	var lines []string
	var currentLine strings.Builder
	currentWidth := 0

	for i, word := range words {
		wordWidth := len(word)

		if currentWidth == 0 {
			if wordWidth > maxWidth {
				lines = append(lines, truncateString(word, maxWidth))
				continue
			}
			currentLine.WriteString(word)
			currentWidth = wordWidth
			continue
		}

		if currentWidth+1+wordWidth > maxWidth {
			lines = append(lines, currentLine.String())
			currentLine.Reset()

			if wordWidth > maxWidth {
				lines = append(lines, truncateString(word, maxWidth))
				currentWidth = 0
				continue
			}

			currentLine.WriteString(word)
			currentWidth = wordWidth
		} else {
			if i > 0 {
				currentLine.WriteString(" ")
				currentWidth++
			}
			currentLine.WriteString(word)
			currentWidth += wordWidth
		}
	}

	if currentLine.Len() > 0 {
		lines = append(lines, currentLine.String())
	}

	return lines
}

func formatUnknownToolParams(raw json.RawMessage) string {
	var params map[string]interface{}
	if err := json.Unmarshal(raw, &params); err != nil {
		return ""
	}

	// Sort keys for consistent output
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var parts []string
	for _, k := range keys {
		v := params[k]
		var valStr string
		switch val := v.(type) {
		case string:
			if len(val) > 30 {
				valStr = val[:27] + "..."
			} else {
				valStr = val
			}
		case bool:
			valStr = fmt.Sprintf("%v", val)
		case float64:
			valStr = fmt.Sprintf("%v", val)
		case []interface{}:
			if len(val) == 0 {
				valStr = "[]"
			} else {
				valStr = fmt.Sprintf("[%d items]", len(val))
			}
		case map[string]interface{}:
			valStr = fmt.Sprintf("{%d keys}", len(val))
		case nil:
			valStr = "null"
		default:
			valStr = fmt.Sprintf("<%T>", val)
		}
		parts = append(parts, k+"="+valStr)
	}

	result := strings.Join(parts, ", ")
	if len(result) > 80 {
		result = result[:77] + "..."
	}
	return result
}
