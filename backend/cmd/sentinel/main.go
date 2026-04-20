package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"tracker/internal/monitor"
	"tracker/internal/store"
)

func main() {
	urls := []string{
		"https://www.netflix.com",
		"https://www.disneyplus.com",
		"https://www.twitch.tv",
		"https://www.max.com",
		"https://www.soundcloud.com",
		"https://www.bbc.com",
		"https://www.pinterest.com",
		"https://www.spotify.com",
		"https://www.vimeo.com",
		"https://www.facebook.com",
		"https://www.instagram.com",
		"https://www.x.com",
		"https://www.reddit.com",
		"https://www.slack.com",
		"https://www.discord.com",
		"https://www.linkedin.com",
		"https://www.zoom.us",
		"https://www.tumblr.com",
		"https://www.meetup.com",
		"https://www.amazon.com",
		"https://www.ebay.com",
		"https://www.robinhood.com",
		"https://www.paypal.com",
		"https://www.stripe.com",
		"https://www.nasdaq.com",
		"https://www.zalando.com",
		"https://www.instacart.com",
		"https://www.airbnb.com",
		"https://www.lyft.com",
		"https://www.delta.com",
		"https://www.toyota.com",
		"https://www.siemens.com",
		"https://www.sysco.com",
		"https://www.nytimes.com",
		"https://www.theguardian.com",
		"https://www.forbes.com",
		"https://www.coursera.org",
		"https://www.duolingo.com",
		"https://www.blackboard.com",
		"https://www.medium.com",
		"https://www.canva.com",
		"https://www.doordash.com",
		"https://www.expedia.com",
		"https://www.etsy.com",
		"https://www.coinbase.com",
		"https://www.imdb.com",
	}

	if len(os.Args) > 1 {
		urls = os.Args[1:]
	}

	db, err := store.New("localhost", 5432, "drew", "password123", "uptime_monitor")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()
	defer db.Close()

	for {
		fmt.Println("=== Sentinel Health Check ===")
		fmt.Println()

		// Run checks concurrently (rate-limited to 10 at a time)
		results := monitor.CheckAll(urls)

		for _, r := range results {
			fmt.Println(r)
		}

		if err := db.SaveChecks(context.Background(), results); err != nil {
			log.Printf("Failed to save results: %v", err)
		} else {
			fmt.Printf("\n✓ Saved %d check(s) to database\n", len(results))
		}
		<-ticker.C
	}

}
