package tui

import (
	"context"
	"fmt"

	"bmaduum/internal/claude"
	"bmaduum/internal/config"

	tea "github.com/charmbracelet/bubbletea"
)

// Runner handles TUI-based workflow execution.
type Runner struct {
	executor claude.Executor
	config   *config.Config
}

// NewRunner creates a new TUI runner.
func NewRunner(executor claude.Executor, cfg *config.Config) *Runner {
	return &Runner{
		executor: executor,
		config:   cfg,
	}
}

// Run executes a workflow with the TUI.
func (r *Runner) Run(ctx context.Context, workflowName, storyKey string) int {
	prompt, err := r.config.GetPrompt(workflowName, storyKey)
	if err != nil {
		fmt.Printf("Error getting prompt: %v\n", err)
		return 1
	}

	model := r.config.GetModel(workflowName)

	// Create TUI model
	modelName := model
	if modelName == "" {
		modelName = "claude"
	}

	tuiModel := NewModel(r.executor, r.config, storyKey, modelName)
	tuiModel.SetStep(1, 1, workflowName)

	// Create Bubble Tea program
	p := tea.NewProgram(
		tuiModel,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	// Run Claude in background
	go func() {
		handler := func(event claude.Event) {
			p.Send(ClaudeEventAdapter(event))
		}

		exitCode, _ := r.executor.ExecuteWithResult(ctx, prompt, handler, model)
		p.Send(CompleteMsg{ExitCode: exitCode})
	}()

	// Run TUI (blocks until complete)
	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf("TUI error: %v\n", err)
		return 1
	}

	m := finalModel.(*Model)

	// Print final summary in regular terminal mode
	fmt.Print(m.RenderFinalSummary())

	return m.GetExitCode()
}

// RunMultiStep executes multiple workflow steps with the TUI.
func (r *Runner) RunMultiStep(ctx context.Context, steps []StepInfo, storyKey string) int {
	// Create TUI model
	tuiModel := NewModel(r.executor, r.config, storyKey, "claude")

	// Create Bubble Tea program
	p := tea.NewProgram(
		tuiModel,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	// Run workflows in background
	go func() {
		for i, step := range steps {
			// Update step info
			p.Send(StepStartMsg{
				Step:     i + 1,
				Total:    len(steps),
				StepName: step.Name,
				StoryKey: storyKey,
			})

			prompt, err := r.config.GetPrompt(step.Name, storyKey)
			if err != nil {
				p.Send(CompleteMsg{ExitCode: 1})
				return
			}

			model := r.config.GetModel(step.Name)

			handler := func(event claude.Event) {
				p.Send(ClaudeEventAdapter(event))
			}

			exitCode, _ := r.executor.ExecuteWithResult(ctx, prompt, handler, model)
			if exitCode != 0 {
				p.Send(CompleteMsg{ExitCode: exitCode})
				return
			}

			p.Send(StepCompleteMsg{Success: true})
		}

		p.Send(CompleteMsg{ExitCode: 0})
	}()

	// Run TUI
	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf("TUI error: %v\n", err)
		return 1
	}

	m := finalModel.(*Model)
	fmt.Print(m.RenderFinalSummary())

	return m.GetExitCode()
}

// StepInfo represents a single workflow step.
type StepInfo struct {
	Name      string
	StoryKey  string
	NextStatus string
}

// ClaudeEventAdapter converts a claude.Event to a tea.Msg.
func ClaudeEventAdapter(event claude.Event) tea.Msg {
	switch {
	case event.SessionStarted:
		return SessionStartEvent{}

	case event.IsText():
		return TextEvent{Text: event.Text}

	case event.IsToolUse():
		return ToolUseEvent{
			Name:        event.ToolName,
			Description: event.ToolDescription,
			Command:     event.ToolCommand,
			FilePath:    event.ToolFilePath,
		}

	case event.IsToolResult():
		return ToolResultEvent{
			Stdout: event.ToolStdout,
			Stderr: event.ToolStderr,
		}

	case event.SessionComplete:
		return SessionCompleteEvent{}

	default:
		return nil
	}
}

