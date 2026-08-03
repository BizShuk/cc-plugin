---
name: project-route
description: >
    Use to route inbound content/files to the correct project under ~/projects/
    and load that project's context. Given a receipt, a stock CSV, a log dump, a
    note, or any artifact whose home project is unclear, this skill resolves it to
    exactly one project via the tag index, then loads that project's README.md and
    CLAUDE.md so the work follows the project's own conventions. Triggers on:
    "which project does this go to", "route this file", "把這個歸到哪個專案",
    "where should this live", "find the project for", or any cross-project
    filing/dispatch decision. Other agents (e.g. hermes) should invoke this FIRST,
    before moving or creating any file.
version: "1.1.0"
allowed-tools: Bash, Read
metadata:
    type: reference
---

# project-route — route content to a ~/projects/ project

Resolve a piece of content to exactly one project, then load that project's own
docs so downstream work obeys the project's conventions. This skill decides
`where`; it does not move files.

## Layout — two levels

`~/projects/` holds projects at two depths:

```tree
~/projects/
├── <project>/            # 專案 (Project) at root
└── <category>/           # 分類 (Category) — a container, not a project
    └── <project>/        # 專案 inside a category
```

- Current categories (all 13): `ai`, `collections`, `data`, `env_setup`,
  `game`, `iphone`, `platform`, `playground`, `product`, `research`, `social`,
  `tools`, `web`. The live list is `router.py categories`; if it disagrees
  with this line, the index is authoritative — rebuild and trust it.
- Some categories are `hybrids`: `collections`/`social`/`web` carry their own
  docs or code and can themselves be routing destinations; `env_setup` has its
  own `go.mod` yet also holds projects. Never infer container-ness from the
  name; check the index.
- A pure-container category (e.g. `game/`) may carry an umbrella `README.md`
  but is never a routing destination itself — route to a project inside it.
  Hybrid categories above are the exception.
- Category depth is fixed at one level; `<category>/<category>/<project>` does
  not exist.
- Project names are unique across the whole tree, so the bare name still
  resolves. `router.py cd <name>` and `router.py cd <category>/<name>` both
  work.
- A project symlinked into a category (e.g. `ai/cc-plugin` → `cc-plugin`) is an
  `alias`. The index keeps the physical path; always file to the physical path.

## Config (single source of truth)

- Index: `~/projects/.project_index/projects.json` — one combined JSON.
    - Top-level: `generated_at` (ISO timestamp), `stale_after_days` (default 3),
      `projects_root`, `categories`, `projects`.
    - Per project: `path`, `category` (`null` when at root), `aliases`, `tags`,
      `purpose`, `stack`, `subprojects`, `has_readme`, `has_claude`,
      `is_category_root`.
    - Per subproject: `name`, `purpose` (extracted from its own `README.md`),
      `has_readme`, `has_claude`.
- Router CLI: `~/projects/.project_index/router.py` — reads the same index.
- `tags` is the routing signal (bilingual, lowercase). `tags` and `categories`
  are hand-curated and are preserved across `router.py rebuild`; the auto-scan
  never overwrites them.

## Staleness — auto-rebuild

The router **auto-rebuilds** when the index's `generated_at` is older than
`stale_after_days` (default 3 days). Every `search` / `list` / `show` / `cd` /
`stack` / `category` call checks the timestamp first; if stale, it re-scans both
levels of `~/projects/` (reading each project's `README.md` + `CLAUDE.md` AND
every direct subfolder's `README.md` + `CLAUDE.md` to extract subproject
purpose), then refreshes `projects.json` and `INDEX.md`. Pass
`--no-auto-rebuild` to skip the check (use the stale index as-is).

New categories are auto-detected: a top-level dir with no own `README.md` /
`CLAUDE.md` and `2+` project-like children becomes a category. A container that
does have its own docs (like `game/`) must be listed by hand in `categories`.

`Do not take the auto-rebuild on trust` — it only fires when the router is
actually invoked, so an index nobody has queried for weeks stays stale and fails
silently. Check the timestamp yourself before relying on a routing answer:

```bash
python3 -c "import json,datetime;d=json.load(open('$HOME/projects/.project_index/projects.json'));g=datetime.datetime.fromisoformat(d['generated_at']);print(g, (datetime.datetime.now()-g).days, 'days old')"
```

If it is older than `stale_after_days` (default 3), run
`python3 ~/projects/.project_index/router.py rebuild` before routing. If the
rebuild cannot run (index missing, not writable, `router.py` absent), say so
in the routing answer and mark the result as based on a stale index — never
present a stale match as if it were current.

## Procedure

1. Extract routing keywords from the content (vendor, currency, domain words,
   filenames, language). Keep both Chinese and English forms — tags carry both.

2. Resolve the project. Prefer the CLI:

    ```bash
    python3 ~/projects/.project_index/router.py search "<keyword>"
    ```

    The CLI scores `name` (10) > `tags` (8) > `purpose` (5) > `stack` (3) >
    `subprojects` (2) > `category` (2) > README/CLAUDE full-text (1). The top hit
    is the candidate; results are labelled `<category>/<name>`. For a precise
    candidate, confirm its path:

    ```bash
    python3 ~/projects/.project_index/router.py cd <project-name>
    # or, when narrowing by category:
    python3 ~/projects/.project_index/router.py category ai
    python3 ~/projects/.project_index/router.py cd ai/msgHub
    ```

3. Decide by confidence:
    - `Single clear top hit` → that is the destination.
    - `Several close scores / ambiguous` → read the candidates' `purpose` from the
      index (or `router.py show <name>`) and pick; if still tied, ask the user.
    - `Category matched, no project` → never file into the category directory
      itself; list its projects (`router.py category <name>`) and pick or ask.
    - `No hit` → do NOT guess. Report "unmatched" and ask, or hold the item.

4. Load the chosen project's context BEFORE acting — always use the `path` from
   the index, which already accounts for the category level:

    ```bash
    P=$(python3 ~/projects/.project_index/router.py cd <name>)
    cat "$P/README.md" "$P/CLAUDE.md" "$P/docs/terminology.md" 2>/dev/null
    ```

    `docs/terminology.md` is the project's term glossary — follow its wording
    when naming files or writing anything into the project.

    Honour that project's CLAUDE.md (naming, folder layout, hard rules). For
    `collections` — a top-level project, not a category, so it is a valid
    destination — that means the `receipts/invoices/agreements` layout and the
    `<source>_<date>_<target>_<amount>.<ext>` naming convention.

5. Propose the destination path and the action. `Confirm before moving,
   renaming, or creating any file` — these are real records. This skill returns
   a routing decision; the caller performs the move only after confirmation.

## Maintaining the index

Rebuild triggers, tag hygiene, and how a category is detected or hand-listed:
[references/maintenance.md](references/maintenance.md). Routing never needs
these commands — only a wrong index does.

## Hard rules

- Resolve to a project FIRST; never file by guesswork.
- Always print the matched project's absolute path clearly to the user once it is found.
- Never route to a `category` directory — it is a container, not a destination.
- Prefer the physical path over a symlinked alias (`ai/cc-plugin` → `cc-plugin`).
- Never move/rename/delete without explicit confirmation.
- The index is the routing signal; per-project README.md + CLAUDE.md +
  docs/terminology.md are the authority on how to file inside the chosen project.

## Related

- `[[project-docs]]` establishes/audits the target project's README.md,
  CLAUDE.md, and docs/terminology.md — run it after routing when the project's
  canonical docs are missing.
