# Books page runbook

The `/books/` page is generated from `data/books.ndjson`, a durable record of
Vegard's Goodreads read-shelf. **The site never touches the network at build
time** — it only reads that file — so refreshing the list is a two-step:
fetch into the file, then rebuild.

## TL;DR (the whole update)

```bash
make add-book      # finished a book? (after marking it read on Goodreads) fetch + rebuild
make serve         # (optional) preview at http://localhost:8000
```

To hide a book you don't want on `/books/`, edit `data/books-excluded.txt`
by hand (see "Hiding a book you don't want to show" below).

Then commit the changed `data/books.ndjson` (and `data/books-excluded.txt`
if you excluded anything) — see "Committing" below.

### Just finished a book

If you've already marked the book as read (with date and rating) on Goodreads,
the one-step version is:

```bash
make add-book      # fetches the latest read-shelf and rebuilds the site
```

This runs `fetch-books` + `run`. The fetcher reports `+ N new` — if it says
`0 new`, Goodreads hasn't propagated the change yet (or the book wasn't moved
to the read shelf); retry in a minute.

### Hiding a book you don't want to show

Deleting a line straight from `data/books.ndjson` does **not** work long-term:
the book is still on your Goodreads read-shelf, so the RSS feed keeps returning
it and the next `fetch-books` will silently re-add it. Instead, add its
Goodreads review link to `data/books-excluded.txt` (one per line; `#`
comments allowed):

```
https://www.goodreads.com/review/show/8550645517?utm_medium=api&utm_source=rss
```

To find a book's review link, look it up in `data/books.ndjson` — the `link`
field is the value to paste.

Then run `make fetch-books` (or `make add-book`). The fetcher reads that file
on every run and skips excluded books (reporting
`x N skipped (excluded, stays hidden)`), and drops any already-stored entries
that match, so they stay hidden forever. Then rebuild with `make run`.

### Full refresh (e.g. on a new machine)

```bash
make fetch-books   # pull latest read-shelf into data/books.ndjson (respects excludes)
make run           # rebuild the site from the (now-updated) ndjson
```

## What each step does

### 1. `make fetch-books`

Runs `scripts/fetch_goodreads.py`, which:

- Fetches the public RSS read-shelf:
  `https://www.goodreads.com/review/list_rss/3400170?shelf=read`
- **Merges** the result into `data/books.ndjson`, keyed by `link` (the
  Goodreads review URL, which is unique on the read shelf).
  - New reads → appended.
  - Existing reads → kept exactly as stored, even if Goodreads later changes
    the rating, title, or another field.
  - Books that scrolled off the RSS feed → **kept** (this is why we merge, not
    overwrite — see "The 100-book cap" below).
- Sorts newest-read first and atomically rewrites the file.
- Prints a diff summary, e.g. `+ 4 new  = 96 unchanged`.

A clean re-run with no new books shows `+ 0 new  = 100 unchanged`.

### 2. `make run`

Builds the `ssg` binary and regenerates the site. The generator's `loadBooks`
(`books.go`) decodes `data/books.ndjson` one object per line and emits
`public/books/index.html` via the `booksListContent` template in
`templates.go`.

If `data/books.ndjson` is absent or empty, `generateBooksList` silently skips
the page (logs a warning) — the build never hard-fails on books.

## Field reference

Each NDJSON line is one book, with these 9 fields (see `books.go`'s
`Book` struct). The fetcher trims everything else at write time.

Rendered on the page:
- `title`, `author_name` — shown in the list
- `link` — title links to the Goodreads review (also the fetcher's merge key)
- `user_rating` — rendered as star glyphs (`★`/`☆`) via `Book.Stars()`
- `read_at_iso` — date finished (ISO); sort key, newest first
- `book_published` — publication year, shown after the rating

Kept but not yet rendered:
- `book_image_url`, `book_medium_image_url`, `book_large_image_url` — cover
  art at three sizes (small/medium/large), available for a future cover view

Fields that used to be stored but were dropped to keep the file small and
noise-free: `book_description` (the bulk of the old size), `average_rating`,
`book_id`, `review_id`, `guid`, raw `user_read_at`/`user_date_added`,
`date_added_iso`, `fetched_at`, `user_review`, `user_shelves`, `isbn`. The
merge key moved from `review_id` to `link` (which embeds the review id and is
unique).

## The 100-book cap (important limitation)

Goodreads' RSS feed returns **only the 100 most-recently-read books**. Two
consequences:

1. **Run it regularly.** Every time you finish a book, fetch soon so it lands
   in the file before it scrolls off. Once a book is older than your 100
   newest, the RSS feed can't see it anymore.
2. **Never delete `data/books.ndjson`.** It is the only place older reads
   survive. Re-running the fetcher against today's RSS cannot recover a book
   that already scrolled off.

Lifetime count on Goodreads is ~330; the file currently holds the 100 newest.
Capturing the full history would need authenticated API access or scraping
the paginated shelf HTML, which is out of scope.

## Verifying it worked

```bash
# page generated?
ls -l public/books/index.html

# how many books rendered (should be ~the line count of the ndjson)?
grep -c 'class="reading-item"' public/books/index.html

# spot-check the newest entry
head -1 data/books.ndjson | python3 -m json.tool | grep -E 'title|read_at_iso|user_rating'
```

If `make fetch-books` fails (Goodreads down, network error):
- The existing `data/books.ndjson` is untouched (the rewrite is atomic via a
  `.tmp` file + `os.replace`).
- Just retry later. The site keeps showing the last-known data.

## Committing

Per the project's git rules, stage explicit paths only — never `git add -A`:

```bash
git add data/books.ndjson
git add data/books-excluded.txt   # only if you excluded a book
git commit -m "chore(books): refresh reading list"
```

If you also changed generator/template code, stage those files explicitly:

```bash
git add books.go generator.go templates.go types.go main.go scripts/fetch_goodreads.py
```

## Troubleshooting

- **`make fetch-books` shows 0 new but you just finished a book:** Goodreads
  may lag marking it as "read". Make sure the shelf is `read` and the
  `user_read_at` date is set on Goodreads.
- **A book's title looks escaped (`&#39;`):** that's intentional HTML-escaping
  by Go's `html/template`, which is safe. Don't unescape it.
- **A book you deleted from `data/books.ndjson` keeps coming back:** it's
  still on your Goodreads read-shelf, so the RSS feed returns it and the
  fetcher re-adds it. Add its review link to `data/books-excluded.txt` to
  hide it permanently (see "Hiding a book you don't want to show" above).
- **Want to force a full re-fetch:** you can't get older books back (see the
  100-cap), but to scrub a stale field just delete the file and re-run — you
  will get only the current 100. Excluded books (`data/books-excluded.txt`)
  stay excluded across such a reset.
