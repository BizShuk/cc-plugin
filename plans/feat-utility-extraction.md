# Feature: 公用程式抽取 (Utility Extraction)

## 問題陳述 (Problem Statement)

兩個相關議題需要將工具邏輯自業務模組剝離：

### 議題 A — 重複的 utility 輔助函數

- [cmd/root.go](../cmd/root.go) 內部實作了 `expandPath` 函數
- [model/store.go](../model/store.go#L198) 中以 `ExpandPath`（首字母大寫）又重新實作了一份
- 兩者功能完全相同，造成代碼冗餘
- 更糟的是 `model` 包應只負責純粹的 Domain Model 定義，不應包含路徑展開等底層系統公用工具

### 議題 B — 臨時檔案路徑未收斂

- [config/config.go](../config/config.go#L26) 中預設的 `stores.mempalace.temp_dir` 指向全域的 `/tmp/mempalace-temp`
- 增加了對系統目錄的依賴
- 不符合專案內部 tmp 自封裝原則

## 受影響檔案 (Affected Files)

- `cmd/root.go` — `expandPath` 函數
- `model/store.go` — `ExpandPath` 函數
- `config/config.go` — `stores.mempalace.temp_dir` 預設值

## 實作步驟 (Implementation Steps)

1. 在 `pkg/utils/` 下建立 `path.go`，將 `expandPath` / `ExpandPath` 統一為 `utils.ExpandPath(homeDir, path string) string`
2. 從 `cmd/root.go` 刪除內部的 `expandPath` 實作，改為呼叫 `utils.ExpandPath`
3. 從 `model/store.go` 刪除 `ExpandPath` 方法，改為引用 `pkg/utils/path.go`
4. 修改 `config/default_settings.json`，將 `stores.mempalace.temp_dir` 預設值改為 `~/.cache/cc-plugin/mempalace-temp`（透過 `homedir.Expand` 展開）
5. 在 `cmd/` 初始化階段將展開後的絕對路徑注入至 `MempalaceWriter`

## 驗證方式 (Verification)

- 執行 `go test ./...` 確保測試綠燈
- 驗證 `MempalaceWriter` 寫入時使用的臨時目錄為 `~/.cache/cc-plugin/mempalace-temp` 而非 `/tmp/...`
- 確認 `model/store.go` 不再 import `go-homedir`

## 來源 (Source Plans)

- [`architecture-cc-plugin-evolution.md`](architecture-cc-plugin-evolution.md) §1 診斷 5, §4 遷移映射
- [`architecture-cc-plugin-decoupling.md`](architecture-cc-plugin-decoupling.md) §1 (`工具與通用邏輯放置不當` + `臨時檔案路徑未收斂`)