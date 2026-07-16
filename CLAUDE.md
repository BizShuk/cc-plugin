# CC-Plugin — 技術脈絡 (Technical Context)

## 專案結構 (Project Structure)

```tree
.
├── .claude-plugin/           # 根 marketplace 註冊表
├── .vscode/                  # repo-relative launch 與 editor 設定
├── cmd/                      # Cobra CLI 子命令
│   ├── root.go               # CLI root，註冊所有子命令
│   ├── export/               # gbrain / claudemem / mempalace 匯出
│   ├── memory/               # 記憶相關子命令與邏輯
│   ├── topology/             # Markdown 知識圖譜驗證、查詢與重建
│   └── sample/api/           # Anthropic SDK + MiniMax 相容端點範例
├── config/                   # Viper 預設與各 AI Agent 設定
├── docs/
│   ├── backlog/              # 尚未進入實作的產品決策
│   ├── memory/               # 歷史操作與決策
│   ├── specs/                # 已完成並接受的設計
│   └── tutorials/            # 教學文件
├── model/                    # 記憶蒸餾領域模型與 SQLite StateStore
├── pkg/
│   ├── topology/             # 純 Markdown topology 解析與圖運算
│   ├── ccstatusline/         # CCStatusline 設定
│   ├── hermes/               # Hermes Agent 設定
│   ├── litellm/              # LiteLLM proxy 設定範本
│   ├── lsp/                  # Marksman 安裝說明
│   ├── system-prompts/       # 外部 system prompt 資源
│   ├── tools/                # 外部工具與 submodule
│   └── usage/                # 使用統計工具設定
├── plugins/                  # 8 個本地 Claude Code plugins
│   ├── experiment/           # 候選技能沙盒
│   ├── explore/              # 摘要、轉檔、專案探索與路由；含 plugin-scoped MCP
│   ├── general/              # 通用技能、feature agent、hooks、output styles、Marksman LSP
│   ├── god/                  # 系統大一統理論與通用算子
│   ├── review/               # 審查、規劃、自演化與 review-coordinator
│   ├── team/                 # Agent team 設計與角色資料
│   ├── tools/                # Apple Calendar/Mail/Notes/Reminders
│   └── ultra-explore/        # 可續跑的知識庫建構管道
├── scripts/pluginmeta/       # plugin/skill/agent metadata 稽核工具
├── ecosystem.config.js      # PM2 常駐程序與 cron
├── main.go                   # Go CLI 入口點
├── go.mod                    # Go 模組定義
└── run.sh                    # 可重入環境初始化腳本（macOS/Unix）
```

## 技術棧 (Tech Stack)

- Language: `Go 1.26.3`, `Python 3.10+` (用於輔助技能與範例)
- CLI Framework: `spf13/cobra`
- Configuration: `spf13/viper` + `viper.SetDefault` 直接宣告預設值
- ORM: `gorm` + `SQLite`（state store, claude-mem 讀取, mempalace 讀取）
- LLM: `Ollama` HTTP API（預設模型 `qwen3:14b-q4_K_M`）
- LSP: `gopls` (Go), `marksman` (Markdown)
- Logging: Go `log/slog`，由 `github.com/bizshuk/gosdk/log` 安裝全域 handler
- Custom SDK: `github.com/bizshuk/gosdk`（config、log）
- Key dependencies: `go-homedir`, `gocsv`, `yaml.v3`

## 關鍵決策 (Key Decisions)

- `Cobra + Viper` 組合：CLI 指令定義與設定管理標準模式，支援環境變數覆蓋
- `Viper 預設值單一來源`：預設值集中於 `config/config.go`，不使用嵌入式 JSON 預設檔
- `GORM + SQLite` 作為狀態儲存：輕量、無需外部資料庫服務、適合單機排程任務
- `Ollama 本地 LLM`：隱私優先，不將記憶資料傳至雲端 API
- `指紋 (Fingerprint) 去重`：透過 SHA-256 雜湊（正規化文本 + 排序實體）避免重複記憶
- `真實性門檻 (Truth Qualification)`：僅經人類確認、第一人稱事實/經驗、或 2+ 來源佐證的候選才寫入 mempalace 作為 Fact
- `agentskills.io 規範`：技能採用 YAML frontmatter + Markdown 格式，支援跨 Agent 安裝
- `軟連結同步`：以 symlink 而非複製來管理跨目錄設定，確保單一來源
- `模組化插件架構 (Modular Plugin Architecture)`：8 個本地 plugin 依職責拆分，skill/agent 由標準目錄自動探索，manifest 不重複列舉檔案。
- `整合 Marksman LSP`：引進 `marksman` Language Server 以提供 Markdown 的補全、診斷與檔案鏈結管理。

