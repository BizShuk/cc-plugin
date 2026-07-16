---
name: changelog
description: >
    Use when generating or maintaining CHANGELOG.md with weekly LLM-narrated
    sections for git repositories under a root directory. Runs a deterministic
    stats pipeline first, then spawns parallel agents per repository. Triggers on:
    "generate changelog", "update CHANGELOG", "summarize git history".
---

# /changelog

ultracode: Generate per-repo `CHANGELOG.md` with weekly LLM-narrated sections.

## Usage

```
/changelog <root-dir>
```

Defaults to current working directory if no root given.

## Workflow

### Phase 1: Deterministic pipeline

First install the package, then run the pipeline for all repos:

```bash
pip install -e .claude/skills/changelog/scripts
repo-changelog run-all <root>
```

This produces per repo: `stats.json`, `_diffs/*.diff`, and a `CHANGELOG.md` skeleton with `<!-- LLM: ... -->` placeholders.

### Phase 2: LLM narrative (parallel subagents)

For each repo with `stats.json`:

1. Read `stats.json` to get the weekly bucket list and top-level stats
2. For each week, read `_diffs/<week>.diff` and the existing commit list in the skeleton
3. **Replace the `<!-- LLM: ... -->` placeholder** with a 3-5 sentence narrative summarizing the feature changes visible in the diff. Write in past tense. Focus on what changed from a user/API perspective, not implementation details. Mention key files/modules only when they clarify the change.
4. If a week's diff is empty (no kept-file changes), write: `_No business-logic changes this week._`
5. If a week's diff exceeds ~100K tokens, add `<!-- warning: large diff (~N tokens) -->` before the narrative.

**Concurrency:** Spawn subagents with `run_in_background: true`, bounded to ~8 concurrent. Each subagent handles one repo (all its weeks).

**Subagent prompt template:**

```
You are writing the weekly narrative sections for a CHANGELOG.md.

Repo: <repo_name>
Output path: <output_path>/CHANGELOG.md

The CHANGELOG.md skeleton already exists at the output path. It contains:
- Header with repo stats
- Committer activity table
- Top 10 commits table
- Weekly sections with commit lists and `<!-- LLM: ... -->` placeholders

For each week, read the corresponding diff file at `<output_path>/_diffs/<week>.diff`.
Replace the `<!-- LLM: ... -->` comment with a 3-5 sentence narrative summarizing the
feature changes in the diff. Focus on what changed from a user/API perspective.
Write in past tense. Keep it concise.

Do NOT change the commit lists, stat tables, or any other part of the file.
Only replace the `<!-- LLM: ... -->` placeholders.
```

### Phase 3: Summary

After all subagents complete, print a summary:

- Total repos processed
- Total weeks with narratives written
- Any repos/week that failed (with `<!-- LLM: pending -->` placeholders)
