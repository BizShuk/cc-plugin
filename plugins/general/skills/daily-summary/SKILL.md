---
name: daily-summary
description: >
  Use when the user asks to summarize the day's work, generate a daily summary
  / 工作日報, recap what was done across AI sessions and repos, or extract today's
  TODOs into Apple Notes. Pulls from claude-mem, Claude sessions, Antigravity
  sessions, Hermes sessions, and git log across ~/projects, then writes a note to
  the iCloud "Daily" folder. Triggers on: "daily summary", "summarize my work
  today", "工作日報", "what did I do today", "今天做了什麼".
version: "1.0.0"
allowed-tools: Read, Bash, Write
metadata:
  type: technique
  platforms: [macos]
---

# Daily Work Summary

## Overview

Aggregate the last 24h of work from five local sources, synthesize a summary +
TODO lists, and write it to Apple Notes (`iCloud` › `Daily`) titled
`Daily Summary - {YYYY-MM-DD}`. If today's note already exists, append to it.

Core rule: always run `scan.sh` first. It tells you which sources actually have
data in the window so you don't waste calls on empty ones.

## When to Use

- "summarize my work today" / "工作日報" / "what did I do today"
- Building a daily recap of AI sessions + code changes
- Extracting today's open TODOs into one place

## Step 0 — Preflight scan (ALWAYS FIRST)

Run the `scan.sh` that sits beside this SKILL.md (arg = window hours):

```bash
bash "${CLAUDE_PLUGIN_ROOT:-$HOME/projects/cc-plugin/plugins/general}/skills/daily-summary/scan.sh" 24
```

Read its output. Process only sources marked `[OK]` with a non-zero count.
Skip `[MISS]`/empty sources — note them in the summary as "no data".

## Sources (explicit paths)

| # | Source | Path | How to pull (window = last 24h) |
|---|--------|------|----------------------------------|
| 1 | claude-mem | `~/.claude-mem/claude-mem.db` | `sqlite3`: `session_summaries` (request/completed/learned/next_steps/files_edited) + `observations`, `WHERE created_at_epoch > <cutoff>`. PRIMARY source. |
| 2 | Claude sessions | `~/.claude/projects/<encoded>/*.jsonl` | Raw transcripts. Only if claude-mem is empty — find `*.jsonl` newer than cutoff. claude-mem already distills these, so usually skip. |
| 3 | Antigravity sessions | `~/.gemini/antigravity/brain/<uuid>/task.md` | Plain-markdown task plans. `[x]`=done, `[ ]`=open TODO. Read files newer than cutoff. (NOT the empty `~/Library/Application Support/Antigravity`.) |
| 3b| Antigravity convos | `~/.gemini/antigravity/conversations/*.db` | SQLite, but `steps` are binary BLOBs — only report which dbs were touched (activity signal), don't try to parse bodies. |
| 4 | Hermes sessions (channel-aware) | `~/.hermes/state.db` | SQLite. `sessions` table: `source` = channel (slack/telegram/whatsapp/discord/signal/email/sms/cron/cli/tui/webui/subagent…), plus `title`, `started_at`/`ended_at` (REAL epoch), `message_count`, `user_id`, `cwd`. `messages` table: `session_id`, `role`, `content`, `timestamp` (REAL epoch). Join on `sessions.id = messages.session_id` to attribute messages to a channel. PRIMARY Hermes source. (Legacy `~/.hermes/sessions/*.jsonl` is dead since ~2026-05 — ignore.) |
| 5 | git log | `~/projects/**/.git` | Per repo: `git -C <repo> log --since="24 hours ago" --pretty="%s"`. Concrete shipped changes (also covers Antigravity-edited files). |

`cutoff` epoch = `now - hours*3600`. The scan script prints it.

### Pull commands

