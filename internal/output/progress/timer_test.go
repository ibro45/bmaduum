// Package progress provides terminal progress display functionality.
package progress

import (
	"testing"
	"time"
)

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		want     string
	}{
		{"zero", 0, "0:00"},
		{"seconds", 30 * time.Second, "0:30"},
		{"minutes", 5 * time.Minute, "5:00"},
		{"minutes seconds", 5*time.Minute + 30*time.Second, "5:30"},
		{"hour", 1 * time.Hour, "1:00:00"},
		{"hour minutes", 1*time.Hour + 30*time.Minute, "1:30:00"},
		{"hour minutes seconds", 1*time.Hour + 30*time.Minute + 15*time.Second, "1:30:15"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatDuration(tt.duration); got != tt.want {
				t.Errorf("formatDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatDurationNatural(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		want     string
	}{
		{"zero", 0, "0s"},
		{"seconds", 30 * time.Second, "30s"},
		{"minutes", 5 * time.Minute, "5m"},
		{"minutes seconds", 5*time.Minute + 30*time.Second, "5m 30s"},
		{"hour", 1 * time.Hour, "1h"},
		{"hour minutes", 1*time.Hour + 30*time.Minute, "1h 30m"},
		{"hour minutes seconds", 1*time.Hour + 30*time.Minute + 15*time.Second, "1h 30m"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatDurationNatural(tt.duration); got != tt.want {
				t.Errorf("formatDurationNatural() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatTokenCount(t *testing.T) {
	tests := []struct {
		name  string
		count int
		want  string
	}{
		{"zero", 0, "0"},
		{"hundreds", 500, "500"},
		{"thousands", 1500, "1.5k"},
		{"ten thousands", 15700, "15.7k"},
		{"hundred thousands", 157000, "157.0k"},
		{"millions", 1500000, "1.50M"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatTokenCount(tt.count); got != tt.want {
				t.Errorf("formatTokenCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDisplayWidth(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want int
	}{
		{"ascii", "hello", 5},
		{"emoji", "⚡", 2},
		{"mixed", "⚡hello", 7},
		{"wide", "你好", 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := displayWidth(tt.s); got != tt.want {
				t.Errorf("displayWidth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTruncateToDisplayWidth(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		maxWidth int
		want     string
	}{
		{"no truncation", "hello", 10, "hello"},
		{"exact fit", "hello", 5, "hello"},
		{"truncate", "hello world", 8, "hello w…"},
		{"emoji truncation", "⚡⚡⚡", 4, "⚡…"},
		{"empty", "", 5, ""},
		{"zero width", "hello", 0, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := truncateToDisplayWidth(tt.s, tt.maxWidth); got != tt.want {
				t.Errorf("truncateToDisplayWidth() = %v, want %v", got, tt.want)
			}
		})
	}
}
