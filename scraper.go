package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/intaek-h/rss/internal/database"
)

func startScraping(
	db *database.Queries,
	concurrency int,
	timeBetweenRequest time.Duration,
) {
	fmt.Println("Starting scraping on", concurrency, "goroutines every", timeBetweenRequest, "duration")

	ticker := time.NewTicker(timeBetweenRequest)

	// for range ticker.C {} 는 for 문이 시작하고 1분 뒤에 로직이 실행된다.
	// for ; ; <-ticker.C {} 는 for 문이 시작하자마자 로직이 실행된다.
	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedToFetch(
			context.Background(),
			int32(concurrency),
		)
		if err != nil {
			fmt.Println("Error fetching feeds", err)
			continue
		}

		wg := sync.WaitGroup{}

		for _, feed := range feeds {
			wg.Add(1)
			go scrape(db, &wg, feed)
		}

		wg.Wait()
	}
}

func scrape(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()

	_, err := db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		fmt.Println("Error marking feed fetched", err)
		return
	}

	rssFeed, err := fetchFeedFromURL(feed.Url)
	if err != nil {
		fmt.Println("Error fetching feed", err)
		return
	}

	for _, item := range rssFeed.Channel.Item {
		fmt.Println("Found post", item.Title, "on", feed.Name)
	}

	fmt.Println("Total", len(rssFeed.Channel.Item), "collected.")

}
