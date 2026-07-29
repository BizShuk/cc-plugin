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
    Triggers on: "consolidate docs", "merge specs", "clean up plans",
    "move changelog", "archive todo", "文件整併", "合併規格", "整理 plans",
    "搬移變更紀錄", "整理已完成待辦", "docs consolidation".
version: "1.1.0"
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
把正典文件 (canonical docs) 越撐越長。本技能把`兩週以前`的歷史壓縮成`每個資料夾一份`的摘要表，
只保留仍然存在於 workspace 的功能，其餘轉為 `README.md` 的淘汰註記；
並把已完成的變更紀錄搬到 `docs/CHANGELOG.md`，讓正典文件只描述`現況`。

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

## 執行程序 (Procedure)

各階段的指令、判定表與輸出樣板見
[references/procedure.md](references/procedure.md)。

| Phase | 動作 | 產物 |
| ----- | ---- | ---- |
| 0 `Preflight` | 確認專案根目錄、git 乾淨、算出兩週門檻 | 門檻日期 |
| 1 `Inventory` | 列出文件來源與變更紀錄條目，判定各自日期 | 來源清單 + 條目清單 |
| 2 `Extract` | 逐份讀取，抽出日期／功能／使用方式／價值 | 四欄表格列 |
| 3 `Verify` | 驗證每列的功能是否仍存在於 workspace | live / deprecated / unknown |
| 4 `Write` | 寫摘要檔 → 補 README 側記 → `git rm` 來源 → 搬變更紀錄 → 清空原章節 | 摘要檔 + 側記 + `docs/CHANGELOG.md` |
| 5 `Report` | 回報門檻、列數、刪除清單、搬移條數與待確認項 | 報告 |

`先寫後刪`：Phase 4 的寫入失敗時，絕不進入刪除步驟；
變更紀錄同理，`docs/CHANGELOG.md` 寫入成功之後才清空 `CLAUDE.md` 與 `README.todo` 的原章節。

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

常見錯誤與失效情境處置見
[references/troubleshooting.md](references/troubleshooting.md)。

## Related

- `[[project-docs]]` 正典文件（README / CLAUDE.md）的建立與稽核；本技能只碰歷史文件、README 側記與變更紀錄章節
- `[[changelog]]` 從 git history `生成` CHANGELOG；本技能只`搬移`人工寫好的紀錄，兩者互補不重疊
- `[[project-route]]` 先確認目標是專案根目錄而非分類目錄
- `[[universal-consolidate]]` 同類產物 N → 1 的通用凝聚算子，本技能是它在文件歷史上的具體化
- `[[tutorial]]` 學習導向文件，不在本技能的整併範圍
