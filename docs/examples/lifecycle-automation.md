# Lifecycle Automation Recipes

Full lifecycle execution from current status to done.

---

## Run Full Lifecycle From Any Status

Execute all remaining workflows until the story is complete.

```bash
bmaduum story AUTH-042
```

Reads current status, runs remaining workflows in sequence, auto-updates status after each step. Stops at `done` or on first failure.

**Example output:**

```
[1/4] create-story
  ... workflow output ...

[2/4] dev-story
  ... workflow output ...

[3/4] code-review
  ... workflow output ...

[4/4] git-commit
  ... workflow output ...

Story AUTH-042 completed successfully
```

---

## Run Multiple Stories

Process multiple stories in sequence:

```bash
bmaduum story AUTH-042 AUTH-043 AUTH-044
```

Each story runs to completion before moving to the next. Done stories are skipped.

---

## Preview Workflow Sequence

See what would run without executing anything.

```bash
bmaduum story --dry-run AUTH-042
```

**Example output:**

```
Dry run for story AUTH-042:
  1. create-story -> ready-for-dev
  2. dev-story -> review
  3. code-review -> done
  4. git-commit -> done
```

Useful before committing to a long-running lifecycle.

---

## Process Entire Epics

Run all stories in an epic:

```bash
# Single epic
bmaduum epic 05

# Multiple epics
bmaduum epic 02 04 06

# All active epics
bmaduum epic all
```

Stories are discovered using the pattern `{epic-id}-{story-number}-*` and processed in numerical order.

---

## Resume After Failure

When a workflow fails, the tool saves state to `.bmad-state.json`. Re-run to continue from current status.

```bash
# First run fails at dev-story
bmaduum story AUTH-042
# Error: workflow failed: dev-story returned exit code 1

# Fix the issue, then re-run
bmaduum story AUTH-042
# Continues from current status (in-progress)
```

State file is cleared automatically on successful completion.

---

## Force Fresh Start

Delete the state file to restart from the story's current status.

```bash
rm .bmad-state.json
bmaduum story AUTH-042
```

Useful when you want to discard partial progress and start over.