```bash
CUT=$(( $(date +%s) - 24*3600 ))
# 1. claude-mem — per-project completed work (dedupe noisy repeats yourself)
sqlite3 -json ~/.claude-mem/claude-mem.db \
  "SELECT project, request, completed, next_steps FROM session_summaries \
   WHERE created_at_epoch > $CUT AND length(completed)>40 \
   ORDER BY project, created_at_epoch;"
# 3. antigravity open TODOs
find ~/.gemini/antigravity/brain -name task.md -newermt "@$CUT" \
  -exec sh -c 'echo "## $1"; grep -nE "^\s*-\s*\[ \]" "$1"' _ {} \;
# 4. hermes sessions — active sessions per channel + their messages in window
sqlite3 -json ~/.hermes/state.db \
  "SELECT s.source AS channel, s.title, s.message_count, \
          datetime(s.started_at,'unixepoch','localtime') AS started \
   FROM sessions s \
   WHERE s.id IN (SELECT DISTINCT session_id FROM messages WHERE timestamp > $CUT) \
   ORDER BY s.started_at DESC;"
# message bodies for those sessions (user asks + assistant outcomes), per channel
sqlite3 -json ~/.hermes/state.db \
  "SELECT s.source AS channel, m.role, substr(m.content,1,200) AS content \
   FROM messages m JOIN sessions s ON s.id=m.session_id \
   WHERE m.timestamp > $CUT AND m.role IN ('user','assistant') AND length(m.content)>0 \
   ORDER BY m.timestamp;"
# 5. git per repo
find ~/projects -maxdepth 3 -name .git -type d | while read -r g; do r=$(dirname "$g");
  c=$(git -C "$r" log --since="24 hours ago" --pretty="  - %s" 2>/dev/null);
  [ -n "$c" ] && printf '### %s\n%s\n' "$(basename "$r")" "$c"; done
```

## Synthesize

Produce three sections (Traditional Chinese, per user's global style):

1. `## 工作摘要` — group by project. Use git commits as the "shipped/confirmed"
   spine; layer claude-mem `completed` + Hermes session outcomes on top for intent
   + uncommitted work. Dedupe claude-mem's repetitive per-prompt rows into one
   bullet per theme.
2. `## TODO（明確）` — open `[ ]` items from Antigravity `task.md` +
   explicit `next_steps` from claude-mem + unfinished asks from Hermes sessions.
3. `## TODO（潛在）` — inferred follow-ups: things left mid-flight, "next would
   be…", failing checks, deferred decisions. Mark clearly as inferred.

End with a one-line source footnote (which sources had data, which were empty).

## Write to Apple Notes

Target: account `iCloud`, folder `Daily`, title `Daily Summary - {YYYY-MM-DD}`
(use local date). Body is Markdown.

```bash
DATE=$(date +%F); TITLE="Daily Summary - $DATE"; ACC=iCloud; FOLDER=Daily
# 1. ensure folder exists (idempotent — ignore error if already there)
notes list -a "$ACC" -f "$FOLDER" -i 2>/dev/null | grep -q . || notes mkdir -a "$ACC" "$FOLDER"
# 2. append if today's note exists, else create
EXIST=$(notes list -a "$ACC" -f "$FOLDER" -n "$TITLE" -i 2>/dev/null | head -1)
BODY=$(cat /tmp/daily-summary.md)   # write your synthesized markdown here first
if [ -n "$EXIST" ]; then
  OLD=$(notes get "$EXIST" --body-only 2>/dev/null)
  printf '%s\n\n---\n%s' "$OLD" "$BODY" | notes edit "$EXIST" -m
else
  printf '%s\n\n%s' "$TITLE" "$BODY" | notes add -a "$ACC" -f "$FOLDER" -m
fi
```

`notes add` treats the first line as the title, rest as body. `-m` = Markdown.

## Common Mistakes

| Mistake | Fix |
|---------|-----|
| Dumping all 900+ claude-mem rows | They are per-prompt and repetitive. Aggregate to one bullet per project theme. |
| Using `~/Library/.../Antigravity` | That store is empty/cloud-bound. Real data is `~/.gemini/antigravity/`. |
| Parsing antigravity conversation `.db` bodies | `steps` are binary BLOBs. Use `brain/*/task.md` for content; dbs only as an activity signal. |
| Reading Hermes from `sessions/*.jsonl` | That dir is legacy (dead since ~2026-05). Use `~/.hermes/state.db` (`sessions`+`messages`). |
| Ignoring Hermes channel attribution | Group Hermes work by `sessions.source` (slack/whatsapp/cron/cli…) so message-channel work is visible, not lumped together. |
| Overwriting an existing daily note | Check for the title first; append with a `---` divider. |
| Skipping the scan | Run `scan.sh` first — it prevents querying empty sources and prints the cutoff. |
