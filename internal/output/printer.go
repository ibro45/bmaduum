// Package output provides terminal output formatting using lipgloss styles.
//
// The package provides structured output for CLI operations including session
// lifecycle, step progress, tool usage display, and batch operation summaries.
// All output is styled using the lipgloss library for consistent terminal rendering.
//
// Key types:
//   - [DefaultPrinter] - Production implementation using lipgloss styles
//
// The core types (Printer, StepResult, StoryResult, ToolParams) are in the
// core sub-package. Use [NewPrinter] for production output to stdout, or
// [NewPrinterWithWriter] to capture output in tests.
package output

import (
	"fmt"
	"io"
	"os"
	"time"

	"bmaduum/internal/output/core"
	"bmaduum/internal/output/diff"
	"bmaduum/internal/output/render"
)

// DefaultPrinter implements [core.Printer] with lipgloss terminal styling.
//
// It composites the specialized renderers from the render package:
//   - SessionRenderer for session lifecycle
//   - ToolRenderer for tool invocations
//   - CycleRenderer for cycle and queue operations
type DefaultPrinter struct {
	out           io.Writer
	session       *render.SessionRenderer
	tool          *render.ToolRenderer
	cycle         *render.CycleRenderer
	styleProvider *defaultStyleProvider
	widthProvider *defaultWidthProvider
	markdown      *MarkdownRenderer
}

// NewPrinter creates a new [DefaultPrinter] that writes to stdout.
func NewPrinter() *DefaultPrinter {
	return NewPrinterWithWriter(os.Stdout)
}

// NewPrinterWithWriter creates a new [DefaultPrinter] with a custom writer.
// This is useful for tests to capture output.
func NewPrinterWithWriter(w io.Writer) *DefaultPrinter {
	return NewPrinterWithConfig(w, DefaultMarkdownConfig())
}

// NewPrinterWithConfig creates a new [DefaultPrinter] with custom markdown configuration.
func NewPrinterWithConfig(w io.Writer, cfg MarkdownConfig) *DefaultPrinter {
	styleProvider := newDefaultStyleProvider()
	widthProvider := newDefaultWidthProvider()
	markdown := NewMarkdownRendererWithConfig(cfg)

	// Create the tool diff adapter
	diffAdapter := newToolDiffAdapter()

	// Create the printer first with minimal fields
	p := &DefaultPrinter{
		out:           w,
		styleProvider: styleProvider,
		widthProvider: widthProvider,
		markdown:      markdown,
	}

	// Now create the renderers with the printer reference
	session := render.NewSessionRenderer(w, styleProvider, widthProvider, markdown.Render)
	tool := render.NewToolRenderer(w, styleProvider, diffAdapter, markdown.Render)
	cycle := render.NewCycleRenderer(&outputWriterAdapter{printer: p}, styleProvider, widthProvider)

	// Assign the renderers to the printer
	p.session = session
	p.tool = tool
	p.cycle = cycle

	return p
}

// SessionStart prints session start indicator.
func (p *DefaultPrinter) SessionStart() {
	p.session.SessionStart()
}

// SessionEnd prints session end with status.
func (p *DefaultPrinter) SessionEnd(duration time.Duration, success bool) {
	p.session.SessionEnd(duration, success)
}

// StepStart prints step start header.
func (p *DefaultPrinter) StepStart(step, total int, name string) {
	// No output - step info is in CommandHeader and ProgressLine
}

// StepEnd prints step completion status.
func (p *DefaultPrinter) StepEnd(duration time.Duration, success bool) {
	// Handled by CommandFooter and ProgressLine
}

// ToolUse prints tool invocation details.
func (p *DefaultPrinter) ToolUse(params core.ToolParams) {
	// Convert core.ToolParams to render.ToolParams
	renderParams := render.ToolParams{
		Name:         params.Name,
		Description:  params.Description,
		Command:      params.Command,
		FilePath:     params.FilePath,
		OldString:    params.OldString,
		NewString:    params.NewString,
		Pattern:      params.Pattern,
		Query:        params.Query,
		URL:          params.URL,
		Path:         params.Path,
		Content:      params.Content,
		InputRaw:     params.InputRaw,
		SubagentType: params.SubagentType,
		Prompt:       params.Prompt,
		NotebookPath: params.NotebookPath,
		CellID:       params.CellID,
		NewSource:    params.NewSource,
		EditMode:     params.EditMode,
		CellType:     params.CellType,
		Questions:    params.Questions,
		Skill:        params.Skill,
		Args:         params.Args,
		Todos:        params.Todos,
	}
	p.tool.ToolUse(renderParams)
}

// ToolResult prints tool execution results.
func (p *DefaultPrinter) ToolResult(stdout, stderr string, truncateLines int) {
	p.tool.ToolResult(stdout, stderr, truncateLines)
}

// Text prints a text message from Claude.
func (p *DefaultPrinter) Text(message string) {
	p.session.Text(message)
}

