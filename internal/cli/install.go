package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mcphub/mcphub/internal/config"
	"github.com/mcphub/mcphub/internal/installer"
	"github.com/mcphub/mcphub/internal/registry"
	"github.com/mcphub/mcphub/internal/store"
	"github.com/mcphub/mcphub/internal/ui"
	"github.com/spf13/cobra"
)

var installClient string

var installCmd = &cobra.Command{
	Use:   "install <server-name>",
	Short: "Install an MCP server",
	Long: `Install an MCP server and auto-configure it in your MCP clients.

Supports npm-based servers (via npx) and remote servers (via URL).
Automatically detects and configures Claude Desktop and Cursor.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		client := registry.NewClient()

		// 1. Resolve package
		ui.PrintInfo(fmt.Sprintf("Searching for %s...", ui.Cyan(name)))
		entry, err := client.GetServer(name)
		if err != nil {
			// Try fuzzy search
			entries, searchErr := client.Search(name, 5)
			if searchErr != nil || len(entries) == 0 {
				return fmt.Errorf("server %q not found in registry", name)
			}
			if len(entries) == 1 {
				entry = &entries[0]
			} else {
				fmt.Printf("\n  Multiple matches found. Did you mean:\n\n")
				for i, e := range entries {
					fmt.Printf("    %d. %s - %s\n", i+1, ui.Bold(e.Server.Name), e.Server.Description)
				}
				fmt.Println()
				return fmt.Errorf("please specify the full server name")
			}
		}

		server := &entry.Server
		fmt.Printf("  %s %s %s\n", ui.Bold(server.ShortName()), ui.Dim("v"+server.Version), ui.Dim(server.Description))

		// 2. Select installer
		inst, pkg, remote, err := installer.Select(server)
		if err != nil {
			return err
		}

		// 3. Collect environment variables
		envVars := make(map[string]string)
		var envList []registry.KeyValueInput
		if pkg != nil {
			envList = pkg.EnvironmentVariables
		}
		for _, ev := range envList {
			if ev.IsRequired {
				fmt.Printf("  %s (%s): ", ui.Yellow(ev.Name), ev.Description)
				reader := bufio.NewReader(os.Stdin)
				val, _ := reader.ReadString('\n')
				val = strings.TrimSpace(val)
				if val == "" && ev.Default != "" {
					val = ev.Default
				}
				envVars[ev.Name] = val
			}
		}

		// 4. Execute installation
		var result *installer.Result
		if remote != nil {
			result = installer.InstallRemote(remote, envVars)
			ui.PrintSuccess(fmt.Sprintf("Configured remote server at %s", remote.URL))
		} else {
			result, err = inst.Install(*pkg, envVars)
			if err != nil {
				return fmt.Errorf("installation failed: %w", err)
			}
			ui.PrintSuccess("Package ready")
		}

		// 5. Configure MCP clients
		shortName := server.ShortName()
		clients := config.DetectedClients()
		if installClient != "" {
			// Filter to specific client
			var filtered []config.ClientConfig
			for _, c := range clients {
				if c.Name == installClient {
					filtered = append(filtered, c)
				}
			}
			clients = filtered
		}

		configuredIn := []string{}
		for _, c := range clients {
			if err := c.AddServer(shortName, result.Entry); err != nil {
				ui.PrintError(fmt.Sprintf("Failed to configure %s: %s", c.Name, err))
				continue
			}
			configuredIn = append(configuredIn, c.Name)
			ui.PrintSuccess(fmt.Sprintf("Configured in %s", ui.Bold(c.Name)))
		}

		if len(configuredIn) == 0 {
			ui.PrintInfo("No MCP clients detected. You can manually add the server to your config.")
		}

		// 6. Update lockfile
		lf, err := store.Load()
		if err != nil {
			return fmt.Errorf("failed to load lockfile: %w", err)
		}

		regType := ""
		identifier := ""
		runtimeHint := ""
		if pkg != nil {
			regType = pkg.RegistryType
			identifier = pkg.Identifier
			runtimeHint = pkg.RuntimeHint
		} else if remote != nil {
			regType = "remote"
			identifier = remote.URL
		}

		lf.Add(store.InstalledPackage{
			Name:         server.Name,
			Version:      server.Version,
			InstalledAt:  time.Now(),
			RegistryType: regType,
			Identifier:   identifier,
			RuntimeHint:  runtimeHint,
			Transport:    result.Entry.ToTransport(),
			EnvVars:      envVars,
			ConfiguredIn: configuredIn,
			InstallPath:  result.InstallPath,
		})

		if err := lf.Save(); err != nil {
			return fmt.Errorf("failed to save lockfile: %w", err)
		}

		fmt.Printf("\n  %s installed successfully!\n\n", ui.Green(shortName))
		return nil
	},
}

func init() {
	installCmd.Flags().StringVar(&installClient, "client", "", "Target specific client (claude-desktop, cursor)")
}
