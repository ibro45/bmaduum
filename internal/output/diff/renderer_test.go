package diff

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenderer_Render(t *testing.T) {
	diff := &Diff{
		OldFile: "test.go",
		NewFile: "test.go",
		Added:   1,
		Deleted: 1,
		Hunks: []Hunk{
			{
				OldStart: 1,
				OldCount: 2,
				NewStart: 1,
				NewCount: 2,
				Lines: []Line{
					{Type: LineTypeContext, Content: "context", OldLineNum: 1, NewLineNum: 1},
					{Type: LineTypeDeleted, Content: "old", OldLineNum: 2},
					{Type: LineTypeAdded, Content: "new", NewLineNum: 2},
				},
			},
		},
	}

	renderer := NewRenderer()
	output := renderer.Render(diff)

	// Should contain line numbers
	assert.Contains(t, output, "1")
	assert.Contains(t, output, "2")

	// Should contain content
	assert.Contains(t, output, "context")
	assert.Contains(t, output, "old")
	assert.Contains(t, output, "new")

	// Should contain markers
	assert.Contains(t, output, "+")
	assert.Contains(t, output, "-")
}

func TestRenderer_Render_NilDiff(t *testing.T) {
	renderer := NewRenderer()
	output := renderer.Render(nil)

	assert.Equal(t, "", output)
}

func TestRenderer_Render_EmptyDiff(t *testing.T) {
	diff := &Diff{
		OldFile: "test.go",
		NewFile: "test.go",
		Hunks:   []Hunk{},
	}

	renderer := NewRenderer()
	output := renderer.Render(diff)

	assert.Equal(t, "", output)
}

func TestRenderer_RenderWithSummary(t *testing.T) {
	diff := &Diff{
		Added:   3,
		Deleted: 2,
		Hunks:   []Hunk{},
	}

	renderer := NewRenderer()
	output := renderer.RenderWithSummary(diff)

	// Should contain summary
	assert.Contains(t, output, "3 lines added")
	assert.Contains(t, output, "2 lines removed")
}

func TestRenderer_RenderWithSummary_NilDiff(t *testing.T) {
	renderer := NewRenderer()
	output := renderer.RenderWithSummary(nil)

	assert.Equal(t, "", output)
}

func TestRenderer_WithOptions(t *testing.T) {
	diff := &Diff{
		Hunks: []Hunk{
			{
				Lines: []Line{
					{Type: LineTypeAdded, Content: "test", NewLineNum: 1},
				},
			},
		},
		Added: 1,
	}

	t.Run("without line numbers", func(t *testing.T) {
		renderer := NewRenderer(WithLineNumbers(false))
		output := renderer.Render(diff)
		// Should still have the content
		assert.Contains(t, output, "test")
	})

	t.Run("without gutter", func(t *testing.T) {
		renderer := NewRenderer(WithGutter(false))
		output := renderer.Render(diff)
		// Should not have gutter marker at start
		assert.False(t, strings.HasPrefix(output, "┃"))
	})

	t.Run("with max content length", func(t *testing.T) {
		longContent := strings.Repeat("x", 100)
		diff := &Diff{
			Hunks: []Hunk{
				{
					Lines: []Line{
						{Type: LineTypeAdded, Content: longContent, NewLineNum: 1},
					},
				},
			},
		}
		renderer := NewRenderer(WithMaxContentLength(10))
		output := renderer.Render(diff)
		// Should be truncated
		assert.Contains(t, output, "...")
	})

	t.Run("with syntax highlight", func(t *testing.T) {
		renderer := NewRenderer(WithSyntaxHighlight("go"))
		assert.True(t, renderer.syntaxHL)
		assert.Equal(t, "go", renderer.language)
	})
}

func TestRenderer_Render_MultipleHunks(t *testing.T) {
	diff := &Diff{
		Hunks: []Hunk{
			{
				OldStart: 1,
				OldCount: 1,
				NewStart: 1,
				NewCount: 1,
				Lines: []Line{
					{Type: LineTypeAdded, Content: "first", NewLineNum: 1},
				},
			},
			{
				OldStart: 10,
				OldCount: 1,
				NewStart: 10,
				NewCount: 1,
				Lines: []Line{
					{Type: LineTypeAdded, Content: "second", NewLineNum: 10},
				},
			},
		},
		Added: 2,
	}

	renderer := NewRenderer()
	output := renderer.Render(diff)

	// Should contain separator between hunks
	assert.Contains(t, output, "...")
	// Should contain both lines
	assert.Contains(t, output, "first")
	assert.Contains(t, output, "second")
}

