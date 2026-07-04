# Feature: 讀取邏輯去重 (Reader Logic Deduplication)

## 問題陳述 (Problem Statement)

`cmd/read_logic.go` 與 `cmd/export/` 子目錄中存在高度重複的讀取邏輯，違反 DRY 原則：

- [cmd/read_logic.go](../cmd/read_logic.go#L16-L58) 的 `readGbrainLogic` 與 [cmd/export/gbrain.go](../cmd/export/gbrain.go#L15-L61) 的 `gbrainRead` 邏輯幾乎 100% 相同，僅差在 `all` 參數（是否從 epoch 0 開始讀取）
- [cmd/export/claudemem.go](../cmd/export/claudemem.go#L15-L56) 中的 `claudeMemRead` 與 [cmd/read_logic.go](../cmd/read_logic.go#L60-L103) 中的 `readClaudeMemLogic` 也高度重複，僅差在是否利用 cursor 過濾
- 當未來需新增資料來源（如 `slack-mem`、`notion-mem`）時，必須在兩處重複實作

## 受影響檔案 (Affected Files)

- `cmd/read_logic.go` — `readGbrainLogic`、`readClaudeMemLogic`
- `cmd/export/gbrain.go` — `gbrainRead`
- `cmd/export/claudemem.go` — `claudeMemRead`

## 實作步驟 (Implementation Steps)

1. 在 `internal/service/reader/` 下建立 `gbrain.go` 與 `claudemem.go`，將核心讀取邏輯下沉
2. 新增 `Reader` 介面統一兩個讀取器的行為契約：
   ```go
   type Reader interface {
       Read(ctx context.Context, store *model.StateStore, fromCursor int64) ([]model.Observation, int64, error)
   }
   ```
3. `cmd/distill.go` 改為呼叫 `internal/service/reader/gbrain.go` 的 `Read()` 與 `internal/service/reader/claudemem.go` 的 `Read()`
4. `cmd/export/gbrain.go` 與 `cmd/export/claudemem.go` 改為呼叫對應的內部服務，刪除重複實作

## 驗證方式 (Verification)

- 執行 `go test ./...` 確保測試綠燈
- 對比重構前後 `cc-plugin export gbrain` 與 `cc-plugin distill` 產出的觀察值總數，確保完全一致
- 驗證 `fromCursor` 參數正確控制讀取起點（export 從 0 開始，distill 從上次 cursor 開始）

## 來源 (Source Plans)

- [`architecture-cc-plugin.md`](architecture-cc-plugin.md) §1, §6 Phase 2
- [`architecture-cc-plugin-evolution.md`](architecture-cc-plugin-evolution.md) §1 診斷 2, §4
- [`architecture-cc-plugin-decoupling.md`](architecture-cc-plugin-decoupling.md) §4 遷移映射
- [`architecture-system-modularization.md`](architecture-system-modularization.md) §1 診斷 4, §3, §6 Phase 3
- [`architecture-system-simplification.md`](architecture-system-simplification.md) §1, §6 Phase 2
- [`architecture-cc-distiller.md`](architecture-cc-distiller.md) §1, §6 第二階段