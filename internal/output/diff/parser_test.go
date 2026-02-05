package diff

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParser_Parse(t *testing.T) {
	input := `--- a/file.go
+++ b/file.go
@@ -1,3 +1,4 @@
 line 1
-old line
+new line
+added line
 line 3`

	parser := NewParser()
	diff, err := parser.Parse(input)
	require.NoError(t, err)

	assert.Equal(t, "file.go", diff.OldFile)
	assert.Equal(t, "file.go", diff.NewFile)
	assert.Equal(t, 2, diff.Added)
	assert.Equal(t, 1, diff.Deleted)
	assert.Len(t, diff.Hunks, 1)

	hunk := diff.Hunks[0]
	assert.Equal(t, 1, hunk.OldStart)
	assert.Equal(t, 3, hunk.OldCount)
	assert.Equal(t, 1, hunk.NewStart)
	assert.Equal(t, 4, hunk.NewCount)
	assert.Len(t, hunk.Lines, 5)

	// Check line types
	assert.Equal(t, LineTypeContext, hunk.Lines[0].Type)
	assert.Equal(t, "line 1", hunk.Lines[0].Content)

	assert.Equal(t, LineTypeDeleted, hunk.Lines[1].Type)
	assert.Equal(t, "old line", hunk.Lines[1].Content)

	assert.Equal(t, LineTypeAdded, hunk.Lines[2].Type)
	assert.Equal(t, "new line", hunk.Lines[2].Content)

	assert.Equal(t, LineTypeAdded, hunk.Lines[3].Type)
	assert.Equal(t, "added line", hunk.Lines[3].Content)

	assert.Equal(t, LineTypeContext, hunk.Lines[4].Type)
	assert.Equal(t, "line 3", hunk.Lines[4].Content)
}

func TestParseUnifiedDiff(t *testing.T) {
	input := `--- a/file.go
+++ b/file.go
@@ -1,3 +1,4 @@
 line 1
-old line
+new line
+added line
 line 3`

	diff, err := ParseUnifiedDiff(input)
	require.NoError(t, err)

	assert.Equal(t, "file.go", diff.OldFile)
	assert.Equal(t, "file.go", diff.NewFile)
	assert.Equal(t, 2, diff.Added)
	assert.Equal(t, 1, diff.Deleted)
}

func TestParser_Parse_MultipleHunks(t *testing.T) {
	input := `--- a/file.go
+++ b/file.go
@@ -1,2 +1,2 @@
 context
-old1
+new1
@@ -10,2 +10,2 @@
 more context
-old2
+new2`

	parser := NewParser()
	diff, err := parser.Parse(input)
	require.NoError(t, err)

	assert.Len(t, diff.Hunks, 2)
	assert.Equal(t, 2, diff.Added)
	assert.Equal(t, 2, diff.Deleted)

	// First hunk
	assert.Equal(t, 1, diff.Hunks[0].OldStart)
	assert.Equal(t, 1, diff.Hunks[0].NewStart)

	// Second hunk
	assert.Equal(t, 10, diff.Hunks[1].OldStart)
	assert.Equal(t, 10, diff.Hunks[1].NewStart)
}

func TestParser_Parse_NoCount(t *testing.T) {
	// When count is 1, it can be omitted: @@ -1 +1 @@
	input := `--- a/file.go
+++ b/file.go
@@ -1 +1 @@
-old
+new`

	parser := NewParser()
	diff, err := parser.Parse(input)
	require.NoError(t, err)

	assert.Len(t, diff.Hunks, 1)
	hunk := diff.Hunks[0]
	assert.Equal(t, 1, hunk.OldCount)
	assert.Equal(t, 1, hunk.NewCount)
}

func TestParser_Parse_EmptyDiff(t *testing.T) {
	input := `--- a/file.go
+++ b/file.go`

	parser := NewParser()
	diff, err := parser.Parse(input)
	require.NoError(t, err)

	assert.Equal(t, "file.go", diff.OldFile)
	assert.Equal(t, "file.go", diff.NewFile)
	assert.Len(t, diff.Hunks, 0)
}

func TestParser_Parse_WithAPrefix(t *testing.T) {
	input := `--- a/file.go
+++ b/file.go
@@ -1,3 +1,4 @@
 line 1
-old line
+new line`

	parser := NewParser()
	diff, err := parser.Parse(input)
	require.NoError(t, err)

	// Should strip the a/ and b/ prefixes
	assert.Equal(t, "file.go", diff.OldFile)
	assert.Equal(t, "file.go", diff.NewFile)
}

