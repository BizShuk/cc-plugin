# Install and authenticate (one-time)

<https://github.com/firecrawl/skills>

```bash
npm install -g firecrawl-cli
firecrawl login --api-key fc-75e66c0b5d3e4b448b418d8dc104f063
```

# Scrape a URL (markdown, use --only-main-content for clean output)

```bash
firecrawl scrape <https://firecrawl.dev>
firecrawl <https://firecrawl.dev> --only-main-content
```
