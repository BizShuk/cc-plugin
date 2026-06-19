# Explore Plugin — Summarization Skills Benchmark

Source: `plugins/explore/references/benchmark.md`
Date: 2026-06-19

## Overview

This guide walks through setting up and running a head-to-head benchmark of the 6 summarization/crawling skills under `plugins/explore/skills/` against a common test page (we use `https://github.com/trending`).

The benchmark measures: speed, major content quality, minor content (chrome) residue, JS/CSS presence in output, and settings/links completeness.

---

## Step 1: Create a Clean Benchmark Workspace

Each skill writes to its own temp directory so outputs never collide:

```bash
SKILLS="content-summarizer firecrawl markitdown playwright-cli scrapling summarize.sh"
for s in $SKILLS; do
  mkdir -p /tmp/skill-test-v2-${s}
done
```

---

## Step 2: Check Tool Availability

Run these pre-flight checks **before** starting the benchmark. If a tool is missing, install it per the instructions below.

### content-summarizer

```bash
# content-summarizer uses markitdown under the hood
markitdown --version
# If missing:
python3 -m venv /tmp/markitdown-venv-v2
/tmp/markitdown-venv-v2/bin/pip install 'markitdown[all]'
export PATH="/tmp/markitdown-venv-v2/bin:$PATH"
```

AI credential needed: **implicit** — the AI agent driving the skill uses its session's LLM (this session). No separate key config.

### firecrawl

```bash
npx -y firecrawl-cli@latest --version   # should print v1.19.13+
npx -y firecrawl-cli@latest --status     # shows credits/auth status
```

AI credential needed: **optional** — `firecrawl scrape` works on a free/guest tier without auth (~1000 credits limit). For full access:

```bash
# Option A: Set env var
export FIRECRAWL_API_KEY=fc-...

# Option B: Browser auth flow (see SKILL.md Path D)
npx -y firecrawl-cli@latest init --all --browser
```

### markitdown

```bash
markitdown --version   # should print 0.1.6+
# If missing (macOS Homebrew + PEP 668 requires a venv):
python3 -m venv /tmp/markitdown-venv
/tmp/markitdown-venv/bin/pip install 'markitdown[all]'
```

AI credential needed: **none** — pure HTTP fetcher, no API keys.

### playwright-cli

```bash
playwright-cli --version              # tries global install
npx --no-install playwright-cli --version  # fallback
# If both fail:
npm install -g @playwright/cli@latest
```

AI credential needed: **none** — uses a local browser. However, this session's Playwright MCP tools can substitute if the binary is unavailable (see Step 5 notes).

### scrapling

```bash
scrapling --version   # should print 0.4.9+
# If missing:
pip install "scrapling[all]>=0.4.8"
scrapling install --force   # fetches Playwright browsers ONCE
```

AI credential needed: **none** — no API keys, no solver credentials (Cloudflare bypass is done via automation). Proxy is optional.

### summarize.sh

```bash
summarize --version   # or: npx -y @steipete/summarize --version
# If missing:
npm install -g @steipete/summarize
```

AI credential needed: **required** — needs `~/.summarize/config.json` with at least `prompt` set and one `apiKeys` entry populated.

Minimal working config (`~/.summarize/config.json`):

```json
{
    "model": { "id": "anthropic/minimax-m3" },
    "apiKeys": {
        "anthropic": "sk-your-key-here"
    },
    "anthropic": {
        "baseUrl": "https://llmbox.bytedance.net/"
    },
    "prompt": "Summarize this page. List every repository with title, full URL, language, stars, description, and any operational details (install command, version requirements, config flags, license, defaults). Write 'not stated' if the page doesn't have it. Do NOT drop any URLs.",
    "output": { "length": "xl", "language": "auto" },
    "cache": { "enabled": false },
    "slides": { "enabled": false }
}
```

> `prompt` must be non-empty — the 2026-06-07 README failure was caused by `"prompt": ""`.
> `slides.sceneThreshold` must be >= 0.1 in summarize v0.19.0+; omit the `slides` block or set it above. Use `prompt` above to force per-repo listings (otherwise the default summary drops all repo URLs).

