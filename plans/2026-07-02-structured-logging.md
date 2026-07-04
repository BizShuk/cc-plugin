# Feature: 結構化日誌導入 (Structured Logging Adoption)

## 問題陳述 (Problem Statement)

當前系統廣泛使用 `fmt.Printf` 與 `fmt.Fprintf(os.Stderr, ...)` 進行簡易輸出：

- 缺乏結構化欄位（時間、嚴重性、模組），難以被日誌聚合系統解析
- 若將 CLI 工具部署於 `Crontab` 作為每日定時任務，將難以與現代觀測系統（如 Elastic Stack 或 Grafana Loki）對接
- 無法依日誌等級 (`debug` / `info` / `warn` / `error`) 過濾輸出

## 受影響檔案 (Affected Files)

- `cmd/*.go` — 所有使用 `fmt.Printf` / `fmt.Fprintf` 輸出的檔案
- `config/config.go` — 新增 `log.level` 與 `log.format` 設定項

## 實作步驟 (Implementation Steps)

1. 引入 Go 1.21+ 標準庫 `log/slog`（無外部依賴）
2. 在 `config.Init()` 中根據 `log.level` 設定建立 `slog.Handler`（預設為 JSON 格式）
3. 將所有 `cmd/*.go` 中的 `fmt.Printf(...)` 替換為 `slog.Info(...)` / `slog.Error(...)`
4. 將所有 `fmt.Fprintf(os.Stderr, ...)` 替換為 `slog.Error(...)` 並以 `fmt.Errorf("context: %w", err)` 包裝錯誤
5. 新增 `log.level` 設定項，支援 `debug` / `info` / `warn` / `error`

## 驗證方式 (Verification)

- 執行 `cc-plugin distill` 確認輸出為合法 JSON 格式
- 設定 `log.level=debug` 後能輸出額外的診斷資訊
- 將輸出導入 `jq` 解析無錯誤
- 若部署於 Crontab，輸出可被 Loki/Filebeat 直接收集

## 來源 (Source Plans)

- [`architecture-system-modularization.md`](architecture-system-modularization.md) §1 診斷 6