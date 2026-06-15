---
name: naming-convention
description: >
    Review identifiers — files, packages, types, functions, variables, configs,
    endpoints — for clarity and a single consistent convention. Flags casing
    drift, synonym sprawl, vague or misleading names, and abbreviations the
    codebase does not otherwise use. Use when asked to review naming. Triggers
    on: "review naming", "naming convention", "rename suggestions", "命名一致",
    "命名規範", "is this name clear".
version: "1.0.0"
allowed-tools: Read, Bash, Glob, Grep
user-invocable: true
disable-model-invocation: false
effort: medium
context: fork
metadata:
    type: review
---

# Naming Convention Review

A name is good when a reader predicts behavior from it and the same idea always
wears the same word. Infer the codebase's existing convention before judging;
the established pattern wins over any external style guide.

## Procedure

1. Sample existing names per category (files, types, funcs, vars, consts) and
   infer the de facto convention — casing, word order, prefixes.
2. Grep for the change's new names and for synonyms of each concept it touches.
3. Score each name against the dimensions below; the convention you inferred in
   step 1 is the baseline, not a generic guide.
4. Output `current → suggested` with a one-line reason. Suggest only net wins.

## Dimensions

| Dimension     | Failure                                                     |
| ------------- | ---------------------------------------------------------- |
| Casing        | `userId` next to `user_id`; convention applied unevenly    |
| Synonym drift | `fetch` / `get` / `load` for the same operation            |
| Vagueness     | `data`, `handle`, `manager`, `util` that say nothing       |
| Misleading    | Name implies behavior the code does not have               |
| Abbreviation  | A shortened form the rest of the codebase spells out        |
| Symmetry      | `open` without `close`; `start` without `stop`              |
| Scope fit     | A long name for a tiny scope, or a terse one for a wide API |

## Output

```
Naming review (inferred convention: Go exported = PascalCase, files = snake)
- cmd/Read_Logic.go     → read_logic.go (file casing drift)
- func GetUser/FetchUser → pick one verb for the same op (synonym drift)
- var d                  → decoded (vague in a 40-line scope)
```

Do not rename for taste. Every suggestion must reduce ambiguity or restore the
codebase's own consistency.

## Related

- `[[consistency]]` for terminology drift beyond identifiers
- `[[folder-structure]]` for the names of directories
