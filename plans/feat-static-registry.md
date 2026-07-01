# Feature: Reader/Writer 靜態註冊表 (Static Registry Pattern)

## 問題陳述 (Problem Statement)

cc-plugin 未來可能需要支援更多記憶來源（例如 `Slack`、`Notion`、`Apple Notes`）或寫入至不同的儲存端。若每次新增來源都需修改核心 `DistillerService`，將違反開放封閉原則 (Open-Closed Principle)。

## 受介面契約 (Interface Contracts)

```go
// pkg/reader/registry.go
var readerRegistry = map[string]func(ctx context.Context, store *model.StateStore, fromCursor int64) ([]model.Observation, int64, error){}

func RegisterReader(name string, factory func(ctx context.Context, store *model.StateStore, fromCursor int64) ([]model.Observation, int64, error)) {
    readerRegistry[name] = factory
}
```

## 受影響檔案 (Affected Files)

- `internal/service/reader/` — 新增 `registry.go`
- `config/settings.json` — 新增 `distiller.sources` 陣列設定
- `internal/service/distiller/orchestrator.go` — 透過註冊表動態選擇 reader

## 實作步驟 (Implementation Steps)

1. 在 `config/settings.json` 定義配置：
   ```json
   {
     "distiller": {
       "sources": ["gbrain", "claude-mem"],
       "stores": ["agentmemory", "mempalace"]
     }
   }
   ```
2. 在 `internal/service/reader/` 建立 `registry.go`，提供 `RegisterReader` 與 `GetReader` 函式
3. 在 `gbrain.go` 與 `claudemem.go` 的 `init()` 中呼叫 `RegisterReader` 註冊自己
4. 在 `DistillerService` 啟動時讀取 `distiller.sources` 設定，從註冊表取得對應 reader 實例
5. 類似地為 `Writer` 設計 `writerRegistry`

## 驗證方式 (Verification)

- 執行 `go test ./...` 確保測試綠燈
- 修改 `settings.json` 將 `sources` 改為 `["gbrain"]`，驗證 distill 是否正確略過 `claude-mem` 來源
- 確認 `DistillerService` 不再 hard-code `gbrain` / `claude-mem` 字串

## 來源 (Source Plans)

- [`architecture-cc-distiller.md`](architecture-cc-distiller.md) §5