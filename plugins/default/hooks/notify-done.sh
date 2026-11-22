#!/usr/bin/env bash
# Stop hook: post a macOS notification "<agent> <project>" with a one-line recap of
# the agent's final response; clicking it focuses the window the agent ran in
# (focus-window.sh). Requires terminal-notifier.
#
# One script, three hosts:
#   Claude Code, Grok  -> hooks/hooks.json (Claude format), host detected from $GROK_HOOK_EVENT
#   Antigravity (agy)  -> hooks.json at the plugin root, called with the argument `agy`
#
# Each host hands the final response over differently:
#   Claude  stdin .last_assistant_message
#   Grok    stdin .lastAssistantMessage
#   agy     last PLANNER_RESPONSE with text in the .transcriptPath JSONL

set -euo pipefail
export LC_ALL=en_US.UTF-8

RECAP_CHARS=120

[[ "$(uname -s)" == Darwin ]] || exit 0

input=$(cat)
field() { jq -r "$1 // empty" <<<"$input"; }

host=${1:-}
if [[ -z $host ]]; then
    if [[ -n ${GROK_HOOK_EVENT:-} ]]; then host=grok; else host=claude; fi
fi

case $host in
claude)
    app="Claude Code"
    message=$(field .last_assistant_message)
    ;;
grok)
    # Grok also fires Stop at session teardown; that is not a finished turn.
    case $(field .reason) in channel_closed | shutdown) exit 0 ;; esac
    app="Grok"
    message=$(field .lastAssistantMessage)
    ;;
agy)
    # agy expects a JSON object on stdout; `{}` lets the agent stop.
    echo '{}'
    [[ $(field .terminationReason) == USER_CANCELED ]] && exit 0
    app="Antigravity"
    message=""
    transcript=$(field .transcriptPath)
    if [[ -f $transcript ]]; then
        message=$(jq -rs '
            map(select(.type == "PLANNER_RESPONSE" and (.content | type) == "string" and .content != ""))
            | last | .content // empty' "$transcript")
    fi
    ;;
*)
    echo "notify-done: unknown host '$host'" >&2
    exit 0
    ;;
esac

# Recap = the first text line of the response. The executive_summary output
# style opens every response with a single-line recap; for other responses this
# falls back to the first line that is not a heading, code block, table or blank.
# List markers, emphasis, inline code and link targets are stripped.
recap=$(printf '%s\n' "$message" |
    awk '
        /^[[:space:]]*```/ { fence = !fence; next }
        fence || /^[[:space:]]*(#|\||$)/ { next }
        { print; exit }' |
    sed -E \
        -e 's/^[[:space:]]*(>[[:space:]]*)?([-*+]|[0-9]+\.)?[[:space:]]*(\[[ xX]\][[:space:]]+)?//' \
        -e 's/\[([^]]*)\]\([^)]*\)/\1/g' \
        -e 's/(\*\*|__|`)//g')
recap=${recap:-回應完成}
((${#recap} > RECAP_CHARS)) && recap="${recap:0:RECAP_CHARS}…"

cwd=$(field .cwd)
project=$(basename -- "$cwd")

# Clicking the notification focuses the app (window) hosting the agent. macOS
# sets __CFBundleIdentifier for processes started from an app, so the host
# terminal or IDE is inherited through the agent CLI.
click=()
if [[ -n ${__CFBundleIdentifier:-} ]]; then
    focus="$(cd -- "$(dirname -- "$0")" && pwd -P)/focus-window.sh"
    click=(-execute "$(printf '%q %q %q' "$focus" "$__CFBundleIdentifier" "$cwd")")
fi

terminal-notifier \
    -title "$app${project:+ $project}" \
    -message "$recap" \
    -group "$host:$cwd" \
    ${click[@]+"${click[@]}"} >/dev/null 2>&1 || true
