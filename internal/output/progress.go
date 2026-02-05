// Package output provides terminal output formatting using lipgloss styles.
package output

import (
	"io"

	"bmaduum/internal/output/progress"
)

// ProgressLine provides a fixed two-line status area with scrolling output above.
//
// The status area consists of:
//   - Activity line (second-to-last): Spinner + activity verb + timer + token count
//   - Status bar (last row): Operation + step info + story key + model + total timer
//
// This type is re-exported from the progress sub-package for backward compatibility.
type ProgressLine = progress.Line

// NewProgressLine creates a new progress line writer.
//
// This function is re-exported from the progress sub-package for backward compatibility.
func NewProgressLine(out io.Writer) *ProgressLine {
	return progress.NewLine(out)
}
