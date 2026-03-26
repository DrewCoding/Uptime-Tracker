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

	for _, i := range results {
		fmt.Println(i)
	}
}
