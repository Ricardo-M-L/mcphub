package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestAddServer(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "claude_desktop_config.json")

	// Write initial config
	initial := map[string]interface{}{
		"mcpServers": map[string]interface{}{
			"existing-server": map[string]interface{}{
				"command": "npx",
				"args":    []interface{}{"-y", "existing-pkg"},
			},
		},
		"otherSetting": "should be preserved",
	}
	data, _ := json.MarshalIndent(initial, "", "  ")
	os.WriteFile(configPath, data, 0o644)

	cc := ClientConfig{Name: "test", Path: configPath, Key: "mcpServers"}
	err := cc.AddServer("new-server", MCPServerEntry{
		Command: "npx",
		Args:    []string{"-y", "new-pkg"},
		Env:     map[string]string{"API_KEY": "test123"},
	})
	if err != nil {
		t.Fatalf("AddServer failed: %v", err)
	}

	// Read back and verify
	result, _ := os.ReadFile(configPath)
	var parsed map[string]interface{}
	json.Unmarshal(result, &parsed)

	// Check other settings preserved
	if parsed["otherSetting"] != "should be preserved" {
		t.Error("otherSetting was not preserved")
	}

	// Check existing server preserved
	servers := parsed["mcpServers"].(map[string]interface{})
	if _, ok := servers["existing-server"]; !ok {
		t.Error("existing-server was removed")
	}

	// Check new server added
	newServer, ok := servers["new-server"].(map[string]interface{})
	if !ok {
		t.Fatal("new-server not found")
	}
	if newServer["command"] != "npx" {
		t.Errorf("expected command=npx, got %v", newServer["command"])
	}

	// Check backup created
	if _, err := os.Stat(configPath + ".bak"); err != nil {
		t.Error("backup file not created")
	}
}

func TestRemoveServer(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	initial := map[string]interface{}{
		"mcpServers": map[string]interface{}{
			"to-remove": map[string]interface{}{"command": "test"},
			"to-keep":   map[string]interface{}{"command": "keep"},
		},
	}
	data, _ := json.MarshalIndent(initial, "", "  ")
	os.WriteFile(configPath, data, 0o644)

	cc := ClientConfig{Name: "test", Path: configPath, Key: "mcpServers"}
	err := cc.RemoveServer("to-remove")
	if err != nil {
		t.Fatalf("RemoveServer failed: %v", err)
	}

	result, _ := os.ReadFile(configPath)
	var parsed map[string]interface{}
	json.Unmarshal(result, &parsed)

	servers := parsed["mcpServers"].(map[string]interface{})
	if _, ok := servers["to-remove"]; ok {
		t.Error("to-remove was not removed")
	}
	if _, ok := servers["to-keep"]; !ok {
		t.Error("to-keep was removed")
	}
}

func TestAddServerCreatesFile(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "subdir", "config.json")

	cc := ClientConfig{Name: "test", Path: configPath, Key: "mcpServers"}
	err := cc.AddServer("new-server", MCPServerEntry{
		URL: "https://example.com/mcp",
	})
	if err != nil {
		t.Fatalf("AddServer failed: %v", err)
	}

	if _, err := os.Stat(configPath); err != nil {
		t.Fatal("config file not created")
	}
}
