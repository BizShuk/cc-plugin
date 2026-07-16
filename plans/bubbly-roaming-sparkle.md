# `cc-plugin tui` — Claude Code 使用統計互動式 TUI

## Context

`review all claude history ... evaluate how to be more effective` 的後續。手動分析已揭示：`2926` 個 session / `1.6G`、`Bash` 佔 `39%` 工具呼叫、`3%` session 破 `1MB`。但這些洞察目前只能靠一次性 shell + python 取得，無法重複探索。

本計畫新增 `cc-plugin tui` 指令：bubbletea 互動式介面，左側類別選單、右側明細面板，串流掃描 `~/.claude/projects/*.jsonl` 與 `~/.claude/history.jsonl`，讓使用統計可即時分層瀏覽（Overview → Sessions/Prompts/Tokens/Projects → 單一項目明細）。cc-plugin 目前無任何互動式程式碼，本指令同時建立 TUI 與 JSONL 串流兩項新慣例。

## 決策（已與使用者確認）

- TUI 函式庫：`bubbletea + lipgloss + bubbles`（charmbracelet）
- 版面：master-detail — 左側類別/項目清單，右側聚焦項目明細；`Enter` 鑽入、`Esc` 返回

## 檔案佈局

全部新增於 `cmd/`（沿用 repo 扁平風格，避免子套件需 re-export cobra）：

| 檔案 | 職責 |
| --- | --- |
| `cmd/tui.go` | Cobra 接線：`TuiCmd()`、flags、`tuiLogic()` |
| `cmd/tui_data.go` | 資料結構（`HistoryEntry`/`SessionRecord`/`Usage`/`*Stat`） |
| `cmd/tui_scan.go` | 串流掃描與聚合：`Scan()`、`loadSessions()`、`loadHistory()`、`tuiDump()` |
| `cmd/tui_model.go` | bubbletea `Model`、`Init`、`Update`、狀態機 |
| `cmd/tui_view.go` | `View()` + lipgloss 渲染 + 手刻長條圖 |
| `cmd/tui_test.go` | 掃描/聚合函式測試（不測 TUI render） |

## 資料結構（`tui_data.go`）

關鍵：`Message` 一律 `json.RawMessage`，只對 `assistant` 解碼 `model`+`usage`，避免 `user` 的多型 content 失敗。

```go
type HistoryEntry struct {
    Display        string          `json:"display"`        // 開頭 `/` = slash-command
    PastedContents json.RawMessage `json:"pastedContents"`
    Timestamp      int64           `json:"timestamp"`      // epoch ms
    Project        string          `json:"project"`
    SessionID      string          `json:"sessionId"`
}

type SessionRecord struct {
    Type      string          `json:"type"`            // user|assistant|ai-title|last-prompt|...
    Timestamp string          `json:"timestamp"`       // RFC3339Nano, "" 於部分 type
    SessionID string          `json:"sessionId"`
    LeafUUID  string          `json:"leafUuid"`        // last-prompt 用
    AITitle   string          `json:"aiTitle"`         // ai-title 用
    Message   json.RawMessage `json:"message"`         // 只對 assistant 解碼
}

type Usage struct {
    InputTokens              int64 `json:"input_tokens"`
    OutputTokens             int64 `json:"output_tokens"`
    CacheCreationInputTokens int64 `json:"cache_creation_input_tokens"`
    CacheReadInputTokens     int64 `json:"cache_read_input_tokens"`
}

type assistantMessage struct {
    Model string `json:"model"`
    Usage Usage  `json:"usage"`
}

type SessionStat struct {
    SessionID, FilePath, Project, EncodedDir, AITitle string
    MTime, LastActive time.Time
    PromptCount int                 // type=="user" 計數
    TypedCount  int                 // 由 history.jsonl join 得（乾淨來源）
    TotalTokens, InputTokens, OutputTokens, CacheCreate, CacheRead int64
    ModelUsage  map[string]int64    // model -> tokens
    SizeBytes   int64
}

type ProjectStat struct {
    EncodedDir, DecodedDir string
    SessionCount int
    TotalTokens, SizeBytes int64
    LastActive   time.Time
    SessionIDs   []string
}

type GlobalStat struct {
    SessionCount, ProjectCount, PromptCount int
    TotalTokens, InputTokens, OutputTokens, CacheCreate, CacheRead, SizeBytes int64
    ModelUsage   map[string]int64
    ProjectUsage map[string]int64
    ScanDuration time.Duration
}

type ScanResult struct {
    Global   GlobalStat
    Sessions []SessionStat  // 依 LastActive desc
    Projects []ProjectStat  // 依 TotalTokens desc
    History  []HistoryEntry
    Warnings []string
}
```

