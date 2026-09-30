#!/usr/bin/env bash
# hermes:gateway — 啟動 Hermes Gateway 前台常駐程序（供 PM2 或前台監控使用）。
#
#   ./pkg/hermes/scripts/gateway.sh          前台啟動 Gateway（宣告 --external-supervisor）
#   ./pkg/hermes/scripts/gateway.sh --force  若 launchd 仍在執行，強制以此程序接管
#
# 管理方式由 PM2 (ecosystem.config.js) 常駐監控。
# 依據 inf-spec 契約，執行期載入指定 Python 虛擬環境 (/Users/shuk/.venv)。

set -euo pipefail

HERMES_HOME="${HERMES_HOME:-$HOME/.hermes}"
VENV_PATH="${VENV_PATH:-/Users/shuk/.venv}"
export HERMES_HOME

# 1. 載入指定的 Python 虛擬環境 (如存在)
if [ -f "$VENV_PATH/bin/activate" ]; then
    # shellcheck source=/dev/null
    source "$VENV_PATH/bin/activate"
fi

# 2. 確保執行路徑包含 ~/.local/bin 與 venv/bin
export PATH="$HOME/.local/bin:$VENV_PATH/bin:$PATH"

# 3. 確保 hermes 指令可用
if ! command -v hermes >/dev/null 2>&1; then
    echo "錯誤：找不到 hermes 指令（請先執行 pnpm run hermes:install）" >&2
    exit 1
fi

# 4. 檢查 launchd 是否仍在守護同一個 gateway 實例
if launchctl list ai.hermes.gateway >/dev/null 2>&1; then
    has_force=false
    for arg in "$@"; do
        if [ "$arg" = "--force" ]; then
            has_force=true
            break
        fi
    done
    if [ "$has_force" = false ]; then
        echo "注意：macOS launchd 正在守護 ai.hermes.gateway。" >&2
        echo "若要完全改由 PM2 接管，建議先卸載 launchd 服務：" >&2
        echo "  hermes gateway stop" >&2
        echo "  launchctl unload ~/Library/LaunchAgents/ai.hermes.gateway.plist" >&2
        echo "或帶 --force 強制前台啟動：" >&2
        echo "  $0 --force" >&2
    fi
fi

# 5. 前台啟動 Hermes Gateway，宣告外部 supervisor（由 PM2 接管生命週期）
exec hermes gateway run --external-supervisor "$@"
