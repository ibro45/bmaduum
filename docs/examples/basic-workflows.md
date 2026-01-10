# Basic Workflow Recipes

Single-story commands for step-by-step development.

---

## Create a Story Definition

Generate story requirements and acceptance criteria from a story key.

```bash
bmad-automate create-story AUTH-042
```

Creates the story definition file. Status moves from `backlog` to `ready-for-dev`.

---

## Implement a Story

Run the development workflow to implement a feature.

```bash
bmad-automate dev-story AUTH-042
```

Claude implements the story, runs tests after each change. Status moves to `review` on success.

---

## Review Code Changes

Run code review on story changes with auto-fix.

```bash
bmad-automate code-review AUTH-042
```

Claude reviews for quality issues, security vulnerabilities, and missing tests. Automatically applies fixes.

---

## Commit and Push

Create a conventional commit and push to remote.

```bash
bmad-automate git-commit AUTH-042
```

Claude creates a commit with conventional format (e.g., `feat(auth): add login endpoint`), then pushes to the current branch.

---

## Run Ad-hoc Prompts

Execute arbitrary prompts for investigation or exploration.

```bash
# Explore the codebase
bmad-automate raw "What files handle user authentication?"

# Find issues
bmad-automate raw "List all TODO comments in this project"

# Generate reports
bmad-automate raw "Summarize test coverage by package"
```

Useful for understanding code before starting work or investigating issues.
