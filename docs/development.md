# 開發環境 (Development Guide)

`AGENTS.md` 開發指南的展開細節; 指令一覽見 [`docs/cli.md`](cli.md)。

## 前置需求 (Prerequisites)

- Go 1.26.3+
- SQLite3
- `skills` CLI（`go install github.com/bizshuk/skills@latest`，用於全域規則安裝與技能管理）
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

# 安裝技能 (全域規則已由 run.sh 連結)
skills add .

# 安裝為 Claude Code 插件
claude --plugin-dir .
```

## 部署 (Deploy)

```bash
# 安裝至 $GOPATH/bin
go install
```

## 設定同步範圍 (Config Sync Scope)

- `scripts/run.sh` 只軟連結 `config/*` 與 `pkg/*` 設定檔至各 agent 家目錄
  （`$HOME/.claude`、`$HOME/.gemini`、`$HOME/.hermes` 等），並同步外部工具設定
  （LiteLLM、CCStatusline、Tokscale）、建立 `tmp/` 反向連結以供調試。
- Hermes 例外：`~/.hermes` 的記憶檔、樣板與 `.env` key 由子專案 `pkg/hermes/scripts/run.sh`
  負責，根層 `run.sh` `不`叫用它，須明確執行 `pnpm run hermes:setup`；安裝上游 Hermes 另跑 `pnpm run hermes:install`（WebUI：`pnpm run hermes:install:webui`）。
- 切換供應商：把 `~/.claude/settings.json` 的連結目標改指向對應的
  `plugins/default/config/settings/claude.<provider>.json`。
- Grok 全域規則目錄是 `~/.grok/rules/`（只掃該層 `*.md`）。
  `plugins/default/output-styles/executive_summary.md` 由 `run.sh` 軟連結至
  `~/.grok/rules/executive_summary.md`。
- 雙相容插件 `plugins/default/`：Claude Code 由 `run.sh` 登錄 marketplace `shuk-cc-plugin` 並安裝；
  Antigravity 由 `run.sh` 軟連結至 `~/.gemini/config/plugins/default`，其 `rules/output-style.md` 軟連結至
  `output-styles/executive_summary.md`（可用 `agy plugin enable/disable default` 開關）。
- `全域規則 (plugins/default/AGENTS.md)` 由 `run.sh` 軟連結至各 agent 設定目錄，只連結 PATH 上偵測到的 CLI
  （`claude`、`agy`、`codex`、`pi`、`hermes`、`grok`）。
- symlink 的例外：`skill` 由 `skills add` 安裝至 `~/.agents/skills`，不由 `run.sh` 連結。
