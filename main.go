package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	// Parse command-line flags
	includeDrafts := flag.Bool("include-drafts", false, "Include draft posts in generation")
	flag.Parse()

	// 1. Parse content directory
	posts, err := loadPosts("content/blog")
	if err != nil {
		log.Fatalf("Error loading posts: %v", err)
	}

	// 2. Filter out drafts (unless includeDrafts is true)
	publishedPosts := filterPublished(posts, *includeDrafts)

	// 3. Sort posts by date (newest first)
	sortPostsByDate(publishedPosts)

	// 4. Load static pages
	aboutPage, err := loadPage("content/about.md")
	if err != nil {
		log.Fatalf("Error loading about page: %v", err)
	}

	// 5. Load books from markdown files
	books, err := loadBooks("content/books")
	if err != nil {
		log.Fatalf("Error loading books: %v", err)
	}

	// 6. Sort books by date_read (most recent first)
	sortBooksByDate(books)

	// 7. Load RSS feed items
	feedItems, err := loadFeedItems("rss-feeds.conf")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Error loading RSS feeds: %v\n", err)
		feedItems = nil
	}

	// 8. Build site structure
	site := Site{
		Posts:     publishedPosts,
		AboutPage: aboutPage,
		Books:     books,
		FeedItems: feedItems,
	}

	// 9. Create output directory
	if err := os.RemoveAll("public"); err != nil && !os.IsNotExist(err) {
		log.Fatalf("Error removing public dir: %v", err)
	}
	if err := os.MkdirAll("public", 0755); err != nil {
		log.Fatalf("Error creating public dir: %v", err)
	}

	// 10. Copy static assets
	if err := copyStaticAssets(); err != nil {
		log.Fatalf("Error copying static assets: %v", err)
	}

	// 11. Generate all HTML files
	if err := generateHomepage(site); err != nil {
		log.Fatalf("Error generating homepage: %v", err)
	}

	if err := generatePostsListing(site); err != nil {
		log.Fatalf("Error generating posts listing: %v", err)
	}

	if err := generateBooksListing(site); err != nil {
		log.Fatalf("Error generating books listing: %v", err)
	}

	if err := generateIndividualPosts(site); err != nil {
		log.Fatalf("Error generating individual posts: %v", err)
	}

	if err := generateIndividualBooks(site); err != nil {
		log.Fatalf("Error generating individual books: %v", err)
	}

	// 12. Generate RSS reading list page
	if err := generateRSSListing(site); err != nil {
		log.Fatalf("Error generating RSS listing: %v", err)
	}

	// 13. Generate RSS feed
	if err := generateRSSFeed(site); err != nil {
		log.Fatalf("Error generating RSS feed: %v", err)
	}

	fmt.Printf("✓ Site generated successfully in public/\n")
	fmt.Printf("✓ Generated %d posts\n", len(publishedPosts))
	fmt.Printf("✓ Generated %d books\n", len(books))
	if len(feedItems) > 0 {
		fmt.Printf("✓ Generated reading list with %d articles at /reading/\n", len(feedItems))
	}
	fmt.Printf("✓ Generated RSS feed at /feed.xml\n")
}
