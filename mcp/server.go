package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/Ricardo-M-L/mcphub/internal/config"
	"github.com/Ricardo-M-L/mcphub/internal/health"
	"github.com/Ricardo-M-L/mcphub/internal/installer"
	"github.com/Ricardo-M-L/mcphub/internal/registry"
	"github.com/Ricardo-M-L/mcphub/internal/store"
)

// MCP JSON-RPC types
type jsonrpcRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type jsonrpcResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Error   *rpcError   `json:"error,omitempty"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Tool definitions
var tools = []map[string]interface{}{
	{
		"name":        "search_servers",
		"description": "Search the MCP registry for servers matching a query. Returns a list of available MCP servers with name, description, version, and transport type.",
		"inputSchema": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"query": map[string]interface{}{
					"type":        "string",
					"description": "Search query (e.g. 'filesystem', 'database', 'github')",
				},
				"limit": map[string]interface{}{
					"type":        "number",
					"description": "Maximum number of results (default: 10)",
				},
			},
			"required": []string{"query"},
		},
	},
	{
		"name":        "install_server",
		"description": "Install an MCP server by its registry name. Auto-configures Claude Desktop and Cursor. For servers requiring environment variables, provide them in the env_vars parameter.",
		"inputSchema": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"name": map[string]interface{}{
					"type":        "string",
					"description": "Full registry name of the MCP server (e.g. 'io.github.user/server-name')",
				},
				"env_vars": map[string]interface{}{
					"type":        "object",
					"description": "Environment variables required by the server (e.g. {\"API_KEY\": \"xxx\"})",
				},
			},
			"required": []string{"name"},
		},
	},
	{
		"name":        "list_installed",
		"description": "List all currently installed MCP servers with their version, transport type, and which clients they are configured in.",
		"inputSchema": map[string]interface{}{
			"type":       "object",
			"properties": map[string]interface{}{},
		},
	},
	{
		"name":        "remove_server",
		"description": "Remove an installed MCP server and its configuration from all MCP clients.",
		"inputSchema": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"name": map[string]interface{}{
					"type":        "string",
					"description": "Name of the installed MCP server to remove",
				},
			},
			"required": []string{"name"},
		},
	},
	{
		"name":        "server_info",
		"description": "Get detailed information about an MCP server from the registry, including packages, remote endpoints, and required environment variables.",
		"inputSchema": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"name": map[string]interface{}{
					"type":        "string",
					"description": "Registry name of the MCP server",
				},
			},
			"required": []string{"name"},
		},
	},
	{
		"name":        "doctor",
		"description": "Diagnose MCP server health and configuration issues. Checks all installed servers for connectivity, validates client configs, and reports problems with fix suggestions.",
		"inputSchema": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"server": map[string]interface{}{
					"type":        "string",
					"description": "Optional: check a specific server only. If empty, checks all servers.",
				},
			},
		},
	},
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	client := registry.NewClient()

	for {
		line, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		var req jsonrpcRequest
		if err := json.Unmarshal(line, &req); err != nil {
			continue
		}

		var resp jsonrpcResponse
		resp.JSONRPC = "2.0"
		resp.ID = req.ID

		switch req.Method {
		case "initialize":
			resp.Result = map[string]interface{}{
				"protocolVersion": "2025-03-26",
				"capabilities": map[string]interface{}{
					"tools": map[string]interface{}{},
				},
				"serverInfo": map[string]interface{}{
					"name":    "mcphub",
					"version": "0.1.0",
				},
			}

		case "notifications/initialized":
			continue // No response needed

		case "tools/list":
			resp.Result = map[string]interface{}{
				"tools": tools,
			}

		case "tools/call":
			resp.Result = handleToolCall(req.Params, client)

		default:
			resp.Error = &rpcError{Code: -32601, Message: "method not found: " + req.Method}
		}

		out, _ := json.Marshal(resp)
		fmt.Fprintf(os.Stdout, "%s\n", out)
	}
}

func handleToolCall(params json.RawMessage, client *registry.Client) interface{} {
	var call struct {
		Name      string          `json:"name"`
		Arguments json.RawMessage `json:"arguments"`
	}
	if err := json.Unmarshal(params, &call); err != nil {
		return toolError("invalid params: " + err.Error())
	}

	switch call.Name {
	case "search_servers":
		return handleSearch(call.Arguments, client)
	case "install_server":
		return handleInstall(call.Arguments, client)
	case "list_installed":
		return handleList()
	case "remove_server":
		return handleRemove(call.Arguments)
	case "server_info":
		return handleInfo(call.Arguments, client)
	case "doctor":
		return handleDoctor(call.Arguments)
	default:
		return toolError("unknown tool: " + call.Name)
	}
}

func handleSearch(args json.RawMessage, client *registry.Client) interface{} {
	var params struct {
		Query string `json:"query"`
		Limit int    `json:"limit"`
	}
	json.Unmarshal(args, &params)
	if params.Limit <= 0 {
		params.Limit = 10
	}

	entries, err := client.SearchAll(params.Query, params.Limit)
	if err != nil {
		return toolError("search failed: " + err.Error())
	}

	var results []map[string]string
	for _, e := range entries {
		s := e.Server
		transport := "stdio"
		if len(s.Remotes) > 0 {
			transport = s.Remotes[0].Type
		} else if len(s.Packages) > 0 {
			transport = s.Packages[0].Transport.Type
		}
		results = append(results, map[string]string{
			"name":        s.Name,
			"shortName":   s.ShortName(),
			"description": s.Description,
			"version":     s.Version,
			"transport":   transport,
		})
	}

	text, _ := json.MarshalIndent(results, "", "  ")
	return toolResult(fmt.Sprintf("Found %d servers:\n%s", len(results), string(text)))
}

func handleInstall(args json.RawMessage, client *registry.Client) interface{} {
	var params struct {
		Name    string            `json:"name"`
		EnvVars map[string]string `json:"env_vars"`
	}
	json.Unmarshal(args, &params)

	entry, err := client.GetServer(params.Name)
	if err != nil {
		entries, searchErr := client.Search(params.Name, 5)
		if searchErr != nil || len(entries) == 0 {
			return toolError("server not found: " + params.Name)
		}
		entry = &entries[0]
	}

	server := &entry.Server
	inst, pkg, remote, err := installer.Select(server)
	if err != nil {
		return toolError(err.Error())
	}

	envVars := params.EnvVars
	if envVars == nil {
		envVars = make(map[string]string)
	}

	var result *installer.Result
	if remote != nil {
		result = installer.InstallRemote(remote, envVars)
	} else {
		result, err = inst.Install(*pkg, envVars)
		if err != nil {
			return toolError("install failed: " + err.Error())
		}
	}

	shortName := server.ShortName()
	clients := config.DetectedClients()
	configuredIn := []string{}
	for _, c := range clients {
		if err := c.AddServer(shortName, result.Entry); err != nil {
			continue
		}
		configuredIn = append(configuredIn, c.Name)
	}

	lf, _ := store.Load()
	regType, identifier := "", ""
	if pkg != nil {
		regType = pkg.RegistryType
		identifier = pkg.Identifier
	} else if remote != nil {
		regType = "remote"
		identifier = remote.URL
	}

	lf.Add(store.InstalledPackage{
		Name:         server.Name,
		Version:      server.Version,
		RegistryType: regType,
		Identifier:   identifier,
		Transport:    result.Entry.ToTransport(),
		EnvVars:      envVars,
		ConfiguredIn: configuredIn,
	})
	lf.Save()

	msg := fmt.Sprintf("Installed %s (v%s)\nConfigured in: %v", shortName, server.Version, configuredIn)
	return toolResult(msg)
}

func handleList() interface{} {
	lf, err := store.Load()
	if err != nil {
		return toolError("failed to load: " + err.Error())
	}

	if len(lf.Packages) == 0 {
		return toolResult("No MCP servers installed.")
	}

	var results []map[string]string
	for _, pkg := range lf.Packages {
		results = append(results, map[string]string{
			"name":         pkg.Name,
			"version":      pkg.Version,
			"transport":    pkg.Transport.Type,
			"configuredIn": fmt.Sprintf("%v", pkg.ConfiguredIn),
		})
	}

	text, _ := json.MarshalIndent(results, "", "  ")
	return toolResult(fmt.Sprintf("%d servers installed:\n%s", len(results), string(text)))
}

func handleRemove(args json.RawMessage) interface{} {
	var params struct {
		Name string `json:"name"`
	}
	json.Unmarshal(args, &params)

	lf, _ := store.Load()
	pkg, ok := lf.Get(params.Name)
	if !ok {
		return toolError("server not installed: " + params.Name)
	}

	shortName := pkg.Name
	for i := len(shortName) - 1; i >= 0; i-- {
		if shortName[i] == '/' {
			shortName = shortName[i+1:]
			break
		}
	}

	for _, clientName := range pkg.ConfiguredIn {
		for _, c := range config.KnownClients() {
			if c.Name == clientName {
				c.RemoveServer(shortName)
			}
		}
	}

	lf.Remove(params.Name)
	lf.Save()

	return toolResult(fmt.Sprintf("Removed %s", params.Name))
}

func handleInfo(args json.RawMessage, client *registry.Client) interface{} {
	var params struct {
		Name string `json:"name"`
	}
	json.Unmarshal(args, &params)

	entry, err := client.GetServer(params.Name)
	if err != nil {
		return toolError("server not found: " + err.Error())
	}

	text, _ := json.MarshalIndent(entry.Server, "", "  ")
	return toolResult(string(text))
}

func handleDoctor(args json.RawMessage) interface{} {
	var params struct {
		Server string `json:"server"`
	}
	json.Unmarshal(args, &params)

	var lines []string

	// Check client configs
	lines = append(lines, "=== Client Configuration ===")
	for _, c := range config.KnownClients() {
		issues := health.CheckClientConfig(c.Name, c.Path, c.Key)
		if len(issues) == 0 {
			lines = append(lines, fmt.Sprintf("✓ %s: OK", c.Name))
		} else {
			for _, issue := range issues {
				lines = append(lines, fmt.Sprintf("✗ %s: %s → %s", c.Name, issue.Message, issue.Suggestion))
			}
		}
	}

	// Check installed servers
	lf, err := store.Load()
	if err != nil {
		return toolError("failed to load lockfile: " + err.Error())
	}

	lines = append(lines, "", "=== Installed Servers ===")

	if len(lf.Packages) == 0 {
		lines = append(lines, "No servers installed.")
	}

	for name, pkg := range lf.Packages {
		if params.Server != "" && name != params.Server && pkg.Name != params.Server {
			continue
		}

		var report *health.Report
		if pkg.Transport.Type == "streamable-http" || pkg.Transport.Type == "sse" {
			report = health.CheckRemoteServer(name, pkg.Transport.URL)
		} else {
			command := "npx"
			cmdArgs := []string{"-y", pkg.Identifier}
			if pkg.RuntimeHint == "uvx" {
				command = "uvx"
				cmdArgs = []string{pkg.Identifier}
			}
			if pkg.Identifier != "" {
				report = health.CheckStdioServer(name, command, cmdArgs, pkg.EnvVars)
			} else {
				report = &health.Report{
					Name:   name,
					Status: health.StatusWarning,
					Issues: []health.Issue{{
						Severity:   health.StatusWarning,
						Message:    "Cannot determine how to start this server",
						Suggestion: "Re-install with: mcphub install " + name,
					}},
				}
			}
		}

		status := "✓"
		if report.Status == health.StatusError {
			status = "✗"
		} else if report.Status == health.StatusWarning {
			status = "!"
		}
		lines = append(lines, fmt.Sprintf("%s %s [%s] (%s)", status, name, report.Transport, report.Duration.Round(1e6)))
		for _, issue := range report.Issues {
			if issue.Suggestion != "" {
				lines = append(lines, fmt.Sprintf("  %s → %s", issue.Message, issue.Suggestion))
			} else {
				lines = append(lines, fmt.Sprintf("  %s", issue.Message))
			}
		}
	}

	return toolResult(fmt.Sprintf("%s", joinLines(lines)))
}

func joinLines(lines []string) string {
	result := ""
	for _, l := range lines {
		result += l + "\n"
	}
	return result
}

func toolResult(text string) map[string]interface{} {
	return map[string]interface{}{
		"content": []map[string]interface{}{
			{"type": "text", "text": text},
		},
	}
}

func toolError(msg string) map[string]interface{} {
	return map[string]interface{}{
		"content": []map[string]interface{}{
			{"type": "text", "text": "Error: " + msg},
		},
		"isError": true,
	}
}
