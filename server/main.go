package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/mcphub/mcphub/server/crawler"
	"github.com/mcphub/mcphub/server/db"
	"github.com/mcphub/mcphub/server/handler"
)

func main() {
	port := flag.String("port", "8080", "Server port")
	dbPath := flag.String("db", "mcphub.db", "SQLite database path")
	syncInterval := flag.Duration("sync-interval", 30*time.Minute, "Registry sync interval")
	flag.Parse()

	// Open database
	database, err := db.Open(*dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer database.Close()

	// Start crawler
	c := crawler.New(database, *syncInterval)
	c.Start()
	defer c.Stop()

	// Setup routes
	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("/api/v1/servers", handler.SearchHandler(database))
	mux.HandleFunc("/api/v1/servers/", handler.ServerHandler(database))
	mux.HandleFunc("/api/v1/stats", handler.StatsHandler(database))

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// CORS middleware
	h := corsMiddleware(mux)

	// Start server
	srv := &http.Server{
		Addr:         ":" + *port,
		Handler:      h,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("MCP Hub Registry API listening on :%s", *port)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down...")
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if strings.ToUpper(r.Method) == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
