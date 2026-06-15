---
name: folder-structure
description: >
    Review a project's directory layout for clarity and convention — misplaced
    files, mixed concerns, inconsistent depth, orphaned or duplicated folders,
    and drift from the language's idiomatic structure. Use when asked to review
    or reorganize folders. Triggers on: "review folder structure", "is this
    layout right", "reorganize directories", "資料夾結構", "目錄結構審查".
version: "1.0.0"
allowed-tools: Read, Bash, Glob, Grep
user-invocable: true
disable-model-invocation: false
effort: medium
context: fork
metadata:
    type: review
---

# Folder Structure Review

Judge the layout against the project's own declared structure first, then
against the language's idiom — not against a generic ideal. Report misplacements
with a concrete target path, not vague preferences.

## Procedure

1. Read the structure section of `CLAUDE.md` (or `README.md`). That declared
   tree is the contract; deviations from it are the primary findings.
2. Generate the actual tree (`find . -type d` minus ignored paths) and diff it
   against the declared one.
3. For each directory, check it against the heuristics below.
4. Output findings as `current path → suggested path` with a one-line reason.

## Heuristics

| Smell                | Rule it violates                                         |
| -------------------- | -------------------------------------------------------- |
| Mixed concerns       | One folder holds unrelated responsibilities              |
| Wrong layer          | A file sits in a layer that should not own it             |
| Inconsistent depth   | Sibling features nested at different depths               |
| Orphan               | A folder nothing references and no doc explains           |
| Duplication          | Two folders doing the same job under different names      |
| Loose file           | A file at a level the convention reserves for subfolders  |
| Undeclared           | A real folder the structure doc never mentions            |

## Convention first

If the project (or its language) already has a standard place for something,
the file belongs there — do not invent a parallel location. Flag any custom
path that competes with an established convention.

## Output

```
Folder structure review (declared tree = source of truth)
- cmd/helpers.go        → cmd/util/ (loose file; siblings are subpackages)
- pkg/old-litellm/      → remove (orphan; nothing imports it)
- DRIFT: plugins/god/ exists but CLAUDE.md tree omits it
```

When the layout is sound, say so and stop. Do not propose churn for taste.

## Related

- `[[naming-convention]]` for the names of the folders themselves
- `[[doc-sync]]` to fix the declared tree once the real one changes