func TestParser_Parse_EmptyLines(t *testing.T) {
	input := "--- a/file.go\n" +
		"+++ b/file.go\n" +
		"@@ -1,3 +1,3 @@\n" +
		" line 1\n" +
		"-old line\n" +
		"+\n" +
		" line 2"

	parser := NewParser()
	diff, err := parser.Parse(input)
	require.NoError(t, err)

	hunk := diff.Hunks[0]
	// We should have 4 lines: context, deleted, added (empty), context
	assert.Equal(t, 4, len(hunk.Lines))
}

func TestParser_LineNumbers(t *testing.T) {
	input := `--- a/file.go
+++ b/file.go
@@ -10,5 +10,5 @@
 context line 10
-old line 11
+new line 11
 context line 12`

	parser := NewParser()
	diff, err := parser.Parse(input)
	require.NoError(t, err)

	hunk := diff.Hunks[0]
	lines := hunk.Lines

	// Context line should have both line numbers
	assert.Equal(t, 10, lines[0].OldLineNum)
	assert.Equal(t, 10, lines[0].NewLineNum)

	// Deleted line should only have old line number
	assert.Equal(t, 11, lines[1].OldLineNum)
	assert.Equal(t, 0, lines[1].NewLineNum)

	// Added line should only have new line number
	assert.Equal(t, 0, lines[2].OldLineNum)
	assert.Equal(t, 11, lines[2].NewLineNum)
}

func TestDiff_Summary(t *testing.T) {
	tests := []struct {
		name     string
		added    int
		deleted  int
		expected string
	}{
		{"no changes", 0, 0, "No changes"},
		{"one added", 1, 0, "1 line added"},
		{"multiple added", 5, 0, "5 lines added"},
		{"one deleted", 0, 1, "1 line removed"},
		{"multiple deleted", 0, 3, "3 lines removed"},
		{"both", 2, 1, "2 lines added, 1 line removed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diff := &Diff{Added: tt.added, Deleted: tt.deleted}
			assert.Equal(t, tt.expected, diff.Summary())
		})
	}
}

func TestIsUnifiedDiff(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name: "valid unified diff",
			input: `--- a/file.go
+++ b/file.go
@@ -1,3 +1,4 @@
 context`,
			expected: true,
		},
		{
			name:     "not a diff",
			input:    "just some regular text",
			expected: false,
		},
		{
			name:     "partial diff markers",
			input:    "--- a/file.go\n+++ b/file.go",
			expected: false, // missing @@ marker
		},
		{
			name:     "has @@ but no file headers",
			input:    "@@ -1,3 +1,4 @@\n context",
			expected: false,
		},
		{
			name:     "diff without a/ prefix",
			input:    "--- file.go\n+++ file.go\n@@ -1,3 +1,4 @@\n context",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, IsUnifiedDiff(tt.input))
		})
	}
}

func TestParser_parseHunkHeader(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *Hunk
	}{
		{
			name:  "full format",
			input: "@@ -1,3 +1,4 @@",
			expected: &Hunk{
				OldStart: 1, OldCount: 3,
				NewStart: 1, NewCount: 4,
			},
		},
		{
			name:  "no old count",
			input: "@@ -1 +1,2 @@",
			expected: &Hunk{
				OldStart: 1, OldCount: 1,
				NewStart: 1, NewCount: 2,
			},
		},
		{
			name:  "no new count",
			input: "@@ -1,2 +1 @@",
			expected: &Hunk{
				OldStart: 1, OldCount: 2,
				NewStart: 1, NewCount: 1,
			},
		},
		{
			name:  "large line numbers",
			input: "@@ -100,50 +150,60 @@",
			expected: &Hunk{
				OldStart: 100, OldCount: 50,
				NewStart: 150, NewCount: 60,
			},
		},
		{
			name:     "invalid format",
			input:    "not a hunk header",
			expected: nil,
		},
	}

	parser := NewParser()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parser.parseHunkHeader(tt.input)
			if tt.expected == nil {
				assert.Nil(t, result)
			} else {
				require.NotNil(t, result)
				assert.Equal(t, tt.expected.OldStart, result.OldStart)
				assert.Equal(t, tt.expected.OldCount, result.OldCount)
				assert.Equal(t, tt.expected.NewStart, result.NewStart)
				assert.Equal(t, tt.expected.NewCount, result.NewCount)
			}
		})
	}
}

func TestPluralize(t *testing.T) {
	assert.Equal(t, "1 line", pluralize(1, "line"))
	assert.Equal(t, "0 lines", pluralize(0, "line"))
	assert.Equal(t, "5 lines", pluralize(5, "line"))
	assert.Equal(t, "100 files", pluralize(100, "file"))
}

func TestItoa(t *testing.T) {
	tests := []struct {
		input    int
		expected string
	}{
		{0, "0"},
		{1, "1"},
		{42, "42"},
		{123, "123"},
		{-1, "-1"},
		{-42, "-42"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			assert.Equal(t, tt.expected, itoa(tt.input))
		})
	}
}

// Edge case tests

