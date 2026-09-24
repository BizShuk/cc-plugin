#!/usr/bin/env bash
# hermes:install — 冪等 (idempotent) 安裝上游 Hermes，不碰設定、不啟動任何服務。
#
#   ./scripts/setup.sh            缺 hermes 才以上游安裝器安裝
#   ./scripts/setup.sh --upgrade  另外以 `hermes update` 升級既有安裝
#
# 設定樣板同步與 tmp/config 連結是 scripts/run.sh（run:setup）的事。

set -euo pipefail

INSTALLER_URL="https://raw.githubusercontent.com/NousResearch/hermes-agent/main/scripts/install.sh"

HERMES_HOME="${HERMES_HOME:-$HOME/.hermes}"
export HERMES_HOME

UPGRADE=false
case "${1:-}" in
    "") ;;
    --upgrade) UPGRADE=true ;;
    *) echo "用法：$0 [--upgrade]" >&2; exit 2 ;;
esac

# 上游把指令連結放在 ~/.local/bin；新 shell 才會吃到安裝器改的 PATH。
export PATH="$HOME/.local/bin:$PATH"

if command -v hermes >/dev/null 2>&1; then
    if [ "$UPGRADE" = true ]; then
        echo "升級 Hermes：hermes update"
        hermes update
    else
        echo "Hermes 已安裝，略過（升級請帶 --upgrade）"
    fi
    exit 0
fi

echo "安裝 Hermes：$INSTALLER_URL"
# --skip-setup：模型與通道設定是互動式的 `hermes setup`，不在安裝腳本裡跑。
curl -fsSL "$INSTALLER_URL" | bash -s -- --skip-setup --hermes-home "$HERMES_HOME"
command -v hermes >/dev/null 2>&1 || {
    echo "錯誤：安裝完成但找不到 hermes 指令（HERMES_HOME=$HERMES_HOME）" >&2
    exit 1
}
echo "下一步：pnpm run run:setup 同步設定（cc-plugin 根層）；首次再跑 hermes setup 與 hermes gateway setup"
