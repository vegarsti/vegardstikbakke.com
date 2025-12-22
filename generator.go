package main

import (
	"bytes"
	"fmt"
	"html/template"
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
		CanonicalURL: "https://vegardstikbakke.com/posts/",
		Posts:        site.Posts,
	}

	if err := os.MkdirAll("public/posts", 0755); err != nil {
		return err
	}

	return renderToFile(tmpl, data, "public/posts/index.html")
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
			CanonicalURL: fmt.Sprintf("https://vegardstikbakke.com/blog/%s/", post.Slug),
			Content:      template.HTML(post.HTMLContent),
		}

		dir := filepath.Join("public/blog", post.Slug)
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
