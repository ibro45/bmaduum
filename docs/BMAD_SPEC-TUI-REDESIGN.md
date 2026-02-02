# BMAD Specification: TUI Redesign - "Claude Theater Mode"

**Document Version:** 1.0
**Date:** 2026-02-02
**Author:** AI Assistant
**Status:** Draft for Review

---

## 1. Executive Summary

### 1.1 Goal
Redesign bmaduum's output layer to provide a **read-only "mission control" dashboard** that authentically replicates Claude Code's visual interface while maintaining full programmatic control and automatic workflow progression.

### 1.2 Key Value Propositions
- **Visual Familiarity**: Users see exactly what they'd see in Claude Code
- **Zero Interaction Friction**: Fully automatic step progression
- **Workflow Visibility**: Always-visible header showing current step/progress
- **Authentic Experience**: Matching Claude's colors, symbols, animations, and behavior

### 1.3 Success Criteria
- [ ] Visual output is indistinguishable from Claude Code CLI (colors, symbols, layout)
- [ ] Token-by-token streaming text animation
- [ ] "Thinking" spinner during processing gaps
- [ ] Always-auto-scroll (content streams to bottom automatically)
- [ ] Zero keyboard shortcuts except Ctrl+C to quit
- [ ] Zero regression in automation (steps proceed automatically)
- [ ] Existing tests pass; new TUI components have 80%+ coverage

---

## 2. Visual Design Specification

### 2.1 Layout Structure

```
┌──────────────────────────────────────────────────────────────────────┐
│ ⚡ bmaduum │ Step 2/4: dev-story │ PROJ-123 │ claude-4 │ ⏱️ 02:34   │  <- Header (1 line)
├──────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  I'll implement the JWT authentication system. Let me start by      │  <- Text streams
│  examining the current project structure...                         │     token-by-token
│                                                                      │
│  ⏺ Bash(ls -la src/)                                                │  <- Tool use (⏺ symbol)
│  ⎿  total 64                                                         │  <- Output (⎿ symbol)
│     drwxr-xr-x  10 user staff   320 Jan  1 00:00 .                  │
│     drwxr-xr-x   5 user staff   160 Jan  1 00:00 ..                 │
│     -rw-r--r--   1 user staff  2048 Jan  1 00:00 index.ts           │  <- All content visible
│                                                                      │
│  ⏺ Read(src/auth/types.ts)                                          │  <- Tool use
│  ⎿  export interface AuthConfig {                                   │  <- Output with
│       tokenExpiry: number;                                          │     syntax highlighting
│       refreshToken: string;                                         │
│     }                                                               │
│                                                                      │
│  [Thinking spinner] Analyzing the codebase structure...             │  <- Processing indicator
│                                                                      │
│  ✓ Session complete                                                 │  <- Completion
│                                                                      │
│  ── Step 2/4 complete ──                                            │  <- Step transition
│                                                                      │
│  ⏺ Bash(npm test)                                                   │  <- Next step begins
│                                                                      │
│  [Content scrolls smoothly...]                                      │
│                                                                      │
└──────────────────────────────────────────────────────────────────────┘
```

**Scrolling**: Mouse wheel/trackpad scroll naturally. Always auto-scrolls to bottom as new content arrives.

### 2.2 Color Palette (Claude Code Authentic)

```go
var ClaudeColors = struct {
    // Primary brand
    Primary     string // #6B4EE6 (Anthropic purple)
    PrimaryDim  string // #5A3FD4

    // Tool display
    ToolIcon    string // #58A6FF (Blue) - for ⏺
    OutputIcon  string // #8B949E (Gray) - for ⎿

    // Text
    Text        string // #E6EDF3 (Off-white)
    TextMuted   string // #8B949E (Gray)
    TextDim     string // #6E7681 (Dark gray)

    // Semantic
    Success     string // #3FB950 (Green)
    Error       string // #F85149 (Red)
    Warning     string // #D29922 (Orange)
    Info        string // #58A6FF (Blue)

    // Background/structure
    Background  string // Terminal default (transparent)
    Border      string // #30363D (Dark border)
    Selection   string // #264F78 (Selection blue)

    // Syntax highlighting (for code blocks)
    Comment     string // #8B949E
    Keyword     string // #FF7B72
    String      string // #A5D6FF
    Function    string // #D2A8FF
    Number      string // #79C0FF
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
```

