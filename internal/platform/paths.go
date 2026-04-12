package platform

import (
	"os"
	"path/filepath"
	"runtime"
)

// MCPHubDir returns the mcphub data directory.
func MCPHubDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".mcphub")
}

// LockfilePath returns the path to mcphub-lock.json.
func LockfilePath() string {
	return filepath.Join(MCPHubDir(), "mcphub-lock.json")
}

// ClaudeDesktopConfigPath returns the Claude Desktop config file path.
func ClaudeDesktopConfigPath() string {
	home, _ := os.UserHomeDir()
	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(home, "Library", "Application Support", "Claude", "claude_desktop_config.json")
	case "windows":
		return filepath.Join(os.Getenv("APPDATA"), "Claude", "claude_desktop_config.json")
	default:
		return filepath.Join(home, ".config", "Claude", "claude_desktop_config.json")
	}
}

// CursorConfigPath returns the Cursor MCP config file path.
func CursorConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".cursor", "mcp.json")
}

// ClaudeCodeConfigPath returns the Claude Code settings path.
func ClaudeCodeConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".claude.json")
}

// WindsurfConfigPath returns the Windsurf MCP config path.
func WindsurfConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".windsurf", "mcp.json")
}

// DetectInstalledClients returns a list of MCP client names that have config files present.
func DetectInstalledClients() []string {
	var clients []string
	paths := map[string]string{
		"claude-desktop": ClaudeDesktopConfigPath(),
		"cursor":         CursorConfigPath(),
		"claude-code":    ClaudeCodeConfigPath(),
	}
	for name, p := range paths {
		if _, err := os.Stat(p); err == nil {
			clients = append(clients, name)
		}
	}
	return clients
}
