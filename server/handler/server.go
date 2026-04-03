package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/mcphub/mcphub/server/db"
)

// ServerHandler handles GET /api/v1/servers/{name}
func ServerHandler(database *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract name from path: /api/v1/servers/{name...}
		name := strings.TrimPrefix(r.URL.Path, "/api/v1/servers/")
		if name == "" {
			http.Error(w, "server name required", http.StatusBadRequest)
			return
		}

		entry, err := database.GetServer(name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if entry == nil {
			http.Error(w, "server not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(entry)
	}
}
