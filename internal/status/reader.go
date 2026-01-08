package status

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// DefaultStatusPath is the default location of the sprint status file.
const DefaultStatusPath = "_bmad-output/implementation-artifacts/sprint-status.yaml"

// Reader reads sprint status from YAML files.
type Reader struct {
	basePath string
}

// NewReader creates a new Reader with the specified base path.
func NewReader(basePath string) *Reader {
	return &Reader{
		basePath: basePath,
	}
}

// Read reads and parses the sprint status file.
func (r *Reader) Read() (*SprintStatus, error) {
	fullPath := filepath.Join(r.basePath, DefaultStatusPath)

	data, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read sprint status: %w", err)
	}

	var status SprintStatus
	if err := yaml.Unmarshal(data, &status); err != nil {
		return nil, fmt.Errorf("failed to read sprint status: %w", err)
	}

	return &status, nil
}

// GetStoryStatus returns the status for a specific story key.
func (r *Reader) GetStoryStatus(storyKey string) (Status, error) {
	sprintStatus, err := r.Read()
	if err != nil {
		return "", err
	}

	status, ok := sprintStatus.DevelopmentStatus[storyKey]
	if !ok {
		return "", fmt.Errorf("story not found: %s", storyKey)
	}

	return status, nil
}

// GetEpicStories returns all story keys belonging to an epic, sorted by story number.
// Story keys are expected to follow the pattern {epicID}-{storyNum}-{description}.
func (r *Reader) GetEpicStories(epicID string) ([]string, error) {
	sprintStatus, err := r.Read()
	if err != nil {
		return nil, err
	}

	// Collect all keys matching the epic ID pattern
	type storyWithNum struct {
		key string
		num int
	}
	var stories []storyWithNum

	prefix := epicID + "-"
	for key := range sprintStatus.DevelopmentStatus {
		if !strings.HasPrefix(key, prefix) {
			continue
		}

		// Extract the story number (second segment)
		// Format: {epicID}-{storyNum}-{rest}
		remainder := strings.TrimPrefix(key, prefix)
		parts := strings.SplitN(remainder, "-", 2)
		if len(parts) < 1 {
			continue
		}

		num, err := strconv.Atoi(parts[0])
		if err != nil {
			// Not a numeric story number, skip
			continue
		}

		stories = append(stories, storyWithNum{key: key, num: num})
	}

	if len(stories) == 0 {
		return nil, fmt.Errorf("no stories found for epic: %s", epicID)
	}

	// Sort by story number
	sort.Slice(stories, func(i, j int) bool {
		return stories[i].num < stories[j].num
	})

	// Extract just the keys
	result := make([]string, len(stories))
	for i, s := range stories {
		result[i] = s.key
	}

	return result, nil
}
