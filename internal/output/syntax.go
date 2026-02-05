// Package output provides terminal output formatting using lipgloss styles.
//
// The package provides structured output for CLI operations including session
// lifecycle, step progress, tool usage display, and batch operation summaries.
// All output is styled using the lipgloss library for consistent terminal rendering.
//
// Key types:
//   - [Printer] - Interface for structured terminal output operations
//   - [DefaultPrinter] - Production implementation using lipgloss styles
//   - [StepResult] - Result of a single workflow step execution
//   - [StoryResult] - Result of processing a story in queue/epic operations
//
// Use [NewPrinter] for production output to stdout, or [NewPrinterWithWriter]
// to capture output in tests by providing a custom io.Writer.
package output

import (
	"fmt"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

// FormatDiff formats a git-style diff with colors.
//
// Parses unified diff format and applies color coding:
//   - Green (+) for added lines
//   - Red (-) for deleted lines
//   - Blue/cyan for diff headers (file paths, line numbers)
//   - Gray for context lines
//
// Example input:
//
//	--- a/file.txt
//	+++ b/file.txt
//	@@ -1,3 +1,4 @@
//	 line 1
//	-old line
//	+new line
//	 line 3
func FormatDiff(diff string) string {
	if !SupportsColor() {
		return diff
	}

	var result strings.Builder
	lines := strings.Split(diff, "\n")

	for i, line := range lines {
		var styled string
		switch {
		case strings.HasPrefix(line, "+++ ") || strings.HasPrefix(line, "--- "):
			// File path header
			styled = diffHeaderStyle.Render(line)
		case strings.HasPrefix(line, "@@ "):
			// Hunk header
			styled = diffMetaStyle.Render(line)
		case strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++"):
			// Added line (green)
			styled = diffAddStyle.Render(line)
		case strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---"):
			// Deleted line (red)
			styled = diffDelStyle.Render(line)
		case strings.HasPrefix(line, " "):
			// Context line (gray/muted)
			styled = diffMetaStyle.Render(line)
		default:
			// Other lines (index, new file mode, etc.)
			styled = line
		}

		result.WriteString(styled)
		if i < len(lines)-1 {
			result.WriteString("\n")
		}
	}

	return result.String()
}

// HighlightCode applies syntax highlighting to code using Chroma.
//
// The language parameter should be a language identifier recognized by
// Chroma (e.g., "go", "python", "javascript", "bash"). If the language
// is not recognized or is empty, uses a fallback lexer.
//
// Returns the syntax-highlighted code as a string with ANSI escape codes.
// If colors are not supported, returns the original code unchanged.
//
// Example:
//
//	highlighted := HighlightCode("func main() {}", "go")
func HighlightCode(code, language string) string {
	if !SupportsColor() || code == "" {
		return code
	}

	var lexer chroma.Lexer
	if language != "" {
		lexer = lexers.Get(language)
	}
	if lexer == nil {
		lexer = lexers.Fallback
	}

	// Configure the lexer for code
	lexer = chroma.Coalesce(lexer)

	// Parse the code
	it, err := lexer.Tokenise(nil, code)
	if err != nil {
		return code
	}

	// Create a GitHub-like style (matches Claude Code's look)
	// Use "github-dark" or fallback to monokai
	theme := styles.Get("github-dark")
	if theme == nil {
		theme = styles.Monokai
	}

	// Format as ANSI
	var buf strings.Builder
	formatter := formatters.Get("terminal256")
	if formatter == nil {
		formatter = formatters.NoOp
	}

	err = formatter.Format(&buf, theme, it)
	if err != nil {
		return code
	}

	return buf.String()
}

// DetectLanguage attempts to detect the programming language from a file path.
//
// Uses common file extensions to map to Chroma lexer names.
// Returns empty string if the language cannot be determined.
//
// Examples:
//
//	DetectLanguage("main.go")      → "go"
//	DetectLanguage("script.py")    → "python"
//	DetectLanguage("file.ts")      → "typescript"
func DetectLanguage(filePath string) string {
	// Get file extension
	if idx := strings.LastIndex(filePath, "."); idx != -1 {
		ext := strings.ToLower(filePath[idx+1:])

		// Map extensions to lexer names
		langMap := map[string]string{
			"go":       "go",
			"py":       "python",
			"js":       "javascript",
			"ts":       "typescript",
			"tsx":      "tsx",
			"jsx":      "jsx",
			"rs":       "rust",
			"c":        "c",
			"h":        "c",
			"cpp":      "cpp",
			"hpp":      "cpp",
			"cc":       "cpp",
			"java":     "java",
			"kt":       "kotlin",
			"swift":    "swift",
			"sh":       "bash",
			"bash":     "bash",
			"zsh":      "bash",
			"fish":     "fish",
			"yaml":     "yaml",
			"yml":      "yaml",
			"json":     "json",
			"toml":     "toml",
			"xml":      "xml",
			"html":     "html",
			"css":      "css",
			"scss":     "scss",
			"sql":      "sql",
			"md":       "markdown",
			"markdown": "markdown",
			"rb":       "ruby",
			"php":      "php",
			"lua":      "lua",
			"r":        "r",
			"scala":    "scala",
			"clj":      "clojure",
			"cljs":     "clojure",
			"ex":       "elixir",
			"exs":      "elixir",
			"erl":      "erlang",
			"hrl":      "erlang",
			"dart":     "dart",
			"vue":      "vue",
			"svelte":   "svelte",
		}

		if lang, ok := langMap[ext]; ok {
			return lang
		}
	}

	return ""
}

// IsDiffOutput checks if output looks like a git diff.
//
// Returns true if the output contains diff markers like "--- a/",
// "+++ b/", or "@@". This can be used to automatically apply
// diff formatting to tool output.
func IsDiffOutput(output string) bool {
	return strings.Contains(output, "--- a/") ||
		strings.Contains(output, "+++ b/") ||
		strings.Contains(output, "@@ ")
}

// TruncateForDisplay truncates output to fit within the terminal width.
//
// If a line is longer than the terminal width, it is truncated with "..."
// appended. If maxLines is specified, only that many lines are returned
// with a "... (+N more lines)" indicator if content was truncated.
//
// This is useful for displaying large tool outputs without overwhelming
// the terminal.
func TruncateForDisplay(output string, maxLines int) string {
	lines := strings.Split(output, "\n")

	if maxLines > 0 && len(lines) > maxLines {
		half := maxLines / 2
		result := strings.Join(lines[:half], "\n")
		omitted := len(lines) - maxLines
		result += "\n   ... (" + strings.Repeat(" ", 1) + fmt.Sprintf("%d", omitted) + " lines omitted) ...\n"
		result += strings.Join(lines[len(lines)-half:], "\n")
		return result
	}

	return output
}
