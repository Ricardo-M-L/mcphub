package cli

import (
	"encoding/json"
	"fmt"

	"github.com/mcphub/mcphub/internal/registry"
	"github.com/mcphub/mcphub/internal/ui"
	"github.com/spf13/cobra"
)

var infoJSON bool

var infoCmd = &cobra.Command{
	Use:   "info <server-name>",
	Short: "Show detailed information about an MCP server",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		client := registry.NewClient()

		entry, err := client.GetServer(name)
		if err != nil {
			return fmt.Errorf("server %q not found: %w", name, err)
		}

		if infoJSON {
			data, _ := json.MarshalIndent(entry.Server, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		s := entry.Server
		fmt.Println()
		fmt.Printf("  %s %s\n", ui.Bold(s.ShortName()), ui.Dim("v"+s.Version))
		fmt.Printf("  %s\n\n", s.Description)
		fmt.Printf("  %s  %s\n", ui.Dim("Registry:"), s.Name)

		if s.Repository != nil {
			fmt.Printf("  %s  %s\n", ui.Dim("Repo:    "), s.Repository.URL)
		}
		if s.WebsiteURL != "" {
			fmt.Printf("  %s  %s\n", ui.Dim("Website: "), s.WebsiteURL)
		}

		if len(s.Packages) > 0 {
			fmt.Printf("\n  %s\n", ui.Bold("Packages:"))
			for _, p := range s.Packages {
				fmt.Printf("    - %s (%s) via %s\n", p.Identifier, p.RegistryType, p.Transport.Type)
			}
		}

		if len(s.Remotes) > 0 {
			fmt.Printf("\n  %s\n", ui.Bold("Remote Endpoints:"))
			for _, r := range s.Remotes {
				fmt.Printf("    - %s (%s)\n", r.URL, r.Type)
			}
		}

		if len(s.Packages) > 0 && len(s.Packages[0].EnvironmentVariables) > 0 {
			fmt.Printf("\n  %s\n", ui.Bold("Environment Variables:"))
			for _, ev := range s.Packages[0].EnvironmentVariables {
				required := ""
				if ev.IsRequired {
					required = ui.Red(" (required)")
				}
				fmt.Printf("    - %s%s  %s\n", ui.Yellow(ev.Name), required, ui.Dim(ev.Description))
			}
		}

		fmt.Printf("\n  %s  mcphub install %s\n\n", ui.Dim("Install:"), s.Name)
		return nil
	},
}

func init() {
	infoCmd.Flags().BoolVar(&infoJSON, "json", false, "Output as JSON")
}
