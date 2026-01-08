# BMAD Automate - Status-Based Workflow Routing

## What This Is

A CLI tool that orchestrates Claude CLI to run automated development workflows. Now enhanced with automatic story status detection from sprint-status.yaml for workflow routing, plus an `epic` command for batch execution.

## Core Value

Eliminate manual workflow selection by automatically routing stories to the correct workflow based on their status in sprint-status.yaml.

## Current State (v1.0)

**Shipped:** 2026-01-08

Status-based workflow routing is complete:

- `run <story>` — Automatically routes to correct workflow based on status
- `queue <story>...` — Routes each story, skips done, fails fast
- `epic <epic-id>` — Runs all stories for an epic in numeric order

Tech stack: Go, Cobra, Viper, yaml.v3
Codebase: 4,951 LOC Go

## Requirements

### Validated

- ✓ CLI command structure with Cobra — existing
- ✓ Configuration via Viper with YAML and env vars — existing
- ✓ Claude CLI subprocess execution with streaming JSON — existing
- ✓ Event-driven output parsing — existing
- ✓ Terminal formatting with Lipgloss — existing
- ✓ Commands: create-story, dev-story, code-review, git-commit, run, queue, epic, raw — v1.0
- ✓ Interface-based design for testability (Executor, Printer) — existing
- ✓ Go template expansion for prompts — existing
- ✓ Status-based workflow routing from sprint-status.yaml — v1.0
- ✓ Run command auto-routing based on status — v1.0
- ✓ Queue command with status routing and done-skip — v1.0
- ✓ Epic command for batch execution with numeric sorting — v1.0
- ✓ Fail-fast on story failure — v1.0

### Active

(None currently — v1.0 milestone complete)

### Out of Scope

- Manual workflow override flag — status always determines workflow
- Updating sprint-status.yaml after story completion — read-only access
- Parallel story execution — sequential only
- Epic status auto-transitions — only story-level routing

## Context

**Sprint Status File:**

- Always located at `_bmad-output/implementation-artifacts/sprint-status.yaml`
- YAML format with `development_status` section
- Story keys follow pattern: `{epic#}-{story#}-{description}` (e.g., `7-1-define-schema`)
- Statuses: `backlog`, `ready-for-dev`, `in-progress`, `review`, `done`

**Workflow Mapping:**
| Status | Workflow |
|--------|----------|
| `backlog` | `/bmad:bmm:workflows:create-story` |
| `ready-for-dev` | `/bmad:bmm:workflows:dev-story` |
| `in-progress` | `/bmad:bmm:workflows:dev-story` |
| `review` | `/bmad:bmm:workflows:code-review` |

## Constraints

- **Tech Stack**: Go with existing Cobra/Viper patterns
- **File Location**: Sprint status always at `_bmad-output/implementation-artifacts/sprint-status.yaml`
- **Claude CLI**: Requires Claude CLI installed and in PATH

## Key Decisions

| Decision                             | Rationale                                             | Outcome |
| ------------------------------------ | ----------------------------------------------------- | ------- |
| Auto-detect only, no manual override | Simplicity — status is source of truth                | ✓ Good  |
| Stop on first failure in epic        | Allows investigation before continuing                | ✓ Good  |
| Sequential execution only            | Stories may have dependencies                         | ✓ Good  |
| Read-only sprint-status.yaml access  | Separation of concerns — status managed elsewhere     | ✓ Good  |
| yaml.v3 instead of Viper for status  | Simpler for single file with known structure          | ✓ Good  |
| Package-level router function        | Pure mapping with no state needed                     | ✓ Good  |
| StatusReader injected via App struct | Testability — allows mock injection                   | ✓ Good  |
| Done stories skipped in queue        | Allows mixed-status batches without failure           | ✓ Good  |
| Epic reuses QueueRunner              | DRY — inherits all routing, skip, and fail-fast logic | ✓ Good  |

---

_Last updated: 2026-01-08 after v1.0 milestone_
