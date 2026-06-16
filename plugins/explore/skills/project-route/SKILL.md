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
version: "1.0.0"
allowed-tools: Bash, Read
metadata:
    type: reference
---

# project-route — route content to a ~/projects/ project

Resolve a piece of content to exactly one project, then load that project's own
docs so downstream work obeys the project's conventions. This skill decides
`where`; it does not move files.

## Config (single source of truth)

- Index: `~/projects/.project_index/projects.json` — one combined JSON.
  - Top-level: `generated_at` (ISO timestamp), `stale_after_days` (default 3),
    `projects_root`, `projects`.
  - Per project: `path`, `tags`, `purpose`, `stack`, `subprojects`,
    `has_readme`, `has_claude`.
  - Per subproject: `name`, `purpose` (extracted from its own `README.md`),
    `has_readme`, `has_claude`.
- Router CLI: `~/projects/.project_index/router.py` — reads the same index.
- `tags` is the routing signal (bilingual, lowercase). `tags` is hand-curated and
  is preserved across `router.py rebuild`; the auto-scan never overwrites it.

## Staleness — auto-rebuild

The router **auto-rebuilds** when the index's `generated_at` is older than
`stale_after_days` (default 3 days). Every `search` / `list` / `show` / `cd` /
`stack` call checks the timestamp first; if stale, it re-scans all of
`~/projects/` (reading each project's `README.md` + `CLAUDE.md` AND every
direct subfolder's `README.md` + `CLAUDE.md` to extract subproject purpose),
then refreshes the JSON. Pass `--no-auto-rebuild` to skip the check (use the
stale index as-is).

## Procedure

1. Extract routing keywords from the content (vendor, currency, domain words,
   filenames, language). Keep both Chinese and English forms — tags carry both.

2. Resolve the project. Prefer the CLI:

   ```bash
   python3 ~/projects/.project_index/router.py search "<keyword>"
   ```

   The CLI scores `name` (10) > `tags` (8) > `purpose` (5) > `stack` (3) >
   `subprojects` (2) > README/CLAUDE full-text (1). The top hit is the candidate.
   For a precise candidate, confirm its path:

   ```bash
   python3 ~/projects/.project_index/router.py cd <project-name>
   ```

3. Decide by confidence:
   - `Single clear top hit` → that is the destination.
   - `Several close scores / ambiguous` → read the candidates' `purpose` from the
     index (or `router.py show <name>`) and pick; if still tied, ask the user.
   - `No hit` → do NOT guess. Report "unmatched" and ask, or hold the item.

4. Load the chosen project's context BEFORE acting:

   ```bash
   cat ~/projects/<name>/README.md ~/projects/<name>/CLAUDE.md 2>/dev/null
   ```

   Honour that project's CLAUDE.md (naming, folder layout, hard rules). For
   `collections`, that means the `receipts/invoices/agreements` layout and the
   `<source>_<date>_<target>_<amount>.<ext>` naming convention.

5. Propose the destination path and the action. `Confirm before moving,
   renaming, or creating any file` — these are real records. This skill returns a
   routing decision; the caller performs the move only after confirmation.

## Maintaining the index

- **You rarely need to rebuild manually** — the router does it for you every
  time the index is >3 days old. Run `python3 ~/projects/.project_index/router.py
  rebuild` only if you want to force a refresh right now (e.g. you just added
  a project, edited tags, or want to bypass the 3-day window).
- A project with no tags is invisible to tag-based routing — add tags so it
  can receive content.

## Hard rules

- Resolve to a project FIRST; never file by guesswork.
- Never move/rename/delete without explicit confirmation.
- The index is the routing signal; per-project README.md + CLAUDE.md are the
  authority on how to file inside the chosen project.
