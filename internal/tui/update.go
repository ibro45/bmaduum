package tui

import (
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

// Update handles messages and updates the model.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyMsg(msg)

	case tea.MouseMsg:
		// Forward mouse events to viewport for scrolling
		var cmd tea.Cmd
		m.Viewport, cmd = m.Viewport.Update(msg)
		return m, cmd

	case tea.WindowSizeMsg:
		return m.handleWindowSizeMsg(msg)

	case typewriterTickMsg:
		return m.handleTypewriterTick()

	case thinkingCheckMsg:
		return m.handleThinkingCheck()

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.Spinner, cmd = m.Spinner.Update(msg)
		return m, cmd

	case CompleteMsg:
		m.ExitCode = msg.ExitCode
		m.Quitting = true
		return m, tea.Quit

	default:
		// Handle custom event messages from the adapter
		return m.handleEventMsg(msg)
	}

	return m, tea.Batch(cmds...)
}

// handleKeyMsg handles keyboard input.
func (m *Model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		m.Quitting = true
		m.Cancel()
		return m, tea.Quit
	default:
		// Forward to viewport for any scrolling keys
		var cmd tea.Cmd
		m.Viewport, cmd = m.Viewport.Update(msg)
		return m, cmd
	}
}

// handleWindowSizeMsg handles terminal resize events.
func (m *Model) handleWindowSizeMsg(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	m.Width = msg.Width
	m.Height = msg.Height

	// Reserve 2 lines for header
	viewportHeight := msg.Height - 2
	if viewportHeight < 1 {
		viewportHeight = 1
	}

	m.Viewport.Width = msg.Width
	m.Viewport.Height = viewportHeight

	// Re-render content
	m.Viewport.SetContent(m.renderContent())

	return m, nil
}

// handleTypewriterTick processes typewriter animation.
func (m *Model) handleTypewriterTick() (tea.Model, tea.Cmd) {
	if !m.Typewriter.Active || len(m.Typewriter.Pending) == 0 {
		return m, typewriterCmd()
	}

	// Display 3-5 characters per tick for efficiency
	batchSize := 4
	if len(m.Typewriter.Pending) < batchSize {
		batchSize = len(m.Typewriter.Pending)
	}

	m.Typewriter.Displayed += batchSize
	m.Typewriter.Pending = m.Typewriter.Pending[batchSize:]

	// Update the current section's rendered content
	if m.CurrentSection != nil {
		m.CurrentSection.Rendered = string([]rune(m.Typewriter.Buffer)[:m.Typewriter.Displayed])
		m.Viewport.SetContent(m.renderContent())
		m.Viewport.GotoBottom()
	}

	return m, typewriterCmd()
}

// handleThinkingCheck shows thinking spinner after inactivity.
func (m *Model) handleThinkingCheck() (tea.Model, tea.Cmd) {
	// Show thinking spinner after 500ms of inactivity
	if !m.Thinking && m.Typewriter.Active && len(m.Typewriter.Pending) == 0 {
		inactive := time.Since(m.LastActivity)
		if inactive > 500*time.Millisecond {
			m.Thinking = true
			m.ThinkText = "Claude is thinking..."
		}
	}

	return m, thinkingDetectorCmd()
}

// Update handles the custom event messages from the adapter
func (m *Model) handleEventMsg(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch e := msg.(type) {
	case TextEvent:
		m.AppendText(e.Text)
		m.LastActivity = time.Now()
		m.Viewport.SetContent(m.renderContent())
		m.Viewport.GotoBottom()
		return m, typewriterCmd()

	case ToolUseEvent:
		m.CompleteText()
		m.AddToolUseSection(e.Name, e.Description, e.Command, e.FilePath)
		m.LastActivity = time.Now()
		m.Viewport.SetContent(m.renderContent())
		m.Viewport.GotoBottom()

	case ToolResultEvent:
		m.AddToolResultSection(e.Stdout, e.Stderr)
		m.LastActivity = time.Now()
		m.Viewport.SetContent(m.renderContent())
		m.Viewport.GotoBottom()

	case StepStartMsg:
		m.SetStep(e.Step, e.Total, e.StepName)
		// Add step divider
		if e.Step > 1 {
			m.Sections = append(m.Sections, OutputSection{
				ID:      "divider",
				Type:    SectionDivider,
				Content: "",
			})
		}

	case SessionStartEvent:
		m.StartTime = time.Now()

	case SessionCompleteEvent:
		m.CompleteText()
		m.Thinking = false
	}

	return m, nil
}

// Quit requests the TUI to quit with the given exit code.
func (m *Model) Quit(exitCode int) tea.Cmd {
	m.ExitCode = exitCode
	return func() tea.Msg {
		return CompleteMsg{ExitCode: exitCode}
	}
}
