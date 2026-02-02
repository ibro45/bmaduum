# Batch Processing Recipes

Multi-story and epic processing patterns.

---

## Process Multiple Stories

Run full lifecycle for each story in sequence.

```bash
bmaduum story AUTH-042 AUTH-043 AUTH-044
```

Each story runs through its complete lifecycle. Skips stories already at `done`. Stops on first failure.

**Example output:**

```
─── Story 1 of 3: AUTH-042
  ... workflow output ...
Story AUTH-042 completed successfully

─── Story 2 of 3: AUTH-043
  ... workflow output ...
Story AUTH-043 completed successfully

All 3 stories processed
```

---

## Preview Batch Execution

See the full execution plan before running.

```bash
bmaduum story --dry-run AUTH-042 AUTH-043 AUTH-044
```

**Example output:**

```
Dry run for story AUTH-042:
  1. dev-story -> review
  2. code-review -> done
  3. git-commit -> done

Dry run for 2 stories:

Story AUTH-043:
  (already complete)

Story AUTH-044:
  1. create-story -> ready-for-dev
  2. dev-story -> review
  3. code-review -> done
  4. git-commit -> done

Total: 7 workflows across 2 stories (1 already complete)
```

---

## Process an Epic

Run all stories matching an epic pattern.

```bash
bmaduum epic 05
```

Finds stories matching `05-*` (e.g., `05-01-auth`, `05-02-dashboard`), sorts by story number, runs each through its complete lifecycle.

---

## Process Multiple Epics

Run multiple epics in sequence.

```bash
bmaduum epic 02 04 06
```

Processes epic 02, then 04, then 06. Each epic's stories are run to completion before moving to the next epic.

---

## Process All Active Epics

Auto-discover and process all epics with non-completed stories.

```bash
bmaduum epic all
```

---

## Preview Epic Execution

See what stories would be processed and their workflows.

```bash
# Single epic
bmaduum epic --dry-run 05

# Multiple epics
bmaduum epic --dry-run 02 04 06

# All active epics
bmaduum epic --dry-run all
```

Shows all matching stories and their remaining lifecycle steps.

---

## Handle Mixed-Status Batches

When processing stories with different statuses, each runs only its remaining workflows.

```bash
# Stories at different statuses:
# AUTH-042: ready-for-dev
# AUTH-043: review
# AUTH-044: done

bmaduum story AUTH-042 AUTH-043 AUTH-044
```

**Result:**

- AUTH-042: Runs dev-story -> code-review -> git-commit
- AUTH-043: Runs code-review -> git-commit
- AUTH-044: Skipped (already done)

Completed stories are always skipped, not re-run.