### 2.3 Typography & Symbols

| Element | Symbol | Unicode | Usage |
|---------|--------|---------|-------|
| Tool invocation | ⏺ | U+23FA | Before tool name |
| Tool output | ⎿ | U+23BF | Before each output line |
| Success | ✓ | U+2713 | Step/session complete |
| Error | ✗ | U+2717 | Step/session failed |
| Spinner | ⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏ | Braille | Thinking indicator |
| Ellipsis | … | U+2026 | Indicator text |
| Header | ⚡ | U+26A1 | bmaduum logo |
| Clock | ⏱️ | U+23F1 | Timer |
| New content | ↓ | U+2193 | Scroll indicator |

### 2.4 Component Specifications

#### 2.4.1 Header Bar

**Height:** 1 line
**Style:** Inverse/reverse video or distinct background
**Content:**
```
[Logo] bmaduum │ [Step indicator] Step 2/4: dev-story │ [Story] PROJ-123 │ [Model] claude-4 │ [Timer] ⏱️ 02:34
```

**Layout:**
- Left: Logo + "bmaduum" brand
- Center-left: Step progress (current/total + workflow name)
- Center: Story key
- Center-right: Model name
- Right: Elapsed timer

**Color:** Background #6B4EE6 (primary), foreground white

#### 2.4.2 Tool Use Display

**Format:**
```
⏺ ToolName(parameters)
```

