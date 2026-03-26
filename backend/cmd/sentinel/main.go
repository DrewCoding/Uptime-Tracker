package main

import (
	"fmt"
	"os"
	"sync"

	"tracker/internal/monitor"
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

	fmt.Println("=== Sentinel Health Check ===")
	fmt.Println()

	var wg sync.WaitGroup

	for _, url := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			result := monitor.Check(u)
			fmt.Println(result)
		}(url)
	}

	wg.Wait()
}
