# content-summarizer — Step 2: Strip Side Content

Raw fetched content typically includes page chrome that dilutes the summary:
navigation menus, site headers, footers, sidebars, cookie banners, share
buttons, related-article blocks, comment sections, and ad placeholders.

Remove these `before` classification and summarization — they waste tokens
and risk polluting the summary with irrelevant text.

## Removal checklist

Scan the fetched Markdown and strip blocks that match these patterns:

| Pattern                  | Examples                                                                     |
| ------------------------ | ---------------------------------------------------------------------------- |
| `Navigation / menus`     | Top nav bars, hamburger menus, breadcrumbs, site-wide link lists             |
| `Site header / branding` | Logo blocks, search bars, login/signup links                                 |
| `Footer`                 | Copyright notices, site maps, "About Us" / "Contact" / "Terms" link clusters |
| `Sidebar`                | Tag clouds, category lists, "Popular posts", newsletter signup forms         |
| `Social / sharing`       | Share buttons, follow links, social embeds                                   |
| `Ads / promotions`       | Banner ads, sponsored content blocks, "You may also like"                    |
| `Cookie / consent`       | Cookie banners, GDPR consent text                                            |
| `Comments`               | User comment sections, "Leave a reply"                                       |
| `Repeated boilerplate`   | Identical blocks appearing at top and bottom (e.g. site tagline)             |

## How to strip

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

## When NOT to strip

- `Index / list pages`: the navigation links ARE the content — do not strip
  them. Only strip site-level chrome (header, footer), not the item list.
- `Documentation / API reference`: sidebars with table-of-contents or
  parameter listings are content, not noise.

After cleaning, the remaining text should be predominantly the page's
primary content. If the result is near-empty, the page may be JavaScript-
rendered — note this and consider a browser-based fetcher.
