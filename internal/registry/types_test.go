package registry

import "testing"

func TestShortName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple name", "filesystem", "filesystem"},
		{"with prefix server-", "io.github.user/server-filesystem", "filesystem"},
		{"with prefix mcp-", "io.github.user/mcp-github", "github"},
		{"with prefix mcp_", "io.github.user/mcp_slack", "slack"},
		{"no prefix", "io.github.user/weather", "weather"},
		{"deep path", "io.github.org/sub/server-postgres", "postgres"},
		{"just slash", "org/tool", "tool"},
		{"no slash", "standalone", "standalone"},
		{"prefix only would be empty", "io.github.user/server-", "server-"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := ServerDetail{Name: tt.input}
			got := s.ShortName()
			if got != tt.expected {
				t.Errorf("ShortName(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
