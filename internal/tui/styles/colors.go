package styles

// ClaudeColors provides the authentic Claude Code color palette.
// These colors match the visual appearance of the official Claude CLI.
var ClaudeColors = struct {
	// Primary brand
	Primary    string // #6B4EE6 (Anthropic purple)
	PrimaryDim string // #5A3FD4

	// Tool display
	ToolIcon   string // #58A6FF (Blue) - for ⏺
	OutputIcon string // #8B949E (Gray) - for ⎿

	// Text
	Text      string // #E6EDF3 (Off-white)
	TextMuted string // #8B949E (Gray)
	TextDim   string // #6E7681 (Dark gray)

	// Semantic
	Success string // #3FB950 (Green)
	Error   string // #F85149 (Red)
	Warning string // #D29922 (Orange)
	Info    string // #58A6FF (Blue)

	// Background/structure
	Background string // Terminal default (transparent)
	Border     string // #30363D (Dark border)
	Selection  string // #264F78 (Selection blue)

	// Syntax highlighting (for code blocks)
	Comment  string // #8B949E
	Keyword  string // #FF7B72
	String   string // #A5D6FF
	Function string // #D2A8FF
	Number   string // #79C0FF
}{
	Primary:    "#6B4EE6",
	PrimaryDim: "#5A3FD4",
	ToolIcon:   "#58A6FF",
	OutputIcon: "#8B949E",
	Text:       "#E6EDF3",
	TextMuted:  "#8B949E",
	TextDim:    "#6E7681",
	Success:    "#3FB950",
	Error:      "#F85149",
	Warning:    "#D29922",
	Info:       "#58A6FF",
	Border:     "#30363D",
	Selection:  "#264F78",
	Comment:    "#8B949E",
	Keyword:    "#FF7B72",
	String:     "#A5D6FF",
	Function:   "#D2A8FF",
	Number:     "#79C0FF",
}
