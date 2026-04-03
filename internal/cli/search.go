package cli

import (
	"encoding/json"
	"fmt"

	"github.com/Ricardo-M-L/mcphub/internal/registry"
	"github.com/Ricardo-M-L/mcphub/internal/ui"
	"github.com/spf13/cobra"
)

var searchJSON bool
var searchLimit int

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search for MCP servers",
	Long:  "Search the MCP registry for servers matching your query.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := args[0]
		client := registry.NewClient()

		entries, err := client.Search(query, searchLimit)
		if err != nil {
			return fmt.Errorf("search failed: %w", err)
		}

		if searchJSON {
			data, _ := json.MarshalIndent(entries, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("\n  Found %s servers matching %s\n\n",
			ui.Bold(fmt.Sprintf("%d", len(entries))),
			ui.Cyan(fmt.Sprintf("%q", query)),
		)

		headers := []string{"Name", "Description", "Version", "Transport"}
		rows := make([][]string, 0, len(entries))
		for _, e := range entries {
			s := e.Server
			transport := "stdio"
			if len(s.Remotes) > 0 {
				transport = s.Remotes[0].Type
			} else if len(s.Packages) > 0 {
				transport = s.Packages[0].Transport.Type
			}
			rows = append(rows, []string{
				s.ShortName(),
				ui.Truncate(s.Description, 50),
				s.Version,
				transport,
			})
		}

		ui.PrintTable(headers, rows)
		fmt.Println()
		return nil
	},
}

func init() {
	searchCmd.Flags().BoolVar(&searchJSON, "json", false, "Output results as JSON")
	searchCmd.Flags().IntVar(&searchLimit, "limit", 20, "Maximum number of results")
}
