# User Guide

A practical guide to using `bmaduum` for automating development workflows.

## Overview

`bmaduum` is a CLI tool that orchestrates Claude AI to automate repetitive development tasks. It handles:

- Creating story definitions from story keys
- Implementing features based on story requirements
- Running code reviews
- Committing and pushing changes
- Processing multiple stories in batch
- Managing entire epics

## Quick Start

### Prerequisites

1. **Go 1.25+** installed
2. **Claude CLI** installed and configured ([installation guide](https://github.com/anthropics/claude-code))
3. **just** command runner (optional but recommended)

### Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/bmaduum.git
cd bmaduum

# Build the binary
just build
# OR
go build -o bmaduum ./cmd/bmaduum

# (Optional) Install globally
go install ./cmd/bmaduum
```

### Verify Installation

```bash
bmaduum --help
```

## Basic Usage

### Recommended: Status-Based Automation

The `story` command automatically runs the appropriate workflows based on the story's current status:

```bash
# Run full lifecycle for a single story
bmaduum story 6-1-setup-project

# Run multiple stories
bmaduum story 6-1-auth 6-2-dashboard 6-3-tests
```

The `story` command reads the story's current status and executes all remaining workflows until completion:

| Story Status    | Remaining Lifecycle                                            |
| --------------- | -------------------------------------------------------------- |
| `backlog`       | create-story -> dev-story -> code-review -> git-commit -> done |
| `ready-for-dev` | dev-story -> code-review -> git-commit -> done                 |
| `in-progress`   | dev-story -> code-review -> git-commit -> done                 |
| `review`        | code-review -> git-commit -> done                              |
| `done`          | No action (story already complete)                             |

### Processing Epics

Process entire epics or all active epics:

```bash
# Single epic
bmaduum epic 6

# Multiple epics
bmaduum epic 2 4 6

# All active epics
bmaduum epic all
```

This finds all stories matching the pattern `{epic-id}-{N}-*` (e.g., `6-1-auth`, `6-2-dashboard`), sorts them by story number, and runs each through its complete lifecycle.

### Individual Workflow Steps (Advanced)

For advanced use cases, you can run individual workflow steps directly:

```bash
bmaduum workflow create-story 6-1-setup-project
bmaduum workflow dev-story 6-1-setup-project
bmaduum workflow code-review 6-1-setup-project
bmaduum workflow git-commit 6-1-setup-project
```

Use these when:
- A workflow fails and you want to retry just that step
- You need to run a step out of sequence for debugging
- You're testing or developing workflow prompts

Most users should use `story` or `epic` commands instead.

### Ad-Hoc Prompts

Run any prompt directly:

```bash
bmaduum raw "List all TODO comments in the codebase"
bmaduum raw "What tests are missing coverage?"
```

## Configuration

### Config File Location

By default, configuration is loaded from `config/workflows.yaml`.

Override with environment variable:

```bash
export BMADUUM_CONFIG_PATH=/path/to/custom/config.yaml
bmaduum story 6-1-setup-project
```

### Customizing Workflows

Edit `config/workflows.yaml` to customize workflow prompts:

```yaml
workflows:
  create-story:
    prompt_template: |
      Create a detailed story definition for {{.StoryKey}}.
      Include acceptance criteria and technical requirements.
      Do not ask clarifying questions.

  dev-story:
    prompt_template: |
      Implement story {{.StoryKey}}.
      Follow existing code patterns.
      Run tests after each change.
      Do not ask questions - use best judgment.

  code-review:
    prompt_template: |
      Review changes for story {{.StoryKey}}.
      Check for:
      - Code quality issues
      - Missing tests
      - Security vulnerabilities
      Auto-fix all issues immediately.

  git-commit:
    prompt_template: |
      Commit changes for {{.StoryKey}}.
      Use conventional commit format.
      Push to current branch.
```

### Template Variables

| Variable        | Description                         |
| --------------- | ----------------------------------- |
| `{{.StoryKey}}` | The story key passed to the command |

### Output Settings

Control how much output is displayed:

```yaml
output:
  truncate_lines: 20 # Max lines for tool output
  truncate_length: 60 # Max chars for command headers
```

### Claude Settings

Customize Claude execution:

```yaml
claude:
  binary_path: claude # Path to Claude binary
  output_format: stream-json # Output format (don't change)
```

Or use environment variables:

```bash
export BMADUUM_CLAUDE_PATH=/usr/local/bin/claude
```

## Sprint Status File

### File Location

The tool reads story status from:

```
_bmad-output/implementation-artifacts/sprint-status.yaml
```

### File Format

```yaml
development_status:
  6-1-setup-project: ready-for-dev
  6-2-add-authentication: in-progress
  6-3-fix-bug: review
  6-4-documentation: done
```

### Valid Status Values

| Status          | Meaning                                 |
| --------------- | --------------------------------------- |
| `backlog`       | Not started, needs story creation       |
| `ready-for-dev` | Story created, ready for implementation |
| `in-progress`   | Currently being implemented             |
| `review`        | Implementation done, needs review       |
| `done`          | Complete                                |

## Workflow Patterns

### Pattern 1: Full Lifecycle (Recommended)

Let the tool handle the entire lifecycle:

```bash
bmaduum story 6-1-setup-project
```

### Pattern 2: Batch Processing

Process multiple stories:

```bash
bmaduum story 6-1-setup 6-2-auth 6-3-tests
```

### Pattern 3: Epic Processing

Process all stories in an epic:

```bash
bmaduum epic 6
```

### Pattern 4: All Active Epics

Process all active epics:

```bash
bmaduum epic all
```

### Pattern 5: Investigation

Use raw prompts for ad-hoc tasks:

```bash
# Understand the codebase
bmaduum raw "Explain the authentication flow"

# Find issues
bmaduum raw "What tests have the most failures?"

# Generate reports
bmaduum raw "Create a summary of recent changes"
```

## Dry Run Mode

Preview what workflows will run without executing them:

```bash
# Single story
bmaduum story --dry-run 6-1-setup-project

# Multiple stories
bmaduum story --dry-run 6-1-setup 6-2-auth

# Epic
bmaduum epic --dry-run 6

# All active epics
bmaduum epic --dry-run all
```

**Example output:**

```
Dry run for story 6-1-setup-project:
  1. create-story -> ready-for-dev
  2. dev-story -> review
  3. code-review -> done
  4. git-commit -> done
```

## Error Recovery

### State Persistence

When a workflow fails, the tool saves execution state to `.bmad-state.json`:

```json
{
	"story_key": "6-1-setup-project",
	"step_index": 2,
	"total_steps": 4,
	"start_status": "backlog"
}
```

### Resuming After Failure

```bash
# Workflow fails at step 2 (dev-story)
bmaduum story 6-1-setup-project
# Error: workflow failed: dev-story returned exit code 1

# Fix the issue, then re-run
bmaduum story 6-1-setup-project
# Continues from current status (no work is lost)
```

The state file is automatically cleared on successful completion.

### Force Fresh Start

Delete the state file to restart from the story's current status:

```bash
rm .bmad-state.json
bmaduum story 6-1-setup-project
```

## Understanding Output

### Tool Invocations

When Claude uses tools, you'll see formatted output:

```
┌─ Bash ─────────────────────────────────────────────────────────
│  List project files
│  $ ls -la
└────────────────────────────────────────────────────────────────

total 48
drwxr-xr-x  12 user  staff   384 Jan  8 10:00 .
drwxr-xr-x   5 user  staff   160 Jan  8 09:00 ..
...
```

### Progress Indicators

| Symbol | Meaning     |
| ------ | ----------- |
| ●      | In progress |
| ✓      | Success     |
| ✗      | Failure     |
| ○      | Skipped     |

### Completion Summary

After processing multiple stories:

```
─── Story 1 of 3: 6-1-setup-project
Story 6-1-setup-project completed successfully

─── Story 2 of 3: 6-2-add-authentication
Story 6-2-add-authentication completed successfully

─── Story 3 of 3: 6-3-fix-bug
Story 6-3-fix-bug is already complete, skipping

All 3 stories processed (2 completed, 1 skipped)
```

## Error Handling

### Exit Codes

| Code | Meaning                             |
| ---- | ----------------------------------- |
| 0    | Success                             |
| 1    | General error                       |
| N    | Claude's exit code (passed through) |

### Common Issues

**Claude not found:**

```
Error: failed to start claude: exec: "claude": executable file not found in $PATH
```

Solution: Install Claude CLI or set `BMADUUM_CLAUDE_PATH`.

**Config not found:**

```
Error: error loading config: open config/workflows.yaml: no such file or directory
```

Solution: Create the config file or set `BMADUUM_CONFIG_PATH`.

**Story not in status file:**

```
Error: story 9-99-not-found not found in sprint status
```

Solution: Add the story to `sprint-status.yaml`.

**Unknown status:**

```
Error: unknown status value: invalid-status
```

Solution: Use a valid status: `backlog`, `ready-for-dev`, `in-progress`, `review`, or `done`.

## Tips and Best Practices

### 1. Use Story Keys That Follow Your Convention

Story keys appear in commits and prompts. Use meaningful identifiers that match your project:

- BMAD-style: `6-1-setup`, `6-2-auth`, `6-3-tests`
- JIRA-style: `PROJ-123`, `PROJ-124`
- Descriptive: `feat-user-profile`, `bug-login-fix`

### 2. Customize Prompts for Your Project

Edit the prompt templates to match your project's conventions, coding standards, and requirements.

### 3. Keep Status File Updated

The status file is the source of truth for automation. Keep it current as stories progress.

### 4. Handle Failures Gracefully

When processing stops on failure:

1. Review the error
2. Fix the issue manually or adjust the story
3. Update the status file if needed
4. Re-run the command (completed workflows will be skipped)

### 5. Use Raw for Exploration

Before starting a story, use raw prompts to understand the codebase:

```bash
bmaduum raw "What files would I need to change to add user authentication?"
```

### 6. Use Dry Run Before Long Operations

Preview what will run before committing to a long-running operation:

```bash
bmaduum story --dry-run 6-1 6-2 6-3
bmaduum epic --dry-run all
```

## Command Reference Summary

### Main Commands

| Command | Purpose | Example |
|---------|---------|---------|
| `story` | Run full lifecycle for one or more stories | `bmaduum story 6-1-setup` |
| `epic` | Run all stories in one or more epics | `bmaduum epic 6` |
| `raw` | Run an arbitrary prompt | `bmaduum raw "explain this"` |
| `version` | Display version info | `bmaduum version` |

### Advanced Commands

| Command | Purpose | Example |
|---------|---------|---------|
| `workflow create-story` | Create story definition | `bmaduum workflow create-story 6-1` |
| `workflow dev-story` | Implement a story | `bmaduum workflow dev-story 6-1` |
| `workflow code-review` | Review code changes | `bmaduum workflow code-review 6-1` |
| `workflow git-commit` | Commit and push changes | `bmaduum workflow git-commit 6-1` |

### Global Flags

| Flag | Available On | Description |
|------|--------------|-------------|
| `--dry-run` | `story`, `epic` | Preview without executing |
| `--auto-retry` | `story`, `epic` | Retry on rate limit errors |
| `--tui` | `story` | Enable interactive TUI mode |
| `--auto-retry` | `workflow` subcommands | Retry on rate limit errors |
