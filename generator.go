package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"html"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// generateHomepage creates the homepage with bio from about.md
func generateHomepage(site Site) error {
	tmpl := getBaseTemplate()
	template.Must(tmpl.Parse(homepageContent))

	// Replace {{MOST_RECENT_POST}} placeholder with actual post link
	content := site.AboutPage.HTMLContent
	if len(site.Posts) > 0 {
		mostRecent := site.Posts[0]
		postLink := fmt.Sprintf(`<a href="/%s/">%s</a>`, mostRecent.Slug, mostRecent.Title)
		content = strings.Replace(content, "{{MOST_RECENT_POST}}", postLink, 1)
	}

	data := struct {
		Title        string
		Description  string
		CanonicalURL string
		Image        string
		OGType       string
		Content      template.HTML
	}{
		Title:        "Vegard Stikbakke",
		Description:  "Vegard Stikbakke — software engineer from Norway",
		CanonicalURL: "https://vegardstikbakke.com/",
		Image:        "/me.jpg",
		OGType:       "website",
		Content:      template.HTML(content),
	}

	return renderToFile(tmpl, data, "public/index.html")
}

// generatePostsListing creates the posts listing page
func generatePostsListing(site Site) error {
	tmpl := getBaseTemplate()
	template.Must(tmpl.Parse(postsListingContent))

	data := struct {
		Title        string
		Description  string
		CanonicalURL string
		Image        string
		OGType       string
		Posts        []Post
	}{
		Title:        "Posts — Vegard Stikbakke",
		Description:  "Blog posts by Vegard Stikbakke",
		CanonicalURL: "https://vegardstikbakke.com/blog/",
		Image:        "",
		OGType:       "website",
		Posts:        site.Posts,
	}

	if err := os.MkdirAll("public/blog", 0755); err != nil {
		return err
	}

	return renderToFile(tmpl, data, "public/blog/index.html")
}

// generateBooksListing creates the books listing page
func generateBooksListing(site Site) error {
	tmpl := getBaseTemplate()
	template.Must(tmpl.Parse(booksListingContent))

	// Convert books to template-friendly format with HTML summaries
	type BookDisplay struct {
		Title         string
		Slug          string
		Author        string
		YearPublished string
		DateRead      string
		Rating        int
		Summary       template.HTML
	}

	books := make([]BookDisplay, len(site.Books))
	for i, book := range site.Books {
		books[i] = BookDisplay{
			Title:         book.Title,
			Slug:          book.Slug,
			Author:        book.Author,
			YearPublished: book.YearPublished,
			DateRead:      book.DateRead,
			Rating:        book.Rating,
			Summary:       template.HTML(book.Summary),
		}
	}

	data := struct {
		Title        string
		Description  string
		CanonicalURL string
		Image        string
		OGType       string
		Books        []BookDisplay
	}{
		Title:        "Books — Vegard Stikbakke",
		Description:  "Books read by Vegard Stikbakke",
		CanonicalURL: "https://vegardstikbakke.com/books/",
		Image:        "",
		OGType:       "website",
		Books:        books,
	}

	if err := os.MkdirAll("public/books", 0755); err != nil {
		return err
	}

	return renderToFile(tmpl, data, "public/books/index.html")
}

// generateIndividualPosts creates individual post pages
func generateIndividualPosts(site Site) error {
	tmpl := getBaseTemplate()
	template.Must(tmpl.Parse(postContent))

	for _, post := range site.Posts {
		data := struct {
			Title        string
			PostTitle    string
			DateString   string
			Description  string
			CanonicalURL string
			Image        string
			OGType       string
			Content      template.HTML
		}{
			Title:        post.Title + " — Vegard Stikbakke",
			PostTitle:    post.Title,
			DateString:   post.DateString,
			Description:  post.Description,
			CanonicalURL: fmt.Sprintf("https://vegardstikbakke.com/%s/", post.Slug),
			Image:        post.Image,
			OGType:       "article",
			Content:      template.HTML(post.HTMLContent),
		}

		dir := filepath.Join("public", post.Slug)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}

		outputPath := filepath.Join(dir, "index.html")
		if err := renderToFile(tmpl, data, outputPath); err != nil {
			return err
		}
	}

	return nil
}

// generateIndividualBooks creates individual book pages
func generateIndividualBooks(site Site) error {
	tmpl := getBaseTemplate()
	template.Must(tmpl.Parse(bookContent))

	for _, book := range site.Books {
		data := struct {
			Title         string
			BookTitle     string
			Author        string
			YearPublished string
			DateRead      string
			Rating        int
			Description   string
			CanonicalURL  string
			Image         string
			OGType        string
			Content       template.HTML
		}{
			Title:         book.Title + " — Vegard Stikbakke",
			BookTitle:     book.Title,
			Author:        book.Author,
			YearPublished: book.YearPublished,
			DateRead:      book.DateRead,
			Rating:        book.Rating,
			Description:   fmt.Sprintf("%s by %s", book.Title, book.Author),
			CanonicalURL:  fmt.Sprintf("https://vegardstikbakke.com/books/%s/", book.Slug),
			Image:         "",
			OGType:        "article",
			Content:       template.HTML(book.Summary),
		}

		dir := filepath.Join("public", "books", book.Slug)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}

		outputPath := filepath.Join(dir, "index.html")
		if err := renderToFile(tmpl, data, outputPath); err != nil {
			return err
		}
	}

	return nil
}

