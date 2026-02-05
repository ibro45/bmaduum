// Package render provides output rendering components for the CLI.
package render

import (
	"testing"
)

func TestIcons(t *testing.T) {
	tests := []struct {
		name  string
		icon  string
		want  rune // First rune should be a specific unicode character
	}{
		{"Success", IconSuccess, '✓'},
		{"Error", IconError, '✗'},
		{"Pending", IconPending, '○'},
		{"InProgress", IconInProgress, '●'},
		{"Tool", IconTool, '⏺'},
		{"Output", IconOutput, '⎿'},
		{"Bmaduum", IconBmaduum, '⚡'},
		{"Clock", IconClock, '⏱'},
		{"Paused", IconPaused, '⏸'},
		{"Retrying", IconRetrying, '↻'},
		{"Separator", IconSeparator, '│'},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.icon) == 0 {
				t.Errorf("icon %s is empty", tt.name)
			}
			if r := []rune(tt.icon)[0]; r != tt.want {
				t.Errorf("icon %s = %c, want %c", tt.name, r, tt.want)
			}
		})
	}
}

func TestIcons_NonEmpty(t *testing.T) {
	icons := []string{
		IconSuccess,
		IconError,
		IconPending,
		IconInProgress,
		IconTool,
		IconOutput,
		IconBmaduum,
		IconClock,
		IconPaused,
		IconRetrying,
		IconSeparator,
	}

	for _, icon := range icons {
		if len(icon) == 0 {
			t.Errorf("icon should not be empty")
		}
	}
}

func TestIcons_UniqueValues(t *testing.T) {
	icons := []string{
		IconSuccess,
		IconError,
		IconPending,
		IconInProgress,
		IconTool,
		IconOutput,
		IconBmaduum,
		IconClock,
		IconPaused,
		IconRetrying,
		IconSeparator,
	}

	seen := make(map[string]bool)
	for _, icon := range icons {
		if seen[icon] {
			t.Errorf("duplicate icon found: %s", icon)
		}
		seen[icon] = true
	}
}
