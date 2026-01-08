---
phase: 02-workflow-router
plan: 01
subsystem: router
tags: [go, workflow, status-mapping, routing]

# Dependency graph
requires:
  - phase: 01-sprint-status-reader
    provides: Status type and constants (StatusBacklog, StatusReadyForDev, StatusInProgress, StatusReview, StatusDone)
provides:
  - GetWorkflow function mapping status to workflow name
  - ErrStoryComplete sentinel error for done stories
  - ErrUnknownStatus sentinel error for invalid status values
affects: [run-command, queue-command, epic-command]

# Tech tracking
tech-stack:
  added: []
  patterns: [table-driven-tests, switch-statement-routing, sentinel-errors]

key-files:
  created:
    - internal/router/router.go
    - internal/router/router_test.go
  modified: []

key-decisions:
  - "Package-level function instead of struct (pure mapping, no state needed)"
  - "Sentinel errors for done and unknown status (allows errors.Is() checks)"

patterns-established:
  - "Switch statement routing for status-to-workflow mapping"
  - "Table-driven tests covering all status constants and edge cases"

issues-created: []

# Metrics
duration: 3min
completed: 2026-01-08
---

# Phase 2 Plan 01: Status to Workflow Router Summary

**Created `internal/router` package with GetWorkflow function mapping story status to workflow names using switch statement routing**

## Performance

- **Duration:** 3 min
- **Started:** 2026-01-08T20:05:00Z
- **Completed:** 2026-01-08T20:08:00Z
- **Tasks:** 2 (TDD: RED → GREEN, no refactor needed)
- **Files modified:** 2

## Accomplishments

- GetWorkflow function routing status.Status to workflow name
- Sentinel errors: ErrStoryComplete (done stories), ErrUnknownStatus (invalid status)
- Complete test coverage: 7 GetWorkflow tests + 3 sentinel error tests (10 total)
- All mapping cases covered: backlog→create-story, ready-for-dev/in-progress→dev-story, review→code-review, done→error

## Task Commits

TDD plan commits (RED → GREEN cycle):

1. **RED: Failing tests** - `4498b41` (test)
   - Tests for all 5 status values
   - Tests for unknown and empty status
   - Tests for sentinel errors

2. **GREEN: Implementation** - `ae01adf` (feat)
   - Switch statement routing
   - Returns workflow name or appropriate error

**Refactor:** Not needed - implementation was already clean

**Plan metadata:** (pending - this commit)

## Files Created/Modified

- `internal/router/router.go` - GetWorkflow function, ErrStoryComplete, ErrUnknownStatus
- `internal/router/router_test.go` - TestGetWorkflow (7 cases), TestSentinelErrors (3 cases)

## Decisions Made

- Package-level function instead of struct (pure mapping with no state)
- Sentinel errors for errors.Is() compatibility in calling code
- Combined ready-for-dev and in-progress in single case (both use dev-story)

## Deviations from Plan

None - plan executed exactly as written.

## Issues Encountered

None

## Next Phase Readiness

- Router package complete and tested
- Ready for Phase 3 (Update Run Command) to integrate router.GetWorkflow
- All verification checks pass (test, lint)

---

_Phase: 02-workflow-router_
_Completed: 2026-01-08_
