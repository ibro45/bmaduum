package tui

import (
	"fmt"
	"strings"

	"bmaduum/internal/tui/styles"

	"github.com/charmbracelet/lipgloss"
)

// View renders the TUI.
func (m *Model) View() string {
	if m.Quitting {
		return ""
	}

	var sb strings.Builder

	// Render header
	sb.WriteString(m.renderHeader())
	sb.WriteString("\n")

	// Render viewport content
	sb.WriteString(m.Viewport.View())

	return sb.String()
}

// renderHeader renders the status header bar.
func (m *Model) renderHeader() string {
	width := m.Width
	if width == 0 {
		width = 80
	}

	// Build header sections
	var leftParts []string
	leftParts = append(leftParts, styles.Symbols.HeaderLogo+" bmaduum")

	if m.StepName != "" {
		stepInfo := fmt.Sprintf("Step %d/%d: %s", m.CurrentStep, m.TotalSteps, m.StepName)
		leftParts = append(leftParts, stepInfo)
	}

	if m.StoryKey != "" {
		leftParts = append(leftParts, m.StoryKey)
	}

	var rightParts []string
	if m.ModelName != "" {
		rightParts = append(rightParts, m.ModelName)
	}
	rightParts = append(rightParts, styles.Symbols.Clock+" "+m.ElapsedTime())

	// Join sections
	leftContent := strings.Join(leftParts, " │ ")
	rightContent := strings.Join(rightParts, " │ ")

	// Calculate padding
	separator := " │ "
	totalContentLen := lipgloss.Width(leftContent) + lipgloss.Width(separator) + lipgloss.Width(rightContent)
	padding := width - totalContentLen

	if padding < 1 {
		padding = 1
	}

	// Build full header line
	fullLine := leftContent + strings.Repeat(" ", padding) + rightContent

	// Apply header style
	headerStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(styles.ClaudeColors.Primary)).
		Foreground(lipgloss.Color("#FFFFFF")).
		Bold(true).
		Width(width)

	return headerStyle.Render(fullLine)
}

// renderContent renders all output sections.
func (m *Model) renderContent() string {
	var sb strings.Builder

	for _, section := range m.Sections {
		switch section.Type {
		case SectionText:
			sb.WriteString(m.renderTextSection(section))
		case SectionToolUse:
			sb.WriteString(m.renderToolUseSection(section))
		case SectionToolResult:
			sb.WriteString(m.renderToolResultSection(section))
		case SectionDivider:
			sb.WriteString(m.renderDivider())
		}
		sb.WriteString("\n")
	}

	// Add thinking spinner if active
	if m.Thinking {
		sb.WriteString("\n")
		spinnerLine := m.Spinner.View() + " " + styles.Default.TextMuted.Render(m.ThinkText)
		sb.WriteString(spinnerLine)
	}

	return sb.String()
}

// renderTextSection renders a text content section.
func (m *Model) renderTextSection(section OutputSection) string {
	content := section.Rendered
	if content == "" && section.Content != "" {
		content = section.Content
	}

	// Wrap text to viewport width minus some margin
	width := m.Width - 4
	if width < 40 {
		width = 40
	}

	lines := wrapText(content, width)
	return strings.Join(lines, "\n")
}

// renderToolUseSection renders a tool invocation section.
func (m *Model) renderToolUseSection(section OutputSection) string {
	var sb strings.Builder

	// Tool icon with name
	toolIcon := styles.Default.ToolIcon.Render(styles.Symbols.ToolInvocation)
	toolName := styles.Default.Text.Render(section.Content)

	sb.WriteString(toolIcon + " " + toolName)

	return sb.String()
}

// renderToolResultSection renders a tool result section.
func (m *Model) renderToolResultSection(section OutputSection) string {
	var sb strings.Builder

	if section.Content == "" {
		return ""
	}

	// Output icon
	outputIcon := styles.Default.OutputIcon.Render(styles.Symbols.ToolOutput)

	// Process lines
	lines := strings.Split(section.Content, "\n")
	for i, line := range lines {
		if i == 0 {
			// First line has the icon
			sb.WriteString(outputIcon + "  " + styles.Default.ToolOutput.Render(line))
		} else {
			// Subsequent lines are indented
			indent := strings.Repeat(" ", 4)
			sb.WriteString(indent + styles.Default.ToolOutput.Render(line))
		}
		if i < len(lines)-1 {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

// renderDivider renders a step transition divider.
func (m *Model) renderDivider() string {
	width := m.Width - 4
	if width < 20 {
		width = 20
	}

	line := strings.Repeat(styles.Symbols.Divider, width/2)
	divider := "── " + line + " ──"

	return styles.Default.SectionDivider.Render(divider)
}

// wrapText wraps text to a maximum width.
func wrapText(text string, width int) []string {
	var lines []string
	paragraphs := strings.Split(text, "\n")

	for _, paragraph := range paragraphs {
		if paragraph == "" {
			lines = append(lines, "")
			continue
		}

		words := strings.Fields(paragraph)
		var currentLine strings.Builder

		for _, word := range words {
			if currentLine.Len()+len(word)+1 > width {
				lines = append(lines, currentLine.String())
				currentLine.Reset()
			}
			if currentLine.Len() > 0 {
				currentLine.WriteString(" ")
			}
			currentLine.WriteString(word)
		}

		if currentLine.Len() > 0 {
			lines = append(lines, currentLine.String())
		}
	}

	return lines
}

// RenderFinalSummary renders a final summary after the TUI exits.
func (m *Model) RenderFinalSummary() string {
	var sb strings.Builder

	success := m.ExitCode == 0
	duration := m.ElapsedTime()

	sb.WriteString("\n")
	if success {
		sb.WriteString(styles.Default.Success.Render(styles.Symbols.Success + " Session complete"))
	} else {
		sb.WriteString(styles.Default.Error.Render(styles.Symbols.Error + " Session failed"))
	}
	sb.WriteString(" (" + duration + ")\n")

	return sb.String()
}