---

## Step 3: Credential Requirements Summary

| Skill                |         Needs API key?          | How to configure                                                     | Install method                                                      |
| -------------------- | :-----------------------------: | -------------------------------------------------------------------- | ------------------------------------------------------------------- |
| `content-summarizer` |    **implicit** (AI session)    | No config needed; the parent AI session's LLM does the summary       | `pip install 'markitdown[all]'` into a venv                         |
| `firecrawl`          | **optional** (guest tier works) | `export FIRECRAWL_API_KEY=fc-...` or `firecrawl init --browser`      | `npx -y firecrawl-cli@latest` (no install)                          |
| `markitdown`         |            **none**             | —                                                                    | `pip install 'markitdown[all]'` into a venv                         |
| `playwright-cli`     |            **none**             | —                                                                    | `npm install -g @playwright/cli@latest`                             |
| `scrapling`          |            **none**             | —                                                                    | `pip install "scrapling[all]>=0.4.8"` + `scrapling install --force` |
| `summarize.sh`       |          **required**           | `~/.summarize/config.json` — `apiKeys`, `prompt`, provider `baseUrl` | `npm install -g @steipete/summarize`                                |

---

## Step 4: Run the Benchmark — V1 (Full Page, No Constraints)

This captures the "default" behavior — what each skill does out of the box.

### markitdown (2.62s)

```bash
time /tmp/markitdown-venv/bin/markitdown https://github.com/trending \
  -o /tmp/skill-test-markitdown/trending.md
# Output: 63,386 B, 511 lines, 2,344 words — full page: nav + trending + footer
```

### scrapling (2.61s)

```bash
time /Users/bytedance/.venv/bin/scrapling extract get \
  "https://github.com/trending" /tmp/skill-test-scrapling/trending.md \
  --ai-targeted
# Output: 64,294 B, 550 lines, 2,339 words — `--ai-targeted` is for prompt-injection, not chrome-strip
```

### firecrawl (~4s)

```bash
npx -y firecrawl-cli@latest scrape "https://github.com/trending" \
  -o /tmp/skill-test-firecrawl/trending.md
# Output: 84,635 B, 473 lines, 2,300 words — largest file; sub-menu expansion
```

### playwright-cli (~26s)

```bash
playwright-cli open https://github.com/trending
playwright-cli snapshot --filename=/tmp/skill-test-pwcli/snapshot.yml
# Then eval() to extract repo cards from the DOM
playwright-cli close
# Output: 37 KB snapshot + extracted JSON (8 repos)
```

### content-summarizer (~3-4 min)

The content-summarizer uses markitdown as the fetcher, then applies the Index strategy from its SKILL.md. The AI agent:

```bash
markitdown https://github.com/trending -o /tmp/skill-test-content-summarizer/trending.md
# Then classifies as "Index page," picks top 8 repos, writes TL;DR+themes+per-item+business-value
```

### summarize.sh (~14s LLM mode)

```bash
summarize "https://github.com/trending" --model anthropic/minimax-m3 \
  2>&1 | tee /tmp/skill-test-summarize-sh/output.md
# Also try extract-only for comparison:
summarize "https://github.com/trending" --extract --format md \
  > /tmp/skill-test-summarize-sh/trending-raw.md
```

---

## Step 5: Run the Benchmark — V2 (Major Only + Links + Settings)

This re-runs with strict constraints: **only major content, all references/links preserved, functional settings recorded, zero chrome.**

### scrapling (1.94s — `--ai-targeted`)

```bash
time scrapling extract get "https://github.com/trending" \
  /tmp/skill-test-v2-scrapling/trending-cards.md \
  -s "article.Box-row" --ai-targeted
# Output: 13,161 B, 307 lines — zero chrome, 145 links, all 16 repos
```

### firecrawl (1.1s — `--only-main-content --include-tags article`)

