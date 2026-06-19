# CC-Plugin — 技術脈絡 (Technical Context)

## 專案結構 (Project Structure)

```tree
.
├── .lsp.json                 # LSP 伺服器設定（gopls, marksman）
├── .mcp.json                 # MCP 伺服器設定（Playwright, codegraph）
├── cmd/                      # Cobra CLI 子命令
│   ├── root.go               # CLI root，註冊所有子命令
│   ├── distill.go            # distill 主命令 — 完整蒸餾管道編排
│   ├── ollama.go             # extract 命令 + OllamaService LLM 呼叫
│   ├── read_logic.go         # gbrain / claude-mem 讀取邏輯
│   ├── write_agentmemory.go  # 寫入 agentmemory API
│   ├── write_mempalace.go    # 寫入 mempalace CLI
│   ├── retain.go             # 過期資料清理
│   ├── reset.go              # 狀態重置
│   ├── state.go              # StateStore 向後相容包裝
│   ├── export/               # export 子命令群
│   │   ├── export.go         # 頂層 export 命令
│   │   ├── claudemem.go      # export claudemem
│   │   ├── gbrain.go         # export gbrain
│   │   └── mempalace.go      # export mempalace（CSV / Markdown）
│   └── sample/               # API 使用範例 (Go with Anthropic SDK)
│       └── api/              # Anthropic SDK + MiniMax 相容端點範例
├── config/                   # 設定管理
│   ├── config.go             # Viper 初始化 + 嵌入預設值
│   ├── default_settings.json # 空白預設
│   ├── CLAUDE.global.md      # 全域 AI Agent 指令
│   ├── settings.json         # Claude Code 使用者設定
│   ├── keybindings.json      # 快捷鍵設定
│   ├── llmbox.json           # LLMBox 設定
│   └── minimax.json          # MiniMax 設定
├── model/                    # 核心領域模型
│   ├── agentmemory.go        # Memory 結構
│   ├── claudemem.go          # ClaudeMemObservation 結構
│   ├── cursor.go             # Cursor, Seen, Distilled ORM 模型
│   ├── distiller.go          # Candidate 結構
│   ├── gbrain.go             # Observation 結構
│   ├── mempalace.go          # Fact 結構
│   └── store.go              # StateStore — SQLite 狀態管理
├── pkg/                      # 外部工具設定範本與資源
│   ├── ccstatusline/         # CCStatusline 設定
│   ├── claude-code-system-prompts/ # Claude Code 系統提示詞
│   ├── hermes/               # Hermes Agent 設定
│   ├── litellm/              # LiteLLM proxy 設定範本
│   ├── lsp/                  # Marksman LSP 設定與 README
│   ├── tools/                # 外部工具依賴設定 (MarkItDown)
│   └── usage/tokscale/       # Tokscale 設定
├── plans/                    # 實作計畫與記憶系統設計文件
│   ├── agent/                # Agent 設計文件
│   └── memory/               # 記憶系統架構文件
├── plugins/                  # 模組化插件目錄 (Modular Plugins)
│   ├── base/                 # 預設基礎插件（must-install）— 提供 Stop/StopFailure 終端機 bell hook
│   │   └── hooks/            # 終端機鈴聲 (stop-bell.sh, hooks.json)
│   ├── apple/                # macOS Apple 整合插件
│   │   └── skills/           # Apple 相關技能 (apple-calendar, apple-email, apple-notes, apple-reminders)
│   ├── explore/              # 探索與抓取插件 (explore, scraping, fetching)
│   │   └── skills/           # 抓取與摘要技能 (content-summarizer, firecrawl, markitdown, playwright-cli, scrapling, summarize.sh)
│   ├── general/              # 通用功能插件
│   │   ├── agents/           # 自訂代理 (feature.md)
│   │   └── skills/           # 通用技能 (anti-sabotage, business-extract, domain, markdownlint, mermaid, model-evaluator, project-explore, topology-builder)
│   ├── media/                # 影片生成與劇本創作插件
│   │   ├── skills/           # 影片生成技能 (character-setting, prompt-to-story-script, scene-to-video-prompt)
│   │   └── voice/            # VoxCPM 聲音樣板與語音克隆指南設定
│   ├── team/                 # AI 代理團隊規劃與設計插件
│   │   └── skills/           # 團隊相關技能 (orchestration-config, role-generator, team-design)
│   ├── god/                  # 系統大一統理論插件 (Grand Unified Theory)
│   │   └── skills/           # 架構哲學技能 (llm-mechanics, domain-exploration, system-roles-laws, fusion-methods, unified-matrix, grand-unified-theory)
│   ├── review/               # 審查插件 (Review) — 一致性、業務改善、結構/命名/文件/依賴/測試/學習審查
│   │   ├── agents/           # 審查代理 (review-coordinator — 編排全部審查技能)
│   │   └── skills/           # 審查技能 (consistency, business-improvement, folder-structure, naming-convention, doc-sync, dependency-hygiene, learning-document)
│   ├── superpowers/          # 核心技能庫 (TDD、除錯、協作模式 — git submodule 來自 obra/superpowers)
│   │   └── skills/           # 流程技能 (brainstorming, dispatching-parallel-agents, using-superpowers, 等)
│   ├── gosdk/                # Go 開發工具包 (基於 github.com/bizshuk/gosdk)
│   │   ├── agents/           # Go 代理 (golang-refactor)
│   │   └── skills/           # Go 技能 (golang-code-quality, golang-mvc, golang-performance-tuning, 等)
│   └── tmp/                  # 臨時與測試用插件
│       ├── hooks/            # PostToolUse hooks (post-tool.sh, hooks.json)
│       ├── monitors/         # Monitors (monitors.json)
│   ├── tmp/                  # 臨時與測試用插件
│   │   ├── hooks/            # PostToolUse hooks (post-tool.sh, hooks.json)
│   │   ├── monitors/         # Monitors (monitors.json)
│   │   └── skills/           # 臨時與草稿技能 (celebrity-quotes, gemini-tmp, security-scanner, social-business-explore)
│   └── understand-anything/  # AI 輔助代碼庫理解插件 (Understand Anything)
│       └── understand-anything-plugin/ # 核心插件目錄
│           ├── agents/       # 分析代理 (project-scanner, file-analyzer, 等 9 個)
│           └── skills/       # 理解技能 (understand, understand-chat, 等 8 個)
├── main.go                   # Go CLI 入口點
├── go.mod                    # Go 模組定義
├── run.sh                    # 環境初始化腳本（macOS/Unix）
├── crontab.txt               # 排程設定（每日 03:00 執行 distill）
└── skills.json               # 技能註冊表
```

