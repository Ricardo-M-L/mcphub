# MCP Hub

The package manager for [Model Context Protocol](https://modelcontextprotocol.io) servers.

Search, install, and manage MCP servers with a single command. Auto-configures Claude Desktop, Cursor, and other MCP clients.

```bash
mcphub install @anthropic/server-filesystem
```

## Features

- **One-command install** - Install any MCP server and auto-configure your clients
- **Smart detection** - Automatically finds Claude Desktop, Cursor, and other MCP clients
- **Registry search** - Browse the entire MCP server ecosystem from your terminal
- **Zero config** - Works out of the box with sensible defaults
- **Safe config updates** - Backs up your config files before modifying them
- **Cross-platform** - macOS, Linux, and Windows support

## Quick Start

### Install mcphub

```bash
# macOS / Linux
curl -fsSL https://mcphub.dev/install.sh | sh

# Homebrew
brew install mcphub/tap/mcphub

# From source
go install github.com/Ricardo-M-L/mcphub/cmd/mcphub@latest
```

### Usage

```bash
# Search for MCP servers
mcphub search filesystem

# Install an MCP server
mcphub install @modelcontextprotocol/server-filesystem

# List installed servers
mcphub list

# Show server details
mcphub info @anthropic/server-filesystem

# Remove a server
mcphub remove @modelcontextprotocol/server-filesystem
```

## How It Works

1. **Search** the [MCP Registry](https://registry.modelcontextprotocol.io) for servers
2. **Install** via npm/npx (for Node.js servers) or configure remote URLs
3. **Auto-configure** detected MCP clients (Claude Desktop, Cursor)
4. **Track** installed servers in `~/.mcphub/mcphub-lock.json`

### Supported MCP Clients

| Client | Config Path | Status |
|--------|------------|--------|
| Claude Desktop | `~/Library/Application Support/Claude/claude_desktop_config.json` | Supported |
| Cursor | `~/.cursor/mcp.json` | Supported |
| Claude Code | `~/.claude/settings.json` | Planned |
| OpenCode | `opencode.json` | Planned |

## Configuration

mcphub stores its data in `~/.mcphub/`:

```
~/.mcphub/
├── mcphub-lock.json    # Installed packages
└── config.json         # Registry settings (optional)
```

### Custom Registry

```bash
# Point to a custom registry
mcphub search --registry https://my-registry.example.com filesystem
```

## Development

### Prerequisites

- Go 1.23+

### Build

```bash
git clone https://github.com/mcphub/mcphub.git
cd mcphub
make build
```

### Test

```bash
make test
```

### Project Structure

```
mcphub/
├── cmd/mcphub/          # CLI entrypoint
├── internal/
│   ├── cli/             # Command definitions
│   ├── registry/        # MCP Registry API client
│   ├── installer/       # Package installers (npm, remote)
│   ├── config/          # MCP client config management
│   ├── store/           # Local package tracking
│   ├── ui/              # Terminal output helpers
│   └── platform/        # OS-specific paths
├── server/              # Registry API server (Go)
└── web/                 # Web UI (Next.js)
```

## Roadmap

- [x] CLI with search, install, list, remove
- [x] Auto-configure Claude Desktop and Cursor
- [ ] Registry API server with caching and full-text search
- [ ] Web UI for browsing MCP servers
- [ ] Quality scoring (maintenance, security, compatibility)
- [ ] `mcphub publish` for server authors
- [ ] SDK auto-generation from MCP tool schemas
- [ ] Claude Code and OpenCode support

## Contributing

We welcome contributions! See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

Apache-2.0 - see [LICENSE](LICENSE) for details.