func TestRenderer_calculateLineNumWidth(t *testing.T) {
	tests := []struct {
		name          string
		diff          *Diff
		expectedWidth int
	}{
		{
			name: "small line numbers",
			diff: &Diff{
				Hunks: []Hunk{
					{NewStart: 1, NewCount: 5, OldStart: 1, OldCount: 5},
				},
			},
			expectedWidth: 3, // At least 3
		},
		{
			name: "large line numbers",
			diff: &Diff{
				Hunks: []Hunk{
					{NewStart: 100, NewCount: 50, OldStart: 1, OldCount: 1},
				},
			},
			expectedWidth: 3, // 150 -> 3 digits
		},
		{
			name: "very large line numbers",
			diff: &Diff{
				Hunks: []Hunk{
					{NewStart: 1000, NewCount: 100, OldStart: 1, OldCount: 1},
				},
			},
			expectedWidth: 4, // 1100 -> 4 digits
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			renderer := NewRenderer()
			width := renderer.calculateLineNumWidth(tt.diff)
			assert.Equal(t, tt.expectedWidth, width)
		})
	}
}

func TestRenderer_renderPlain(t *testing.T) {
	diff := &Diff{
		Hunks: []Hunk{
			{
				Lines: []Line{
					{Type: LineTypeContext, Content: "context", OldLineNum: 1, NewLineNum: 1},
					{Type: LineTypeAdded, Content: "added", NewLineNum: 2},
					{Type: LineTypeDeleted, Content: "deleted", OldLineNum: 2},
				},
			},
		},
	}

	renderer := NewRenderer()
	output := renderer.renderPlain(diff)

	// Should have markers
	assert.Contains(t, output, " ")
	assert.Contains(t, output, "+")
	assert.Contains(t, output, "-")

	// Should have content
	assert.Contains(t, output, "context")
	assert.Contains(t, output, "added")
	assert.Contains(t, output, "deleted")

	// Should not have ANSI codes (plain text)
	assert.NotContains(t, output, "\x1b[")
}

func TestRenderer_renderHunkSeparator(t *testing.T) {
	t.Run("with gutter", func(t *testing.T) {
		renderer := NewRenderer(WithGutter(true))
		separator := renderer.renderHunkSeparator()
		// Should start with spaces for gutter
		assert.True(t, strings.HasPrefix(separator, "  "))
		assert.Contains(t, separator, "...")
	})

	t.Run("without gutter", func(t *testing.T) {
		renderer := NewRenderer(WithGutter(false))
		separator := renderer.renderHunkSeparator()
		// Should not start with gutter prefix (line number style + spaces, not gutter spaces)
		// Since lineNumWidth is at least 3, the separator should start with at least 3 spaces
		// But without gutter, it shouldn't have the exact "  " + lineNumStyle prefix pattern
		assert.Contains(t, separator, "...")
		// The separator without gutter should be shorter than with gutter
		rendererWithGutter := NewRenderer(WithGutter(true))
		separatorWithGutter := rendererWithGutter.renderHunkSeparator()
		assert.Less(t, len(separator), len(separatorWithGutter))
	})
}

func TestRenderer_renderLine(t *testing.T) {
	tests := []struct {
		name             string
		line             Line
		shouldContain    []string
		shouldNotContain []string
	}{
		{
			name:             "context line",
			line:             Line{Type: LineTypeContext, Content: "context", OldLineNum: 1, NewLineNum: 1},
			shouldContain:    []string{"context", " "},
			shouldNotContain: []string{"\x1b[48;2;"}, // No background color
		},
		{
			name:             "added line",
			line:             Line{Type: LineTypeAdded, Content: "added", NewLineNum: 2},
			shouldContain:    []string{"added", "+"},
			shouldNotContain: []string{},
		},
		{
			name:             "deleted line",
			line:             Line{Type: LineTypeDeleted, Content: "deleted", OldLineNum: 2},
			shouldContain:    []string{"deleted", "-"},
			shouldNotContain: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			renderer := NewRenderer()
			var buf strings.Builder
			renderer.renderLine(&buf, &tt.line)
			output := buf.String()

			for _, s := range tt.shouldContain {
				assert.Contains(t, output, s)
			}
			for _, s := range tt.shouldNotContain {
				assert.NotContains(t, output, s)
			}
		})
	}
}

