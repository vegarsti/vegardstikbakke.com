package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

// FeedItem represents a single item from an RSS/Atom feed
type FeedItem struct {
	Title     string
	Link      string
	PubDate   time.Time
	FeedTitle string
	FeedURL   string
}

// RSS feed structures
type rssFeed struct {
	Channel rssChannel `xml:"channel"`
}

type rssChannel struct {
	Title string    `xml:"title"`
	Items []rssItem `xml:"item"`
}

type rssItem struct {
	Title   string `xml:"title"`
	Link    string `xml:"link"`
	PubDate string `xml:"pubDate"`
}

// Atom feed structures
type atomFeed struct {
	Title   string      `xml:"title"`
	Entries []atomEntry `xml:"entry"`
}

type atomEntry struct {
	Title     string   `xml:"title"`
	Link      atomLink `xml:"link"`
	Published string   `xml:"published"`
	Updated   string   `xml:"updated"`
}

type atomLink struct {
	Href string `xml:"href,attr"`
}

// FeedSource identifies an RSS/Atom feed and an optional display-name override.
type FeedSource struct {
	URL         string
	DisplayName string
}

// loadFeedSources reads feed URLs and optional display-name overrides from a config file.
// Each non-comment line has the form: URL [| Display Name].
func loadFeedSources(filename string) ([]FeedSource, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var feeds []FeedSource
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Skip empty lines and comments.
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "|", 2)
		feeds = append(feeds, FeedSource{
			URL:         strings.TrimSpace(parts[0]),
			DisplayName: "",
		})
		if len(parts) == 2 {
			feeds[len(feeds)-1].DisplayName = strings.TrimSpace(parts[1])
		}
	}
	return feeds, scanner.Err()
}

// fetchAllFeeds fetches all feeds concurrently and returns combined items
func fetchAllFeeds(feeds []FeedSource) ([]FeedItem, error) {
	var allItems []FeedItem
	var mu sync.Mutex
	var wg sync.WaitGroup
	var errors []error
	var errMu sync.Mutex

	for _, feed := range feeds {
		wg.Add(1)
		go func(feed FeedSource) {
			defer wg.Done()
			items, err := fetchFeed(feed.URL)
			if err != nil {
				errMu.Lock()
				errors = append(errors, fmt.Errorf("fetching %s: %w", feed.URL, err))
				errMu.Unlock()
				return
			}
			if feed.DisplayName != "" {
				for i := range items {
					items[i].FeedTitle = feed.DisplayName
				}
			}
			mu.Lock()
			allItems = append(allItems, items...)
			mu.Unlock()
		}(feed)
	}

	wg.Wait()

	// Log errors but don't fail completely
	for _, err := range errors {
		fmt.Fprintf(os.Stderr, "Warning: %v\n", err)
	}

	// Sort by date, most recent first
	sort.Slice(allItems, func(i, j int) bool {
		return allItems[i].PubDate.After(allItems[j].PubDate)
	})

	return allItems, nil
}

// fetchFeed fetches and parses a single feed
func fetchFeed(url string) ([]FeedItem, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Try parsing as RSS first
	var rss rssFeed
	if err := xml.Unmarshal(body, &rss); err == nil && len(rss.Channel.Items) > 0 {
		return parseRSSItems(rss, url), nil
	}

	// Try parsing as Atom
	var atom atomFeed
	if err := xml.Unmarshal(body, &atom); err == nil && len(atom.Entries) > 0 {
		return parseAtomItems(atom, url), nil
	}

	return nil, fmt.Errorf("unable to parse feed format")
}

func parseRSSItems(rss rssFeed, feedURL string) []FeedItem {
	items := make([]FeedItem, 0, len(rss.Channel.Items))
	for _, item := range rss.Channel.Items {
		items = append(items, FeedItem{
			Title:     item.Title,
			Link:      item.Link,
			PubDate:   parseFeedDate(item.PubDate),
			FeedTitle: rss.Channel.Title,
			FeedURL:   feedURL,
		})
	}
	return items
}

func parseAtomItems(atom atomFeed, feedURL string) []FeedItem {
	items := make([]FeedItem, 0, len(atom.Entries))
	for _, entry := range atom.Entries {
		dateStr := entry.Published
		if dateStr == "" {
			dateStr = entry.Updated
		}
		items = append(items, FeedItem{
			Title:     entry.Title,
			Link:      entry.Link.Href,
			PubDate:   parseFeedDate(dateStr),
			FeedTitle: atom.Title,
			FeedURL:   feedURL,
		})
	}
	return items
}

func parseFeedDate(dateStr string) time.Time {
	dateStr = strings.TrimSpace(dateStr)
	if dateStr == "" {
		return time.Time{}
	}

	// Common date formats used in RSS/Atom feeds
	formats := []string{
		time.RFC1123Z,
		time.RFC1123,
		time.RFC3339,
		time.RFC3339Nano,
		"Mon, 2 Jan 2006 15:04:05 -0700",
		"Mon, 2 Jan 2006 15:04:05 MST",
		"2006-01-02T15:04:05-07:00",
		"2006-01-02T15:04:05Z",
		"2006-01-02 15:04:05",
		"2006-01-02",
		"02 Jan 2006 15:04:05 -0700",
		"02 Jan 2006 15:04:05 MST",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t
		}
	}

	return time.Time{}
}

// loadFeedItems loads all feed items from the configured feeds
func loadFeedItems(configFile string) ([]FeedItem, error) {
	feeds, err := loadFeedSources(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // No feeds config, skip RSS aggregation
		}
		return nil, err
	}

	if len(feeds) == 0 {
		return nil, nil
	}

	return fetchAllFeeds(feeds)
}
