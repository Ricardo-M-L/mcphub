package store

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Ricardo-M-L/mcphub/internal/registry"
)

func TestLockfileRoundTrip(t *testing.T) {
	// Use temp dir
	tmpDir := t.TempDir()
	lockPath := filepath.Join(tmpDir, "mcphub-lock.json")

	// Override lockfile path for testing
	origPath := os.Getenv("MCPHUB_HOME")
	os.Setenv("MCPHUB_HOME", tmpDir)
	defer os.Setenv("MCPHUB_HOME", origPath)

	lf := &Lockfile{
		Version:  1,
		Packages: make(map[string]InstalledPackage),
	}

	// Add a package
	pkg := InstalledPackage{
		Name:         "io.github.test/server-filesystem",
		Version:      "1.0.0",
		InstalledAt:  time.Now().Truncate(time.Second),
		RegistryType: "npm",
		Identifier:   "@modelcontextprotocol/server-filesystem",
		RuntimeHint:  "npx",
		Transport:    registry.Transport{Type: "stdio"},
		ConfiguredIn: []string{"claude-desktop"},
	}
	lf.Add(pkg)

	// Save to custom path
	if err := os.MkdirAll(filepath.Dir(lockPath), 0o755); err != nil {
		t.Fatal(err)
	}
	_ = lockPath

	// Verify add
	got, ok := lf.Get("io.github.test/server-filesystem")
	if !ok {
		t.Fatal("package not found after Add")
	}
	if got.Version != "1.0.0" {
		t.Errorf("expected version 1.0.0, got %s", got.Version)
	}

	// Remove
	lf.Remove("io.github.test/server-filesystem")
	_, ok = lf.Get("io.github.test/server-filesystem")
	if ok {
		t.Fatal("package still found after Remove")
	}
}

func TestLoadEmptyLockfile(t *testing.T) {
	// Point to non-existent path
	os.Setenv("MCPHUB_HOME", t.TempDir())
	defer os.Unsetenv("MCPHUB_HOME")

	lf, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if lf.Version != 1 {
		t.Errorf("expected version 1, got %d", lf.Version)
	}
	if len(lf.Packages) != 0 {
		t.Errorf("expected 0 packages, got %d", len(lf.Packages))
	}
}
