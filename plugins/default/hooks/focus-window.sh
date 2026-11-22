#!/usr/bin/env bash
# Notification click action: bring the app window the agent ran in to the front.
#
# Usage: focus-window.sh <host-app-bundle-id> <agent-cwd>
#
# Antigravity IDE: focus the window whose workspace folder is the deepest one
# containing <agent-cwd>; opening an already-open folder focuses its window
# instead of opening a new one. The open windows come from the IDE's
# windowsState in storage.json. Any other host app is just activated.

set -euo pipefail

IDE_BUNDLE=com.google.antigravity-ide
IDE_STORAGE="$HOME/Library/Application Support/Antigravity IDE/User/globalStorage/storage.json"

bundle=$1
cwd=$2

if [[ $bundle == "$IDE_BUNDLE" && -f $IDE_STORAGE ]]; then
    folder=$(/usr/bin/jq -r --arg cwd "$cwd" '
        [.windowsState.openedWindows[]?.folder // empty
         | sub("^file://"; "")
         | select(. as $f | $cwd == $f or ($cwd | startswith($f + "/")))]
        | max_by(length) // empty' "$IDE_STORAGE")
    if [[ -n $folder ]]; then
        exec /usr/bin/open -b "$IDE_BUNDLE" "$folder"
    fi
fi

exec /usr/bin/open -b "$bundle"
