# 專案工作區 (Projects Workspace) — 管理方向 (Management Direction)

本檔為 `~/projects/` 全工作區的方向性規範 (direction)，所有 repo 的 session 均會繼承。
各 repo 細節以其自身 `CLAUDE.md` 為準；本檔只定規則與分層，不放具體內容。

## 工作規則 (Working Rules)

### 上下文 (Context)

- 載入 `@./CLAUDE.md` 作為專案結構
- 載入 `@./README.md` 作為業務核心

### 習慣 (Behaviors)

- 遇到執行錯誤時，先嘗試修復，最多重試 5 次；若仍無法解決則明確報錯並停止
- git worktree branch should be name with `w-<feature>` or `agent-<feature>`

### 慣例 (Convention)

既有慣例 (convention) 不得額外建立自訂選項。
例如 `gosdk` 已有 app log dir / app config dir / app data dir，service layer 不得再建 config path 指向 data dir。

### 演化 (Evolution)

- 多次遭遇相同錯誤/問題時，將解法記錄至 Memory

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

| 檔案                        | 必要性 | 職責                                                                                   |
| --------------------------- | ------ | -------------------------------------------------------------------------------------- |
| `README.md`                 | 必備   | 業務定義 (business definition)、domain flow                                            |
| `CLAUDE.md`                 | 必備   | 技術脈絡 (technical context)、結構、關鍵決策                                           |
| `AGENTS.md`                 | 必備   | 軟連結 `AGENTS.md -> CLAUDE.md`（一律建立，不例外）                                    |
| `run.sh`                    | 選備   | 預設執行程序，可隨時執行 (default process, can run always)                             |
| `ecosystem.config.js`       | 選備   | 常駐程序或 cron 任務，置於 repo 根目錄由 pm2 管理                                      |
| `README.todo`               | 必備   | 待辦事項 (pending todo item)                                                           |
| `~/.config/<app_name>/`     | 必備   | 透過gosdk config.Default 固定配置文件                                                  |
| `~/.config/<app_name>/data` | 必備   | 透過gosdk config.Default 固定配置文件                                                  |
| `~/.config/<app_name>/logs` | 必備   | pm2 task logs                                                                          |
| `plans/`                    | 選備   | 進行中計畫，命名 `YYYY-MM-DD-<topic>.md` topic name should be meaningful to the change |
| `docs/terminology.md`       | 必備   | 術語表 (terminology)：領域名詞、縮寫、狀態值的單一定義來源                             |
| `docs/memory/`              | 必備   | 歷史操作跟決策 retrospective                                                           |
| `docs/backlog/`             | 選備   | 待辦想法 (pending ideas)                                                               |
| `docs/specs/`               | 選備   | 既有設計與規格 (existing design)，統一 `YYYY-MM-DD-<topic>.md`                         |
| `docs/tutorials/`           | 選備   | 領域知識學習與概念導覽 (domain tutorials)；專案架構/流程/環境等常規指南放 `docs/` 根層 |
| `scripts/`                  | 選備   | 專案相關腳本 (project related script)                                                  |
| `tmp/`                      | 選備   | 實例專屬之資料與設定 (data/config per instance, not source code/logic)                 |
