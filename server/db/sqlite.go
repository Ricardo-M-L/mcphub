package db

import (
	"database/sql"
	_ "embed"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Ricardo-M-L/mcphub/internal/registry"
	_ "modernc.org/sqlite"
)

//go:embed migrations/001_init.sql
var initSQL string

// DB wraps SQLite operations for the registry server.
type DB struct {
	conn *sql.DB
}

// Open creates or opens a SQLite database and runs migrations.
func Open(path string) (*DB, error) {
	conn, err := sql.Open("sqlite", path+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	if _, err := conn.Exec(initSQL); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}
	return &DB{conn: conn}, nil
}

// Close closes the database connection.
func (d *DB) Close() error {
	return d.conn.Close()
}

// UpsertServer inserts or updates a server entry.
func (d *DB) UpsertServer(entry registry.ServerEntry) error {
	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal server: %w", err)
	}
	_, err = d.conn.Exec(`
		INSERT INTO servers (name, version, title, description, data, updated_at)
		VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(name) DO UPDATE SET
			version = excluded.version,
			title = excluded.title,
			description = excluded.description,
			data = excluded.data,
			updated_at = CURRENT_TIMESTAMP
	`, entry.Server.Name, entry.Server.Version, entry.Server.Title, entry.Server.Description, string(data))
	return err
}

// SearchResult holds a search result with relevance info.
type SearchResult struct {
	Entry registry.ServerEntry
	Rank  float64
}

// Search performs full-text search on the servers table.
func (d *DB) Search(query string, limit, offset int) ([]registry.ServerEntry, int, error) {
	if limit <= 0 {
		limit = 20
	}

	// Count total matches
	var total int
	err := d.conn.QueryRow(`
		SELECT COUNT(*) FROM servers_fts WHERE servers_fts MATCH ?
	`, query).Scan(&total)
	if err != nil {
		// Fall back to LIKE search if FTS fails
		return d.searchLike(query, limit, offset)
	}

	rows, err := d.conn.Query(`
		SELECT s.data FROM servers s
		JOIN servers_fts fts ON s.rowid = fts.rowid
		WHERE servers_fts MATCH ?
		ORDER BY rank
		LIMIT ? OFFSET ?
	`, query, limit, offset)
	if err != nil {
		return d.searchLike(query, limit, offset)
	}
	defer rows.Close()

	return d.scanEntries(rows), total, nil
}

func (d *DB) searchLike(query string, limit, offset int) ([]registry.ServerEntry, int, error) {
	pattern := "%" + query + "%"

	var total int
	d.conn.QueryRow(`
		SELECT COUNT(*) FROM servers WHERE name LIKE ? OR title LIKE ? OR description LIKE ?
	`, pattern, pattern, pattern).Scan(&total)

	rows, err := d.conn.Query(`
		SELECT data FROM servers
		WHERE name LIKE ? OR title LIKE ? OR description LIKE ?
		ORDER BY name
		LIMIT ? OFFSET ?
	`, pattern, pattern, pattern, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	return d.scanEntries(rows), total, nil
}

// GetServer retrieves a single server by name.
func (d *DB) GetServer(name string) (*registry.ServerEntry, error) {
	var data string
	err := d.conn.QueryRow(`SELECT data FROM servers WHERE name = ?`, name).Scan(&data)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var entry registry.ServerEntry
	if err := json.Unmarshal([]byte(data), &entry); err != nil {
		return nil, err
	}
	return &entry, nil
}

// ListAll returns all servers.
func (d *DB) ListAll(limit, offset int) ([]registry.ServerEntry, int, error) {
	if limit <= 0 {
		limit = 50
	}
	var total int
	d.conn.QueryRow(`SELECT COUNT(*) FROM servers`).Scan(&total)

	rows, err := d.conn.Query(`SELECT data FROM servers ORDER BY name LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	return d.scanEntries(rows), total, nil
}

// Stats returns registry statistics.
type Stats struct {
	TotalServers int       `json:"totalServers"`
	LastSync     time.Time `json:"lastSync"`
	SyncDuration int       `json:"syncDurationMs"`
}

// GetStats returns registry statistics.
func (d *DB) GetStats() (*Stats, error) {
	stats := &Stats{}
	d.conn.QueryRow(`SELECT COUNT(*) FROM servers`).Scan(&stats.TotalServers)

	var syncedAt sql.NullTime
	var duration sql.NullInt64
	d.conn.QueryRow(`SELECT synced_at, duration_ms FROM sync_log ORDER BY id DESC LIMIT 1`).Scan(&syncedAt, &duration)
	if syncedAt.Valid {
		stats.LastSync = syncedAt.Time
	}
	if duration.Valid {
		stats.SyncDuration = int(duration.Int64)
	}
	return stats, nil
}

// LogSync records a sync event.
func (d *DB) LogSync(count int, duration time.Duration) error {
	_, err := d.conn.Exec(`INSERT INTO sync_log (server_count, duration_ms) VALUES (?, ?)`,
		count, duration.Milliseconds())
	return err
}

func (d *DB) scanEntries(rows *sql.Rows) []registry.ServerEntry {
	var entries []registry.ServerEntry
	for rows.Next() {
		var data string
		if err := rows.Scan(&data); err != nil {
			continue
		}
		var entry registry.ServerEntry
		if err := json.Unmarshal([]byte(data), &entry); err != nil {
			continue
		}
		entries = append(entries, entry)
	}
	return entries
}
