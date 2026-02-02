// Package sprint provides story file-based status management.
//
// The sprint package implements StoryStatusManager which treats story files as
// the source of truth for development status, with sprint-status.yaml acting
// as a derived cache. This eliminates sync drift between the two sources.
//
// Key types:
//   - [StoryStatusManager] - Reads/writes story status from/to story files
//   - [Rebuilder] - Rebuilds sprint-status.yaml from all story files
//
// Story file format:
//   # Story {epic}.{number}: {title}
//
//   Status: {backlog|ready-for-dev|in-progress|review|done}
//
// The status is parsed from the first line matching "Status: <value>".
package sprint

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"bmad-automate/internal/status"
)

// DefaultStoryDir is the canonical location of story files relative to project root.
const DefaultStoryDir = "_bmad-output/implementation-artifacts/stories"

// StoryStatusManager implements lifecycle.StatusReader and lifecycle.StatusWriter
// using story files as the source of truth.
//
// Story files are the authoritative source of status. The sprint-status.yaml
// file is treated as a derived cache that is updated as a side effect of
// story file updates.
type StoryStatusManager struct {
	storyDir   string
	sprintPath string
	// cacheWriter is used to sync updates to sprint-status.yaml
	cacheWriter *status.Writer
}

// NewStoryStatusManager creates a StoryStatusManager with the specified paths.
//
// The storyDir is the directory containing story .md files.
// The sprintPath is the path to sprint-status.yaml (used as cache).
// Pass empty strings to use default paths relative to current directory.
func NewStoryStatusManager(storyDir, sprintPath string) *StoryStatusManager {
	if storyDir == "" {
		storyDir = DefaultStoryDir
	}
	if sprintPath == "" {
		sprintPath = status.DefaultStatusPath
	}

	// Extract base path from sprintPath to create cache writer
	basePath := ""
	if sprintPath != status.DefaultStatusPath {
		// If custom sprint path, extract base path
		dir := filepath.Dir(sprintPath)
		// Walk up to find project root (assumed to be where _bmad-output lives)
		for dir != "." && dir != "/" {
			if strings.HasSuffix(dir, "_bmad-output") {
				basePath = filepath.Dir(dir)
				break
			}
			dir = filepath.Dir(dir)
		}
	}

	return &StoryStatusManager{
		storyDir:    storyDir,
		sprintPath:  sprintPath,
		cacheWriter: status.NewWriter(basePath),
	}
}

// GetStoryStatus reads the status from a story file.
//
// It looks for the first line starting with "Status:" and returns the value.
// The storyKey is mapped to a filename by adding .md extension.
func (m *StoryStatusManager) GetStoryStatus(storyKey string) (status.Status, error) {
	storyPath := m.findStoryFile(storyKey)

	st, err := m.readStoryStatus(storyPath)
	if err != nil {
		return "", fmt.Errorf("failed to get status for %s: %w", storyKey, err)
	}

	return st, nil
}

// UpdateStatus updates the story file and syncs to sprint-status.yaml.
//
// The update process:
//  1. Updates the story file's Status line atomically
//  2. Syncs the change to sprint-status.yaml via status.Writer
func (m *StoryStatusManager) UpdateStatus(storyKey string, newStatus status.Status) error {
	// Validate the new status
	if !newStatus.IsValid() {
		return fmt.Errorf("invalid status: %s", newStatus)
	}

	storyPath := m.findStoryFile(storyKey)

	// Update story file (source of truth)
	if err := m.writeStoryStatus(storyPath, newStatus); err != nil {
		return fmt.Errorf("failed to update story file for %s: %w", storyKey, err)
	}

	// Sync to sprint-status.yaml (cache)
	if err := m.cacheWriter.UpdateStatus(storyKey, newStatus); err != nil {
		return fmt.Errorf("failed to sync status to sprint-status.yaml for %s: %w", storyKey, err)
	}

	return nil
}