// renderToFile renders a template to a file
func renderToFile(tmpl *template.Template, data interface{}, filepath string) error {
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return err
	}

	return os.WriteFile(filepath, buf.Bytes(), 0644)
}

// copyStaticAssets copies files from static/ directory to public/
func copyStaticAssets() error {
	staticDir := "static"

	// Check if static directory exists
	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		return nil // No static directory, nothing to copy
	}

	return filepath.Walk(staticDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the root static directory itself
		if path == staticDir {
			return nil
		}

		// Calculate relative path
		relPath, err := filepath.Rel(staticDir, path)
		if err != nil {
			return err
		}

		destPath := filepath.Join("public", relPath)

		// Create directories
		if info.IsDir() {
			return os.MkdirAll(destPath, 0755)
		}

		// Copy files
		return copyFile(path, destPath)
	})
}

// copyFile copies a single file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

// generateRSSListing creates the RSS reading list page
func generateRSSListing(site Site) error {
	if len(site.FeedItems) == 0 {
		return nil // No items to display
	}

	tmpl := getBaseTemplate()
	template.Must(tmpl.Parse(rssListingContent))

	// Convert to display format with date strings
	type ItemDisplay struct {
		Title      string
		Link       string
		DateString string
		FeedTitle  string
	}

	items := make([]ItemDisplay, len(site.FeedItems))
	for i, item := range site.FeedItems {
		dateStr := ""
		if !item.PubDate.IsZero() {
			dateStr = item.PubDate.Format("2006-01-02")
		}
		items[i] = ItemDisplay{
			Title:      item.Title,
			Link:       item.Link,
			DateString: dateStr,
			FeedTitle:  item.FeedTitle,
		}
	}

	data := struct {
		Title        string
		Description  string
		CanonicalURL string
		Image        string
		OGType       string
		Items        []ItemDisplay
	}{
		Title:        "Reading — Vegard Stikbakke",
		Description:  "Articles from blogs I follow",
		CanonicalURL: "https://vegardstikbakke.com/reading/",
		Image:        "",
		OGType:       "website",
		Items:        items,
	}

	if err := os.MkdirAll("public/reading", 0755); err != nil {
		return err
	}

	return renderToFile(tmpl, data, "public/reading/index.html")
}

// RSS structs for XML generation
type RSSFeed struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Atom    string   `xml:"xmlns:atom,attr"`
	Channel RSSChannel
}

type RSSChannel struct {
	XMLName       xml.Name `xml:"channel"`
	Title         string   `xml:"title"`
	Link          string   `xml:"link"`
	Description   string   `xml:"description"`
	Language      string   `xml:"language"`
	LastBuildDate string   `xml:"lastBuildDate"`
	AtomLink      AtomLink `xml:"atom:link"`
	Items         []RSSItem
}

type AtomLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
	Type string `xml:"type,attr"`
}

type RSSItem struct {
	XMLName     xml.Name `xml:"item"`
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	PubDate     string   `xml:"pubDate"`
	GUID        string   `xml:"guid"`
	Description string   `xml:"description"`
}

// generateRSSFeed creates an RSS feed at /feed.xml
func generateRSSFeed(site Site) error {
	// Build items from posts
	items := make([]RSSItem, 0, len(site.Posts))
	for _, post := range site.Posts {
		// Format date as RFC1123Z (required for RSS)
		pubDate := post.Date.Format(time.RFC1123Z)

		// Create post URL
		postURL := fmt.Sprintf("https://vegardstikbakke.com/%s/", post.Slug)

		// Use description if available, otherwise use truncated content
		description := post.Description
		if description == "" {
			// Use raw HTML content as description
			description = post.HTMLContent
		}

		items = append(items, RSSItem{
			Title:       html.EscapeString(post.Title),
			Link:        postURL,
			PubDate:     pubDate,
			GUID:        postURL,
			Description: description, // Will be escaped by xml.Marshal
		})
	}

	// Determine last build date (use most recent post date)
	lastBuildDate := time.Now().Format(time.RFC1123Z)
	if len(site.Posts) > 0 {
		lastBuildDate = site.Posts[0].Date.Format(time.RFC1123Z)
	}

	feed := RSSFeed{
		Version: "2.0",
		Atom:    "http://www.w3.org/2005/Atom",
		Channel: RSSChannel{
			Title:         "Vegard Stikbakke",
			Link:          "https://vegardstikbakke.com/",
			Description:   "Software engineer from Norway",
			Language:      "en-us",
			LastBuildDate: lastBuildDate,
			AtomLink: AtomLink{
				Href: "https://vegardstikbakke.com/feed.xml",
				Rel:  "self",
				Type: "application/rss+xml",
			},
			Items: items,
		},
	}

	// Marshal to XML
	output, err := xml.MarshalIndent(feed, "", "  ")
	if err != nil {
		return err
	}

	// Add XML header
	xmlContent := []byte(xml.Header + string(output))

	// Write to file
	return os.WriteFile("public/feed.xml", xmlContent, 0644)
}