## 技術棧 (Tech Stack)

- Language: `Go 1.25`, `Python 3.10+` (用於輔助技能與範例)
- CLI Framework: `spf13/cobra`
- Configuration: `spf13/viper` + 嵌入式 JSON 預設（`go:embed`）
- ORM: `gorm` + `SQLite`（state store, claude-mem 讀取, mempalace 讀取）
- LLM: `Ollama` HTTP API（預設模型 `qwen3:14b-q4_K_M`）
- LSP: `gopls` (Go), `marksman` (Markdown)
- Custom SDK: `github.com/bizshuk/gosdk`（config 模組）
- Key dependencies: `go-homedir`, `gocsv`, `zap`

## 關鍵決策 (Key Decisions)

- `Cobra + Viper` 組合：CLI 指令定義與設定管理標準模式，支援環境變數覆蓋
- `GORM + SQLite` 作為狀態儲存：輕量、無需外部資料庫服務、適合單機排程任務
- `Ollama 本地 LLM`：隱私優先，不將記憶資料傳至雲端 API
- `指紋 (Fingerprint) 去重`：透過 SHA-256 雜湊（正規化文本 + 排序實體）避免重複記憶
- `真實性門檻 (Truth Qualification)`：僅經人類確認、第一人稱事實/經驗、或 2+ 來源佐證的候選才寫入 mempalace 作為 Fact
- `agentskills.io 規範`：技能採用 YAML frontmatter + Markdown 格式，支援跨 Agent 安裝
- `軟連結同步`：以 symlink 而非複製來管理跨目錄設定，確保單一來源
- `模組化插件架構 (Modular Plugin Architecture)`：將技能、代理、掛鉤與監控器拆分為 `apple`、`explore`、`general` 與 `tmp` 獨立插件目錄，便於分類管理與跨插件打包。
- `整合 Marksman LSP`：引進 `marksman` Language Server 以提供 Markdown 的補全、診斷與檔案鏈結管理。

## 模組對應 (Module Mapping)

| 業務領域 (Domain) | 套件/模組 (Package/Module)                                                                                   | 進入點 (Entry Point)                           |
| ----------------- | ------------------------------------------------------------------------------------------------------------ | ---------------------------------------------- |
| 記憶蒸餾管道      | `cmd/`, `model/`                                                                                             | `DistillCmd()`                                 |
| LLM 提取          | `cmd/ollama.go`                                                                                              | `ExtractCmd()`, `OllamaService.Extract()`      |
| 讀取來源          | `cmd/read_logic.go`                                                                                          | `readGbrainLogic()`, `readClaudeMemLogic()`    |
| 寫入儲存          | `cmd/write_*.go`                                                                                             | `WriteAgentMemoryCmd()`, `WriteMempalaceCmd()` |
| 資料匯出          | `cmd/export/`                                                                                                | `ExportCmd()`                                  |
| 狀態管理          | `model/store.go`, `model/cursor.go`                                                                          | `NewStateStore()`                              |
| 環境初始化        | `run.sh`, `config/`                                                                                          | `config.Init()`                                |
| AI 技能           | `plugins/` (base, apple, explore, general, god, review, tmp, superpowers, gosdk, media, team, understand-anything) | 各 `SKILL.md`                                  |
| 程式碼審查        | `plugins/review/skills/`                                                                                     | 各 `SKILL.md` (consistency 等 7 項)            |
| AI 代理           | `plugins/general/agents/`, `plugins/review/agents/`                                                          | `feature.md`, `review-coordinator.md` |

