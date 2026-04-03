package cli

import (
	"fmt"

	"github.com/mcphub/mcphub/internal/config"
	"github.com/mcphub/mcphub/internal/store"
	"github.com/mcphub/mcphub/internal/ui"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:     "remove <server-name>",
	Short:   "Remove an installed MCP server",
	Aliases: []string{"rm", "uninstall"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		lf, err := store.Load()
		if err != nil {
			return fmt.Errorf("failed to load lockfile: %w", err)
		}

		pkg, ok := lf.Get(name)
		if !ok {
			// Try short name match
			for key, p := range lf.Packages {
				if p.Name == name || key == name {
					pkg = p
					name = key
					ok = true
					break
				}
			}
			if !ok {
				return fmt.Errorf("server %q is not installed", name)
			}
		}

		// Remove from client configs
		shortName := ""
		for i := len(pkg.Name) - 1; i >= 0; i-- {
			if pkg.Name[i] == '/' {
				shortName = pkg.Name[i+1:]
				break
			}
		}
		if shortName == "" {
			shortName = pkg.Name
		}

		for _, clientName := range pkg.ConfiguredIn {
			for _, c := range config.KnownClients() {
				if c.Name == clientName {
					if err := c.RemoveServer(shortName); err != nil {
						ui.PrintError(fmt.Sprintf("Failed to remove from %s: %s", clientName, err))
					} else {
						ui.PrintSuccess(fmt.Sprintf("Removed from %s", clientName))
					}
				}
			}
		}

		// Remove from lockfile
		lf.Remove(name)
		if err := lf.Save(); err != nil {
			return fmt.Errorf("failed to save lockfile: %w", err)
		}

		ui.PrintSuccess(fmt.Sprintf("%s removed", ui.Bold(name)))
		return nil
	},
}
