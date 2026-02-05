// Package progress provides terminal progress display functionality.
package progress

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"bmaduum/internal/output/terminal"
)

// Line manages the fixed status area at the bottom of the terminal.
// It provides a two-line status display with scrolling output above:
//   - Activity line (second-to-last): Spinner + activity verb + timer + token count
//   - Status bar (last row): Operation + step info + story key + model + total timer
type Line struct {
	// Terminal control
	term      *terminal.Terminal
	cursor    *terminal.Cursor
	activity  *ActivityLine
	statusBar *StatusBar

	// State management
	mu       sync.Mutex
	state    State
	enabled  bool
	initOnce sync.Once

	// Lifecycle
	ctx            context.Context
	cancel         context.CancelFunc
	resizeChan     chan os.Signal
	verbChangeTick int
}

// NewLine creates a new progress Line instance.
func NewLine(out io.Writer) *Line {
	term := terminal.New(out)

	return &Line{
		term:      term,
		cursor:    terminal.NewCursor(out),
		activity:  NewActivityLine(out),
		statusBar: NewStatusBar(out),
		enabled:   term.FileDescriptor() >= 0 && terminal.IsTTY(out.(*os.File)),
	}
}

// Init initializes the fixed status area at the bottom of the terminal.
func (l *Line) Init() {
	if !l.enabled {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	l.initOnce.Do(func() {
		_, height := l.term.Size()
		if height < 3 {
			l.enabled = false
			return
		}

		// Set scrolling region from row 1 to (height-2)
		// This reserves the bottom 2 rows for the status area (activity line + status bar)
		if err := l.term.SetScrollRegion(1, height-2); err != nil {
			l.enabled = false
			return
		}

		// Create context for goroutines
		l.ctx, l.cancel = context.WithCancel(context.Background())

		// Pick random starting verb
		l.state.VerbIdx = rand.Intn(len(thinkingVerbs))

		// Initialize start time
		l.state.StartTime = time.Now()

		// Start animation ticker
		go l.animationTicker()

		// Handle terminal resize
		go l.handleResize()

		// Initial render
		l.render()
	})
}

// animationTicker updates the spinner animation every 100ms.
func (l *Line) animationTicker() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-l.ctx.Done():
			return
		case <-ticker.C:
			l.mu.Lock()
			l.state.SpinnerIdx = (l.state.SpinnerIdx + 1) % len(spinnerFrames)
			l.verbChangeTick++
			if l.verbChangeTick >= 30 { // Change verb every ~3 seconds
				l.verbChangeTick = 0
				l.state.VerbIdx = rand.Intn(len(thinkingVerbs))
			}
			l.render()
			l.mu.Unlock()
		}
	}
}

// handleResize listens for SIGWINCH and reconfigures the scroll region.
func (l *Line) handleResize() {
	l.resizeChan = make(chan os.Signal, 1)
	signal.Notify(l.resizeChan, syscall.SIGWINCH)

	for {
		select {
		case <-l.ctx.Done():
			signal.Stop(l.resizeChan)
			close(l.resizeChan)
			return
		case <-l.resizeChan:
			// Debounce: wait a bit for resize to settle
			time.Sleep(50 * time.Millisecond)

			// Drain any queued resize signals
			for len(l.resizeChan) > 0 {
				<-l.resizeChan
			}

			l.mu.Lock()
			_, oldHeight := l.term.Size()
			newWidth, newHeight := l.term.UpdateSize()

			if newWidth != 0 || newHeight != oldHeight {
				l.handleResizeInner(oldHeight, newWidth, newHeight)
			}
			l.mu.Unlock()
		}
	}
}

// handleResizeInner handles the actual resize logic (caller must hold lock).
func (l *Line) handleResizeInner(oldHeight, newWidth, newHeight int) {
	out := l.term.Writer()

	// 1. Fully reset scroll region first
	fmt.Fprint(out, terminal.ResetScrollRegion)

	// 2. Clear the OLD status area locations (both lines)
	fmt.Fprintf(out, terminal.MoveToFormat, oldHeight-1, 1)
	fmt.Fprint(out, terminal.ClearLine)
	fmt.Fprintf(out, terminal.MoveToFormat, oldHeight, 1)
	fmt.Fprint(out, terminal.ClearLine)

	// 3. Clear the NEW status area locations too
	fmt.Fprintf(out, terminal.MoveToFormat, newHeight-1, 1)
	fmt.Fprint(out, terminal.ClearLine)
	fmt.Fprintf(out, terminal.MoveToFormat, newHeight, 1)
	fmt.Fprint(out, terminal.ClearLine)

	// 4. APT trick: newline + save cursor
	fmt.Fprint(out, "\n")
	fmt.Fprint(out, terminal.SaveCursor)

	// 5. Set new scroll region (reserve 2 rows for status area)
	if newHeight >= 3 {
		fmt.Fprintf(out, terminal.SetScrollRegionFormat, 1, newHeight-2)
	}

	// 6. Restore cursor and move up (stay in scroll region)
	fmt.Fprint(out, terminal.RestoreCursor)
	fmt.Fprint(out, terminal.CursorUp)

	// 7. Render status area at new position
	l.render()
}

// SetStepInfo sets the step information for the progress display.
func (l *Line) SetStepInfo(step, total int, stepName, storyKey, model string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	stepChanged := l.state.StepName != stepName
	l.state.Step = step
	l.state.Total = total
	l.state.StepName = stepName
	l.state.StoryKey = storyKey
	l.state.Model = model

	now := time.Now()
	if l.state.StartTime.IsZero() {
		l.state.StartTime = now
	}
	if stepChanged || l.state.StepStartTime.IsZero() {
		l.state.StepStartTime = now
		l.state.ActivityStart = now // Start activity timer immediately
		// Reset thinking time tracking for new step
		l.state.ThinkingStart = now
		l.state.ThinkingDuration = 0
		l.state.HadFirstResponse = false
	}
	l.render()
}

