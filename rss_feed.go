package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	feed := RSSFeed{}
	client := http.Client{Timeout: time.Minute}

	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &feed, err
	}

	req.Header.Add("User-Agent", "gator")
	res, err := client.Do(req)
	if err != nil {
		return &feed, err
	}
	defer res.Body.Close()

	if res.StatusCode > 399 {
		if res.StatusCode == 404 {
			return &feed, fmt.Errorf("invalid invalid feed: %s", feedURL)
		}
		return &feed, fmt.Errorf("bad status code: %v", res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return &feed, err
	}

	err = xml.Unmarshal(data, &feed)
	if err != nil {
		return &feed, err
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	for i := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(feed.Channel.Item[i].Title)
		feed.Channel.Item[i].Description = html.UnescapeString(feed.Channel.Item[i].Description)
	}

	return &feed, nil
}
