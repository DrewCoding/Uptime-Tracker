package monitor

import (
	"fmt"
	"net/http"
	"time"
)

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

	resp, err := http.Get(url)
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
