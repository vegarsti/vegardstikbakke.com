---
title: My open source contributions
slug: open-source
date: "2022-01-10"
---

I've worked full time as a software engineer for around two and a half years, since May 2019.
Almost all of the code I've written has been in closed-source repos, but occasionally I've ran into something somewhere else that I was able to fix.
I compiled a list of the pull requests (PRs) I've submitted that have been accepted upstream.
There are 15 PRs in total, and most of them are typos in documentation, but I have a few code fixes as well.

## Python

A bug fix, bpo-38932, to a test library in the Python standard library, which [was included in Python 3.9](https://docs.python.org/3/whatsnew/changelog.html).
- <https://github.com/python/cpython/pull/17409>

Two fixes to [typeshed](https://github.com/python/typeshed), which is where external type annotations for the Python standard library live.
- <https://github.com/python/typeshed/pull/4191>
- <https://github.com/python/typeshed/pull/4358>

## Others

- [alexbrainman/odbc](https://github.com/alexbrainman/odbc) is an [ODBC](https://github.com/alexbrainman/odbc) driver in Go. It failed to compile on ARM-based Macs due to the Homebrew path for these being different. <https://github.com/alexbrainman/odbc/pull/168>
- Documentation fix to FastAPI, which is an excellent API framework in Python. <https://github.com/tiangolo/fastapi/pull/1015>

## Typos

I'm really good at spotting typos, so I've fixed a number of these.

- <https://github.com/emcconville/wand/pull/464>
- <https://github.com/googleapis/python-storage/pull/84>
- <https://github.com/sourcegraph/about/pull/932>
- <https://github.com/elastic/elasticsearch-py/pull/1277>
- <https://github.com/openethereum/openethereum.github.io/pull/11>
- <https://github.com/benbjohnson/litestream.io/pull/16>
- <https://github.com/benhoyt/benhoyt.github.com/pull/5>
- <https://github.com/lib/pq/pull/1034>
- <https://github.com/ethereum/go-ethereum/pull/22769>
- <https://github.com/ethereum/go-ethereum/pull/22744>
