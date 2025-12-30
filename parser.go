package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	"gopkg.in/yaml.v3"
)

// loadPosts reads all markdown files from a directory and parses them
func loadPosts(dir string) ([]Post, error) {
	files, err := filepath.Glob(filepath.Join(dir, "*.md"))
	if err != nil {
		return nil, err
	}

	var posts []Post
	for _, file := range files {
		post, err := parsePost(file)
		if err != nil {
			return nil, fmt.Errorf("error parsing %s: %w", file, err)
		}
		posts = append(posts, post)
	}

	return posts, nil
}

// parsePost reads a markdown file and extracts frontmatter and content
func parsePost(filePath string) (Post, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return Post{}, err
	}

	// Split frontmatter from content
	frontmatter, markdown := extractFrontmatter(content)

	// Parse YAML frontmatter
	var fm Frontmatter
	if err := yaml.Unmarshal([]byte(frontmatter), &fm); err != nil {
		return Post{}, fmt.Errorf("error parsing frontmatter: %w", err)
	}

	// Validate required fields
	if fm.Title == "" {
		return Post{}, fmt.Errorf("post missing required field: title")
	}

	// If slug is missing, generate from filename
	slug := fm.Slug
	if slug == "" {
		// Extract filename without extension
		filename := filepath.Base(filePath)
		slug = strings.TrimSuffix(filename, ".md")
	}

	// Convert markdown to HTML
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(markdown), &buf); err != nil {
		return Post{}, fmt.Errorf("error converting markdown: %w", err)
	}

	// Parse date flexibly
	parsedDate := parseDate(fm.Date)

	return Post{
		Title:       fm.Title,
		Slug:        slug,
		Date:        parsedDate,
		DateString:  fm.Date,
		Draft:       fm.Draft,
		Description: fm.Description,
		Image:       fm.Image,
		HTMLContent: buf.String(),
		RawContent:  markdown,
	}, nil
}

// extractFrontmatter splits a file into frontmatter and content
func extractFrontmatter(content []byte) (frontmatter, markdown string) {
	str := string(content)

	// Check if file starts with "---"
	if !strings.HasPrefix(str, "---\n") && !strings.HasPrefix(str, "---\r\n") {
		return "", str
	}

	// Skip the opening "---\n" or "---\r\n"
	rest := str
	if strings.HasPrefix(str, "---\r\n") {
		rest = str[5:]
	} else {
		rest = str[4:]
	}

	// Find the closing "---" which can be followed by newline or EOF
	// Look for newline followed by three dashes
	closingPatterns := []string{"\n---\n", "\n---\r\n", "\n---"}
	var parts []string
	var found bool

	for _, pattern := range closingPatterns {
		idx := strings.Index(rest, pattern)
		if idx >= 0 {
			frontmatterContent := rest[:idx]
			markdownContent := ""

			// Get content after the closing ---
			afterClosing := idx + len(pattern)
			if afterClosing < len(rest) {
				markdownContent = rest[afterClosing:]
			}

			return frontmatterContent, strings.TrimSpace(markdownContent)
		}
	}

	// If no closing found, treat entire file as content
	if !found {
		return "", str
	}

	return parts[0], strings.TrimSpace(parts[1])
}

// parseDate attempts to parse various date formats
func parseDate(dateStr string) time.Time {
	if dateStr == "" {
		return time.Time{}
	}

	// Try common formats
	formats := []string{
		"2006-01-02",                 // "2020-03-21"
		"2006-01-02T15:04:05Z07:00",  // ISO 8601 with timezone
		"2006-01-02T15:04:05",        // ISO 8601 without timezone
		"2006-01-02 15:04:05",        // Space separated
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t
		}
	}

	// If parsing fails, return zero time
	return time.Time{}
}

// filterPublished filters out draft posts unless includeDrafts is true
func filterPublished(posts []Post, includeDrafts bool) []Post {
	if includeDrafts {
		return posts
	}

	var published []Post
	for _, post := range posts {
		if !post.Draft {
			published = append(published, post)
		}
	}
	return published
}

// sortPostsByDate sorts posts by date, newest first
func sortPostsByDate(posts []Post) {
	sort.Slice(posts, func(i, j int) bool {
		// Posts without dates should go to the end
		if posts[i].Date.IsZero() {
			return false
		}
		if posts[j].Date.IsZero() {
			return true
		}
		return posts[i].Date.After(posts[j].Date)
	})
}

// loadPage loads a static page (like about.md)
func loadPage(filePath string) (Page, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return Page{}, err
	}

	frontmatter, markdown := extractFrontmatter(content)

	var fm Frontmatter
	if err := yaml.Unmarshal([]byte(frontmatter), &fm); err != nil {
		return Page{}, fmt.Errorf("error parsing frontmatter: %w", err)
	}

	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(markdown), &buf); err != nil {
		return Page{}, fmt.Errorf("error converting markdown: %w", err)
	}

	return Page{
		Title:       fm.Title,
		HTMLContent: buf.String(),
	}, nil
}

// loadBooks reads all markdown files from content/books/ and parses them
func loadBooks(dir string) ([]Book, error) {
	files, err := filepath.Glob(filepath.Join(dir, "*.md"))
	if err != nil {
		return nil, err
	}

	var books []Book
	for _, file := range files {
		book, err := parseBook(file)
		if err != nil {
			return nil, fmt.Errorf("error parsing %s: %w", file, err)
		}
		books = append(books, book)
	}

	return books, nil
}

// parseBook reads a markdown file and extracts book metadata and summary
func parseBook(filePath string) (Book, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return Book{}, err
	}

	// Split frontmatter from content
	frontmatter, markdown := extractFrontmatter(content)

	// Parse YAML frontmatter
	var fm BookFrontmatter
	if err := yaml.Unmarshal([]byte(frontmatter), &fm); err != nil {
		return Book{}, fmt.Errorf("error parsing frontmatter: %w", err)
	}

	// Validate required fields
	if fm.Title == "" {
		return Book{}, fmt.Errorf("book missing required field: title")
	}
	if fm.Author == "" {
		return Book{}, fmt.Errorf("book missing required field: author")
	}

	// Generate slug from filename
	filename := filepath.Base(filePath)
	slug := strings.TrimSuffix(filename, ".md")

	// Convert markdown summary to HTML
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(markdown), &buf); err != nil {
		return Book{}, fmt.Errorf("error converting markdown: %w", err)
	}

	return Book{
		Title:         fm.Title,
		Slug:          slug,
		Author:        fm.Author,
		YearPublished: fm.YearPublished,
		DateRead:      fm.DateRead,
		Rating:        fm.Rating,
		Summary:       buf.String(),
	}, nil
}

// sortBooksByDate sorts books by date_read, most recent first
func sortBooksByDate(books []Book) {
	sort.Slice(books, func(i, j int) bool {
		// Books without dates should go to the end
		if books[i].DateRead == "" {
			return false
		}
		if books[j].DateRead == "" {
			return true
		}
		// Simple string comparison works for YYYY-MM or YYYY formats
		return books[i].DateRead > books[j].DateRead
	})
}
