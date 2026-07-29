---
name: docs-consolidation
description: >
    Collapse a project's accumulated design and planning history into one file
    per folder — `docs/specs/*.md` into `docs/specs/<YYYY-MM-DD>-Summary.md`
    and `plans/*.md` into `plans/<YYYY-MM-DD>-Refresh.md` — as a
    date / feature / how-to-use / value table, and migrate the hand-written
    change log in `CLAUDE.md` plus completed `## Archive` items in
    `README.todo` into an append-only `docs/CHANGELOG.md`. Only documents and
    entries older than two weeks are consolidated; features no longer present
    in the workspace are dropped from the table and recorded as a side note in
    `README.md`. The previous consolidated file is absorbed and removed when
    the next one is generated. Use when specs or plans have piled up, when
    `CLAUDE.md` or `README.todo` has grown a long history tail, after shipping
    a batch of features, or before onboarding someone to a long-lived repo.
    Also runs a `scope cleanup` mode that removes from `README.md` and
    `CLAUDE.md` what does not belong there — another repo's internals,
    history already recorded in `docs/CHANGELOG.md`, measured numbers that
    rot on the next commit, duplicated trees and sections, machine-specific
    absolute paths, and executable assertions that no one runs — relocating
    each to its rightful owner or automating it into a test or script.
    Triggers on: "consolidate docs", "merge specs", "clean up plans",
    "move changelog", "archive todo", "文件整併", "合併規格", "整理 plans",
    "搬移變更紀錄", "整理已完成待辦", "docs consolidation", "scope cleanup",
    "doc scope", "文件瘦身", "範疇清理", "這兩份文件不該有什麼",
    "what should not be in README".
version: "1.2.0"
allowed-tools: Read, Write, Edit, Bash, Glob, Grep
user-invocable: true
disable-model-invocation: false
effort: high
context: fork
metadata:
    type: review
    platforms: [macos, linux]
---

# docs-consolidation

`docs/specs/` 與 `plans/` 會隨時間累積成一堆單次性文件，每份都只記錄「當時想做什麼」，
沒有人回頭讀；`CLAUDE.md` 的變更紀錄章節與 `README.todo` 的 `## Archive` 同樣只增不減，
把正典文件（canonical docs）越撐越長。本技能把`兩週以前`的歷史壓縮成`每個資料夾一份`的摘要表，
只保留仍然存在於 workspace 的功能，其餘轉為 `README.md` 的淘汰註記；
並把已完成的變更紀錄搬到 `docs/CHANGELOG.md`，讓正典文件只描述`現況`。

本技能有`兩個模式`，可分別執行也可接續執行：

| 模式                       | 對象                                  | 問題                       |
| -------------------------- | ------------------------------------- | -------------------------- |
| `整併 (Consolidate)`：預設 | `docs/specs/`、`plans/`、變更紀錄章節 | 歷史文件太多，壓縮成摘要表 |
| `範疇清理 (Scope Cleanup)` | `README.md`、`CLAUDE.md` 本身         | 正典文件裡有`不該在那裡`的內容 |

兩者互補：`整併`處理`歷史文件的數量`，`範疇清理`處理`正典文件的內容歸屬`。
使用者問「這兩份文件不該有什麼」、「文件瘦身」、「範疇清理」時走後者。

| 來源 (Source)                 | 產出 (Output)                        | 舊產物 (Previous) |
| ----------------------------- | ------------------------------------ | ----------------- |
| `docs/specs/*.md`             | `docs/specs/<YYYY-MM-DD>-Summary.md` | 被吸收後刪除      |
| `plans/*.md`                  | `plans/<YYYY-MM-DD>-Refresh.md`      | 被吸收後刪除      |
| `CLAUDE.md` 變更紀錄章節      | `docs/CHANGELOG.md`                  | 搬移後從原檔移除  |
| `README.todo` 的 `## Archive` | `docs/CHANGELOG.md`                  | 搬移後從原檔移除  |
| 已淘汰功能                    | `README.md` 的側記章節               | 累加，不刪        |

`<YYYY-MM-DD>` 是`執行當天`的日期，不是來源文件的日期。
`docs/CHANGELOG.md` 是`唯一且累加 (append-only)` 的檔案，不隨執行日期換檔名。

## When to Use

- `docs/specs/` 或 `plans/` 超過 5 份文件，且多數已實作完畢
- 交付一批功能之後，想把設計歷史收斂成一份可讀的清單
- 接手長期 repo，需要一張「有什麼功能、怎麼用、值多少」的表
- `CLAUDE.md` 長出一段越寫越長的變更紀錄，或 `README.todo` 的 `## Archive` 佔滿整個檔案
- 定期維護（如每季）壓縮文件目錄

`範疇清理`模式另適用於：

