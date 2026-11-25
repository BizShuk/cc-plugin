#!/usr/bin/env bash
# PostToolUse hook: auto-fix and lint Markdown files created or edited by the agent.
#
# One script, three hosts:
#   Claude Code, Grok  -> hooks/hooks.json (Claude format), host detected from $GROK_HOOK_EVENT
#   Antigravity (agy)  -> hooks.json at the plugin root, called with the argument `agy`
#
# Behavior:
#   - Ensures `markdownlint` (`markdownlint-cli`) is installed globally; installs it via npm if missing.
#   - Resolves the target file path from the host's PostToolUse JSON payload on stdin.
#   - Skips non-Markdown, missing, or vendored (.git / node_modules) files.
#   - Respects a repo's own .markdownlint* config if present; otherwise uses hooks/markdownlint.json.
#   - Runs `markdownlint --fix`, then checks for remaining unfixable issues.
#   - Claude / Grok: returns `hookSpecificOutput.additionalContext` if issues remain.
#   - agy: outputs `{}` on stdout (agy PostToolUse contract does not feed text back to the model).

set -euo pipefail

SCRIPT_DIR="$(cd -- "$(dirname -- "$0")" && pwd -P)"
DEFAULT_CONFIG="$SCRIPT_DIR/markdownlint.json"

# Hooks often launch with a minimal PATH; expose pnpm/npm/homebrew/nvm global bins.
export PATH="$HOME/Library/pnpm/bin:$HOME/.local/bin:/opt/homebrew/bin:/usr/local/bin:$PATH"
if [[ -d "$HOME/.nvm/versions/node" ]]; then
    latest_node=$(ls "$HOME/.nvm/versions/node" 2>/dev/null | sed 's/^v//' | sort -t. -k1,1n -k2,2n -k3,3n | tail -1 || true)
    if [[ -n $latest_node && -d "$HOME/.nvm/versions/node/v${latest_node}/bin" ]]; then
        export PATH="$HOME/.nvm/versions/node/v${latest_node}/bin:$PATH"
    fi
fi

input=$(cat)

host=${1:-}
if [[ -z $host ]]; then
    if [[ -n ${GROK_HOOK_EVENT:-} ]]; then host=grok; else host=claude; fi
fi

finish_agy() {
    if [[ $host == agy ]]; then
        echo '{}'
    fi
    exit 0
}

if ! command -v jq >/dev/null 2>&1; then
    finish_agy
fi

field() { jq -r "$1 // empty" <<<"$input"; }

case $host in
claude)
    file_path=$(field '.tool_input.file_path // .tool_input.path')
    ;;
grok)
    file_path=$(field '.tool_input.file_path // .tool_input.path // .toolInput.file_path // .toolInput.path // .tool_input.target_file // .toolInput.target_file')
    ;;
agy)
    file_path=$(field '.toolCall.args.Path // .toolCall.args.TargetFile // .toolCall.args.path')
    ;;
*)
    echo "markdownlint hook: unknown host '$host'" >&2
    exit 0
    ;;
esac

if [[ -z $file_path ]]; then
    finish_agy
fi

# Resolve relative paths against the session's cwd if provided.
if [[ $file_path != /* ]]; then
    cwd=$(field '.cwd // .workspaceRoot')
    if [[ -n $cwd ]]; then
        file_path="$cwd/$file_path"
    fi
fi

case ${file_path##*.} in
md|MD|markdown|MARKDOWN) ;;
*) finish_agy ;;
esac

case $file_path in
*/.git/*|*/node_modules/*) finish_agy ;;
esac

[[ -f $file_path ]] || finish_agy

# Ensure markdownlint-cli is installed globally without falling back to npx.
if ! command -v markdownlint >/dev/null 2>&1; then
    if command -v npm >/dev/null 2>&1; then
        npm install -g markdownlint-cli >/dev/null 2>&1 || true
    fi
fi

if ! command -v markdownlint >/dev/null 2>&1; then
    finish_agy
fi

# Prefer project-local markdownlint config if one exists in an ancestor directory.
has_project_config() {
    local dir
    dir=$(cd -- "$(dirname -- "$1")" 2>/dev/null && pwd -P) || return 1
    while true; do
        for name in \
            .markdownlint.json \
            .markdownlint.jsonc \
            .markdownlint.yaml \
            .markdownlint.yml \
            .markdownlint.cjs \
            .markdownlint.mjs \
            .markdownlintrc
        do
            if [[ -f "$dir/$name" ]]; then
                return 0
            fi
        done
        [[ $dir == / ]] && break
        dir=$(dirname -- "$dir")
    done
    return 1
}

cfg_args=()
if ! has_project_config "$file_path" && [[ -f $DEFAULT_CONFIG ]]; then
    cfg_args=(-c "$DEFAULT_CONFIG")
fi

markdownlint "${cfg_args[@]}" --fix -- "$file_path" >/dev/null 2>&1 || true

lint_out=""
if ! lint_out=$(markdownlint "${cfg_args[@]}" -- "$file_path" 2>&1); then
    if [[ $host == agy ]]; then
        finish_agy
    fi
    if [[ -n $lint_out ]]; then
        jq -n --arg ctx "markdownlint found remaining issues in $file_path (auto-fix already ran):
$lint_out" '{
            hookSpecificOutput: {
                hookEventName: "PostToolUse",
                additionalContext: $ctx
            }
        }'
    fi
    exit 0
fi

finish_agy
