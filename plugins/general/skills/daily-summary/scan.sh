#!/usr/bin/env bash
# scan.sh — preflight source scan for the daily-summary skill.
# Checks every source listed below, reports availability + how many items fall
# inside the time window. Run this FIRST; the skill reads its output to decide
# which sources to actually pull. Read-only: never writes anything.
#
# Usage:  ./scan.sh [WINDOW_HOURS]      (default 24)

set -uo pipefail
WINDOW_HOURS="${1:-24}"
NOW="$(date +%s)"
CUTOFF=$(( NOW - WINDOW_HOURS * 3600 ))
SINCE_HUMAN="$(date -r "$CUTOFF" '+%Y-%m-%d %H:%M' 2>/dev/null || date '+%Y-%m-%d %H:%M')"

# ---- explicit source paths (edit here if a path moves) ----------------------
CLAUDE_MEM_DB="$HOME/.claude-mem/claude-mem.db"
CLAUDE_SESSIONS_DIR="$HOME/.claude/projects"
AG_CONVO_DIR="$HOME/.gemini/antigravity/conversations"
AG_BRAIN_DIR="$HOME/.gemini/antigravity/brain"
HERMES_STATE_DB="$HOME/.hermes/state.db"          # live channel-aware store (sessions+messages)
HERMES_SESSIONS_DIR="$HOME/.hermes/sessions"      # legacy jsonl (dead since ~2026-05)
PROJECTS_DIR="$HOME/projects"
# Apple Notes target (output sink — verified, not scanned for content)
NOTES_ACCOUNT="iCloud"
NOTES_FOLDER="Daily"
# ----------------------------------------------------------------------------

ok()   { printf '  [OK]   %s\n' "$*"; }
miss() { printf '  [MISS] %s\n' "$*"; }
warn() { printf '  [WARN] %s\n' "$*"; }

echo "================ daily-summary preflight scan ================"
echo "window      : last ${WINDOW_HOURS}h  (since ${SINCE_HUMAN})"
echo "cutoff epoch: ${CUTOFF}"
echo

# 1. claude-mem -------------------------------------------------------------
echo "1) claude-mem            $CLAUDE_MEM_DB"
if [ -f "$CLAUDE_MEM_DB" ] && command -v sqlite3 >/dev/null 2>&1; then
  ss=$(sqlite3 "$CLAUDE_MEM_DB" "SELECT count(*) FROM session_summaries WHERE created_at_epoch > $CUTOFF;" 2>/dev/null)
  ob=$(sqlite3 "$CLAUDE_MEM_DB" "SELECT count(*) FROM observations      WHERE created_at_epoch > $CUTOFF;" 2>/dev/null)
  pj=$(sqlite3 "$CLAUDE_MEM_DB" "SELECT count(DISTINCT project) FROM session_summaries WHERE created_at_epoch > $CUTOFF;" 2>/dev/null)
  ok "session_summaries=${ss:-0}  observations=${ob:-0}  projects=${pj:-0}"
else
  miss "db not found or sqlite3 missing"
fi
echo

# 2. claude sessions (raw jsonl transcripts) --------------------------------
echo "2) claude sessions       $CLAUDE_SESSIONS_DIR/<encoded>/*.jsonl"
if [ -d "$CLAUDE_SESSIONS_DIR" ]; then
  n=$(find "$CLAUDE_SESSIONS_DIR" -name '*.jsonl' -newermt "@$CUTOFF" 2>/dev/null | wc -l | tr -d ' ')
  ok "jsonl files touched in window=${n}"
else
  miss "dir not found"
fi
echo

# 3. antigravity sessions ---------------------------------------------------
echo "3) antigravity tasks     $AG_BRAIN_DIR/<uuid>/task.md"
if [ -d "$AG_BRAIN_DIR" ]; then
  tasks=$(find "$AG_BRAIN_DIR" -name 'task.md' -newermt "@$CUTOFF" 2>/dev/null)
  tn=$(printf '%s\n' "$tasks" | grep -c . )
  open_items=0
  while IFS= read -r f; do
    [ -n "$f" ] || continue
    c=$(grep -cE '^\s*-\s*\[ \]' "$f" 2>/dev/null)
    open_items=$(( open_items + c ))
  done <<< "$tasks"
  ok "task.md updated in window=${tn}  open TODO items=${open_items}"
