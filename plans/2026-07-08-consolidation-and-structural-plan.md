# 2026-07-08 一致性清理與結構性整合計畫 (Consolidation & Structural Plan)

> 本計畫整合目前仍有效的 CLI 基礎設施、插件、文件與部署工作。
> 已不在專案方向內的領域計畫不再納入工作清單。

## 結論 (Conclusion)

- `plans/` 僅保留仍具獨立價值的設定、日誌、拓撲、插件、TUI 與 Agent 設計工作。
- 三組重複的 feature/architecture 文件需要合併，避免同一工作有兩個真理來源。
- 插件 manifest、README 與實際技能目錄仍需持續保持一致。
- 部署設定需移除本機綁定，確保新環境可直接建置與初始化。

## 範圍 (Scope)

包含：

- CLI 基礎設施：設定外化、結構化日誌、Topology 套件與命令。
- 插件生態：manifest、frontmatter、技能目錄與 README 一致性。
- 文件治理：`plans/`、`docs/specs/`、`README.md`、`README.todo` 的單一來源。
- 部署可移植性：`.gitmodules`、`ecosystem.config.js`、`run.sh`、編輯器設定。
- 獨立產品計畫：Claude Code 使用統計 TUI。
- Agent 設計：feature agent 與專門 agent 的邊界。

排除：

- 已有獨立技術脈絡的外部插件或 submodule 內部實作。
- `data/`、`localtest/`、`tmp/` 等 instance-specific 資料。
- 新增第三方服務或新業務領域。

## 現有計畫盤點 (Current Plan Inventory)

```tree
plans/
├── 2026-06-11-topology-cli.md
├── 2026-07-02-architecture-config-externalization.md
├── 2026-07-02-architecture-structured-logging.md
├── 2026-07-02-structured-logging.md
├── 2026-07-02-architecture-topology-extraction.md
├── 2026-07-02-topology-extraction.md
├── 2026-07-02-skill-md-renaming.md
├── 2026-07-06-architecture-skill-md-renaming.md
├── 2026-07-08-consolidation-and-structural-plan.md
├── bubbly-roaming-sparkle.md
└── agent/
    ├── 2026-05-11-refactor-agent-vs-specialized-agents.md
    └── 2026-05-14-feature-agent-design.md
```

## Phase 1 — CLI 基礎設施 (C 編號)

- [ ] `C1` 依 `2026-07-02-architecture-config-externalization.md` 將預設設定集中於嵌入式設定檔，讓 Go 程式只負責載入與覆寫規則。
- [ ] `C2` 合併 structured-logging 兩份計畫，以 architecture 版本為主，保留可執行的驗證步驟。
- [ ] `C3` 實作 `config.InitLogger`，區分人類可讀輸出與診斷 log，避免把 CLI 資料輸出誤轉成 log。
- [ ] `C4` 合併 topology-extraction 兩份計畫，並與既有 `2026-06-11-topology-cli.md` 對齊命令名稱與公開介面。
- [ ] `C5` 將 Topology 邏輯移至獨立套件，提供 `verify`、`rewrite`、`unlinked` 子命令。

驗證：

- [ ] `go test ./... -count=1`
- [ ] `go vet ./...`
- [ ] `go build -o cc-plugin main.go`
- [ ] `cc-plugin topology --help`

## Phase 2 — 插件生態清理 (P 編號)

- [ ] `P1` 建立 manifest 同步工具，掃描 `plugins/*/skills/*/SKILL.md` 與 `plugins/*/agents/*.md`。
- [ ] `P2` 驗證每個 skill 的 `name` 與目錄一致，`description` 含觸發詞，frontmatter 符合 tier 規範。
- [ ] `P3` 清理空目錄、建置產物與鬆散技能檔案；技能一律使用獨立子目錄與 `SKILL.md`。
- [ ] `P4` 確認每個本地插件都有 `README.md`，並與 manifest、實際目錄同步。
- [ ] `P5` 合併 skill-md-renaming 兩份計畫，以 architecture 版本保留設計脈絡，將執行狀態寫回同一份文件。
- [ ] `P6` 驗證 `.claude-plugin/marketplace.json` 的本地路徑與外部來源均可解析。

## Phase 3 — 文件單一來源 (D 編號)

- [ ] `D1` 合併三組重複計畫：structured-logging、topology-extraction、skill-md-renaming。
- [ ] `D2` 為計畫加入狀態欄位；完成並接受的設計移至 `docs/specs/YYYY-MM-DD-<topic>.md`。
- [ ] `D3` 移除指向已不存在 `architecture-*.md` 的連結，來源證據改指向現存文件或 git history。
- [ ] `D4` 更新 `README.todo`，只保留現存計畫連結與真實進度。
- [ ] `D5` 更新 `README.md` 與 `CLAUDE.md` 的 plans/plugin 結構，避免文件與檔案樹不一致。
- [ ] `D6` `plans/` 只保留 WIP 或未來工作；歷史決策移入 `docs/memory/`，既有設計移入 `docs/specs/`。

## Phase 4 — 部署可移植性 (X 編號)

- [ ] `X1` 驗證 `.gitmodules` 每個條目只有一個正確路徑，未使用項目明確移除。
- [ ] `X2` 將 `ecosystem.config.js` 與編輯器設定中的絕對使用者路徑改為 repo-relative 或環境變數。
- [ ] `X3` 為 `run.sh` 加入 `set -euo pipefail`，移除重複區塊，並為既有設定檔提供備份守衛。
- [ ] `X4` 驗證 `.mcp.json`、`.lsp.json`、VS Code launch/settings 與實際入口一致。
- [ ] `X5` 移除根目錄空 scaffolding 與失效 `.gitignore` 規則。

## Phase 5 — 獨立產品與 Agent 計畫 (U/A 編號)

- [ ] `U1` 審核 `bubbly-roaming-sparkle.md` 的 TUI 資料模型、掃描上限與端對端驗證，再決定是否進入實作。
- [ ] `A1` 合併兩份 Agent 設計中的重複角色定義，保留可驗證的委派邊界與輸出契約。
- [ ] `A2` Agent 設計若已被插件實作，將完成內容移至 `docs/specs/`，不在 `plans/` 長期保留。

## 執行順序 (Execution Order)

```tree
Phase 1 CLI 基礎設施
└── Phase 2 插件生態
    └── Phase 3 文件單一來源
        └── Phase 4 部署可移植性

Phase 5 TUI / Agent 計畫可獨立評估
```

## 風險 (Risks)

| 風險 | 影響 | 緩解 |
| :--- | :--- | :--- |
| 合併計畫時遺失歷史決策 | 中 | 依賴 git history，不複製已失效內容 |
| 日誌改造改變 CLI stdout 契約 | 中 | 資料輸出留在 stdout，診斷資訊走 logger |
| Topology 搬遷破壞 import | 中 | 先搬測試，再改 package 與 CLI 接線 |
| manifest 自動同步覆蓋手寫 metadata | 中 | 僅管理 `skills`/`agents` 欄位並先提供 dry-run |
| `run.sh` 改動影響既有 symlink | 中 | 實作備份守衛與可重入測試 |

## 驗收 (Acceptance Criteria)

- [ ] `plans/` 無已刪檔案的連結。
- [ ] 同一主題只有一份 active plan。
- [ ] `README.todo`、`CLAUDE.md` 與實際 `plans/` 樹一致。
- [ ] `go test ./... -count=1`、`go vet ./...`、`go build ./...` 通過。
- [ ] 每個本地插件的 README、manifest、skills 與 agents 清單一致。
- [ ] 新環境 clone 後可執行 `go build ./...` 與 `./run.sh`。
