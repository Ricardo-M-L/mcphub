package installer

import (
	"fmt"

	"github.com/mcphub/mcphub/internal/config"
	"github.com/mcphub/mcphub/internal/registry"
)

// Result holds the output of an installation.
type Result struct {
	Entry       config.MCPServerEntry
	InstallPath string
}

// Installer handles installing a specific type of MCP server.
type Installer interface {
	Install(pkg registry.Package, envVars map[string]string) (*Result, error)
}

// Select picks the best installer for a given server.
// Priority: remote > npm > binary
func Select(server *registry.ServerDetail) (Installer, *registry.Package, *registry.Remote, error) {
	// Prefer remote (zero-install)
	for _, r := range server.Remotes {
		return &RemoteInstaller{}, nil, &r, nil
	}

	// Then try packages
	for i := range server.Packages {
		pkg := &server.Packages[i]
		switch pkg.RegistryType {
		case "npm":
			return &NpmInstaller{}, pkg, nil, nil
		}
	}

	if len(server.Packages) > 0 {
		return &NpmInstaller{}, &server.Packages[0], nil, nil
	}

	return nil, nil, nil, fmt.Errorf("no supported installation method found for %s", server.Name)
}
