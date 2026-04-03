package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of mcphub",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("mcphub v%s\n", version)
	},
}
