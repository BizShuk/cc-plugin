#!/usr/bin/env bash
#
# aliases.sh — 各家 LLM CLI 的 shell alias 入口 (LLM CLI shell aliases)
#
# 用法 (Usage): 在 ~/.bash_aliases 內 source 本檔案。
#   source ~/projects/ai/cc-plugin/scripts/aliases.sh
#
# 路徑自我定位 (self-locating)：CC_PLUGIN_DIR 由本檔實際位置推導，
# repo 搬家後不需回頭改 alias 字串。
#
# Token 變數 (TIKTOK_API_KEY, TIKTOK_API_KEY2, PROXY_API_KEY 等) 存於
# git-ignored 的 ~/.bash_local，alias 內只引用變數名，token 不入 git。
#
# 基礎 claudew / claudem 已提升為 ~/bin/claudew 與 ~/bin/claudem 實體
# script file（取代 alias 以避免 alias 無法接參數的限制）。

CC_PLUGIN_DIR="$(CDPATH= cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd -P)"
export CC_PLUGIN_DIR

# [Antigravity]
alias agy="agy --dangerously-skip-permissions"

# [Claude] — Anthropic 原生端 (native endpoint)
alias claude='claude --allow-dangerously-skip-permissions --settings '"$CC_PLUGIN_DIR"'/plugins/default/config/settings/claude.json'

# [Claude] — 代理端 (proxy / third-party endpoints)
# alias claudep='ANTHROPIC_AUTH_TOKEN=$PROXY_API_KEY claude --allow-dangerously-skip-permissions --settings '"$CC_PLUGIN_DIR"'/plugins/default/config/settings/claude.proxy.json'

