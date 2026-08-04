# 開發環境 (Development Guide)

`CLAUDE.md` 開發指南的展開細節; 指令一覽見 [`docs/cli.md`](cli.md)。

## 前置需求 (Prerequisites)

- Go 1.26.3+
- SQLite3
- Ollama（用於 LLM 提取，預設 `http://localhost:11434`）
- `mempalace` CLI（用於事實寫入）
- `jq`（用於 hook 腳本解析 JSON）
- `marksman`（選用；`plugins/general/.lsp.json` 的 Markdown LSP）
- `codegraph`（選用；`plugins/explore/.mcp.json` 的 MCP server）

## 安裝 (Installation)

```bash
# 複製專案（含 vendored submodule）
git clone https://github.com/bizshuk/cc-plugin.git
git submodule update --init --recursive

# 初始化環境（建立軟連結、同步設定）
chmod +x scripts/run.sh && ./scripts/run.sh

# 安裝為 Claude Code 插件
claude --plugin-dir .
```

## 部署 (Deploy)

```bash
# 安裝至 $GOPATH/bin
go install

# 排程執行（每日 03:00）
crontab -e
# 加入: 0 3 * * * $HOME/go/bin/cc-plugin distill >> $HOME/.config/cc-plugin/logs/run.log 2>&1
```

## 設定同步範圍 (Config Sync Scope)

- `scripts/run.sh` 只軟連結 `config/*` 與 `pkg/*` 設定檔至各 agent 家目錄
  （`$HOME/.claude`、`$HOME/.gemini`、`$HOME/.hermes` 等），並同步外部工具設定
  （LiteLLM、CCStatusline、Tokscale）、建立 `tmp/` 反向連結以供調試。
- 切換供應商：把 `~/.claude/settings.json` 的連結目標改指向對應的
  `config/<provider>.json`。
- symlink 的兩類例外，一律不由 `run.sh` 連結：
  - `skill` 由 `skills add` 安裝至 `~/.agents/skills`
  - `全域規則 (config/CLAUDE.global.md)` 由 `skills install` 複製成實體檔至各
    agent 設定目錄
