package main

import "time"

// Frontmatter represents the YAML metadata at the top of markdown files
type Frontmatter struct {
	Title           string `yaml:"title"`
	Slug            string `yaml:"slug"`
	Date            string `yaml:"date"`
	Draft           bool   `yaml:"draft"`
	Starred         bool   `yaml:"starred"`
	Layout          string `yaml:"layout"`
	Description     string `yaml:"description"`
	Image           string `yaml:"image"`
	CollapsibleCode *bool  `yaml:"collapsible_code"`
}

// Post represents a blog post with its content and metadata
type Post struct {
	Title           string
	Slug            string
	Date            time.Time
	DateString      string
	Draft           bool
	Starred         bool
	Description     string
	Image           string
	CollapsibleCode bool
	HTMLContent     string
	RawContent      string
}

// Page represents a static page (like About)
type Page struct {
	Title       string
	HTMLContent string
}

// Site represents the entire site structure
type Site struct {
	Posts     []Post // Published posts only (for listings, RSS, etc.)
	AllPosts  []Post // All posts including drafts (for individual page generation)
	AboutPage Page
	FeedItems []FeedItem
}