// Divider prints a visual divider.
func (p *DefaultPrinter) Divider() {
	p.session.Divider()
}

// CycleHeader prints the header for a full cycle run.
func (p *DefaultPrinter) CycleHeader(storyKey string) {
	p.cycle.CycleHeader(storyKey)
}

// CycleSummary prints the summary after a successful cycle.
func (p *DefaultPrinter) CycleSummary(storyKey string, steps []core.StepResult, totalDuration time.Duration) {
	// Convert core.StepResult to render.StepResult
	renderSteps := make([]render.StepResult, len(steps))
	for i, s := range steps {
		renderSteps[i] = render.StepResult{
			Name:     s.Name,
			Duration: s.Duration,
			Success:  s.Success,
		}
	}
	p.cycle.CycleSummary(storyKey, renderSteps, totalDuration)
}

// CycleFailed prints failure information when a cycle fails.
func (p *DefaultPrinter) CycleFailed(storyKey string, failedStep string, duration time.Duration) {
	p.cycle.CycleFailed(storyKey, failedStep, duration)
}

// QueueHeader prints the header for a queue run.
func (p *DefaultPrinter) QueueHeader(count int, stories []string) {
	p.cycle.QueueHeader(count, stories)
}

// QueueStoryStart prints the header for starting a story in a queue.
func (p *DefaultPrinter) QueueStoryStart(index, total int, storyKey string) {
	p.cycle.QueueStoryStart(index, total, storyKey)
}

// QueueSummary prints the summary after a queue completes.
func (p *DefaultPrinter) QueueSummary(results []core.StoryResult, allKeys []string, totalDuration time.Duration) {
	// Convert core.StoryResult to render.StoryResult
	renderResults := make([]render.StoryResult, len(results))
	for i, r := range results {
		renderResults[i] = render.StoryResult{
			Key:      r.Key,
			Success:  r.Success,
			Duration: r.Duration,
			FailedAt: r.FailedAt,
			Skipped:  r.Skipped,
		}
	}
	p.cycle.QueueSummary(renderResults, allKeys, totalDuration)
}

// CommandHeader prints a nice box with command information.
func (p *DefaultPrinter) CommandHeader(label, prompt string, truncateLength int) {
	p.session.CommandHeader(label, prompt, truncateLength)
}

// CommandFooter prints the footer after a command completes.
func (p *DefaultPrinter) CommandFooter(duration time.Duration, success bool, exitCode int) {
	p.session.CommandFooter(duration, success, exitCode)
}

// defaultStyleProvider implements render.StyleProvider using lipgloss styles.
type defaultStyleProvider struct{}

func newDefaultStyleProvider() *defaultStyleProvider {
	return &defaultStyleProvider{}
}

func (p *defaultStyleProvider) RenderHeader(s string) string {
	return headerStyle.Render(s)
}

func (p *defaultStyleProvider) RenderSuccess(s string) string {
	return successStyle.Render(s)
}

func (p *defaultStyleProvider) RenderError(s string) string {
	return errorStyle.Render(s)
}

func (p *defaultStyleProvider) RenderMuted(s string) string {
	return mutedStyle.Render(s)
}

func (p *defaultStyleProvider) RenderDivider(s string) string {
	return dividerStyle.Render(s)
}

func (p *defaultStyleProvider) RenderBullet(s string) string {
	return bulletStyle.Render(s)
}

func (p *defaultStyleProvider) RenderToolName(s string) string {
	return toolNameStyle.Render(s)
}

func (p *defaultStyleProvider) RenderToolParams(s string) string {
	return toolParamsStyle.Render(s)
}

func (p *defaultStyleProvider) RenderToolOutput(s string) string {
	return toolOutputStyle.Render(s)
}

func (p *defaultStyleProvider) RenderQuestionHeader(s string) string {
	return questionHeaderStyle.Render(s)
}

func (p *defaultStyleProvider) RenderDiffSummary(s string) string {
	return diffSummaryStyle.Render(s)
}

func (p *defaultStyleProvider) RenderText(s string) string {
	return textStyle.Render(s)
}

// defaultWidthProvider implements width providers.
type defaultWidthProvider struct{}

func newDefaultWidthProvider() *defaultWidthProvider {
	return &defaultWidthProvider{}
}

func (p *defaultWidthProvider) TerminalWidth() int {
	return TerminalWidth()
}

// toolDiffAdapter implements render.DiffRenderer for tool output.
type toolDiffAdapter struct {
	renderer *diff.Renderer
}

func newToolDiffAdapter() *toolDiffAdapter {
	return &toolDiffAdapter{
		renderer: diff.NewRenderer(
			diff.WithLineNumbers(false),
			diff.WithGutter(true),
			diff.WithHighlighter(HighlightCode),
		),
	}
}

