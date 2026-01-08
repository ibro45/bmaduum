---
phase: 03-update-run-command
plan: 01
subsystem: cli
tags: [cobra, status-routing, yaml]

# Dependency graph
requires:
  - phase: 01-sprint-status-reader
    provides: status.Reader for reading sprint-status.yaml
  - phase: 02-workflow-router
    provides: router.GetWorkflow for status-to-workflow mapping
provides:
  - status-based routing in run command
  - automatic workflow selection based on story status
affects: [04-update-queue-command, 05-epic-command]

# Tech tracking
tech-stack:
  added: []
  patterns: [status-reader-injection, error-sentinel-handling]

key-files:
  created: [internal/cli/run_test.go]
  modified:
    [internal/cli/run.go, internal/cli/root.go, internal/cli/cli_test.go]

key-decisions:
  - "StatusReader injected via App struct rather than created inline"
  - "Done status returns success (exit 0) with message, not error"

patterns-established:
  - "Status reader injection via App struct for testability"

issues-created: []

# Metrics
duration: 3min
completed: 2026-01-08
---

# Phase 3 Plan 01: Update Run Command Summary

**Run command now routes to appropriate workflow based on story status from sprint-status.yaml**

## Performance

- **Duration:** 3 min
- **Started:** 2026-01-08T20:15:09Z
- **Completed:** 2026-01-08T20:18:39Z
- **Tasks:** 2
- **Files modified:** 4

## Accomplishments

- Run command reads story status from sprint-status.yaml
- Automatic routing: backlog→create-story, ready-for-dev/in-progress→dev-story, review→code-review
- Done stories return success with completion message
- Updated command description to reflect new behavior

## Task Commits

Each task was committed atomically:

1. **Task 1: Update run command to use status-based routing** - `5ba38ac` (feat)
2. **Task 2: Add tests for status-based routing** - `2794f27` (test)

**Plan metadata:** TBD (docs: complete plan)

## Files Created/Modified

- `internal/cli/run.go` - Status-based routing logic
- `internal/cli/root.go` - Added StatusReader to App struct
- `internal/cli/run_test.go` - Tests for status-based routing (new)
- `internal/cli/cli_test.go` - Updated setup functions with StatusReader

## Decisions Made

- StatusReader injected via App struct for testability (consistent with other dependencies)
- Done status treated as success, not error (story is complete, nothing to do is not a failure)

## Deviations from Plan

None - plan executed exactly as written.

## Issues Encountered

None

## Next Phase Readiness

- Run command updated and fully tested
- Pattern established for queue command update in Phase 4
- StatusReader and router modules proven in integration

---

_Phase: 03-update-run-command_
_Completed: 2026-01-08_
