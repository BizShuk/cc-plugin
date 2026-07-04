# 架構重構執行藍圖 (Refactoring Execution Roadmap)

本藍圖將 `plans/2026-07-02-*.md` 中各項功能對應到四個可獨立交付且可回滾的階段。每個階段皆須保留獨立的 git commit 並通過驗證測試。

## Phase 1 — 建立單一 store 連線與設定純淨化 (30% 進度)

**涵蓋功能：**
- [`2026-07-02-single-state-store.md`](2026-07-02-single-state-store.md) — 單一 StateStore 連線與 viper 解耦
- [`2026-07-02-config-externalization.md`](2026-07-02-config-externalization.md) — 硬編碼預設移至 `default_settings.json`
- [`2026-07-02-resource-leak-fix.md`](2026-07-02-resource-leak-fix.md) — 修復 defer-in-loop 資源洩漏（可與 Phase 2 一併處理）

**目標：**
- 修改 `NewStateStore` 接受 `dbPath` 參數，移除 `viper` 直接依賴
- 調整 `readClaudeMemLogic` 與 `retainLogic` 使引導同一個 `store` 實例，消滅重複連線
- 將 `config.go` 預設設定移入 `default_settings.json`

**驗證：**
- 執行 `go test ./model/...` 與 `go test ./cmd/...` 確保測試綠燈
- 功能無 regression，手動執行 `cc-plugin distill` 正常

---

## Phase 2 — 消除重複代碼與服務層初步抽取 (55% 進度)

**涵蓋功能：**
- [`2026-07-02-reader-deduplication.md`](2026-07-02-reader-deduplication.md) — gbrain/claudemem 讀取邏輯合併至 `internal/service/reader/`
- [`2026-07-02-utility-extraction.md`](2026-07-02-utility-extraction.md) — `ExpandPath` 重複整合至 `pkg/utils/path.go`，temp_dir 收斂

**目標：**
- 將 `readGbrainLogic` 提取至 `internal/service/reader/gbrain.go`，讓 `cmd/export/gbrain.go` 與 `cmd/distill.go` 共用
- 將 `ExpandPath`/`expandPath` 統一至 `pkg/utils/`，並修正 `stores.mempalace.temp_dir` 預設值

**驗證：**
- 執行單元測試
- 對比 export 與 distill 產出的觀察值總數與重構前一致
- 大量資料匯入時檔案描述符 (File Descriptor) 無持續攀升

---

## Phase 3 — 建立核心服務層與介面抽離 (80% 進度)

**涵蓋功能：**
- [`2026-07-02-service-layer-extraction.md`](2026-07-02-service-layer-extraction.md) — 業務邏輯自 `cmd/` 移至 `internal/service/`
- [`2026-07-02-static-registry.md`](2026-07-02-static-registry.md) — Reader/Writer 靜態註冊表
- [`2026-07-02-structured-logging.md`](2026-07-02-structured-logging.md) — `fmt.Printf` 遷移至 `slog`

**目標：**
- 在 `internal/service/` 建立 `DistillerService`、`ReaderService`、`WriterService`
- 將 `cmd/distill.go` 的蒸餾編排邏輯轉移至服務層
- 定義 `Extractor`、`Reader`、`Writer` 介面契約
- 透過靜態註冊表動態載入 reader/writer
- 將日誌輸出改為結構化格式

**驗證：**
- 對 `DistillerService` 補上特徵測試 (Characterization Test) 與 Mock 測試
- 修改 `settings.json` 切換 source，驗證 distill 正確略過未註冊來源
- 輸出為合法 JSON 格式，可被 `jq` 解析

---

## Phase 4 — 規範檔案更名與整合驗證 (100% 進度)

**涵蓋功能：**
- [`2026-07-02-skill-md-renaming.md`](2026-07-02-skill-md-renaming.md) — `anti-sabotage.md` 重命名為 `SKILL.md` 並加上 frontmatter
- [`2026-07-02-topology-extraction.md`](2026-07-02-topology-extraction.md) — `model/topology` 遷移至 `pkg/topology`，新增 `cmd/topology` CLI

**目標：**
- 重命名 `plugins/general/skills/anti-sabotage/anti-sabotage.md` 為 `SKILL.md`，補上標準 YAML frontmatter
- 將 `model/topology.go`、`model/topology_ops.go` 與相關測試移至 `pkg/topology/`
- 新增 `cc-plugin topology` 子命令（`verify`、`rewrite`、`unlinked`）
- 更新 `README.md` 與 `CLAUDE.md` 中的專案結構描述

**驗證：**
- 執行 `npx skills add .` 確認無異常
- 執行 `go test ./...` 驗證整個工作區 Go 套件編譯綠燈
- 執行 `cc-plugin topology verify --root <path>` 驗證拓撲檢查功能
- 全鏈路集成測試：`cc-plugin distill` 成功且所有測試綠燈

---

## 風險與回滾策略 (Risks & Rollback)

每個 Phase 結束時必須：

1. 為本次重構建立獨立的 `git commit`
2. 確認 `go test ./...` 全綠
3. 手動執行關鍵指令驗證外部行為無變化

若驗證失敗，立即執行：

```bash
git reset --hard HEAD~1
```

回滾至上一個穩定狀態，並檢視失敗原因後再開新分支重試。

詳細風險與對策請參閱各 feature 檔案中的「驗證方式」段落。