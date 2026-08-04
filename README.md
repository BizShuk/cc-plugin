# CC-Plugin 全域設定配置庫 (CC-Plugin Global Configuration Repository)

本專案是一個針對 `Claude Code` 與其他 AI 編碼代理的全域設定配置庫，提供集中化的設定管理、客製化插件 (Plugins)、自訂技能 (Skills) 與專屬代理 (Agents) 配置，並內建一套 Go 語言實作的記憶蒸餾管道 (Distiller Pipeline) 用於整合多個 AI 記憶來源。

## 業務領域 (Business Domains)

### 記憶蒸餾管道 (Memory Distillation Pipeline)

從多個 AI 記憶來源（`gbrain`、`claude-mem`）自動讀取觀察值、透過本地 LLM（`Ollama`）提取候選記憶，再分流寫入兩個記憶儲存庫（`agentmemory` API、`mempalace` CLI），最後清理過期資料。

`領域流程 (Domain Flow):`

1. `distill` 主命令啟動管道 → 從 `gbrain` 與 `claude-mem` 讀取新增觀察值
2. 透過本地 LLM 提取候選記憶 → 分類為 `Memory` 與 `Fact`
3. 寫入 `agentmemory` API（所有記憶）與 `mempalace mine`（通過真實性門檻的事實）
4. 更新遊標與蒸餾狀態 → 清理過期資料與檔案

`核心實體 (Key Entities):` `Observation`, `Candidate`, `Memory`, `Fact`, `Cursor`, `Seen`, `Distilled`

---

### 資料匯出 (Data Export)

提供從 `gbrain`、`claude-mem`、`mempalace` 三個儲存來源匯出原始資料的能力，支援增量匯出（基於遊標）與全量匯出。

`領域流程 (Domain Flow):`

1. 使用者執行 `cc-plugin export <子命令>` → 選擇匯出 `gbrain`、`claudemem` 或 `mempalace`
2. 讀取遊標（增量模式）或從頭開始（`--all` 全量模式）
3. `mempalace` 子命令支援類別清單（CSV）與完整 Markdown 結構匯出（`--data`）

`核心實體 (Key Entities):` `DrawerRow`, `Observation`

---

### 環境初始化與配置同步 (Environment Initialization & Config Sync)

透過 `scripts/run.sh`（macOS/Unix）將本庫的設定檔與範本軟連結至使用者的家目錄資料夾（`$HOME/.claude`、`$HOME/.gemini`、`$HOME/.hermes` 等），同步外部工具設定（LiteLLM、CCStatusline、Tokscale）。切換供應商的方式是把 `~/.claude/settings.json` 的連結目標改指向對應的 `config/<provider>.json`。

`領域流程 (Domain Flow):`

1. 執行 `scripts/run.sh` → 建立家目錄結構
2. 軟連結全域設定檔（`CLAUDE.global.md`、`settings.json`）→ 至 Claude Code、Gemini CLI、Codex、Hermes
3. 同步外部工具設定並建立調試用反向連結（同步範圍見 [`docs/development.md`](docs/development.md)）

`核心實體 (Key Entities):` `Active settings`, `Provider settings`, `Global rule`（定義見 [`docs/terminology.md`](docs/terminology.md)）

---

### AI 技能與代理生態 (AI Skills & Agents Ecosystem)

提供可跨 AI 編碼代理共用的自訂技能集與專屬代理定義，劃分為八個本地模組化插件目錄；
插件的分界與清單由 [`plugins/README.md`](plugins/README.md) 單一擁有。

`領域流程 (Domain Flow):`

1. 開發者在對應的 `plugins/<name>/skills/` 目錄下建立 `SKILL.md`（符合 agentskills.io 規範）
2. 使用 `npx skills add .` 掃描並註冊技能至 `skills.json`，並安裝至多個 AI Agent（Antigravity、Claude Code、Gemini CLI 等）
3. `plugin.json` 保留空的 `skills`／`agents` 陣列，由標準目錄自動探索；hooks、MCP/LSP 與其他 metadata 仍由 manifest 宣告

`核心實體 (Key Entities):` `SKILL.md`, `plugin.json`, `hooks.json`, `monitors.json`, `skills.json`

---

## 領域關聯 (Domain Relationships)

- `記憶蒸餾管道` 的輸出（`Memory`、`Fact`）寫入外部記憶儲存庫，而 `資料匯出` 則可從同一儲存庫反向匯出資料
- `環境初始化` 負責將 `AI 技能與代理` 的設定檔同步至各個 AI Agent 的家目錄
- `資料匯出` 與 `記憶蒸餾管道` 共用同一份遊標狀態以支援增量操作

## 使用方式 (Usage)

```bash
cc-plugin distill                      # 記憶蒸餾管道（讀取 → 提取 → 寫入 → 清理）
cc-plugin export mempalace             # 資料匯出（gbrain / claudemem / mempalace）
cc-plugin topology verify              # Topology 圖譜驗證
./scripts/run.sh && npx skills add .   # 環境初始化與技能安裝
```

完整 CLI 參考（含 flags、增量／全量模式）見 [`docs/cli.md`](docs/cli.md)。

待辦與改善項目見 [`README.todo`](README.todo)；已淘汰功能見
[`docs/specs/2026-07-22-Summary.md`](docs/specs/2026-07-22-Summary.md) 的「已淘汰」章節。
