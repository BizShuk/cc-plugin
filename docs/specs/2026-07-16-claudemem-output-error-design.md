# Claude-mem 匯出錯誤處理設計

## 目標

`cc-plugin export claudemem` 只有在 JSON 完整寫入命令輸出後才能推進 `claude-mem` 遊標。

## 設計

- 使用 `json.Encoder` 將觀察資料寫入 `cmd.OutOrStdout()`。
- 保留現有兩個空白的縮排 JSON 格式。
- `Encode` 失敗時以 `%w` 包裝並回傳錯誤，不執行遊標更新。
- 本次不變更時間戳遊標、SQLite 連線或空陣列行為。

## 驗收

- 輸出 writer 回傳錯誤時，命令回傳錯誤。
- 同一情境下 `claude-mem` 遊標保持原值。
- 現有專案測試維持通過。
