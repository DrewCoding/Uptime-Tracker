package main

import (
	"fmt"
	"os"
	"sync"

	"tracker/internal/monitor"
	// "tracker/internal/store"
)

func main() {
	urls := []string{
		"https://www.google.com",
		"https://github.com",
		"https://httpstat.us/500",
	}

	if len(os.Args) > 1 {
		urls = os.Args[1:]
	}

	// db, err := store.New("localhost", 5432, "drew", "password123", "uptime_monitor")
	// if err != nil {
	// 	log.Fatalf("Failed to connect to database: %v", err)
	// }

	fmt.Println("=== Sentinel Health Check ===")
	fmt.Println()

	// Run checks concurrently
	var wg sync.WaitGroup
	results := make([]monitor.HealthCheck, len(urls))

	for i, url := range urls {
		wg.Add(1)
		go func(idx int, u string) {
			defer wg.Done()
			results[idx] = monitor.Check(u)
		}(i, url)
	}

	wg.Wait()

	for _, r := range results {
		fmt.Println(r)
	}
}
