#!/usr/bin/env bash

set -euo pipefail

REPO_ROOT="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd -P)"
TMP_DIR="$REPO_ROOT/tmp"
STATUSLINE_DIR="$REPO_ROOT/plugins/default/config/statusline"

backup_and_link() {
    local source=$1
    local target=$2
    local backup="${target}.bak"

    mkdir -p "$(dirname -- "$target")"
    if [[ -L "$target" ]]; then
        ln -sfn -- "$source" "$target"
        return 0
    fi
    if [[ -e "$target" ]]; then
        if [[ ! -e "$backup" && ! -L "$backup" ]]; then
            mv -- "$target" "$backup"
        else
            rm -rf -- "$target"
        fi
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
    "$HOME/.gemini/antigravity-cli" \
    "$HOME/.claude-mem" \
    "$HOME/.hermes" \
    "$HOME/.grok" \
    "$HOME/.local/bin" \
    "$TMP_DIR"

# Global rule (plugins/default/AGENTS.md): link into each agent CLI found on PATH.
GLOBAL_RULE="$REPO_ROOT/plugins/default/AGENTS.md"
GLOBAL_RULE_TARGETS=(
    "claude:$HOME/.claude/CLAUDE.md"
    "agy:$HOME/.gemini/GEMINI.md"
    "codex:$HOME/.codex/AGENTS.md"
    "pi:$HOME/.pi/agent/AGENTS.md"
    "hermes:$HOME/.hermes/AGENTS.md"
    "grok:$HOME/.grok/rules/AGENTS.md"
)
for entry in "${GLOBAL_RULE_TARGETS[@]}"; do
    if command -v "${entry%%:*}" >/dev/null 2>&1; then
        backup_and_link "$GLOBAL_RULE" "${entry#*:}"
    fi
done
# Skills are NOT linked here -- `skills add` owns ~/.agents/skills.

backup_and_link "$REPO_ROOT/plugins/default/config/settings/claude.json" "$HOME/.claude/settings.json"
backup_and_link "$REPO_ROOT/plugins/default/config/settings/claude.keybinding.json" "$HOME/.claude/keybindings.json"
backup_and_link "$REPO_ROOT/plugins/default/config/settings/codex.toml" "$HOME/.codex/config.toml"
backup_and_link "$REPO_ROOT/plugins/default/config/settings/grok.toml" "$HOME/.grok/config.toml"
# Grok has no outputStyle setting. ~/.grok/rules/*.md (this directory only) loads as global rules.
backup_and_link \
    "$REPO_ROOT/plugins/default/output-styles/executive_summary.md" \
    "$HOME/.grok/rules/executive_summary.md"
backup_and_link "$REPO_ROOT/plugins/default/config/settings/antigravity-cli.json" "$HOME/.gemini/antigravity-cli/settings.json"
backup_and_link "$REPO_ROOT/plugins/default/config/settings/antigravity-cli.keybinding.json" "$HOME/.gemini/antigravity-cli/keybindings.json"
backup_and_link \
    "$REPO_ROOT/plugins/default" \
    "$HOME/.gemini/config/plugins/default"

# Claude Code: register this repo as a local marketplace and install its plugins.
# Local plugins are read from their folder; re-running is idempotent (add/install once, then update).
MARKETPLACE="shuk-cc-plugin"
CLAUDE_PLUGINS=(default)

install_claude_plugins() {
    if ! command -v claude >/dev/null 2>&1; then
        echo "skip Claude plugins: claude CLI not found" >&2
        return 0
    fi
    if claude plugin marketplace list | grep -q "❯ $MARKETPLACE\$"; then
        claude plugin marketplace update "$MARKETPLACE"
    else
        claude plugin marketplace add "$REPO_ROOT"
    fi
    local installed plugin
    installed=$(claude plugin list --json)
    for plugin in "${CLAUDE_PLUGINS[@]}"; do
        if grep -q "\"$plugin@$MARKETPLACE\"" <<<"$installed"; then
            claude plugin update "$plugin@$MARKETPLACE"
        else
            claude plugin install "$plugin@$MARKETPLACE"
        fi
    done
}

install_claude_plugins

# Hermes is NOT set up here -- run it explicitly: `pnpm run hermes:setup`
# (install / upgrade / webui are likewise explicit `hermes:*` tasks).

install_sample_if_missing \
    "$REPO_ROOT/plugins/experiment/skills/summarize-sh/config.sample.json" \
    "$HOME/.summarize/config.json"

backup_and_link \
    "$STATUSLINE_DIR/claude-statusline.json" \
    "$HOME/.config/ccstatusline/settings.json"
chmod +x "$STATUSLINE_DIR/agy-statusline"
backup_and_link \
    "$STATUSLINE_DIR/agy-statusline" \
    "$HOME/.local/bin/agy-statusline"
chmod +x "$STATUSLINE_DIR/grok-statusline"
backup_and_link \
    "$STATUSLINE_DIR/grok-statusline" \
    "$HOME/.grok/statusline.sh"
backup_and_link \
    "$REPO_ROOT/plugins/default/config/settings/tokscale.json" \
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
    "$HOME/.grok" \
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