## 模組對應 (Module Mapping)

| 業務領域 (Domain) | 套件/模組 (Package/Module)                                                  | 進入點 (Entry Point)                                    |
| ----------------- | --------------------------------------------------------------------------- | ------------------------------------------------------- |
| 記憶蒸餾管道      | `cmd/memory/`, `model/`                                                     | `DistillCmd()`                                          |
| LLM 提取          | `cmd/memory/ollama.go`                                                      | `ExtractCmd()`, `OllamaService.Extract()`               |
| 讀取來源          | `cmd/memory/read_logic.go`                                                  | `readGbrainLogic()`, `readClaudeMemLogic()`             |
| 寫入儲存          | `cmd/memory/write_*.go`                                                     | `WriteAgentMemoryCmd()`, `WriteMempalaceCmd()`          |
| 資料匯出          | `cmd/export/`                                                               | `ExportCmd()`                                           |
| Topology 圖譜     | `cmd/topology/`, `pkg/topology/`                                            | `TopologyCmd()`, `LoadTopology()`                       |
| 狀態管理          | `model/store.go`, `model/cursor.go`                                         | `NewStateStore()`                                       |
| 狀態包裝          | `cmd/memory/state.go`                                                       | `NewStateStore()`                                       |
| 環境初始化        | `run.sh`, `config/`                                                         | `config.Init()`                                         |
| AI 技能           | `plugins/` (experiment, explore, general, god, review, team, tools, ultra-explore) | 各 `SKILL.md`                                      |
| 知識庫建構        | `plugins/ultra-explore/skills/`, `plugins/ultra-explore/agents/`            | `ultra-explore` 入口 + kb-\* 10 項, `kb-coordinator.md` |
| 審查、規劃與演化  | `plugins/review/skills/`                                                    | `auto-evolving` 單一演化閉環 + 各專項 `SKILL.md`        |
| AI 代理           | `plugins/general/agents/`, `plugins/review/agents/`                         | `feature.md`, `review-coordinator.md`                   |
| Plugin metadata   | `scripts/pluginmeta/`                                                       | `go run ./scripts/pluginmeta`                           |

## 開發指南 (Development Guide)

### 前置需求 (Prerequisites)

- Go 1.26.3+
- SQLite3
- Ollama（用於 LLM 提取，預設 `http://localhost:11434`）
- `mempalace` CLI（用於事實寫入）
- `jq`（用於 hook 腳本解析 JSON）
- `marksman`（選用；`plugins/general/.lsp.json` 的 Markdown LSP）
- `codegraph`（選用；`plugins/explore/.mcp.json` 的 MCP server）

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

主要測試套件：`cmd/export`, `cmd/memory`, `cmd/topology`, `pkg/topology`, `scripts/pluginmeta`

### 部署 (Deploy)

```bash
# 安裝至 $GOPATH/bin
go install

# 排程執行（每日 03:00）
crontab -e
# 加入: 0 3 * * * $HOME/go/bin/cc-plugin distill >> $HOME/.config/cc-plugin/logs/run.log 2>&1
```

## 慣例 (Conventions)

- Naming: Go 檔案以功能命名（`distill.go`, `read_logic.go`, `write_mempalace.go`），命令函數統一使用 `XxxCmd()` 格式
- Error handling: 使用 `fmt.Errorf("context: %w", err)` 包裝錯誤鏈，頂層由 Cobra 統一輸出至 stderr
- Logging: 診斷事件使用 `log/slog`；CLI 資料輸出保留 Cobra/stdout/file writer，避免污染 JSON/CSV/Markdown
- Testing: 測試檔案與實作同目錄，使用 `_test.go` 後綴
- Configuration: 設定路徑統一使用 `~` 前綴，由 `go-homedir` 展開；預設值寫在 `config.go`
- Skills: 遵循 `agentskills.io` 規範，YAML frontmatter 必須包含 `name` 與 `description`
- Plugin Manifest: `skills`／`agents` 維持空陣列，由 `plugins/<name>/skills/` 與 `plugins/<name>/agents/` 自動探索；不得重複列舉檔案
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
- `description` 採用 `>` 折疊式（`>` 或 `|`），長度 ≤ 1024 字元，且`必須`包含觸發詞（"Use when...", "Triggers on..."）
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
