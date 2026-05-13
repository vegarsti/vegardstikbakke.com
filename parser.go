package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"gopkg.in/yaml.v3"
)

var headingRe = regexp.MustCompile(`<h([1-6]) id="([^"]+)">`)

// addHeadingAnchors post-processes HTML to inject anchor links into headings
func addHeadingAnchors(htmlContent string) string {
	return headingRe.ReplaceAllStringFunc(htmlContent, func(match string) string {
		groups := headingRe.FindStringSubmatch(match)
		level, id := groups[1], groups[2]
		return fmt.Sprintf(`<h%s id="%s"><a class="heading-anchor" href="#%s" aria-label="Link to this section">#</a>`, level, id, id)
	})
}

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
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.Table,
			extension.Footnote,
			highlighting.NewHighlighting(
				highlighting.WithStyle("github"),
				highlighting.WithFormatOptions(
					chromahtml.WithClasses(true),
				),
			),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)
	var buf bytes.Buffer
	if err := md.Convert([]byte(markdown), &buf); err != nil {
		return Post{}, fmt.Errorf("error converting markdown: %w", err)
	}

	// Parse date flexibly
	parsedDate := parseDate(fm.Date)

	// Default collapsible_code to true
	collapsibleCode := true
	if fm.CollapsibleCode != nil {
		collapsibleCode = *fm.CollapsibleCode
	}

	return Post{
		Title:           fm.Title,
		Slug:            slug,
		Date:            parsedDate,
		DateString:      fm.Date,
		Draft:           fm.Draft,
		Starred:         fm.Starred,
		Description:     fm.Description,
		Image:           fm.Image,
		CollapsibleCode: collapsibleCode,
		HTMLContent:     addHeadingAnchors(buf.String()),
		RawContent:      markdown,
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
		"2006-01-02",                // "2020-03-21"
		"2006-01-02T15:04:05Z07:00", // ISO 8601 with timezone
		"2006-01-02T15:04:05",       // ISO 8601 without timezone
		"2006-01-02 15:04:05",       // Space separated
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

	md := goldmark.New(
		goldmark.WithExtensions(
			extension.Table,
			extension.Footnote,
			highlighting.NewHighlighting(
				highlighting.WithStyle("github"),
				highlighting.WithFormatOptions(
					chromahtml.WithClasses(true),
				),
			),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)
	var buf bytes.Buffer
	if err := md.Convert([]byte(markdown), &buf); err != nil {
		return Page{}, fmt.Errorf("error converting markdown: %w", err)
	}

	return Page{
		Title:       fm.Title,
		HTMLContent: addHeadingAnchors(buf.String()),
	}, nil
}
