package registry

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientSearch(t *testing.T) {
	// Mock server
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("search") != "filesystem" {
			t.Errorf("expected search=filesystem, got %s", r.URL.Query().Get("search"))
		}

		resp := RegistryResponse{
			Servers: []ServerEntry{
				{
					Server: ServerDetail{
						Name:        "io.github.test/server-filesystem",
						Description: "A test filesystem server",
						Version:     "1.0.0",
					},
				},
			},
			Metadata: ResponseMetadata{Count: 1},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	client := NewClientWithURL(srv.URL)
	entries, err := client.Search("filesystem", 10)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 result, got %d", len(entries))
	}
	if entries[0].Server.Name != "io.github.test/server-filesystem" {
		t.Errorf("unexpected name: %s", entries[0].Server.Name)
	}
}

func TestClientGetServer(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := RegistryResponse{
			Servers: []ServerEntry{
				{
					Server: ServerDetail{
						Name:        "io.github.test/server-github",
						Description: "GitHub MCP server",
						Version:     "2.0.0",
					},
				},
			},
			Metadata: ResponseMetadata{Count: 1},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	client := NewClientWithURL(srv.URL)
	entry, err := client.GetServer("io.github.test/server-github")
	if err != nil {
		t.Fatalf("GetServer failed: %v", err)
	}
	if entry.Server.Version != "2.0.0" {
		t.Errorf("unexpected version: %s", entry.Server.Version)
	}
}

func TestClientGetServerNotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := RegistryResponse{
			Servers:  []ServerEntry{},
			Metadata: ResponseMetadata{Count: 0},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	client := NewClientWithURL(srv.URL)
	_, err := client.GetServer("nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent server")
	}
}
