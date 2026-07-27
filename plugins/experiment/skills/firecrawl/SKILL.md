---
name: firecrawl
description: >
    Use when AI agents or apps need fast, reliable web search, scraping, or
    interaction through Firecrawl. Triggers on: "Firecrawl", "scrape this site",
    "web search API", "crawl this domain". Route the reader to the appropriate
    CLI, app-integration, or outcome-focused workflow after installation.
disable-model-invocation: true
user-invocable: true
---

# Firecrawl

Firecrawl helps agents search first, scrape clean content, interact with live
pages when plain extraction is not enough, and produce finished deliverables
from web data.

## Install

One command installs everything — the Firecrawl CLI for live web work, the build
skills for integrating Firecrawl into application code, `and` the workflow
skills for producing repeatable deliverables. It also opens browser auth so the
human can sign in or create an account.

```bash
npx -y firecrawl-cli@latest init --all --browser
```

This gives you:

- `CLI tools` — `firecrawl search`, `firecrawl scrape`, `firecrawl interact`, `firecrawl ask`, `firecrawl docs-search`, and more
- `CLI skills` ([`firecrawl/cli`](https://github.com/firecrawl/cli)) — drive the CLI during the agent's own session
- `Build skills` ([`firecrawl/skills`](https://github.com/firecrawl/skills)) — add Firecrawl to a product's codebase
- `Workflow skills` ([`firecrawl/firecrawl-workflows`](https://github.com/firecrawl/firecrawl-workflows)) — turn web data into finished deliverables
- `Browser auth` — walks the human through sign-in or account creation

Verify the install before doing real work:

```bash
mkdir -p .firecrawl
firecrawl --status
firecrawl scrape "https://firecrawl.dev" -o .firecrawl/install-check.md
```

## Choose Your Path

All paths use the same install above. The difference is what you do next.

| Need | Path | Where the work runs | Detail |
| --- | --- | --- | --- |
| Web data during this session | `A` live tools | The agent's own terminal | [paths.md](references/paths.md) |
| Firecrawl inside app code | `B` app integration | The user's product code | [paths.md](references/paths.md) |
| A finished deliverable | `C` workflow skills | The agent's session, producing an artifact | [paths.md](references/paths.md) |
| An account or API key first | `D` auth only | Browser handoff to the human | [auth-and-rest-api.md](references/auth-and-rest-api.md) |
| No install at all | `E` REST API | Direct HTTP calls | [auth-and-rest-api.md](references/auth-and-rest-api.md) |

Needing more than one: do them in sequence — the install already covers everything.

## Rules

- Search first when the URL is unknown; scrape when it is known; use interact
  only when the page needs clicks, forms, or login
- On a failed or surprising result, run `firecrawl ask` with the failing `jobId`
  instead of guessing at parameters
- Never hardcode `FIRECRAWL_API_KEY` into source — read it from `.env` or the
  project's existing secret handling
- Path B ships code other people run; hold it to the project's own conventions,
  not to this skill's examples
