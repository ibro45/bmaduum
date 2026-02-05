// Package diff provides unified diff parsing and rendering.
//
// The package handles standard unified diff format as produced by git diff,
// parsing it into structured data and providing rich terminal rendering with
// colors, line numbers, and optional syntax highlighting.
//
// Key types:
//   - [Diff] - Complete file diff with metadata and hunks
//   - [DiffHunk] - Contiguous section of changes
//   - [DiffLine] - Single line with type and content
//   - [Parser] - Unified diff parser
//   - [Renderer] - Terminal diff renderer
package diff

// LineType represents the type of a diff line.
type LineType int

const (
	// LineTypeContext represents an unchanged line (starts with space).
	LineTypeContext LineType = iota
	// LineTypeAdded represents an added line (starts with +).
	LineTypeAdded
	// LineTypeDeleted represents a deleted line (starts with -).
	LineTypeDeleted
)

// Line represents a single line in a diff.
type Line struct {
	Type       LineType // Context, Added, or Deleted
	Content    string   // Line content without the prefix (+, -, space)
	OldLineNum int      // Line number in old file (0 if not applicable)
	NewLineNum int      // Line number in new file (0 if not applicable)
}

// Hunk represents a contiguous section of changes in a diff.
type Hunk struct {
	OldStart int    // Starting line number in old file
	OldCount int    // Number of lines in old file
	NewStart int    // Starting line number in new file
	NewCount int    // Number of lines in new file
	Lines    []Line // Lines in this hunk
}

// Diff represents a complete file diff with metadata and hunks.
type Diff struct {
	OldFile string // Path to old file (from --- line)
	NewFile string // Path to new file (from +++ line)
	Hunks   []Hunk // All hunks in this diff
	Added   int    // Total lines added
	Deleted int    // Total lines deleted
}

// Summary returns a human-readable summary of the diff.
func (d *Diff) Summary() string {
	if d.Added == 0 && d.Deleted == 0 {
		return "No changes"
	}
	if d.Added > 0 && d.Deleted > 0 {
		return pluralize(d.Added, "line") + " added, " + pluralize(d.Deleted, "line") + " removed"
	}
	if d.Added > 0 {
		return pluralize(d.Added, "line") + " added"
	}
	return pluralize(d.Deleted, "line") + " removed"
}

// pluralize returns "n item" or "n items" based on count.
func pluralize(n int, singular string) string {
	if n == 1 {
		return "1 " + singular
	}
	return itoa(n) + " " + singular + "s"
}

// itoa converts an int to string without importing strconv.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	if n < 0 {
		return "-" + itoa(-n)
	}
	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	return string(digits)
}
