// Package render provides output rendering components for the CLI.
//
// This package contains specialized renderers for different types of output
// including tool invocations, box drawing, and cycle/queue summaries.
package render

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strings"

	"bmaduum/internal/claude"
	"bmaduum/internal/output/core"
)

// lineNumberArrowPattern matches Claude CLI's line number format: "    1→content"
// Captures: optional whitespace, digits, arrow character, then content
var lineNumberArrowPattern = regexp.MustCompile(`^(\s*)(\d+)→(.*)$`)

// Re-export core types for convenience.
type (
	// ToolParams contains parameters for a tool invocation.
	ToolParams = core.ToolParams

	// StepResult represents the result of a single workflow step execution.
	StepResult = core.StepResult

	// StoryResult represents the result of processing a story in queue or epic operations.
	StoryResult = core.StoryResult
)

// StyleProvider provides styling functions for rendered output.
type StyleProvider interface {
	RenderHeader(s string) string
	RenderSuccess(s string) string
	RenderError(s string) string
	RenderMuted(s string) string
	RenderBullet(s string) string
	RenderToolName(s string) string
	RenderToolParams(s string) string
	RenderToolOutput(s string) string
	RenderQuestionHeader(s string) string
	RenderDiffSummary(s string) string
}

// DiffRenderer handles rendering diffs for tool output.
type DiffRenderer interface {
	RenderDiff(filePath, oldStr, newStr string) string
	RenderWriteDiff(filePath, content string) string
	RenderNotebookEdit(params core.ToolParams) string
	RenderUnifiedDiff(output string) string
	IsUnifiedDiff(output string) bool
	FormatDiff(output string) string
}

// ToolRenderer handles tool invocation and result display.
type ToolRenderer struct {
	out           io.Writer
	styles        StyleProvider
	diffRenderer  DiffRenderer
	truncateLines int
}

// NewToolRenderer creates a new tool renderer.
func NewToolRenderer(out io.Writer, styles StyleProvider, diffRenderer DiffRenderer, _ interface{}) *ToolRenderer {
	// The markdown parameter is kept for API compatibility but not used
	// since tool rendering doesn't need markdown formatting
	return &ToolRenderer{
		out:           out,
		styles:        styles,
		diffRenderer:  diffRenderer,
		truncateLines: 50, // Default truncation
	}
}

// SetTruncateLines sets the maximum number of lines to display for tool output.
func (r *ToolRenderer) SetTruncateLines(n int) {
	r.truncateLines = n
}

// Writeln writes a formatted line to the output.
func (r *ToolRenderer) Writeln(format string, args ...interface{}) {
	fmt.Fprintf(r.out, format+"\n", args...)
}

