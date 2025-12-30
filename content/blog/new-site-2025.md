---
title: Rewriting my site using AI
slug: new-site-2025
date: 2025-12-30
---

Or: [You're telling me a vibe coded this site?](https://knowyourmeme.com/memes/youre-telling-me-a-shrimp-fried-this-rice)

This site is now built using a [small static site builder](https://github.com/vegarsti/vegardstikbakke.com) that Claude Code wrote for me.
Iterating on the site is so much fun now!

If you're a software engineer and not currently using AI for coding, I hope this post can entice you into giving it a try.
My hypothesis is that you're underestimating how useful it is!

## Year of vibes

In 2025, AI tooling for software engineers did so many leaps forward that it's hard to fathom.
My workflow has gone from copy/pasting snippets between ChatGPT and my editor + using GitHub Copilot in VSCode, to Cursor, to now mostly using agents[^1] like [Amp](https://ampcode.com/) and [Claude Code](https://claude.com/product/claude-code).
I'm now regularly using agents at work, both in new codebases and in old.
I use it for

- writing new application code
- writing tests
- modifying existing code
- debugging
- finding terminal commands

Writing a simple static site builder for this small blog is a perfect task for a coding agent, as they are especially good at starting from scratch.
Here's how I did it.

## First Claude session

To begin, I created a new directory and copied over my existing blog posts.
I opened `claude` and gave it this prompt:

```
I want to create a static site from scratch. The code to build it needs to
 be in Go. I want it to be very minimal: The site builder needs to be able
to take md blog posts and embed them. I have posts in the content repo.
Their slugs need to be preserved.
```

Some back and forth to refine its implementation plan, and a few minutes later I had a fully functioning site builder rendering my existing posts.

```
⏺ Perfect! I've successfully built your minimal static site generator in Go. Here's what was created:

  Generated Files

  - types.go - Data structures (Post, Page, Book, Site)
  - parser.go - Markdown parsing and frontmatter extraction
  - templates.go - HTML templates with minimal CSS
  - generator.go - HTML generation functions
  - main.go - Main orchestration logic
  - go.mod - Go module with dependencies

  Features Implemented

  ✓ Homepage - Displays bio from content/about.md
  ✓ Posts listing - Shows all published blog posts at /posts/
  ✓ Books page - Displays placeholder book data with ratings at /books/
  ✓ Individual posts - Each post at /blog/{slug}/
  ✓ Navigation - "Vegard Stikbakke" | "Posts" | "Books" on all pages
  ✓ Draft filtering - Excludes posts with draft: true (18 published posts)
  ✓ Slug preservation - Uses slug from frontmatter (falls back to filename)
  ✓ Clean styling - Minimal, readable CSS
```

You can see the full session [here](/prompts-for-new-site/initial.txt).

## More sessions

I then iterated on the functionality in a few more sessions, and you can browse those here:

- [Adding fonts and SEO by comparing to old site](/prompts-for-new-site/compare-features.txt)
- [Tweaking styling of the Posts list page](/prompts-for-new-site/styling.txt)
- [Deploying to Netlify](/prompts-for-new-site/netlify.txt)
- [Adding previews for draft posts when built locally](/prompts-for-new-site/draft.txt)
- [Adding system aware light/dark mode](/prompts-for-new-site/dark-mode.txt)

![Claude commit](/claude-commit.png)

Happy prompting!

[^1]: An LLM agent runs tools in a loop to achieve a goal ([via Simon Willison](https://simonwillison.net/2025/Sep/18/agents/)).