```bash
npx -y firecrawl-cli@latest scrape "https://github.com/trending" \
  --include-tags article --only-main-content \
  -o /tmp/skill-test-v2-firecrawl/trending-main.md
# Then strip "Built by" + "Sponsor" lines:
# Output (after post-process): 5,328 B, 110 lines — zero chrome, 48 links
```

### markitdown (1.2s — post-process strip)

```bash
time markitdown https://github.com/trending \
  -o /tmp/skill-test-v2-markitdown/trending-full.md

# Post-process (Python):
python3 << 'PYEOF'
import re
text = open('/tmp/skill-test-v2-markitdown/trending-full.md').read()
m = re.search(r'^## \[google-research[^\n]*timesfm\]', text, re.MULTILINE)
major_start = m.start() if m else 0
m = re.search(r'^## Footer', text, re.MULTILINE)
major_end = m.start() if m else len(text)
major = text[major_start:major_end]
# Strip "Built by" avatar runs
major = re.sub(r'\n\s*\nBuilt by\n(\[!\[.*?\n){1,6}', '', major)
# Strip long language-filter lines
lines = [l for l in major.split('\n') if len(l) < 1000]
open('/tmp/skill-test-v2-markitdown/trending-major.md', 'w').write('\n'.join(lines))
PYEOF
# Output: 5,388 B, 196 lines — zero chrome, 67 links, all 16 repos
```

### playwright-cli (~126s — eval scoped to article.Box-row)

Use the browser to navigate and extract only the 16 cards into structured JSON:

```js
// In the browser's eval():
JSON.stringify(
    Array.from(document.querySelectorAll("article.Box-row")).map((row) => {
        const titleEl = row.querySelector("h2 a");
        return {
            title: titleEl?.textContent?.trim().replace(/\s+/g, " "),
            url: titleEl
                ? "https://github.com" + titleEl.getAttribute("href")
                : null,
            desc: row.querySelector("p")?.textContent?.trim(),
            lang: row
                .querySelector('[itemprop="programmingLanguage"]')
                ?.textContent?.trim(),
            links: Array.from(row.querySelectorAll("a[href]")).map((a) => ({
                text: a.textContent.trim() || null,
                href: a.getAttribute("href")
            }))
        };
    }),
    null,
    2
);
```

Output: 15,849 B JSON — structured, machine-readable, 148 link refs, zero chrome.

### content-summarizer (~4s — markitdown + Index strategy)

Same as markitdown's post-process above, then:

```bash
# After cleaning to trending-major.md, apply the SKILL.md Index strategy:
# 1. Read trending-major.md
# 2. Classify as Index/list page
# 3. Produce: source → TL;DR → themes → top 8 per-item → business-value ideas
# Output summary is inline in the agent's report; the cleaned major-content file
# is /tmp/skill-test-v2-content-summarizer/trending-major.md (4,563 B, 166 lines)
```

### summarize.sh (~18s — LLM with explicit prompt)

```bash
# Use the exact prompt from Step 2 (summarize.sh config) to force per-repo listings:
summarize "https://github.com/trending" --model anthropic/minimax-m3 \
  > /tmp/skill-test-v2-summarize-sh/llm-output-v2.md
# Output: 5,751 B, 113 lines — 16 repos, 16 URLs, 6 with real functional settings
```

---

## Step 6: Measure and Collect Metrics

After all runs complete, collect metrics for each skill:

| Metric                     | How to measure                                                                                         |
| -------------------------- | ------------------------------------------------------------------------------------------------------ | ------- | ------ | ------ | -------- | -------------------------- |
| **Speed** (seconds)        | Use `time` before the command; record `real`                                                           |
| **Major content size**     | `wc -clw /tmp/skill-test-v2-<skill>/trending-major.md` (or equivalent)                                 |
| **Minor content (chrome)** | `wc -clw /tmp/skill-test-<skill>/trending.md` (v1 full-page) and diff                                  |
| **JS/CSS in output**       | `grep -ciE '(javascript:                                                                               | <script | <style | \.js\b | window\. | addEventListener)' <file>` |
| **Links preserved**        | `grep -oP '\[.\*?\]\(https?://[^\)]+\)' <file>                                                         | wc -l`  |
| **Repos captured**         | Count `## [owner` headings or equivalent                                                               |
| **Functional settings**    | Manual — read the summary; count repos with `install`, `version`, `config`, `license`, `requires` info |

