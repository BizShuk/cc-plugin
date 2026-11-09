#!/bin/bash
# setup-links.sh — Converge agent entry files for a project
#
# Usage: ./setup-links.sh [workspace_root]
#   workspace_root: path to the project root (default: current directory)
#
# Target layout:
#   AGENTS.md     real file, the only technical-context entry (Claude reads it directly)
#   CLAUDE.md     must not exist
#   .geminiignore -> .gitignore (symlink)
#
# Legacy layouts are converged in place:
#   AGENTS.md -> CLAUDE.md (symlink)  → AGENTS.md becomes a real file, CLAUDE.md removed
#   CLAUDE.md -> AGENTS.md (symlink)  → CLAUDE.md removed
#   CLAUDE.md real, AGENTS.md absent  → CLAUDE.md renamed to AGENTS.md
#   both real with different content  → WARN and skip (merge by hand)
# Idempotent: re-running on a converged repo changes nothing.

set -euo pipefail

WORKSPACE="${1:-.}"
cd "$WORKSPACE"

# ── helpers ──────────────────────────────────────────────────────────
# Replace a symlink with a real copy of what it points to.
materialize() {
    local file="$1"
    if [ ! -L "$file" ]; then
        return 0
    fi
    if [ ! -e "$file" ]; then
        rm "$file"
        echo "🗑  $file removed (dangling symlink)"
        return 0
    fi
    local tmp
    tmp="$(mktemp "$file.XXXXXX")"
    cat "$file" > "$tmp"
    chmod 644 "$tmp"
    rm "$file"
    mv "$tmp" "$file"
    echo "📄 $file materialized (symlink -> real file)"
}

converge_agents() {
    materialize "AGENTS.md"

    if [ ! -e "AGENTS.md" ] && [ -f "CLAUDE.md" ] && [ ! -L "CLAUDE.md" ]; then
        mv "CLAUDE.md" "AGENTS.md"
        echo "🔁 CLAUDE.md -> AGENTS.md (renamed)"
        return 0
    fi

    if [ -L "CLAUDE.md" ]; then
        rm "CLAUDE.md"
        echo "🗑  CLAUDE.md removed (was a symlink)"
    elif [ -f "CLAUDE.md" ]; then
        if [ -f "AGENTS.md" ] && cmp -s "CLAUDE.md" "AGENTS.md"; then
            rm "CLAUDE.md"
            echo "🗑  CLAUDE.md removed (identical to AGENTS.md)"
        else
            echo "⚠️  WARN: CLAUDE.md and AGENTS.md are both real files with different content, skipping."
        fi
    else
        echo "⏭  CLAUDE.md absent (already converged)"
    fi
}

create_symlink() {
    local target="$1"   # existing file or dir (e.g. .gitignore)
    local link="$2"     # symlink to create (e.g. .geminiignore)

    if [ -L "$link" ]; then
        if [ "$(readlink "$link")" = "$target" ]; then
            echo "⏭  $link -> $target (already a symlink)"
            return 0
        fi
        rm "$link"
    elif [ -e "$link" ]; then
        echo "⚠️  WARN: $link already exists as a regular entry, skipping symlink."
        return 0
    fi

    if [ ! -e "$target" ]; then
        echo "⏭  $link skipped ($target not found)"
        return 0
    fi

    ln -s "$target" "$link"
    echo "✅ $link -> $target (created)"
}

# ── main ─────────────────────────────────────────────────────────────
echo "── setup-links: $(pwd) ──"

converge_agents
create_symlink ".gitignore" ".geminiignore"

echo "── done ──"
