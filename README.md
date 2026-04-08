# MCP Hub

The package manager for [Model Context Protocol](https://modelcontextprotocol.io) servers.

Search, install, and manage MCP servers with a single command. Auto-configures Claude Desktop, Cursor, Claude Code, and other MCP clients.

```bash
mcphub search filesystem
mcphub install io.github.user/server-filesystem
```

## Features

- **One-command install** - Install any MCP server and auto-configure your clients
- **Smart detection** - Automatically finds Claude Desktop, Cursor, and other MCP clients
- **Registry search** - Browse the entire MCP server ecosystem from your terminal
- **MCP server mode** - Use mcphub as an MCP service inside Claude Code or Cursor
- **Zero config** - Works out of the box with sensible defaults
- **Safe config updates** - Backs up your config files before modifying them
- **Cross-platform** - macOS, Linux, and Windows support
- **Single binary** - No runtime dependencies

---

## Installation

Choose any method:

### curl (macOS / Linux)

```bash
curl -fsSL https://raw.githubusercontent.com/Ricardo-M-L/mcphub/master/install.sh | sh
```

### Homebrew

```bash
brew tap Ricardo-M-L/mcphub
brew install mcphub
```

### npm

```bash
npm install -g @ricardo.m.lu/mcphub
```

### Go

```bash
go install github.com/Ricardo-M-L/mcphub/cmd/mcphub@latest
```

### From source

```bash
git clone https://github.com/Ricardo-M-L/mcphub.git
cd mcphub
make build
# Binary at ./bin/mcphub
```

---

## Usage

### Search for MCP servers

```bash
mcphub search filesystem
mcphub search database
mcphub search github
mcphub search slack

# JSON output
mcphub search filesystem --json

# Limit results
mcphub search database --limit 5
```

### View server details

```bash
mcphub info io.github.user/server-filesystem
mcphub info io.github.user/server-filesystem --json
```

### Install an MCP server

```bash
mcphub install io.github.user/server-filesystem
```

What happens when you run install:

1. Queries the [MCP Registry](https://registry.modelcontextprotocol.io) for the server
2. Determines the install method (npm/npx or remote URL)
3. Prompts for required environment variables (API keys, etc.)
4. Auto-detects installed MCP clients (Claude Desktop, Cursor)
5. Writes the server config into each client's config file (with backup)
6. Records the install in `~/.mcphub/mcphub-lock.json`

Target a specific client:

```bash
mcphub install io.github.user/server-filesystem --client claude-desktop
mcphub install io.github.user/server-filesystem --client cursor
```

### List installed servers

```bash
mcphub list
mcphub list --json
```

### Remove a server

```bash
mcphub remove io.github.user/server-filesystem
```

This removes the server from all configured MCP clients and the lockfile.

### Publish your own MCP server

```bash
# Create a manifest
mcphub init

# Validate
mcphub publish --dry-run

# Publish
mcphub publish
```

---

## Use as MCP Service (Claude Code / Cursor)

mcphub can run as an MCP server itself, so you can search and install MCP servers directly from your AI assistant.

### Setup

#### 1. Build the MCP server binary

```bash
cd mcphub
go build -o bin/mcphub-mcp ./mcp
```

Or if you installed via Go:

```bash
go install github.com/Ricardo-M-L/mcphub/mcp@latest
```

#### 2. Configure your MCP client

**Claude Code** (simplest - one command):

```bash
claude mcp add mcphub mcphub-mcp
```

Or manually add to `~/.claude/settings.json`:

```json
{
  "mcpServers": {
    "mcphub": {
      "command": "/path/to/mcphub-mcp"
    }
  }
}
```

**Claude Desktop** - Add to `~/Library/Application Support/Claude/claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "mcphub": {
      "command": "/path/to/mcphub-mcp"
    }
  }
}
```

**Cursor** - Add to `~/.cursor/mcp.json`:

```json
{
  "mcpServers": {
    "mcphub": {
      "command": "/path/to/mcphub-mcp"
    }
  }
}
```

#### 3. Use in conversation

After restarting your client, you can say things like:

- "Search for database MCP servers"
- "Install the filesystem MCP server"
- "What MCP servers do I have installed?"
- "Remove the github MCP server"

### Available MCP Tools

| Tool | Description |
|------|-------------|
| `search_servers` | Search the MCP registry for servers |
| `install_server` | Install a server and auto-configure clients |
| `list_installed` | List all installed MCP servers |
| `remove_server` | Remove a server from all clients |
| `server_info` | Get detailed info about a server |

---

## How It Works

```
                                    MCP Registry
                                  (registry.modelcontextprotocol.io)
                                         |
                                    HTTP API
                                         |
    Terminal ──── mcphub CLI ────────────┤
                      |                  |
                      |            Search / Info
                      |
                Install / Remove
                      |
            ┌─────────┼─────────┐
            |         |         |
      Claude       Cursor    Claude
      Desktop                  Code
            |         |         |
        config     config    config
         .json      .json    .json
```

- **mcphub does NOT download packages** - For npm-based servers, it writes `npx -y <package>` into your client config. npx handles download and caching at runtime.
- **mcphub does NOT run MCP servers** - It only configures your MCP clients to run them.
- **Config safety** - Always creates `.bak` backups before modifying config files.
- **Lockfile** - Tracks installed servers at `~/.mcphub/mcphub-lock.json`.

---

## Supported MCP Clients

| Client | Config Path | Auto-detect |
|--------|------------|-------------|
| Claude Desktop | `~/Library/Application Support/Claude/claude_desktop_config.json` | Yes |
| Cursor | `~/.cursor/mcp.json` | Yes |
| Claude Code | `~/.claude/settings.json` | Planned |
| OpenCode | `opencode.json` | Planned |

---

## Project Structure

```
mcphub/
├── cmd/mcphub/          # CLI entrypoint
├── mcp/                 # MCP server mode (stdio transport)
├── internal/
│   ├── cli/             # Command definitions (search, install, list, remove, info, publish)
│   ├── registry/        # MCP Registry API client
│   ├── installer/       # Package installers (npm, remote)
│   ├── config/          # MCP client config read/write (with backup)
│   ├── store/           # Local lockfile management
│   ├── ui/              # Terminal output formatting
│   └── platform/        # OS-specific file paths
├── server/              # Registry API server (Go + SQLite + FTS5)
│   ├── handler/         # REST API handlers
│   ├── crawler/         # Upstream registry sync
│   └── db/              # SQLite with full-text search
├── web/                 # Web UI (Next.js + TypeScript)
├── npm/                 # npm package wrapper
├── install.sh           # curl installer script
├── Dockerfile           # Container build
├── Makefile             # Build targets
└── .github/workflows/   # CI/CD + Release automation
```

---

## Development

### Prerequisites

- Go 1.23+
- Node.js 20+ (for web UI)

### Build

```bash
# CLI
make build

# MCP server
go build -o bin/mcphub-mcp ./mcp

# Registry API server
go build -o bin/mcphub-server ./server

# Web UI
cd web && npm install && npx next build
```

### Test

```bash
make test
```

### Run locally

```bash
# CLI
./bin/mcphub search filesystem

# MCP server (stdio)
echo '{"jsonrpc":"2.0","id":1,"method":"tools/list"}' | ./bin/mcphub-mcp

# Registry API server
./bin/mcphub-server --port 8080
```

---

## Roadmap

- [x] CLI: search, install, list, remove, info, publish, init
- [x] Auto-configure Claude Desktop and Cursor
- [x] MCP server mode for Claude Code / Cursor integration
- [x] Registry API server with SQLite + FTS5 full-text search
- [x] Web UI for browsing MCP servers
- [x] Quality scoring API (completeness, installability, documentation, security)
- [x] CI/CD with GitHub Actions + cross-platform release builds
- [x] 5 distribution methods: curl, Homebrew, npm, Go, source
- [x] 14 unit tests passing
- [ ] Claude Code auto-detection support
- [ ] SDK auto-generation from MCP tool schemas
- [ ] Server health monitoring
- [ ] Community ratings and reviews

---

## Contributing

We welcome contributions! See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

Apache-2.0 - see [LICENSE](LICENSE) for details.