- 使用者問「這兩份文件不該有什麼」、要求文件瘦身或範疇稽核
- `README.md` 與 `CLAUDE.md` 出現重複章節或兩份會分岔的結構樹
- 正典文件裡有描述`外部 repo` 的章節（本 repo 無法 build 或 test 它）
- 文件裡有整段可執行的驗證指令，但沒有任何 CI 或測試在跑
- 大型重構之後，正典文件仍在描述`曾經如何`

不適用：`docs/backlog/`（尚未實作的想法，沒有「是否仍存在」可驗證）、
`docs/memory/`（歷史決策本身就是保存目的）、`docs/tutorials/`（教學文件用
`[[tutorial]]`）、從 git history `生成`變更紀錄（用 `ultra-explore` 的 `changelog`；
本技能只`搬移已經寫好`的紀錄，不從 commit 推導）。

## Scope Rules

- `兩週門檻 (Two-week cutoff)`：只處理文件日期`早於` `today - 14d` 的檔案。
  門檻內的文件`原封不動`，它們還在活躍使用中。
- `文件日期 (Document date)` 判定順序：檔名 `YYYY-MM-DD-` 前綴 → frontmatter
  `date:` → `git log --diff-filter=A --format=%ad --date=short -- <file>`（首次提交）
  → 檔案 mtime。用到後兩者時在報告中標註。
- 上一份 `-Summary.md` / `-Refresh.md` `一律納入`來源（不受兩週門檻限制），
  其表格內容合併進新檔後刪除舊檔。
- `README.md`、資料夾內的 `README.md` 不是來源，不刪。
- `CLAUDE.md` 與 `README.todo` 是`部分來源`：只搬出指定章節的`條目`，
  檔案本身`絕不刪除`，章節外的內容一個字都不動。
    - `CLAUDE.md` 只取變更紀錄類章節：標題含 `變更紀錄` / `更新紀錄` / `更新歷史` /
      `Changelog` / `Change Log` / `History` / `Release Notes`。
    - `關鍵決策 (Key Decisions)`、`技術棧`、`模組對應`、`慣例` 等`描述現況`的章節
      `不是`變更紀錄，一律不動 — 它們是 `CLAUDE.md` 存在的理由。
    - `README.todo` 只取 `## Archive` 章節內已勾選（`- [x]`）的項目；
      未勾選項目與其他章節不動。
- `docs/CHANGELOG.md` 只增不減：既有內容不得改寫、重排或刪除，新條目`合併`進去。

## 範疇判準 (Scope Ownership)

`範疇清理`模式的唯一判準是一個問題：**「這句話會因為什麼而變成假的？」**

| 失效原因                                          | 歸屬                                                        | 典型徵狀                                             |
| ------------------------------------------------- | ----------------------------------------------------------- | ---------------------------------------------------- |
| 別的 repo 改動                                    | 那個 repo；只留`消費端契約`                                 | 章節以「本節描述`外部 repo`」開頭                    |
| 時間經過（已經發生的事）                          | `docs/CHANGELOG.md`                                         | `已移除` / `已解體` / `已併回` / `不再` / `原 X 是`  |
| 還沒做完                                          | `README.todo`                                               | 「不另立 X」「暫不支援 Y」其實是待辦                 |
| 一次 commit（行數、檔案數、byte 上限、module 數） | 刪除                                                        | `333→101 行`、`共 10 個 module`、`上限 128 MiB`      |
| 程式碼改了，且可用指令驗證                        | 測試或 `scripts/`                                           | 文件裡整段可執行的 `grep` / `go list` 斷言           |
| 換一台機器                                        | 相對路徑                                                    | `/Users/<name>/...`、`~/projects/<other-repo>/...`   |
| 都不會失效（不變式）                              | 留在 `CLAUDE.md`                                            | 「只有 X 可以 import Y」                             |

`README.md` 是`為什麼用它、怎麼開始`；`CLAUDE.md` 是`邊界是什麼、誰擁有什麼`。
一個事實只能有`一個` owner —— 兩份重複的結構樹必然分岔，且`兩份都會不準`。

`消費端契約 (consumer-side contract)` 是唯一可以描述外部 repo 的內容：誰能 import 它、
誰擁有哪張對照表、優先序是什麼。它的對立面是`實作細節`（對方的檔案權限、預設 port、
內部路由表）—— 判準是`本 repo 能不能 build 或 test 它`，不能就不寫。

## 執行程序 (Procedure)

各階段的指令、判定表與輸出樣板見
[references/procedure.md](references/procedure.md)；
`範疇清理`模式見 [references/scope-cleanup.md](references/scope-cleanup.md)。

