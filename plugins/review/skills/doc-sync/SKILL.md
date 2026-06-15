---
name: doc-sync
description: >
    Review whether docs match the code they describe — README.md, CLAUDE.md,
    structure trees, module maps, build/run commands, config paths, and code
    comments. Reports drift as doc-claim vs actual, without rewriting unasked.
    Use when verifying docs are current. Triggers on: "docs out of date",
    "doc sync", "does the README match", "文件同步", "更新文件了嗎",
    "check the structure tree".
version: "1.0.0"
allowed-tools: Read, Bash, Glob, Grep
user-invocable: false
disable-model-invocation: false
effort: medium
context: fork
metadata:
    type: review
---

# Doc Sync Review

Documentation drifts the moment code changes without it. This skill finds the
drift; it does not rewrite docs unless asked. Every finding pairs a doc claim
with the contradicting reality.

## Procedure

1. Read `README.md` (business scope) and `CLAUDE.md` (structure, tech, build).
2. Extract every verifiable claim — file/dir paths, commands, module mappings,
   entry points, config paths, version pins, counts.
3. Verify each against the repo: does the path exist, does the command run, does
   the named function exist, does the tree match `find`?
4. Report drift as `doc says X → actually Y`, grouped by file.

## What drifts most

| Claim type     | How to verify                                             |
| -------------- | --------------------------------------------------------- |
| Structure tree | Diff the documented tree against the real directory tree  |
| Module map     | Each named entry point still exists at the named path     |
| Build / run    | The documented command exists and is still the real one   |
| Config paths   | The `~/...` paths the doc cites match what the code reads |
| Counts / lists | "skills: [...]" arrays match the folders on disk          |
| Comments       | A comment describing behavior the code no longer has      |

## Output

```
Doc sync review
- CLAUDE.md tree: omits plugins/god/ and plugins/superpower/ (both exist)
- CLAUDE.md module map: DistillCmd() — confirmed at cmd/distill.go
- README.md: "uses MySQL" → code uses SQLite (model/store.go)
```

If asked to fix, update the doc to match the code (code is the source of truth),
then re-verify. Do not silently change scope while syncing.

## Related

- `[[folder-structure]]` when the real tree, not the doc, is what is wrong
- `[[consistency]]` for code-to-code contradictions rather than doc-to-code
