package main

import (
	"encoding/xml"
	"io"
	"net/http"
)

type RSSFeed struct {
	Channel FeedChannel `xml:"channel"`
}

type FeedChannel struct {
	Title         string    `xml:"title"`
	Description   string    `xml:"description"`
	LastBuildDate string    `xml:"lastBuildDate"`
	Items         []RssItem `xml:"items"`
}

type RssItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
	Guid        string `xml:"guid"`
	Description string `xml:"description"`
}

func FetchRssFeed(client *http.Client, url string, rssFeed *RSSFeed) error {
	res, err := client.Get(url)

	if err != nil {
		return err
	}
	defer res.Body.Close()

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = xml.Unmarshal(bytes, rssFeed)
	if err != nil {
		return err
	}

	return nil
}