**Style:**
- ⏺ symbol: Blue (#58A6FF), bold
- ToolName: White, bold
- Parameters: Gray (#8B949E), normal weight

**Example:**
```
⏺ Bash(ls -la src/)
⏺ Read(src/auth/types.ts)
⏺ Edit(src/config.ts)
```

#### 2.4.3 Tool Output Display

**Format:**
```
⎿  [output line 1]
   [output line 2]
   [output line 3]
   [output line 4]
   [... all remaining lines visible ...]
```

**Style:**
- ⎿ symbol: Gray (#8B949E)
- First line: Aligned with symbol
- Subsequent lines: Indented 3 spaces
- **All content visible** - no truncation (scroll to see earlier output)

**Behavior:**
- Show all output lines from Claude
- Mouse wheel to scroll up and see earlier content

#### 2.4.4 Text Streaming Display

**Format:**
```
[Claude's text appears here, streaming character by character]
```

**Style:**
- Normal text color (#E6EDF3)
- Paragraphs separated by blank line
- No special prefix symbols

**Animation:**
- Characters appear one at a time
- Delay: ~5ms per character (adjustable)
- Batch: 3-5 characters per frame for efficiency
- Pause at punctuation (.,!?) for 100ms for natural feel

#### 2.4.5 Thinking/Processing Indicator

**Format:**
```
[spinner] [status text]
```

**Style:**
- Spinner: Braille pattern, rotating
- Text: Muted gray (#8B949E), italic

**States:**
- "Claude is thinking..."
- "Executing Bash(...)"
- "Waiting for tool result..."
- "Processing..."

**Timing:**
- Appear after 500ms of no output
- Disappear immediately when output resumes

#### 2.4.6 Content Display (Auto-Expand)

**Decision:** All content is displayed in full (no collapsing)

**Rationale:** Simplified UX - users can scroll with mouse to see everything

**Implementation:**
- Tool output: Show all lines
- No truncation (beyond what Claude already does)
- No "... +N lines" indicators
- Scroll with mouse wheel to see earlier output

#### 2.4.7 Scroll Status (Footer)

**Format:**
```
⚡ bmaduum │ Step 2/4 │ PROJ-123 │ ⏱️ 02:34
```

**Note:** No footer needed - all status is in the sticky header. Mouse wheel to scroll.

---

## 3. Architecture Specification

### 3.1 Package Structure

```
internal/
├── tui/                          # NEW: Bubble Tea TUI implementation
│   ├── model.go                  # Main Bubble Tea model
│   ├── update.go                 # Update function (event handling)
│   ├── view.go                   # View function (rendering)
│   ├── init.go                   # Initialization and commands
│   │
│   ├── components/               # UI components
│   │   ├── header.go             # Status header bar
│   │   ├── output_section.go     # Collapsible output section
│   │   ├── tool_use.go           # Tool invocation display
│   │   ├── text_stream.go        # Streaming text with typewriter effect
│   │   ├── spinner.go            # Thinking/processing indicator
│   │   └── viewport.go           # Viewport with auto-scroll
│   │
│   ├── styles/                   # Visual styling
│   │   ├── colors.go             # Claude color palette
│   │   ├── symbols.go            # Unicode symbols
│   │   └── lipgloss.go           # Pre-defined lipgloss styles
│   │
│   └── events/                   # Custom TUI events
│       ├── content.go            # New content arrived
│       ├── step.go               # Step transition events
│       └── scroll.go             # Scroll state changes
│
├── claude/                       # EXISTING: Modified for streaming
│   ├── types.go                  # ADD: Token-level streaming support
│   ├── client.go                 # MODIFY: Support token callbacks
│   └── parser.go                 # MODIFY: Emit partial text events
│
├── output/                       # EXISTING: Extended
│   ├── printer.go                # MODIFY: Interface for TUI integration
│   └── claude_printer.go         # NEW: Claude-style formatting
│
└── workflow/                     # EXISTING: Unchanged interface
    └── workflow.go               # Runner connects to new TUI
```

### 3.2 Component Details

#### 3.2.1 TUIModel (Main Model)

```go
type TUIModel struct {
    // Header state
    CurrentStep    int
    TotalSteps     int
    StepName       string
    StoryKey       string
    Model          string
    StartTime      time.Time

    // Content state
    Sections       []OutputSection    // All output sections
    CurrentSection *OutputSection     // Currently being built

    // Viewport
    Viewport       viewport.Model

    // Animation
    Typewriter     TypewriterState    // For character-by-character text
    Spinner        spinner.Model      // Thinking indicator
    Thinking       bool               // Currently showing spinner
    LastActivity   time.Time          // For thinking detection

    // Runtime
    Executor       claude.Executor
    Config         *config.Config
    Width          int
    Height         int
    Err            error
    Quitting       bool
}

type OutputSection struct {
    ID          string
    Type        SectionType      // Text, ToolUse, ToolResult
    Content     string           // Raw content
    Lines       []string         // Split lines
    Rendered    string           // Final rendered string
    Language    string           // For syntax highlighting (code blocks)
}

type TypewriterState struct {
    Buffer      string           // Full text to display
    Displayed   int              // Characters already shown
    Pending     []rune           // Characters waiting to be shown
    Speed       time.Duration    // Delay between characters
}
```

#### 3.2.2 Event Flow

```
┌─────────────────────────────────────────────────────────────────┐
│                     Event Flow Diagram                          │
└─────────────────────────────────────────────────────────────────┘

claude.Executor
       │
       │ stream-json events
       ▼
┌──────────────┐     ┌──────────────┐     ┌─────────────────────┐
│ JSON Parser  │────▶│ Event Router │────▶│ TUI Update Function │
└──────────────┘     └──────────────┘     └─────────────────────┘
                                                  │
       ┌──────────────────────────────────────────┼──────────┐
       │                                          │          │
       ▼                                          ▼          ▼
┌─────────────┐  ┌──────────────┐  ┌──────────┐  ┌──────┐  ┌─────────┐
│ Text Event  │  │ ToolUseEvent │  │ ToolResult│  │ Init │  │ Complete│
└─────────────┘  └──────────────┘  └──────────┘  └──────┘  └─────────┘
       │                                          │          │
       ▼                                          ▼          ▼
┌─────────────────────┐  ┌────────────────┐  ┌──────────────────────┐
│ Typewriter Animation│  │ Create Section │  │ Step Transition      │
│ (character by char) │  │ (auto-expand)  │  │ (auto-proceed)       │
└─────────────────────┘  └────────────────┘  └──────────────────────┘
       │                                          │          │
       ▼                                          ▼          ▼
┌─────────────────────┐  ┌────────────────┐  ┌──────────────────────┐
│ Viewport.SetContent │  │ Viewport.Update│  │ Next Step / Quit     │
│ GotoBottom (always) │  │ (always follow)│  │                      │
└─────────────────────┘  └────────────────┘  └──────────────────────┘
```

#### 3.2.3 Key Handler (Simplified)

| Key | Action |
|-----|--------|
| `ctrl+c` | Quit application (only shortcut) |

**Scrolling**: Mouse wheel/trackpad only (no keyboard scrolling)

#### 3.2.4 Auto-Scroll Logic (Always On)

```go
func (m *TUIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        // Only handle quit
        switch msg.String() {
        case "ctrl+c":
            return m, tea.Quit
        }

    case tea.MouseMsg:
        // Handle mouse scrolling (forward to viewport)
        var cmd tea.Cmd
        m.Viewport, cmd = m.Viewport.Update(msg)
        return m, cmd

    case ContentMsg:
        // New content arrived - always scroll to bottom
        m.addContent(msg.Content)
        m.Viewport.GotoBottom()
        return m, nil
    }
}
```

### 3.3 Integration Points

#### 3.3.1 Modified: claude.Event

Add token-level streaming support:

```go
// EventTypeToken for character-by-character streaming
type EventType string

const (
    EventTypeSystem      EventType = "system"
    EventTypeAssistant   EventType = "assistant"
    EventTypeToolUse     EventType = "tool_use"
    EventTypeToolResult  EventType = "tool_result"
    EventTypeToken       EventType = "token"        // NEW: Individual token
    EventTypeResult      EventType = "result"
)

type Event struct {
    Type            EventType

    // Existing fields
    SessionStarted  bool
    SessionComplete bool
    Text            string
    ToolName        string
    ToolInput       map[string]interface{}
    ToolStdout      string
    ToolStderr      string

    // NEW: Token streaming
    Token           string     // Single character/token
    IsCompleteText  bool       // This event contains the complete text
}
```

#### 3.3.2 Modified: workflow.Runner

Integrate with TUI instead of Printer:

```go
func (r *Runner) RunSingleTUI(ctx context.Context, workflowName, storyKey string, tuiModel *tui.TUIModel) int {
    prompt, err := r.config.GetPrompt(workflowName, storyKey)
    if err != nil {
        return 1
    }

    model := r.config.GetModel(workflowName)

    // Create TUI program with our model
    p := tea.NewProgram(tuiModel, tea.WithAltScreen(), tea.WithMouseCellMotion())

    // Run Claude in background, feeding events to TUI
    go func() {
        handler := func(event claude.Event) {
            p.Send(tui.ClaudeEventMsg{Event: event})
        }

        exitCode, _ := r.executor.ExecuteWithResult(ctx, prompt, handler, model)
        p.Send(tui.CompleteMsg{ExitCode: exitCode})
    }()

    // Run TUI (blocks until complete)
    finalModel, err := p.Run()
    if err != nil {
        return 1
    }

    return finalModel.(*TUIModel).ExitCode
}
```

---

## 4. Implementation Phases

### Phase 1: Foundation (Week 1)

**Goal:** Basic TUI shell with header and viewport

**Deliverables:**
- [ ] Create `internal/tui/` package structure
- [ ] Implement `TUIModel` with basic Bubble Tea integration
- [ ] Build header component with step/story/timer display
- [ ] Integrate `viewport` bubble for scrollable content
- [ ] Basic keyboard handling (quit only)

**Testing:**
- Unit tests for header rendering
- Unit tests for viewport integration
- Manual test: `go run ./cmd/bmaduum` with mock events

**Files Created/Modified:**
- `internal/tui/model.go`
- `internal/tui/view.go`
- `internal/tui/update.go`
- `internal/tui/components/header.go`
- `internal/tui/styles/colors.go`
- `internal/tui/styles/symbols.go`

### Phase 2: Content Display (Week 1-2)

**Goal:** Display Claude events in authentic Claude Code style

**Deliverables:**
- [ ] `OutputSection` component with type detection
- [ ] Tool use display (⏺ symbol, proper formatting)
- [ ] Tool result display (⎿ symbol, line indentation)
- [ ] Text display (paragraph formatting)
- [ ] Syntax highlighting integration (Chroma)

**Testing:**
- Component tests for each section type
- Visual regression tests (if feasible)
- Test with real Claude events captured to file

**Files Created:**
- `internal/tui/components/output_section.go`
- `internal/tui/components/tool_use.go`
- `internal/tui/components/text_stream.go`
- `internal/output/claude_printer.go`

### Phase 3: Animations (Week 2)

**Goal:** Typewriter effect and thinking spinner

**Deliverables:**
- [ ] Typewriter animation (character-by-character streaming)
- [ ] Thinking spinner with gap detection
- [ ] All content auto-expands (no collapsing)

**Testing:**
- Animation timing tests
- Thinking state detection tests

**Files Created:**
- `internal/tui/components/text_stream.go` (extended)
- `internal/tui/components/spinner.go`

### Phase 4: Scrolling (Week 2-3)

**Goal:** Mouse-based scrolling with auto-follow

**Deliverables:**
- [ ] Auto-scroll to bottom on new content (always on)
- [ ] Mouse wheel/trackpad scrolling support

**Testing:**
- Viewport integration tests

**Files Created:**
- `internal/tui/components/viewport.go` (wrapper)

### Phase 5: Integration & Polish (Week 3)

**Goal:** Wire into existing workflow system, test end-to-end

**Deliverables:**
- [ ] Modify `workflow.Runner` to support TUI mode
- [ ] CLI flag `--tui` to enable new interface
- [ ] Graceful fallback to old printer if TUI fails
- [ ] State persistence integration (save/restore scroll position)
- [ ] Performance optimization (large output handling)

**Testing:**
- End-to-end tests with mock executor
- Integration tests with real Claude CLI (optional, manual)
- Performance tests (10k+ lines of output)
- Cross-platform testing (macOS, Linux)

**Files Modified:**
- `internal/cli/run.go`
- `internal/cli/queue.go`
- `internal/cli/epic.go`
- `internal/workflow/workflow.go`

### Phase 6: Documentation & Release (Week 4)

**Goal:** Document, demo, release

**Deliverables:**
- [ ] User documentation (TUI controls, features)
- [ ] Architecture documentation
- [ ] Demo video/GIF
- [ ] Blog post (optional)
- [ ] Release notes

---

## 5. Testing Strategy

### 5.1 Unit Testing

**Components to Test:**
- Header rendering with various states
- Output section rendering
- Typewriter animation state machine
- Color/style application

**Example:**
```go
func TestOutputSection_Render(t *testing.T) {
    section := OutputSection{
        Type:    SectionToolResult,
        Content: "line1\nline2\nline3\nline4\nline5",
        Lines:   []string{"line1", "line2", "line3", "line4", "line5"},
    }

    view := section.View(80)

    assert.Contains(t, view, "line1")
    assert.Contains(t, view, "line5")
    assert.Contains(t, view, "⎿") // Output symbol
}
```

### 5.2 Integration Testing

**Approach:**
- Use `claude.MockExecutor` with recorded event sequences
- Capture real Claude output to fixture files
- Replay through TUI and verify rendering

**Fixtures:**
- `test/fixtures/simple_text.jsonl` - Simple text response
- `test/fixtures/tool_use.jsonl` - Tool invocation and result
- `test/fixtures/long_output.jsonl` - Large output for performance testing
- `test/fixtures/multi_step.jsonl` - Multiple workflow steps

### 5.3 Visual Testing

**Approach:**
- Manual comparison screenshots with Claude Code
- Terminal recording (asciinema) for regression testing
- Color accuracy verification

### 5.4 Performance Testing

**Scenarios:**
- 10,000 lines of output (should remain responsive)
- Rapid event streaming (100 events/second)
- Large tool output (100KB stdout)

**Metrics:**
- Frame rate (target: 30fps minimum)
- Memory usage (target: <100MB for 10k lines)
- Startup time (target: <100ms)

---

## 6. Open Questions & Decisions

### Q1: Syntax Highlighting Approach

**Options:**
1. **Chroma** (github.com/alecthomas/chroma) - Native Go, ANSI output
2. **Glamour** (charmbracelet) - Markdown rendering, might be overkill
3. **Manual** - Simple keyword coloring, faster but less accurate

**Recommendation:** Chroma with `terminal256` formatter

### Q2: Mouse Support

**Question:** What mouse support is needed?

**Decision:** Mouse wheel/trackpad scrolling only

**Rationale:** No clicking needed since all content auto-expands. Mouse wheel is universally supported in terminals for TUIs.

### Q3: Step Transition Animation

**Question:** Should step transitions have special visual treatment?

**Options:**
1. Simple divider line
2. Full-screen flash/slide
3. Compact inline transition

**Recommendation:** Simple divider line with step completion summary

### Q4: Configuration Options

**Question:** What should be user-configurable?

**Proposed:**
- `TUI_ENABLED` - Enable/disable TUI (fallback to old printer)
- `TUI_TYPEWRITER_SPEED` - Character delay in ms (default: 5)
- `TUI_TYPEWRITER_SPEED` - Character delay (ms)

### Q5: Error Handling in TUI

**Question:** How should fatal errors be displayed?

**Options:**
1. Stay in TUI, show error in viewport
2. Exit TUI, print error to terminal
3. Modal dialog in TUI

**Recommendation:** Option 1 for non-fatal, Option 2 for fatal errors

---

## 7. Dependencies

### New Dependencies

```bash
# Core TUI framework (already using lipgloss, adding bubbletea)
go get github.com/charmbracelet/bubbletea
go get github.com/charmbracelet/bubbles/viewport
go get github.com/charmbracelet/bubbles/spinner

# Syntax highlighting
go get github.com/alecthomas/chroma/v2
go get github.com/alecthomas/chroma/v2/formatters
go get github.com/alecthomas/chroma/v2/lexers
go get github.com/alecthomas/chroma/v2/styles
```

### Existing Dependencies (No Change)

- `github.com/charmbracelet/lipgloss` - Styling (already used)
- `github.com/spf13/cobra` - CLI commands
- `github.com/spf13/viper` - Configuration

---

## 8. Risks & Mitigations

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| Performance issues with large output | Medium | High | Virtual scrolling, content virtualization |
| Terminal compatibility issues | Medium | Medium | Graceful fallback to old printer |
| Bubble Tea learning curve | Low | Low | Phased implementation, comprehensive tests |
| User prefers old output | Low | Medium | Keep `--no-tui` flag for backward compatibility |
| Stream-json format changes | Low | High | Version detection, graceful degradation |

---

## 9. Success Metrics

| Metric | Target | Measurement |
|--------|--------|-------------|
| Visual similarity to Claude Code | 90%+ | Side-by-side comparison screenshots |
| Test coverage (new code) | 80%+ | `go test -cover` |
| Performance (10k lines) | <100MB memory | Benchmark tests |
| Frame rate | 30fps+ | Visual inspection |
| User satisfaction | Positive | Post-implementation feedback |

---

## 10. Appendix

### A. Claude Code Stream-JSON Format Reference

```json
{"type":"init","session_id":"abc123"}
{"type":"message","role":"assistant","content":[{"type":"text","text":"I'll help"}]}
{"type":"message","role":"assistant","content":[{"type":"tool_use","name":"Bash","input":{"command":"ls"}}]}
{"type":"tool_result","output":"file1\nfile2","error":null}
{"type":"result","result":"Done","status":"success"}
```

### B. Keyboard Shortcuts Reference

| Key | Action | Context |
|-----|--------|---------|
| `ctrl+c` | Quit | Global |

**Scrolling**: Mouse wheel/trackpad only (no keyboard shortcuts)

### C. Color Reference Table

See Section 2.2 for full color palette.

---

**End of Specification**