// ToolUse prints tool invocation details in Claude Code style.
// Format: ● ToolName(params) with green bullet, bold name, muted params.
// For Edit tools, also displays the diff being applied.
// For Write tools, displays the content as an all-additions diff.
// For Task tools, displays subagent type and truncated prompt.
// For NotebookEdit tools, displays path, cell info, and source diff.
// For AskUserQuestion tools, displays questions and options.
// For Skill tools, displays skill name and arguments.
// For unknown tools, displays formatted raw JSON parameters.
func (r *ToolRenderer) ToolUse(params ToolParams) {
	bullet := r.styles.RenderBullet(IconTool)
	toolName := r.styles.RenderToolName(params.Name)

	// Handle specific tool types with enhanced display
	switch params.Name {
	case "Edit":
		r.printToolHeader(bullet, toolName, params.FilePath)
		if params.OldString != "" || params.NewString != "" {
			r.renderEditDiff(params.FilePath, params.OldString, params.NewString)
		}
		return

	case "Write":
		r.printToolHeader(bullet, toolName, params.FilePath)
		if params.Content != "" {
			r.renderWriteContent(params.FilePath, params.Content)
		}
		return

	case "Task":
		paramStr := params.SubagentType
		if params.Prompt != "" {
			// Truncate prompt to first 60 chars
			prompt := params.Prompt
			if len(prompt) > 60 {
				prompt = prompt[:57] + "..."
			}
			paramStr += ": " + prompt
		}
		r.printToolHeader(bullet, toolName, paramStr)
		return

	case "NotebookEdit":
		paramStr := params.NotebookPath
		if params.CellID != "" {
			paramStr += " [" + params.CellID + "]"
		}
		if params.EditMode != "" {
			paramStr += " (" + params.EditMode + ")"
		}
		r.printToolHeader(bullet, toolName, paramStr)
		if params.NewSource != "" {
			r.renderNotebookEdit(params)
		}
		return

	case "AskUserQuestion":
		qCount := len(params.Questions)
		paramStr := fmt.Sprintf("%d question", qCount)
		if qCount != 1 {
			paramStr += "s"
		}
		r.printToolHeader(bullet, toolName, paramStr)
		if qCount > 0 {
			r.renderQuestions(params.Questions)
		}
		return

	case "Skill":
		paramStr := params.Skill
		if params.Args != "" {
			paramStr += " " + params.Args
		}
		r.printToolHeader(bullet, toolName, paramStr)
		return

	case "TodoWrite":
		count := len(params.Todos)
		if count == 0 {
			r.printToolHeader(bullet, toolName, "")
		} else {
			paramStr := fmt.Sprintf("%d todo", count)
			if count != 1 {
				paramStr += "s"
			}
			r.printToolHeader(bullet, toolName, paramStr)
			r.renderTodos(params.Todos)
		}
		return
	}

	// Default handling for known tools with standard params
	var paramStr string
	switch {
	case params.Command != "":
		// Bash tool - show command
		paramStr = params.Command
	case params.FilePath != "":
		// Read - show file path
		paramStr = params.FilePath
	case params.Pattern != "":
		// Glob, Grep - show pattern (and optionally path)
		if params.Path != "" {
			paramStr = params.Pattern + " in " + params.Path
		} else {
			paramStr = params.Pattern
		}
	case params.Query != "":
		// WebSearch - show query
		paramStr = params.Query
	case params.URL != "":
		// WebFetch - show URL
		paramStr = params.URL
	default:
		// Unknown tool - try to format raw JSON params
		if len(params.InputRaw) > 0 {
			paramStr = formatUnknownToolParams(params.InputRaw)
		}
	}

	r.printToolHeader(bullet, toolName, paramStr)
}

// printToolHeader prints the standard tool header line with Claude Code formatting.
// Format: "  ⏺ ToolName(params)" with 2-space leading indent.
func (r *ToolRenderer) printToolHeader(bullet, toolName, paramStr string) {
	if paramStr != "" {
		rendered := r.styles.RenderToolParams("(" + paramStr + ")")
		r.Writeln("%s%s %s%s", IndentToolUse, bullet, toolName, rendered)
	} else {
		r.Writeln("%s%s %s", IndentToolUse, bullet, toolName)
	}
}

// formatUnknownToolParams formats raw JSON input as key=value pairs.
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
		// Format value based on type
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
			// Handle arrays - show count or brief content
			if len(val) == 0 {
				valStr = "[]"
			} else {
				valStr = fmt.Sprintf("[%d items]", len(val))
			}
		case map[string]interface{}:
			// Handle nested objects
			valStr = fmt.Sprintf("{%d keys}", len(val))
		case nil:
			valStr = "null"
		default:
			// For other complex types, just indicate the type
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

// renderEditDiff renders the old_string -> new_string diff for Edit tools.
func (r *ToolRenderer) renderEditDiff(filePath, oldStr, newStr string) {
	bracket := r.styles.RenderToolOutput(IconOutput)

	// Use the diff renderer to render the diff
	diffStr := r.diffRenderer.RenderDiff(filePath, oldStr, newStr)

	// Print summary line if available
	lines := strings.Split(diffStr, "\n")
	if len(lines) > 0 {
		summary := lines[0]
		if strings.Contains(summary, "line") || strings.Contains(summary, "No changes") {
			fmt.Fprintf(r.out, "%s%s%s%s\n", IndentToolResult, bracket, SpaceAfterBracket, r.styles.RenderDiffSummary(summary))
			// Print remaining lines
			for _, line := range lines[1:] {
				if line != "" {
					fmt.Fprintf(r.out, "%s   %s\n", IndentToolResult, line)
				}
			}
			return
		}
	}

	// Print the full diff
	for _, line := range lines {
		if line != "" {
			fmt.Fprintf(r.out, "%s   %s\n", IndentToolResult, line)
		}
	}
}