// SetOperation sets the operation context (e.g., "Epic 6", "Story 2/3").
func (l *Line) SetOperation(operation string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.state.Operation = operation
	l.render()
}

// RecordFirstResponse records when the first response arrives (for thinking time).
// Call this when text or tool use is first received. Safe to call multiple times.
func (l *Line) RecordFirstResponse() {
	l.mu.Lock()
	defer l.mu.Unlock()

	if !l.state.HadFirstResponse && !l.state.ThinkingStart.IsZero() {
		l.state.ThinkingDuration = time.Since(l.state.ThinkingStart)
		l.state.HadFirstResponse = true
	}
}

// SetCurrentTool sets the current tool/activity being executed.
func (l *Line) SetCurrentTool(tool string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	activityChanged := l.state.CurrentTool != tool
	l.state.CurrentTool = tool
	if activityChanged || l.state.ActivityStart.IsZero() {
		l.state.ActivityStart = time.Now()
	}
	l.render()
}

// AddTokens adds token counts to the running total and refreshes the display.
func (l *Line) AddTokens(inputTokens, outputTokens int) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.state.InputTokens += inputTokens
	l.state.OutputTokens += outputTokens
	l.render()
}

// IncrementToolCount increments the tool usage counter.
func (l *Line) IncrementToolCount() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.state.ToolCount++
	l.render()
}

// Refresh updates the status bar display.
func (l *Line) Refresh() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.render()
}

// Update is a convenience method that updates all state and refreshes.
func (l *Line) Update(
	step, total int,
	stepName, storyKey, model string,
	elapsed time.Duration,
	currentTool string,
) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.state.Step = step
	l.state.Total = total
	l.state.StepName = stepName
	l.state.StoryKey = storyKey
	l.state.Model = model
	l.state.CurrentTool = currentTool
	l.render()
}

// UpdateWithRateLimit shows rate limit status.
func (l *Line) UpdateWithRateLimit(
	step, total int,
	stepName, storyKey string,
	elapsed time.Duration,
	resetTime time.Time,
) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.state.Step = step
	l.state.Total = total
	l.state.StepName = stepName
	l.state.StoryKey = storyKey
	l.state.CurrentTool = iconPaused + " Reset in " + formatDuration(time.Until(resetTime))
	l.render()
}

// Clear removes the status area and resets scrolling to normal.
func (l *Line) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()

	out := l.term.Writer()
	_, height := l.term.Size()

	// Stop goroutines
	if l.cancel != nil {
		l.cancel()
	}

	// Reset scrolling region
	fmt.Fprint(out, terminal.ResetScrollRegion)

	// Clear both status area lines
	fmt.Fprintf(out, terminal.MoveToFormat, height-1, 1)
	fmt.Fprint(out, terminal.ClearLine)
	fmt.Fprintf(out, terminal.MoveToFormat, height, 1)
	fmt.Fprint(out, terminal.ClearLine)

	// Move cursor to where output should continue
	fmt.Fprintf(out, terminal.MoveToFormat, height, 1)
}

// Done shows completion message with token info then clears.
func (l *Line) Done(success bool, duration time.Duration) {
	if !l.enabled {
		return
	}

	l.mu.Lock()
	_, height := l.term.Size()
	totalTokens := l.state.InputTokens + l.state.OutputTokens
	toolCount := l.state.ToolCount

	if err := l.activity.RenderDone(height, success, duration, totalTokens, toolCount, l.state.VerbIdx); err != nil {
		l.mu.Unlock()
		return
	}
	l.mu.Unlock()

	time.Sleep(500 * time.Millisecond)
	l.Clear()
}

// StartTime returns the start time.
func (l *Line) StartTime() time.Time {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.state.StartTime
}

// render writes both the activity line and status bar (caller must hold lock).
func (l *Line) render() {
	if !l.enabled {
		return
	}

	width, height := l.term.Size()
	if height < 3 {
		return
	}

	out := l.term.Writer()

	// Save cursor, hide it during update
	fmt.Fprint(out, terminal.SaveCursor)
	fmt.Fprint(out, terminal.HideCursor)

	// Render activity line (second-to-last row)
	activityState := ActivityState{
		VerbIdx:          l.state.VerbIdx,
		ActivityStart:    l.state.ActivityStart,
		CurrentTool:      l.state.CurrentTool,
		InputTokens:      l.state.InputTokens,
		OutputTokens:     l.state.OutputTokens,
		HadFirstResponse: l.state.HadFirstResponse,
		ThinkingDuration: l.state.ThinkingDuration,
		ToolCount:        l.state.ToolCount,
	}
	if err := l.activity.Render(height, activityState); err != nil {
		// Continue anyway
	}

	// Render status bar (last row)
	statusState := StatusState{
		Step:          l.state.Step,
		Total:         l.state.Total,
		StepName:      l.state.StepName,
		StoryKey:      l.state.StoryKey,
		Model:         l.state.Model,
		Operation:     l.state.Operation,
		StepStartTime: l.state.StepStartTime,
		StartTime:     l.state.StartTime,
		Width:         width,
	}
	if err := l.statusBar.Render(height, statusState); err != nil {
		// Continue anyway
	}

	// Restore cursor
	fmt.Fprint(out, terminal.ShowCursor)
	fmt.Fprint(out, terminal.RestoreCursor)
}
