package cli

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"bmaduum/internal/status"
)

// StatusUpdate represents a status update for testing.
type StatusUpdate struct {
	StoryKey  string
	NewStatus status.Status
}

// MockWorkflowRunner is a mock for testing.
type MockWorkflowRunner struct {
	// ExecutedWorkflows records all workflow executions in order.
	ExecutedWorkflows []string
	// FailOnWorkflow specifies which workflow should fail (returns exit code 1).
	FailOnWorkflow string
}

func (m *MockWorkflowRunner) RunSingle(ctx context.Context, workflowName, storyKey string) int {
	m.ExecutedWorkflows = append(m.ExecutedWorkflows, workflowName)
	if m.FailOnWorkflow == workflowName {
		return 1
	}
	return 0
}

func (m *MockWorkflowRunner) RunRaw(ctx context.Context, prompt string) int {
	return 0
}

// MockStatusWriter is a mock for testing.
type MockStatusWriter struct {
	// Updates records all status updates.
	Updates []StatusUpdate
}

func (m *MockStatusWriter) UpdateStatus(storyKey string, newStatus status.Status) error {
	m.Updates = append(m.Updates, StatusUpdate{StoryKey: storyKey, NewStatus: newStatus})
	return nil
}

// createSprintStatusFile creates a sprint-status.yaml file in a temporary directory for testing.
func createSprintStatusFile(t *testing.T, tmpDir string, content string) {
	t.Helper()

	statusDir := filepath.Join(tmpDir, "_bmad-output", "implementation-artifacts")
	err := os.MkdirAll(statusDir, 0755)
	if err != nil {
		t.Fatalf("failed to create status directory: %v", err)
	}

	statusPath := filepath.Join(statusDir, "sprint-status.yaml")
	err = os.WriteFile(statusPath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to write status file: %v", err)
	}
}