| Phase         | 動作                                                                | 產物                                |
| ------------- | ------------------------------------------------------------------- | ----------------------------------- |
| 0 `Preflight` | 確認專案根目錄、git 乾淨、算出兩週門檻                              | 門檻日期                            |
| 1 `Inventory` | 列出文件來源與變更紀錄條目，判定各自日期                            | 來源清單 + 條目清單                 |
| 2 `Extract`   | 逐份讀取，抽出日期／功能／使用方式／價值                            | 四欄表格列                          |
| 3 `Verify`    | 驗證每列的功能是否仍存在於 workspace                                | live / deprecated / unknown         |
| 4 `Write`     | 寫摘要檔 → 補 README 側記 → `git rm` 來源 → 搬變更紀錄 → 清空原章節 | 摘要檔 + 側記 + `docs/CHANGELOG.md` |
| 5 `Report`    | 回報門檻、列數、刪除清單、搬移條數與待確認項                        | 報告                                |

`先寫後刪`：Phase 4 的寫入失敗時，絕不進入刪除步驟；
變更紀錄同理，`docs/CHANGELOG.md` 寫入成功之後才清空 `CLAUDE.md` 與 `README.todo` 的原章節。

`範疇清理`模式的階段（細節見 `references/scope-cleanup.md`）：

| Phase           | 動作                                                       | 產物                     |
| --------------- | ---------------------------------------------------------- | ------------------------ |
| S0 `Audit`      | 逐段套範疇判準，每筆標`失效原因`與`目的地`                 | 處置表（doc 說 X → 實際 Y） |
| S1 `Automate`   | 可執行的斷言先落成測試／腳本，並`注入違規證明它會紅`       | 測試檔／`scripts/`       |
| S2 `Verify dst` | 確認要搬的內容`目的地已有`；已有則是`刪除`不是搬移         | 前提查核結果             |
| S3 `Cut`        | 以 anchor 文字（非行號）逐段刪改                           | 瘦身後的正典文件         |
| S4 `Sweep`      | `重讀全文`找殘留：懸空引用、被刪章節的交叉連結、孤立表格列 | 殘留清單                 |
| S5 `Lint`       | 連結解析、機器路徑、外部細節、測試與腳本全綠               | 驗收輸出                 |

`先自動化再刪除`：S1 未完成前不得進入 S3 —— 否則斷言會出現無人把關的空窗期。

## Rules

- 章節標題用繁體中文加英文括號；表格內容跟隨來源文件的原始語言
- 表格`固定四欄`：日期、功能、使用方式、價值 — 不增欄不改序
- `使用方式`必須是可執行的指令或明確觸發詞，不得寫「見原文件」
- 兩週門檻內的文件與變更紀錄條目`一律不動`，包含不得讀進摘要表
- 存在性`不確定`時一律保留並標 `⚠️ 待確認`，不得逕行刪除
- 一律 `git rm` 逐檔刪除，不用 `rm`、不用萬用字元
- 不使用粗體強調，改用 `backtick`
- 摘要表不得虛構功能；來源文件沒寫的欄位寫「未記載」
- `docs/CHANGELOG.md` 只增不減，且必須可重複執行 —— 同 `日期 + 變更` 一律去重
- 變更紀錄條目`原文照搬`，不重寫措辭、不合併相似項
- `CLAUDE.md` 與 `README.todo` 只做`就地移除已搬走的條目`，檔案本身不刪、其他章節不動
- 日期`未定`的條目留在原檔，不搬進 `docs/CHANGELOG.md`

`範疇清理`模式另加：

- 稽核先於改寫：每筆先寫成 `doc 說 X → 實際 Y`，經確認才動手
- 刪除前必須`實測目的地已有該內容`，否則是遺失不是搬移
- 文件裡的斷言`預設為錯`，一律先實跑；把錯的斷言原樣搬進測試等於固化錯誤
- guard test 必須`注入違規證明它會紅`，天生全綠的守護測試等於沒寫
- 編輯用 anchor 文字定位，`不用行號` —— 前面的刪除會讓後面的行號全部失效
- 刪完`重讀全文`，不只看改動處：被刪章節的交叉引用會變成懸空
- 驗證工具必須唯讀；會寫檔的建置指令要導向暫存目錄
- 不預告最終行數 —— 刪完才知道，硬湊數字只能靠刪真契約

常見錯誤與失效情境處置見
[references/troubleshooting.md](references/troubleshooting.md)。

## Related

- `[[project-docs]]` 正典文件（README / CLAUDE.md）的建立與稽核；本技能只碰歷史文件、README 側記與變更紀錄章節
- `[[changelog]]` 從 git history `生成` CHANGELOG；本技能只`搬移`人工寫好的紀錄，兩者互補不重疊
- `[[project-route]]` 先確認目標是專案根目錄而非分類目錄
- `[[universal-consolidate]]` 同類產物 N → 1 的通用凝聚算子，本技能是它在文件歷史上的具體化
- `[[tutorial]]` 學習導向文件，不在本技能的整併範圍
