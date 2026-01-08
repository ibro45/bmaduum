// Package status provides functionality for reading sprint status from YAML files.
package status

// Status represents a story's development status.
type Status string

// Status constants for known development statuses.
const (
	StatusBacklog     Status = "backlog"
	StatusReadyForDev Status = "ready-for-dev"
	StatusInProgress  Status = "in-progress"
	StatusReview      Status = "review"
	StatusDone        Status = "done"
)

// IsValid returns true if the status is a known valid status value.
func (s Status) IsValid() bool {
	switch s {
	case StatusBacklog, StatusReadyForDev, StatusInProgress, StatusReview, StatusDone:
		return true
	default:
		return false
	}
}

// SprintStatus represents the structure of a sprint-status.yaml file.
type SprintStatus struct {
	DevelopmentStatus map[string]Status `yaml:"development_status"`
}
