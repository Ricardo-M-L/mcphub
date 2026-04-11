package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "0.1.0"

// rootCmd is the base command for mcphub.
var rootCmd = &cobra.Command{
	Use:   "mcphub",
	Short: "The package manager for MCP servers",
	Long: `MCP Hub - The package manager for Model Context Protocol servers.

Search, install, and manage MCP servers with a single command.
Auto-configures Claude Desktop, Cursor, and other MCP clients.`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(infoCmd)
	rootCmd.AddCommand(publishCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(doctorCmd)
	rootCmd.AddCommand(versionCmd)
}

// Execute runs the CLI.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
