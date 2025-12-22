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

	// 5. Create placeholder book data
	books := []Book{
		{
			Title:    "The Power Broker",
			Author:   "Robert A. Caro",
			DateRead: "2023-06",
			Rating:   5,
			Summary:  "An epic biography of Robert Moses and the fall of New York. A masterpiece of investigative journalism and narrative non-fiction.",
		},
		{
			Title:    "UNIX: A History and a Memoir",
			Author:   "Brian Kernighan",
			DateRead: "2023-08",
			Rating:   5,
			Summary:  "A personal history of Unix from one of its creators. Fascinating insights into the development of computing at Bell Labs.",
		},
		{
			Title:    "The Mythical Man-Month",
			Author:   "Frederick P. Brooks Jr.",
			DateRead: "2023-10",
			Rating:   4,
			Summary:  "Classic book on software engineering and project management. Still relevant decades after publication.",
		},
		{
			Title:    "Dealers of Lightning",
			Author:   "Michael A. Hiltzik",
			DateRead: "2024-02",
			Rating:   5,
			Summary:  "The story of Xerox PARC and how a research lab shaped the future of computing. Incredible innovation and missed opportunities.",
		},
	}

	// 6. Build site structure
	site := Site{
		Posts:     publishedPosts,
		AboutPage: aboutPage,
		Books:     books,
	}

	// 7. Create output directory
	if err := os.RemoveAll("public"); err != nil && !os.IsNotExist(err) {
		log.Fatalf("Error removing public dir: %v", err)
	}
	if err := os.MkdirAll("public", 0755); err != nil {
		log.Fatalf("Error creating public dir: %v", err)
	}

	// 8. Copy static assets
	if err := copyStaticAssets(); err != nil {
		log.Fatalf("Error copying static assets: %v", err)
	}

	// 9. Generate all HTML files
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

	// 10. Generate RSS feed
	if err := generateRSSFeed(site); err != nil {
		log.Fatalf("Error generating RSS feed: %v", err)
	}

	fmt.Printf("✓ Site generated successfully in public/\n")
	fmt.Printf("✓ Generated %d posts\n", len(publishedPosts))
	fmt.Printf("✓ Generated %d books\n", len(books))
	fmt.Printf("✓ Generated RSS feed at /feed.xml\n")
}
