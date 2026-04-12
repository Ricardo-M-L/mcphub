package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Ricardo-M-L/mcphub/internal/platform"
	"github.com/Ricardo-M-L/mcphub/internal/ui"
	"github.com/spf13/cobra"
)

// Rating represents a user's rating of an MCP server.
type Rating struct {
	Server    string    `json:"server"`
	Score     int       `json:"score"`
	Comment   string    `json:"comment,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}

// RatingsFile stores all local ratings.
type RatingsFile struct {
	Ratings []Rating `json:"ratings"`
}

var rateScore int
var rateComment string

var rateCmd = &cobra.Command{
	Use:   "rate <server-name>",
	Short: "Rate an installed MCP server (1-5 stars)",
	Long: `Rate an MCP server you've used. Ratings are stored locally
and can be shared to help others discover quality servers.

Examples:
  mcphub rate io.github.xxx/server-filesystem --score 5
  mcphub rate io.github.xxx/server-filesystem --score 4 --comment "Great but slow on large dirs"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		if rateScore < 1 || rateScore > 5 {
			return fmt.Errorf("score must be between 1 and 5")
		}

		rating := Rating{
			Server:    name,
			Score:     rateScore,
			Comment:   rateComment,
			CreatedAt: time.Now(),
		}

		// Load existing ratings
		ratingsPath := filepath.Join(platform.MCPHubDir(), "ratings.json")
		var rf RatingsFile

		data, err := os.ReadFile(ratingsPath)
		if err == nil {
			json.Unmarshal(data, &rf)
		}

		// Remove existing rating for same server
		filtered := make([]Rating, 0, len(rf.Ratings))
		for _, r := range rf.Ratings {
			if r.Server != name {
				filtered = append(filtered, r)
			}
		}
		filtered = append(filtered, rating)
		rf.Ratings = filtered

		// Save
		os.MkdirAll(platform.MCPHubDir(), 0o755)
		out, _ := json.MarshalIndent(rf, "", "  ")
		if err := os.WriteFile(ratingsPath, out, 0o644); err != nil {
			return fmt.Errorf("failed to save rating: %w", err)
		}

		stars := ""
		for i := 0; i < rateScore; i++ {
			stars += "★"
		}
		for i := rateScore; i < 5; i++ {
			stars += "☆"
		}

		ui.PrintSuccess(fmt.Sprintf("Rated %s %s (%d/5)", ui.Bold(name), ui.Yellow(stars), rateScore))
		if rateComment != "" {
			fmt.Printf("    %s %s\n", ui.Dim("Comment:"), rateComment)
		}

		return nil
	},
}

var ratingsCmd = &cobra.Command{
	Use:   "ratings",
	Short: "Show your ratings of MCP servers",
	RunE: func(cmd *cobra.Command, args []string) error {
		ratingsPath := filepath.Join(platform.MCPHubDir(), "ratings.json")
		data, err := os.ReadFile(ratingsPath)
		if err != nil {
			fmt.Println(ui.Dim("  No ratings yet. Rate a server with: mcphub rate <server> --score 5"))
			return nil
		}

		var rf RatingsFile
		json.Unmarshal(data, &rf)

		if len(rf.Ratings) == 0 {
			fmt.Println(ui.Dim("  No ratings yet."))
			return nil
		}

		fmt.Printf("\n  %s\n\n", ui.Bold("Your Ratings"))
		for _, r := range rf.Ratings {
			stars := ""
			for i := 0; i < r.Score; i++ {
				stars += "★"
			}
			for i := r.Score; i < 5; i++ {
				stars += "☆"
			}
			fmt.Printf("    %s %s  %s\n", ui.Yellow(stars), ui.Bold(r.Server), ui.Dim(r.CreatedAt.Format("2006-01-02")))
			if r.Comment != "" {
				fmt.Printf("      %s\n", r.Comment)
			}
		}
		fmt.Println()

		return nil
	},
}

func init() {
	rateCmd.Flags().IntVar(&rateScore, "score", 0, "Rating score (1-5)")
	rateCmd.MarkFlagRequired("score")
	rateCmd.Flags().StringVar(&rateComment, "comment", "", "Optional comment")
}
