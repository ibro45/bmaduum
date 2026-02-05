package progress

import "testing"

func TestShortenModel(t *testing.T) {
	tests := []struct {
		name     string
		model    string
		expected string
	}{
		{
			name:     "empty string",
			model:    "",
			expected: "",
		},
		{
			name:     "opus with date suffix",
			model:    "claude-opus-4-5-20251101",
			expected: "opus-4-5",
		},
		{
			name:     "sonnet with date suffix",
			model:    "claude-sonnet-4-5-20250929",
			expected: "sonnet-4-5",
		},
		{
			name:     "simple opus",
			model:    "claude-opus-4-5",
			expected: "opus-4-5",
		},
		{
			name:     "simple sonnet",
			model:    "claude-sonnet-4-5",
			expected: "sonnet-4-5",
		},
		{
			name:     "haiku",
			model:    "claude-haiku-3-5",
			expected: "haiku-3-5",
		},
		{
			name:     "no prefix",
			model:    "opus-4-5",
			expected: "opus-4-5",
		},
		{
			name:     "unknown model",
			model:    "gpt-4-turbo",
			expected: "gpt-4-turbo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shortenModel(tt.model)
			if got != tt.expected {
				t.Errorf("shortenModel(%q) = %q, want %q", tt.model, got, tt.expected)
			}
		})
	}
}
