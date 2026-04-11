package health

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

// Status represents the health status of an MCP server.
type Status int

const (
	StatusOK      Status = iota // Server is healthy
	StatusWarning               // Server has minor issues
	StatusError                 // Server is broken
)

func (s Status) String() string {
	switch s {
	case StatusOK:
		return "OK"
	case StatusWarning:
		return "WARNING"
	case StatusError:
		return "ERROR"
	}
	return "UNKNOWN"
}

// Issue represents a diagnosed problem with an MCP server.
type Issue struct {
	Severity    Status `json:"severity"`
	Message     string `json:"message"`
	Suggestion  string `json:"suggestion"`
}

// Report is the health check result for a single MCP server.
type Report struct {
	Name       string  `json:"name"`
	Status     Status  `json:"status"`
	Issues     []Issue `json:"issues"`
	Duration   time.Duration `json:"duration"`
	Transport  string  `json:"transport"`
}

// CheckStdioServer checks a stdio-based MCP server by spawning it and sending an initialize request.
func CheckStdioServer(name, command string, args []string, env map[string]string) *Report {
	start := time.Now()
	report := &Report{
		Name:      name,
		Transport: "stdio",
	}

	// 1. Check if command exists
	cmdPath, err := exec.LookPath(command)
	if err != nil {
		report.Status = StatusError
		report.Issues = append(report.Issues, Issue{
			Severity:   StatusError,
			Message:    fmt.Sprintf("Command not found: %s", command),
			Suggestion: suggestInstall(command),
		})
		report.Duration = time.Since(start)
		return report
	}
	_ = cmdPath

	// 2. Check environment variables
	for k, v := range env {
		if v == "" {
			report.Issues = append(report.Issues, Issue{
				Severity:   StatusWarning,
				Message:    fmt.Sprintf("Environment variable %s is empty", k),
				Suggestion: fmt.Sprintf("Set %s in your environment or re-run: mcphub install <server>", k),
			})
		}
	}

	// 3. Try to start the server and send initialize
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, command, args...)

	// Set environment
	cmd.Env = os.Environ()
	for k, v := range env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		report.Status = StatusError
		report.Issues = append(report.Issues, Issue{
			Severity:   StatusError,
			Message:    fmt.Sprintf("Failed to create stdin pipe: %s", err),
			Suggestion: "Check file permissions and system resources",
		})
		report.Duration = time.Since(start)
		return report
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		report.Status = StatusError
		report.Issues = append(report.Issues, Issue{
			Severity:   StatusError,
			Message:    fmt.Sprintf("Failed to create stdout pipe: %s", err),
			Suggestion: "Check file permissions and system resources",
		})
		report.Duration = time.Since(start)
		return report
	}

	cmd.Stderr = nil // Ignore stderr

	if err := cmd.Start(); err != nil {
		report.Status = StatusError
		report.Issues = append(report.Issues, Issue{
			Severity:   StatusError,
			Message:    fmt.Sprintf("Failed to start server: %s", err),
			Suggestion: suggestInstall(command),
		})
		report.Duration = time.Since(start)
		return report
	}

	// Send initialize request
	initReq := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "initialize",
		"params": map[string]interface{}{
			"protocolVersion": "2025-03-26",
			"capabilities":   map[string]interface{}{},
			"clientInfo": map[string]interface{}{
				"name":    "mcphub-doctor",
				"version": "0.2.0",
			},
		},
	}

	reqBytes, _ := json.Marshal(initReq)
	reqBytes = append(reqBytes, '\n')

	_, err = stdin.Write(reqBytes)
	if err != nil {
		report.Status = StatusError
		report.Issues = append(report.Issues, Issue{
			Severity:   StatusError,
			Message:    fmt.Sprintf("Failed to send initialize request: %s", err),
			Suggestion: "Server may have crashed on startup. Check server logs.",
		})
		report.Duration = time.Since(start)
		cmd.Process.Kill()
		return report
	}

	// Read response with timeout
	scanner := bufio.NewScanner(stdout)
	responseCh := make(chan string, 1)
	go func() {
		if scanner.Scan() {
			responseCh <- scanner.Text()
		} else {
			responseCh <- ""
		}
	}()

	select {
	case resp := <-responseCh:
		if resp == "" {
			report.Status = StatusError
			report.Issues = append(report.Issues, Issue{
				Severity:   StatusError,
				Message:    "Server returned empty response to initialize",
				Suggestion: "Server may not implement MCP protocol correctly",
			})
		} else {
			// Parse response
			var result map[string]interface{}
			if err := json.Unmarshal([]byte(resp), &result); err != nil {
				report.Status = StatusError
				report.Issues = append(report.Issues, Issue{
					Severity:   StatusError,
					Message:    fmt.Sprintf("Invalid JSON response: %s", err),
					Suggestion: "Server output is not valid JSON-RPC. It may be printing debug info to stdout.",
				})
			} else if errObj, ok := result["error"]; ok {
				report.Status = StatusError
				report.Issues = append(report.Issues, Issue{
					Severity:   StatusError,
					Message:    fmt.Sprintf("Server returned error: %v", errObj),
					Suggestion: "Check server configuration and required environment variables",
				})
			} else if _, ok := result["result"]; ok {
				// Success! Check server info
				if resMap, ok := result["result"].(map[string]interface{}); ok {
					if serverInfo, ok := resMap["serverInfo"].(map[string]interface{}); ok {
						report.Issues = append(report.Issues, Issue{
							Severity:   StatusOK,
							Message:    fmt.Sprintf("Server: %s v%s", serverInfo["name"], serverInfo["version"]),
							Suggestion: "",
						})
					}
				}
			}
		}

	case <-ctx.Done():
		report.Status = StatusError
		report.Issues = append(report.Issues, Issue{
			Severity:   StatusError,
			Message:    "Server timed out (10s) during initialize handshake",
			Suggestion: "Server may be hanging on startup. Check if it needs network access or is waiting for input.",
		})
	}

	// Cleanup
	stdin.Close()
	cmd.Process.Kill()
	cmd.Wait()

	// If no errors were found, mark as OK
	if report.Status == 0 {
		hasError := false
		for _, issue := range report.Issues {
			if issue.Severity == StatusError {
				hasError = true
				break
			}
		}
		if hasError {
			report.Status = StatusError
		} else {
			report.Status = StatusOK
		}
	}

	report.Duration = time.Since(start)
	return report
}

