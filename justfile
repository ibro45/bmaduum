# Use PowerShell on Windows
set windows-shell := ["powershell.exe", "-NoLogo", "-Command"]

# Binary name
binary_name := "bmaduum"

# Default recipe - show help
default:
    @just --list

# Build the application
build:
    go build -o {{binary_name}} ./cmd/bmaduum

# Install the binary to $GOPATH/bin
install:
    go install ./cmd/bmaduum

# Run all tests
test:
    go test ./...

# Run tests with coverage report
test-coverage:
    go test -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html
    @echo "Coverage report generated: coverage.html"

# Run tests with verbose output
test-verbose:
    go test -v ./...

# Run tests for a specific package (e.g., just test-pkg ./internal/claude)
test-pkg pkg:
    go test -v {{pkg}}

# Clean build artifacts
[unix]
clean:
    rm -f {{binary_name}}
    rm -f coverage.out coverage.html

[windows]
clean:
    Remove-Item -Force -ErrorAction SilentlyContinue {{binary_name}}.exe, coverage.out, coverage.html

# Run linter (requires golangci-lint)
lint:
    golangci-lint run ./...

# Format code
fmt:
    go fmt ./...

# Run go vet
vet:
    go vet ./...

# Run fmt, vet, and test
check: fmt vet test

# Build and run with arguments (e.g., just run --help)
[unix]
run *args: build
    ./{{binary_name}} {{args}}

[windows]
run *args: build
    .\{{binary_name}}.exe {{args}}

# Build with version info (for testing ldflags)
build-version version="dev":
    go build -ldflags="-X 'main.version={{version}}'" -o {{binary_name}} ./cmd/bmaduum

# ============================================================================
# Release Workflow
# ============================================================================

# Run GoReleaser in snapshot mode (for local testing)
release-snapshot:
    goreleaser release --snapshot

# Release: bump version, create tag, and push (triggers GitHub Action)
# Usage: just release [patch|minor|major] [--dry-run]
release level="patch" dry-run="":
    #!/usr/bin/env bash
    set -eo pipefail

    # Validate level
    case "{{level}}" in
        patch|minor|major) ;;
        *) echo "Error: level must be patch, minor, or major"; exit 1 ;;
    esac

    # Get current version
    current=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
    current_tag=$current
    current=${current#v}
    IFS='.' read -r major minor patch <<< "$current"

    # Default missing parts to 0
    major=${major:-0}
    minor=${minor:-0}
    patch=${patch:-0}

    # Calculate next version
    case "{{level}}" in
        major) ((++major)); minor=0; patch=0 ;;
        minor) ((++minor)); patch=0 ;;
        patch) ((++patch)) ;;
    esac
    new="v$major.$minor.$patch"

    echo "Release: $current_tag -> $new"
    echo

    # Dry run mode
    dry_run="{{dry-run}}"
    if [ -n "$dry_run" ]; then
        echo "DRY RUN - No changes will be made"
        echo
        echo "Would create tag: $new"
        echo
        echo "Commits to include:"
        git log "$current_tag"..HEAD --pretty=format:"- %s (%h)" --reverse 2>/dev/null || echo "No commits found."
        echo
        echo "Run 'just release {{level}}' to perform the release."
        exit 0
    fi

    # Safety checks
    if [ -n "$(git status --porcelain)" ]; then
        echo "Error: Working directory is not clean"
        echo "Please commit or stash changes first."
        exit 1
    fi

    # Create and push tag
    echo "Creating tag: $new"
    git tag -a "$new" -m "Release $new"

    echo "Pushing tag to origin..."
    git push origin "$new"

    echo
    echo "Release $new tagged and pushed!"
    echo
    echo "GitHub Actions will now:"
    echo "  1. Run GoReleaser to build binaries"
    echo "  2. Generate changelog from commits"
    echo "  3. Create GitHub release"
    echo "  4. Publish to package managers"