## 掃描器（`tui_scan.go`）

- **串流**：`bufio.NewScanner` + `scanner.Buffer(make([]byte,1<<20), 10<<20)`（10MB 上限），逐行 `json.Unmarshal`，**只累積 `SessionStat`、不保留逐筆 record** — 防 145MB 大檔 OOM。
- **`ScanOptions`**：`ProjectsDir`、`HistoryPath`、`MaxFileBytes`、`Progress chan<- ScanProgress`。
- **type 分派**：`assistant`→解碼 `Message` 聚合 token；`user`→`PromptCount++`；`ai-title`→首筆獲勝；`last-prompt`→用首輪建立的 `leafTs map[string]time.Time` 解析精確 `LastActive`；其餘 skip。
- **TypedCount**：掃完 history 後以 `map[sessionID]count` O(n+m) join 回填（取代從 session file 偵測 `isCompactSummary`，較穩）。
- **時間防呆**：`parseTimestamp(s) (time.Time, bool)`，失敗 fallback mtime 並記 `Warnings`。
- **背景載入 + spinner**：`scanCmd(opts) tea.Cmd` 回傳 `scanDoneMsg`/`scanErrMsg`；載入時顯示 `bubbles/spinner` + 進度行。

## TUI 架構（`tui_model.go` + `tui_view.go`）

狀態機：`stateLoading → stateOverview → stateItemList → stateItemDetail`，以 `stack []navFrame` 支援 `Esc` 返回。類別：`Overview / Sessions / Prompts / Tokens / Projects`。

- `Model`：`width/height`、`result *ScanResult`、`state`、`category`、`stack`、`spinner`、五個 `list.Model`、`viewport.Model`、`chartRows []barRow`。
- `Init()`：`tea.Batch(spinner.Tick, scanCmd(opts))`。
- `Update()`：處理 `WindowSizeMsg`/`KeyMsg`/`scanDoneMsg`/list 內部訊息。
- `View()`：`lipgloss.JoinHorizontal(Top, left, right)`，左 35% 寬。

按鍵：`q/ctrl+c` 離開、`Esc` 返回（頂層則離開）、`Enter/→` 鑽入、`←` 返回、`↑↓` 列表導覽、`r` 重掃、`?` 說明。

長條圖手刻（Tokens 類別，`t` 切換 project/model 視圖）：
```go
func renderBar(label string, value, max int64, width int) string {
    pct := float64(value) / float64(max) // max==0 退化處理
    filled := int(pct * float64(width))
    bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
    return fmt.Sprintf("%-20s %s %s", label, bar, humanTokens(value))
}
```
`humanTokens(1234567)`→`"1.2M"`，lipgloss 上色。

## Cobra 接線（`tui.go`）+ 註冊

```go
func TuiCmd() *cobra.Command {
    var projectsDir, historyPath string
    var dump bool
    var maxFileMB int
    cmd := &cobra.Command{
        Use:   "tui",
        Short: "Interactive terminal UI to browse Claude Code usage statistics",
        RunE: func(cmd *cobra.Command, args []string) error {
            opts := ScanOptions{
                ProjectsDir:     expandPath(projectsDir),   // cmd/root.go:13
                HistoryPath:     expandPath(historyPath),
                MaxFileBytes:    int64(maxFileMB) << 20,
            }
            if dump { return tuiDump(opts) }
            return tuiLogic(opts)
        },
    }
    cmd.Flags().StringVar(&projectsDir, "projects-dir",
        viper.GetString("sources.claude.projects_dir"), "...")
    cmd.Flags().StringVar(&historyPath, "history-path",
        viper.GetString("sources.claude.history_path"), "...")
    cmd.Flags().IntVar(&maxFileMB, "max-file-mb", 0, "...")
    cmd.Flags().BoolVar(&dump, "dump", false, "Print stats as JSON and exit (no TUI)")
    return cmd
}
```
`tuiLogic`：`tea.NewProgram(newModel(opts), tea.WithAltScreen()).Run()`。

