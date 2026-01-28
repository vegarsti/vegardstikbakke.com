package main

import "time"

// Frontmatter represents the YAML metadata at the top of markdown files
type Frontmatter struct {
	Title       string `yaml:"title"`
	Slug        string `yaml:"slug"`
	Date        string `yaml:"date"`
	Draft       bool   `yaml:"draft"`
	Layout      string `yaml:"layout"`
	Description string `yaml:"description"`
	Image       string `yaml:"image"`
}

// BookFrontmatter represents the YAML metadata for book files
type BookFrontmatter struct {
	Title         string `yaml:"title"`
	Author        string `yaml:"author"`
	YearPublished string `yaml:"year_published"`
	DateRead      string `yaml:"date_read"`
	Rating        int    `yaml:"rating"`
}

// Post represents a blog post with its content and metadata
type Post struct {
	Title       string
	Slug        string
	Date        time.Time
	DateString  string
	Draft       bool
	Description string
	Image       string
	HTMLContent string
	RawContent  string
}

// Page represents a static page (like About)
type Page struct {
	Title       string
	HTMLContent string
}

// Book represents a book with metadata
type Book struct {
	Title         string
	Slug          string
	Author        string
	YearPublished string
	DateRead      string
	Rating        int
	Summary       string
}

// Site represents the entire site structure
type Site struct {
	Posts     []Post
	AboutPage Page
	Books     []Book
	FeedItems []FeedItem
}
