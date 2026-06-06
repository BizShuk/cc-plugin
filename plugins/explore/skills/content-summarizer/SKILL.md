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
CLI that converts files and URLs to clean Markdown.

## When to Use

- User shares a URL → summarize
- User shares a file (PDF, DOCX, etc.) → summarize
- User shares a video link → transcript required; summarize transcript
- User shares an index/list page → recursive item summaries

## Know your setup (before acting)

Check what tools are available before acting. The preferred content extraction
tool is `markitdown` — verify it exists by running `markitdown --version`.

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
and adapt to what is available. When something is missing, see "When a
capability is missing".

## Workflow

Copy this checklist and track progress:

```md
- [ ] Step 1: Confirm the content is fetchable (capability check)
- [ ] Step 2: Strip side content (menus, footers, nav, ads)
- [ ] Step 3: Classify the content shape (Index / Article / Video / Document)
- [ ] Step 4: Summarize using the matching strategy
- [ ] Step 5: Extract 1-3 business-value / opportunity ideas (brainstorm)
- [ ] Step 6: Assemble the output; offer calendar / notes follow-ups
```

## Step 1 — Capability check (do this FIRST)

Never claim to have read something that was not actually retrieved.

### Fetching with `markitdown`

`markitdown` is the preferred tool for both URLs and local files:

```bash
# URL → Markdown
markitdown https://example.com/article

# Local file → Markdown
markitdown path/to/file.pdf

# Pipe content
cat file.docx | markitdown

# Save output
markitdown input.pptx -o output.md
```

After running `markitdown`, inspect the output:

- Substantial readable text returned → proceed.
- Error, login/paywall page, cookie-consent shell, near-empty body, or
  obvious JavaScript-only page → `STOP and warn`. Do not guess or fabricate.
- For YouTube URLs, `markitdown` extracts the transcript directly — no
  separate transcript tool needed (requires `youtube-transcription` extra).

### When `markitdown` is not available

Fall back to any available web-fetching capability. If none exist:

1. Ask the user to paste the text directly into the chat.
2. Provide a clean-reader version of the URL.
3. Suggest installing markitdown: `pip install 'markitdown[all]'`.
4. For video: request the transcript or captions.

`Known hard cases — warn early, before attempting a full summary:`

- X/Twitter, Instagram, Facebook, private LinkedIn posts: usually blocked or login-gated.
- Paywalled news and members-only articles.
- Google Docs / Notion / files requiring sign-in.
- YouTube without the `youtube-transcription` extra installed.

## Step 2 — Strip side content

Raw fetched content typically includes page chrome that dilutes the summary:
navigation menus, site headers, footers, sidebars, cookie banners, share
buttons, related-article blocks, comment sections, and ad placeholders.

Remove these `before` classification and summarization — they waste tokens
and risk polluting the summary with irrelevant text.

### Removal checklist

Scan the fetched Markdown and strip blocks that match these patterns:

| Pattern | Examples |
| ------- | -------- |
| `Navigation / menus` | Top nav bars, hamburger menus, breadcrumbs, site-wide link lists |
| `Site header / branding` | Logo blocks, search bars, login/signup links |
| `Footer` | Copyright notices, site maps, "About Us" / "Contact" / "Terms" link clusters |
| `Sidebar` | Tag clouds, category lists, "Popular posts", newsletter signup forms |
| `Social / sharing` | Share buttons, follow links, social embeds |
| `Ads / promotions` | Banner ads, sponsored content blocks, "You may also like" |
| `Cookie / consent` | Cookie banners, GDPR consent text |
| `Comments` | User comment sections, "Leave a reply" |
| `Repeated boilerplate` | Identical blocks appearing at top and bottom (e.g. site tagline) |

### How to strip

1. `Tool-level filtering (preferred)` — if the fetching tool supports it,
   filter at fetch time:
   - `scrapling`: use `--ai-targeted` flag (auto-removes ads and non-main
     content) or `--css-selector "main"` / `--css-selector "article"` to
     extract only the primary content container.
   - `markitdown`: does not have built-in filtering — proceed to manual
     removal.
2. `Manual removal` — after fetching, scan the Markdown output and delete
   obvious noise sections. Look for:
   - Dense clusters of short links at the very top or bottom of the document.
   - Repeated separator patterns (`---`, `***`) flanking non-content blocks.
   - Sections whose headings are generic site-chrome labels (`Menu`,
     `Navigation`, `Footer`, `Related Posts`, `Comments`).
3. `CSS selector pre-filtering` — when re-fetching is cheap, try fetching
   with a targeted selector (`main`, `article`, `#content`, `.post-body`)
   to grab only the primary content container.

### When NOT to strip

- `Index / list pages`: the navigation links ARE the content — do not strip
  them. Only strip site-level chrome (header, footer), not the item list.
- `Documentation / API reference`: sidebars with table-of-contents or
  parameter listings are content, not noise.

After cleaning, the remaining text should be predominantly the page's
primary content. If the result is near-empty, the page may be JavaScript-
rendered — note this and consider a browser-based fetcher.

## Step 3 — Classify the content shape

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

## Step 4 — Summarize

### Index strategy (recursive summary)

The point is NOT to relist the page. Dig into each item.

1. Extract the items (title + link).
2. Default to the `top 5-8` by prominence/relevance. If more exist, summarize
   these, state how many remain, and offer to continue. Confirm before fetching
   a large batch.
3. For each item: fetch it (re-run the Step 1 check per item) and write
   `2-4 sentences` — what it is and why it matters. If the item is technical,
   also flag any standout configuration option, setting, or operational caveat.
   Skip items that can't be fetched and note which were skipped.
