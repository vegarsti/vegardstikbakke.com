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

    <!-- Open Graph / Facebook -->
    <meta property="og:type" content="{{if .OGType}}{{.OGType}}{{else}}website{{end}}">
    <meta property="og:title" content="{{.Title}}">
    {{if .Description}}<meta property="og:description" content="{{.Description}}">{{end}}
    {{if .CanonicalURL}}<meta property="og:url" content="{{.CanonicalURL}}">{{end}}
    {{if .Image}}<meta property="og:image" content="https://vegardstikbakke.com{{.Image}}">{{end}}

    <!-- Twitter -->
    <meta name="twitter:card" content="summary_large_image">
    <meta name="twitter:site" content="@vegardstikbakke">
    <meta name="twitter:creator" content="@vegardstikbakke">
    <meta name="twitter:title" content="{{.Title}}">
    {{if .Description}}<meta name="twitter:description" content="{{.Description}}">{{end}}
    {{if .Image}}<meta name="twitter:image" content="https://vegardstikbakke.com{{.Image}}">{{end}}
    <link rel="alternate" type="application/rss+xml" title="Vegard Stikbakke" href="/feed.xml">
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <style>
        @import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600&family=Source+Code+Pro&display=swap');

        :root {
            --bg-color: #ffffff;
            --text-color: #111;
            --text-secondary: #666;
            --link-color: #0066cc;
            --border-color: #ddd;
            --code-bg: #f4f4f4;
            --badge-bg: #f0f0f0;
            --badge-text: #666;
        }

        @media (prefers-color-scheme: dark) {
            :root {
                --bg-color: #1a1a1a;
                --text-color: #e4e4e4;
                --text-secondary: #a0a0a0;
                --link-color: #66b3ff;
                --border-color: #404040;
                --code-bg: #2d2d2d;
                --badge-bg: #333;
                --badge-text: #999;
            }
        }

        body {
            max-width: 800px;
            margin: 40px auto;
            padding: 0 20px;
            font-family: 'Inter', -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
            line-height: 1.6;
            background-color: var(--bg-color);
            color: var(--text-color);
        }
        nav {
            margin-bottom: 40px;
            padding-bottom: 20px;
            border-bottom: 1px solid var(--border-color);
        }
        nav a {
            margin-right: 20px;
            text-decoration: none;
            color: var(--text-color);
            font-weight: 500;
        }
        nav a:hover {
            text-decoration: underline;
        }
        h1 { font-size: 2em; margin-bottom: 0.5em; }
        h2 { font-size: 1.5em; margin-top: 1.5em; }
        .post-list {
            display: grid;
            grid-template-columns: 1fr 110px;
            gap: 5px;
            margin-top: 2em;
        }
        .post-title {
            font-weight: 500;
            letter-spacing: -0.015em;
            text-overflow: ellipsis;
            white-space: nowrap;
            overflow: hidden;
        }
        .post-title a {
            color: var(--text-color);
            text-decoration: none;
        }
        .post-title a:hover {
            text-decoration: underline;
        }
        .post-date {
            color: var(--text-secondary);
            font-variant-numeric: tabular-nums;
        }
        @media (max-width: 640px) {
            .post-list {
                grid-template-columns: 1fr;
                gap: 0.5em;
            }
        }
        .draft-badge {
            display: inline-block;
            background: var(--badge-bg);
            color: var(--badge-text);
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
            grid-template-columns: 1fr 0.75fr 110px;
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
        .book-title a {
            color: var(--text-color);
            text-decoration: none;
        }
        .book-title a:hover {
            text-decoration: underline;
        }
        .book-author {
            color: var(--text-secondary);
            letter-spacing: -0.015em;
        }
        .book-date {
            color: var(--text-secondary);
            font-variant-numeric: tabular-nums;
        }
        .book-rating {
            color: var(--text-color);
            white-space: nowrap;
        }
        @media (max-width: 640px) {
            .books {
                grid-template-columns: 1fr;
                gap: 0.5em;
            }
        }
        pre {
            background: var(--code-bg);
            padding: 10px;
            overflow-x: auto;
            border-radius: 4px;
            font-family: 'Source Code Pro', monospace;
        }
        code {
            font-family: 'Source Code Pro', monospace;
        }
        pre code { padding: 0; }
        a { color: var(--link-color); }
        .profile-image {
            width: 180px;
            height: 180px;
            object-fit: cover;
            float: right;
            margin-left: 1.5em;
            margin-bottom: 1em;
        }
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
        <!-- TODO: Uncomment to show books again
        <a href="/books/">Books</a>
        -->
        <a href="https://github.com/vegarsti">GitHub</a>
        <a href="https://twitter.com/vegardstikbakke">Twitter</a>
    </nav>
    <main>
        {{template "content" .}}
    </main>
</body>
</html>`

// Homepage template (shows bio from about.md)
var homepageContent = `{{define "content"}}
<img src="/me.jpg" alt="Vegard Stikbakke" class="profile-image">
{{.Content}}
{{end}}`

// Posts listing template
var postsListingContent = `{{define "content"}}
<h1>Posts</h1>
<div class="post-list">
{{range .Posts}}
    <div class="post-title"><a href="/{{.Slug}}/">{{.Title}}</a>{{if .Draft}}<span class="draft-badge">Draft</span>{{end}}</div>
    <div class="post-date">{{.DateString}}</div>
{{end}}
</div>
{{end}}`

// Books listing template
var booksListingContent = `{{define "content"}}
<h1>Books</h1>
<div class="books">
{{range .Books}}
    <div class="book-title"><a href="/books/{{.Slug}}/">{{.Title}}</a></div>
    <div class="book-author">{{.Author}}</div>
    <div class="book-date">{{.DateRead}}</div>
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
