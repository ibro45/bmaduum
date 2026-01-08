package cli

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewExitError(t *testing.T) {
	tests := []struct {
		code     int
		expected string
	}{
		{0, "exit status 0"},
		{1, "exit status 1"},
		{127, "exit status 127"},
		{255, "exit status 255"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			err := NewExitError(tt.code)
			assert.Equal(t, tt.code, err.Code)
			assert.Equal(t, tt.expected, err.Error())
		})
	}
}

func TestIsExitError(t *testing.T) {
	t.Run("with ExitError", func(t *testing.T) {
		err := NewExitError(42)
		code, ok := IsExitError(err)
		assert.True(t, ok)
		assert.Equal(t, 42, code)
	})

	t.Run("with standard error", func(t *testing.T) {
		err := errors.New("standard error")
		code, ok := IsExitError(err)
		assert.False(t, ok)
		assert.Equal(t, 0, code)
	})

	t.Run("with nil", func(t *testing.T) {
		code, ok := IsExitError(nil)
		assert.False(t, ok)
		assert.Equal(t, 0, code)
	})
}

func TestExitError_ImplementsError(t *testing.T) {
	var _ error = (*ExitError)(nil)
}
