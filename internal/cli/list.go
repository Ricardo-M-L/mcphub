package cli

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Ricardo-M-L/mcphub/internal/store"
	"github.com/Ricardo-M-L/mcphub/internal/ui"
	"github.com/spf13/cobra"
)

var listJSON bool

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed MCP servers",
	Aliases: []string{"ls"},
	RunE: func(cmd *cobra.Command, args []string) error {
		lf, err := store.Load()
		if err != nil {
			return fmt.Errorf("failed to load lockfile: %w", err)
		}

		if listJSON {
			data, _ := json.MarshalIndent(lf.Packages, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("\n  %s installed MCP servers\n\n", ui.Bold(fmt.Sprintf("%d", len(lf.Packages))))

		if len(lf.Packages) == 0 {
			fmt.Println(ui.Dim("  No servers installed. Run 'mcphub search' to find servers."))
			fmt.Println()
			return nil
		}

		headers := []string{"Name", "Version", "Transport", "Configured In"}
		rows := make([][]string, 0, len(lf.Packages))
		for _, pkg := range lf.Packages {
			rows = append(rows, []string{
				pkg.Name,
				pkg.Version,
				pkg.Transport.Type,
				strings.Join(pkg.ConfiguredIn, ", "),
			})
		}

		ui.PrintTable(headers, rows)
		fmt.Println()
		return nil
	},
}

func init() {
	listCmd.Flags().BoolVar(&listJSON, "json", false, "Output as JSON")
}
