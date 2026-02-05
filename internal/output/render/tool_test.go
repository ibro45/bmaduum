package render

import "testing"

func TestStripLineNumberArrows(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "no arrows - plain text",
			input:    "hello world\nfoo bar",
			expected: "hello world\nfoo bar",
		},
		{
			name:     "single line with arrow",
			input:    "     1→package main",
			expected: "1  package main",
		},
		{
			name:     "multiple lines with arrows",
			input:    "     1→package main\n     2→\n     3→import \"fmt\"",
			expected: "1  package main\n2  \n3  import \"fmt\"",
		},
		{
			name:     "mixed lines with and without arrows",
			input:    "     1→first line\nno arrow here\n     2→second line",
			expected: "1  first line\nno arrow here\n2  second line",
		},
		{
			name:     "large line numbers",
			input:    "   100→line one hundred\n   101→line one oh one",
			expected: "100  line one hundred\n101  line one oh one",
		},
		{
			name:     "varying whitespace before numbers",
			input:    "  1→short\n 10→medium\n100→long",
			expected: "  1  short\n 10  medium\n100  long",
		},
		{
			name:     "content with arrow character that isn't line number",
			input:    "x → y means x maps to y",
			expected: "x → y means x maps to y",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stripLineNumberArrows(tt.input)
			if got != tt.expected {
				t.Errorf("stripLineNumberArrows() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestTruncateOutput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		maxLines int
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			maxLines: 5,
			expected: "",
		},
		{
			name:     "within limit",
			input:    "line1\nline2\nline3",
			maxLines: 5,
			expected: "line1\nline2\nline3",
		},
		{
			name:     "exactly at limit",
			input:    "line1\nline2\nline3",
			maxLines: 3,
			expected: "line1\nline2\nline3",
		},
		{
			name:     "exceeds limit",
			input:    "line1\nline2\nline3\nline4\nline5",
			maxLines: 3,
			expected: "line1\nline2\nline3\n… +2 lines",
		},
		{
			name:     "zero maxLines returns all",
			input:    "line1\nline2\nline3",
			maxLines: 0,
			expected: "line1\nline2\nline3",
		},
		{
			name:     "negative maxLines returns all",
			input:    "line1\nline2\nline3",
			maxLines: -1,
			expected: "line1\nline2\nline3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := truncateOutput(tt.input, tt.maxLines)
			if got != tt.expected {
				t.Errorf("truncateOutput() = %q, want %q", got, tt.expected)
			}
		})
	}
}