註冊於 `cmd/root.go init()`（`RootCmd.AddCommand(TuiCmd())`，接在 `export.ExportCmd()` 後）。

## viper 預設（`config/config.go:27` 後新增兩行）

```go
viper.SetDefault("sources.claude.projects_dir", "~/.claude/projects")
viper.SetDefault("sources.claude.history_path", "~/.claude/history.jsonl")
```

## 依賴

```bash
go get github.com/charmbracelet/bubbletea@latest
go get github.com/charmbracelet/bubbles@latest
go get github.com/charmbracelet/lipgloss@latest
go mod tidy
```

## 測試（`tui_test.go`，沿用 `state_test.go` 的 `os.MkdirTemp` + `defer RemoveAll` 風格）

只測資料層，不測 `Update`/`View`：
- `TestParseSessionFile_Assistant` — token 四欄聚合正確
- `TestParseSessionFile_AITitle` — 重複 ai-title 首筆獲勝
- `TestParseSessionFile_LastPrompt` — leafUuid 解析 `LastActive`
- `TestParseSessionFile_TokenExtrasIgnored` — `server_tool_use`/`cache_creation.*` 容忍不解碼失敗
- `TestParseSessionFile_PolymorphicUserContent` — string 與 array content 皆 `PromptCount++`
- `TestLoadHistory_Basic` — slash vs prompt 可由 `Display[0]` 區分
- `TestScan_JoinsTypedCount` — history join 回填 TypedCount
- `TestScanResult_GlobalAggregation` — 全域聚合 = 各 session 之和、ProjectCount 正確
- `TestHumanTokens` — 0/999/1500/1.5M/12.3B 後綴

## 驗證

1. `go test ./cmd/... -count=1` — 全綠
2. `go vet ./...` — 乾淨
3. `go run . tui --dump` — 印 JSON 統計（無 TUI，驗證掃描層對真實 `~/.claude/` 端對端可用）
4. `go run . tui` — Overview 在 ~3–5s 內顯示（spinner 期間不空白）
   - Sessions 列表依 `LastActive` desc；`Enter` 進明細見 token/model 分解
   - Tokens 見長條圖，`t` 切 project/model
   - Projects 鑽入僅見該專案 sessions
   - Prompts 中 slash-command 視覺區分
   - `Esc` 返回、`q` 離開、`r` 重掃、`?` 說明

## 實作順序

1. **A 資料層**：`tui_data.go` + `tui_scan.go` + `tui_test.go` + config 預設 → 測試全綠
2. **B dump 模式**：`tui.go` 的 `--dump` flag → 對真實資料驗證掃描
3. **C TUI 外殼**：`tui_model.go` 最小版（spinner + Overview）→ 確認整合
4. **D 導覽**：類別選單 → 列表 → 明細，逐類別接（Sessions → Prompts → Projects → Tokens）
5. **E 拋光**：說明覆蓋層、按鍵提示、錯誤狀態、重掃

## 風險

| 風險 | 緩解 |
| --- | --- |
| 145MB 檔爆 10MB scanner buffer | 單行即一筆 JSON record，極少破 10MB；超過則 scanner 報錯記 `Warnings` 並 skip 該 record |
| `user` 多型 content 解碼失敗 | `Message` 保持 `json.RawMessage`，只對 `assistant` 解碼 |
| 首掃 5s 空白 | spinner + `ScanProgress` 進度行 |
| `time.Parse(RFC3339Nano)` 失敗 | `parseTimestamp` 回 `bool`，fallback mtime + Warnings |
