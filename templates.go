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
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.5.1/css/all.min.css" integrity="sha512-DTOQO9RWCH3ppGqcWaEA1BIZOC6xxalwEsw9c2QQeAIftl+Vegovlnee1c9QX4TctnWMn13TZye+giMm8e2LwA==" crossorigin="anonymous" referrerpolicy="no-referrer" />
    <style>
        @import url('https://cdn.jsdelivr.net/npm/@fontsource/iosevka@5.0.4/index.min.css');
        @font-face {
            font-family: 'Geist';
            src: url('https://cdn.jsdelivr.net/npm/geist@1.2.0/dist/fonts/geist-sans/Geist-Regular.woff2') format('woff2');
            font-weight: 400;
            font-style: normal;
            font-display: swap;
        }
        @font-face {
            font-family: 'Geist';
            src: url('https://cdn.jsdelivr.net/npm/geist@1.2.0/dist/fonts/geist-sans/Geist-Medium.woff2') format('woff2');
            font-weight: 500;
            font-style: normal;
            font-display: swap;
        }
        @font-face {
            font-family: 'Geist';
            src: url('https://cdn.jsdelivr.net/npm/geist@1.2.0/dist/fonts/geist-sans/Geist-SemiBold.woff2') format('woff2');
            font-weight: 600;
            font-style: normal;
            font-display: swap;
        }

        :root {
            --bg-color: #fff;
            --text-color: #1d1d27;
            --text-secondary: #73738b;
            --link-color: #8a9e6b;
            --link-hover: #496495;
            --code-bg: #E0E7D7;
            --code-border: #b6b6c294;
            --border: #b6b6c2;
            --badge-bg: #4a7ddd;
            --badge-text: #fff;
        }

        @media (prefers-color-scheme: dark) {
            :root {
                --bg-color: #1d1d27;
                --text-color: #e8e8ed;
                --text-secondary: #a0a0b0;
                --link-color: #B7C396;
                --link-hover: #8bb8f8;
                --code-bg: #2a2c27;
                --code-border: #3a3a45;
                --border: #3a3a45;
                --badge-bg: #5c9cf5;
                --badge-text: #1d1d27;
            }
        }

        body {
            max-width: 700px;
            margin: 0 auto;
            padding: 40px 20px;
            font-family: 'Geist', -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
            font-size: 16px;
            letter-spacing: -0.01em;
            line-height: 1.6;
            background-color: var(--bg-color);
            color: var(--text-color);
        }

        nav {
            margin-bottom: 40px;
            padding-bottom: 20px;
            border-bottom: 1px solid var(--text-color);
        }

        nav a {
            margin-right: 20px;
            text-decoration: none;
            color: var(--text-color);
            font-weight: 500;
        }

        nav a:hover {
            color: var(--text-color);
            text-decoration: underline;
        }

        h1 { font-size: 2em; margin-bottom: 0.5em; }
        h2 { font-size: 1.5em; margin-top: 1.5em; }
        h3 { font-size: 1.25em; margin-top: 1.5em; }

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
            color: var(--text-color);
            font-variant-numeric: tabular-nums;
        }

        @media (max-width: 640px) {
            .post-list {
                grid-template-columns: 1fr;
                gap: 0.5em;
            }
            h1 {
                font-size: 1.75em;
            }
        }

        .reading-list {
            margin-top: 2em;
        }

        .reading-item {
            margin-bottom: 1em;
        }

        .reading-title {
            font-weight: 500;
            letter-spacing: -0.015em;
        }

        .reading-title a {
            color: var(--text-color);
            text-decoration: none;
        }

        .reading-title a:hover {
            text-decoration: underline;
        }

        .reading-meta {
            color: var(--text-secondary);
            font-size: 0.9em;
        }

        .reading-source {
            font-weight: 500;
        }

        .draft-badge {
            display: inline-block;
            background: var(--badge-bg);
            color: var(--badge-text);
            font-size: 0.7em;
            font-weight: 600;
            padding: 2px 6px;
            border-radius: 0.4em;
            margin-left: 8px;
            text-transform: uppercase;
            letter-spacing: 0.5px;
            vertical-align: middle;
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
            font-weight: 400;
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
        }

        .book-date {
            color: var(--text-secondary);
            font-variant-numeric: tabular-nums;
            font-family: 'Iosevka', monospace;
        }

        .book-rating {
            color: var(--link-color);
            white-space: nowrap;
        }

        @media (max-width: 640px) {
            .books {
                grid-template-columns: 1fr;
                gap: 0.5em;
            }
        }

        pre, code {
            font-family: 'Iosevka', monospace;
            font-size: 14px;
            border-radius: 0.4em;
        }

        pre {
            background: var(--code-bg);
            padding: 20px;
            overflow-x: auto;
            border: 1px solid var(--code-border);
        }

        code {
            background: var(--code-bg);
            color: var(--text-color);
            padding: 2px 6px;
            border-radius: 0.4em;
        }

        pre code {
            padding: 0;
            background: transparent;
        }

        /* Collapsible code blocks */
        .code-block-collapsible {
            position: relative;
            margin: 1.5em 0;
        }

        .code-block-collapsible.collapsed pre {
            max-height: 400px;
            overflow: hidden;
            position: relative;
        }

        .code-block-collapsible.collapsed pre::after {
            content: '';
            position: absolute;
            bottom: 0;
            left: 0;
            right: 0;
            height: 80px;
            background: linear-gradient(to bottom, transparent, var(--code-bg));
            pointer-events: none;
        }

        .code-block-collapsible pre {
            transition: max-height 0.3s ease;
            margin-bottom: 0;
        }

        .code-toggle {
            display: block;
            width: 100%;
            padding: 8px 15px;
            margin-top: 0;
            background: var(--code-bg);
            border: 1px solid var(--border);
            color: var(--text-color);
            font-family: 'Geist', sans-serif;
            font-size: 0.875rem;
            cursor: pointer;
            text-align: center;
            transition: background 0.2s ease, color 0.2s ease;
            border-radius: 0.4em;
        }

        .code-toggle:hover {
            background: var(--link-color);
            color: #fff;
        }

        .code-toggle:focus {
            outline: 2px solid var(--link-color);
            outline-offset: 2px;
        }

        .code-toggle:active {
            transform: translateY(1px);
        }

        @media (prefers-reduced-motion: reduce) {
            .code-block-collapsible pre {
                transition: none;
            }
        }

        a {
            color: var(--link-color);
            text-decoration: none;
        }

        a:hover {
            color: var(--link-color);
            text-decoration: underline;
        }

        .profile-image {
            width: 180px;
            height: 180px;
            object-fit: cover;
            float: right;
            margin-left: 2em;
            margin-bottom: 1em;
            border: 1px solid var(--text-color);
        }

        main img {
            max-width: 100%;
            height: auto;
            display: block;
            margin: 2em 0;
        }

        main {
            position: relative;
        }

        /* Blockquotes */
        blockquote {
            color: var(--text-secondary);
            margin-left: 14px;
            margin-right: 0px;
            border-left-color: var(--border);
            border-left-style: solid;
            border-left-width: 2px;
        }

        blockquote > p {
            color: var(--text-secondary);
            padding-left: 14px;
        }

        /* Horizontal rules */
        hr {
            border: none;
            border-top: 1px solid var(--border);
            margin-top: 48px;
            margin-bottom: 48px;
        }

        /* Lists */
        ul {
            padding-top: 3px;
            padding-bottom: 3px;
            padding-left: 24px;
        }

        ul li::before {
            content: none;
        }

        li {
            padding-top: 3px;
            padding-bottom: 3px;
            line-height: 25px;
        }

        ol {
            padding-left: 1.5em;
        }

        /* Social links with icons */
        .social-links {
            list-style: none;
            padding-left: 0;
            margin-top: 1em;
        }

        .social-links li {
            margin: 0.5em 0;
        }

        .social-links li::before {
            content: none;
        }

        .social-links a {
            display: flex;
            align-items: center;
            gap: 0.5em;
            color: var(--text-color);
            text-decoration: none;
            font-weight: 500;
        }

        .social-links a:hover {
            text-decoration: underline;
        }

        .social-links i {
            font-size: 1.2em;
            width: 1.2em;
            text-align: center;
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
    </nav>
    <main>
        {{template "content" .}}
    </main>
    <script src="/collapsible-code.js"></script>
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
<p style="margin-top: 2em;"><a href="/feed.xml">RSS feed</a></p>
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

// RSS reading list template
var rssListingContent = `{{define "content"}}
<h1>Reading</h1>
<p>Articles from blogs I follow, aggregated from their RSS feeds.</p>
<div class="reading-list">
{{range .Items}}
    <div class="reading-item">
        <div class="reading-title"><a href="{{.Link}}" target="_blank" rel="noopener">{{.Title}}</a></div>
        <div class="reading-meta"><span class="reading-source">{{.FeedTitle}}</span> · <span class="reading-date">{{.DateString}}</span></div>
    </div>
{{end}}
</div>
{{end}}`

// 404 page template
var notFoundContent = `{{define "content"}}
<h1>404 — Page Not Found</h1>
<p>The page you're looking for doesn't exist.</p>
<p><a href="/">← Back to home</a></p>
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
