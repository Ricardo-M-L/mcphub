package registry

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	defaultBaseURL = "https://registry.modelcontextprotocol.io"
	apiVersion     = "v0.1"
	defaultTimeout = 15 * time.Second
)

// Client talks to the MCP registry API.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a registry client with default settings.
func NewClient() *Client {
	return &Client{
		baseURL: defaultBaseURL,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}
}

// NewClientWithURL creates a registry client pointing to a custom registry.
func NewClientWithURL(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}
}

// Search queries the registry for servers matching the query string.
func (c *Client) Search(query string, limit int) ([]ServerEntry, error) {
	params := url.Values{}
	params.Set("search", query)
	params.Set("version", "latest")
	if limit > 0 {
		params.Set("limit", fmt.Sprintf("%d", limit))
	}
	return c.fetchServers(params)
}

// GetServer fetches a specific server by its registry name.
func (c *Client) GetServer(name string) (*ServerEntry, error) {
	params := url.Values{}
	params.Set("name", name)
	params.Set("version", "latest")

	servers, err := c.fetchServers(params)
	if err != nil {
		return nil, err
	}
	if len(servers) == 0 {
		return nil, fmt.Errorf("server %q not found", name)
	}
	return &servers[0], nil
}

// ListAll fetches all servers from the registry using cursor-based pagination.
func (c *Client) ListAll() ([]ServerEntry, error) {
	var all []ServerEntry
	params := url.Values{}
	params.Set("version", "latest")
	params.Set("limit", "96")

	for {
		resp, err := c.fetchResponse(params)
		if err != nil {
			return all, err
		}
		all = append(all, resp.Servers...)

		if resp.Metadata.NextCursor == "" {
			break
		}
		params.Set("cursor", resp.Metadata.NextCursor)
	}
	return all, nil
}

func (c *Client) fetchServers(params url.Values) ([]ServerEntry, error) {
	resp, err := c.fetchResponse(params)
	if err != nil {
		return nil, err
	}
	return resp.Servers, nil
}

func (c *Client) fetchResponse(params url.Values) (*RegistryResponse, error) {
	reqURL := fmt.Sprintf("%s/%s/servers?%s", c.baseURL, apiVersion, params.Encode())

	resp, err := c.httpClient.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("registry request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("registry returned %d: %s", resp.StatusCode, string(body))
	}

	var result RegistryResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse registry response: %w", err)
	}
	return &result, nil
}
