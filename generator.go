package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
)

// generateHomepage creates the homepage with bio from about.md
func generateHomepage(site Site) error {
	tmpl := getBaseTemplate()
	template.Must(tmpl.Parse(homepageContent))

	data := struct {
		Title        string
		Description  string
		CanonicalURL string
		Content      template.HTML
	}{
		Title:        "Vegard Stikbakke",
		Description:  "Vegard Stikbakke — software engineer from Norway",
		CanonicalURL: "https://vegardstikbakke.com/",
		Content:      template.HTML(site.AboutPage.HTMLContent),
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
		Posts        []Post
	}{
		Title:        "Posts - Vegard Stikbakke",
		Description:  "Blog posts by Vegard Stikbakke",
		CanonicalURL: "https://vegardstikbakke.com/blog/",
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

	data := struct {
		Title        string
		Description  string
		CanonicalURL string
		Books        []Book
	}{
		Title:        "Books - Vegard Stikbakke",
		Description:  "Books read by Vegard Stikbakke",
		CanonicalURL: "https://vegardstikbakke.com/books/",
		Books:        site.Books,
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
			Content      template.HTML
		}{
			Title:        post.Title + " - Vegard Stikbakke",
			PostTitle:    post.Title,
			DateString:   post.DateString,
			Description:  post.Description,
			CanonicalURL: fmt.Sprintf("https://vegardstikbakke.com/%s/", post.Slug),
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
