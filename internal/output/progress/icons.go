// Package progress provides terminal progress display functionality.
package progress

import "bmaduum/internal/output/render"

// Icons for status indicators in terminal output.
// Re-exported from render package for convenience.
const (
	iconSuccess    = render.IconSuccess    // Completed successfully
	iconError      = render.IconError      // Failed
	iconPending    = render.IconPending    // Not yet started
	iconInProgress = render.IconInProgress // Currently running
	iconTool       = render.IconTool       // Tool invocation (filled circle)
	iconOutput     = render.IconOutput     // Tool output (right angle bracket)
	iconBmaduum    = render.IconBmaduum    // Bmaduum logo (high voltage)
	iconClock      = render.IconClock      // Clock/timer (no variation selector for consistency)
	iconPaused     = render.IconPaused     // Paused rate limit
	iconRetrying   = render.IconRetrying   // Retrying (simpler arrow)
	iconSeparator  = render.IconSeparator  // Vertical separator
)

// Spinner frames for animated thinking indicator.
// Claude Code uses a simple static dot, no animation.
var spinnerFrames = []string{"·"}

// Thinking verbs for Claude Code-style status messages.
// Randomly selected during "thinking" states for whimsy.
var thinkingVerbs = []string{
	"Thinking",
	"Pondering",
	"Analyzing",
	"Processing",
	"Cogitating",
	"Contemplating",
	"Reasoning",
	"Evaluating",
	"Computing",
	"Deliberating",
	"Mulling",
	"Considering",
	"Musing",
	"Ruminating",
	"Synthesizing",
	"Clauding", // Easter egg
}

// Past-tense versions of thinking verbs for completion messages.
// Must match the order of thinkingVerbs for correct pairing.
var pastTenseVerbs = []string{
	"Thought",
	"Pondered",
	"Analyzed",
	"Processed",
	"Cogitated",
	"Contemplated",
	"Reasoned",
	"Evaluated",
	"Computed",
	"Deliberated",
	"Mulled",
	"Considered",
	"Mused",
	"Ruminated",
	"Synthesized",
	"Clauded", // Easter egg past tense
}

// getThinkingVerb returns a thinking verb by index (with bounds checking).
func getThinkingVerb(idx int) string {
	if idx < 0 || idx >= len(thinkingVerbs) {
		return thinkingVerbs[0]
	}
	return thinkingVerbs[idx]
}

// getPastTenseVerb returns a past-tense verb by index (with bounds checking).
func getPastTenseVerb(idx int) string {
	if idx < 0 || idx >= len(pastTenseVerbs) {
		return pastTenseVerbs[0]
	}
	return pastTenseVerbs[idx]
}
