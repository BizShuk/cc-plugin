# content-summarizer — Step 1: Capability Check & Fetching

Never claim to have read something that was not actually retrieved.

## Fetching with `markitdown`

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

## When `markitdown` is not available

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

## When a capability is missing

First confirm what you already have. The most common missing piece is
`markitdown` itself — suggest:

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

When you do invoke a known tool, use its fully qualified name
(`ServerName:tool_name`).
