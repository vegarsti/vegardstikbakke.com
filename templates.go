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

        /* Catppuccin Latte */
        :root {
            color-scheme: light;
            --bg-color: #eff1f5;
            --text-color: #4c4f69;
            --text-secondary: #6c6f85;
            --link-color: #1e66f5;
            --link-hover: #209fb5;
            --link-visited: #8839ef;
            --code-bg: #e6e9ef;
            --code-border: #ccd0da;
            --border: #bcc0cc;
            --badge-bg: #40a02b;
            --badge-text: #eff1f5;
            --accent-light: #e6e9ef;
            --accent-muted: #179299;
        }

        /* Catppuccin Mocha */
        @media (prefers-color-scheme: dark) {
            :root {
                color-scheme: dark;
                --bg-color: #1e1e2e;
                --text-color: #cdd6f4;
                --text-secondary: #a6adc8;
                --link-color: #89b4fa;
                --link-hover: #74c7ec;
                --link-visited: #cba6f7;
                --code-bg: #313244;
                --code-border: #45475a;
                --border: #45475a;
                --badge-bg: #a6e3a1;
                --badge-text: #1e1e2e;
                --accent-light: #181825;
                --accent-muted: #94e2d5;
            }
        }

        body {
            max-width: 640px;
            margin: 0 auto;
            padding: 64px 24px 48px;
            font-family: 'Geist', -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
            font-size: 16px;
            letter-spacing: -0.01em;
            line-height: 1.65;
            background-color: var(--bg-color);
            color: var(--text-color);
        }

        .site-header {
            display: flex;
            align-items: flex-start;
            justify-content: space-between;
            gap: 24px;
            margin-bottom: 52px;
        }

        .site-name {
            display: inline-block;
            color: var(--text-color);
            font-size: 1.125rem;
            font-weight: 600;
            letter-spacing: -0.025em;
            line-height: 1.2;
            text-decoration: none;
        }

        .site-name:hover { color: var(--text-color); }

        .site-tagline {
            margin: 4px 0 0;
            color: var(--text-secondary);
            font-size: 0.875rem;
            line-height: 1.4;
        }

        .site-socials {
            display: flex;
            align-items: center;
            gap: 0.75rem;
            margin-top: 0.65rem;
        }

        .site-socials a {
            position: relative;
            color: var(--link-color);
            font-size: 0.875rem;
            line-height: 1;
        }

        .site-socials a:hover { color: var(--link-color); }

        .site-nav {
            display: flex;
            align-items: center;
            gap: 14px;
            padding-top: 1px;
        }

        .site-nav a {
            color: var(--text-secondary);
            font-size: 0.875rem;
            text-decoration: none;
        }

        .site-nav a:hover {
            color: var(--text-color);
            text-decoration: underline;
            text-underline-offset: 3px;
        }

        .home-intro::after {
            content: '';
            display: table;
            clear: both;
        }

        .section-link {
            margin: 0.75rem 0 0;
            font-size: 0.875rem;
        }

        .section-link a {
            color: var(--text-secondary);
            text-decoration: none;
        }

        .section-link a:hover {
            color: var(--text-color);
            text-decoration: underline;
            text-underline-offset: 3px;
        }

        .work-list {
            margin: 0 -12px;
        }

        .work-row {
            display: grid;
            grid-template-columns: 110px minmax(0, 1fr) 90px;
            align-items: baseline;
            column-gap: 1.5rem;
            padding: 0.7rem 12px;
            border-radius: 6px;
            color: var(--text-color);
            text-decoration: none;
            transition: background-color 0.15s ease;
        }

        .work-row:hover {
            color: var(--text-color);
            background: var(--accent-light);
            text-decoration: none;
        }

        .work-company {
            position: relative;
            display: block;
            padding-left: calc(16px + 0.6rem);
            font-weight: 500;
            letter-spacing: -0.015em;
        }

        .work-icon {
            position: absolute;
            top: 50%;
            left: 0;
            width: 16px;
            height: 16px;
            margin: 0;
            object-fit: contain;
            opacity: 0.55;
            transform: translateY(-50%);
        }

        .dune-work-icon { filter: invert(1); }

        @media (prefers-color-scheme: dark) {
            .dune-work-icon { filter: none; }
            .cognite-work-icon { filter: invert(1); }
        }

        .work-role,
        .work-period {
            color: var(--text-secondary);
            font-size: 0.875rem;
        }

        .editorial-row.book-row {
            grid-template-columns: minmax(0, 1fr) minmax(100px, auto) auto;
        }

        .post-star-slot {
            color: var(--text-color);
            font-size: 0.9em;
            line-height: 1;
            text-align: center;
        }

        .book-author {
            color: var(--text-secondary);
            font-size: 0.875rem;
        }

        .work-period {
            font-variant-numeric: tabular-nums;
            text-align: right;
            white-space: nowrap;
        }

        h1 { font-size: 1.75rem; letter-spacing: -0.035em; line-height: 1.2; margin-bottom: 0.75em; position: relative; }
        h2 { font-size: 1.25rem; letter-spacing: -0.025em; line-height: 1.3; margin-top: 2em; position: relative; }
        h3 { font-size: 1.1rem; letter-spacing: -0.015em; margin-top: 1.75em; position: relative; }

        .section-label {
            display: flex;
            align-items: baseline;
            gap: 0.75rem;
            margin: 3.5rem 0 1rem;
            color: var(--text-secondary);
            font-size: 0.8125rem;
            font-weight: 600;
            letter-spacing: 0;
            line-height: 1.2;
        }

        .section-label::after {
            content: '';
            flex: 1;
            border-bottom: 1px dotted currentColor;
            opacity: 0.5;
        }

        .heading-anchor {
            position: absolute;
            left: -1.2em;
            opacity: 0;
            font-weight: normal;
            text-decoration: none !important;
            color: var(--text-secondary);
            transition: opacity 0.15s;
        }

        h1:hover > .heading-anchor,
        h2:hover > .heading-anchor,
        h3:hover > .heading-anchor {
            opacity: 1;
        }

        .editorial-list {
            margin: 0 -12px;
        }

        .post-list {
            display: grid;
            grid-template-columns: 1.25em minmax(0, 1fr) 110px;
            column-gap: 0.35em;
            row-gap: 5px;
            align-items: baseline;
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

        .post-list .post-date {
            color: var(--text-color);
            font-variant-numeric: tabular-nums;
        }

        .posts-rss {
            margin-top: 2em;
        }

        .editorial-row {
            display: grid;
            grid-template-columns: minmax(0, 1fr) auto;
            align-items: baseline;
            column-gap: 0.75rem;
            padding: 0.45rem 12px;
            border-radius: 6px;
            color: var(--text-color);
            font-size: 0.9375rem;
            line-height: 1.45;
            text-decoration: none;
            transition: background-color 0.15s ease;
        }

        .editorial-row:hover {
            color: var(--text-color);
            background: var(--accent-light);
            text-decoration: none;
        }

        .editorial-title {
            min-width: 0;
            font-weight: 500;
            letter-spacing: -0.015em;
        }

        .editorial-leader { display: none; }

        .editorial-meta {
            color: var(--text-secondary);
            font-size: 0.875rem;
            font-variant-numeric: tabular-nums;
            white-space: nowrap;
        }

        .editorial-subtitle {
            color: var(--text-secondary);
            font-size: 0.875rem;
        }

        .post-star {
            color: var(--text-color);
        }

        .external-link-icon {
            margin-left: 0.15em;
            font-size: 0.75em;
        }

        .post-date { color: var(--text-color); font-variant-numeric: tabular-nums; }

        @media (max-width: 640px) {
            body { padding: 40px 20px; }
            .site-header { margin-bottom: 40px; }
            .site-nav { gap: 10px; }
            .editorial-row { grid-template-columns: minmax(0, 1fr) auto; row-gap: 0.15rem; }
            .editorial-leader { display: none; }
            .post-list { grid-template-columns: 1.25em minmax(0, 1fr); column-gap: 0.35em; row-gap: 0; }
            .post-list .post-date { grid-column: 2; margin-bottom: 0.5em; }
            .editorial-subtitle { grid-column: 1; }
            .work-row { grid-template-columns: 85px minmax(0, 1fr); column-gap: 1rem; row-gap: 0.2rem; }
            .work-role { grid-column: 2; grid-row: 2; }
            .work-period { grid-column: 2; grid-row: 1; justify-self: end; }
            .editorial-row.book-row { grid-template-columns: minmax(0, 1fr) auto; }
            .book-author { grid-column: 1; grid-row: 2; }
            h1 { font-size: 1.5rem; }
        }

        @media (max-width: 480px) {
            .site-header { flex-direction: column; gap: 14px; }
            .site-nav { padding-top: 0; }
        }

        .reading-list { margin: 0 -12px; }

        .reading-item { margin: 0; }

        .reading-title { font-weight: 500; letter-spacing: -0.015em; }

        .reading-title a { color: var(--text-color); text-decoration: none; }

        .reading-title a:hover { text-decoration: underline; }

        .reading-meta { color: var(--text-secondary); font-size: 0.875em; }

        .reading-source {
            font-weight: 500;
        }

        .draft-badge {
            display: inline-block;
            background: var(--badge-bg);
            color: var(--badge-text);
            font-size: 0.65em;
            font-weight: 600;
            padding: 1px 5px;
            border-radius: 999px;
            margin-left: 0.45em;
            text-transform: uppercase;
            letter-spacing: 0.04em;
            vertical-align: 0.12em;
        }

        pre, code {
            font-family: 'Iosevka', monospace;
            font-size: 1em;
            border-radius: 0.4em;
            tab-size: 4;
            -webkit-text-size-adjust: 100%;
            text-size-adjust: 100%;
        }

        pre {
            background: var(--code-bg);
            padding: 20px;
            overflow-x: auto;
            border: 1px solid var(--code-border);
            font-size: 14px;
        }

        code {
            background: var(--code-bg);
            color: var(--text-color);
            padding: 2px 6px;
            border-radius: 0.4em;
        }

        a code {
            color: var(--link-color);
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
            margin: 0 0 1rem 2rem;
            border-radius: 0;
            border: 1px solid var(--text-color);
        }

        .profile-image-dark {
            display: none;
        }

        @media (prefers-color-scheme: dark) {
            .profile-image-light { display: none; }
            .profile-image-dark { display: block; }
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
            display: flex;
            flex-wrap: wrap;
            gap: 0.4rem 1rem;
            list-style: none;
            margin: 1.5rem 0 0;
            padding: 0;
        }

        .social-links li { padding: 0; }

        .social-links li::before { content: none; }

        .social-links a {
            display: flex;
            align-items: center;
            gap: 0.35em;
            color: var(--text-secondary);
            font-size: 0.875rem;
            text-decoration: none;
        }

        .social-links a:hover {
            color: var(--text-color);
            text-decoration: underline;
            text-underline-offset: 3px;
        }

        .social-links i { font-size: 1em; width: 1em; text-align: center; }

        /* Syntax highlighting - Catppuccin Latte */
        .chroma { background-color: var(--code-bg); }
        .chroma .lntd { vertical-align: top; padding: 0; margin: 0; border: 0; }
        .chroma .lntable { border-spacing: 0; padding: 0; margin: 0; border: 0; width: auto; overflow: auto; display: block; }
        .chroma .hl { background-color: #ccd0da; display: block; width: 100%; }
        .chroma .lnt { margin-right: 0.4em; padding: 0 0.4em; color: #8c8fa1; }
        .chroma .ln { margin-right: 0.4em; padding: 0 0.4em; color: #8c8fa1; }
        .chroma .line { display: flex; }
        .chroma .k { color: #8839ef; } /* Keyword - Mauve */
        .chroma .kc { color: #fe640b; } /* KeywordConstant - Peach */
        .chroma .kd { color: #8839ef; } /* KeywordDeclaration - Mauve */
        .chroma .kn { color: #8839ef; } /* KeywordNamespace - Mauve */
        .chroma .kp { color: #8839ef; } /* KeywordPseudo - Mauve */
        .chroma .kr { color: #8839ef; } /* KeywordReserved - Mauve */
        .chroma .kt { color: #df8e1d; } /* KeywordType - Yellow */
        .chroma .na { color: #df8e1d; } /* NameAttribute - Yellow */
        .chroma .nb { color: #d20f39; } /* NameBuiltin - Red */
        .chroma .nc { color: #df8e1d; } /* NameClass - Yellow */
        .chroma .no { color: #df8e1d; } /* NameConstant - Yellow */
        .chroma .nd { color: #1e66f5; } /* NameDecorator - Blue */
        .chroma .ni { color: #4c4f69; } /* NameEntity - Text */
        .chroma .ne { color: #fe640b; } /* NameException - Peach */
        .chroma .nf { color: #1e66f5; } /* NameFunction - Blue */
        .chroma .nl { color: #209fb5; } /* NameLabel - Sapphire */
        .chroma .nn { color: #df8e1d; } /* NameNamespace - Yellow */
        .chroma .nt { color: #8839ef; } /* NameTag - Mauve */
        .chroma .nv { color: #4c4f69; } /* NameVariable - Text */
        .chroma .s { color: #40a02b; } /* String - Green */
        .chroma .sa { color: #40a02b; } /* StringAffix - Green */
        .chroma .sb { color: #40a02b; } /* StringBacktick - Green */
        .chroma .sc { color: #40a02b; } /* StringChar - Green */
        .chroma .dl { color: #40a02b; } /* StringDelimiter - Green */
        .chroma .sd { color: #6c6f85; } /* StringDoc - Subtext0 */
        .chroma .s2 { color: #40a02b; } /* StringDouble - Green */
        .chroma .se { color: #fe640b; } /* StringEscape - Peach */
        .chroma .sh { color: #40a02b; } /* StringHeredoc - Green */
        .chroma .si { color: #40a02b; } /* StringInterpol - Green */
        .chroma .sx { color: #40a02b; } /* StringOther - Green */
        .chroma .sr { color: #fe640b; } /* StringRegex - Peach */
        .chroma .s1 { color: #40a02b; } /* StringSingle - Green */
        .chroma .ss { color: #40a02b; } /* StringSymbol - Green */
        .chroma .m { color: #fe640b; } /* Number - Peach */
        .chroma .mb { color: #fe640b; } /* NumberBin - Peach */
        .chroma .mf { color: #fe640b; } /* NumberFloat - Peach */
        .chroma .mh { color: #fe640b; } /* NumberHex - Peach */
        .chroma .mi { color: #fe640b; } /* NumberInteger - Peach */
        .chroma .il { color: #fe640b; } /* NumberIntegerLong - Peach */
        .chroma .mo { color: #fe640b; } /* NumberOct - Peach */
        .chroma .o { color: #04a5e5; } /* Operator - Sky */
        .chroma .ow { color: #04a5e5; } /* OperatorWord - Sky */
        .chroma .p { color: #7c7f93; } /* Punctuation - Overlay2 */
        .chroma .c { color: #8c8fa1; } /* Comment - Overlay1 */
        .chroma .ch { color: #8c8fa1; } /* CommentHashbang - Overlay1 */
        .chroma .cm { color: #8c8fa1; } /* CommentMultiline - Overlay1 */
        .chroma .c1 { color: #8c8fa1; } /* CommentSingle - Overlay1 */
        .chroma .cs { color: #8c8fa1; } /* CommentSpecial - Overlay1 */
        .chroma .cp { color: #fe640b; } /* CommentPreproc - Peach */
        .chroma .cpf { color: #40a02b; } /* CommentPreprocFile - Green */
        .chroma .gd { color: #d20f39; background-color: #e6e9ef; } /* GenericDeleted - Red */
        .chroma .ge { font-style: italic; } /* GenericEmph */
        .chroma .gi { color: #40a02b; background-color: #e6e9ef; } /* GenericInserted - Green */
        .chroma .gs { font-weight: bold; } /* GenericStrong */
        .chroma .gu { color: #df8e1d; font-weight: bold; } /* GenericSubheading - Yellow */

        /* Syntax highlighting - Catppuccin Mocha */
        @media (prefers-color-scheme: dark) {
            .chroma .hl { background-color: #45475a; }
            .chroma .lnt { color: #7f849c; }
            .chroma .ln { color: #7f849c; }
            .chroma .k { color: #cba6f7; } /* Keyword - Mauve */
            .chroma .kc { color: #fab387; } /* KeywordConstant - Peach */
            .chroma .kd { color: #cba6f7; } /* KeywordDeclaration - Mauve */
            .chroma .kn { color: #cba6f7; } /* KeywordNamespace - Mauve */
            .chroma .kp { color: #cba6f7; } /* KeywordPseudo - Mauve */
            .chroma .kr { color: #cba6f7; } /* KeywordReserved - Mauve */
            .chroma .kt { color: #f9e2af; } /* KeywordType - Yellow */
            .chroma .na { color: #f9e2af; } /* NameAttribute - Yellow */
            .chroma .nb { color: #f38ba8; } /* NameBuiltin - Red */
            .chroma .nc { color: #f9e2af; } /* NameClass - Yellow */
            .chroma .no { color: #f9e2af; } /* NameConstant - Yellow */
            .chroma .nd { color: #89b4fa; } /* NameDecorator - Blue */
            .chroma .ni { color: #cdd6f4; } /* NameEntity - Text */
            .chroma .ne { color: #fab387; } /* NameException - Peach */
            .chroma .nf { color: #89b4fa; } /* NameFunction - Blue */
            .chroma .nl { color: #74c7ec; } /* NameLabel - Sapphire */
            .chroma .nn { color: #f9e2af; } /* NameNamespace - Yellow */
            .chroma .nt { color: #cba6f7; } /* NameTag - Mauve */
            .chroma .nv { color: #cdd6f4; } /* NameVariable - Text */
            .chroma .s { color: #a6e3a1; } /* String - Green */
            .chroma .sa { color: #a6e3a1; } /* StringAffix - Green */
            .chroma .sb { color: #a6e3a1; } /* StringBacktick - Green */
            .chroma .sc { color: #a6e3a1; } /* StringChar - Green */
            .chroma .dl { color: #a6e3a1; } /* StringDelimiter - Green */
            .chroma .sd { color: #a6adc8; } /* StringDoc - Subtext0 */
            .chroma .s2 { color: #a6e3a1; } /* StringDouble - Green */
            .chroma .se { color: #fab387; } /* StringEscape - Peach */
            .chroma .sh { color: #a6e3a1; } /* StringHeredoc - Green */
            .chroma .si { color: #a6e3a1; } /* StringInterpol - Green */
            .chroma .sx { color: #a6e3a1; } /* StringOther - Green */
            .chroma .sr { color: #fab387; } /* StringRegex - Peach */
            .chroma .s1 { color: #a6e3a1; } /* StringSingle - Green */
            .chroma .ss { color: #a6e3a1; } /* StringSymbol - Green */
            .chroma .m { color: #fab387; } /* Number - Peach */
            .chroma .mb { color: #fab387; } /* NumberBin - Peach */
            .chroma .mf { color: #fab387; } /* NumberFloat - Peach */
            .chroma .mh { color: #fab387; } /* NumberHex - Peach */
            .chroma .mi { color: #fab387; } /* NumberInteger - Peach */
            .chroma .il { color: #fab387; } /* NumberIntegerLong - Peach */
            .chroma .mo { color: #fab387; } /* NumberOct - Peach */
            .chroma .o { color: #89dceb; } /* Operator - Sky */
            .chroma .ow { color: #89dceb; } /* OperatorWord - Sky */
            .chroma .p { color: #9399b2; } /* Punctuation - Overlay2 */
            .chroma .c { color: #7f849c; } /* Comment - Overlay1 */
            .chroma .ch { color: #7f849c; } /* CommentHashbang - Overlay1 */
            .chroma .cm { color: #7f849c; } /* CommentMultiline - Overlay1 */
            .chroma .c1 { color: #7f849c; } /* CommentSingle - Overlay1 */
            .chroma .cs { color: #7f849c; } /* CommentSpecial - Overlay1 */
            .chroma .cp { color: #fab387; } /* CommentPreproc - Peach */
            .chroma .cpf { color: #a6e3a1; } /* CommentPreprocFile - Green */
            .chroma .gd { color: #f38ba8; background-color: #313244; } /* GenericDeleted - Red */
            .chroma .gi { color: #a6e3a1; background-color: #313244; } /* GenericInserted - Green */
            .chroma .gu { color: #f9e2af; } /* GenericSubheading - Yellow */
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
    <header class="site-header">
        <div class="site-identity">
            <a class="site-name" href="/">Vegard Stikbakke</a>
            <div class="site-socials" aria-label="Social links">
                <a href="https://github.com/vegarsti" aria-label="GitHub"><i class="fab fa-github" aria-hidden="true"></i></a>
                <a href="https://twitter.com/vegardstikbakke" aria-label="Twitter"><i class="fab fa-twitter" aria-hidden="true"></i></a>
                <a href="https://www.linkedin.com/in/vegardstikbakke/" aria-label="LinkedIn"><i class="fab fa-linkedin" aria-hidden="true"></i></a>
                <a href="https://www.goodreads.com/user/show/3400170-vegard-stikbakke" aria-label="Goodreads"><i class="fab fa-goodreads" aria-hidden="true"></i></a>
                <a href="mailto:vegard.stikbakke@gmail.com" aria-label="Email"><i class="fas fa-envelope" aria-hidden="true"></i></a>
            </div>
        </div>
        <nav class="site-nav" aria-label="Main navigation">
            <a href="/blog/">Posts</a>
            <a href="/books/">Books</a>
        </nav>
    </header>
    <main>
        {{template "content" .}}
    </main>
    {{if .CollapsibleCode}}<script src="/collapsible-code.js"></script>{{end}}
</body>
</html>`

// Homepage template (shows a concise bio and curated recent activity)
var homepageContent = `{{define "content"}}
<div class="home-intro">
    <img src="/me.jpg" alt="Vegard Stikbakke" class="profile-image profile-image-light">
    <img src="/me-pixel.png" alt="Vegard Stikbakke" class="profile-image profile-image-dark">
    {{.Content}}
</div>

<section aria-labelledby="selected-work">
    <h2 id="selected-work" class="section-label">Work</h2>
    <div class="work-list">
        <a class="work-row" href="https://earendil.com" target="_blank" rel="noopener">
            <span class="work-company"><img class="work-icon" src="/work/earendil.png" alt="">Earendil</span>
            <span class="work-role">Member of Technical Staff</span>
            <span class="work-period">2026–present</span>
        </a>
        <a class="work-row" href="https://dune.com" target="_blank" rel="noopener">
            <span class="work-company"><img class="work-icon dune-work-icon" src="/work/dune.png" alt="">Dune</span>
            <span class="work-role">Senior Software Engineer</span>
            <span class="work-period">2021–2026</span>
        </a>
        <a class="work-row" href="https://www.cognite.com" target="_blank" rel="noopener">
            <span class="work-company"><img class="work-icon cognite-work-icon" src="/work/cognite.png" alt="">Cognite</span>
            <span class="work-role">Software Engineer</span>
            <span class="work-period">2019–2020</span>
        </a>
    </div>
</section>

<section aria-labelledby="recent-writing">
    <h2 id="recent-writing" class="section-label">Occasional blog posts</h2>
    <div class="editorial-list">
    {{range .RecentPosts}}
        <a class="editorial-row" href="{{if .ExternalURL}}{{.ExternalURL}}{{else}}/{{.Slug}}/{{end}}">
            <span class="editorial-title">{{.Title}}</span>
            <span class="editorial-leader" aria-hidden="true"></span>
            <span class="editorial-meta">{{.DateString}}</span>
        </a>
    {{end}}
    </div>
    <p class="section-link"><a href="/blog/">All posts →</a></p>
</section>

{{if .RecentBooks}}
<section aria-labelledby="recent-books">
    <h2 id="recent-books" class="section-label">Books I've read recently</h2>
    <div class="editorial-list">
    {{range .RecentBooks}}
        <a class="editorial-row book-row" href="{{.Link}}" target="_blank" rel="noopener">
            <span class="editorial-title">{{.Title}}</span>
            <span class="book-author">{{.AuthorName}}</span>
            <span class="editorial-meta">{{.ReadAtISO}}</span>
        </a>
    {{end}}
    </div>
    <p class="section-link"><a href="/books/">All books →</a></p>
</section>
{{end}}
{{end}}`

// Posts listing template
var postsListingContent = `{{define "content"}}
<div class="posts-page">
<h1>Posts</h1>
<p class="posts-intro">A star (<span class="post-star" title="Starred post" aria-label="Starred post">★</span>) marks posts I'm particularly happy with.</p>
<div class="post-list">
{{range .Posts}}
    <div class="post-star-slot">{{if .Starred}}<span class="post-star" title="Starred post" aria-label="Starred post">★</span>{{end}}</div>
    <div class="post-title"><a href="{{if .ExternalURL}}{{.ExternalURL}}{{else}}/{{.Slug}}/{{end}}">{{.Title}}{{if .ExternalURL}} <i class="fa-solid fa-arrow-up-right-from-square external-link-icon" aria-hidden="true"></i>{{end}}</a>{{if .Draft}}<span class="draft-badge">Draft</span>{{end}}</div>
    <div class="post-date">{{.DateString}}</div>
{{end}}
</div>
<p class="posts-rss"><a href="/feed.xml">RSS feed</a></p>
</div>
{{end}}`

// Individual post template
var postContent = `{{define "content"}}
<h1>{{.PostTitle}}{{if .Draft}} <span class="draft-badge">Draft</span>{{end}}</h1>
{{if .DateString}}<p class="post-date">{{.DateString}}</p>{{end}}
{{.Content}}
{{end}}`

// RSS reading list template
var rssListingContent = `{{define "content"}}
<h1>Reading</h1>
<p>Articles from blogs I follow.</p>
<div class="editorial-list">
{{range .Items}}
    <a class="editorial-row" href="{{.Link}}" target="_blank" rel="noopener">
        <span class="editorial-title">{{.Title}}</span>
        <span class="editorial-leader" aria-hidden="true"></span>
        <span class="editorial-meta">{{.FeedTitle}} · {{.DateString}}</span>
    </a>
{{end}}
</div>
{{end}}`

// Books listing template
var booksListContent = `{{define "content"}}
<h1>Books</h1>
{{if .YearlyBookPosts}}
<section aria-labelledby="yearly-reading">
    <h2 id="yearly-reading" class="section-label">Yearly reading</h2>
    <div class="editorial-list">
    {{range .YearlyBookPosts}}
        <a class="editorial-row" href="/{{.Slug}}/">
            <span class="editorial-title">{{.Title}}</span>
            <span class="editorial-leader" aria-hidden="true"></span>
            <span class="editorial-meta">{{.DateString}}</span>
        </a>
    {{end}}
    </div>
</section>
{{end}}
<section aria-labelledby="all-books">
    <h2 id="all-books" class="section-label">All books</h2>
    <div class="editorial-list">
{{range .Books}}
    <a class="editorial-row book-row" href="{{.Link}}" target="_blank" rel="noopener">
        <span class="editorial-title">{{.Title}}</span>
        <span class="book-author">{{.AuthorName}}</span>
        <span class="editorial-meta">{{.ReadAtISO}}</span>
    </a>
{{end}}
    </div>
</section>
{{end}}`

// 404 page template
var notFoundContent = `{{define "content"}}
<h1>404 — Page Not Found</h1>
<p>The page you're looking for doesn't exist.</p>
<p><a href="/">← Back to home</a></p>
{{end}}`

func getBaseTemplate() *template.Template {
	return template.Must(template.New("base").Parse(baseTemplate))
}
