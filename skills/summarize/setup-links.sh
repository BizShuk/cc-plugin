#!/bin/bash
# setup-links.sh — Create symbolic links for multi-agent project setup
#
# Usage: ./setup-links.sh [workspace_root]
#   workspace_root: path to the project root (default: current directory)
#
# Creates:
#   GEMINI.md     -> CLAUDE.md
#   AGENTS.md     -> CLAUDE.md
#   .geminiignore -> .gitignore

set -euo pipefail

WORKSPACE="${1:-.}"
cd "$WORKSPACE"

# ── helpers ──────────────────────────────────────────────────────────
create_symlink() {
    local target="$1"   # existing file (e.g. CLAUDE.md)
    local link="$2"     # symlink to create (e.g. GEMINI.md)

    if [ -L "$link" ]; then
        echo "⏭  $link -> $(readlink "$link") (already a symlink)"
        return 0
    fi

    if [ -e "$link" ]; then
        echo "⚠️  WARN: $link already exists as a regular file, skipping symlink."
        return 0
    fi

    if [ ! -f "$target" ]; then
        echo "⏭  $link skipped ($target not found)"
        return 0
    fi

    ln -s "$target" "$link"
    echo "✅ $link -> $target (created)"
}

# ── main ─────────────────────────────────────────────────────────────
echo "── setup-links: $(pwd) ──"

create_symlink "CLAUDE.md"  "GEMINI.md"
create_symlink "CLAUDE.md"  "AGENTS.md"
create_symlink ".gitignore" ".geminiignore"

echo "── done ──"
