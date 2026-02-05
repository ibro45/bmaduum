// Package progress provides terminal progress display functionality.
package progress

import (
	"testing"
)

func TestIcons(t *testing.T) {
	tests := []struct {
		name  string
		icon  string
		want  rune // First rune should be a specific unicode character
	}{
		{"success", iconSuccess, '✓'},
		{"error", iconError, '✗'},
		{"pending", iconPending, '○'},
		{"in progress", iconInProgress, '●'},
		{"tool", iconTool, '⏺'},
		{"output", iconOutput, '⎿'},
		{"bmaduum", iconBmaduum, '⚡'},
		{"clock", iconClock, '⏱'},
		{"paused", iconPaused, '⏸'},
		{"retrying", iconRetrying, '↻'},
		{"separator", iconSeparator, '│'},
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

func TestSpinnerFrames(t *testing.T) {
	if len(spinnerFrames) == 0 {
		t.Error("spinnerFrames should not be empty")
	}
	if spinnerFrames[0] != "·" {
		t.Errorf("spinnerFrames[0] = %s, want ·", spinnerFrames[0])
	}
}

func TestThinkingVerbs(t *testing.T) {
	if len(thinkingVerbs) == 0 {
		t.Error("thinkingVerbs should not be empty")
	}
	// Check that we have the expected easter egg
	found := false
	for _, v := range thinkingVerbs {
		if v == "Clauding" {
			found = true
			break
		}
	}
	if !found {
		t.Error("thinkingVerbs should contain 'Clauding' easter egg")
	}
}

func TestPastTenseVerbs(t *testing.T) {
	if len(pastTenseVerbs) == 0 {
		t.Error("pastTenseVerbs should not be empty")
	}
	// Check that past tense verbs match thinking verbs count
	if len(pastTenseVerbs) != len(thinkingVerbs) {
		t.Errorf("pastTenseVerbs length = %d, want %d", len(pastTenseVerbs), len(thinkingVerbs))
	}
	// Check that we have the expected easter egg
	found := false
	for _, v := range pastTenseVerbs {
		if v == "Clauded" {
			found = true
			break
		}
	}
	if !found {
		t.Error("pastTenseVerbs should contain 'Clauded' easter egg")
	}
}

func TestGetThinkingVerb_BoundsChecking(t *testing.T) {
	tests := []struct {
		name string
		idx  int
		want string
	}{
		{"valid index", 0, thinkingVerbs[0]},
		{"last index", len(thinkingVerbs) - 1, thinkingVerbs[len(thinkingVerbs)-1]},
		{"negative", -1, thinkingVerbs[0]},
		{"out of bounds", len(thinkingVerbs), thinkingVerbs[0]},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getThinkingVerb(tt.idx); got != tt.want {
				t.Errorf("getThinkingVerb(%d) = %v, want %v", tt.idx, got, tt.want)
			}
		})
	}
}

func TestGetPastTenseVerb_BoundsChecking(t *testing.T) {
	tests := []struct {
		name string
		idx  int
		want string
	}{
		{"valid index", 0, pastTenseVerbs[0]},
		{"last index", len(pastTenseVerbs) - 1, pastTenseVerbs[len(pastTenseVerbs)-1]},
		{"negative", -1, pastTenseVerbs[0]},
		{"out of bounds", len(pastTenseVerbs), pastTenseVerbs[0]},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getPastTenseVerb(tt.idx); got != tt.want {
				t.Errorf("getPastTenseVerb(%d) = %v, want %v", tt.idx, got, tt.want)
			}
		})
	}
}

func TestVerbPairing(t *testing.T) {
	// Verify that thinking verbs and past tense verbs are properly paired
	if len(thinkingVerbs) != len(pastTenseVerbs) {
		t.Fatalf("thinkingVerbs count (%d) != pastTenseVerbs count (%d)", len(thinkingVerbs), len(pastTenseVerbs))
	}

	// Check specific pairings
	pairings := []struct {
		present string
		past    string
	}{
		{"Thinking", "Thought"},
		{"Pondering", "Pondered"},
		{"Analyzing", "Analyzed"},
		{"Clauding", "Clauded"},
	}

	for _, p := range pairings {
		found := false
		for i := 0; i < len(thinkingVerbs); i++ {
			if thinkingVerbs[i] == p.present && pastTenseVerbs[i] == p.past {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("verb pairing not found: %s -> %s", p.present, p.past)
		}
	}
}