---

## Step 7: Compare Against This Benchmark

Run the same test on a different URL to compare results. We used `https://github.com/trending`. Other good candidates:

| Page type       | URL                                                                                                        | Notes                                        |
| --------------- | ---------------------------------------------------------------------------------------------------------- | -------------------------------------------- |
| Blog article    | `https://platform.claude.com/docs/en/build-with-claude/prompt-engineering/claude-prompting-best-practices` | The original 2026-06-07 test page            |
| JS-heavy SPA    | `https://twitter.com/explore`                                                                              | Requires browser — test which skills survive |
| Index/list page | `https://news.ycombinator.com`                                                                             | Another trending-style page for comparison   |
| PDF             | any PDF URL                                                                                                | Tests markitdown's PDF → markdown pipeline   |

---

## Bench mark Results — 2026-06-19

### Run 1: Default Behavior (Full Page)

| Skill                | Tool                            |  Time |         Output | Lines |              Chrome?               |
| -------------------- | ------------------------------- | ----: | -------------: | ----: | :--------------------------------: |
| `markitdown`         | markitdown 0.1.6                |  2.6s |       63,386 B |   511 | yes (nav + footer + language list) |
| `scrapling`          | scrapling 0.4.9                 |  2.6s |       64,294 B |   550 | yes (--ai-targeted collapsed some) |
| `firecrawl`          | firecrawl v1.19.13              |   ~4s |       84,635 B |   473 |  yes (sub-menus expanded, worst)   |
| `playwright-cli`     | MCP playwright (binary missing) |   26s | 37 KB snapshot |     — |     yes then stripped via eval     |
| `content-summarizer` | markitdown + LLM                | ~3-4m | inline summary |     — |  no (LLM breaks chrome naturally)  |
| `summarize.sh`       | summarize v0.19.0 + minimax-m3  | 14.4s |        3,876 B |    60 |   no (LLM) / yes (raw --extract)   |

### Run 2: Major-Only + Links + Settings

| Skill                | Tool                                             |     Time |      Major size | Chrome stripped |   Links | JS/CSS | Settings |
| -------------------- | ------------------------------------------------ | -------: | --------------: | :-------------: | ------: | :----: | :------: |
| `scrapling`          | `-s "article.Box-row" --ai-targeted`             | **1.9s** |        13,161 B |        0        | **145** |   ✗    |   0/16   |
| `firecrawl`          | `--include-tags article --only-main-content`     | **1.1s** |         5,328 B |        0        |      48 |   ✗    |   0/16   |
| `markitdown`         | fetch + post-process (slice+regex)               | **1.2s** |         5,388 B |    345 lines    |      67 |   ✗    |   0/16   |
| `content-summarizer` | markitdown + post-process + Index strategy       |   **4s** |        12,608 B |    345 lines    |      52 |   ✗    |   0/16   |
| `summarize.sh`       | `--model anthropic/minimax-m3` + explicit prompt |  **18s** |         5,751 B |        0        |      16 |   ✗    | **6/16** |
| `playwright-cli`     | MCP playwright + eval()                          |     126s | 15,849 B (JSON) |        0        |     148 |   ✗    |   0/16   |

#### Notes on the results

- **Chrome stripped column**: `0` means the tool's flag/eval approach removed chrome natively — no manual work. A number means manual post-process lines needed.
- **JS/CSS**: All 6 skills produce zero JS/CSS in their v2 major-only output. GitHub Trending is server-rendered. For JS-heavy SPAs, only `playwright-cli` and `scrapling fetch`/`stealthy-fetch` would work.
- **Settings column**: The GitHub Trending page card-level content does NOT expose install commands, version requirements, config flags, or licenses. Only `summarize.sh` (via LLM inference from repo descriptions) extracted any operational notes (6 of 16 repos — e.g., "Single static binary, zero dependencies").
- **Firecrawl's SKILL.md mentions `firecrawl ask`** — but the CLI v1.19.13 does not expose this subcommand. This is a doc/CLI drift.
- **summarize.sh's default prompt drops ALL repo URLs** — only the explicit `--prompt` (shown in Step 2) forces per-repo listings with URLs. The bundled `config.sample.json` has `"prompt": ""` and empty `apiKeys[]` — both must be filled.

