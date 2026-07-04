# Feature: 設定檔外化與純淨化 (Configuration Externalization)

## 問題陳述 (Problem Statement)

預設設定值被硬編碼於 Go 程式碼中，違反「設定與程式碼分離」原則：

- [config/config.go](../config/config.go#L18-L27) 的 `Init()` 函數將所有預設設定寫死（資料庫路徑、Ollama 預設模型、儲存 API 端點等）
- [config/default_settings.json](../config/default_settings.json) 內容為空 `{}`
- 外部管理（修改預設值）需重新編譯 Go 程式
- 多個模型與服務元件（如 [`model/store.go`](../model/store.go#L27) 的 `NewStateStore`）直接呼叫 `viper.GetString` 獲取配置，強耦合全域設定

## 受影響檔案 (Affected Files)

- `config/config.go` — `Init()` 函數中的硬編碼預設值
- `config/default_settings.json` — 空的預設設定檔
- `model/store.go` — `viper` 全域依賴

## 實作步驟 (Implementation Steps)

1. 將 `config.go` 中寫死的預設設定值遷移至 `default_settings.json`：
   - `state.db_path` — SQLite 狀態檔路徑
   - `ollama.endpoint` — Ollama HTTP API 端點
   - `ollama.model` — 預設 LLM 模型
   - `agentmemory.endpoint` — agentmemory API 端點
   - `mempalace.cli` — mempalace CLI 路徑
2. 修改 `config.Init()` 使用 `viper.SetDefault()` 從 `default_settings.json` 載入預設值
3. 移除 `model/store.go` 內部的 `viper.GetString` 直接呼叫，改為透過建構式注入 `dbPath`（見 [2026-07-02-single-state-store.md](2026-07-02-single-state-store.md)）
4. 驗證 `default_settings.json` 與環境變數覆蓋機制（環境變數優先）

## 驗證方式 (Verification)

- 執行 `cc-plugin distill --config /path/to/custom.json` 確認自訂設定可生效
- 設定環境變數 `CC_OLLAMA_MODEL` 後，`viper.GetString("ollama.model")` 應回傳環境變數值
- 執行 `go test ./...` 確保測試綠燈

## 來源 (Source Plans)

- [`architecture-cc-plugin.md`](architecture-cc-plugin.md) §1, §6 Phase 1
- [`architecture-cc-plugin-reorganization.md`](architecture-cc-plugin-reorganization.md) §1
- [`architecture-cc-plugin-evolution.md`](architecture-cc-plugin-evolution.md) §1 診斷 3, §6 Phase 1
- [`architecture-cc-plugin-decoupling.md`](architecture-cc-plugin-decoupling.md) §1
- [`architecture-system-modularization.md`](architecture-system-modularization.md) §1, §6 第三階段
- [`architecture-system-simplification.md`](architecture-system-simplification.md) §1, §6 Phase 1
- [`architecture-cc-distiller.md`](architecture-cc-distiller.md) §1, §6 第三階段