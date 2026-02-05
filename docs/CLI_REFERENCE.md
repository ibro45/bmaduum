# CLI Reference

Complete command-line interface reference for `bmaduum`.

## Synopsis

```
bmaduum [command] [arguments] [flags]
```

## Description

BMAD Automation CLI orchestrates Claude AI to run development workflows including story creation, implementation, code review, and git operations.

## Global Behavior

All commands:

- Load configuration from `config/workflows.yaml` (or `BMADUUM_CONFIG_PATH`)
- Execute Claude CLI with `--dangerously-skip-permissions` and `--output-format stream-json`
- Display styled terminal output with progress indicators
- Return appropriate exit codes (0 for success, non-zero for failure)

---

## Commands

### story

Run full lifecycle for one or more stories from their current status to done.

**Usage:**

```bash
bmaduum story [--dry-run] [--auto-retry] <story-key> [story-key...]
```

**Arguments:**
| Argument | Required | Description |
|----------|----------|-------------|
| story-key | Yes (1+) | One or more story identifiers |

**Flags:**
| Flag | Description |
|------|-------------|
| `--dry-run` | Preview workflow sequence without execution |
| `--auto-retry` | Automatically retry on rate limit errors |

**Examples:**

```bash
# Run full lifecycle for a single story
bmaduum story 6-1-setup-project

# Run full lifecycle for multiple stories
bmaduum story 6-1-setup 6-2-auth 6-3-tests

# Preview what would run
bmaduum story --dry-run 6-1-setup 6-2-auth
```

**Behavior:**

1. Processes each story through its **full lifecycle** to completion
2. Auto-updates status after each successful workflow step
3. Skips stories with status `done`
4. Stops on first failure
5. For multiple stories, shows progress indicators

**Lifecycle Routing:**

| Story Status    | Remaining Lifecycle                                            |
| --------------- | -------------------------------------------------------------- |
| `backlog`       | create-story -> dev-story -> code-review -> git-commit -> done |
| `ready-for-dev` | dev-story -> code-review -> git-commit -> done                 |
| `in-progress`   | dev-story -> code-review -> git-commit -> done                 |
| `review`        | code-review -> git-commit -> done                              |
| `done`          | No action (story already complete)                             |

---

### epic

Run full lifecycle for all stories in one or more epics, or all active epics.

**Usage:**

```bash
# Single or multiple epics
bmaduum epic [--dry-run] [--auto-retry] <epic-id> [epic-id...]

# All active epics
bmaduum epic [--dry-run] [--auto-retry] all
```

**Arguments:**
| Argument | Required | Description |
|----------|----------|-------------|
| epic-id | Yes (1+) | One or more epic identifiers, or `all` for all active epics |

**Flags:**
| Flag | Description |
|------|-------------|
| `--dry-run` | Preview workflow sequence without execution |
| `--auto-retry` | Automatically retry on rate limit errors |

**Examples:**

```bash
# Run full lifecycle for all stories in a single epic
bmaduum epic 6

# Run multiple epics
bmaduum epic 2 4 6

# Run all active epics
bmaduum epic all

# Preview what would run
bmaduum epic --dry-run 2 4 6
bmaduum epic --dry-run all
```

**Story Discovery:**

Stories are discovered from `sprint-status.yaml` using the pattern:

```
{epic-id}-{story-number}-*
```

For epic `6`, this matches:

- `6-1-implement-auth`
- `6-2-add-dashboard`
- `6-3-fix-navigation`

Stories are sorted by story number and processed in order.

**When using `all`:**

The `all` argument auto-discovers all epics with non-completed stories and processes them in numerical order.

---

### workflow (Advanced)

Run individual BMAD workflow steps directly. These are the same workflow commands used in BMAD-METHOD and are automatically executed by `story` and `epic` commands.

**Usage:**

```bash
bmaduum workflow <workflow-name> <story-key>
bmaduum workflow <subcommand> --help
```

**Available workflows:**
- `create-story`: Create a story definition from backlog
- `dev-story`: Implement a story (ready-for-dev or in-progress)
- `code-review`: Review code changes (review status)
- `git-commit`: Commit and push changes after review

**Subcommands:**
| Subcommand | Description |
|------------|-------------|
| `create-story` | Create a story definition from backlog |
| `dev-story` | Implement a story through development |
| `code-review` | Review code changes for a story |
| `git-commit` | Commit and push changes for a story |

**Flags:**
| Flag | Description |
|------|-------------|
| `--auto-retry` | Automatically retry on rate limit errors |

**Examples:**

```bash
# Using parent command syntax
bmaduum workflow create-story 6-1-setup
bmaduum workflow dev-story 6-1-setup
bmaduum workflow code-review 6-1-setup
bmaduum workflow git-commit 6-1-setup
```

