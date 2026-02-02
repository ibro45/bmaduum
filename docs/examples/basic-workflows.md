# Basic Workflow Recipes

Single-story commands for step-by-step development.

**Note:** These are advanced commands. Most users should use `story` or `epic` commands which automatically run the appropriate workflows based on story status.

---

## Create a Story Definition

Generate story requirements and acceptance criteria from a story key.

```bash
bmaduum workflow create-story AUTH-042
```

Creates the story definition file. Status moves from `backlog` to `ready-for-dev`.

---

## Implement a Story

Run the development workflow to implement a feature.

```bash
bmaduum workflow dev-story AUTH-042
```

Claude implements the story, runs tests after each change. Status moves to `review` on success.

---

## Review Code Changes

Run code review on story changes with auto-fix.

```bash
bmaduum workflow code-review AUTH-042
```

Claude reviews for quality issues, security vulnerabilities, and missing tests. Automatically applies fixes.

---

## Commit and Push

Create a conventional commit and push to remote.

```bash
bmaduum workflow git-commit AUTH-042
```

Claude creates a commit with conventional format (e.g., `feat(auth): add login endpoint`), then pushes to the current branch.

---

## Run Ad-hoc Prompts

Execute arbitrary prompts for investigation or exploration.

```bash
# Explore the codebase
bmaduum raw "What files handle user authentication?"

# Find issues
bmaduum raw "List all TODO comments in this project"

# Generate reports
bmaduum raw "Summarize test coverage by package"
```

Useful for understanding code before starting work or investigating issues.

---

## Recommended: Status-Based Automation

Instead of running individual workflows, use the `story` command which automatically executes the right workflows based on story status:

```bash
# Run full lifecycle for a story
bmaduum story AUTH-042

# Run full lifecycle for multiple stories
bmaduum story AUTH-042 AUTH-043 AUTH-044

# Run an entire epic
bmaduum epic 05

# Run all active epics
bmaduum epic all
```
