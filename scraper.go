package main

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
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

		description := sql.NullString{
			String: item.Description,
			Valid:  item.Description != "",
		}

		pubAt, err := time.Parse(time.RFC1123, item.PubDate)
		if err != nil {
			fmt.Println("Error parsing time", err, "for", item.PubDate)
			continue
		}

		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Description: description,
			Url:         item.Link,
			PublishedAt: pubAt,
			FeedID:      feed.ID,
		})
		if err != nil {
			// 같은 아티클이 이미 존재하는 경우 URL 이 동일해서 막히니깐 굳이 로그에 안찍는다.
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			fmt.Println("Error creating post", err)
		}
	}

	fmt.Println("Total", len(rssFeed.Channel.Item), "collected.")

}
