# 內容歸屬與檔案命名 (Content Ownership and File Naming)

本檔是`內容歸屬判準`與 `plans/`／`docs/specs/` 檔名規範的單一 owner。
`config/CLAUDE.global.md`（全域規則）與本技能的 `SKILL.md` 都只留指標指向這裡。

## 內容歸屬 (Content Ownership)

`README.md` 是`為什麼用它、怎麼開始`；`CLAUDE.md` 是`邊界是什麼、誰擁有什麼`。
一個事實只能有`一個` owner —— 重複的兩份必然分岔，結果是兩份都不準。

歸屬用單一問題判定：「這句話會因為什麼而變成假的？」

| 失效原因                                          | 歸屬                                                                   | 典型徵狀                                            |
| ------------------------------------------------- | ---------------------------------------------------------------------- | --------------------------------------------------- |
| 別的 repo 改動                                    | 那個 repo；本檔只留`消費端契約`（誰能 import、誰擁有哪張對照、優先序） | 章節以「本節描述`外部 repo`」開頭                   |
| 時間經過（已經發生的事）                          | `docs/CHANGELOG.md`                                                    | `已移除` / `已解體` / `已併回` / `不再` / `原 X 是` |
| 還沒做完                                          | `README.todo`                                                          | 「不另立 X」「暫不支援 Y」其實是待辦                |
| 一次 commit（行數、檔案數、byte 上限、module 數） | 刪除；量測值隨每次 commit 改變，不寫入文件                             | `333→101 行`、`共 10 個 module`、`上限 128 MiB`     |
| 程式碼改了，且可用指令驗證                        | 測試或 `scripts/`；文件只留一句規則 + 測試名稱                         | 文件裡整段可執行的 `grep` / `go list` 斷言          |
| 細節層級改動（CLI flag、依賴版本、安裝步驟）      | `docs/` 根層指南；正典檔留 quick start／一行指令 + 指標（見下節）      | `使用方式`／`開發指南`章節長成完整參考              |
| 換一台機器（絕對路徑）                            | 改相對路徑或 `$(git rev-parse --show-toplevel)`                        | `/Users/<name>/...`、`~/projects/<other-repo>/...`  |
| 都不會失效（不變式）                              | 留在 `CLAUDE.md`                                                       | 「只有 X 可以 import Y」                            |

- 結構樹、ownership、架構決策由 `CLAUDE.md` 單一擁有，`README.md` 用一行指過去。
- `README.md` 出現需要先讀原始碼才懂的型別名，就是越界 —— `刪除`，不下放：
  entry point 由 `CLAUDE.md` 模組對應表單一擁有，領域名稱即對照鍵。
- 清單另有 owner 檔時（如 `plugins/README.md` 擁有 plugin 清單），
  正典檔改`一行指標`，不複製清單。
- 寫得出 pass/fail 的斷言不放 Markdown ——`沒人執行的斷言會腐爛成錯誤資訊`。
  先把斷言自動化，才有資格刪文件。
- 稽核與瘦身用本技能的`範疇清理 (Scope Cleanup)` 模式，階段細節見
  [scope-cleanup.md](scope-cleanup.md)。

`消費端契約 (consumer-side contract)` 是唯一可以描述外部 repo 的內容：誰能 import 它、
誰擁有哪張對照表、優先序是什麼。它的對立面是`實作細節`（對方的檔案權限、預設 port、
內部路由表）—— 判準是`本 repo 能不能 build 或 test 它`，不能就不寫。

## 細節下放 (Detail Demotion)

`仍是現況`但改動頻率高於正典檔其餘內容的細節，不刪除而是`下放`到 `docs/` 根層指南：

| 來源章節                                                        | 目的地                | 正典檔保留                          |
| --------------------------------------------------------------- | --------------------- | ----------------------------------- |
| `README.md` `使用方式`的完整 CLI／API 參考                      | `docs/cli.md`         | 每領域 1-2 個 quick start 例 + 指標 |
| `CLAUDE.md` `開發指南`的展開細節（前置需求、clone、部署、排程） | `docs/development.md` | build／test／deploy 各一行指令 + 指標 |

- 判定門檻：章節超過約 `25 行`，或單一子命令／步驟的說明超過一個 code block。
- 內容`原文搬移`不重寫；遵循`先寫後刪`（scope-cleanup 的 Phase S2）。
- 目的地檔名固定為 `docs/cli.md` 與 `docs/development.md`，各 repo 一致，
  不另創同義檔名（`usage.md`、`setup.md`、`install.md`）。

## 檔案命名 (File Naming)

### `plans/` 與 `docs/specs/`

統一採用 `YYYY-MM-DD-<topic>.md` 格式：

- 日期：本地時區（Asia/Taipei），與既有 `docs/memory/` convention 一致
- `<topic>`：kebab-case 英文 topic name，必須 meaningful to the change（例：`windows-11-desktop-receiver`）
- 範例：✅ `2026-07-25-windows-11-desktop-receiver.md`；❌ `delegated-wishing-spark.md`
- 反例（system-generated slug，禁止使用）：`hashed-dancing-pascal.md`、`partitioned-rolling-gosling.md`

例外：plan-mode 系統啟動時自動產生的暫存檔（隨機 slug 命名）僅在 plan mode 內部使用，
`離開 plan mode 前必須` rename 為正式檔名，並同步更新所有 cross-reference
（README、CLAUDE.md、specs、plans）。
