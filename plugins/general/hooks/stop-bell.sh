#!/bin/bash
# stop-bell.sh — Terminal bell notification for Claude Code Stop hook.
#
# Strategy: Claude Code runs hooks as detached child processes (no inherited
# TTY), so writing BEL to stdout/stderr hits a pipe and never rings the bell.
# Workaround: locate the parent process's TTY device and write BEL directly
# to that device. Fallbacks: stdout, stderr, macOS system beep.

set -u

# Consume stdin if present (Claude Code pipes JSON; we don't need it).
if [ ! -t 0 ]; then
    cat > /dev/null
fi

BEL_TARGET="none"

# 1) Find Claude Code's TTY via PPID, then write BEL to /dev/<tty>.
#    `ps -o tty= -p $PPID` returns the controlling TTY name (e.g., ttys003).
if command -v ps >/dev/null 2>&1; then
    PARENT_TTY=$(ps -o tty= -p "$PPID" 2>/dev/null | tr -d ' \n')
    if [ -n "$PARENT_TTY" ] && [ "$PARENT_TTY" != "??" ] && [ -e "/dev/$PARENT_TTY" ]; then
        if (printf '\a' > "/dev/$PARENT_TTY") 2>&-; then
            BEL_TARGET="/dev/$PARENT_TTY"
        fi
    fi
fi

# 2) stdout — TTY-bound in some Claude Code versions.
if [ "$BEL_TARGET" = "none" ]; then
    if printf '\a' 2>/dev/null; then
        BEL_TARGET="stdout"
    fi
fi

# 3) stderr — last-resort.
if [ "$BEL_TARGET" = "none" ]; then
    printf '\a' >&2 2>/dev/null || true
    BEL_TARGET="stderr"
fi

# 4) macOS system beep via osascript (audio fallback, ignores TTY).
if [ "$BEL_TARGET" = "stderr" ] && [ "$(uname -s 2>/dev/null)" = "Darwin" ]; then
    if command -v osascript >/dev/null 2>&1; then
        osascript -e 'beep' >/dev/null 2>&1 || true
        BEL_TARGET="osascript"
    fi
fi

exit 0
