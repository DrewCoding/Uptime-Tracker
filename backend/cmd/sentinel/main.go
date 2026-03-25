package main

import (
	"fmt"
	"os"

	"tracker/internal/monitor"
)

func main() {
	urls := []string{
		"https://www.google.com",
		"https://github.com",
		"https://httpstat.us/500",
	}

	// Allow overriding default URLs via command-line args
	if len(os.Args) > 1 {
		urls = os.Args[1:]
	}

	fmt.Println("=== Sentinel Health Check ===")
	fmt.Println()

	for _, url := range urls {
		result := monitor.Check(url)
		fmt.Println(result)
	}
}
