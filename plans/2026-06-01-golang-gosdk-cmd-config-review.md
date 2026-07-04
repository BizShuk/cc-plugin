# Plan: golang-gosdk Code Review — `cmd/` & `config/`

## Context

使用者要求使用 `@golang-gosdk` 技能對 `cmd/` 和 `config/` 目錄進行整體程式碼審查。這是一個例行性代碼品質檢查，非實作任務。

---

## Review 範圍

### `cmd/` 目錄架構

| 元件 | 說明 |
| --- | --- |
| 10 個子命令 | `distill`(default), `retain`, `read-gbrain`, `read-claudemem`, `write-agentmemory`, `write-mempalace`, `extract`, `reset`, `export`, `export mempalace` |
| 框架 | spf13/cobra + spf13/viper |
| 狀態管理 | SQLite via GORM (`state.go`) |
| LLM 整合 | Ollama API (`ollama.go`) |

### `config/` 目錄架構

| 元件 | 說明 |
| --- | --- |
| 設定檔 | `settings.json`, `llmbox.json`, `minimax.json`, `keybindings.json`, `CLAUDE.global.md`, `default_settings.json` |
| 載入機制 | `viper` + `//go:embed` 編譯期嵌入 |
| 動態連結 | `run.sh` 執行時建立 symlink 至 `~/.claude/` |

---

## 審查重點

### 1. `cmd/` 審查項目

- [ ] **錯誤處理** — `distill.go` 中 pipeline 各階段的 error propagation 是否完整
- [ ] **Context 傳遞** — `state.go` / `ollama.go` 是否正確使用 `context.Context`
- [ ] **依賴注入** — `OllamaService` 是否可被置換（測試友好）
- [ ] **SQLite GORM 使用** — transaction 正確性、cursor 競態條件
- [ ] **Cobra 命令結構** — `root.go` AddCommand 順序與穩定性
- [ ] **資源清理** — `reset.go` 是否有可能造成資料遺失的風險

### 2. `config/` 審查項目

- [ ] **Viper 預設值** — `config.go` 中 `SetDefault` 是否覆蓋潛在的環境變數
- [ ] **敏感資訊** — `settings.json` / `llmbox.json` / `minimax.json` 是否包含明文 API key
- [ ] **路徑擴展** — `~` 路徑在 Windows/macOS 的相容性處理
- [ ] **go:embed 安全** — `default_settings.json` 是否為空或佔位符（確認中）

---

## 驗證方式

1. 逐一閱讀 `cmd/` 下各 `.go` 檔案
2. 使用 `golang-code-quality` 技能標準審查
3. 確認測試覆蓋率（如有 `*_test.go`）

---

## 輸出格式

審查完成後提供：
- 發現問題清單（ severity + 說明 ）
- 建議改進方向
- 無問題的模組確認

---

## 詳細審查計劃（來自 Plan agent）

### Critical 問題（生產環境前必須修復）

1. **state.go Close() bug** — `db.DB()` error 被忽略後可能 panic
2. **Context 傳遞缺失** — HTTP 和 DB 操作無法被取消
3. **exec.Command injection** — `write_mempalace.go` 中的命令注入風險

### High 優先級

4. 定義 error sentinel 以區分可復原與致命錯誤
5. 啟動時的設定驗證
6. distill pipeline 狀態操作的 transaction 包裝

### Medium 優先級

8. 使用 `homedir` package 替換 `expandPath`（platform 相容性）
9. 添加 pipeline 的整合測試
10. 移除測試中的全域 viper 狀態
11. OllamaService 中 HTTP client 重用

### Low 優先級

12. 清理空的 `model/models.go`
13. CSV escaping in export
14. 結構化日誌替代 `fmt.Fprintf`
15. 外部 API 呼叫的重試邏輯

### 關鍵檔案清單

- `../cmd/state.go`
- `../cmd/ollama.go`
- `../cmd/distill.go`
- `../config/config.go`
- `../cmd/root.go`
