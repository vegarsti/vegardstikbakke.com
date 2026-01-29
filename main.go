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
	skipRSSIfExists := flag.Bool("skip-rss-if-exists", false, "Skip RSS fetching if reading page already exists")
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

	// 7. Check if we should skip RSS fetching (preserve existing reading page)
	skipRSS := *skipRSSIfExists
	if _, err := os.Stat("public/reading/index.html"); os.IsNotExist(err) {
		skipRSS = false // Must fetch if reading page doesn't exist
	}

	// 8. Load RSS feed items (unless skipping)
	var feedItems []FeedItem
	if skipRSS {
		fmt.Println("Skipping RSS fetch (reading page already exists)")
	} else {
		feedItems, err = loadFeedItems("rss-feeds.conf")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Error loading RSS feeds: %v\n", err)
			feedItems = nil
		}
	}

	// 9. Build site structure
	site := Site{
		Posts:     publishedPosts,
		AboutPage: aboutPage,
		Books:     books,
		FeedItems: feedItems,
	}

	// 10. Create output directory (preserve reading/ if skipping RSS)
	if skipRSS {
		// Preserve reading directory by removing everything else
		entries, err := os.ReadDir("public")
		if err != nil && !os.IsNotExist(err) {
			log.Fatalf("Error reading public dir: %v", err)
		}
		for _, entry := range entries {
			if entry.Name() != "reading" {
				if err := os.RemoveAll("public/" + entry.Name()); err != nil {
					log.Fatalf("Error removing %s: %v", entry.Name(), err)
				}
			}
		}
	} else {
		if err := os.RemoveAll("public"); err != nil && !os.IsNotExist(err) {
			log.Fatalf("Error removing public dir: %v", err)
		}
	}
	if err := os.MkdirAll("public", 0755); err != nil {
		log.Fatalf("Error creating public dir: %v", err)
	}

	// 11. Copy static assets
	if err := copyStaticAssets(); err != nil {
		log.Fatalf("Error copying static assets: %v", err)
	}

	// 12. Generate all HTML files
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

	// 13. Generate RSS reading list page (skip if preserving existing)
	if !skipRSS {
		if err := generateRSSListing(site); err != nil {
			log.Fatalf("Error generating RSS listing: %v", err)
		}
	}

	// 14. Generate RSS feed
	if err := generateRSSFeed(site); err != nil {
		log.Fatalf("Error generating RSS feed: %v", err)
	}

	// 15. Generate 404 page
	if err := generate404Page(); err != nil {
		log.Fatalf("Error generating 404 page: %v", err)
	}

	fmt.Printf("✓ Site generated successfully in public/\n")
	fmt.Printf("✓ Generated %d posts\n", len(publishedPosts))
	fmt.Printf("✓ Generated %d books\n", len(books))
	if len(feedItems) > 0 {
		fmt.Printf("✓ Generated reading list with %d articles at /reading/\n", len(feedItems))
	}
	fmt.Printf("✓ Generated RSS feed at /feed.xml\n")
}
