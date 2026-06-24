#!/usr/bin/env python3
"""
Fetch Vegard's Goodreads read-shelf via the public RSS feed and merge into a
durable NDJSON record at data/books.ndjson.

Why NDJSON: book descriptions contain commas/quotes/newlines, so CSV is
painful. NDJSON is one JSON object per line -> trivially appendable,
git-diffable line-by-line, and streamable.

Why merge (not overwrite): the RSS feed caps at the 100 most-recent reads.
Once a book scrolls off the feed it's gone from Goodreads' RSS forever, so we
must accumulate. Keying by the Goodreads review id means re-ratings and edits
update in place while genuinely new reads append.

Run:
    python3 scripts/fetch_goodreads.py
    make fetch-books
"""
import json
import os
import sys
import urllib.request
from datetime import datetime
from xml.etree import ElementTree as ET

USER_ID = "3400170"
SHELF = "read"
DATA_FILE = os.path.join(os.path.dirname(__file__), "..", "data", "books.ndjson")
EXCLUDE_FILE = os.path.join(os.path.dirname(__file__), "..", "data", "books-excluded.txt")
FEED_URL = f"https://www.goodreads.com/review/list_rss/{USER_ID}?shelf={SHELF}"
USER_AGENT = "vegardstikbakke.com/1.0 (static-site book fetcher)"

# Fields persisted to data/books.ndjson, in storage order. Only what the site
# renders plus the cover URLs and ISBN; descriptions/IDs/raw dates are dropped
# to keep the file small and noise-free. The merge key is `link` (it embeds
# the Goodreads review id and is unique on the read shelf).
FIELDS = [
    "title", "author_name", "user_rating", "read_at_iso", "book_published",
    "link",
    "book_image_url", "book_medium_image_url", "book_large_image_url",
]
# Raw RSS tags we read transiently to build the stored record.
RAW_TAGS = [
    "title", "author_name", "book_published", "user_rating",
    "user_read_at", "link",
    "book_image_url", "book_medium_image_url", "book_large_image_url",
]


def parse_rss_date(s):
    """'Thu, 7 May 2026 00:00:00 +0000' -> '2026-05-07' (or None)."""
    if not s:
        return None
    try:
        return datetime.strptime(s, "%a, %d %b %Y %H:%M:%S %z").date().isoformat()
    except ValueError:
        return None


def fetch_feed():
    req = urllib.request.Request(FEED_URL, headers={"User-Agent": USER_AGENT})
    with urllib.request.urlopen(req, timeout=30) as resp:
        return ET.parse(resp)


def item_to_book(item):
    def g(tag):
        e = item.find(tag)
        return e.text.strip() if e is not None and e.text else None

    raw = {t: g(t) for t in RAW_TAGS}
    book = {f: raw.get(f) for f in FIELDS}
    # Normalize the read date to ISO for stable sorting/storage; drop the raw
    # RSS form (user_read_at) since the ISO value is what we keep.
    book["read_at_iso"] = parse_rss_date(raw.get("user_read_at"))
    return book


def load_existing(path):
    books = {}
    if not os.path.exists(path):
        return books
    with open(path, "r", encoding="utf-8") as f:
        for line in f:
            line = line.strip()
            if not line:
                continue
            try:
                b = json.loads(line)
            except json.JSONDecodeError as e:
                print(f"  ! skipping malformed line: {e}", file=sys.stderr)
                continue
            books[b.get("link")] = b
    return books


def load_excluded(path):
    """Review links the user has chosen to hide from /books/.

    One link per line; # comments and blank lines ignored. A book in this set is
    never (re)added from the RSS feed, and is dropped from data/books.ndjson if
    present — so manual exclusions survive future fetches instead of being
    silently re-added by the RSS merge.
    """
    excluded = set()
    if not os.path.exists(path):
        return excluded
    with open(path, "r", encoding="utf-8") as f:
        for line in f:
            line = line.split("#", 1)[0].strip()
            if line:
                excluded.add(line)
    return excluded


def sort_key(b):
    # Newest read first; books with no read_at sink to the bottom.
    return (b.get("read_at_iso") or "0000-00-00", b.get("link") or "")


def main():
    os.makedirs(os.path.dirname(DATA_FILE), exist_ok=True)
    existing = load_existing(DATA_FILE)
    excluded = load_excluded(EXCLUDE_FILE)
    if excluded:
        print(f"Loaded {len(excluded)} excluded book(s) from {os.path.relpath(EXCLUDE_FILE)}")
    print(f"Loaded {len(existing)} existing books from {os.path.relpath(DATA_FILE)}")

    # Drop any previously-stored book that has since been excluded.
    dropped = [lk for lk in existing if lk in excluded]
    for lk in dropped:
        del existing[lk]

    print(f"Fetching {FEED_URL} ...")
    tree = fetch_feed()
    items = tree.findall(".//item")
    print(f"Got {len(items)} items from feed")

    skipped_excluded = 0
    new_links, updated_links = [], []
    for item in items:
        book = item_to_book(item)
        key = book["link"]
        if key in excluded:
            skipped_excluded += 1
            continue
        prev = existing.get(key)
        if prev is None:
            new_links.append(key)
        elif prev != book:
            updated_links.append(key)
        existing[key] = book  # upsert

    ordered = sorted(existing.values(), key=sort_key, reverse=True)
    tmp = DATA_FILE + ".tmp"
    with open(tmp, "w", encoding="utf-8") as f:
        for b in ordered:
            f.write(json.dumps(b, ensure_ascii=False) + "\n")
    os.replace(tmp, DATA_FILE)

    print(f"\nWrote {len(ordered)} books to {os.path.relpath(DATA_FILE)}")
    print(f"  + {len(new_links)} new")
    if new_links:
        for key in new_links[:10]:
            print(f"      • {existing[key]['title']}")
        if len(new_links) > 10:
            print(f"      … and {len(new_links) - 10} more")
    print(f"  ~ {len(updated_links)} updated (re-rated / edited)")
    print(f"  = {len(ordered) - len(new_links) - len(updated_links)} unchanged")
    if dropped:
        print(f"  - {len(dropped)} removed (now excluded)")
    if skipped_excluded:
        print(f"  x {skipped_excluded} skipped (excluded, stays hidden)")


if __name__ == "__main__":
    main()