else
  miss "brain dir not found"
fi
echo "   antigravity convo db  $AG_CONVO_DIR/*.db  (binary trajectories)"
if [ -d "$AG_CONVO_DIR" ]; then
  cdb=$(find "$AG_CONVO_DIR" -name '*.db' -newermt "@$CUTOFF" 2>/dev/null | wc -l | tr -d ' ')
  if [ "${cdb:-0}" -gt 0 ]; then ok "conversation dbs touched=${cdb}"; else warn "no conversation db touched in window (activity may be older)"; fi
else
  miss "conversations dir not found"
fi
echo

# 4. hermes sessions (channel-aware: slack/telegram/whatsapp/cron/cli/...) --
echo "4) hermes sessions       $HERMES_STATE_DB  (sessions+messages by channel)"
if [ -f "$HERMES_STATE_DB" ] && command -v sqlite3 >/dev/null 2>&1; then
  hm=$(sqlite3 "$HERMES_STATE_DB" "SELECT count(*) FROM messages WHERE timestamp > $CUTOFF;" 2>/dev/null)
  hs=$(sqlite3 "$HERMES_STATE_DB" "SELECT count(DISTINCT m.session_id) FROM messages m WHERE m.timestamp > $CUTOFF;" 2>/dev/null)
  chans=$(sqlite3 "$HERMES_STATE_DB" "SELECT group_concat(source||':'||n,' ') FROM (SELECT s.source, count(DISTINCT m.session_id) n FROM messages m JOIN sessions s ON s.id=m.session_id WHERE m.timestamp > $CUTOFF GROUP BY s.source ORDER BY n DESC);" 2>/dev/null)
  if [ "${hm:-0}" -gt 0 ]; then ok "active sessions=${hs:-0}  messages=${hm}  channels=[${chans}]"; else warn "no hermes messages in window (activity may be older)"; fi
  # legacy jsonl fallback (informational only)
  if [ -d "$HERMES_SESSIONS_DIR" ]; then
    lj=$(find "$HERMES_SESSIONS_DIR" -name '*.jsonl' -newermt "@$CUTOFF" 2>/dev/null | wc -l | tr -d ' ')
    [ "${lj:-0}" -gt 0 ] && warn "legacy jsonl also touched=${lj} (usually empty — prefer state.db)"
  fi
else
  miss "state.db not found or sqlite3 missing"
fi
echo

# 5. git log across ~/projects ---------------------------------------------
echo "5) git log               $PROJECTS_DIR/**/.git"
if [ -d "$PROJECTS_DIR" ] && command -v git >/dev/null 2>&1; then
  hits=0
  while IFS= read -r g; do
    repo=$(dirname "$g")
    n=$(git -C "$repo" log --since="${WINDOW_HOURS} hours ago" --oneline 2>/dev/null | wc -l | tr -d ' ')
    if [ "${n:-0}" -gt 0 ]; then
      printf '         %-28s %s commits\n' "$(basename "$repo")" "$n"
      hits=$(( hits + 1 ))
    fi
  done < <(find "$PROJECTS_DIR" -maxdepth 3 -name .git -type d 2>/dev/null)
  if [ "$hits" -gt 0 ]; then ok "${hits} repo(s) with commits in window"; else warn "no repo with commits in window"; fi
else
  miss "projects dir or git missing"
fi
echo

# 6. Apple Notes output sink ------------------------------------------------
echo "6) Apple Notes sink      account=${NOTES_ACCOUNT} folder=${NOTES_FOLDER}  (via 'notes' CLI)"
if command -v notes >/dev/null 2>&1; then
  if notes accounts 2>/dev/null | grep -qi "name: ${NOTES_ACCOUNT}"; then
    ok "account '${NOTES_ACCOUNT}' present"
  else
    warn "account '${NOTES_ACCOUNT}' not found in 'notes accounts'"
  fi
  if notes list -a "$NOTES_ACCOUNT" -f "$NOTES_FOLDER" -i 2>/dev/null | grep -q .; then
    ok "folder '${NOTES_FOLDER}' exists"
  else
    warn "folder '${NOTES_FOLDER}' not found — skill will 'notes mkdir' it"
  fi
else
  miss "'notes' CLI not installed"
fi
echo "=================================================================="
