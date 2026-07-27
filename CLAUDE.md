# CC-Plugin — 技術脈絡 (Technical Context)

## 關鍵決策 (Key Decisions)

- `Cobra + Viper` 組合：CLI 指令定義與設定管理標準模式，支援環境變數覆蓋
- `Viper 預設值單一來源`：預設值集中於 `config/config.go`
- `Home path resolution`：所有 `~` 路徑統一使用 `gosdk/config.ExpandHome`；不保留
  自訂 home expansion wrapper 或 `go-homedir` dependency
- `GORM + SQLite` 作為狀態儲存：輕量、無需外部資料庫服務、適合單機排程任務
- `Ollama 本地 LLM`：隱私優先，不將記憶資料傳至雲端 API；預設模型
  `qwen3:14b-q4_K_M`，預設端點 `http://localhost:11434`
- `指紋 (Fingerprint) 去重`：透過 SHA-256 雜湊（正規化文本 + 排序實體）避免重複記憶
- `真實性門檻 (Truth Qualification)`：僅經人類確認、第一人稱事實/經驗、或 2+ 來源佐證的候選才寫入 mempalace 作為 Fact
- `agentskills.io 規範`：技能採用 YAML frontmatter + Markdown 格式，支援跨 Agent 安裝
- `軟連結同步`：以 symlink 而非複製來管理跨目錄設定，確保單一來源
- `模組化插件架構 (Modular Plugin Architecture)`：本地 plugin 依職責拆分，skill/agent 由標準目錄自動探索，manifest 不重複列舉檔案。
- `LSP 整合`：`gopls` (Go) 與 `marksman` (Markdown) 提供補全、診斷與檔案鏈結管理。
- `Claude-mem 匯出 ID 遊標`：`export claudemem` 使用獨立 `claude-mem-export` 狀態與 `observations.id` autoincrement 順序，避免 timestamp 相同或回填造成漏匯，並以 SQLite `mode=ro` 讀取來源。

## 模組對應 (Module Mapping)

| 業務領域 (Domain) | 套件/模組 (Package/Module) | 進入點 (Entry Point) |
| ----------------- | -------------------------- | -------------------- |
| 記憶蒸餾管道 | `cmd/memory/`, `model/` | `memory.DistillCmd` |
| LLM 提取 | `cmd/memory/ollama.go` | `memory.ExtractCmd`, `OllamaService.Extract()` |
| 讀取來源 | `cmd/memory/read_logic.go` | `readGbrainLogic()`, `readClaudeMemLogic()` |
| 寫入儲存 | `cmd/memory/write_*.go` | `memory.WriteAgentMemoryCmd`, `memory.WriteMempalaceCmd` |
| 資料匯出 | `cmd/export/` | `export.ExportCmd`, `export.ClaudeMemCmd` |
| Topology 圖譜 | `cmd/topology/`, `pkg/topology/` | `topology.TopologyCmd`, `LoadTopology()` |
| 狀態管理 | `model/store.go`, `model/cursor.go` | `NewStateStore()`, `GetCursorPosition()`, `SetCursorPosition()` |
| 環境初始化 | `run.sh`, `config/` | `config.Init()` |
| AI 技能 | `plugins/` | 各 `SKILL.md` |
| 知識庫建構 | `plugins/ultra-explore/skills/`, `plugins/ultra-explore/agents/` | `ultra-explore`, `kb-coordinator.md` |
| 審查、規劃與演化 | `plugins/review/skills/` | `auto-evolving` 與各專項 `SKILL.md` |
| AI 代理 | `plugins/general/agents/`, `plugins/review/agents/` | `feature.md`, `review-coordinator.md` |

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

### 部署 (Deploy)

```bash
# 安裝至 $GOPATH/bin
go install

# 排程執行（每日 03:00）
crontab -e
# 加入: 0 3 * * * $HOME/go/bin/cc-plugin distill >> $HOME/.config/cc-plugin/logs/run.log 2>&1
```

## 慣例 (Conventions)

- Naming: Go 檔案以功能命名（`distill.go`, `read_logic.go`, `write_mempalace.go`）；Cobra command 使用 package-level exported `XxxCmd` 變數，flags 在 `init()` 綁定
- Error handling: 使用 `fmt.Errorf("context: %w", err)` 包裝錯誤鏈，頂層由 Cobra 統一輸出至 stderr
- Logging: 診斷事件使用 `log/slog`；CLI 資料輸出保留 Cobra/stdout/file writer，避免污染 JSON/CSV/Markdown
- Testing: 測試檔案與實作同目錄，使用 `_test.go` 後綴
- Configuration: 設定路徑統一使用 `~` 前綴並展開；預設值寫在 `config/config.go`
- Skills: 遵循 `agentskills.io` 規範，YAML frontmatter 必須包含 `name` 與 `description`；
  完整 frontmatter tier 規範見 `.claude/skills/skill-frontmatter/SKILL.md`
- Plugin Manifest: `skills`／`agents` 維持空陣列，由 `plugins/<name>/skills/` 與 `plugins/<name>/agents/` 自動探索；不得重複列舉檔案
- 鬆散技能檔案禁止：`plugins/<plugin>/skills/` 頂層只放子目錄，所有 `SKILL.md` 必須位於獨立子目錄內
- 插件說明文件 (Plugin README)：位於 `plugins/` 目錄下的每個插件 (Plugin) 都必須在其資料夾內擁有一個 `README.md` 用以說明該插件的用途與使用方法；更新插件 (Plugin) 時亦必須同步更新對應的 `README.md`