// renderWriteContent renders file content as an all-additions diff.
func (r *ToolRenderer) renderWriteContent(filePath, content string) {
	bracket := r.styles.RenderToolOutput(IconOutput)

	// Use the diff renderer to render the write diff
	diffStr := r.diffRenderer.RenderWriteDiff(filePath, content)

	// Print summary line if available
	lines := strings.Split(diffStr, "\n")
	if len(lines) > 0 {
		summary := lines[0]
		if strings.Contains(summary, "line") || strings.Contains(summary, "No changes") {
			fmt.Fprintf(r.out, "%s%s%s%s\n", IndentToolResult, bracket, SpaceAfterBracket, r.styles.RenderDiffSummary(summary))
			// Print remaining lines (limit to reasonable size)
			maxLines := 30
			for i, line := range lines[1:] {
				if i >= maxLines {
					fmt.Fprintf(r.out, "%s   %s\n", IndentToolResult, r.styles.RenderToolOutput(fmt.Sprintf("... +%d more lines", len(lines)-maxLines-1)))
					break
				}
				if line != "" {
					fmt.Fprintf(r.out, "%s   %s\n", IndentToolResult, line)
				}
			}
			return
		}
	}

	// Print the full diff
	maxLines := 30
	for i, line := range lines {
		if i >= maxLines {
			fmt.Fprintf(r.out, "%s   %s\n", IndentToolResult, r.styles.RenderToolOutput(fmt.Sprintf("... +%d more lines", len(lines)-maxLines)))
			break
		}
		if line != "" {
			fmt.Fprintf(r.out, "%s   %s\n", IndentToolResult, line)
		}
	}
}

// renderNotebookEdit renders notebook cell edit information.
func (r *ToolRenderer) renderNotebookEdit(params ToolParams) {
	// Use the diff renderer to render the notebook edit
	diffStr := r.diffRenderer.RenderNotebookEdit(params)

	// Print output
	lines := strings.Split(diffStr, "\n")
	for _, line := range lines {
		if line != "" {
			fmt.Fprintf(r.out, "%s   %s\n", IndentToolResult, line)
		}
	}
}

// renderTodos renders TodoWrite todo items with status indicators.
func (r *ToolRenderer) renderTodos(todos []claude.TodoItem) {
	bracket := r.styles.RenderToolOutput(IconOutput)

	for _, todo := range todos {
		// Status icon
		var statusIcon string
		switch todo.Status {
		case "completed":
			statusIcon = r.styles.RenderSuccess("✓")
		case "in_progress":
			statusIcon = r.styles.RenderBullet("●")
		default: // pending
			statusIcon = r.styles.RenderMuted("○")
		}

		// Content - use ActiveForm for in_progress, otherwise Content
		content := todo.Content
		if todo.Status == "in_progress" && todo.ActiveForm != "" {
			content = todo.ActiveForm
		}

		fmt.Fprintf(r.out, "%s%s %s %s\n", IndentToolResult, bracket, statusIcon, r.styles.RenderToolOutput(content))
	}
}

// renderQuestions renders AskUserQuestion questions and options.
func (r *ToolRenderer) renderQuestions(questions []claude.Question) {
	bracket := r.styles.RenderToolOutput(IconOutput)

	for i, q := range questions {
		// Print question header
		header := ""
		if q.Header != "" {
			header = "[" + q.Header + "] "
		}
		fmt.Fprintf(r.out, "%s%s%s%s%s\n", IndentToolResult, bracket, SpaceAfterBracket, r.styles.RenderQuestionHeader(header), q.Question)

		// Print options
		for j, opt := range q.Options {
			prefix := fmt.Sprintf("  %d. ", j+1)
			optLine := prefix + opt.Label
			if opt.Description != "" {
				optLine += " - " + opt.Description
			}
			fmt.Fprintf(r.out, "%s   %s\n", IndentToolResult, r.styles.RenderToolOutput(optLine))
		}

		// Add spacing between questions
		if i < len(questions)-1 {
			fmt.Fprintln(r.out)
		}
	}
}