**When to use:**

- A workflow fails and you want to retry just that step
- You need to run a step out of the normal sequence for debugging
- You're testing or developing workflow prompts

**Note:** Most users should use `story` or `epic` commands instead.

---

### raw

Execute an arbitrary prompt with Claude.

**Usage:**

```bash
bmaduum raw <prompt>
```

**Arguments:**
| Argument | Required | Description |
|----------|----------|-------------|
| prompt | Yes | The prompt text (can be multiple words) |

**Example:**

```bash
bmaduum raw "List all Go files in the project"
bmaduum raw Explain the architecture of this codebase
```

**Behavior:**

1. Joins all arguments into a single prompt
2. Executes Claude directly with the prompt
3. Does not use any workflow templates

---

### version

Display version information.

**Usage:**

```bash
bmaduum version
```

**Example:**

```bash
bmaduum version
# Output: bmaduum 1.0.0 (build: abc1234)
```

---

## Exit Codes

| Code | Meaning                                              |
| ---- | ---------------------------------------------------- |
| 0    | Success                                              |
| 1    | General error (config load failure, unknown command) |
| N    | Claude exit code (passed through from Claude CLI)    |

---

## Environment Variables

| Variable           | Description                | Default                   |
| ------------------ | -------------------------- | ------------------------- |
| `BMADUUM_CONFIG_PATH` | Path to configuration file | `./config/workflows.yaml` |
| `BMADUUM_CLAUDE_PATH` | Path to claude command/binary | `claude` (from PATH)  |

---

## Configuration File

The default configuration file is `config/workflows.yaml`:

```yaml
workflows:
  create-story:
    prompt_template: "Create story: {{.StoryKey}}"

  dev-story:
    prompt_template: "Work on story: {{.StoryKey}}"

  code-review:
    prompt_template: "Review story: {{.StoryKey}}"

  git-commit:
    prompt_template: "Commit changes for {{.StoryKey}}"

full_cycle:
  steps:
    - create-story
    - dev-story
    - code-review
    - git-commit

claude:
  output_format: stream-json
  binary_path: claude

output:
  truncate_lines: 20 # Max lines to show for tool output
  truncate_length: 60 # Max chars for command header
```

### Template Variables

| Variable        | Description                         |
| --------------- | ----------------------------------- |
| `{{.StoryKey}}` | The story key passed to the command |

---

## Sprint Status File

The `story` and `epic` commands read story status from:

```
_bmad-output/implementation-artifacts/sprint-status.yaml
```

**Format:**

```yaml
development_status:
  6-1-setup-project: ready-for-dev
  6-2-add-authentication: in-progress
  6-3-fix-bug: review
  6-4-documentation: done
```

**Valid Status Values:**

- `backlog` - Story not yet started
- `ready-for-dev` - Story ready for implementation
- `in-progress` - Story being implemented
- `review` - Story in code review
- `done` - Story complete

---

## State File

The lifecycle executor persists execution state for error recovery.

**Location:**

```
.bmad-state.json   # In working directory (hidden file)
```

**Format:**

```json
{
	"story_key": "6-1-setup-project",
	"step_index": 2,
	"total_steps": 4,
	"start_status": "backlog"
}
```

**Fields:**
| Field | Description |
|-------|-------------|
| `story_key` | The story being processed |
| `step_index` | 0-based index of the current/failed step |
| `total_steps` | Total steps in the lifecycle sequence |
| `start_status` | The story's status when execution began |

**Lifecycle:**

1. **Saved on failure** - State is written when a workflow step fails
2. **Used on resume** - On re-run, execution continues from current status
3. **Cleared on success** - State file is deleted after successful lifecycle completion

---

## Examples

### Status-Based Automation (Recommended)

```bash
# Run full lifecycle for a single story
bmaduum story 6-1-setup-project

# Process multiple stories
bmaduum story 6-1-setup 6-2-auth 6-3-tests

# Process an entire epic
bmaduum epic 6

# Process all active epics
bmaduum epic all
```

### Individual Workflow Steps (Advanced)

```bash
# Run a specific workflow step
bmaduum workflow create-story 6-1-setup
bmaduum workflow dev-story 6-1-setup
bmaduum workflow code-review 6-1-setup
bmaduum workflow git-commit 6-1-setup
```

### Ad-Hoc Tasks

```bash
# Run arbitrary prompts
bmaduum raw "What is the test coverage?"
bmaduum raw "Find all TODO comments"
```

### Custom Configuration

```bash
# Use custom config file
BMADUUM_CONFIG_PATH=/path/to/config.yaml bmaduum story 6-1-setup

# Use custom Claude binary
BMADUUM_CLAUDE_PATH=/usr/local/bin/claude bmaduum workflow dev-story 6-1-setup
```
