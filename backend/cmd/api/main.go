package main

import (
	"fmt"
	"log"
	"net/http"

	"tracker/internal/api"
	"tracker/internal/store"
)

func main() {
	db, err := store.New("localhost", 5432, "drew", "password123", "uptime_monitor")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	srv := api.New(db)

	addr := ":8080"
	fmt.Printf("🚀 Sentinel API server listening on %s\n", addr)
	if err := http.ListenAndServe(addr, srv.Handler()); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
