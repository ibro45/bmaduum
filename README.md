# bmaduum

A CLI tool for automating [BMAD-METHOD](https://github.com/bmad-code-org/BMAD-METHOD) development workflows with Claude AI.

**bmaduum** orchestrates Claude AI to automate development workflows—creating stories, implementing features, reviewing code, and managing git operations based on your project's sprint status.

> **Warning:** This tool runs Claude CLI with `--dangerously-skip-permissions`, meaning Claude can read, write, and execute commands **without asking for confirmation**. Only use in trusted repositories and isolated environments.

## Installation

Requires [Claude CLI](https://github.com/anthropics/claude-code) installed and configured.

```bash
git clone https://github.com/ibro45/bmaduum.git
cd bmaduum
go install ./cmd/bmaduum
```

> **Note:** Installs to `~/go/bin`. Add to PATH: `export PATH="$HOME/go/bin:$PATH"`

## Usage

```bash
# Run a story through its full lifecycle
bmaduum story 6-1-setup-project

# Process multiple stories
bmaduum story 6-1 6-2 6-3

# Process an epic (all stories matching 6-*)
bmaduum epic 6

# Process all active epics
bmaduum epic all

# Preview without executing
bmaduum story --dry-run 6-1
bmaduum epic --dry-run all

# Run arbitrary prompt
bmaduum raw "List all Go files"
```

### Lifecycle

Stories progress through workflows based on their status in `sprint-status.yaml`:

| Status | Remaining Workflows |
|--------|---------------------|
| `backlog` | create-story → dev-story → code-review → git-commit |
| `ready-for-dev` | dev-story → code-review → git-commit |
| `in-progress` | dev-story → code-review → git-commit |
| `review` | code-review → git-commit |
| `done` | skipped |

### Flags

| Flag | Description |
|------|-------------|
| `--dry-run` | Preview workflows without executing |
| `--auto-retry` | Retry on rate limit errors |

## Configuration

Configuration is optional. Defaults work out of the box.

```bash
# Custom config file
export BMADUUM_CONFIG_PATH=./my-config.yaml

# Custom Claude binary path
export BMADUUM_CLAUDE_PATH=/opt/homebrew/bin/claude
```

See [docs/CLI_REFERENCE.md](docs/CLI_REFERENCE.md) for full configuration options.

## Sprint Status File

The tool reads story status from:

```
_bmad-output/implementation-artifacts/sprint-status.yaml
```

```yaml
development_status:
  6-1-setup-project: done
  6-2-add-feature: in-progress
  6-3-fix-bug: backlog
```

## Development

```bash
just build    # Build binary
just test     # Run tests
just lint     # Run linter
```

See [docs/DEVELOPMENT.md](docs/DEVELOPMENT.md) for details.

## License

MIT - see [LICENSE](LICENSE)
