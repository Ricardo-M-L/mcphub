package store

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/mcphub/mcphub/internal/platform"
	"github.com/mcphub/mcphub/internal/registry"
)

// InstalledPackage tracks a locally installed MCP server.
type InstalledPackage struct {
	Name         string            `json:"name"`
	Version      string            `json:"version"`
	InstalledAt  time.Time         `json:"installedAt"`
	RegistryType string            `json:"registryType"`
	Identifier   string            `json:"identifier"`
	RuntimeHint  string            `json:"runtimeHint,omitempty"`
	Transport    registry.Transport `json:"transport"`
	EnvVars      map[string]string `json:"envVars,omitempty"`
	ConfiguredIn []string          `json:"configuredIn"`
	InstallPath  string            `json:"installPath,omitempty"`
}

// Lockfile persists installed packages at ~/.mcphub/mcphub-lock.json.
type Lockfile struct {
	Version  int                         `json:"version"`
	Packages map[string]InstalledPackage  `json:"packages"`
}

// Load reads the lockfile from disk. Returns an empty lockfile if not found.
func Load() (*Lockfile, error) {
	path := platform.LockfilePath()
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &Lockfile{Version: 1, Packages: make(map[string]InstalledPackage)}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read lockfile: %w", err)
	}

	var lf Lockfile
	if err := json.Unmarshal(data, &lf); err != nil {
		return nil, fmt.Errorf("failed to parse lockfile: %w", err)
	}
	if lf.Packages == nil {
		lf.Packages = make(map[string]InstalledPackage)
	}
	return &lf, nil
}

// Save writes the lockfile to disk, creating directories as needed.
func (lf *Lockfile) Save() error {
	path := platform.LockfilePath()
	if err := os.MkdirAll(platform.MCPHubDir(), 0o755); err != nil {
		return fmt.Errorf("failed to create mcphub directory: %w", err)
	}

	data, err := json.MarshalIndent(lf, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal lockfile: %w", err)
	}
	return os.WriteFile(path, data, 0o644)
}

// Add records an installed package.
func (lf *Lockfile) Add(pkg InstalledPackage) {
	lf.Packages[pkg.Name] = pkg
}

// Remove deletes a package from the lockfile.
func (lf *Lockfile) Remove(name string) {
	delete(lf.Packages, name)
}

// Get retrieves an installed package by name.
func (lf *Lockfile) Get(name string) (InstalledPackage, bool) {
	pkg, ok := lf.Packages[name]
	return pkg, ok
}
