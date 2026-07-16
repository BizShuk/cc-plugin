---
name: sort-todo
description: >
    Use when asked to sort, prioritize, re-prioritize, triage, or re-organize a
    `.todo` file such as README.todo. Assigns each item a P0/P1/P2 from business,
    UX, and system lenses, then regroups items under feature-domain sections.
    Triggers only for files with the `.todo` extension.
disable-model-invocation: true
---

# Sort TODO

## Overview

Re-prioritize and re-organize a `.todo` file (the project's is `README.todo`) so every actionable item carries a justified `P0/P1/P2` and lives under a feature-domain section (not a date-stamped "Iteration N sweep"). Priority is a judgment call across three lenses — business, UX, system — never a guess.

Core principle: **the file is parsed live by the TODO panel.** Sorting that breaks the line format silently breaks the panel. Format-correctness is non-negotiable; see [Hard format rules](#hard-format-rules).

**Scope guard — `.todo` only.** This skill operates exclusively on files whose name ends in `.todo`. If the target is a `.md`, `.markdown`, or any other extension, STOP and tell the user this skill only sorts `.todo` files — do not touch it, even if it contains a checkbox/TODO list. A `README.md` is never in scope; `README.todo` is.

## When to Use

- "sort / triage / re-prioritize the todo list", "organize README.todo", "what should I work on next"
- After a planning sweep added many untagged items
- When sections are date-stamped sweeps (`### Iteration 4 …`) instead of feature names

Do NOT use for: any file that is not `.todo` (e.g. `README.md`, `*.markdown` — out of scope, leave untouched), toggling one item done (panel does that), adding a single TODO (`superset.todoNew`), editing `plans/` specs.

## Workflow

0. **Check the extension first.** Confirm the target file ends in `.todo`. If not, stop — this skill does not apply (see [Scope guard](#overview)).
1. **Read the whole file** (the `.todo` file, e.g. `README.todo`) — every item, in every section, including `## Archive` and `## Plans`. Sorting a subset produces wrong priorities (you can't rank what you didn't read).
2. **Open referenced plans when ambiguous.** Items link `— see [plans/…md]`. If a one-line item is too thin to score, read its plan before assigning priority. Don't invent value the plan contradicts.
3. **Score each unchecked item** across the three lenses → derive `P0/P1/P2` (see [Priority rubric](#priority-rubric)).
4. **Re-home into feature-domain sections** (see [Section taxonomy](#section-taxonomy)). Within each section order `P0 → P1 → P2`, untagged last, completed `[x]` items last of all (kept in place; the panel hides them by default).
5. **Leave `## Archive` untouched.** Do NOT reorder, add to, remove from, or re-prioritize anything inside `## Archive`. Read it (step 1) but treat it as frozen. Completed items elsewhere are NOT moved into it.
6. **Apply to `README.todo` only — do not commit.** Write the sorted file directly with the edit tools. Never run `git add` / `git commit` / `git push`; leave the change in the working tree for the user to review.
7. **Report the diff in prose**: what moved, what was re-prioritized and the lens that drove it, what you couldn't decide. Then, as a final list, surface any **valuable item currently sitting in `## Archive`** (see [Archive surfacing](#archive-surfacing)) — list it, do not move it.

## Priority rubric

Rate each unchecked item on three lenses, then combine. Score the *item as written*, not the feature you wish it were.

| Lens | High signal | Low signal |
| --- | --- | --- |
| Business (B) | Gates a release; blocks adoption; correctness bug users hit | Cosmetic; internal-only nicety |
| UX (U) | Removes daily friction; improves discoverability of an existing surface | Edge-case convenience; power-user-only |
| System (S) | Foundation many other items depend on; stability/maintainability debt | Isolated; touches one file; no dependents |

Combine:

- `P0` — any lens **High** AND (blocks shipping **or** is a correctness/lifecycle bug **or** is a foundation ≥2 other items depend on). Unchecked bug-style items with no tag (e.g. a parsing error) default to **at least P1, usually P0**.
- `P1` — clear value, ≥1 lens High/Med, but nothing is blocked waiting on it.
- `P2` — nice-to-have, isolated, low reach. Most `[chore]` / `[refactor]` polish and speculative features land here (unless they unblock a release — see the dependency override).

Dependency override: if item A is a prerequisite for B/C/D, A inherits the **max** priority of its dependents. A `chore` that unblocks a release (CI, VSIX packaging) is business-High even though it's "just a chore".

When two items tie, the one with more dependents or a written plan ranks first.

## Section taxonomy

Replace date-sweep headings with feature domains drawn from the codebase (`src/<feature>/`). Recommended set for this project — create a section only when it has items:

```tree
README.todo
├── (top, no heading)   # cross-cutting P0 bugs not tied to one feature
├── ## Terminals        # [feature]/bug — panel, highlight, PTY, jump-to, fuzzy
├── ## mDNS             # [feature]/bug — discovery, dedup, detail cache, one-click connect
├── ## Topology         # [feature]/bug — scan, routing parse
├── ## TODO Panel       # [feature]/bug — the todo view itself
├── ## Tree Preview     # [feature]/bug — markdown tree block
├── ## Platform & UX    # [feature]/bug — settings webview, caches reset, layout persistence, reveal-in-tree, diagnostics
├── ## Architecture     # ALL non-functional work — chore, refactor, build/packaging, CI, deps, docs, tooling, tests, perf infra, baseline
├── ## Plans            # links to plans/ implementation docs (leave as-is unless re-themed)
└── ## Archive          # FROZEN — never reorder/add/remove; only read it
```

Rule: a new section name must name a **feature or concern**, never a date or sweep number.

**Functional vs non-functional routing.** Decide first whether an item is functional (changes user-facing behavior) or not:

- **Non-functional → `## Architecture`.** Everything that is not a user-facing feature: `[chore]`, `[refactor]`, build/packaging, CI/CD, dependencies, version baselines, docs, tooling, test infrastructure, performance/structural work. These don't add or change a runtime feature, so they live together regardless of which feature's code they touch — not scattered across feature sections.
- **Functional → feature-domain section.** `[feature]` items and bugs (untagged actionable defects in runtime behavior) → the matching feature section. If a feature item fits two domains, file it under the one whose code it changes most.

Test: "Does completing this change what the user can do or observe in the extension?" No → `## Architecture`.

## Archive surfacing

`## Archive` is read-only during a sort (see Workflow step 5), but it may hide items worth reviving. After rewriting the file, scan Archive once and surface — in the final response only, not in the file — any item that scores `P0`/`P1` under the [Priority rubric](#priority-rubric) (e.g. a `[feature]` with a written plan, or a still-unchecked `[ ]` parked there). For each, give: the item text, the lens that makes it valuable, and the priority it would earn if revived. Do not move or edit it; the user decides whether to pull it out.

## Hard format rules

Breaking any of these breaks the live panel parser (`src/todo/todoStore.ts`, `todoTreeProvider.ts`):

- Priority tag is **first** in the item text, right after the checkbox: `- [ ] [P0] [feature] …`. The panel's priority regex only matches a *leading* `[Px]` / `(Px)`. `- [ ] [feature] [P0]` will NOT be recognized.
- Canonical line: `` `- [ ] [Px] [type] description — see [link](path)` `` (priority → type tag → text → reference).
- Keep checkbox state (`[ ]` vs `[x]`) exactly; sorting never silently completes or un-completes an item.
- Sections are `##` / `###` headings only. One blank line between a heading and its items.
- Preserve `[feature]` / `[chore]` / `[refactor]` tags and every markdown link target.
- Only `P0`, `P1`, `P2` are valid (or none). No `P3`, no bare numbers.

## Common Mistakes

| Mistake | Fix |
| --- | --- |
| Priority tag placed after `[feature]` | Tag goes first: `[P0] [feature]`, panel won't read it otherwise |
| Sorting only the visible/top items | Read the entire file first; priorities are relative across all items |
| Keeping `### Iteration N (date sweep)` headings | Re-theme to feature domains |
| Filing a chore/CI/docs/refactor item under a feature section | All non-functional work → `## Architecture`; ask "does it change what the user can do?" first |
| Guessing priority from the title | Open the linked plan when the item is thin |
| Assigning P0 to everything urgent-sounding | P0 is reserved: High lens + blocking/bug/foundation |
| Touching `## Archive` (sorting/adding/removing) | Archive is frozen; only read it, then surface valuable items in the response |
| Committing the sorted file | Apply to the working tree only; never `git commit`/`push` |
| Dropping `— see [link]` references when moving | Move the whole line verbatim, links intact |

## Output language

Per project convention, write any new section names and the change summary in Traditional Chinese with the English term in round brackets where natural (e.g. `## 終端機 (Terminals)`), and use `backtick` rather than bold for highlight. Item text itself is preserved as written.