4. Go `one level deep only` (index → item). Do not recurse into items-of-items
   unless asked.
5. Open with a short `roll-up`: the 2-3 themes or patterns across the items.

### Article strategy (concise)

Optimize for speed and signal. One fetch, no deep recursion.

- `TL;DR`: 1-2 sentences.
- `Key points`: 3-6 bullets, each a complete thought.
- `Pros / Cons`: include ONLY when the piece is evaluative (product, proposal,
  argument, recommendation). If it is purely informational, skip pros/cons —
  do not force it.
- `Configuration & fine print`: when the content includes technical or
  operational detail, list the easy-to-miss specifics — special configuration
  options, key settings or flags, defaults, prerequisites, limits, and version-
  or platform-specific caveats, plus any non-obvious steps. Capture exact names
  and values (a parameter, a flag, a limit), not vague paraphrases. Include
  only when such details exist.
- Do not pad the factual summary.

## Step 5 — Business value (1-3 ideas)

After the summary, add `1-3` concrete business-value or opportunity ideas drawn
from the content. This is the high-leverage part: move beyond "what it says" to
"so what — how could this matter".

- Produce `1-3` items, calibrated to how much signal the content offers. One
  sharp idea beats three generic ones.
- Make each idea `specific and tied to this content`. Bad: "could improve
  efficiency". Good: name the angle, who it helps, and the concrete move.
- Brainstorm across angles, picking what fits: a product/feature idea, a market
  or customer opportunity, a competitive/strategic implication, a process or
  cost improvement, a risk to watch, or a concrete next action.
- `Label these as your own inference, not the source.` They are extrapolations —
  keep them clearly separate from the faithful summary so they are never
  mistaken for claims the source made.
- Tailor to the user's context when known (role, company, goals); otherwise
  keep ideas broadly applicable.
- If the content has no plausible business angle, say so briefly instead of
  forcing ideas.

## Output format

Lead with the source so the summary is always traceable. Keep formatting light.
Match length to type: Articles short; Index pages as long as the item count requires.

```md
Source: [title](url) ·or· file: filename.pdf

TL;DR: ...

Key points

- ...

(Index only) By item

- Item title — 2-4 sentence summary

(if evaluative) Pros / Cons

-   - ...
-   - ...

(if technical / operational) Configuration & fine print

- setting / flag / option — what it does, default, caveat

Business value (ideas / inference, not from the source)

1. ...
2. ...
   (1-3 items)

(if any) Dates & action items

- ...
```

## Step 6 — Calendar & Notes follow-ups

Offer these whenever they apply; also perform them on direct request. Always
embed the source link or file reference in whatever gets saved.

`Calendar` — if the content contains dated events (webinar, deadline, meetup,
launch, conference session):

- Use your available calendar capability, preferring `Apple Calendar` when an
  Apple calendar integration is present. If you have no calendar capability,
  say so.
- Always confirm date, time, and timezone before creating. Put the source URL
  in the event notes.

`Notes` — to save the summary itself:

- Use your available notes capability, preferring `Apple Notes` when an Apple
  Notes integration is present.
- Title the note clearly (source title + date) and place the source link / file
  name at the top.
- If you have no notes capability, say so and fall back to: returning the
  summary as a clearly formatted, copy-ready block, or saving it as a file.
  Never silently drop the request.

## When a capability is missing

First confirm what you already have (see "Know your setup"). The most common
missing piece is `markitdown` itself — suggest:

```bash
pip install 'markitdown[all]'
```

Or install only the extras you need:

```bash
pip install 'markitdown[pdf,docx,pptx,xlsx,youtube-transcription]'
```

For other missing capabilities:

- `JavaScript-heavy / SPA pages`: a browser-automation capability (e.g.
  Playwright) — `markitdown` fetches static HTML only.
- `Login-gated or interactive pages`: a browser-automation capability.
- `Apple Notes / Calendar` (e.g. on a desktop client): a connector or local
  integration that exposes those apps.
- `Search-first cases`: a web-search capability, when the source must be found
  before it can be summarized.
- `OCR on embedded images in documents`: install `markitdown-ocr` plugin and
  use `markitdown --use-plugins` with an LLM client configured.

When you do invoke a known tool, use its fully qualified name (`ServerName:tool_name`).

## Common Mistakes

| Mistake                                                    | Fix                                     |
| ---------------------------------------------------------- | --------------------------------------- |
| Summarizing URL content without fetching                   | Always fetch first; warn if unreachable |
| Generating plausible-sounding summaries from URL structure | `STOP and warn` — never fabricate       |
| Applying Index strategy to an article                      | Quick test: continuous text = Article   |
| Forcing pros/cons on purely informational content          | Only use when piece is evaluative       |
| Recursing into items-of-items on index pages               | One level deep only, unless asked       |
| Inventing calendar/notes entries without confirmation      | Always confirm date/time/timezone first |

## Red Flags — STOP

If you catch yourself doing any of these, stop and correct:

- About to summarize without fetching (fabrication risk)
- URL looks like an article, you're treating it as fetchable without verifying
- Writing pros/cons for a how-to guide
- Digging three levels deep on an index page
- Creating calendar event without confirming date/time

## Principles

- `Know your setup`: discover available capabilities; don't assume tool names.
- `One source of truth`: every summary names its source.
- `Separate fact from idea`: the summary stays faithful to the source;
  business-value ideas are clearly labeled inference.
- `Don't fabricate`: if retrieval failed, an honest "couldn't fetch this" is
  the correct answer.
- `Respect the budget`: cap recursion, confirm before large batches, keep
  Articles short.
