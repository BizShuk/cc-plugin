#!/bin/bash
# setup-links.sh — Create symbolic links for multi-agent project setup
#
# Usage: ./setup-links.sh [workspace_root]
#   workspace_root: path to the project root (default: current directory)
#
# Real sources are agent-agnostic; Claude-specific names are symlinks:
#   CLAUDE.md     -> AGENTS.md
#   .claude       -> .agents
#   .geminiignore -> .gitignore
#
# Legacy layouts (CLAUDE.md / .claude as the real entry) are inverted in place.
# Idempotent: re-running on a converged repo changes nothing.

set -euo pipefail

WORKSPACE="${1:-.}"
cd "$WORKSPACE"

# ── helpers ──────────────────────────────────────────────────────────
is_empty_dir() {
    [ -d "$1" ] && [ -z "$(ls -A "$1")" ]
}

# Make $real the real file and $link a symlink to it.
converge_file() {
    local real="$1"   # e.g. AGENTS.md
    local link="$2"   # e.g. CLAUDE.md

    # Legacy: link is the real file, real is a symlink back to it (or absent).
    if [ -f "$link" ] && [ ! -L "$link" ]; then
        if [ -L "$real" ] || [ ! -e "$real" ]; then
            rm -f "$real"
            mv "$link" "$real"
            echo "🔁 $link -> $real (moved content into $real)"
        else
            echo "⚠️  WARN: both $real and $link are regular files, skipping."
            return 0
        fi
    fi

    create_symlink "$real" "$link"
}

# Make $real the real directory and $link a symlink to it.
converge_dir() {
    local real="$1"   # e.g. .agents
    local link="$2"   # e.g. .claude

    if [ -d "$link" ] && [ ! -L "$link" ]; then
        if [ -L "$real" ]; then
            rm "$real"
        fi
        if [ ! -e "$real" ] || is_empty_dir "$real"; then
            [ -d "$real" ] && rmdir "$real"
            mv "$link" "$real"
            echo "🔁 $link/ -> $real/ (moved content into $real/)"
        elif is_empty_dir "$link"; then
            rmdir "$link"
        else
            echo "⚠️  WARN: both $real/ and $link/ are non-empty directories, merge manually."
            return 0
        fi
    fi

    [ -d "$real" ] || return 0
    create_symlink "$real" "$link"
}

create_symlink() {
    local target="$1"   # existing file or dir (e.g. AGENTS.md)
    local link="$2"     # symlink to create (e.g. CLAUDE.md)

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

converge_file "AGENTS.md" "CLAUDE.md"
converge_dir  ".agents"   ".claude"
create_symlink ".gitignore" ".geminiignore"

echo "── done ──"
