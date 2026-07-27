---
name: content-summarizer
description: >
    Use when the user shares a URL, file, PDF, social post, or video and asks to
    summarize, extract key points, get a TL;DR, brief them, or "tell me what this
    says". Also use when the user shares a list/index page (trending projects, event
    line-ups, link round-ups, top-N lists) and wants an overview of its items.
    Triggers on: "summarize this", "TL;DR", "brief me", "what does this say",
    "key points from", "read this for me".
---

# Content Summarizer

Turn a web page, document, social post, or video into key points, plus a few
business-value ideas. Every summary must be traceable to a source (URL or file name).

Primary content extraction tool: `markitdown` (microsoft/markitdown) — a Python
CLI that converts files and URLs to clean Markdown. Verify it exists with
`markitdown --version` before acting.

## When to Use

- User shares a URL → summarize
- User shares a file (PDF, DOCX, etc.) → summarize
- User shares a video link → transcript required; summarize transcript
- User shares an index/list page → recursive item summaries

`markitdown` natively handles:

| Category  | Formats                                                  |
| --------- | -------------------------------------------------------- |
| Web       | URLs (fetches and converts HTML to Markdown)             |
| Documents | PDF, DOCX, PPTX, XLSX, XLS, EPUB                         |
| Data      | CSV, JSON, XML                                           |
| Media     | Images (EXIF + OCR), Audio (EXIF + speech transcription) |
| Video     | YouTube URLs (transcript extraction)                     |
| Archives  | ZIP (iterates over contents)                             |
| Other     | Outlook messages (.msg), plain HTML files                |

Also check for calendar and notes capabilities (Apple Calendar, Apple Notes)
and adapt to what is available.

## Workflow

Copy this checklist and track progress:

```md
- [ ] Step 1: Confirm the content is fetchable (capability check)
- [ ] Step 2: Strip side content (menus, footers, nav, ads)
- [ ] Step 3: Classify the content shape (Index / Article / Video / Document / Links)
- [ ] Step 4: Summarize using the matching strategy
- [ ] Step 5: Extract 1-3 business-value / opportunity ideas (brainstorm)
- [ ] Step 6: Assemble the output; offer calendar / notes follow-ups
```

| Step | Detail |
| --- | --- |
| 1 — fetch, verify, fall back | [references/fetching.md](references/fetching.md) |
| 2 — strip page chrome | [references/cleaning.md](references/cleaning.md) |
| 3 — classify | below |
| 4-6 — strategies, business value, output template | [references/output-format.md](references/output-format.md) |

## Step 3 — Classify the content shape

This choice drives everything downstream, so make it explicitly.

- `Index / list page` → a hub whose value is its links, not its own prose
  (trending repos, "top 10" lists, conference agendas, newsletter round-ups,
  search results). Use the `Index strategy`.
- `Single article / post` → one self-contained piece (blog post, news story,
  essay, product page, a single social post or thread). Use the `Article strategy`.
- `Video` → apply the `Article strategy` to the transcript, but add timestamps
  for key moments. If using `markitdown` on a YouTube URL, the transcript is
  extracted automatically.
- `Document (PDF/DOCX/PPTX/XLSX/etc.)` → run through `markitdown` first to
  get Markdown, then treat as Article unless it is clearly a directory of
  separate items (then treat as Index).
- `Data file (CSV/JSON/XML)` → `markitdown` converts these to readable
  Markdown tables/structure. Treat as Article.

Quick test: many short titled links pointing elsewhere = Index; continuous body text = Article.

## Red Flags — STOP

If you catch yourself doing any of these, stop and correct:

- About to summarize without fetching (fabrication risk)
- Treating a URL as fetchable without verifying the fetch actually returned text
- Writing pros/cons for a how-to guide
- Digging three levels deep on an index page
- Creating a calendar event without confirming date/time/timezone

## Principles

- `Know your setup`: discover available capabilities; don't assume tool names.
- `One source of truth`: every summary names its source.
- `Separate fact from idea`: the summary stays faithful to the source;
  business-value ideas are clearly labeled inference.
- `Don't fabricate`: if retrieval failed, an honest "couldn't fetch this" is
  the correct answer.
- `Respect the budget`: cap recursion, confirm before large batches, keep
  Articles short.

## Related

- `[[markitdown]]` the underlying conversion CLI, when only the Markdown is wanted
- `[[apple-calendar]]` / `[[apple-notes]]` for the Step 6 follow-ups
