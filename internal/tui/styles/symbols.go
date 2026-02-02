package styles

// Symbols provides Unicode characters used in the Claude Code interface.
var Symbols = struct {
	ToolInvocation string // ⏺ U+23FA - Tool invocation marker
	ToolOutput     string // ⎿ U+23BF - Tool output marker
	Success        string // ✓ U+2713 - Success/completion marker
	Error          string // ✗ U+2717 - Error marker
	Ellipsis       string // … U+2026 - Ellipsis indicator
	HeaderLogo     string // ⚡ U+26A1 - bmaduum logo
	Clock          string // ⏱️ U+23F1 - Timer icon
	NewContent     string // ↓ U+2193 - Scroll indicator
	Bullet         string // • U+2022 - Bullet point
	Divider        string // ─ U+2500 - Horizontal line
}{
	ToolInvocation: "⏺",
	ToolOutput:     "⎿",
	Success:        "✓",
	Error:          "✗",
	Ellipsis:       "…",
	HeaderLogo:     "⚡",
	Clock:          "⏱️",
	NewContent:     "↓",
	Bullet:         "•",
	Divider:        "─",
}

// SpinnerFrames provides the Braille pattern for the thinking spinner.
var SpinnerFrames = []string{
	"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏",
}
