package installer

import (
	"github.com/Ricardo-M-L/mcphub/internal/config"
	"github.com/Ricardo-M-L/mcphub/internal/registry"
)

// RemoteInstaller handles remotely-hosted MCP servers.
// No download needed, just configure the URL.
type RemoteInstaller struct{}

func (r *RemoteInstaller) Install(pkg registry.Package, envVars map[string]string) (*Result, error) {
	return nil, nil
}

// InstallRemote configures a remote MCP server by URL.
func InstallRemote(remote *registry.Remote, envVars map[string]string) *Result {
	return &Result{
		Entry: config.MCPServerEntry{
			URL: remote.URL,
			Env: envVars,
		},
	}
}
