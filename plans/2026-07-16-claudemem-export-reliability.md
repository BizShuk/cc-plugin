# Claude-mem 匯出可靠性實作計畫

> `Agentic worker` 必須在同一 session 依 TDD 順序執行，步驟使用 checkbox 追蹤。

`Goal:` 使 Claude-mem 增量匯出不漏掉相同或回填 timestamp 的新資料，並完整處理遊標、SQLite 與 JSON 錯誤邊界。

`Architecture:` `Cursor` 加入 `LastID`，匯出以獨立 `claude-mem-export` 狀態追蹤 autoincrement ID。讀取器以 SQLite `mode=ro` 開啟並關閉來源 DB，命令在 JSON 成功寫入後才更新遊標。

`Tech Stack:` Go 1.25、Cobra、GORM、SQLite、testing

## 全域限制

- 只修正 `cmd/export/claudemem.go` review 的剩餘項目與必要的 state model。
- 現有 `GetCursor` / `SetCursor` timestamp API 保持相容。
- 錯誤使用 `fmt.Errorf("context: %w", err)` 保留 error chain。
- 不覆寫其他未提交變更，不建立 commit。

---

### Task 1: ID 遊標狀態

`Files:`

- Modify: `model/cursor.go`
- Modify: `model/store.go`
- Create: `model/store_cursor_test.go`

`Interfaces:`

- Produces: `model.CursorPosition{LastTS int64, LastID int64}`
- Produces: `(*StateStore).GetCursorPosition(string) (CursorPosition, error)`
- Produces: `(*StateStore).SetCursorPosition(string, CursorPosition) error`

- [x] Step 1: 新增測試，將 `{LastTS: 100, LastID: 7}` 寫入後讀回並比對兩欄位。
- [x] Step 2: 執行 `go test ./model -run TestStateStoreCursorPosition -count=1`，預期因 API 尚未存在而編譯失敗。
- [x] Step 3: 在 `Cursor` 加入 `LastID int64` 與 `default:0`，實作查詢及 upsert 兩個欄位的 methods；現有 timestamp methods 維持原契約。
- [x] Step 4: 重跑聚焦測試，預期通過。

### Task 2: 唯讀 ID 增量讀取器

`Files:`

- Modify: `cmd/export/claudemem.go`
- Modify: `cmd/export/claudemem_test.go`

`Interfaces:`

- Consumes: `GetCursorPosition("claude-mem-export")`
- Produces: `claudeMemRead(*model.StateStore, bool) ([]model.Observation, model.CursorPosition, error)`

- [x] Step 1: 新增測試，先匯出 ID `1` / timestamp `100`，再插入 ID `2` / timestamp `100`，斷言第二次只匯出 ID `2`。
- [x] Step 2: 新增測試，將來源設為不存在的 DB，斷言命令回傳錯誤且路徑仍不存在。
- [x] Step 3: 分別執行兩個聚焦測試，預期因 timestamp 查詢漏資料與 SQLite 建立檔案而失敗。
- [x] Step 4: 保留 `Observation.SourceID` 字串契約，在 export reader 以 `strconv.ParseInt` 將來源 ID 轉為遊標值，解析失敗時回傳 wrapped error。
- [x] Step 5: 將來源路徑組為 `mode=ro` URI，取得底層 `*sql.DB`，查詢 `id > last_id ORDER BY id ASC`，查詢後明確關閉。
- [x] Step 6: 讓讀取器返回最後一筆 `{LastTS, LastID}`，並重跑兩個聚焦測試，預期通過。

### Task 3: 命令錯誤與 JSON 契約

`Files:`

- Modify: `cmd/export/claudemem.go`
- Modify: `cmd/export/claudemem_test.go`

`Interfaces:`

- Consumes: `GetCursorPosition` / `SetCursorPosition`
- Produces: 空匯出 JSON `[]\n`

- [x] Step 1: 新增測試，在 writer 寫入成功後刪除 state cursor table，斷言命令回傳 `get claude-mem export cursor` 錯誤。
- [x] Step 2: 新增空來源資料庫測試，斷言輸出為 `[]\n`。
- [x] Step 3: 分別執行兩個聚焦測試，預期因 cursor error 被吞掉與 nil slice 編碼為 `null` 而失敗。
- [x] Step 4: 將 observations 初始化為非 nil 空 slice，並立即檢查、包裝 `GetCursorPosition` 錯誤。
- [x] Step 5: 使用 `SetCursorPosition("claude-mem-export", maxPosition)` 更新遊標，重跑兩個聚焦測試，預期通過。

### Task 4: 文件與完整驗證

`Files:`

- Modify: `README.md`
- Modify: `CLAUDE.md`

- [x] Step 1: 在資料匯出說明記錄 Claude-mem 使用獨立 autoincrement ID 遊標，升級後首次會完整匯出。
- [x] Step 2: 同步 `CLAUDE.md` 關鍵決策與 model mapping。
- [x] Step 3: 執行 `gofmt -w model/cursor.go model/store.go model/store_cursor_test.go cmd/export/claudemem.go cmd/export/claudemem_test.go`。
- [x] Step 4: 執行 `go test ./... -count=1`、`go vet ./...` 與 `git diff --check`，預期全部 exit code `0`。
