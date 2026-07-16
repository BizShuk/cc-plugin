#!/usr/bin/env bash

set -euo pipefail

REPO_ROOT="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)"
TMP_DIR="$REPO_ROOT/tmp"

backup_and_link() {
    local source=$1
    local target=$2
    local backup="${target}.bak"

    mkdir -p "$(dirname -- "$target")"
    if [[ -e "$target" && ! -L "$target" ]]; then
        if [[ -e "$backup" || -L "$backup" ]]; then
            printf 'Refusing to overwrite existing backup: %s\n' "$backup" >&2
            return 1
        fi
        mv -- "$target" "$backup"
    fi
    ln -sfn -- "$source" "$target"
}

link_to_tmp_if_present() {
    local source=$1
    if [[ -e "$source" || -L "$source" ]]; then
        ln -sfn -- "$source" "$TMP_DIR/$(basename -- "$source")"
    fi
}

install_sample_if_missing() {
    local sample=$1
    local target=$2
    mkdir -p "$(dirname -- "$target")"
    if [[ ! -e "$target" ]]; then
        cp -- "$sample" "$target"
    fi
}

mkdir -p \
    "$HOME/.claude" \
    "$HOME/.codex" \
    "$HOME/.gemini" \
    "$HOME/.claude-mem" \
    "$HOME/.hermes" \
    "$TMP_DIR"

backup_and_link "$REPO_ROOT/config/CLAUDE.global.md" "$HOME/.claude/CLAUDE.md"
backup_and_link "$REPO_ROOT/config/CLAUDE.global.md" "$HOME/.gemini/GEMINI.md"
backup_and_link "$REPO_ROOT/config/CLAUDE.global.md" "$HOME/.codex/AGENTS.md"
backup_and_link "$REPO_ROOT/config/settings.json" "$HOME/.claude/settings.json"
backup_and_link "$REPO_ROOT/config/config.toml" "$HOME/.codex/config.toml"

backup_and_link "$REPO_ROOT/config/CLAUDE.global.md" "$HOME/.hermes/AGENTS.md"
backup_and_link "$REPO_ROOT/pkg/hermes/MEMORY.md" "$HOME/.hermes/MEMORY.md"
backup_and_link "$REPO_ROOT/pkg/hermes/USER.md" "$HOME/.hermes/USER.md"

install_sample_if_missing \
    "$REPO_ROOT/pkg/litellm/litellm_config.sample.yaml" \
    "$HOME/.config/litellm/litellm_config.yaml"
install_sample_if_missing \
    "$REPO_ROOT/plugins/experiment/skills/summarize-sh/config.sample.json" \
    "$HOME/.summarize/config.json"

backup_and_link \
    "$REPO_ROOT/pkg/ccstatusline/settings.json" \
    "$HOME/.config/ccstatusline/settings.json"
backup_and_link \
    "$REPO_ROOT/pkg/usage/tokscale/settings.json" \
    "$HOME/.config/tokscale/settings.json"

for path in \
    "$HOME/.gemini" \
    "$HOME/.claude" \
    "$HOME/.codex" \
    "$HOME/.claude-mem" \
    "$HOME/.claude.json" \
    "$HOME/.cli-proxy-api" \
    "$HOME/.auth2api" \
    "$HOME/.hermes" \
    "$HOME/.gbrain" \
    "$HOME/.mempalace" \
    "$HOME/.config/cc-plugin" \
    "$HOME/.agentmemory" \
    "$HOME/.paperclip" \
    "$HOME/.config/opencode" \
    "$HOME/.config/litellm/litellm_config.yaml" \
    "$HOME/.summarize" \
    "$HOME/.config/ccstatusline" \
    "$HOME/.config/tokscale"
do
    link_to_tmp_if_present "$path"
done
