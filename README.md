# Static Site Generator

A minimal static site generator written in Go that converts markdown blog posts to HTML.

## Features

- Converts markdown posts to HTML with clean URLs
- Preserves slugs from frontmatter
- Filters out draft posts
- Minimal, readable styling
- Homepage, posts listing, books page, and individual post pages

## Quick Start with Make

The easiest way to build and run the site locally:

```bash
# Build the generator and generate the site
make run

# Serve the site locally at http://localhost:8000
make serve
```

### Available Make Targets

- `make build` - Build the static site generator binary
- `make generate` - Generate the static site in `public/`
- `make serve` - Serve the generated site locally on port 8000
- `make run` - Build the generator and generate the site
- `make clean` - Remove generated files (binary and `public/` directory)
- `make help` - Show available targets

## Manual Building (without Make)

If you prefer not to use Make:

```bash
# Build the static site generator
go build -o ssg

# Generate the static site in public/
./ssg
```

This will:
- Read markdown files from `content/blog/`
- Filter out posts with `draft: true`
- Generate HTML in `public/` directory

## Serving Locally

The generated site uses absolute paths, so you need to serve it with a web server:

```bash
# Using Make
make serve

# Or using Python directly
cd public
python3 -m http.server 8000
```

Then open http://localhost:8000 in your browser.

## Content Structure

Blog posts should have YAML frontmatter:

```yaml
---
title: My Post Title
slug: my-post-slug
date: "2024-01-01"
draft: false
---

Your markdown content here...
```

- `title` (required): Post title
- `slug` (optional): URL slug (defaults to filename)
- `date` (optional): Publication date
- `draft` (optional): Set to `true` to exclude from site

## Output Structure

```
public/
├── index.html              # Homepage with bio
├── posts/index.html        # Posts listing
├── books/index.html        # Books page
└── blog/
    └── {slug}/index.html   # Individual posts
```