// GetEpicStories returns all story keys belonging to an epic, sorted by story number.
//
// Story keys are discovered by reading the story directory and matching files
// with the pattern {epicID}-{number}-*.md. Results are sorted numerically.
func (m *StoryStatusManager) GetEpicStories(epicID string) ([]string, error) {
	entries, err := os.ReadDir(m.storyDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read story directory: %w", err)
	}

	var stories []struct {
		key string
		num int
	}

	prefix := epicID + "-"
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		name := strings.TrimSuffix(entry.Name(), ".md")
		if !strings.HasPrefix(name, prefix) {
			continue
		}

		// Extract story number (second segment)
		remainder := strings.TrimPrefix(name, prefix)
		parts := strings.SplitN(remainder, "-", 2)
		if len(parts) < 1 {
			continue
		}

		num, err := parseStoryNumber(parts[0])
		if err != nil {
			// Not a numeric story number, skip
			continue
		}

		stories = append(stories, struct {
			key string
			num int
		}{key: name, num: num})
	}

	if len(stories) == 0 {
		return nil, fmt.Errorf("no stories found for epic: %s", epicID)
	}

	// Sort by story number
	for i := 0; i < len(stories)-1; i++ {
		for j := i + 1; j < len(stories); j++ {
			if stories[i].num > stories[j].num {
				stories[i], stories[j] = stories[j], stories[i]
			}
		}
	}

	result := make([]string, len(stories))
	for i, s := range stories {
		result[i] = s.key
	}

	return result, nil
}

// GetAllEpics returns all epic IDs with stories, sorted numerically.
//
// Epic IDs are extracted from story filenames (format: {epicID}-{storyNum}-{description}.md).
func (m *StoryStatusManager) GetAllEpics() ([]string, error) {
	entries, err := os.ReadDir(m.storyDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read story directory: %w", err)
	}

	epicMap := make(map[string]int)

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		name := strings.TrimSuffix(entry.Name(), ".md")
		parts := strings.SplitN(name, "-", 2)
		if len(parts) < 1 {
			continue
		}

		epicID := parts[0]
		num, _ := parseStoryNumber(epicID)
		epicMap[epicID] = num
	}

	if len(epicMap) == 0 {
		return []string{}, nil
	}

	// Convert to slice and sort
	type epicInfo struct {
		id  string
		num int
	}

	epics := make([]epicInfo, 0, len(epicMap))
	for id, num := range epicMap {
		epics = append(epics, epicInfo{id: id, num: num})
	}

	// Sort by numeric value
	for i := 0; i < len(epics)-1; i++ {
		for j := i + 1; j < len(epics); j++ {
			if epics[i].num > epics[j].num {
				epics[i], epics[j] = epics[j], epics[i]
			}
		}
	}

	result := make([]string, len(epics))
	for i, e := range epics {
		result[i] = e.id
	}

	return result, nil
}

// findStoryFile maps a story key to its file path.
func (m *StoryStatusManager) findStoryFile(storyKey string) string {
	return filepath.Join(m.storyDir, storyKey+".md")
}

// readStoryStatus reads the status from a story file.
//
// It scans for the first line starting with "Status:" (case-insensitive)
// and returns the value after the colon.
func (m *StoryStatusManager) readStoryStatus(path string) (status.Status, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("failed to open story file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		// Look for Status: line (case-insensitive)
		if strings.HasPrefix(strings.ToLower(trimmed), "status:") {
			parts := strings.SplitN(trimmed, ":", 2)
			if len(parts) == 2 {
				statusValue := strings.TrimSpace(parts[1])
				return status.Status(statusValue), nil
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("failed to read story file: %w", err)
	}

	return "", fmt.Errorf("status line not found in story file")
}

// writeStoryStatus updates the status line in a story file atomically.
//
// The update is performed by:
//  1. Reading the entire file
//  2. Finding and replacing the Status line
//  3. Writing to a temp file
//  4. Renaming temp file to original (atomic)
func (m *StoryStatusManager) writeStoryStatus(path string, newStatus status.Status) error {
	// Read existing file
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read story file: %w", err)
	}

	lines := strings.Split(string(data), "\n")

	// Find and update the Status line
	found := false
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(strings.ToLower(trimmed), "status:") {
			// Preserve indentation from original line
			colonIdx := strings.Index(line, ":")
			if colonIdx >= 0 {
				prefix := line[:colonIdx]
				lines[i] = prefix + ": " + string(newStatus)
				found = true
				break
			}
		}
	}

	if !found {
		return fmt.Errorf("status line not found in story file")
	}

	// Join lines back together
	updatedData := strings.Join(lines, "\n")

	// Write atomically (temp file + rename)
	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, []byte(updatedData), 0644); err != nil {
		return fmt.Errorf("failed to write temp file: %w", err)
	}

	if err := os.Rename(tmpPath, path); err != nil {
		// Clean up temp file on rename failure
		os.Remove(tmpPath)
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	return nil
}

// parseStoryNumber attempts to parse a string as a story/epic number.
// Returns the number or an error if parsing fails.
func parseStoryNumber(s string) (int, error) {
	var num int
	_, err := fmt.Sscanf(s, "%d", &num)
	return num, err
}
