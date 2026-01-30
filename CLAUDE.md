# CLAUDE.md

This document provides guidance for AI assistants (like Claude Code) working with this codebase.

## Project Overview

This is a minimal static site generator written in Go that converts markdown blog posts to HTML. It's a personal blog/website with a focus on simplicity and readability.

## Architecture

The project consists of several Go files that work together:

- `main.go` - Entry point, orchestrates the generation process
- `parser.go` - Parses markdown files with YAML frontmatter
- `generator.go` - Generates HTML pages from parsed content
- `templates.go` - Contains HTML templates for different page types
- `types.go` - Defines data structures used throughout the codebase

Content is stored in `content/` directory:
- `content/blog/*.md` - Blog posts with YAML frontmatter
- `content/about.md` - About page
- `content/favorite-books.md` - Books page

Static assets are in `static/` directory.

Generated output goes to `public/` directory.

## Building and Running

Use the Makefile for common tasks:

```bash
make run      # Build generator and generate site
make serve    # Generate site and serve on localhost:8000
make build    # Build the ssg binary only
make generate # Generate site using existing binary
make clean    # Remove generated files
```

Or manually:
```bash
go build -o ssg
./ssg
```

## Key Patterns and Conventions

### Content Format

Blog posts must have YAML frontmatter:
```yaml
---
title: Post Title
slug: url-slug
date: "2024-01-01"
draft: false
---

Markdown content here...
```

- `draft: true` posts are filtered out during generation
- `slug` determines the URL path (defaults to filename if not specified)
- Generated URLs follow the pattern: `/blog/{slug}/`

### Code Organization

The generator follows a pipeline pattern:
1. Parse markdown files with frontmatter (parser.go)
2. Filter out drafts
3. Sort posts by date
4. Generate HTML pages using templates (generator.go)
5. Copy static assets to public directory

### Output Structure

```
public/
├── index.html              # Homepage with bio
├── posts/index.html        # Posts listing
├── books/index.html        # Books page
└── blog/
    └── {slug}/index.html   # Individual posts
```

## Important Notes for AI Assistants

1. **Don't modify content files** unless explicitly asked - the `content/` directory contains blog posts that should be preserved.

2. **The `old/` directory** contains legacy content from a previous Hugo-based setup. It's kept for reference but not used in generation.

3. **Testing changes**: After modifying Go code, run `make serve` to build, generate, and serve the site locally for testing.

4. **Keep it minimal**: This project intentionally avoids complexity. Don't add unnecessary features, dependencies, or abstractions unless specifically requested.

5. **Static assets**: Files in `static/` are copied as-is to `public/`. This includes CSS, fonts, images, etc.

6. **URL paths**: The generated site uses absolute paths starting with `/`, so it must be served with a web server (not opened directly as files).

## Common Tasks

### Adding a new blog post
Create a new `.md` file in `content/blog/` with proper frontmatter, then run `make generate`.

### Modifying page templates
Edit the template strings in `templates.go`. Templates use Go's `html/template` package.

### Changing site styling
All CSS is defined inline in `templates.go` within the `<style>` tag in `baseTemplate`. The `static/style.css` file exists but is empty/unused. CSS custom properties (variables) are defined in `:root` with dark mode variants in `@media (prefers-color-scheme: dark)`.

### Debugging generation issues
Run `./ssg` directly to see any error output from the generator.

## Dependencies

- Go 1.18+ (uses generics)
- Standard library only - no external dependencies
- Python 3 (optional, for local development server via `make serve`)

## Git Workflow

The main branch is `main`. Recent commits show the project was recently set up with:
- Plausible analytics integration
- Favicon support
- SEO optimizations
- Custom fonts

When making changes, follow the existing commit message style (imperative mood, concise descriptions).
