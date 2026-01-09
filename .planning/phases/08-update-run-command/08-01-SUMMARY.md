---
phase: 08-update-run-command
plan: 01
subsystem: cli
tags: [go, cobra, lifecycle, tdd, interfaces]

# Dependency graph
requires:
  - phase: 07-story-lifecycle-executor
    provides: lifecycle.Executor for orchestrating full story workflows
provides:
  - run command executes full lifecycle instead of single workflow
  - App struct with interface-based dependency injection
  - StatusWriter integration
affects: [09-update-epic-command, 10-update-queue-command]

# Tech tracking
tech-stack:
  added: []
  patterns: [interface-based DI in App struct, TDD red-green-refactor]

key-files:
  created: []
  modified:
    - internal/cli/run.go
    - internal/cli/root.go
    - internal/cli/run_test.go
    - internal/workflow/queue.go

key-decisions:
  - "Changed App.Runner from *workflow.Runner to WorkflowRunner interface"
  - "Changed App.StatusReader from *status.Reader to StatusReader interface"
  - "Added StatusWriter interface and field to App struct"
  - "Removed obsolete single-workflow tests, replaced with lifecycle tests"

patterns-established:
  - "Interface-based dependency injection in App struct for testability"

issues-created: []

# Metrics
duration: 5min
completed: 2026-01-09
---

# Phase 8 Plan 1: Update Run Command Summary

**Run command updated to execute full story lifecycle using lifecycle.Executor with interface-based dependency injection**

## Performance

- **Duration:** 5 min
- **Started:** 2026-01-09T02:13:55Z
- **Completed:** 2026-01-09T02:18:46Z
- **Tasks:** 1 (TDD feature)
- **Files modified:** 4

## Accomplishments

- Run command now executes complete lifecycle from current status to done
- App struct uses interfaces (WorkflowRunner, StatusReader, StatusWriter) for testability
- Comprehensive table-driven tests verify all lifecycle scenarios
- Status updated in sprint-status.yaml after each workflow

## TDD Commits

**RED:** `2fd64c9` - test(08-01): add failing tests for run command lifecycle execution

- Tests for full lifecycle from backlog (4 workflows)
- Tests for partial lifecycle from ready-for-dev (3 workflows)
- Tests for partial lifecycle from review (2 workflows)
- Tests for done story (no workflows, message printed)
- Tests for workflow failure mid-lifecycle

**GREEN:** `7456ede` - feat(08-01): implement run command with lifecycle executor

- Updated run.go to use lifecycle.Executor
- Added StatusWriter to App struct
- Changed Runner and StatusReader to interfaces
- Updated help text to describe full lifecycle behavior

**REFACTOR:** No separate refactor commit needed - code was clean after GREEN phase

## Files Created/Modified

- `internal/cli/run.go` - Changed from single-workflow routing to lifecycle.Executor
- `internal/cli/root.go` - Added interfaces and StatusWriter field
- `internal/cli/run_test.go` - New lifecycle tests, removed obsolete tests
- `internal/workflow/queue.go` - Changed StatusReader param to interface type

## Decisions Made

- Changed App fields from concrete types to interfaces for testability
- Defined interfaces at the consumer (cli package) following Go conventions
- Removed obsolete single-workflow tests as they no longer test actual behavior

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 3 - Blocking] Updated workflow.QueueRunner.RunQueueWithStatus signature**

- **Found during:** GREEN phase
- **Issue:** Other commands (queue, epic) failed to compile because StatusReader was changed to interface
- **Fix:** Changed `*status.Reader` to `StatusReader` interface in queue.go
- **Files modified:** internal/workflow/queue.go
- **Verification:** Build passes, all tests pass
- **Committed in:** 7456ede

### Deferred Enhancements

None

---

**Total deviations:** 1 auto-fixed (blocking), 0 deferred
**Impact on plan:** Blocking fix was necessary for compilation. No scope creep.

## Issues Encountered

None - implementation went smoothly.

## Next Phase Readiness

- Run command fully functional with lifecycle execution
- Ready for Phase 9: Update Epic Command to use lifecycle executor
- Pattern established for updating queue command in Phase 10

---

_Phase: 08-update-run-command_
_Completed: 2026-01-09_
