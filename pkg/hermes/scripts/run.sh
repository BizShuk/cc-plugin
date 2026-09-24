#!/usr/bin/env bash
# run:setup — 冪等 (idempotent) 設定前置作業：連結人設、記憶與 config.yaml、同步樣板、建立 tmp/config 連結。
# 人設由 HERMES_PERSONA 指定（預設 yuna-v2），SOUL.md 以 symlink 指向 persona/<name>/SOUL.md。
# 不安裝 Hermes（那是 scripts/setup.sh / hermes:install 的事），也不啟動任何服務。
#
# Hermes 家目錄由上游固定為 ~/.hermes（可用 HERMES_HOME 覆寫，與上游安裝器同義）。
# repo 的 config/ 與家目錄同形：缺檔才複製，已存在的檔案一律不動；
# config/.env.example 只把 .env 缺少的 key 補上（空值），既有值絕不覆寫。

set -euo pipefail

PROJECT_ROOT="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)"
TEMPLATE_DIR="$PROJECT_ROOT/config"
HERMES_HOME="${HERMES_HOME:-$HOME/.hermes}"
WORKSPACE_LINK="$PROJECT_ROOT/tmp/config"
HERMES_PERSONA="${HERMES_PERSONA:-yuna-v2}"
export HERMES_HOME

# 版控檔以 symlink 放進家目錄；既有實體檔先備份成 .bak。
link_file() {
    local src="$1" dst="$2"
    if [ -e "$dst" ] && [ ! -L "$dst" ] && [ ! -e "$dst.bak" ]; then
        mv -- "$dst" "$dst.bak"
        echo "已備份：$dst.bak"
    fi
    ln -sfn -- "$src" "$dst"
}

# 版控的 config/config.yaml、config/memory/ 與 config/workspace/（→ ~/workspace）。
link_tracked() {
    link_file "$TEMPLATE_DIR/config.yaml" "$HERMES_HOME/config.yaml"
    link_file "$TEMPLATE_DIR/memory" "$HERMES_HOME/memories"
    link_file "$TEMPLATE_DIR/workspace" "$HOME/workspace"
}

# 啟用中的人設：persona/<HERMES_PERSONA>/SOUL.md → $HERMES_HOME/SOUL.md。
link_persona() {
    local soul="$PROJECT_ROOT/persona/$HERMES_PERSONA/SOUL.md"
    if [ ! -f "$soul" ]; then
        echo "找不到人設：$soul" >&2
        exit 1
    fi
    link_file "$soul" "$HERMES_HOME/SOUL.md"
    echo "已啟用人設：$HERMES_PERSONA"
}

# config/<path> → $HERMES_HOME/<path>，缺檔才複製。
sync_templates() {
    [ -d "$TEMPLATE_DIR" ] || return 0
    local src rel dst
    while IFS= read -r -d '' src; do
        rel="${src#"$TEMPLATE_DIR"/}"
        case "$rel" in .env.example | config.yaml | memory/* | workspace/*) continue ;; esac
        dst="$HERMES_HOME/$rel"
        if [ ! -e "$dst" ]; then
            mkdir -p "$(dirname -- "$dst")"
            cp "$src" "$dst"
            echo "已複製樣板：$dst"
        fi
    done < <(find "$TEMPLATE_DIR" -type f -print0)
}

# config/.env.example 的 KEY= 在 .env 缺少時才補上空值；權限固定 0600。
merge_env_keys() {
    local example="$TEMPLATE_DIR/.env.example" env_file="$HERMES_HOME/.env" key
    touch "$env_file"
    chmod 600 "$env_file"
    [ -f "$example" ] || return 0
    while IFS= read -r key; do
        if ! grep -qE "^[[:space:]]*(export[[:space:]]+)?${key}=" "$env_file"; then
            printf '%s=\n' "$key" >>"$env_file"
            echo "已補上 .env key：$key"
        fi
    done < <(sed -nE 's/^[[:space:]]*(export[[:space:]]+)?([A-Za-z_][A-Za-z0-9_]*)=.*/\2/p' "$example")
}

mkdir -p "$HERMES_HOME" "$PROJECT_ROOT/tmp"
link_tracked
link_persona
sync_templates
merge_env_keys
ln -sfn "$HERMES_HOME" "$WORKSPACE_LINK"
echo "已建立軟連結：$WORKSPACE_LINK -> $HERMES_HOME"
