package diff

import (
	"regexp"
	"strconv"
	"strings"
)

// hunkHeaderRegex matches unified diff hunk headers like "@@ -1,3 +1,4 @@" or "@@ -1 +1,2 @".
var hunkHeaderRegex = regexp.MustCompile(`^@@ -(\d+)(?:,(\d+))? \+(\d+)(?:,(\d+))? @@`)

// Parser parses unified diff format into structured Diff objects.
type Parser struct{}

// NewParser creates a new diff parser.
func NewParser() *Parser {
	return &Parser{}
}

// Parse parses a unified diff string into a structured Diff.
//
// Handles standard unified diff format as produced by git diff:
//
//	--- a/file.txt
//	+++ b/file.txt
//	@@ -1,3 +1,4 @@
//	 context line
//	-deleted line
//	+added line
//	 context line
func (p *Parser) Parse(diffText string) (*Diff, error) {
	diff := &Diff{}
	lines := strings.Split(diffText, "\n")

	var currentHunk *Hunk
	oldLine, newLine := 0, 0

	for _, line := range lines {
		switch {
		case strings.HasPrefix(line, "--- "):
			// Old file header
			diff.OldFile = strings.TrimPrefix(line, "--- ")
			// Strip a/ prefix if present
			diff.OldFile = strings.TrimPrefix(diff.OldFile, "a/")

		case strings.HasPrefix(line, "+++ "):
			// New file header
			diff.NewFile = strings.TrimPrefix(line, "+++ ")
			// Strip b/ prefix if present
			diff.NewFile = strings.TrimPrefix(diff.NewFile, "b/")

		case strings.HasPrefix(line, "@@ "):
			// Hunk header
			hunk := p.parseHunkHeader(line)
			if hunk != nil {
				currentHunk = hunk
				oldLine = hunk.OldStart
				newLine = hunk.NewStart
				diff.Hunks = append(diff.Hunks, *currentHunk)
				// Update pointer to the actual hunk in the slice
				currentHunk = &diff.Hunks[len(diff.Hunks)-1]
			}

		case len(line) > 0 && line[0] == '+' && !strings.HasPrefix(line, "+++"):
			// Added line
			if currentHunk != nil {
				currentHunk.Lines = append(currentHunk.Lines, Line{
					Type:       LineTypeAdded,
					Content:    line[1:], // Remove the + prefix
					NewLineNum: newLine,
				})
				newLine++
				diff.Added++
			}

		case len(line) > 0 && line[0] == '-' && !strings.HasPrefix(line, "---"):
			// Deleted line
			if currentHunk != nil {
				currentHunk.Lines = append(currentHunk.Lines, Line{
					Type:       LineTypeDeleted,
					Content:    line[1:], // Remove the - prefix
					OldLineNum: oldLine,
				})
				oldLine++
				diff.Deleted++
			}

		case len(line) > 0 && line[0] == ' ':
			// Context line
			if currentHunk != nil {
				currentHunk.Lines = append(currentHunk.Lines, Line{
					Type:       LineTypeContext,
					Content:    line[1:], // Remove the space prefix
					OldLineNum: oldLine,
					NewLineNum: newLine,
				})
				oldLine++
				newLine++
			}

		case line == "":
			// Empty context line (could be a blank line in the diff)
			if currentHunk != nil {
				currentHunk.Lines = append(currentHunk.Lines, Line{
					Type:       LineTypeContext,
					Content:    "",
					OldLineNum: oldLine,
					NewLineNum: newLine,
				})
				oldLine++
				newLine++
			}
		}
	}

	return diff, nil
}

// parseHunkHeader parses a hunk header line like "@@ -1,3 +1,4 @@".
func (p *Parser) parseHunkHeader(line string) *Hunk {
	matches := hunkHeaderRegex.FindStringSubmatch(line)
	if matches == nil {
		return nil
	}

	hunk := &Hunk{}

	// Parse old start and count
	hunk.OldStart, _ = strconv.Atoi(matches[1])
	if matches[2] != "" {
		hunk.OldCount, _ = strconv.Atoi(matches[2])
	} else {
		hunk.OldCount = 1
	}

	// Parse new start and count
	hunk.NewStart, _ = strconv.Atoi(matches[3])
	if matches[4] != "" {
		hunk.NewCount, _ = strconv.Atoi(matches[4])
	} else {
		hunk.NewCount = 1
	}

	return hunk
}

// IsUnifiedDiff checks if text looks like a unified diff.
//
// Returns true if the text contains diff markers like "--- a/",
// "+++ b/", or hunk headers "@@".
func IsUnifiedDiff(text string) bool {
	return (strings.Contains(text, "--- a/") || strings.Contains(text, "--- ")) &&
		(strings.Contains(text, "+++ b/") || strings.Contains(text, "+++ ")) &&
		strings.Contains(text, "@@ ")
}

// ParseUnifiedDiff is a convenience function that creates a parser and parses the input.
func ParseUnifiedDiff(diffText string) (*Diff, error) {
	return NewParser().Parse(diffText)
}
