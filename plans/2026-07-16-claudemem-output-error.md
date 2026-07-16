# Claude-mem 匯出錯誤處理實作計畫

> `agentic worker` 需在同一 session 依 TDD 順序執行下列步驟。

`Goal:` 輸出 JSON 失敗時回傳錯誤，並阻止 `claude-mem` 遊標推進。

`Architecture:` 命令透過 Cobra 提供的 `cmd.OutOrStdout()` 取得可注入 writer，再以 `json.Encoder` 寫入。輸出成功後才沿用現有遊標更新邏輯。

`Tech Stack:` Go 1.25、Cobra、GORM、SQLite、testing

## 全域限制

- 只修正 review 的第一項。
- 錯誤以 `fmt.Errorf("context: %w", err)` 保留 error chain。
- 不覆寫工作區其他未提交變更。

---

### Task 1: 防止輸出失敗後更新遊標

`Files:`

- Create: `cmd/export/claudemem_test.go`
- Modify: `cmd/export/claudemem.go:77-81`

`Interfaces:`

- Consumes: `(*cobra.Command).OutOrStdout() io.Writer`
- Produces: `ClaudeMemCmd()` 在 writer 失敗時回傳 wrapped error

- [x] Step 1: 新增 `TestClaudeMemCmdDoesNotAdvanceCursorWhenOutputFails`，使用固定回傳 `write failed` 的 writer，並斷言命令失敗與遊標為 `0`。
- [x] Step 2: 執行 `go test ./cmd/export -run TestClaudeMemCmdDoesNotAdvanceCursorWhenOutputFails -count=1`，確認因現行 `fmt.Println` 繞過 Cobra writer 而失敗。
- [x] Step 3: 將 `json.MarshalIndent` 與 `fmt.Println` 替換為設定縮排的 `json.NewEncoder(cmd.OutOrStdout())`，`Encode` 失敗時回傳 `fmt.Errorf("write claude-mem export: %w", err)`。
- [x] Step 4: 重跑聚焦測試，預期通過。
- [x] Step 5: 執行 `gofmt -w cmd/export/claudemem.go cmd/export/claudemem_test.go`、`go test ./... -count=1` 與 `go vet ./...`。
