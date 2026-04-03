package handler

import (
	"encoding/json"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/Ricardo-M-L/mcphub/internal/registry"
	"github.com/Ricardo-M-L/mcphub/server/db"
)

// QualityScore represents the quality assessment of an MCP server.
type QualityScore struct {
	Overall       float64            `json:"overall"`
	Breakdown     map[string]float64 `json:"breakdown"`
	ServerName    string             `json:"serverName"`
	ComputedAt    time.Time          `json:"computedAt"`
}

// ScoreHandler handles GET /api/v1/servers/{name}/score
func ScoreHandler(database *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimPrefix(r.URL.Path, "/api/v1/servers/")
		name = strings.TrimSuffix(name, "/score")

		entry, err := database.GetServer(name)
		if err != nil || entry == nil {
			http.Error(w, "server not found", http.StatusNotFound)
			return
		}

		score := computeScore(entry)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(score)
	}
}

func computeScore(entry *registry.ServerEntry) QualityScore {
	s := entry.Server
	breakdown := make(map[string]float64)

	// 1. Completeness (0-25): has description, version, repo, packages/remotes
	completeness := 0.0
	if s.Description != "" {
		completeness += 5
	}
	if len(s.Description) > 50 {
		completeness += 5
	}
	if s.Version != "" {
		completeness += 5
	}
	if s.Repository != nil && s.Repository.URL != "" {
		completeness += 5
	}
	if len(s.Packages) > 0 || len(s.Remotes) > 0 {
		completeness += 5
	}
	breakdown["completeness"] = completeness

	// 2. Installability (0-25): has clear install path
	installability := 0.0
	if len(s.Packages) > 0 {
		installability += 15
		pkg := s.Packages[0]
		if pkg.RuntimeHint != "" {
			installability += 5
		}
		if pkg.Transport.Type != "" {
			installability += 5
		}
	} else if len(s.Remotes) > 0 {
		installability += 20
		if s.Remotes[0].URL != "" {
			installability += 5
		}
	}
	breakdown["installability"] = installability

	// 3. Documentation (0-25): env vars documented, has title
	documentation := 0.0
	if s.Title != "" {
		documentation += 10
	}
	if len(s.Packages) > 0 {
		pkg := s.Packages[0]
		if len(pkg.EnvironmentVariables) > 0 {
			documentation += 10
			allDocumented := true
			for _, ev := range pkg.EnvironmentVariables {
				if ev.Description == "" {
					allDocumented = false
				}
			}
			if allDocumented {
				documentation += 5
			}
		} else {
			documentation += 15 // No env vars needed = fully documented
		}
	}
	breakdown["documentation"] = documentation

	// 4. Security (0-25): secrets marked, transport type
	security := 10.0 // Base score
	if len(s.Packages) > 0 {
		pkg := s.Packages[0]
		for _, ev := range pkg.EnvironmentVariables {
			if ev.IsSecret {
				security += 5
				break
			}
		}
		if pkg.Transport.Type == "stdio" {
			security += 10 // Local transport = more secure
		} else {
			security += 5
		}
	} else {
		security += 10
	}
	breakdown["security"] = math.Min(security, 25)

	// Calculate overall
	overall := 0.0
	for _, v := range breakdown {
		overall += v
	}

	return QualityScore{
		Overall:    overall,
		Breakdown:  breakdown,
		ServerName: s.Name,
		ComputedAt: time.Now(),
	}
}