// ToolResult prints tool execution results in Claude Code style.
// Format: "    ⎿  content" with 4-space leading indent and 2 spaces after bracket.
// Errors show with red ✗ icon.
// Line number arrows (N→) from Claude CLI are converted to space-padded format.
func (r *ToolRenderer) ToolResult(stdout, stderr string, truncateLines int) {
	bracket := r.styles.RenderToolOutput(IconOutput)

	if stdout == "" && stderr == "" {
		// Show "(No content)" when there's no output
		r.Writeln("%s%s%s%s", IndentToolResult, bracket, SpaceAfterBracket, r.styles.RenderToolOutput("(No content)"))
		return
	}

	if stdout != "" {
		// Clean up line number arrows from Claude CLI Read output
		output := stripLineNumberArrows(stdout)
		output = truncateOutput(output, truncateLines)

		// Check if this looks like a unified diff
		if r.diffRenderer.IsUnifiedDiff(output) {
			// Try to render with rich formatting
			rendered := r.diffRenderer.RenderUnifiedDiff(output)

			// Print rendered output
			lines := strings.Split(rendered, "\n")
			for i, line := range lines {
				if i == 0 && strings.Contains(line, "line") {
					// First line is the summary
					fmt.Fprintf(r.out, "%s%s%s%s\n", IndentToolResult, bracket, SpaceAfterBracket, r.styles.RenderDiffSummary(line))
				} else if line != "" {
					// Continuation lines align with content after bracket
					fmt.Fprintf(r.out, "%s   %s\n", IndentToolResult, line)
				}
			}
			return
		}

		// Print with continuation marker (muted gray)
		lines := strings.Split(output, "\n")
		for i, line := range lines {
			if i == 0 {
				fmt.Fprintf(r.out, "%s%s%s%s\n", IndentToolResult, bracket, SpaceAfterBracket, r.styles.RenderToolOutput(line))
			} else if line != "" {
				// Continuation lines align with content after bracket (4 + 1 + 2 = 7 spaces)
				fmt.Fprintf(r.out, "%s   %s\n", IndentToolResult, r.styles.RenderToolOutput(line))
			} else {
				fmt.Fprintln(r.out)
			}
		}
	}
	if stderr != "" {
		// Show errors with red X icon
		errorIcon := r.styles.RenderError(IconError)
		fmt.Fprintf(r.out, "%s%s%s%s %s\n", IndentToolResult, bracket, errorIcon, SpaceAfterBracket, r.styles.RenderError(stderr))
	}
}

// truncateOutput truncates output to maxLines, showing first portion with count.
// Format matches Claude Code: shows first lines then "… +N lines".
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

// stripLineNumberArrows removes Claude CLI's line number arrow format from output.
// Converts "    1→content" to "  1  content" (space-padded line numbers).
// This makes the output match Claude Code's native display format.
func stripLineNumberArrows(output string) string {
	lines := strings.Split(output, "\n")
	hasArrows := false

	// First pass: check if output has arrow format and find max line number width
	maxLineNum := 0
	for _, line := range lines {
		if match := lineNumberArrowPattern.FindStringSubmatch(line); match != nil {
			hasArrows = true
			lineNumStr := match[2]
			if len(lineNumStr) > maxLineNum {
				maxLineNum = len(lineNumStr)
			}
		}
	}

	if !hasArrows {
		return output
	}

	// Second pass: convert arrow format to space-padded format
	result := make([]string, len(lines))
	for i, line := range lines {
		if match := lineNumberArrowPattern.FindStringSubmatch(line); match != nil {
			lineNum := match[2]
			content := match[3]
			// Format: right-aligned line number + 2 spaces + content
			// e.g., "  1  content" or "100  content"
			result[i] = fmt.Sprintf("%*s  %s", maxLineNum, lineNum, content)
		} else {
			result[i] = line
		}
	}

	return strings.Join(result, "\n")
}
