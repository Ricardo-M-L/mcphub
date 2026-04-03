package installer

import (
	"fmt"
	"os/exec"

	"github.com/Ricardo-M-L/mcphub/internal/config"
	"github.com/Ricardo-M-L/mcphub/internal/registry"
)

// NpmInstaller handles npm/npx-based MCP servers.
// Most npm MCP servers use npx which downloads and caches automatically.
type NpmInstaller struct{}

func (n *NpmInstaller) Install(pkg registry.Package, envVars map[string]string) (*Result, error) {
	runtime := pkg.RuntimeHint
	if runtime == "" {
		runtime = "npx"
	}

	// Verify runtime is available
	if _, err := exec.LookPath(runtime); err != nil {
		return nil, fmt.Errorf("%s is not installed. Please install Node.js first: https://nodejs.org", runtime)
	}

	// Build args
	args := []string{}
	if runtime == "npx" {
		args = append(args, "-y")
	}
	args = append(args, pkg.Identifier)

	// Append package arguments
	for _, arg := range pkg.PackageArguments {
		if arg.Value != "" {
			args = append(args, arg.Value)
		} else if arg.Default != "" {
			args = append(args, arg.Default)
		}
	}

	return &Result{
		Entry: config.MCPServerEntry{
			Command: runtime,
			Args:    args,
			Env:     envVars,
		},
	}, nil
}
