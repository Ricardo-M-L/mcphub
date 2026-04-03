# Contributing to MCP Hub

Thank you for your interest in contributing to MCP Hub! This document provides guidelines to help you get started.

## Getting Started

### Prerequisites

- Go 1.23 or later
- Node.js 20+ (for web UI development)
- Git

### Setup

```bash
git clone https://github.com/mcphub/mcphub.git
cd mcphub
make build
make test
```

### Running Locally

```bash
# Build and run
go run ./cmd/mcphub search filesystem

# Or build first
make build
./bin/mcphub search filesystem
```

## Types of Contributions

### Bug Fixes

- Check the [issue tracker](https://github.com/mcphub/mcphub/issues) for known bugs
- Look for issues labeled `good first issue` or `help wanted`
- Open an issue before starting work on complex fixes

### New Features

- Open a feature request issue first to discuss the design
- Keep PRs focused on a single feature
- Include tests for new functionality

### Documentation

- Fix typos, improve examples, add missing docs
- Documentation PRs don't need an issue

## Pull Request Process

### 1. Branch Naming

Use conventional prefixes:

```
feat/add-docker-support
fix/config-backup-race-condition
docs/update-installation-guide
refactor/simplify-installer-interface
test/add-config-manager-tests
```

### 2. Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
feat: add Docker-based MCP server installer
fix: prevent config corruption when file is locked
docs: add troubleshooting section to README
test: add integration tests for npm installer
```

### 3. Code Style

- Run `make fmt` before committing
- Run `make lint` to check for issues
- Follow Go conventions (effective Go, Go Code Review Comments)
- Keep functions focused and short
- Use meaningful variable names
- Add comments only where the logic isn't self-evident

### 4. Testing

- Add unit tests for new code in `*_test.go` files
- Use table-driven tests where appropriate
- Mock external dependencies (HTTP calls, file system)
- Run `make test` before submitting

### 5. PR Description

Use this template:

```markdown
## What

Brief description of the change.

## Why

The problem this solves or feature this adds.

## How

Key implementation decisions.

## Testing

How you verified this works.
```

## Architecture Overview

```
cmd/mcphub/main.go          # Entrypoint
internal/cli/                # Cobra commands (search, install, list, remove)
internal/registry/client.go  # HTTP client for MCP Registry API
internal/installer/          # Package installation logic
internal/config/manager.go   # MCP client config read/write
internal/store/store.go      # Local lockfile management
internal/ui/                 # Terminal output formatting
internal/platform/paths.go   # OS-specific file paths
```

### Key Design Principles

1. **Config safety** - Always backup before modifying client configs
2. **Zero infrastructure** - CLI works by querying the upstream registry directly
3. **Single binary** - No runtime dependencies for the CLI
4. **Preserve user config** - Use `map[string]interface{}` when reading configs to keep unknown fields

## Code of Conduct

This project follows the [Contributor Covenant Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

## Questions?

Open a [Discussion](https://github.com/mcphub/mcphub/discussions) or file an issue.
