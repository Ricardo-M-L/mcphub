package handler

import (
	"encoding/json"
	"net/http"

	"github.com/mcphub/mcphub/server/db"
)

// StatsHandler handles GET /api/v1/stats
func StatsHandler(database *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stats, err := database.GetStats()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stats)
	}
}
