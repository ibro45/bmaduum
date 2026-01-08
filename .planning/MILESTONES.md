# Project Milestones: BMAD Automate

## v1.0 Status-Based Workflow Routing (Shipped: 2026-01-08)

**Delivered:** Automatic workflow routing based on sprint-status.yaml, eliminating manual workflow selection.

**Phases completed:** 1-5 (5 plans total)

**Key accomplishments:**

- Sprint status reader package parsing YAML with Status type and validation
- Workflow router mapping status values to workflow names (backlog→create-story, ready-for-dev/in-progress→dev-story, review→code-review)
- Run command with automatic status-based workflow routing
- Queue command with status-based routing and done-story skipping
- New epic command for batch-running all stories in an epic with numeric sorting

**Stats:**

- 29 files created/modified (+2,636 lines, -249 lines)
- 4,951 lines of Go total
- 5 phases, 5 plans, 10 tasks
- Same-day completion (2026-01-08)

**Git range:** `docs(01)` → `docs(05-01)`

**What's next:** TBD - milestone complete

---
