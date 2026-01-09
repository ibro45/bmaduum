// Package state provides lifecycle execution state persistence for resume functionality.
package state

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// StateFileName is the name of the state file in the working directory.
const StateFileName = ".bmad-state.json"

// ErrNoState is returned by Load when no state file exists.
var ErrNoState = errors.New("no state file exists")

// State represents the current lifecycle execution state.
type State struct {
	StoryKey    string `json:"story_key"`
	StepIndex   int    `json:"step_index"`
	TotalSteps  int    `json:"total_steps"`
	StartStatus string `json:"start_status"`
}

// Manager handles state persistence operations.
type Manager struct {
	dir string
}

// NewManager creates a new state manager for the given directory.
func NewManager(dir string) *Manager {
	return &Manager{dir: dir}
}

// statePath returns the full path to the state file.
func (m *Manager) statePath() string {
	return filepath.Join(m.dir, StateFileName)
}

// Save persists the state to disk atomically.
func (m *Manager) Save(state State) error {
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}

	// Write to temp file first for atomic operation
	tmpPath := m.statePath() + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return err
	}

	// Rename temp file to final location (atomic on POSIX)
	return os.Rename(tmpPath, m.statePath())
}

// Load reads the state from disk.
func (m *Manager) Load() (State, error) {
	data, err := os.ReadFile(m.statePath())
	if err != nil {
		if os.IsNotExist(err) {
			return State{}, ErrNoState
		}
		return State{}, err
	}

	var state State
	if err := json.Unmarshal(data, &state); err != nil {
		return State{}, err
	}

	return state, nil
}

// Clear removes the state file if it exists.
func (m *Manager) Clear() error {
	err := os.Remove(m.statePath())
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// Exists returns true if a state file exists.
func (m *Manager) Exists() bool {
	_, err := os.Stat(m.statePath())
	return err == nil
}
