package gatorfeed

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
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

func FetchFeed(ctx context.Context, feedUrl string) (*RSSFeed, error) {
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", feedUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating a get request with context: %w", err)
	}
	req.Header.Set("User-Agent", "gator")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending get request to '%v': %w")
	}
	defer resp.Body.Close()
	var feed RSSFeed
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading all bytes of response body: %w", err)
	}
	err = xml.Unmarshal(bytes, &feed)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling data from bytes: %w", err)
	}
	CleanFeed(&feed)
	return &feed, nil
}

func CleanFeed(feed *RSSFeed) {
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for _, rssitem := range feed.Channel.Item {
		rssitem.Title = html.UnescapeString(rssitem.Title)
		rssitem.Description = html.UnescapeString(rssitem.Description)
	}
}
