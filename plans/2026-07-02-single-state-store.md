# Feature: 單一 StateStore 連線 (Single StateStore Connection)

## 問題陳述 (Problem Statement)

同一個 SQLite 資料庫在同一處理流程中被重複開啟和關閉多次，違反單一連線 (Single Connection) 原則：

- [cmd/distill.go](../cmd/distill.go#L35) 已建立 `store` 實例，但 [cmd/read_logic.go](../cmd/read_logic.go#L60) 中的 `readClaudeMemLogic` 卻在內部又自行呼叫 `NewStateStore` 並於結束時 `store.Close`
- [cmd/retain.go](../cmd/retain.go#L31) 內部的 `retainLogic` 也是自行開啟並關閉 `store` 連線
- 這造成 SQLite 連線頻繁開啟與關閉，浪費系統資源
- 在高併發或排程重疊時容易引發 `database is locked` 錯誤

此外，`model/store.go` 中的 [`NewStateStore`](../model/store.go#L27) 直接呼叫 `viper.GetString("state.db_path")`，導致領域層強耦合全域設定，無法在不載入設定檔的情況下獨立進行單元測試。

## 受影響檔案 (Affected Files)

- `cmd/distill.go` — 主流程建立 store
- `cmd/read_logic.go` — 內部自行 `NewStateStore`
- `cmd/retain.go` — 內部自行 `NewStateStore`
- `model/store.go` — `NewStateStore` 簽名與 viper 依賴

## 實作步驟 (Implementation Steps)

1. 修改 `model/store.go` 的 `NewStateStore` 簽名為 `NewStateStore(dbPath string)`，移除 `github.com/spf13/viper` 直接依賴
2. 在 `cmd/distill.go` 初始化階段建立單一 `StateStore` 實例
3. 修改 `readClaudeMemLogic` 與 `retainLogic` 函式簽名，使其接受外部傳入的 `*StateStore` 參數，移除其內部的 `NewStateStore` 與 `Close`
4. 將 `dbPath` 解析從 `model/` 層移至 `cmd/` 層（由 CLI 負責讀 viper 並呼叫 `homedir.Expand`）

## 驗證方式 (Verification)

- 執行 `go test ./model/...` 與 `go test ./cmd/...` 確保測試綠燈
- 背景執行 1000 筆以上資料時，無 SQLite `database is locked` 錯誤
- 單元測試可使用 `:memory:` SQLite 而無需依賴全域設定

## 來源 (Source Plans)

- [`architecture-cc-plugin.md`](architecture-cc-plugin.md) §1, §6 Phase 1
- [`architecture-cc-plugin-reorganization.md`](architecture-cc-plugin-reorganization.md) §1
- [`architecture-cc-plugin-evolution.md`](architecture-cc-plugin-evolution.md) §1 診斷 1, §6 Phase 1
- [`architecture-cc-plugin-decoupling.md`](architecture-cc-plugin-decoupling.md) §1, §6 Phase 2
- [`architecture-system-modularization.md`](architecture-system-modularization.md) §1 診斷 3, §3
- [`architecture-system-simplification.md`](architecture-system-simplification.md) §1, §6 Phase 1
- [`architecture-cc-distiller.md`](architecture-cc-distiller.md) §1, §6 第一階段