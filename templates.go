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
            color: var(--text-color);
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

        /* Syntax highlighting - Light mode (GitHub-inspired) */
        .chroma { background-color: var(--code-bg); }
        .chroma .lntd { vertical-align: top; padding: 0; margin: 0; border: 0; }
        .chroma .lntable { border-spacing: 0; padding: 0; margin: 0; border: 0; width: auto; overflow: auto; display: block; }
        .chroma .hl { background-color: #ffffcc; display: block; width: 100%; }
        .chroma .lnt { margin-right: 0.4em; padding: 0 0.4em; color: #7f7f7f; }
        .chroma .ln { margin-right: 0.4em; padding: 0 0.4em; color: #7f7f7f; }
        .chroma .line { display: flex; }
        .chroma .k { color: #d73a49; } /* Keyword */
        .chroma .kc { color: #005cc5; } /* KeywordConstant */
        .chroma .kd { color: #d73a49; } /* KeywordDeclaration */
        .chroma .kn { color: #d73a49; } /* KeywordNamespace */
        .chroma .kp { color: #d73a49; } /* KeywordPseudo */
        .chroma .kr { color: #d73a49; } /* KeywordReserved */
        .chroma .kt { color: #d73a49; } /* KeywordType */
        .chroma .na { color: #005cc5; } /* NameAttribute */
        .chroma .nb { color: #005cc5; } /* NameBuiltin */
        .chroma .nc { color: #6f42c1; } /* NameClass */
        .chroma .no { color: #005cc5; } /* NameConstant */
        .chroma .nd { color: #6f42c1; } /* NameDecorator */
        .chroma .ni { color: #24292e; } /* NameEntity */
        .chroma .ne { color: #6f42c1; } /* NameException */
        .chroma .nf { color: #6f42c1; } /* NameFunction */
        .chroma .nl { color: #005cc5; } /* NameLabel */
        .chroma .nn { color: #24292e; } /* NameNamespace */
        .chroma .nt { color: #22863a; } /* NameTag */
        .chroma .nv { color: #e36209; } /* NameVariable */
        .chroma .s { color: #032f62; } /* String */
        .chroma .sa { color: #032f62; } /* StringAffix */
        .chroma .sb { color: #032f62; } /* StringBacktick */
        .chroma .sc { color: #032f62; } /* StringChar */
        .chroma .dl { color: #032f62; } /* StringDelimiter */
        .chroma .sd { color: #6a737d; } /* StringDoc */
        .chroma .s2 { color: #032f62; } /* StringDouble */
        .chroma .se { color: #032f62; } /* StringEscape */
        .chroma .sh { color: #032f62; } /* StringHeredoc */
        .chroma .si { color: #032f62; } /* StringInterpol */
        .chroma .sx { color: #032f62; } /* StringOther */
        .chroma .sr { color: #032f62; } /* StringRegex */
        .chroma .s1 { color: #032f62; } /* StringSingle */
        .chroma .ss { color: #032f62; } /* StringSymbol */
        .chroma .m { color: #005cc5; } /* Number */
        .chroma .mb { color: #005cc5; } /* NumberBin */
        .chroma .mf { color: #005cc5; } /* NumberFloat */
        .chroma .mh { color: #005cc5; } /* NumberHex */
        .chroma .mi { color: #005cc5; } /* NumberInteger */
        .chroma .il { color: #005cc5; } /* NumberIntegerLong */
        .chroma .mo { color: #005cc5; } /* NumberOct */
        .chroma .o { color: #d73a49; } /* Operator */
        .chroma .ow { color: #d73a49; } /* OperatorWord */
        .chroma .p { color: #24292e; } /* Punctuation */
        .chroma .c { color: #6a737d; } /* Comment */
        .chroma .ch { color: #6a737d; } /* CommentHashbang */
        .chroma .cm { color: #6a737d; } /* CommentMultiline */
        .chroma .c1 { color: #6a737d; } /* CommentSingle */
        .chroma .cs { color: #6a737d; } /* CommentSpecial */
        .chroma .cp { color: #d73a49; } /* CommentPreproc */
        .chroma .cpf { color: #032f62; } /* CommentPreprocFile */
        .chroma .gd { color: #b31d28; background-color: #ffeef0; } /* GenericDeleted */
        .chroma .ge { font-style: italic; } /* GenericEmph */
        .chroma .gi { color: #22863a; background-color: #f0fff4; } /* GenericInserted */
        .chroma .gs { font-weight: bold; } /* GenericStrong */
        .chroma .gu { color: #6f42c1; font-weight: bold; } /* GenericSubheading */

        /* Syntax highlighting - Dark mode */
        @media (prefers-color-scheme: dark) {
            .chroma .hl { background-color: #3b3b00; }
            .chroma .k { color: #ff7b72; } /* Keyword */
            .chroma .kc { color: #79c0ff; } /* KeywordConstant */
            .chroma .kd { color: #ff7b72; } /* KeywordDeclaration */
            .chroma .kn { color: #ff7b72; } /* KeywordNamespace */
            .chroma .kp { color: #ff7b72; } /* KeywordPseudo */
            .chroma .kr { color: #ff7b72; } /* KeywordReserved */
            .chroma .kt { color: #ff7b72; } /* KeywordType */
            .chroma .na { color: #79c0ff; } /* NameAttribute */
            .chroma .nb { color: #79c0ff; } /* NameBuiltin */
            .chroma .nc { color: #d2a8ff; } /* NameClass */
            .chroma .no { color: #79c0ff; } /* NameConstant */
            .chroma .nd { color: #d2a8ff; } /* NameDecorator */
            .chroma .ni { color: #c9d1d9; } /* NameEntity */
            .chroma .ne { color: #d2a8ff; } /* NameException */
            .chroma .nf { color: #d2a8ff; } /* NameFunction */
            .chroma .nl { color: #79c0ff; } /* NameLabel */
            .chroma .nn { color: #c9d1d9; } /* NameNamespace */
            .chroma .nt { color: #7ee787; } /* NameTag */
            .chroma .nv { color: #ffa657; } /* NameVariable */
            .chroma .s { color: #a5d6ff; } /* String */
            .chroma .sa { color: #a5d6ff; } /* StringAffix */
            .chroma .sb { color: #a5d6ff; } /* StringBacktick */
            .chroma .sc { color: #a5d6ff; } /* StringChar */
            .chroma .dl { color: #a5d6ff; } /* StringDelimiter */
            .chroma .sd { color: #8b949e; } /* StringDoc */
            .chroma .s2 { color: #a5d6ff; } /* StringDouble */
            .chroma .se { color: #a5d6ff; } /* StringEscape */
            .chroma .sh { color: #a5d6ff; } /* StringHeredoc */
            .chroma .si { color: #a5d6ff; } /* StringInterpol */
            .chroma .sx { color: #a5d6ff; } /* StringOther */
            .chroma .sr { color: #a5d6ff; } /* StringRegex */
            .chroma .s1 { color: #a5d6ff; } /* StringSingle */
            .chroma .ss { color: #a5d6ff; } /* StringSymbol */
            .chroma .m { color: #79c0ff; } /* Number */
            .chroma .mb { color: #79c0ff; } /* NumberBin */
            .chroma .mf { color: #79c0ff; } /* NumberFloat */
            .chroma .mh { color: #79c0ff; } /* NumberHex */
            .chroma .mi { color: #79c0ff; } /* NumberInteger */
            .chroma .il { color: #79c0ff; } /* NumberIntegerLong */
            .chroma .mo { color: #79c0ff; } /* NumberOct */
            .chroma .o { color: #ff7b72; } /* Operator */
            .chroma .ow { color: #ff7b72; } /* OperatorWord */
            .chroma .p { color: #c9d1d9; } /* Punctuation */
            .chroma .c { color: #8b949e; } /* Comment */
            .chroma .ch { color: #8b949e; } /* CommentHashbang */
            .chroma .cm { color: #8b949e; } /* CommentMultiline */
            .chroma .c1 { color: #8b949e; } /* CommentSingle */
            .chroma .cs { color: #8b949e; } /* CommentSpecial */
            .chroma .cp { color: #ff7b72; } /* CommentPreproc */
            .chroma .cpf { color: #a5d6ff; } /* CommentPreprocFile */
            .chroma .gd { color: #ffa198; background-color: #490202; } /* GenericDeleted */
            .chroma .gi { color: #7ee787; background-color: #04260f; } /* GenericInserted */
            .chroma .gu { color: #d2a8ff; } /* GenericSubheading */
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
