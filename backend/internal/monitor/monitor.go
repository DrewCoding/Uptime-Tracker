package monitor

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"sync"
	"time"
)

const maxConcurrent = 10

// - Use a 10-second timeout to prevent hanging on unresponsive hosts
var client = &http.Client{
	Timeout: 10 * time.Second,
	Transport: &http.Transport{
		// Force HTTP/1.1 — setting TLSNextProto to an empty map disables
		// the automatic HTTP/2 upgrade that Go's default transport performs.
		TLSNextProto: make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),

		MaxIdleConns:        maxConcurrent,
		MaxIdleConnsPerHost: 2,
		IdleConnTimeout:     30 * time.Second,
	},
}

type HealthCheck struct {
	URL        string
	StatusCode int
	LatencyMs  int64
	Err        error
	CheckedAt  time.Time
}

// Check performs an HTTP GET to the given URL and measures the round-trip latency.
func Check(url string) HealthCheck {
	start := time.Now()

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return HealthCheck{
			URL:       url,
			LatencyMs: 0,
			Err:       err,
			CheckedAt: start,
		}
	}

	// Set a realistic User-Agent so firewalls/WAFs don't reject the request.
	req.Header.Set("User-Agent", "Sentinel-Monitor/1.0")

	resp, err := client.Do(req)
	latency := time.Since(start).Milliseconds()

	if err != nil {
		return HealthCheck{
			URL:       url,
			LatencyMs: latency,
			Err:       err,
			CheckedAt: start,
		}
	}
	defer resp.Body.Close()

	return HealthCheck{
		URL:        url,
		StatusCode: resp.StatusCode,
		LatencyMs:  latency,
		CheckedAt:  start,
	}
}

// CheckAll runs health checks on all provided URLs concurrently,
// limiting parallelism to maxConcurrent in-flight requests.
func CheckAll(urls []string) []HealthCheck {
	results := make([]HealthCheck, len(urls))
	sem := make(chan struct{}, maxConcurrent)

	var wg sync.WaitGroup
	for i, url := range urls {
		wg.Add(1)
		go func(idx int, u string) {
			defer wg.Done()
			sem <- struct{}{} // acquire slot
			results[idx] = Check(u)
			<-sem // release slot
		}(i, url)
	}
	wg.Wait()
	return results
}

func (r HealthCheck) String() string {
	if r.Err != nil {
		return fmt.Sprintf("[OFFLINE]  %s | error: %v | %dms", r.URL, r.Err, r.LatencyMs)
	}
	status := "ONLINE"
	if r.StatusCode != http.StatusOK {
		status = "ERROR"
	}
	return fmt.Sprintf("[%s]  %s | %d | %dms", status, r.URL, r.StatusCode, r.LatencyMs)
}
