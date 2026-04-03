package config

import "github.com/mcphub/mcphub/internal/registry"

// MCPServerEntry is the common MCP server config format used by Claude Desktop and Cursor.
type MCPServerEntry struct {
	Command string            `json:"command,omitempty"`
	Args    []string          `json:"args,omitempty"`
	Env     map[string]string `json:"env,omitempty"`
	URL     string            `json:"url,omitempty"`
}

// ToTransport converts the entry to a registry Transport.
func (e MCPServerEntry) ToTransport() registry.Transport {
	if e.URL != "" {
		return registry.Transport{Type: "streamable-http", URL: e.URL}
	}
	return registry.Transport{Type: "stdio"}
}
