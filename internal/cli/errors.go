package cli

import "fmt"

// ExitError represents a command execution failure with an exit code.
// It implements the error interface and is used by RunE functions
// to signal non-zero exit without calling os.Exit() directly.
type ExitError struct {
	Code int
}

// Error implements the error interface.
func (e *ExitError) Error() string {
	return fmt.Sprintf("exit status %d", e.Code)
}

// NewExitError creates an ExitError with the given exit code.
func NewExitError(code int) *ExitError {
	return &ExitError{Code: code}
}

// IsExitError checks if an error is an ExitError and returns its code.
// Returns (code, true) if it's an ExitError, (0, false) otherwise.
func IsExitError(err error) (int, bool) {
	if exitErr, ok := err.(*ExitError); ok {
		return exitErr.Code, true
	}
	return 0, false
}
