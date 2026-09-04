package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestFetchAllFeedsAppliesDisplayNameOverride(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/rss+xml")
		_, _ = w.Write([]byte(`<?xml version="1.0"?><rss><channel><title></title><item><title>Bug blindness</title><link>https://danluu.com/bug-blindness/</link><pubDate>Mon, 2 Jan 2006 15:04:05 -0700</pubDate></item></channel></rss>`))
	}))
	defer server.Close()

	items, err := fetchAllFeeds([]FeedSource{{URL: server.URL, DisplayName: "Dan Luu"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 {
		t.Fatalf("got %d items, want 1", len(items))
	}
	if items[0].FeedTitle != "Dan Luu" {
		t.Errorf("FeedTitle = %q, want %q", items[0].FeedTitle, "Dan Luu")
	}
}

func TestLoadFeedSources(t *testing.T) {
	config := "# A comment\n\nhttps://example.com/feed.xml\nhttps://danluu.com/atom.xml | Dan Luu\n"
	path := filepath.Join(t.TempDir(), "feeds.conf")
	if err := os.WriteFile(path, []byte(config), 0644); err != nil {
		t.Fatal(err)
	}

	feeds, err := loadFeedSources(path)
	if err != nil {
		t.Fatal(err)
	}

	want := []FeedSource{
		{URL: "https://example.com/feed.xml"},
		{URL: "https://danluu.com/atom.xml", DisplayName: "Dan Luu"},
	}
	if len(feeds) != len(want) {
		t.Fatalf("got %d feeds, want %d", len(feeds), len(want))
	}
	for i := range want {
		if feeds[i] != want[i] {
			t.Errorf("feed %d = %#v, want %#v", i, feeds[i], want[i])
		}
	}
}