## 開發指南 (Development Guide)

### 前置需求 (Prerequisites)

- Go 1.25+
- SQLite3
- Ollama（用於 LLM 提取，預設 `http://localhost:11434`）
- `mempalace` CLI（用於事實寫入）
- `jq`（用於 hook 腳本解析 JSON）

### 安裝 (Installation)

```bash
# 複製專案
git clone https://github.com/bizshuk/cc-plugin.git

# 初始化環境（建立軟連結、同步設定）
chmod +x run.sh && ./run.sh

# 安裝為 Claude Code 插件
claude --plugin-dir .
```

### 建置 (Build)

```bash
go build -o cc-plugin main.go
```

### 測試 (Test)

```bash
go test ./... -count=1
```

測試檔案：`cmd/distill_test.go`, `cmd/main_test.go`, `cmd/ollama_test.go`, `cmd/state_test.go`, `cmd/export/mempalace_test.go`

### 部署 (Deploy)

```bash
# 安裝至 $GOPATH/bin
go install

# 排程執行（每日 03:00）
crontab -e
# 加入: 0 3 * * * $HOME/go/bin/cc-plugin >> $HOME/.distiller/logs/run.log 2>&1
```

## 慣例 (Conventions)

- Naming: Go 檔案以功能命名（`distill.go`, `read_logic.go`, `write_mempalace.go`），命令函數統一使用 `XxxCmd()` 格式
- Error handling: 使用 `fmt.Errorf("context: %w", err)` 包裝錯誤鏈，頂層由 Cobra 統一輸出至 stderr
- Logging: 使用 `fmt.Printf` / `fmt.Fprintf(os.Stderr, ...)` 進行簡易輸出（未引入結構化日誌）
- Testing: 測試檔案與實作同目錄，使用 `_test.go` 後綴
- Configuration: 設定路徑統一使用 `~` 前綴，由 `go-homedir` 展開；預設值寫在 `config.go`
- Skills: 遵循 `agentskills.io` 規範，YAML frontmatter 必須包含 `name` 與 `description`
- Plugin Manifest: 新增 skill 時需同步更新 `plugins/<name>/.claude-plugin/plugin.json` 的 `skills` 陣列（目錄型 `"./skills/skill-name"`，單檔型 `"./skills/skill-name.md"`），可選加 `keywords`
- 鬆散技能檔案禁止：`plugins/<plugin>/skills/` 頂層只放子目錄，所有 `SKILL.md` 必須位於獨立子目錄內
- 插件說明文件 (Plugin README)：位於 `plugins/` 目錄下的每個插件 (Plugin) 都必須在其資料夾內擁有一個 `README.md` 用以說明該插件的用途與使用方法；更新插件 (Plugin) 時亦必須同步更新對應的 `README.md`

### SKILL Frontmatter 規範 (Frontmatter Spec)

YAML frontmatter 分三個 tier，由簡至詳擇一使用：

| Tier       | 必填欄位              | 選填欄位                                                                                                  | 適用情境                       |
| ---------- | --------------------- | --------------------------------------------------------------------------------------------------------- | ------------------------------ |
| `minimal`  | `name`, `description` | —                                                                                                         | 參考文件、靜態知識             |
| `standard` | `name`, `description` | `version`, `allowed-tools`                                                                                | 一般 CLI 工具技能              |
| `full`     | `name`, `description` | `version`, `allowed-tools`, `user-invocable`, `disable-model-invocation`, `effort`, `context`, `metadata` | 需要控制模型呼叫行為的進階技能 |

額外規範：

- `name` 必須與所在子目錄名稱一致，使用 `kebab-case`
- `description` 採用 `>` 折疊式（`>` 或 `|`），長度 ≤ 1024 字元，且**必須**包含觸發詞（"Use when...", "Triggers on..."）
- `versio` 之類的拼字錯誤禁止（CI 將以 `yaml.Unmarshal` 驗證）
- Rules-style frontmatter（`trigger: always_on` + `globs` + `scope`）僅 `consistency` 與 `go-convention` 兩個常駐技能使用，其餘不得混用
- 標準 frontmatter 範例（`full` tier）：

```yaml
---
name: my-skill
description: >
    Use when ... Triggers on: "foo", "bar".
version: "1.0.0"
allowed-tools: Read, Bash, Glob
user-invocable: true
disable-model-invocation: false
effort: medium
context: fork
metadata:
    type: reference
    platforms: [macos, linux]
---
```
