package registry

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// GitHubSearchResult represents the GitHub search API response.
type GitHubSearchResult struct {
	TotalCount int          `json:"total_count"`
	Items      []GitHubRepo `json:"items"`
}

// GitHubRepo represents a GitHub repository from search results.
type GitHubRepo struct {
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	HTMLURL     string `json:"html_url"`
	StarCount   int    `json:"stargazers_count"`
	Language    string `json:"language"`
	UpdatedAt   string `json:"updated_at"`
}

// SearchGitHub searches GitHub for MCP server repositories.
// Used as a fallback when the official registry returns no results.
func (c *Client) SearchGitHub(query string, limit int) ([]ServerEntry, error) {
	if limit <= 0 {
		limit = 10
	}

	// Search GitHub for MCP server repos
	ghQuery := fmt.Sprintf("mcp server %s in:name,description,readme", query)
	params := url.Values{}
	params.Set("q", ghQuery)
	params.Set("sort", "stars")
	params.Set("order", "desc")
	params.Set("per_page", fmt.Sprintf("%d", limit))

	reqURL := fmt.Sprintf("https://api.github.com/search/repositories?%s", params.Encode())

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "mcphub")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GitHub search failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API returned %d: %s", resp.StatusCode, string(body))
	}

	var result GitHubSearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse GitHub response: %w", err)
	}

	// Convert GitHub repos to ServerEntry format
	var entries []ServerEntry
	for _, repo := range result.Items {
		desc := repo.Description
		if desc == "" {
			desc = "MCP server from GitHub"
		}

		entries = append(entries, ServerEntry{
			Server: ServerDetail{
				Name:        fmt.Sprintf("github.com/%s", repo.FullName),
				Description: fmt.Sprintf("%s (%d stars)", desc, repo.StarCount),
				Repository: &Repository{
					URL:    repo.HTMLURL,
					Source: "github",
				},
			},
			Meta: map[string]interface{}{
				"source":   "github",
				"stars":    repo.StarCount,
				"language": repo.Language,
			},
		})
	}

	return entries, nil
}

// SearchAll searches both the official registry and GitHub, merging results.
func (c *Client) SearchAll(query string, limit int) ([]ServerEntry, error) {
	// First search official registry
	entries, err := c.Search(query, limit)
	if err != nil {
		entries = nil // Don't fail, try GitHub
	}

	// If registry returned few results, supplement with GitHub
	if len(entries) < 3 {
		ghEntries, ghErr := c.SearchGitHub(query, limit-len(entries))
		if ghErr == nil {
			entries = append(entries, ghEntries...)
		}
	}

	// Cap at limit
	if len(entries) > limit {
		entries = entries[:limit]
	}

	return entries, nil
}
