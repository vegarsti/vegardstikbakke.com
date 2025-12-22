package main

import "time"

// Frontmatter represents the YAML metadata at the top of markdown files
type Frontmatter struct {
	Title  string `yaml:"title"`
	Slug   string `yaml:"slug"`
	Date   string `yaml:"date"`
	Draft  bool   `yaml:"draft"`
	Layout string `yaml:"layout"`
}

// Post represents a blog post with its content and metadata
type Post struct {
	Title       string
	Slug        string
	Date        time.Time
	DateString  string
	Draft       bool
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
	Title    string
	Author   string
	DateRead string
	Rating   int
	Summary  string
}

// Site represents the entire site structure
type Site struct {
	Posts     []Post
	AboutPage Page
	Books     []Book
}
