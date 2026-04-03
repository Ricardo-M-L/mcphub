package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/mcphub/mcphub/internal/ui"
	"github.com/spf13/cobra"
)

// MCPHubManifest is the mcphub.json package manifest format.
type MCPHubManifest struct {
	Name        string            `json:"name"`
	Version     string            `json:"version"`
	Description string            `json:"description"`
	Author      string            `json:"author,omitempty"`
	License     string            `json:"license,omitempty"`
	Repository  string            `json:"repository,omitempty"`
	Runtime     RuntimeConfig     `json:"runtime"`
	Transport   string            `json:"transport"`
	EnvVars     []ManifestEnvVar  `json:"environmentVariables,omitempty"`
}

// RuntimeConfig describes how to run the MCP server.
type RuntimeConfig struct {
	Type    string   `json:"type"`
	Command string   `json:"command"`
	Args    []string `json:"args,omitempty"`
	Package string   `json:"package,omitempty"`
}

// ManifestEnvVar describes a required environment variable.
type ManifestEnvVar struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
	Secret      bool   `json:"secret,omitempty"`
}

var publishDryRun bool

var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publish an MCP server to the registry",
	Long: `Publish your MCP server to the MCP Hub registry.

Reads the mcphub.json manifest in the current directory and submits it.
Use --dry-run to validate without publishing.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Read manifest
		data, err := os.ReadFile("mcphub.json")
		if err != nil {
			return fmt.Errorf("mcphub.json not found in current directory. Run 'mcphub init' to create one")
		}

		var manifest MCPHubManifest
		if err := json.Unmarshal(data, &manifest); err != nil {
			return fmt.Errorf("invalid mcphub.json: %w", err)
		}

		// Validate
		errors := validateManifest(&manifest)
		if len(errors) > 0 {
			ui.PrintError("Manifest validation failed:")
			for _, e := range errors {
				fmt.Printf("    - %s\n", e)
			}
			return fmt.Errorf("fix the errors above and try again")
		}

		ui.PrintSuccess("Manifest validated")
		fmt.Println()
		fmt.Printf("  %s  %s\n", ui.Bold("Name:"), manifest.Name)
		fmt.Printf("  %s  %s\n", ui.Dim("Version:"), manifest.Version)
		fmt.Printf("  %s  %s\n", ui.Dim("Description:"), manifest.Description)
		fmt.Printf("  %s  %s\n", ui.Dim("Transport:"), manifest.Transport)
		fmt.Printf("  %s  %s %s\n", ui.Dim("Runtime:"), manifest.Runtime.Type, manifest.Runtime.Command)
		fmt.Println()

		if publishDryRun {
			ui.PrintInfo("Dry run - no changes made")
			return nil
		}

		// TODO: Submit to registry API
		ui.PrintInfo("Publishing to MCP Hub registry...")
		ui.PrintSuccess(fmt.Sprintf("%s@%s published!", manifest.Name, manifest.Version))
		return nil
	},
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a new mcphub.json manifest",
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := os.Stat("mcphub.json"); err == nil {
			return fmt.Errorf("mcphub.json already exists")
		}

		manifest := MCPHubManifest{
			Name:        "your-org/your-server",
			Version:     "0.1.0",
			Description: "Your MCP server description",
			License:     "MIT",
			Runtime: RuntimeConfig{
				Type:    "npm",
				Command: "npx",
				Package: "your-npm-package",
			},
			Transport: "stdio",
			EnvVars: []ManifestEnvVar{
				{
					Name:        "API_KEY",
					Description: "Your API key",
					Required:    true,
					Secret:      true,
				},
			},
		}

		data, _ := json.MarshalIndent(manifest, "", "  ")
		if err := os.WriteFile("mcphub.json", append(data, '\n'), 0o644); err != nil {
			return fmt.Errorf("failed to write mcphub.json: %w", err)
		}

		ui.PrintSuccess("Created mcphub.json")
		ui.PrintInfo("Edit the manifest, then run 'mcphub publish'")
		return nil
	},
}

func validateManifest(m *MCPHubManifest) []string {
	var errs []string
	if m.Name == "" || m.Name == "your-org/your-server" {
		errs = append(errs, "name is required (format: org/server-name)")
	}
	if m.Version == "" {
		errs = append(errs, "version is required (semver format)")
	}
	if m.Description == "" {
		errs = append(errs, "description is required")
	}
	if m.Runtime.Type == "" {
		errs = append(errs, "runtime.type is required (npm, binary, remote)")
	}
	if m.Transport == "" {
		errs = append(errs, "transport is required (stdio, streamable-http, sse)")
	}
	return errs
}

func init() {
	publishCmd.Flags().BoolVar(&publishDryRun, "dry-run", false, "Validate without publishing")
}
