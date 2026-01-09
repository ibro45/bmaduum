---
phase: 12-dry-run-mode
plan: 01
type: tdd-summary
subsystem: lifecycle
tags: [executor, dry-run, preview]
started: 2026-01-08T00:00:00Z
completed: 2026-01-08T00:00:00Z
---

# Phase 12 Plan 01: GetSteps Method Summary

**Added GetSteps method to lifecycle executor enabling dry-run preview of workflow steps without execution.**

## RED Phase

Wrote table-driven tests in `executor_test.go` for `TestGetSteps` covering:

- Story in backlog returns all 4 steps (create-story, dev-story, code-review, git-commit)
- Story in ready-for-dev returns 3 steps (dev-story, code-review, git-commit)
- Story in review returns 2 steps (code-review, git-commit)
- Story already done returns `router.ErrStoryComplete`
- Status read error propagates correctly

Tests also verify that GetSteps does NOT:

- Execute any workflows (runner.Calls is empty)
- Update any status (writer.Calls is empty)

Tests failed with: `executor.GetSteps undefined (type *Executor has no field or method GetSteps)`

## GREEN Phase

Implemented `GetSteps` method in `executor.go`:

```go
func (e *Executor) GetSteps(storyKey string) ([]router.LifecycleStep, error) {
    currentStatus, err := e.statusReader.GetStoryStatus(storyKey)
    if err != nil {
        return nil, err
    }

    steps, err := router.GetLifecycle(currentStatus)
    if err != nil {
        return nil, err
    }

    return steps, nil
}
```

The method reuses existing logic patterns from `Execute` but stops before the execution loop.

## REFACTOR Phase

None needed. The implementation is minimal and clear. While there is some code similarity between `Execute` and `GetSteps` in the first two steps (reading status and getting lifecycle), extracting a helper method would add complexity for minimal benefit. The current structure maintains clarity.

## Commits

| Phase | Hash      | Message                                            |
| ----- | --------- | -------------------------------------------------- |
| RED   | `f1ae09f` | test(12-01): add failing tests for GetSteps method |
| GREEN | `b698ab4` | feat(12-01): implement GetSteps method             |

## Files Modified

- `internal/lifecycle/executor.go` - Added GetSteps method (18 lines)
- `internal/lifecycle/executor_test.go` - Added TestGetSteps tests (93 lines)

## Verification

```
$ go test ./internal/lifecycle/... -v
=== RUN   TestNewExecutor
--- PASS: TestNewExecutor (0.00s)
=== RUN   TestExecute
--- PASS: TestExecute (0.00s)
=== RUN   TestGetSteps
=== RUN   TestGetSteps/story_in_backlog_returns_all_4_steps
=== RUN   TestGetSteps/story_in_ready-for-dev_returns_3_steps
=== RUN   TestGetSteps/story_in_review_returns_2_steps
=== RUN   TestGetSteps/story_already_done_returns_ErrStoryComplete
=== RUN   TestGetSteps/get_status_error_propagates
--- PASS: TestGetSteps (0.00s)
PASS
ok      bmad-automate/internal/lifecycle        0.225s
```

## Deviations

None.

## Next Step

Ready for 12-02-PLAN.md: Add --dry-run flag to commands