func (a *toolDiffAdapter) RenderDiff(filePath, oldStr, newStr string) string {
	d := createEditDiff(oldStr, newStr)
	lang := DetectLanguage(filePath)
	a.renderer = diff.NewRenderer(
		diff.WithLineNumbers(false),
		diff.WithGutter(true),
		diff.WithSyntaxHighlight(lang),
		diff.WithHighlighter(HighlightCode),
	)
	return a.renderer.RenderWithSummary(d)
}

func (a *toolDiffAdapter) RenderWriteDiff(filePath, content string) string {
	d := createWriteDiff(content)
	lang := DetectLanguage(filePath)
	a.renderer = diff.NewRenderer(
		diff.WithLineNumbers(true),
		diff.WithGutter(true),
		diff.WithSyntaxHighlight(lang),
		diff.WithHighlighter(HighlightCode),
	)
	return a.renderer.RenderWithSummary(d)
}

func (a *toolDiffAdapter) RenderNotebookEdit(params core.ToolParams) string {
	// Create info line
	info := ""
	if params.EditMode != "" {
		info += "mode: " + params.EditMode
	}
	if params.CellType != "" {
		if info != "" {
			info += ", "
		}
		info += "type: " + params.CellType
	}
	var result string
	if info != "" {
		result = toolOutputStyle.Render(iconOutput) + " " + toolOutputStyle.Render(info) + "\n"
	}

	// Create diff for the new source
	d := createWriteDiff(params.NewSource)
	lang := "python"
	if params.CellType == "markdown" {
		lang = "markdown"
	}

	a.renderer = diff.NewRenderer(
		diff.WithLineNumbers(true),
		diff.WithGutter(true),
		diff.WithSyntaxHighlight(lang),
		diff.WithHighlighter(HighlightCode),
	)

	return result + a.renderer.Render(d)
}

func (a *toolDiffAdapter) RenderUnifiedDiff(output string) string {
	d, err := diff.ParseUnifiedDiff(output)
	if err != nil || len(d.Hunks) == 0 {
		// Fall back to basic formatting
		return FormatDiff(output)
	}

	// Detect language from file path
	lang := ""
	if d.NewFile != "" {
		lang = DetectLanguage(d.NewFile)
	}

	a.renderer = diff.NewRenderer(
		diff.WithLineNumbers(true),
		diff.WithGutter(true),
		diff.WithSyntaxHighlight(lang),
		diff.WithHighlighter(HighlightCode),
	)

	return a.renderer.RenderWithSummary(d)
}

func (a *toolDiffAdapter) IsUnifiedDiff(output string) bool {
	return diff.IsUnifiedDiff(output)
}

func (a *toolDiffAdapter) FormatDiff(output string) string {
	return FormatDiff(output)
}

// outputWriterAdapter allows the DefaultPrinter to be used as an OutputWriter for CycleRenderer.
type outputWriterAdapter struct {
	printer *DefaultPrinter
}

func (a *outputWriterAdapter) Writeln(format string, args ...interface{}) {
	fmt.Fprintf(a.printer.out, format+"\n", args...)
}

func (a *outputWriterAdapter) Divider() {
	a.printer.Divider()
}

// createEditDiff creates a Diff from old/new strings for Edit tool display.
func createEditDiff(oldStr, newStr string) *diff.Diff {
	d := &diff.Diff{}

	var lines []diff.Line

	// Add deleted lines (from old_string)
	if oldStr != "" {
		oldLines := splitLines(oldStr)
		for i, line := range oldLines {
			lines = append(lines, diff.Line{
				Type:       diff.LineTypeDeleted,
				Content:    line,
				OldLineNum: i + 1,
			})
			d.Deleted++
		}
	}

	// Add added lines (from new_string)
	if newStr != "" {
		newLines := splitLines(newStr)
		for i, line := range newLines {
			lines = append(lines, diff.Line{
				Type:       diff.LineTypeAdded,
				Content:    line,
				NewLineNum: i + 1,
			})
			d.Added++
		}
	}

	if len(lines) > 0 {
		d.Hunks = []diff.Hunk{{Lines: lines}}
	}

	return d
}

// createWriteDiff creates an all-additions Diff from content for Write tool display.
func createWriteDiff(content string) *diff.Diff {
	d := &diff.Diff{}

	var lines []diff.Line
	contentLines := splitLines(content)
	for i, line := range contentLines {
		lines = append(lines, diff.Line{
			Type:       diff.LineTypeAdded,
			Content:    line,
			NewLineNum: i + 1,
		})
		d.Added++
	}

	if len(lines) > 0 {
		d.Hunks = []diff.Hunk{{Lines: lines}}
	}

	return d
}

// splitLines splits a string into lines, preserving the final empty line if present.
func splitLines(s string) []string {
	if s == "" {
		return []string{}
	}
	lines := make([]string, 0, 1)
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	// Add the final line if there's remaining content or if the string ends with newline
	if start < len(s) || (len(s) > 0 && s[len(s)-1] == '\n') {
		lines = append(lines, s[start:])
	}
	return lines
}
