package health

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestCheckRemoteServer_OK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))
	defer srv.Close()

	report := CheckRemoteServer("test-server", srv.URL)
	if report.Status != StatusOK {
		t.Errorf("expected StatusOK, got %s", report.Status)
	}
	if report.Name != "test-server" {
		t.Errorf("expected name test-server, got %s", report.Name)
	}
}

func TestCheckRemoteServer_500(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	report := CheckRemoteServer("test-server", srv.URL)
	if report.Status != StatusError {
		t.Errorf("expected StatusError, got %s", report.Status)
	}
}

func TestCheckRemoteServer_404(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	report := CheckRemoteServer("test-server", srv.URL)
	if report.Status != StatusWarning {
		t.Errorf("expected StatusWarning, got %s", report.Status)
	}
}

func TestCheckRemoteServer_EmptyURL(t *testing.T) {
	report := CheckRemoteServer("test-server", "")
	if report.Status != StatusError {
		t.Errorf("expected StatusError, got %s", report.Status)
	}
}

func TestCheckRemoteServer_Unreachable(t *testing.T) {
	report := CheckRemoteServer("test-server", "http://127.0.0.1:1")
	if report.Status != StatusError {
		t.Errorf("expected StatusError, got %s", report.Status)
	}
}

func TestCheckStdioServer_CommandNotFound(t *testing.T) {
	report := CheckStdioServer("test", "nonexistent-binary-xyz-123", nil, nil)
	if report.Status != StatusError {
		t.Errorf("expected StatusError, got %s", report.Status)
	}
	if len(report.Issues) == 0 {
		t.Fatal("expected at least 1 issue")
	}
	if report.Issues[0].Message == "" {
		t.Error("expected non-empty issue message")
	}
}

func TestCheckStdioServer_EchoServer(t *testing.T) {
	// Use a simple echo-like command that outputs valid JSON
	// Create a temp script that responds with valid MCP initialize response
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "mock-mcp.sh")
	response := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"result": map[string]interface{}{
			"protocolVersion": "2025-03-26",
			"serverInfo": map[string]interface{}{
				"name":    "test-server",
				"version": "1.0.0",
			},
		},
	}
	respBytes, _ := json.Marshal(response)
	script := fmt.Sprintf("#!/bin/sh\nread line\necho '%s'\n", string(respBytes))
	os.WriteFile(scriptPath, []byte(script), 0o755)

	report := CheckStdioServer("test-server", "sh", []string{scriptPath}, nil)
	if report.Status != StatusOK {
		t.Errorf("expected StatusOK, got %s. Issues: %+v", report.Status, report.Issues)
	}
}

func TestCheckClientConfig_ValidConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	config := map[string]interface{}{
		"mcpServers": map[string]interface{}{
			"test": map[string]interface{}{
				"command": "npx",
				"args":    []string{"-y", "test-pkg"},
			},
		},
	}
	data, _ := json.MarshalIndent(config, "", "  ")
	os.WriteFile(configPath, data, 0o644)

	issues := CheckClientConfig("test-client", configPath, "mcpServers")
	for _, issue := range issues {
		if issue.Severity == StatusError {
			t.Errorf("unexpected error: %s", issue.Message)
		}
	}
}

func TestCheckClientConfig_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")
	os.WriteFile(configPath, []byte("{invalid json"), 0o644)

	issues := CheckClientConfig("test-client", configPath, "mcpServers")
	if len(issues) == 0 {
		t.Fatal("expected issues for invalid JSON")
	}
	if issues[0].Severity != StatusError {
		t.Errorf("expected StatusError, got %s", issues[0].Severity)
	}
}

func TestCheckClientConfig_MissingCommand(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	config := map[string]interface{}{
		"mcpServers": map[string]interface{}{
			"broken": map[string]interface{}{
				"args": []string{"test"},
			},
		},
	}
	data, _ := json.MarshalIndent(config, "", "  ")
	os.WriteFile(configPath, data, 0o644)

	issues := CheckClientConfig("test-client", configPath, "mcpServers")
	hasError := false
	for _, issue := range issues {
		if issue.Severity == StatusError {
			hasError = true
		}
	}
	if !hasError {
		t.Error("expected error for missing command/url")
	}
}

func TestCheckClientConfig_NonExistent(t *testing.T) {
	issues := CheckClientConfig("test", "/nonexistent/path.json", "mcpServers")
	if len(issues) != 0 {
		t.Errorf("expected no issues for non-existent file, got %d", len(issues))
	}
}

func TestSuggestInstall(t *testing.T) {
	tests := []struct {
		command  string
		contains string
	}{
		{"npx", "Node.js"},
		{"node", "Node.js"},
		{"uvx", "uv"},
		{"python", "Python"},
		{"docker", "Docker"},
		{"/usr/local/bin/something", "Binary not found"},
		{"unknown", "PATH"},
	}

	for _, tt := range tests {
		result := suggestInstall(tt.command)
		if result == "" {
			t.Errorf("suggestInstall(%q) returned empty", tt.command)
		}
	}
}

func TestStatusString(t *testing.T) {
	if StatusOK.String() != "OK" {
		t.Errorf("expected OK, got %s", StatusOK.String())
	}
	if StatusWarning.String() != "WARNING" {
		t.Errorf("expected WARNING, got %s", StatusWarning.String())
	}
	if StatusError.String() != "ERROR" {
		t.Errorf("expected ERROR, got %s", StatusError.String())
	}
}
