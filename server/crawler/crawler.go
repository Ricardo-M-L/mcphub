package crawler

import (
	"log"
	"time"

	"github.com/mcphub/mcphub/internal/registry"
	"github.com/mcphub/mcphub/server/db"
)

// Crawler syncs MCP servers from the upstream registry into the local database.
type Crawler struct {
	client   *registry.Client
	database *db.DB
	interval time.Duration
	stopCh   chan struct{}
}

// New creates a new Crawler.
func New(database *db.DB, interval time.Duration) *Crawler {
	return &Crawler{
		client:   registry.NewClient(),
		database: database,
		interval: interval,
		stopCh:   make(chan struct{}),
	}
}

// Start begins periodic syncing in a goroutine.
func (c *Crawler) Start() {
	go func() {
		// Initial sync
		c.sync()

		ticker := time.NewTicker(c.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				c.sync()
			case <-c.stopCh:
				return
			}
		}
	}()
}

// Stop halts the crawler.
func (c *Crawler) Stop() {
	close(c.stopCh)
}

func (c *Crawler) sync() {
	start := time.Now()
	log.Println("[crawler] syncing from upstream registry...")

	entries, err := c.client.ListAll()
	if err != nil {
		log.Printf("[crawler] sync failed: %v", err)
		return
	}

	count := 0
	for _, entry := range entries {
		if err := c.database.UpsertServer(entry); err != nil {
			log.Printf("[crawler] failed to upsert %s: %v", entry.Server.Name, err)
			continue
		}
		count++
	}

	duration := time.Since(start)
	c.database.LogSync(count, duration)
	log.Printf("[crawler] synced %d servers in %s", count, duration.Round(time.Millisecond))
}
