package cli

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Ricardo-M-L/mcphub/internal/config"
	"github.com/Ricardo-M-L/mcphub/internal/health"
	"github.com/Ricardo-M-L/mcphub/internal/store"
	"github.com/Ricardo-M-L/mcphub/internal/ui"
	"github.com/spf13/cobra"
)

var doctorJSON bool
var doctorServer string

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Diagnose MCP server health and configuration issues",
	Long: `Check all installed MCP servers and client configurations for problems.

Performs:
  - Client config validation (JSON syntax, missing fields)
  - Server connectivity tests (spawn + MCP handshake for stdio, HTTP for remote)
  - Environment variable checks
  - Runtime dependency checks (node, python, docker)

Examples:
  mcphub doctor                 # Check everything
  mcphub doctor --server xxx    # Check a specific server
  mcphub doctor --json          # JSON output`,
	RunE: func(cmd *cobra.Command, args []string) error {
		start := time.Now()

		if !doctorJSON {
			fmt.Println()
			fmt.Printf("  %s\n\n", ui.Bold("MCP Hub Doctor"))
		}

		var allReports []*health.Report
		var allConfigIssues []health.Issue

		// 1. Check client configurations
		if !doctorJSON {
			fmt.Printf("  %s Checking client configurations...\n", ui.Cyan("i"))
		}

		clients := config.KnownClients()
		for _, c := range clients {
			issues := health.CheckClientConfig(c.Name, c.Path, c.Key)
			if len(issues) > 0 {
				allConfigIssues = append(allConfigIssues, issues...)
			}
			if !doctorJSON {
				hasError := false
				for _, issue := range issues {
					if issue.Severity == health.StatusError {
						hasError = true
					}
				}
				if hasError {
					fmt.Printf("    %s %s\n", ui.Red("✗"), c.Name)
					for _, issue := range issues {
						printIssue(issue)
					}
				} else {
					fmt.Printf("    %s %s\n", ui.Green("✓"), c.Name)
				}
			}
		}

		if !doctorJSON {
			fmt.Println()
		}

		// 2. Check installed servers
		lf, err := store.Load()
		if err != nil {
			return fmt.Errorf("failed to load lockfile: %w", err)
		}

		if len(lf.Packages) == 0 {
			if !doctorJSON {
				fmt.Printf("  %s No MCP servers installed. Run 'mcphub search' to find servers.\n\n", ui.Dim("i"))
			}
			return nil
		}

		if !doctorJSON {
			fmt.Printf("  %s Checking %d installed server(s)...\n\n", ui.Cyan("i"), len(lf.Packages))
		}

		for name, pkg := range lf.Packages {
			// Skip if --server flag is set and doesn't match
			if doctorServer != "" && name != doctorServer && pkg.Name != doctorServer {
				continue
			}

			var report *health.Report

			if pkg.Transport.Type == "streamable-http" || pkg.Transport.Type == "sse" {
				// Remote server — HTTP check
				report = health.CheckRemoteServer(name, pkg.Transport.URL)
			} else {
				// Stdio server — spawn and test
				command := ""
				var args []string

				// Reconstruct command from what we know
				switch pkg.RuntimeHint {
				case "npx":
					command = "npx"
					args = []string{"-y", pkg.Identifier}
				case "uvx":
					command = "uvx"
					args = []string{pkg.Identifier}
				default:
					if pkg.Identifier != "" {
						command = "npx"
						args = []string{"-y", pkg.Identifier}
					} else {
						report = &health.Report{
							Name:      name,
							Status:    health.StatusWarning,
							Transport: "stdio",
							Issues: []health.Issue{{
								Severity:   health.StatusWarning,
								Message:    "Cannot determine how to start this server",
								Suggestion: "Re-install with: mcphub install " + name,
							}},
						}
					}
				}

				if report == nil {
					report = health.CheckStdioServer(name, command, args, pkg.EnvVars)
				}
			}

			allReports = append(allReports, report)

			if !doctorJSON {
				printReport(report)
			}
		}

		// JSON output
		if doctorJSON {
			output := map[string]interface{}{
				"configIssues": allConfigIssues,
				"servers":      allReports,
				"duration":     time.Since(start).String(),
			}
			data, _ := json.MarshalIndent(output, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		// Summary
		fmt.Println()
		okCount := 0
		warnCount := 0
		errCount := 0
		for _, r := range allReports {
			switch r.Status {
			case health.StatusOK:
				okCount++
			case health.StatusWarning:
				warnCount++
			case health.StatusError:
				errCount++
			}
		}

		fmt.Printf("  %s  %s healthy", ui.Bold("Summary:"), ui.Green(fmt.Sprintf("%d", okCount)))
		if warnCount > 0 {
			fmt.Printf(", %s warnings", ui.Yellow(fmt.Sprintf("%d", warnCount)))
		}
		if errCount > 0 {
			fmt.Printf(", %s errors", ui.Red(fmt.Sprintf("%d", errCount)))
		}
		fmt.Printf("  (%s)\n\n", time.Since(start).Round(time.Millisecond))

		return nil
	},
}

func printReport(r *health.Report) {
	var icon string
	switch r.Status {
	case health.StatusOK:
		icon = ui.Green("✓")
	case health.StatusWarning:
		icon = ui.Yellow("!")
	case health.StatusError:
		icon = ui.Red("✗")
	}

	fmt.Printf("    %s %s %s  (%s)\n", icon, ui.Bold(r.Name), ui.Dim("["+r.Transport+"]"), r.Duration.Round(time.Millisecond))

	for _, issue := range r.Issues {
		printIssue(issue)
	}
	fmt.Println()
}

func printIssue(issue health.Issue) {
	switch issue.Severity {
	case health.StatusOK:
		fmt.Printf("      %s %s\n", ui.Green("✓"), issue.Message)
	case health.StatusWarning:
		fmt.Printf("      %s %s\n", ui.Yellow("!"), issue.Message)
		if issue.Suggestion != "" {
			fmt.Printf("        %s %s\n", ui.Dim("→"), issue.Suggestion)
		}
	case health.StatusError:
		fmt.Printf("      %s %s\n", ui.Red("✗"), issue.Message)
		if issue.Suggestion != "" {
			fmt.Printf("        %s %s\n", ui.Dim("→"), issue.Suggestion)
		}
	}
}

func init() {
	doctorCmd.Flags().BoolVar(&doctorJSON, "json", false, "Output as JSON")
	doctorCmd.Flags().StringVar(&doctorServer, "server", "", "Check a specific server only")
}