### Sample Output Files

All sample files are stored under `plugins/explore/references/benchmark.sample/`.

**V2 (major-only):**

| Skill                | Sample file                                                                                                   | Size     | Content                                                                   |
| -------------------- | ------------------------------------------------------------------------------------------------------------- | -------- | ------------------------------------------------------------------------- |
| `content-summarizer` | [`v2/content-summarizer-summary.md`](benchmark.sample/v2/content-summarizer-summary.md)                       | 12,608 B | Full Index strategy: TL;DR + themes + 8 per-item + 3 business-value ideas |
| `content-summarizer` | [`v2/content-summarizer-trending-major.md`](benchmark.sample/v2/content-summarizer-trending-major.md)         | 4,563 B  | Cleaned raw cards (markitdown → post-processed)                           |
| `firecrawl`          | [`v2/firecrawl-trending-major.md`](benchmark.sample/v2/firecrawl-trending-major.md)                          | 5,328 B  | Cleaned raw cards (`--only-main-content`)                                 |
| `firecrawl`          | [`v2/firecrawl-trending-main.md`](benchmark.sample/v2/firecrawl-trending-main.md)                            | 14,974 B | Raw scrape before post-process                                            |
| `markitdown`         | [`v2/markitdown-trending-major.md`](benchmark.sample/v2/markitdown-trending-major.md)                        | 5,388 B  | Cleaned raw cards (post-processed)                                        |
| `markitdown`         | [`v2/markitdown-trending-full.md`](benchmark.sample/v2/markitdown-trending-full.md)                          | 63,386 B | Raw full-page before post-process                                         |
| `playwright-cli`     | [`v2/playwright-cli-trending.json`](benchmark.sample/v2/playwright-cli-trending.json)                        | 15,849 B | Structured JSON, 148 link refs                                            |
| `scrapling`          | [`v2/scrapling-trending-cards.md`](benchmark.sample/v2/scrapling-trending-cards.md)                          | 13,161 B | Raw cards, zero chrome, 145 links (`--ai-targeted`)                       |
| `scrapling`          | [`v2/scrapling-trending-cards-raw.md`](benchmark.sample/v2/scrapling-trending-cards-raw.md)                  | 13,177 B | Raw cards without `--ai-targeted`                                         |
| `summarize.sh`       | [`v2/summarize.sh-summary.md`](benchmark.sample/v2/summarize.sh-summary.md)                                   | 5,751 B  | LLM-formatted, 16 repos, 16 URLs, 6 with settings                         |
| `summarize.sh`       | [`v2/summarize.sh-llm-output-v2.md`](benchmark.sample/v2/summarize.sh-llm-output-v2.md)                      | 5,751 B  | LLM output with explicit prompt                                           |
| `summarize.sh`       | [`v2/summarize.sh-llm-output.md`](benchmark.sample/v2/summarize.sh-llm-output.md)                            | 3,699 B  | LLM output with default prompt (URLs dropped!)                            |
| `summarize.sh`       | [`v2/summarize.sh-raw-extract.md`](benchmark.sample/v2/summarize.sh-raw-extract.md)                          | 6,149 B  | Raw `--extract --format md` (no LLM)                                      |

**V1 (full-page, for before/after comparison):**