// CheckRemoteServer checks a remote (streamable-http/sse) MCP server by making an HTTP request.
func CheckRemoteServer(name, url string) *Report {
	start := time.Now()
	report := &Report{
		Name:      name,
		Transport: "remote",
	}

	if url == "" {
		report.Status = StatusError
		report.Issues = append(report.Issues, Issue{
			Severity:   StatusError,
			Message:    "No URL configured",
			Suggestion: "Re-install the server with: mcphub install <server-name>",
		})
		report.Duration = time.Since(start)
		return report
	}

	// HTTP health check
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		report.Status = StatusError
		report.Issues = append(report.Issues, Issue{
			Severity:   StatusError,
			Message:    fmt.Sprintf("Cannot reach server: %s", err),
			Suggestion: "Check if the server is running and the URL is correct",
		})
		report.Duration = time.Since(start)
		return report
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 500 {
		report.Status = StatusError
		report.Issues = append(report.Issues, Issue{
			Severity:   StatusError,
			Message:    fmt.Sprintf("Server returned HTTP %d", resp.StatusCode),
			Suggestion: "Server is experiencing internal errors",
		})
	} else if resp.StatusCode >= 400 {
		report.Status = StatusWarning
		report.Issues = append(report.Issues, Issue{
			Severity:   StatusWarning,
			Message:    fmt.Sprintf("Server returned HTTP %d", resp.StatusCode),
			Suggestion: "Server may require authentication or the endpoint has changed",
		})
	} else {
		report.Status = StatusOK
		report.Issues = append(report.Issues, Issue{
			Severity:   StatusOK,
			Message:    fmt.Sprintf("Server reachable (HTTP %d)", resp.StatusCode),
			Suggestion: "",
		})
	}

	report.Duration = time.Since(start)
	return report
}

// CheckClientConfig validates an MCP client's configuration file.
func CheckClientConfig(name, path, key string) []Issue {
	var issues []Issue

	// Check if file exists
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return issues // File doesn't exist, client not installed
	}
	if err != nil {
		issues = append(issues, Issue{
			Severity:   StatusError,
			Message:    fmt.Sprintf("%s config unreadable: %s", name, err),
			Suggestion: "Check file permissions",
		})
		return issues
	}

	// Check if valid JSON
	var config map[string]interface{}
	if err := json.Unmarshal(data, &config); err != nil {
		issues = append(issues, Issue{
			Severity:   StatusError,
			Message:    fmt.Sprintf("%s config has invalid JSON: %s", name, err),
			Suggestion: fmt.Sprintf("Fix the JSON syntax in %s", path),
		})
		return issues
	}

	// Check mcpServers section
	servers, ok := config[key]
	if !ok {
		issues = append(issues, Issue{
			Severity:   StatusWarning,
			Message:    fmt.Sprintf("%s config has no '%s' section", name, key),
			Suggestion: "No MCP servers configured in this client",
		})
		return issues
	}

	serverMap, ok := servers.(map[string]interface{})
	if !ok {
		issues = append(issues, Issue{
			Severity:   StatusError,
			Message:    fmt.Sprintf("%s '%s' is not a valid object", name, key),
			Suggestion: fmt.Sprintf("Fix the '%s' section in %s", key, path),
		})
		return issues
	}

	// Validate each server entry
	for serverName, entry := range serverMap {
		entryMap, ok := entry.(map[string]interface{})
		if !ok {
			issues = append(issues, Issue{
				Severity:   StatusError,
				Message:    fmt.Sprintf("%s: server '%s' has invalid config", name, serverName),
				Suggestion: "Each server entry must be a JSON object",
			})
			continue
		}

		// Check for command or url
		_, hasCommand := entryMap["command"]
		_, hasURL := entryMap["url"]
		if !hasCommand && !hasURL {
			issues = append(issues, Issue{
				Severity:   StatusError,
				Message:    fmt.Sprintf("%s: server '%s' has no 'command' or 'url'", name, serverName),
				Suggestion: fmt.Sprintf("Add a 'command' or 'url' field to '%s' in %s", serverName, path),
			})
		}
	}

	return issues
}

func suggestInstall(command string) string {
	switch {
	case command == "npx" || command == "node":
		return "Install Node.js: https://nodejs.org or brew install node"
	case command == "uvx" || command == "uv":
		return "Install uv: curl -LsSf https://astral.sh/uv/install.sh | sh"
	case command == "python" || command == "python3":
		return "Install Python: https://python.org or brew install python"
	case command == "docker":
		return "Install Docker: https://docker.com or brew install --cask docker"
	case strings.HasPrefix(command, "/"):
		return fmt.Sprintf("Binary not found at %s. Re-install the server.", command)
	default:
		return fmt.Sprintf("Install %s or check your PATH", command)
	}
}
