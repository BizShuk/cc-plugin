#!/usr/bin/env bash
# hermes:install:webui — 冪等 (idempotent) 安裝上游 Hermes WebUI，不啟動服務。
#
#   ./scripts/setup-webui.sh            缺 WebUI 才 clone
#   ./scripts/setup-webui.sh --upgrade  另外 fast-forward 既有 clone
#
# Python 相依由 WebUI 首次啟動時的 bootstrap 自行解析（優先沿用 Hermes agent venv）。

set -euo pipefail

REPO_URL="https://github.com/nesquena/hermes-webui.git"
HERMES_HOME="${HERMES_HOME:-$HOME/.hermes}"
WEBUI_DIR="$HERMES_HOME/hermes-webui"

UPGRADE=false
case "${1:-}" in
    "") ;;
    --upgrade) UPGRADE=true ;;
    *) echo "用法：$0 [--upgrade]" >&2; exit 2 ;;
esac

export PATH="$HOME/.local/bin:$PATH"
command -v hermes >/dev/null 2>&1 || {
    echo "錯誤：找不到 hermes，先跑 pnpm run hermes:install" >&2
    exit 1
}

if [ -d "$WEBUI_DIR/.git" ]; then
    if [ "$UPGRADE" = true ]; then
        echo "升級 Hermes WebUI：$WEBUI_DIR"
        git -C "$WEBUI_DIR" pull --ff-only
    else
        echo "Hermes WebUI 已安裝，略過（升級請帶 --upgrade）"
    fi
else
    echo "安裝 Hermes WebUI：$REPO_URL → $WEBUI_DIR"
    git clone --depth 1 "$REPO_URL" "$WEBUI_DIR"
fi
echo "啟動：$WEBUI_DIR/ctl.sh start（預設 http://127.0.0.1:8787）"
