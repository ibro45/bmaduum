# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Development Commands

```bash
just build              # Build binary to ./bmaduum
just test               # Run all tests
just test-verbose       # Run tests with verbose output
just test-pkg ./internal/claude  # Test specific package
just test-coverage      # Generate coverage.html
just lint               # Run golangci-lint
just check              # Run fmt, vet, and test
just run --help         # Build and run with arguments
just release-snapshot   # Build release locally (snapshot)
just release            # Full release with GoReleaser
```

## Architecture

This is a CLI tool that orchestrates Claude CLI to run automated development workflows. It spawns Claude as a subprocess, parses its streaming JSON output, and displays formatted results.

### Package Dependencies

```
cmd/bmaduum/main.go
         |
         v
    internal/cli (Cobra commands)
         |
         ├──> internal/workflow (orchestration)
         |         |
         |         ├──> internal/claude (Claude execution + JSON parsing)
         |         |
         |         └──> internal/output (terminal formatting)
         |
         ├──> internal/lifecycle (lifecycle orchestration)
         |         |
         |         ├──> internal/router (GetLifecycle for step sequences)
         |         |
         |         └──> internal/state (execution state persistence)
         |
         ├──> internal/ratelimit (rate limit detection from stderr)
         |
         ├──> internal/status (sprint status reading)
         |
         ├──> internal/router (workflow routing)
         |
         └──> internal/config (Viper configuration)
```

### Key Interfaces for Testing

- **`claude.Executor`** - Interface for running Claude CLI. Use `MockExecutor` in tests to avoid spawning real processes.
- **`output.Printer`** - Interface for terminal output. Use `NewPrinterWithWriter(buf)` to capture output in tests.

### Data Flow

1. CLI command receives story key
2. `config.Config.GetPrompt()` expands Go template with `{{.StoryKey}}`
3. `workflow.Runner` calls `claude.Executor.ExecuteWithResult()`
4. `claude.Parser` reads streaming JSON, emits `Event` structs
5. `output.Printer` formats and displays events
6. Lifecycle executor runs multiple workflows with state persistence

### Configuration

Workflow prompts are in `config/workflows.yaml` using Go templates. Config loads via Viper with env var overrides (`BMADUUM_` prefix).

### Claude CLI Integration

The executor always passes `--dangerously-skip-permissions` and `--output-format stream-json`. Each JSON line from stdout is parsed into `StreamEvent` structs, then converted to the higher-level `Event` type with convenience methods (`IsText()`, `IsToolUse()`, `IsToolResult()`).

### Rate Limit Detection

The `internal/ratelimit` package provides rate limit detection from Claude CLI's stderr output:

- `Detector` identifies rate limit error messages and extracts reset times
- `State` provides thread-safe rate limit state management
- Used in conjunction with `--auto-retry` flag for automatic retry with intelligent wait times

### State Persistence

The `internal/state` package enables resume functionality when lifecycle execution fails:

- State is saved to `.bmad-state.json` in the working directory
- Atomic writes (temp file + rename) ensure crash safety
- Automatically cleared on successful completion
- Tracks story key, step index, total steps, and start status
