package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
)

// Book is a single entry from data/books.ndjson, fetched from the Goodreads
// read-shelf RSS feed by scripts/fetch_goodreads.py. The site reads only this
// file (no network at build time), so a Goodreads outage never breaks a build.
// Only the fields the page renders plus cover URLs are kept; the fetcher
// trims everything else (descriptions, raw dates, ids) at write time.
type Book struct {
	Title           string `json:"title"`
	AuthorName      string `json:"author_name"`
	UserRating      string `json:"user_rating"`
	ReadAtISO       string `json:"read_at_iso"`
	BookPublished   string `json:"book_published"`
	Link            string `json:"link"`
	BookImageURL    string `json:"book_image_url"`
	BookMediumImage string `json:"book_medium_image_url"`
	BookLargeImage  string `json:"book_large_image_url"`
}

// loadBooks reads data/books.ndjson (one JSON object per line) and returns the
// books sorted newest-read first. A missing file is not an error: it simply
// means the books page won't be generated.
func loadBooks(path string) ([]Book, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()

	var books []Book
	scanner := bufio.NewScanner(f)
	// Book descriptions can be long; raise the per-line buffer to 1 MiB.
	scanner.Buffer(make([]byte, 0, 64*1024), 1<<20)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var b Book
		if err := json.Unmarshal([]byte(line), &b); err != nil {
			return nil, fmt.Errorf("books.ndjson: malformed line: %w", err)
		}
		books = append(books, b)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Sort newest-read first; books with no read_at sink to the bottom.
	sort.SliceStable(books, func(i, j int) bool {
		return books[i].ReadAtISO > books[j].ReadAtISO
	})

	return books, nil
}

// stars renders a "5" rating (or "0"/empty) as filled/empty star glyphs.
// "0" (unrated) renders as empty stars only.
func (b Book) Stars() string {
	r := strings.TrimSpace(b.UserRating)
	var n int
	if r != "" {
		fmt.Sscanf(r, "%d", &n)
	}
	if n < 0 {
		n = 0
	}
	if n > 5 {
		n = 5
	}
	return strings.Repeat("★", n) + strings.Repeat("☆", 5-n)
}
