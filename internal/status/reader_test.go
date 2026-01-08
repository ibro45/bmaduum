package status

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewReader(t *testing.T) {
	reader := NewReader("/some/path")

	assert.NotNil(t, reader)
	assert.Equal(t, "/some/path", reader.basePath)
}

func TestReader_Read_Success(t *testing.T) {
	tmpDir := t.TempDir()

	// Create the nested directory structure
	statusDir := filepath.Join(tmpDir, "_bmad-output", "implementation-artifacts")
	err := os.MkdirAll(statusDir, 0755)
	require.NoError(t, err)

	// Create a valid sprint-status.yaml
	statusContent := `development_status:
  7-1-define-schema: ready-for-dev
  7-2-create-api: in-progress
  7-3-build-ui: backlog
`
	statusPath := filepath.Join(statusDir, "sprint-status.yaml")
	err = os.WriteFile(statusPath, []byte(statusContent), 0644)
	require.NoError(t, err)

	reader := NewReader(tmpDir)
	status, err := reader.Read()

	require.NoError(t, err)
	assert.NotNil(t, status)
	assert.Len(t, status.DevelopmentStatus, 3)
	assert.Equal(t, StatusReadyForDev, status.DevelopmentStatus["7-1-define-schema"])
	assert.Equal(t, StatusInProgress, status.DevelopmentStatus["7-2-create-api"])
	assert.Equal(t, StatusBacklog, status.DevelopmentStatus["7-3-build-ui"])
}

func TestReader_Read_FileNotFound(t *testing.T) {
	tmpDir := t.TempDir()

	reader := NewReader(tmpDir)
	status, err := reader.Read()

	assert.Error(t, err)
	assert.Nil(t, status)
	assert.Contains(t, err.Error(), "failed to read sprint status")
}

func TestReader_Read_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()

	// Create the nested directory structure
	statusDir := filepath.Join(tmpDir, "_bmad-output", "implementation-artifacts")
	err := os.MkdirAll(statusDir, 0755)
	require.NoError(t, err)

	// Create an invalid YAML file
	invalidContent := `development_status:
  - this is not a map
    missing: colon
`
	statusPath := filepath.Join(statusDir, "sprint-status.yaml")
	err = os.WriteFile(statusPath, []byte(invalidContent), 0644)
	require.NoError(t, err)

	reader := NewReader(tmpDir)
	status, err := reader.Read()

	assert.Error(t, err)
	assert.Nil(t, status)
	assert.Contains(t, err.Error(), "failed to read sprint status")
}

func TestReader_GetStoryStatus_Found(t *testing.T) {
	tmpDir := t.TempDir()

	// Create the nested directory structure
	statusDir := filepath.Join(tmpDir, "_bmad-output", "implementation-artifacts")
	err := os.MkdirAll(statusDir, 0755)
	require.NoError(t, err)

	statusContent := `development_status:
  7-1-define-schema: ready-for-dev
  7-2-create-api: in-progress
`
	statusPath := filepath.Join(statusDir, "sprint-status.yaml")
	err = os.WriteFile(statusPath, []byte(statusContent), 0644)
	require.NoError(t, err)

	reader := NewReader(tmpDir)
	status, err := reader.GetStoryStatus("7-1-define-schema")

	require.NoError(t, err)
	assert.Equal(t, StatusReadyForDev, status)
}

func TestReader_GetStoryStatus_NotFound(t *testing.T) {
	tmpDir := t.TempDir()

	// Create the nested directory structure
	statusDir := filepath.Join(tmpDir, "_bmad-output", "implementation-artifacts")
	err := os.MkdirAll(statusDir, 0755)
	require.NoError(t, err)

	statusContent := `development_status:
  7-1-define-schema: ready-for-dev
`
	statusPath := filepath.Join(statusDir, "sprint-status.yaml")
	err = os.WriteFile(statusPath, []byte(statusContent), 0644)
	require.NoError(t, err)

	reader := NewReader(tmpDir)
	status, err := reader.GetStoryStatus("nonexistent-story")

	assert.Error(t, err)
	assert.Equal(t, Status(""), status)
	assert.Contains(t, err.Error(), "story not found: nonexistent-story")
}

func TestReader_GetStoryStatus_MultipleStories(t *testing.T) {
	tmpDir := t.TempDir()

	// Create the nested directory structure
	statusDir := filepath.Join(tmpDir, "_bmad-output", "implementation-artifacts")
	err := os.MkdirAll(statusDir, 0755)
	require.NoError(t, err)

	statusContent := `development_status:
  7-1-define-schema: ready-for-dev
  7-2-create-api: in-progress
  7-3-build-ui: backlog
  7-4-add-tests: review
  7-5-deploy: done
`
	statusPath := filepath.Join(statusDir, "sprint-status.yaml")
	err = os.WriteFile(statusPath, []byte(statusContent), 0644)
	require.NoError(t, err)

	reader := NewReader(tmpDir)

	tests := []struct {
		storyKey string
		want     Status
	}{
		{"7-1-define-schema", StatusReadyForDev},
		{"7-2-create-api", StatusInProgress},
		{"7-3-build-ui", StatusBacklog},
		{"7-4-add-tests", StatusReview},
		{"7-5-deploy", StatusDone},
	}

	for _, tt := range tests {
		t.Run(tt.storyKey, func(t *testing.T) {
			status, err := reader.GetStoryStatus(tt.storyKey)
			require.NoError(t, err)
			assert.Equal(t, tt.want, status)
		})
	}
}

func TestReader_GetStoryStatus_FileNotFound(t *testing.T) {
	tmpDir := t.TempDir()

	reader := NewReader(tmpDir)
	status, err := reader.GetStoryStatus("any-story")

	assert.Error(t, err)
	assert.Equal(t, Status(""), status)
	assert.Contains(t, err.Error(), "failed to read sprint status")
}
