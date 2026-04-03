package registry

// RegistryResponse is the top-level response from the MCP registry API.
type RegistryResponse struct {
	Servers  []ServerEntry    `json:"servers"`
	Metadata ResponseMetadata `json:"metadata"`
}

// ResponseMetadata contains pagination info.
type ResponseMetadata struct {
	NextCursor string `json:"nextCursor"`
	Count      int    `json:"count"`
}

// ServerEntry wraps a server detail with optional metadata.
type ServerEntry struct {
	Server ServerDetail           `json:"server"`
	Meta   map[string]interface{} `json:"_meta,omitempty"`
}

// ServerDetail holds the full MCP server information.
type ServerDetail struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Title       string      `json:"title,omitempty"`
	Version     string      `json:"version,omitempty"`
	Repository  *Repository `json:"repository,omitempty"`
	WebsiteURL  string      `json:"websiteUrl,omitempty"`
	Packages    []Package   `json:"packages,omitempty"`
	Remotes     []Remote    `json:"remotes,omitempty"`
}

// Repository points to the source code.
type Repository struct {
	URL       string `json:"url"`
	Source    string `json:"source,omitempty"`
	ID        string `json:"id,omitempty"`
	Subfolder string `json:"subfolder,omitempty"`
}

// Package describes an installable distribution of the MCP server.
type Package struct {
	RegistryType         string          `json:"registryType"`
	Identifier           string          `json:"identifier"`
	Version              string          `json:"version,omitempty"`
	RuntimeHint          string          `json:"runtimeHint,omitempty"`
	Transport            Transport       `json:"transport"`
	EnvironmentVariables []KeyValueInput `json:"environmentVariables,omitempty"`
	PackageArguments     []Argument      `json:"packageArguments,omitempty"`
	RuntimeArguments     []Argument      `json:"runtimeArguments,omitempty"`
}

// Remote describes a remotely-hosted MCP server endpoint.
type Remote struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

// Transport describes how the MCP client communicates with the server.
type Transport struct {
	Type string `json:"type"`
	URL  string `json:"url,omitempty"`
}

// KeyValueInput describes an environment variable needed by the server.
type KeyValueInput struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	IsRequired  bool   `json:"isRequired,omitempty"`
	IsSecret    bool   `json:"isSecret,omitempty"`
	Format      string `json:"format,omitempty"`
	Value       string `json:"value,omitempty"`
	Default     string `json:"default,omitempty"`
}

// Argument describes a command-line argument for the server.
type Argument struct {
	Type      string `json:"type"`
	Name      string `json:"name,omitempty"`
	ValueHint string `json:"valueHint,omitempty"`
	Value     string `json:"value,omitempty"`
	Default   string `json:"default,omitempty"`
}

// ShortName derives a short display name from the full registry name.
// e.g. "io.github.user/server-filesystem" -> "filesystem"
func (s *ServerDetail) ShortName() string {
	name := s.Name
	for i := len(name) - 1; i >= 0; i-- {
		if name[i] == '/' {
			name = name[i+1:]
			break
		}
	}
	// Strip common prefixes
	for _, prefix := range []string{"server-", "mcp-", "mcp_"} {
		if len(name) > len(prefix) && name[:len(prefix)] == prefix {
			return name[len(prefix):]
		}
	}
	return name
}
