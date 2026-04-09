# Contributing to MCP Hub

Thank you for your interest in contributing to MCP Hub! Please read this guide carefully before submitting.

## Rules

1. **Issue first** — Open an issue before working on any non-trivial change. PRs without a linked issue may be closed.
2. **One PR, one thing** — Keep PRs focused. Don't mix bug fixes with features.
3. **CI must pass** — All checks (build, test, lint) must be green before review.
4. **PR title must follow conventional commits** — `feat:`, `fix:`, `docs:`, etc.
5. **No AI-generated walls of text** — Write concise, human descriptions.
6. **Maintainer approval required** — All PRs require at least 1 approval from a maintainer before merging.
7. **No direct pushes to master** — All changes go through PRs.

## Getting Started

### Prerequisites

- Go 1.23 or later
- Node.js 20+ (for web UI development)
- Git

### Setup

```bash
git clone https://github.com/Ricardo-M-L/mcphub.git
cd mcphub
make build
make test
```

### Running Locally

```bash
go run ./cmd/mcphub search filesystem
```

## Types of Contributions

### Bug Fixes

- Check the [issue tracker](https://github.com/Ricardo-M-L/mcphub/issues) for known bugs
- Look for issues labeled `good first issue` or `help wanted`
- **Open an issue first** describing the bug and your fix approach

### New Features

- **Open a feature request issue first** to discuss the design
- Wait for maintainer approval before starting work
- Include tests for new functionality

### Documentation

- Fix typos, improve examples, add missing docs
- Small doc fixes don't need an issue

## Pull Request Process

### 1. Fork & Branch

```bash
# Fork the repo on GitHub, then:
git clone https://github.com/YOUR-USERNAME/mcphub.git
cd mcphub
git remote add upstream https://github.com/Ricardo-M-L/mcphub.git
git checkout -b feat/your-feature
```

### 2. Branch Naming

```
feat/add-docker-support
fix/config-backup-race-condition
docs/update-installation-guide
refactor/simplify-installer-interface
test/add-config-manager-tests
```

### 3. PR Title (Conventional Commits)

```
feat: add Docker-based MCP server installer
fix: prevent config corruption when file is locked
docs: add troubleshooting section to README
test: add integration tests for npm installer
```

**Invalid titles will be automatically rejected by CI.**

### 4. PR Description

Use the PR template. Must include:

- **What** — Brief description of the change
- **Why** — The problem this solves
- **How** — Key implementation decisions
- **Testing** — How you verified it works
- **Closes #** — Link to the issue

### 5. Code Requirements

- Run `make fmt` before committing
- Run `make test` — all tests must pass
- Run `go vet ./...` — no warnings
- Add tests for new code
- Keep functions focused and short
- Comments only where logic isn't self-evident

### 6. Review Process

1. CI checks run automatically (build, test, lint, PR title)
2. Maintainer reviews the code
3. Address review feedback with new commits (don't force-push during review)
4. Maintainer approves and merges

**Timeline:** Expect 1-3 days for initial review. Complex PRs may take longer.

### 7. What Will Get Your PR Closed

- No linked issue (for non-trivial changes)
- AI-generated spam descriptions
- Unrelated changes mixed in
- Failing CI
- No response to review feedback for 7 days
- Changes that break existing functionality without tests

## Architecture

```
cmd/mcphub/main.go          # CLI entrypoint
mcp/server.go               # MCP server mode (stdio)
internal/cli/                # Cobra commands
internal/registry/client.go  # MCP Registry API client
internal/registry/github.go  # GitHub search fallback
internal/installer/          # Package installation logic
internal/config/manager.go   # MCP client config management
internal/store/store.go      # Local lockfile
internal/ui/                 # Terminal output
internal/platform/paths.go   # OS-specific paths
server/                      # Registry API server
web/                         # Next.js Web UI
```

### Design Principles

1. **Config safety** — Always backup before modifying client configs
2. **Zero infrastructure** — CLI queries upstream registry directly
3. **Single binary** — No runtime dependencies
4. **Preserve user config** — Use `map[string]interface{}` to keep unknown fields

## Code of Conduct

Be respectful. No harassment, spam, or bad-faith contributions. Violations result in a permanent ban.

## Questions?

Open a [Discussion](https://github.com/Ricardo-M-L/mcphub/discussions) or file an issue.
