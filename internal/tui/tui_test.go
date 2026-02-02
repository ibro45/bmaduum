package tui

import (
	"testing"
	"time"

	"bmaduum/internal/tui/styles"

	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/assert"
)

func TestClaudeColors(t *testing.T) {
	// Verify color palette is properly defined
	assert.Equal(t, "#6B4EE6", styles.ClaudeColors.Primary)
	assert.Equal(t, "#58A6FF", styles.ClaudeColors.ToolIcon)
	assert.Equal(t, "#8B949E", styles.ClaudeColors.OutputIcon)
	assert.Equal(t, "#E6EDF3", styles.ClaudeColors.Text)
	assert.Equal(t, "#3FB950", styles.ClaudeColors.Success)
	assert.Equal(t, "#F85149", styles.ClaudeColors.Error)
}

func TestSymbols(t *testing.T) {
	// Verify symbols are properly defined
	assert.Equal(t, "⏺", styles.Symbols.ToolInvocation)
	assert.Equal(t, "⎿", styles.Symbols.ToolOutput)
	assert.Equal(t, "✓", styles.Symbols.Success)
	assert.Equal(t, "✗", styles.Symbols.Error)
	assert.Equal(t, "⚡", styles.Symbols.HeaderLogo)
	assert.Equal(t, "⏱️", styles.Symbols.Clock)
}

func TestSpinnerFrames(t *testing.T) {
	// Verify spinner frames are defined
	assert.Len(t, styles.SpinnerFrames, 10)
	assert.Equal(t, "⠋", styles.SpinnerFrames[0])
	assert.Equal(t, "⠏", styles.SpinnerFrames[9])
}

func TestDefaultStyles(t *testing.T) {
	// Verify default styles are properly initialized
	s := styles.DefaultStyles()

	// Header should have purple background
	headerBg := s.Header.GetBackground()
	assert.Equal(t, lipgloss.Color("#6B4EE6"), headerBg)

	// Success should have green foreground
	successFg := s.Success.GetForeground()
	assert.Equal(t, lipgloss.Color("#3FB950"), successFg)

	// Error should have red foreground
	errorFg := s.Error.GetForeground()
	assert.Equal(t, lipgloss.Color("#F85149"), errorFg)
}

func TestModel_SetStep(t *testing.T) {
	m := &Model{}

	m.SetStep(2, 4, "dev-story")

	assert.Equal(t, 2, m.CurrentStep)
	assert.Equal(t, 4, m.TotalSteps)
	assert.Equal(t, "dev-story", m.StepName)
}

func TestModel_AddTextSection(t *testing.T) {
	m := NewModel(nil, nil, "PROJ-123", "claude")

	m.AddTextSection()

	assert.Len(t, m.Sections, 1)
	assert.Equal(t, SectionText, m.Sections[0].Type)
	assert.NotNil(t, m.CurrentSection)
	assert.True(t, m.Typewriter.Active)
}

func TestModel_AddToolUseSection(t *testing.T) {
	m := NewModel(nil, nil, "PROJ-123", "claude")

	m.AddToolUseSection("Bash", "", "ls -la", "")

	assert.Len(t, m.Sections, 1)
	assert.Equal(t, SectionToolUse, m.Sections[0].Type)
	assert.Contains(t, m.Sections[0].Content, "Bash")
	assert.Contains(t, m.Sections[0].Content, "ls -la")
}

func TestModel_AddToolResultSection(t *testing.T) {
	m := NewModel(nil, nil, "PROJ-123", "claude")

	m.AddToolResultSection("output line 1\noutput line 2", "")

	assert.Len(t, m.Sections, 1)
	assert.Equal(t, SectionToolResult, m.Sections[0].Type)
	assert.Equal(t, "output line 1\noutput line 2", m.Sections[0].Content)
	assert.Len(t, m.Sections[0].Lines, 2)
}

func TestModel_AppendText(t *testing.T) {
	m := NewModel(nil, nil, "PROJ-123", "claude")
	m.AddTextSection()

	m.AppendText("Hello")

	assert.Equal(t, "Hello", m.Typewriter.Buffer)
	assert.Equal(t, []rune("Hello"), m.Typewriter.Pending)
	assert.True(t, m.Typewriter.Active)
}

func TestModel_CompleteText(t *testing.T) {
	m := NewModel(nil, nil, "PROJ-123", "claude")
	m.AddTextSection()
	m.AppendText("Hello World")

	m.CompleteText()

	assert.False(t, m.Typewriter.Active)
	assert.Equal(t, 11, m.Typewriter.Displayed) // "Hello World" = 11 chars
}

func TestModel_ElapsedTime(t *testing.T) {
	m := NewModel(nil, nil, "PROJ-123", "claude")
	m.StartTime = time.Now().Add(-2*time.Minute - 34*time.Second)

	elapsed := m.ElapsedTime()

	assert.Equal(t, "02:34", elapsed)
}

func TestFormatToolUse(t *testing.T) {
	tests := []struct {
		name        string
		toolName    string
		description string
		command     string
		filePath    string
		expected    string
	}{
		{
			name:     "command only",
			toolName: "Bash",
			command:  "ls -la",
			expected: "Bash(ls -la)",
		},
		{
			name:     "filepath only",
			toolName: "Read",
			filePath: "src/main.go",
			expected: "Read(src/main.go)",
		},
		{
			name:        "description only",
			toolName:    "Edit",
			description: "Update function",
			expected:    "Edit(Update function)",
		},
		{
			name:     "name only",
			toolName: "Think",
			expected: "Think",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatToolUse(tt.toolName, tt.description, tt.command, tt.filePath)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSplitLines(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected []string
	}{
		{
			name:     "multiple lines",
			content:  "line1\nline2\nline3",
			expected: []string{"line1", "line2", "line3"},
		},
		{
			name:     "single line",
			content:  "single",
			expected: []string{"single"},
		},
		{
			name:     "empty",
			content:  "",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := splitLines(tt.content)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestModel_GetExitCode(t *testing.T) {
	m := NewModel(nil, nil, "PROJ-123", "claude")

	assert.Equal(t, 0, m.GetExitCode())

	m.ExitCode = 1
	assert.Equal(t, 1, m.GetExitCode())
}

func TestModel_IsComplete(t *testing.T) {
	m := NewModel(nil, nil, "PROJ-123", "claude")

	assert.False(t, m.IsComplete())

	m.Quitting = true
	assert.True(t, m.IsComplete())
}
