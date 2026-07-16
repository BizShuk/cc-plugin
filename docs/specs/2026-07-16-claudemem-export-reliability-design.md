# Claude-mem 匯出可靠性設計

## 目標

完成 `cmd/export/claudemem.go` review 的剩餘修正：穩定增量遊標、完整錯誤處理、唯讀 SQLite 連線與穩定空陣列 JSON 契約。

## 遊標方案

實際 Claude-mem schema 將 `observations.id` 定義為 `INTEGER PRIMARY KEY AUTOINCREMENT`。

1. `ID 遊標`：依插入順序增量查詢，不受相同或回填時間戳影響。採用。
2. `(timestamp, id)` 複合遊標：可解決相同 timestamp，但仍可能漏掉時間戳小於遊標的後寫資料。
3. `timestamp >= cursor`：無需 schema 變更，但會永久重複匯出最後一批資料。

## 資料與流程

- `Cursor` 新增 `last_id` 欄位，預設為 `0`，由 GORM `AutoMigrate` 升級現有 state DB。
- 匯出使用獨立 source key `claude-mem-export`，不再與 distill 的 `claude-mem` timestamp 遊標共用同一筆狀態。
- 增量查詢使用 `id > last_id ORDER BY id ASC`。
- 第一次升級時新遊標為 `0`，會安全地重新匯出現有資料一次，不從舊 timestamp 推測 ID。
- JSON 寫入成功後才讀取與更新遊標；遊標讀取錯誤不得忽略。

## SQLite 生命週期

- 將來源路徑組成 SQLite URI，以 `mode=ro` 開啟。
- 查詢後明確關閉底層 `*sql.DB`；查詢或關閉錯誤使用 `%w` 包裝。
- 資料庫不存在時，匯出回傳錯誤且不建立空檔案。

## JSON 契約

- 無資料時回傳 `[]`，不回傳 `null`。
- 保留現有兩個空白的縮排格式與結尾換行。

## 測試

- 後寫入相同 timestamp 的新 ID 會在下次匯出出現。
- 遊標讀取失敗會回傳對應錯誤。
- 不存在的來源 DB 不會被匯出命令建立。
- 空資料庫輸出 `[]`。
- 先前「輸出失敗不推進遊標」的測試繼續通過。