func TestParser_Parse_OnlyAdditions(t *testing.T) {
	input := `--- a/file.go
+++ b/file.go
@@ -1,1 +1,2 @@
 line 1
+new line`

	parser := NewParser()
	diff, err := parser.Parse(input)
	require.NoError(t, err)

	assert.Equal(t, 1, diff.Added)
	assert.Equal(t, 0, diff.Deleted)
}

func TestParser_Parse_OnlyDeletions(t *testing.T) {
	input := `--- a/file.go
+++ b/file.go
@@ -1,2 +1,1 @@
 line 1
-old line`

	parser := NewParser()
	diff, err := parser.Parse(input)
	require.NoError(t, err)

	assert.Equal(t, 0, diff.Added)
	assert.Equal(t, 1, diff.Deleted)
}

func TestParser_Parse_NewFile(t *testing.T) {
	input := `--- /dev/null
+++ b/newfile.go
@@ -0,0 +1,3 @@
+line 1
+line 2
+line 3`

	parser := NewParser()
	diff, err := parser.Parse(input)
	require.NoError(t, err)

	assert.Equal(t, "/dev/null", diff.OldFile)
	assert.Equal(t, "newfile.go", diff.NewFile)
	assert.Equal(t, 3, diff.Added)
	assert.Equal(t, 0, diff.Deleted)
}

func TestParser_Parse_DeletedFile(t *testing.T) {
	input := `--- a/oldfile.go
+++ /dev/null
@@ -1,3 +0,0 @@
-line 1
-line 2
-line 3`

	parser := NewParser()
	diff, err := parser.Parse(input)
	require.NoError(t, err)

	assert.Equal(t, "oldfile.go", diff.OldFile)
	assert.Equal(t, "/dev/null", diff.NewFile)
	assert.Equal(t, 0, diff.Added)
	assert.Equal(t, 3, diff.Deleted)
}

func TestParser_Parse_BinaryFile(t *testing.T) {
	input := `Binary files a/file.bin and b/file.bin differ`

	parser := NewParser()
	diff, err := parser.Parse(input)
	require.NoError(t, err)

	// Should parse but with no hunks
	assert.Equal(t, "", diff.OldFile)
	assert.Equal(t, "", diff.NewFile)
	assert.Len(t, diff.Hunks, 0)
}

func TestIsUnifiedDiff_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "empty string",
			input:    "",
			expected: false,
		},
		{
			name:     "only file headers",
			input:    "--- a/file.go\n+++ b/file.go",
			expected: false,
		},
		{
			name:     "only hunk header",
			input:    "@@ -1,3 +1,4 @@",
			expected: false,
		},
		{
			name:     "minimal valid diff",
			input:    "--- a/go\n+++ b/go\n@@ -1 +1 @@",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, IsUnifiedDiff(tt.input))
		})
	}
}

func TestParser_Parse_LongLines(t *testing.T) {
	longContent := strings.Repeat("x", 1000)
	input := "--- a/file.go\n" +
		"+++ b/file.go\n" +
		"@@ -1,1 +1,1 @@\n" +
		"-" + longContent + "\n" +
		"+" + longContent + "\n"

	parser := NewParser()
	diff, err := parser.Parse(input)
	require.NoError(t, err)

	assert.Equal(t, 1, diff.Added)
	assert.Equal(t, 1, diff.Deleted)
	assert.Equal(t, longContent, diff.Hunks[0].Lines[0].Content)
	assert.Equal(t, longContent, diff.Hunks[0].Lines[1].Content)
}

func TestParser_Parse_SpecialCharacters(t *testing.T) {
	// Test parsing of special characters in diff content
	// Note: newline and carriage return in the content will be consumed by string splitting
	// This test verifies that tab characters and other special chars are handled correctly
	input := "--- a/file.go\n" +
		"+++ b/file.go\n" +
		"@@ -1,1 +1,1 @@\n" +
		"-\tfoo\n" +
		"+\tbar\n"

	parser := NewParser()
	diff, err := parser.Parse(input)
	require.NoError(t, err)

	// Tab characters should be preserved in content
	assert.Equal(t, "\tfoo", diff.Hunks[0].Lines[0].Content)
	assert.Equal(t, "\tbar", diff.Hunks[0].Lines[1].Content)
}

func TestParser_Parse_UnicodeCharacters(t *testing.T) {
	input := "--- a/file.go\n" +
		"+++ b/file.go\n" +
		"@@ -1,1 +1,1 @@\n" +
		"-Hello 世界\n" +
		"+Hello 🌍\n"

	parser := NewParser()
	diff, err := parser.Parse(input)
	require.NoError(t, err)

	assert.Equal(t, "Hello 世界", diff.Hunks[0].Lines[0].Content)
	assert.Equal(t, "Hello 🌍", diff.Hunks[0].Lines[1].Content)
}
