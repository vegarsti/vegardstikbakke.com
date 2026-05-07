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
    <script>
        (function() {
            try {
                var theme = localStorage.getItem('theme');
                if (theme === 'light' || theme === 'dark') {
                    document.documentElement.dataset.theme = theme;
                } else {
                    delete document.documentElement.dataset.theme;
                }
            } catch (e) {}
        })();
    </script>
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

        :root[data-theme="dark"] {
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

        /* Catppuccin Mocha */
        @media (prefers-color-scheme: dark) {
            :root:not([data-theme="light"]) {
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
            display: flex;
            align-items: center;
            gap: 20px;
            margin-bottom: 40px;
            padding-bottom: 20px;
            border-bottom: 1px solid var(--text-color);
        }

        nav a {
            text-decoration: none;
            color: var(--text-color);
            font-weight: 500;
        }

        nav a:hover {
            color: var(--text-color);
            text-decoration: underline;
        }

        .theme-toggle {
            display: inline-flex;
            align-items: center;
            justify-content: center;
            width: 2rem;
            height: 2rem;
            margin-left: auto;
            border: 1px solid var(--border);
            border-radius: 999px;
            background: transparent;
            color: var(--text-color);
            cursor: pointer;
            font: inherit;
            line-height: 1;
            transition: background 0.15s ease, border-color 0.15s ease, transform 0.15s ease;
        }

        .theme-toggle:hover {
            background: var(--accent-light);
            border-color: var(--text-secondary);
        }

        .theme-toggle:focus {
            outline: 2px solid var(--link-color);
            outline-offset: 2px;
        }

        .theme-toggle:active {
            transform: translateY(1px);
        }

        h1 { font-size: 2em; margin-bottom: 0.5em; position: relative; }
        h2 { font-size: 1.5em; margin-top: 1.5em; position: relative; }
        h3 { font-size: 1.25em; margin-top: 1.5em; position: relative; }

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
            color: var(--link-color);
            text-decoration: none;
            font-weight: 500;
        }

        .social-links a:hover {
            color: var(--text-color);
            text-decoration: none;
        }

        .social-links i {
            font-size: 1.2em;
            width: 1.2em;
            text-align: center;
        }

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
    <nav>
        <a href="/">Vegard Stikbakke</a>
        <a href="/blog/">Posts</a>
        <button class="theme-toggle" type="button" aria-label="Theme: auto. Switch to light mode" title="Theme: auto">◐</button>
    </nav>
    <main>
        {{template "content" .}}
    </main>
    <script>
        (function() {
            var button = document.querySelector('.theme-toggle');
            if (!button) return;

            var modes = ['auto', 'light', 'dark'];

            function storedTheme() {
                try {
                    var theme = localStorage.getItem('theme');
                    return modes.indexOf(theme) === -1 ? 'auto' : theme;
                } catch (e) {
                    return 'auto';
                }
            }

            function applyTheme(theme) {
                if (theme === 'light' || theme === 'dark') {
                    document.documentElement.dataset.theme = theme;
                } else {
                    delete document.documentElement.dataset.theme;
                }
            }

            function updateButton(theme) {
                var nextTheme = modes[(modes.indexOf(theme) + 1) % modes.length];
                var labels = { auto: '◐', light: '☀', dark: '☾' };
                button.textContent = labels[theme];
                button.setAttribute('title', 'Theme: ' + theme);
                button.setAttribute('aria-label', 'Theme: ' + theme + '. Switch to ' + nextTheme + ' mode');
            }

            applyTheme(storedTheme());
            updateButton(storedTheme());

            button.addEventListener('click', function() {
                var current = storedTheme();
                var nextTheme = modes[(modes.indexOf(current) + 1) % modes.length];
                applyTheme(nextTheme);
                try { localStorage.setItem('theme', nextTheme); } catch (e) {}
                updateButton(nextTheme);
            });
        })();
    </script>
    {{if .CollapsibleCode}}<script src="/collapsible-code.js"></script>{{end}}
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

// Individual post template
var postContent = `{{define "content"}}
<h1>{{.PostTitle}}{{if .Draft}} <span class="draft-badge">Draft</span>{{end}}</h1>
{{if .DateString}}<p class="post-date">{{.DateString}}</p>{{end}}
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
	return template.Must(template.New("base").Parse(baseTemplate))
}