func TestApplyDiffStyle(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		bgCode   string
		fgCode   string
		expected string
	}{
		{
			name:     "both bg and fg",
			content:  "test",
			bgCode:   "\x1b[48;2;13;40;24m",
			fgCode:   "\x1b[38;2;63;185;80m",
			expected: "\x1b[48;2;13;40;24m\x1b[38;2;63;185;80m" + "test" + "\x1b[0m",
		},
		{
			name:     "only bg",
			content:  "test",
			bgCode:   "\x1b[48;2;13;40;24m",
			fgCode:   "",
			expected: "\x1b[48;2;13;40;24m" + "test" + "\x1b[0m",
		},
		{
			name:     "content with reset",
			content:  "test\x1b[0mmore",
			bgCode:   "\x1b[48;2;13;40;24m",
			fgCode:   "\x1b[38;2;63;185;80m",
			expected: "\x1b[48;2;13;40;24m\x1b[38;2;63;185;80m" + "test" + "\x1b[0m\x1b[48;2;13;40;24m\x1b[38;2;63;185;80m" + "more" + "\x1b[0m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := applyDiffStyle(tt.content, tt.bgCode, tt.fgCode)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRenderer_renderHunk(t *testing.T) {
	hunk := &Hunk{
		Lines: []Line{
			{Type: LineTypeContext, Content: "context", OldLineNum: 1, NewLineNum: 1},
			{Type: LineTypeAdded, Content: "added", NewLineNum: 2},
		},
	}

	renderer := NewRenderer()
	var buf strings.Builder
	renderer.renderHunk(&buf, hunk)
	output := buf.String()

	assert.Contains(t, output, "context")
	assert.Contains(t, output, "added")
}

// Edge case tests

func TestRenderer_Render_OnlyContextLines(t *testing.T) {
	diff := &Diff{
		Hunks: []Hunk{
			{
				Lines: []Line{
					{Type: LineTypeContext, Content: "line 1", OldLineNum: 1, NewLineNum: 1},
					{Type: LineTypeContext, Content: "line 2", OldLineNum: 2, NewLineNum: 2},
				},
			},
		},
	}

	renderer := NewRenderer()
	output := renderer.Render(diff)

	assert.Contains(t, output, "line 1")
	assert.Contains(t, output, "line 2")
}

func TestRenderer_Render_OnlyAddedLines(t *testing.T) {
	diff := &Diff{
		Hunks: []Hunk{
			{
				Lines: []Line{
					{Type: LineTypeAdded, Content: "new 1", NewLineNum: 1},
					{Type: LineTypeAdded, Content: "new 2", NewLineNum: 2},
				},
			},
		},
		Added: 2,
	}

	renderer := NewRenderer()
	output := renderer.Render(diff)

	assert.Contains(t, output, "new 1")
	assert.Contains(t, output, "new 2")
	assert.Contains(t, output, "+")
}

func TestRenderer_Render_OnlyDeletedLines(t *testing.T) {
	diff := &Diff{
		Hunks: []Hunk{
			{
				Lines: []Line{
					{Type: LineTypeDeleted, Content: "old 1", OldLineNum: 1},
					{Type: LineTypeDeleted, Content: "old 2", OldLineNum: 2},
				},
			},
		},
		Deleted: 2,
	}

	renderer := NewRenderer()
	output := renderer.Render(diff)

	assert.Contains(t, output, "old 1")
	assert.Contains(t, output, "old 2")
	assert.Contains(t, output, "-")
}

func TestRenderer_Render_EmptyLineContent(t *testing.T) {
	diff := &Diff{
		Hunks: []Hunk{
			{
				Lines: []Line{
					{Type: LineTypeContext, Content: "", OldLineNum: 1, NewLineNum: 1},
				},
			},
		},
	}

	renderer := NewRenderer()
	output := renderer.Render(diff)

	// Should render even with empty content
	assert.NotEmpty(t, output)
}

func TestRenderer_Render_LongContent(t *testing.T) {
	longContent := strings.Repeat("x", 10000)
	diff := &Diff{
		Hunks: []Hunk{
			{
				Lines: []Line{
					{Type: LineTypeAdded, Content: longContent, NewLineNum: 1},
				},
			},
		},
	}

	renderer := NewRenderer()
	output := renderer.Render(diff)

	// Should handle long content
	assert.Contains(t, output, longContent)
}

func TestRenderer_Render_SpecialCharacters(t *testing.T) {
	// Test with carriage return and newline (tab may be converted to spaces by lipgloss)
	specialContent := "\n\r"
	diff := &Diff{
		Hunks: []Hunk{
			{
				Lines: []Line{
					{Type: LineTypeContext, Content: specialContent, OldLineNum: 1, NewLineNum: 1},
				},
			},
		},
	}

	renderer := NewRenderer()
	output := renderer.Render(diff)

	// Should preserve newline and carriage return
	assert.Contains(t, output, "\n")
	assert.Contains(t, output, "\r")
}

func TestRenderer_Render_UnicodeContent(t *testing.T) {
	unicodeContent := "Hello 世界 🌍"
	diff := &Diff{
		Hunks: []Hunk{
			{
				Lines: []Line{
					{Type: LineTypeContext, Content: unicodeContent, OldLineNum: 1, NewLineNum: 1},
				},
			},
		},
	}

	renderer := NewRenderer()
	output := renderer.Render(diff)

	// Should preserve unicode
	assert.Contains(t, output, unicodeContent)
}

func TestRenderer_RenderWithMaxContentLength(t *testing.T) {
	longContent := strings.Repeat("x", 100)
	diff := &Diff{
		Hunks: []Hunk{
			{
				Lines: []Line{
					{Type: LineTypeAdded, Content: longContent, NewLineNum: 1},
				},
			},
		},
	}

	renderer := NewRenderer(WithMaxContentLength(10))
	output := renderer.Render(diff)

	// Should be truncated with ...
	assert.Contains(t, output, "...")
	assert.NotContains(t, output, longContent)
}

func TestRenderer_RenderWithoutLineNumbers(t *testing.T) {
	diff := &Diff{
		Hunks: []Hunk{
			{
				Lines: []Line{
					{Type: LineTypeAdded, Content: "test", NewLineNum: 1},
				},
			},
		},
	}

	renderer := NewRenderer(WithLineNumbers(false))
	output := renderer.Render(diff)

	// Should have content and marker but no line number
	assert.Contains(t, output, "test")
	assert.Contains(t, output, "+")
	// No numeric digits at the start (except possibly as part of content)
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if line != "" {
			// First non-space character should be + or space, not a digit
			trimmed := strings.TrimLeft(line, " ")
			if len(trimmed) > 0 {
				assert.NotRegexp(t, `^\d`, trimmed)
			}
		}
	}
}

func TestRenderer_RenderWithGutter(t *testing.T) {
	diff := &Diff{
		Hunks: []Hunk{
			{
				Lines: []Line{
					{Type: LineTypeAdded, Content: "test", NewLineNum: 1},
					{Type: LineTypeContext, Content: "context", OldLineNum: 1, NewLineNum: 1},
				},
			},
		},
	}

	renderer := NewRenderer(WithGutter(true))
	output := renderer.Render(diff)

	// Should have gutter marker for added line
	assert.Contains(t, output, "┃")
}

func TestRenderer_RenderWithoutGutter(t *testing.T) {
	diff := &Diff{
		Hunks: []Hunk{
			{
				Lines: []Line{
					{Type: LineTypeAdded, Content: "test", NewLineNum: 1},
				},
			},
		},
	}

	renderer := NewRenderer(WithGutter(false))
	output := renderer.Render(diff)

	// Should not have gutter marker
	assert.NotContains(t, output, "┃")
}
