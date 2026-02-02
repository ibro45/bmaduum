package tui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"bmaduum/internal/claude"
	"bmaduum/internal/config"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// OutputSection represents a single output block in the TUI.
type OutputSection struct {
	ID       string
	Type     SectionType
	Content  string
	Lines    []string
	Rendered string
	Language string // For syntax highlighting
}

// SectionType identifies the kind of output section.
type SectionType int

const (
	SectionText SectionType = iota
	SectionToolUse
	SectionToolResult
	SectionDivider
)

// TypewriterState manages character-by-character text animation.
type TypewriterState struct {
	Buffer    string
	Displayed int
	Pending   []rune
	Speed     time.Duration
	Active    bool
}

// Model is the main Bubble Tea model for the TUI.
type Model struct {
	// Header state
	CurrentStep int
	TotalSteps  int
	StepName    string
	StoryKey    string
	ModelName   string
	StartTime   time.Time

	// Content state
	Sections       []OutputSection
	CurrentSection *OutputSection
	ContentBuilder strings.Builder

	// Viewport for scrolling
	Viewport viewport.Model

	// Animation state
	Typewriter TypewriterState
	Spinner    spinner.Model
	Thinking   bool
	ThinkText  string
	LastActivity time.Time

	// Runtime
	Executor claude.Executor
	Config   *config.Config
	Ctx      context.Context
	Cancel   context.CancelFunc

	// Dimensions
	Width  int
	Height int

	// State
	Err      error
	Quitting bool
	ExitCode int
}

// NewModel creates a new TUI model with default settings.
func NewModel(executor claude.Executor, cfg *config.Config, storyKey, modelName string) *Model {
	ctx, cancel := context.WithCancel(context.Background())

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#8B949E"))

	return &Model{
		StoryKey:     storyKey,
		ModelName:    modelName,
		StartTime:    time.Now(),
		Executor:     executor,
		Config:       cfg,
		Ctx:          ctx,
		Cancel:       cancel,
		Spinner:      s,
		LastActivity: time.Now(),
		Typewriter: TypewriterState{
			Speed: 5 * time.Millisecond,
		},
		Sections: make([]OutputSection, 0),
	}
}

// Init initializes the TUI model.
func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		m.Spinner.Tick,
		typewriterCmd(),
		thinkingDetectorCmd(),
	)
}

// Helper function to create a typewriter tick command.
func typewriterCmd() tea.Cmd {
	return tea.Tick(5*time.Millisecond, func(t time.Time) tea.Msg {
		return typewriterTickMsg{Time: t}
	})
}

type typewriterTickMsg struct {
	Time time.Time
}

// Helper function to detect when we should show thinking spinner.
func thinkingDetectorCmd() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return thinkingCheckMsg{Time: t}
	})
}

type thinkingCheckMsg struct {
	Time time.Time
}

// SetStep updates the current step information.
func (m *Model) SetStep(step, total int, name string) {
	m.CurrentStep = step
	m.TotalSteps = total
	m.StepName = name
}

// AddTextSection adds a new text section for streaming content.
func (m *Model) AddTextSection() {
	section := OutputSection{
		ID:      fmt.Sprintf("text-%d", len(m.Sections)),
		Type:    SectionText,
		Content: "",
	}
	m.Sections = append(m.Sections, section)
	m.CurrentSection = &m.Sections[len(m.Sections)-1]
	m.Typewriter.Active = true
	m.Typewriter.Displayed = 0
}

// AddToolUseSection adds a tool invocation section.
func (m *Model) AddToolUseSection(name, description, command, filePath string) {
	section := OutputSection{
		ID:      fmt.Sprintf("tool-%d", len(m.Sections)),
		Type:    SectionToolUse,
		Content: formatToolUse(name, description, command, filePath),
	}
	m.Sections = append(m.Sections, section)
	m.CurrentSection = nil // Tool use sections are complete immediately
	m.Typewriter.Active = false
}

// AddToolResultSection adds a tool result section.
func (m *Model) AddToolResultSection(stdout, stderr string) {
	section := OutputSection{
		ID:      fmt.Sprintf("result-%d", len(m.Sections)),
		Type:    SectionToolResult,
		Content: formatToolResult(stdout, stderr),
		Lines:   splitLines(stdout + stderr),
	}
	m.Sections = append(m.Sections, section)
	m.CurrentSection = nil
	m.Typewriter.Active = false
}

// AppendText appends text to the current text section with typewriter animation.
func (m *Model) AppendText(text string) {
	if m.CurrentSection == nil || m.CurrentSection.Type != SectionText {
		m.AddTextSection()
	}

	m.Typewriter.Buffer += text
	m.Typewriter.Pending = []rune(m.Typewriter.Buffer[m.Typewriter.Displayed:])
	m.LastActivity = time.Now()
	m.Thinking = false
}

// CompleteText marks the current text as complete (disable typewriter).
func (m *Model) CompleteText() {
	m.Typewriter.Active = false
	m.Typewriter.Displayed = len([]rune(m.Typewriter.Buffer))
	m.Typewriter.Pending = nil
}

// SetThinking sets the thinking state.
func (m *Model) SetThinking(thinking bool, status string) {
	m.Thinking = thinking
	m.ThinkText = status
	if thinking {
		m.LastActivity = time.Now()
	}
}

// formatToolUse formats a tool use display.
func formatToolUse(name, description, command, filePath string) string {
	var parts []string
	if name != "" {
		parts = append(parts, name)
	}
	if command != "" {
		parts = append(parts, fmt.Sprintf("(%s)", command))
	} else if filePath != "" {
		parts = append(parts, fmt.Sprintf("(%s)", filePath))
	} else if description != "" {
		parts = append(parts, fmt.Sprintf("(%s)", description))
	}
	return strings.Join(parts, "")
}

// formatToolResult formats a tool result display.
func formatToolResult(stdout, stderr string) string {
	var result strings.Builder
	if stdout != "" {
		result.WriteString(stdout)
	}
	if stderr != "" {
		if result.Len() > 0 {
			result.WriteString("\n")
		}
		result.WriteString(stderr)
	}
	return result.String()
}

// splitLines splits content into lines.
func splitLines(content string) []string {
	if content == "" {
		return nil
	}
	return strings.Split(content, "\n")
}

// ElapsedTime returns the formatted elapsed time since start.
func (m *Model) ElapsedTime() string {
	elapsed := time.Since(m.StartTime)
	minutes := int(elapsed.Minutes())
	seconds := int(elapsed.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

// IsComplete returns true if the TUI session is complete.
func (m *Model) IsComplete() bool {
	return m.Quitting
}

// GetExitCode returns the final exit code.
func (m *Model) GetExitCode() int {
	return m.ExitCode
}
