# Global Rule

`用繁體中文 + English Terminology`

## Principles

- one file one responsibility. one package/folder one domain
- 遇到執行錯誤時，先嘗試修復，最多重試 5 次；若仍無法解決則明確報錯並停止, 多次遭遇相同錯誤/問題時，將解法記錄至 Memory
- git worktree branch should be name with `w-<feature>` or `agent-<feature>`

### 上下文 (Context)

- 載入 `@./CLAUDE.md` 作為專案結構
- 載入 `@./README.md` 作為業務核心

## 目錄佈局 (Directory Layout)

`~/projects/` 是`兩層 (two-level)` 結構：專案可放在根目錄，也可放在分類目錄之下。

```tree
~/projects/
├── <project>/              # 專案 (Project)：repo 根目錄，具備統一介面
└── <category>/             # 分類 (Category)：純容器
    └── <project>/          # 分類下的專案，統一介面規則完全相同
```

分類 (Category) 現有 13 個：`ai`、`collections`、`data`、`env_setup`、`game`、
`iphone`、`platform`、`playground`、`product`、`research`、`social`、`tools`、`web`。

規則：

- 分類目錄`可以`有自身 `README.md`／`CLAUDE.md` 作為領域導覽
- 分類深度`固定一層`：不得出現 `<category>/<category>/<project>`。

## 統一介面 (Unified Interface)

每個 repo（含 monorepo 內的子專案）必須具備：

| 檔案                        | 必要性 | 職責                                                                                      |
| --------------------------- | ------ | ----------------------------------------------------------------------------------------- |
| `README.md`                 | 必備   | 業務定義 (business definition)、domain flow                                               |
| `CLAUDE.md`                 | 必備   | 技術脈絡 (technical context)、結構、關鍵決策                                              |
| `AGENTS.md`                 | 必備   | 軟連結 `AGENTS.md -> CLAUDE.md`（一律建立，不例外）                                       |
| `run.sh`                    | 選備   | 預設執行metadata setup(not running program)，可隨時執行 (default process, can run always) |
| `ecosystem.config.js`       | 選備   | 常駐程序或 cron 任務，置於 repo 根目錄由 pm2 管理                                         |
| `README.todo`               | 必備   | 待辦事項 (pending todo item)                                                              |
| `~/.config/<app_name>/`     | 必備   | 透過gosdk config.Default 固定配置文件                                                     |
| `~/.config/<app_name>/data` | 必備   | 透過gosdk config.Default 固定配置文件                                                     |
| `~/.config/<app_name>/logs` | 必備   | pm2 task logs                                                                             |
| `plans/`                    | 選備   | 進行中計畫，命名 `YYYY-MM-DD-<topic>.md` topic name should be meaningful to the change    |
| `docs/terminology.md`       | 必備   | 術語表 (terminology)：領域名詞、縮寫、狀態值的單一定義來源                                |
| `docs/memory/`              | 必備   | 歷史操作跟決策 retrospective                                                              |
| `docs/backlog/`             | 選備   | 待辦想法 (pending ideas)                                                                  |
| `docs/specs/`               | 選備   | 既有設計與規格 (existing design)，統一 `YYYY-MM-DD-<topic>.md`                            |
| `docs/tutorials/`           | 選備   | 領域知識學習與概念導覽 (domain tutorials)；專案架構/流程/環境等常規指南放 `docs/` 根層    |
| `scripts/`                  | 選備   | 專案相關腳本 (project related script)                                                     |
| `tmp/`                      | 選備   | 實例專屬之資料與設定 (data/config per instance, not source code/logic)                    |

- golang-dev skill for golang structure
- pm2 skill for pm2 structure

### 內容歸屬 (Content Ownership)

`README.md` 是`為什麼用它、怎麼開始`；`CLAUDE.md` 是`邊界是什麼、誰擁有什麼`。
一個事實只能有`一個` owner —— 重複的兩份必然分岔，結果是兩份都不準。

歸屬用單一問題判定：「這句話會因為什麼而變成假的？」

| 失效原因                                          | 歸屬                                                                   |
| ------------------------------------------------- | ---------------------------------------------------------------------- |
| 別的 repo 改動                                    | 那個 repo；本檔只留`消費端契約`（誰能 import、誰擁有哪張對照、優先序） |
| 時間經過（已經發生的事）                          | `docs/CHANGELOG.md`                                                    |
| 還沒做完                                          | `README.todo`                                                          |
| 一次 commit（行數、檔案數、byte 上限、module 數） | 刪除；量測值屬於程式碼與測試                                           |
| 程式碼改了，且可用指令驗證                        | 測試或 `scripts/`；文件只留一句規則 + 測試名稱                         |
| 換一台機器（絕對路徑）                            | 改相對路徑或 `$(git rev-parse --show-toplevel)`                        |
| 都不會失效（不變式）                              | 留在 `CLAUDE.md`                                                       |

- 結構樹、ownership、架構決策由 `CLAUDE.md` 單一擁有，`README.md` 用一行指過去。
- `README.md` 出現需要先讀原始碼才懂的型別名，就是越界。
- 寫得出 pass/fail 的斷言不放 Markdown ——`沒人執行的斷言會腐爛成錯誤資訊`。
  先把斷言自動化，才有資格刪文件。
- 稽核與瘦身用 `[[docs-consolidation]]` 的`範疇清理 (Scope Cleanup)` 模式。

### 檔案命名 (File Naming)

#### `plans/` 與 `docs/specs/`

統一採用 `YYYY-MM-DD-<topic>.md` 格式：

- 日期：本地時區（Asia/Taipei），與既有 `docs/memory/` convention 一致
- `<topic>`：kebab-case 英文 topic name，必須 meaningful to the change（例：`windows-11-desktop-receiver`）
- 範例：✅ `2026-07-25-windows-11-desktop-receiver.md`；❌ `delegated-wishing-spark.md`
- 反例（system-generated slug，禁止使用）：`hashed-dancing-pascal.md`、`partitioned-rolling-gosling.md`

例外：plan-mode 系統啟動時自動產生的暫存檔（隨機 slug 命名）僅在 plan mode 內部使用，**離開 plan mode 前必須** rename 為正式檔名，並同步更新所有 cross-reference（README、CLAUDE.md、specs、plans）。
