package main

import "html/template"

// Base layout template with navigation
var baseTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    {{if .Description}}<meta name="description" content="{{.Description}}">{{end}}
    <title>{{.Title}}</title>
    {{if .CanonicalURL}}<link rel="canonical" href="{{.CanonicalURL}}">{{end}}
    <link rel="icon" type="image/png" href="/favicon.png">
    <link rel="alternate" type="application/rss+xml" title="Vegard Stikbakke" href="/feed.xml">
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <style>
        @import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600&family=Source+Code+Pro&display=swap');

        body {
            max-width: 800px;
            margin: 40px auto;
            padding: 0 20px;
            font-family: 'Inter', -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
            line-height: 1.6;
            color: #333;
        }
        nav {
            margin-bottom: 40px;
            padding-bottom: 20px;
            border-bottom: 1px solid #ddd;
        }
        nav a {
            margin-right: 20px;
            text-decoration: none;
            color: #333;
            font-weight: 500;
        }
        nav a:hover {
            text-decoration: underline;
        }
        h1 { font-size: 2em; margin-bottom: 0.5em; }
        h2 { font-size: 1.5em; margin-top: 1.5em; }
        .post-list { list-style: none; padding: 0; }
        .post-list li { margin-bottom: 1em; }
        .post-date { color: #666; font-size: 0.9em; }
        .draft-badge {
            display: inline-block;
            background: #f0f0f0;
            color: #666;
            font-size: 0.7em;
            font-weight: 600;
            padding: 2px 6px;
            border-radius: 3px;
            margin-left: 8px;
            text-transform: uppercase;
            letter-spacing: 0.5px;
        }
        .books {
            display: grid;
            grid-template-columns: 1fr 0.75fr min-content 70px;
            gap: 1em 1.5em;
            margin-top: 2em;
        }
        .book-title, .book-author {
            text-overflow: ellipsis;
            white-space: nowrap;
            overflow: hidden;
        }
        .book-title {
            font-weight: 500;
            letter-spacing: -0.015em;
        }
        .book-author {
            color: #666;
            letter-spacing: -0.015em;
        }
        .book-date {
            color: #666;
            font-variant-numeric: tabular-nums;
        }
        .book-rating {
            color: #f39c12;
            white-space: nowrap;
        }
        @media (max-width: 640px) {
            .books {
                grid-template-columns: 1fr;
                gap: 0.5em;
            }
        }
        pre {
            background: #f4f4f4;
            padding: 10px;
            overflow-x: auto;
            border-radius: 4px;
            font-family: 'Source Code Pro', monospace;
        }
        code {
            background: #f4f4f4;
            padding: 2px 5px;
            border-radius: 3px;
            font-family: 'Source Code Pro', monospace;
        }
        pre code { padding: 0; }
        a { color: #0066cc; }
    </style>
    <!-- Privacy-friendly analytics by Plausible -->
    <script async src="https://plausible.io/js/pa-fY9B0CzGV3CMqDY5kvN_3.js"></script>
    <script>
      window.plausible=window.plausible||function(){(plausible.q=plausible.q||[]).push(arguments)},plausible.init=plausible.init||function(i){plausible.o=i||{}};
      plausible.init()
    </script>
</head>
<body>
    <nav>
        <a href="/">Vegard Stikbakke</a>
        <a href="/blog/">Posts</a>
        <a href="/books/">Books</a>
    </nav>
    <main>
        {{template "content" .}}
    </main>
</body>
</html>`

// Homepage template (shows bio from about.md)
var homepageContent = `{{define "content"}}
{{.Content}}
{{end}}`

// Posts listing template
var postsListingContent = `{{define "content"}}
<h1>Posts</h1>
<ul class="post-list">
{{range .Posts}}
    <li>
        <a href="/{{.Slug}}/">{{.Title}}</a>
        {{if .Draft}}<span class="draft-badge">Draft</span>{{end}}
        {{if .DateString}}<span class="post-date"> — {{.DateString}}</span>{{end}}
    </li>
{{end}}
</ul>
{{end}}`

// Books listing template
var booksListingContent = `{{define "content"}}
<h1>Books</h1>
<div class="books">
{{range .Books}}
    <div class="book-title"><a href="/books/{{.Slug}}/">{{.Title}}</a></div>
    <div class="book-author">{{.Author}}</div>
    <div class="book-date">{{.DateRead}}</div>
    <div class="book-rating">{{printf "%s" (stars .Rating)}}</div>
{{end}}
</div>
{{end}}`

// Individual post template
var postContent = `{{define "content"}}
<h1>{{.PostTitle}}</h1>
{{if .DateString}}<p class="post-date">{{.DateString}}</p>{{end}}
{{.Content}}
{{end}}`

// Individual book template
var bookContent = `{{define "content"}}
<h1>{{.BookTitle}}</h1>
<p class="post-date">
    by {{.Author}}
    {{if .YearPublished}} ({{.YearPublished}}){{end}}
    {{if .DateRead}} • Read: {{.DateRead}}{{end}}
    {{if .Rating}} • <span class="book-rating">{{printf "%s" (stars .Rating)}}</span>{{end}}
</p>
{{.Content}}
{{end}}`

func getBaseTemplate() *template.Template {
	tmpl := template.New("base")
	tmpl = tmpl.Funcs(template.FuncMap{
		"stars": func(rating int) string {
			result := ""
			for i := 0; i < rating; i++ {
				result += "★"
			}
			for i := rating; i < 5; i++ {
				result += "☆"
			}
			return result
		},
	})
	return template.Must(tmpl.Parse(baseTemplate))
}
