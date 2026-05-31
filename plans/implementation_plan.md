# 使用 GORM 重構狀態儲存與 claude-mem 讀取機制 (Use GORM to Refactor State Storage and claude-mem Reading Mechanism)

將原本使用原生 `database/sql` 實作的 `cmd/state.go` 及 `cmd/read_claudemem.go` 重構為使用 `GORM` 物件關係對映 (Object-Relational Mapping, ORM) 與 `SQLite` 驅動，使其架構更簡潔、更易於維護，並完整相容現行的測試。

## 使用者審查項目 (User Review Required)

> [!IMPORTANT]
> 此次修改將全面替換原生 `database/sql` 操作為 `GORM` 框架，涉及的核心檔案有 `state.go` 與 `read_claudemem.go`。
> 我們將使用 `GORM` 的自動遷移 `AutoMigrate` 功能取代原生的 `initSchema` SQL 字串。

## 開放性問題 (Open Questions)

> [!NOTE]
> 目前沒有懸而未決的開放性問題。實作將完全相容現有的 API 界面 and 測試邏輯。

---

## 預期變更內容 (Proposed Changes)

### 依賴管理 (Dependency Management)

#### [MODIFY] [go.mod](file:///Users/shuk/projects/cc-plugin/go.mod)
將 `gorm.io/gorm` 和 `gorm.io/driver/sqlite` 從間接依賴 (`indirect`) 升級為直接依賴 (`require`)。

---

### 狀態儲存模組 (State Store Module)

#### [MODIFY] [state.go](file:///Users/shuk/projects/cc-plugin/cmd/state.go)
- 使用 `gorm.io/gorm` 和 `gorm.io/driver/sqlite` 取代原生的 `database/sql`。
- 定義三個 `GORM` 模型模型：
  - `Cursor` 代表 `cursor` 資料表。
  - `Seen` 代表 `seen` 資料表（複合主鍵：`Fingerprint` 與 `Source`）。
  - `Distilled` 代表 `distilled` 資料表（複合主鍵：`Source` 與 `SourceID`）。
- 使用 `db.AutoMigrate` 初始化資料庫 Schema。
- 使用 `GORM` 的查詢與寫入 API 來實作以下方法：
  - `GetCursor` / `SetCursor`
  - `RecordSeen`
  - `AlreadyDistilled` / `MarkDistilled`
  - `DueForPrune` / `DropDistilled`
- 使用 `Close()` 方法：因為 `gorm.DB` 通常不需手動關閉，但為了維持向後相容性，我們會透過底層的 `sql.DB` 取得連接並進行 `Close`。

#### [MODIFY] [state_test.go](file:///Users/shuk/projects/cc-plugin/cmd/state_test.go)
- 調整測試中的資料庫初始化邏輯，使其搭配 `GORM` 版本進行運作。

---

### Claude 記憶讀取模組 (Claude Memory Reader Module)

#### [MODIFY] [read_claudemem.go](file:///Users/shuk/projects/cc-plugin/cmd/read_claudemem.go)
- 定義 `ClaudeMemObservation` 作為對應外部 `claude-mem` 的模型。
- 使用 `GORM` 開啟外部資料庫，以物件化方式查詢 `observations` 內容。

#### [MODIFY] [main_test.go](file:///Users/shuk/projects/cc-plugin/cmd/main_test.go)
- 調整 `TestReadClaudeMemLogic` 的 mock 資料庫建立方式與讀取比對，以配合 `GORM` 版本的實作。

---

## 驗證計畫 (Verification Plan)

### 自動化測試 (Automated Tests)
我們將在終端機執行 Go 單元測試，確保所有測試皆能正確通過：
```bash
go test -v ./cmd/...
```
確保 `TestStateStore`、`TestReadClaudeMemLogic` 還有 `TestReadGbrainLogic` 均能正常執行並通過。

### 手動驗證 (Manual Verification)
- 由於此重構僅限於底層資料庫存取層，只需確保單元測試 100% 通過，且專案建置正常即可：
```bash
go build -o /dev/null main.go
```
