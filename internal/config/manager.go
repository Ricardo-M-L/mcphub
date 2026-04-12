package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Ricardo-M-L/mcphub/internal/platform"
)

// ClientConfig provides read/write access to an MCP client's configuration.
type ClientConfig struct {
	Name string
	Path string
	Key  string // JSON key for the MCP servers map (e.g. "mcpServers")
}

// KnownClients returns configs for all supported MCP clients.
func KnownClients() []ClientConfig {
	return []ClientConfig{
		{Name: "claude-desktop", Path: platform.ClaudeDesktopConfigPath(), Key: "mcpServers"},
		{Name: "cursor", Path: platform.CursorConfigPath(), Key: "mcpServers"},
		{Name: "claude-code", Path: platform.ClaudeCodeConfigPath(), Key: "mcpServers"},
		{Name: "windsurf", Path: platform.WindsurfConfigPath(), Key: "mcpServers"},
	}
}

// DetectedClients returns only the clients that have config files on disk.
func DetectedClients() []ClientConfig {
	var result []ClientConfig
	for _, c := range KnownClients() {
		if _, err := os.Stat(c.Path); err == nil {
			result = append(result, c)
		}
	}
	return result
}

// AddServer adds an MCP server entry to a client's config file.
// It preserves all existing config by operating on a raw map.
func (c *ClientConfig) AddServer(shortName string, entry MCPServerEntry) error {
	raw, err := c.readRaw()
	if err != nil {
		return err
	}

	servers, _ := raw[c.Key].(map[string]interface{})
	if servers == nil {
		servers = make(map[string]interface{})
	}

	entryMap := make(map[string]interface{})
	if entry.Command != "" {
		entryMap["command"] = entry.Command
	}
	if len(entry.Args) > 0 {
		entryMap["args"] = entry.Args
	}
	if len(entry.Env) > 0 {
		entryMap["env"] = entry.Env
	}
	if entry.URL != "" {
		entryMap["url"] = entry.URL
	}

	servers[shortName] = entryMap
	raw[c.Key] = servers

	return c.writeRaw(raw)
}

// RemoveServer removes an MCP server entry from a client's config file.
func (c *ClientConfig) RemoveServer(shortName string) error {
	raw, err := c.readRaw()
	if err != nil {
		return err
	}

	servers, _ := raw[c.Key].(map[string]interface{})
	if servers == nil {
		return nil
	}

	delete(servers, shortName)
	raw[c.Key] = servers

	return c.writeRaw(raw)
}

func (c *ClientConfig) readRaw() (map[string]interface{}, error) {
	data, err := os.ReadFile(c.Path)
	if os.IsNotExist(err) {
		return make(map[string]interface{}), nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read %s config: %w", c.Name, err)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse %s config: %w", c.Name, err)
	}
	return raw, nil
}

func (c *ClientConfig) writeRaw(raw map[string]interface{}) error {
	if err := os.MkdirAll(filepath.Dir(c.Path), 0o755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Backup existing file
	if _, err := os.Stat(c.Path); err == nil {
		backup := c.Path + ".bak"
		data, _ := os.ReadFile(c.Path)
		_ = os.WriteFile(backup, data, 0o644)
	}

	out, err := json.MarshalIndent(raw, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	out = append(out, '\n')
	return os.WriteFile(c.Path, out, 0o644)
}
