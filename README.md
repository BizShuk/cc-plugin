# CC-Plugin 全域設定配置庫 (CC-Plugin Global Configuration Repository)

本專案是一個針對 `Claude Code` 與其他 AI 編碼代理的全域設定配置庫，提供集中化的設定管理、客製化插件 (Plugins)、自訂技能 (Skills) 與專屬代理 (Agents) 配置，並內建一套 Go 語言實作的記憶蒸餾管道 (Distiller Pipeline) 用於整合多個 AI 記憶來源。

## 業務領域 (Business Domains)

### 記憶蒸餾管道 (Memory Distillation Pipeline)

從多個 AI 記憶來源（`gbrain`、`claude-mem`）自動讀取觀察值、透過本地 LLM（`Ollama`）提取候選記憶，再分流寫入兩個記憶儲存庫（`agentmemory` API、`mempalace` CLI），最後清理過期資料。

`領域流程 (Domain Flow):`

1. `distill` 主命令啟動管道 → 呼叫 `readGbrainLogic()` 與 `readClaudeMemLogic()` 讀取新增觀察值
2. 呼叫 `OllamaService.Extract()` 透過 LLM 提取候選記憶 → 分類為 `Memory` 與 `Fact`
3. 寫入 `agentmemory` API（所有記憶）與 `mempalace mine`（通過真實性門檻的事實）
4. 更新 `StateStore` 中的遊標與蒸餾狀態 → 由 `retain` 清理過期資料與檔案

`核心實體 (Key Entities):` `Observation`, `Candidate`, `Memory`, `Fact`, `Cursor`, `Seen`, `Distilled`

`相關處理器 (Related Handlers):` `DistillCmd()`, `ExtractCmd()`, `WriteAgentMemoryCmd()`, `WriteMempalaceCmd()`, `RetainCmd()`, `ResetCmd()`

---

### 資料匯出 (Data Export)

提供從 `gbrain`、`claude-mem`、`mempalace` 三個儲存來源匯出原始資料的能力，支援增量匯出（基於遊標）與全量匯出。

`領域流程 (Domain Flow):`

1. 使用者執行 `cc-plugin export <子命令>` → 選擇匯出 `gbrain`、`claudemem` 或 `mempalace`
2. 讀取 `StateStore` 遊標（增量模式）或從 epoch 0 開始（`--all` 模式）
3. `mempalace` 子命令支援類別清單（CSV）與完整 Markdown 結構匯出（`--data`）

`核心實體 (Key Entities):` `DrawerRow`, `Observation`

`相關處理器 (Related Handlers):` `ExportCmd()`, `GbrainCmd()`, `ClaudeMemCmd()`, `MempalaceCmd()`

---

### 環境初始化與配置同步 (Environment Initialization & Config Sync)

透過 `run.sh`（macOS/Unix）將本庫的設定檔與範本軟連結至使用者的家目錄資料夾（`$HOME/.claude`、`$HOME/.gemini`、`$HOME/.hermes` 等），同步外部工具設定（LiteLLM、CCStatusline、Tokscale）。

`領域流程 (Domain Flow):`

1. 執行 `run.sh` → 建立家目錄結構
2. 軟連結全域設定檔（`CLAUDE.global.md`、`settings.json`）→ 至 Claude Code 與 Gemini CLI
3. 複製或連結外部工具設定（LiteLLM、CCStatusline、Tokscale）
4. 建立本專案 `tmp/` 下的反向連結以供調試

`核心實體 (Key Entities):` `CLAUDE.global.md`, `settings.json`, `litellm_config.yaml`

`相關處理器 (Related Handlers):` `run.sh`

---

### AI 技能與代理生態 (AI Skills & Agents Ecosystem)

提供可跨 AI 編碼代理共用的自訂技能集與專屬代理定義，劃分為十一個模組化插件目錄，並透過 `npx skills` CLI 安裝至 55+ 個支援的 AI Agent：

