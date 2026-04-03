package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Ricardo-M-L/mcphub/server/db"
)

// SearchHandler handles GET /api/v1/servers
func SearchHandler(database *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

		if limit <= 0 {
			limit = 20
		}

		var result interface{}
		var total int

		if query != "" {
			entries, t, err := database.Search(query, limit, offset)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			result = entries
			total = t
		} else {
			entries, t, err := database.ListAll(limit, offset)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			result = entries
			total = t
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"servers": result,
			"metadata": map[string]interface{}{
				"total":  total,
				"limit":  limit,
				"offset": offset,
			},
		})
	}
}