| Skill                | Sample file                                                                                    | Size     | Content                                                             |
| -------------------- | ---------------------------------------------------------------------------------------------- | -------- | ------------------------------------------------------------------- |
| `firecrawl`          | [`v1/firecrawl.md`](benchmark.sample/v1/firecrawl.md)                                          | 84,635 B | Full page: nav sub-menus + trending + footer                        |
| `markitdown`         | [`v1/markitdown.md`](benchmark.sample/v1/markitdown.md)                                        | 63,386 B | Full page: nav + language list + trending + footer                  |
| `scrapling`          | [`v1/scrapling.md`](benchmark.sample/v1/scrapling.md)                                          | 64,294 B | Full page: `--ai-targeted` (prompt-injection only)                  |
| `content-summarizer` | [`v1/content-summarizer.md`](benchmark.sample/v1/content-summarizer.md)                        | 63,386 B | Raw markitdown output used as input to Index strategy               |
| `summarize.sh`       | [`v1/summarize.sh-output.md`](benchmark.sample/v1/summarize.sh-output.md)                      | 3,876 B  | Default prompt LLM summary                                          |
| `summarize.sh`       | [`v1/summarize.sh-trending-raw.md`](benchmark.sample/v1/summarize.sh-trending-raw.md)          | 6,149 B  | Raw `--extract --format md`                                         |

### V1 → V2 Reduction (Chrome Stripped)

| Skill          | Full page (v1) | Chrome stripped to (v2) | Reduction |
| -------------- | -------------: | ----------------------: | --------: |
| `firecrawl`    |      84,635 B  |                 5,328 B | **93.7%** |
| `markitdown`   |      63,386 B  |                 5,388 B | **91.5%** |
| `scrapling`    |      64,294 B  |                13,161 B | **79.5%** |
| `summarize.sh` |       3,876 B  |                 5,751 B | n/a (LLM, different output style) |

---

## Top 16 Repos (as of 2026-06-19)

All 6 skills agree on the same 16 trending repos. `DeusData/codebase-memory-mcp` led with 2,322 stars/day.

|  #  | Repo                                                                       |    Lang    | Stars Today | Cumulative |  Forks |
| :-: | -------------------------------------------------------------------------- | :--------: | ----------: | ---------: | -----: |
|  1  | `DeusData/codebase-memory-mcp`                                             |     C      |       2,322 |      7,319 |    576 |
|  2  | `obra/superpowers`                                                         |   Shell    |       1,429 |    232,675 | 20,664 |
|  3  | `Kilo-Org/kilocode`                                                        | TypeScript |       1,345 |     22,412 |  2,715 |
|  4  | `google-research/timesfm`                                                  |   Python   |         844 |     23,546 |  2,235 |
|  5  | `makeplane/plane`                                                          | TypeScript |         613 |     51,936 |  4,610 |
|  6  | `freeCodeCamp/freeCodeCamp`                                                | TypeScript |         417 |    449,658 | 45,143 |
|  7  | `n0-computer/iroh`                                                         |    Rust    |         369 |     10,081 |    464 |
|  8  | `alibaba/zvec`                                                             |    C++     |         259 |     11,353 |    658 |
|  9  | `Universal-Debloater-Alliance/universal-android-debloater-next-generation` |    Rust    |         244 |      7,979 |      — |
| 10  | `zai-org/GLM-5`                                                            |     —      |         202 |      4,270 |    437 |
| 11  | `withastro/flue`                                                           | TypeScript |         162 |      5,583 |    308 |
| 12  | `Kong/insomnia`                                                            | TypeScript |          18 |     38,756 |      — |
| 13  | `dotnet/aspnetcore`                                                        |     C#     |          14 |     38,111 |      — |
| 14  | `owainlewis/awesome-artificial-intelligence`                               |     —      |          40 |     14,498 |      — |
| 15  | `Lightricks/LTX-2`                                                         |   Python   |          51 |      7,538 |      — |
| 16  | `LibreTranslate/LibreTranslate`                                            |   Python   |          51 |     15,055 |      — |

---

## See Also

- Plugin README: `plugins/explore/README.md` — skill comparison table (2026-06-07)
- Plugin manifest: `plugins/explore/.claude-plugin/plugin.json`
- SKILL.md files: `plugins/explore/skills/*/SKILL.md` — per-skill documentation
- Test sample directory: `ls /tmp/skill-test-v2-*` — all sample output files from this benchmark
