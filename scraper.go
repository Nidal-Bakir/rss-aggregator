package main

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/Nidal-Bakir/rss-aggregator/internal/database"
	"github.com/google/uuid"
)

func startScraper(DB *database.Queries, concurrentRequestsCount int32, scraperIntervale time.Duration) {

	log.Printf("Start scraping with %d concurrent requests and interval: %s\n",
		concurrentRequestsCount,
		scraperIntervale.String(),
	)

	client := &http.Client{Timeout: time.Second * 10}

	ticker := time.NewTicker(scraperIntervale)

	for ; ; <-ticker.C {

		feeds, err := DB.GetFeedsOrderedByLastSync(
			context.Background(),
			database.GetFeedsOrderedByLastSyncParams{
				Limit: concurrentRequestsCount,
			},
		)

		if err != nil || len(feeds) == 0 {
			if err != nil {
				log.Println("Error while getting the feeds from the database for scraping:", err)
			} else {
				log.Println("No Feeds to scrap, the database returned an empty slice")
			}
			continue
		}

		scrapeFeeds(DB, &feeds, client)

	}
}

func scrapeFeeds(DB *database.Queries, feeds *[]database.GetFeedsOrderedByLastSyncRow, client *http.Client) {
	wg := &sync.WaitGroup{}
	defer log.Printf("waiting for all feeds to be done count:%d\n", len(*feeds))
	defer wg.Wait()

	for _, feed := range *feeds {
		go func(url string, id uuid.UUID) {
			defer wg.Done()
			wg.Add(1)

			var rssFeed RSSFeed
			if err := scrapeFeed(url, &rssFeed, client); err != nil {
				return
			}
			updateRssFeedDataInDB(id, DB, &rssFeed)

		}(feed.Url, feed.ID)
	}

}

func scrapeFeed(url string, rssFeed *RSSFeed, client *http.Client) error {
	log.Println("Scaping feed:", url)

	err := FetchRssFeed(client, url, rssFeed)

	if err != nil {
		log.Println("Error while fetching the feed from remote server", err)
		return err
	}

	return nil
}

func updateRssFeedDataInDB(id uuid.UUID, DB *database.Queries, rssFeed *RSSFeed) {
	if err := DB.MakeFeedAsSynced(context.Background(), id); err != nil {
		log.Println("Error while marking the feed as synced", err)
		return
	}

	log.Printf("The Feed %s, collected and we found %d posts(items)",
		rssFeed.Channel.Title, len(rssFeed.Channel.Items))
}
