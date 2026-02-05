// Package progress provides terminal progress display functionality.
package progress

import (
	"testing"
	"time"
)

func TestState_DefaultValues(t *testing.T) {
	s := State{}
	if s.Step != 0 {
		t.Errorf("Step should default to 0, got %d", s.Step)
	}
	if s.Total != 0 {
		t.Errorf("Total should default to 0, got %d", s.Total)
	}
	if s.InputTokens != 0 {
		t.Errorf("InputTokens should default to 0, got %d", s.InputTokens)
	}
	if s.OutputTokens != 0 {
		t.Errorf("OutputTokens should default to 0, got %d", s.OutputTokens)
	}
	if s.SpinnerIdx != 0 {
		t.Errorf("SpinnerIdx should default to 0, got %d", s.SpinnerIdx)
	}
	if s.VerbIdx != 0 {
		t.Errorf("VerbIdx should default to 0, got %d", s.VerbIdx)
	}
}

func TestState_TimeFieldsAreZero(t *testing.T) {
	s := State{}
	if !s.StartTime.IsZero() {
		t.Error("StartTime should be zero")
	}
	if !s.StepStartTime.IsZero() {
		t.Error("StepStartTime should be zero")
	}
	if !s.ActivityStart.IsZero() {
		t.Error("ActivityStart should be zero")
	}
	if !s.ThinkingStart.IsZero() {
		t.Error("ThinkingStart should be zero")
	}
	if s.ThinkingDuration != 0 {
		t.Errorf("ThinkingDuration should be 0, got %v", s.ThinkingDuration)
	}
}

func TestActivityState(t *testing.T) {
	now := time.Now()
	as := ActivityState{
		VerbIdx:          5,
		ActivityStart:    now,
		CurrentTool:      "Edit",
		InputTokens:      1000,
		OutputTokens:     2000,
		HadFirstResponse: true,
		ThinkingDuration: 2 * time.Second,
	}

	if as.VerbIdx != 5 {
		t.Errorf("VerbIdx = %d, want 5", as.VerbIdx)
	}
	if as.CurrentTool != "Edit" {
		t.Errorf("CurrentTool = %s, want Edit", as.CurrentTool)
	}
	if as.InputTokens != 1000 {
		t.Errorf("InputTokens = %d, want 1000", as.InputTokens)
	}
	if as.OutputTokens != 2000 {
		t.Errorf("OutputTokens = %d, want 2000", as.OutputTokens)
	}
	if !as.HadFirstResponse {
		t.Error("HadFirstResponse should be true")
	}
	if as.ThinkingDuration != 2*time.Second {
		t.Errorf("ThinkingDuration = %v, want 2s", as.ThinkingDuration)
	}
}

func TestStatusState(t *testing.T) {
	now := time.Now()
	ss := StatusState{
		Step:          2,
		Total:         5,
		StepName:      "create-story",
		StoryKey:      "7-1-test",
		Model:         "claude-opus-4-5",
		Operation:     "Epic 3",
		StepStartTime: now,
		StartTime:     now.Add(-1 * time.Minute),
		Width:         80,
	}

	if ss.Step != 2 {
		t.Errorf("Step = %d, want 2", ss.Step)
	}
	if ss.Total != 5 {
		t.Errorf("Total = %d, want 5", ss.Total)
	}
	if ss.StepName != "create-story" {
		t.Errorf("StepName = %s, want create-story", ss.StepName)
	}
	if ss.StoryKey != "7-1-test" {
		t.Errorf("StoryKey = %s, want 7-1-test", ss.StoryKey)
	}
	if ss.Model != "claude-opus-4-5" {
		t.Errorf("Model = %s, want claude-opus-4-5", ss.Model)
	}
	if ss.Operation != "Epic 3" {
		t.Errorf("Operation = %s, want Epic 3", ss.Operation)
	}
	if ss.Width != 80 {
		t.Errorf("Width = %d, want 80", ss.Width)
	}
}

func TestResult(t *testing.T) {
	r := Result{
		Success:  true,
		Duration: 5 * time.Second,
	}

	if !r.Success {
		t.Error("Success should be true")
	}
	if r.Duration != 5*time.Second {
		t.Errorf("Duration = %v, want 5s", r.Duration)
	}
}
