# CC-Plugin 全域設定配置庫 (CC-Plugin Global Configuration Repository)

本專案是一個針對 `Claude Code` 與其他 AI 編碼代理的全域設定配置庫，提供集中化的設定管理、客製化插件 (Plugins)、自訂技能 (Skills) 與專屬代理 (Agents) 配置，並內建 Go 語言實作的 `claude-mem` 資料匯出與 topology 圖譜工具。

## 業務領域 (Business Domains)

### 資料匯出 (Data Export)

從 `claude-mem` 匯出觀察值，支援增量匯出（基於遊標）與全量匯出。

`領域流程 (Domain Flow):`

1. 使用者執行 `cc-plugin export claudemem`
2. 讀取遊標（增量模式）或從頭開始（`--all` 全量模式）

`核心實體 (Key Entities):` `Observation`, `Cursor`

---

### 環境初始化與配置同步 (Environment Initialization & Config Sync)

透過 `scripts/run.sh`（macOS/Unix）將本庫的設定檔與範本軟連結至使用者的家目錄資料夾（`$HOME/.claude`、`$HOME/.gemini`、`$HOME/.hermes` 等），同步外部工具設定（LiteLLM、CCStatusline、Tokscale）。切換供應商的方式是把 `~/.claude/settings.json` 的連結目標改指向對應的 `config/<provider>.json`。

`領域流程 (Domain Flow):`

1. 執行 `scripts/run.sh` → 建立家目錄結構
2. 軟連結全域設定檔（`CLAUDE.global.md`、`settings.json`）→ 至 Claude Code、Gemini CLI、Codex；
   Hermes 不在此流程內, 須另外明確執行 `pnpm run hermes:setup`（見 [`pkg/hermes/`](pkg/hermes/README.md)）
3. 同步外部工具設定並建立調試用反向連結（同步範圍見 [`docs/development.md`](docs/development.md)）

`核心實體 (Key Entities):` `Active settings`, `Provider settings`, `Global rule`（定義見 [`docs/terminology.md`](docs/terminology.md)）

---

### AI 技能與代理生態 (AI Skills & Agents Ecosystem)

提供可跨 AI 編碼代理共用的自訂技能集與專屬代理定義，劃分為八個本地模組化插件目錄；
插件的分界與清單由 [`plugins/README.md`](plugins/README.md) 單一擁有。

`領域流程 (Domain Flow):`

1. 開發者在對應的 `plugins/<name>/skills/` 目錄下建立 `SKILL.md`（符合 agentskills.io 規範）
2. 使用 `skills add .` 掃描並註冊技能至 `skills.json`，並安裝至多個 AI Agent（Antigravity、Claude Code、Gemini CLI 等）
3. `plugin.json` 的 `skills`／`agents` 只列非預設路徑或外部來源；預設目錄由自動探索，不必列舉。hooks、MCP/LSP 與其他 metadata 仍由 manifest 宣告

`核心實體 (Key Entities):` `SKILL.md`, `plugin.json`, `hooks.json`, `monitors.json`, `skills.json`

---

## 領域關聯 (Domain Relationships)

- `環境初始化` 負責將 `AI 技能與代理` 的設定檔同步至各個 AI Agent 的家目錄

## 使用方式 (Usage)

```bash
cc-plugin export claudemem             # 資料匯出（claude-mem 觀察值, 增量）
cc-plugin topology verify              # Topology 圖譜驗證
./scripts/run.sh && skills install && skills add .   # 環境初始化、全域規則與技能安裝
```

完整 CLI 參考（含 flags、增量／全量模式）見 [`docs/cli.md`](docs/cli.md)。

待辦與改善項目見 [`README.todo`](README.todo)；已淘汰功能見
[`docs/specs/2026-07-22-Summary.md`](docs/specs/2026-07-22-Summary.md) 的「已淘汰」章節。
