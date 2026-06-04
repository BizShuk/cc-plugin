# CC-Plugin — 技術脈絡 (Technical Context)

## 專案結構 (Project Structure)

```tree
.
├── .claude-plugin/
│   └── plugin.json           # 插件 manifest（定義 skills, agents, hooks, monitors, MCP, LSP）
├── agents/                   # 自訂 AI 代理 (Custom Agents)
│   ├── feature.md            # 通用功能實作代理
│   └── README.md
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
│   └── export/               # export 子命令群
│       ├── export.go         # 頂層 export 命令
│       ├── claudemem.go      # export claudemem
│       ├── gbrain.go         # export gbrain
│       └── mempalace.go      # export mempalace（CSV / Markdown）
├── config/                   # 設定管理
│   ├── config.go             # Viper 初始化 + 嵌入預設值
│   ├── default_settings.json # 空白預設（設定寫死在 config.go）
│   ├── CLAUDE.global.md      # 全域 AI Agent 指令（軟連結至 $HOME）
│   ├── settings.json         # Claude Code 使用者設定
│   ├── keybindings.json      # 快捷鍵設定
│   ├── llmbox.json           # LLMBox 設定
│   └── minimax.json          # MiniMax 設定
├── hooks/                    # Hooks 配置與腳本
│   ├── hooks.json            # PostToolUse hook 定義
│   └── post-tool.sh          # Go 檔案自動 fmt + lint
├── model/                    # 核心領域模型
│   ├── agentmemory.go        # Memory 結構
│   ├── claudemem.go          # ClaudeMemObservation 結構
│   ├── cursor.go             # Cursor, Seen, Distilled ORM 模型
│   ├── distiller.go          # Candidate 結構
│   ├── gbrain.go             # Observation 結構
│   ├── mempalace.go          # Fact 結構
│   ├── store.go              # StateStore — SQLite 狀態管理
│   └── model/                # 嵌入式模型資源
│       └── bge-m3-gguf.md    # BGE-M3 模型文件
├── monitors/                 # Monitors 配置
│   └── monitors.json         # error-log / access-log 監控
├── pkg/                      # 外部工具設定範本與資源
│   ├── ccstatusline/         # CCStatusline 設定
│   ├── claude-code-system-prompts/ # Claude Code 系統提示詞
│   ├── hermes/               # Hermes Agent 設定（MEMORY.md, USER.md）
│   ├── litellm/              # LiteLLM proxy 設定範本
│   ├── summarize.sh/         # summarize CLI 設定
│   └── usage/tokscale/       # Tokscale 設定
├── plans/                    # 實作計畫與記憶系統設計文件
│   ├── agent/                # Agent 設計文件
│   └── memory/               # 記憶系統架構文件
├── skills/                   # 自訂技能集
│   ├── apple-calendar/       # Apple Calendar 管理
│   ├── apple-email/          # Apple Email 管理
│   ├── apple-notes/          # Apple Notes 管理
│   ├── apple-reminders/      # Apple Reminders 管理
│   ├── content-summarizer/   # 內容摘要（URL/文件/影片 → 重點 + 商業價值）
│   ├── domain/               # 業務領域技能
│   ├── firecrawl/            # Firecrawl 網頁擷取
│   ├── model-evaluator/      # 模型評估
│   ├── playwright-cli/       # Playwright 瀏覽器自動化
│   ├── project-explore/      # 專案探索與文件化
│   └── superpower/           # 技能路由（自動發現並載入適用技能）
├── .lsp.json                 # LSP 伺服器設定（gopls）
├── .mcp.json                 # MCP 伺服器設定（Playwright, codegraph）
├── main.go                   # Go CLI 入口點
├── go.mod                    # Go 模組定義
├── run.sh                    # 環境初始化腳本（macOS/Unix）
└── crontab.txt               # 排程設定（每日 03:00 執行 distill）
```

## 技術棧 (Tech Stack)

- Language: `Go 1.25`
- CLI Framework: `spf13/cobra`
- Configuration: `spf13/viper` + 嵌入式 JSON 預設（`go:embed`）
- ORM: `gorm` + `SQLite`（state store, claude-mem 讀取, mempalace 讀取）
- LLM: `Ollama` HTTP API（預設模型 `qwen3:14b-q4_K_M`）
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

## 模組對應 (Module Mapping)

| 業務領域 (Domain) | 套件/模組 (Package/Module) | 進入點 (Entry Point) |
| ----------------- | -------------------------- | -------------------- |
| 記憶蒸餾管道 | `cmd/`, `model/` | `DistillCmd()` |
| LLM 提取 | `cmd/ollama.go` | `ExtractCmd()`, `OllamaService.Extract()` |
| 讀取來源 | `cmd/read_logic.go` | `readGbrainLogic()`, `readClaudeMemLogic()` |
| 寫入儲存 | `cmd/write_*.go` | `WriteAgentMemoryCmd()`, `WriteMempalaceCmd()` |
| 資料匯出 | `cmd/export/` | `ExportCmd()` |
| 狀態管理 | `model/store.go`, `model/cursor.go` | `NewStateStore()` |
| 環境初始化 | `run.sh`, `config/` | `config.Init()` |
| AI 技能 | `skills/` | 各 `SKILL.md` |
| AI 代理 | `agents/` | `feature.md` |

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