- `base` (預設基礎插件 — must-install：Stop/StopFailure 終端機 bell hook)
- `apple` (macOS 整合)
- `explore` (抓取與摘要：MarkItDown, Scrapling, Firecrawl, Playwright)
- `last30days-skill` (趨勢與研究：Matt Van Horn 的 Last 30 Days 趨勢搜尋插件，支援 Reddit, X, YouTube 等多來源研究)
- `general` (通用技能：domain, project-explore, model-evaluator)
- `god` (系统大一統理論：LLM 力學、領域探索、融合方法)
- `review` (審查插件：consistency、business-improvement、folder-structure、naming-convention、doc-sync、dependency-hygiene、learning-document)
- `media` (影片生成與劇本創作：prompt-to-story-script, scene-to-video-prompt, character-setting)
- `team` (代理團隊規劃與設計：team-design, role-generator, orchestration-config)
- `tmp` (臨時測試)
- `superpowers` (核心流程：TDD、除錯、協作模式 — vendored from obra/superpowers)
- `gosdk` (Go 開發：code-quality、mvc、performance、network、dead-code、naming、zap→slog 遷移)
- `understand-anything` (AI 輔助代碼庫理解：代碼分析、架構分析與說明等 8 個技能與 9 個分析代理)

`領域流程 (Domain Flow):`

1. 開發者在對應的 `plugins/<name>/skills/` 目錄下建立 `SKILL.md`（符合 agentskills.io 規範）
2. 使用 `npx skills add .` 掃描並註冊技能至 `skills.json`，並安裝至多個 AI Agent（Antigravity、Claude Code、Gemini CLI 等）
3. 各插件的 manifest（`plugins/<name>/.claude-plugin/plugin.json`）定義專屬的 hooks、monitors、MCP/LSP 整合

`核心實體 (Key Entities):` `SKILL.md`, `plugin.json`, `hooks.json`, `monitors.json`, `skills.json`

`相關處理器 (Related Handlers):` `feature` agent, `post-tool.sh` hook

---

## 領域關聯 (Domain Relationships)

- `記憶蒸餾管道` 的輸出（`Memory`、`Fact`）寫入外部記憶儲存庫，而 `資料匯出` 則可從同一儲存庫反向匯出資料
- `環境初始化` 負責將 `AI 技能與代理` 的設定檔同步至各個 AI Agent 的家目錄
- `資料匯出` 與 `記憶蒸餾管道` 共用 `StateStore`（遊標機制）以支援增量操作

## 使用方式 (Usage)

### 記憶蒸餾

```bash
# 執行完整蒸餾管道（讀取 → 提取 → 寫入 → 清理）
cc-plugin distill

# 僅提取記憶（從 stdin 讀取 JSON 觀察值）
cc-plugin extract < observations.json

# 清理狀態（重置遊標、已見、已蒸餾紀錄）
cc-plugin reset
```

### 資料匯出

```bash
# 匯出 mempalace 類別清單
cc-plugin export mempalace

# 匯出 mempalace 完整 Markdown 結構
cc-plugin export mempalace --data -o ./export

# 匯出 gbrain 觀察值（增量）
cc-plugin export gbrain

# 匯出 claude-mem 觀察值（全量）
cc-plugin export claudemem --all
```

### 環境初始化

```bash
# 初始化軟連結與設定同步
chmod +x run.sh && ./run.sh

# 安裝技能至 AI Agents
npx skills add .
```

## 改善建議 (Improvement Suggestions)

Based on codebase analysis:

- [ ] `readClaudeMemLogic()` 在 `cmd/read_logic.go` 中重複建立 `StateStore`，應接收外部傳入的 store 以避免連線浪費
- [ ] `cmd/export/gbrain.go` 與 `cmd/read_logic.go` 中 `readGbrainLogic` / `gbrainRead` 功能幾乎重複，應整合為共用函數
- [ ] `plugins/general/skills/anti-sabotage-skill.md` 是一個散落的技能草稿，未轉換為正式的 `SKILL.md` 目錄結構
- [ ] `config/default_settings.json` 為空 (`{}`)，預設設定全部寫死在 `config.go` 中，建議遷移至 JSON 以利外部修改
- [ ] `cmd/write_agentmemory.go` 中 `resp.Body` 的 `defer resp.Close()` 在迴圈內使用可能造成資源洩漏
