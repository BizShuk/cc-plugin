# Feature: 迴圈內 defer 資源洩漏修復 (Defer-in-Loop Resource Leak Fix)

## 問題陳述 (Problem Statement)

Go 的 `defer` 關鍵字僅於函式返回時執行，而非於當前迴圈迭代結束時執行：

- [cmd/write_agentmemory.go](../cmd/write_agentmemory.go#L34) 的 `writeAgentMemoryLogic` 中，`defer resp.Body.Close()` 被放置於 `for` 迴圈內
- 當 `memories` 陣列龐大時，多個 HTTP 回應主體同時保持開啟
- 會導致：
  - HTTP 連線與檔案描述符 (File Descriptor) 累積
  - 系統觸發 `too many open files` 錯誤
  - 高機率引發記憶體洩漏

## 受影響檔案 (Affected Files)

- `cmd/write_agentmemory.go` — `writeAgentMemoryLogic` 第 34 行

## 實作步驟 (Implementation Steps)

1. 將 `writeAgentMemoryLogic` 中的 `http.Post` 呼叫封裝至匿名函數 (IIFE, Immediately Invoked Function Expression)，確保 `resp.Body.Close()` 在每次迭代結束時立即被呼叫：
   ```go
   for _, m := range memories {
       func() error {
           resp, err := http.Post(...)
           if err != nil {
               return err
           }
           defer resp.Body.Close()  // 此 defer 在匿名函數返回時立即執行
           // ... 處理 resp
           return nil
       }()
   }
   ```
2. 或改用手動 `resp.Body.Close()` 在每次迭代結束時立即呼叫
3. 補充單元測試，模擬大量記憶寫入（1000+ 筆）並驗證檔案描述符未持續攀升

## 驗證方式 (Verification)

- 執行 `go test ./...` 確認測試綠燈
- 使用 `lsof -p <pid>` 或 `lsof -i` 觀察 distill 執行時的檔案描述符狀態，確認批次寫入後無持續攀升
- 背景執行 1000 筆以上記憶寫入，無 `too many open files` 錯誤

## 來源 (Source Plans)

- [`architecture-cc-plugin.md`](architecture-cc-plugin.md) §1, §6 Phase 2
- [`architecture-cc-plugin-reorganization.md`](architecture-cc-plugin-reorganization.md) §1
- [`architecture-cc-plugin-evolution.md`](architecture-cc-plugin-evolution.md) §1 診斷 4, §6 Phase 2
- [`architecture-cc-plugin-decoupling.md`](architecture-cc-plugin-decoupling.md) §1, §6 Phase 2
- [`architecture-system-modularization.md`](architecture-system-modularization.md) §1 診斷 5, §3, §6 Phase 1
- [`architecture-system-simplification.md`](architecture-system-simplification.md) §1, §6 Phase 2
- [`architecture-cc-distiller.md`](architecture-cc-distiller.md) §1, §6 第一階段