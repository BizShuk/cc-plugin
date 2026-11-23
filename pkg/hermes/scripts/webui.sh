#!/usr/bin/env bash
# hermes:webui — 前台啟動 Hermes WebUI（供 PM2 常駐監控使用）。
#
#   ./pkg/hermes/scripts/webui.sh    前台啟動 WebUI（bootstrap --foreground, 預設 http://127.0.0.1:8787）
#
# 管理方式由 PM2 (ecosystem.config.js) 常駐監控；Python 沿用指定虛擬環境 (/Users/shuk/.venv)。

set -euo pipefail

HERMES_HOME="${HERMES_HOME:-$HOME/.hermes}"
VENV_PATH="${VENV_PATH:-/Users/shuk/.venv}"
WEBUI_DIR="$HERMES_HOME/hermes-webui"
export HERMES_HOME

# 載入指定的 Python 虛擬環境 (如存在)
if [ -f "$VENV_PATH/bin/activate" ]; then
    # shellcheck source=/dev/null
    source "$VENV_PATH/bin/activate"
fi

export HERMES_WEBUI_PYTHON="${HERMES_WEBUI_PYTHON:-$VENV_PATH/bin/python3}"
export PATH="$HOME/.local/bin:$VENV_PATH/bin:$PATH"

if [ ! -f "$WEBUI_DIR/bootstrap.py" ]; then
    echo "錯誤：找不到 $WEBUI_DIR/bootstrap.py（請先執行 pnpm run hermes:install:webui）" >&2
    exit 1
fi

exec "$HERMES_WEBUI_PYTHON" "$WEBUI_DIR/bootstrap.py" --foreground --no-browser "$@"
