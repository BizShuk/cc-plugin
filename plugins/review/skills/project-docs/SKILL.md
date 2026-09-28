---
name: project-docs
description: >
    Owns a project's canonical docs end to end — README.md (business domains),
    AGENTS.md (technical context, CLAUDE.md symlinks to it), docs/terminology.md (glossary),
    README.business.md (business value) — plus their content ownership and
    their history. Six modes: audit (verify every doc claim against the repo,
    reported as `doc says X → actually Y`), refresh, bootstrap, business,
    consolidate (fold `docs/specs/` and `plans/` older than two weeks into one
    summary table, move hand-written change logs into `docs/CHANGELOG.md`), and
    scope (strip from README/CLAUDE whatever another file owns, demote CLI and
    dev-setup detail to `docs/cli.md` / `docs/development.md`). Read-only by
    default. Triggers on: "explore project", "summarize codebase", "doc sync",
    "docs out of date", "does the README match", "extract business",
    "consolidate docs", "merge specs", "clean up plans", "scope cleanup",
    "文件同步", "更新文件了嗎", "術語表", "terminology", "業務萃取", "上下游分析",
    "文件整併", "合併規格", "整理 plans", "範疇清理", "文件瘦身", "下放 docs".
version: "6.0.0"
allowed-tools: Read, Bash, Glob, Grep, Write, Edit
user-invocable: true
disable-model-invocation: false
effort: high
context: fork
metadata:
    type: review
---

# project-docs

一個專案的正典文件由本技能統一負責：建立、稽核、更新、`內容歸屬`與`歷史壓縮`。
程式碼是真理來源；文件同步至程式碼，從不反向。

三個軸互補，對象都是同一組文件：`真實性`（文件說的還是不是真的）、
`歸屬`（這句話該不該在這個檔案裡）、`數量`（歷史文件與變更紀錄的累積）。

---

## 文件目標與格式 (Document Goals & Format)

目錄佈局與統一介面由 [unified-interface.md](references/unified-interface.md) 單一擁有。
各文件的模板與寫入規則詳見 [references/](references/)。
`哪句話該放哪個檔案`由 [content-ownership.md](references/content-ownership.md) 單一擁有
（含 `plans/`／`docs/specs/` 檔名規範與細節下放門檻）—— 寫或稽核文件前先讀它。

| 檔案 | 問題 | 目標 | 樣板 |
| ---- | ---- | ---- | ---- |
| `README.md` | `WHAT` | 業務領域、領域流程、實體、使用情境 | [readme.md](references/readme.md) |
| `AGENTS.md` | `HOW` | 專案結構、技術棧、模組對應、建置/部署、慣例 | [agents.md](references/agents.md) |
| `docs/terminology.md` | `WHICH` | 術語單一定義來源：領域名詞、縮寫、狀態值 | [docs-terminology.md](references/docs-terminology.md) |
| `README.business.md` | `WHY` | 業務價值：上下游、約束、風險、核心/非核心 | [readme-business.md](references/readme-business.md) |
| `CLAUDE.md` | — | symlink → `AGENTS.md`（必備） | [agents.md](references/agents.md) |
| `README.todo` | — | 待辦事項（必備，改善建議的唯一去處） | [readme-todo.md](references/readme-todo.md) |
| `docs/memory/` | — | 歷史決策（必備，僅回報缺漏） | — |
| `docs/tutorials/` | — | 領域知識導覽（選備，交由 `[[tutorial]]`） | [docs-structure.md](references/docs-structure.md) |
| `plans/` | — | 進行中計畫（選備，`consolidate` 的來源） | [content-ownership.md](references/content-ownership.md) |
| `docs/specs/` | — | 既有設計與規格（選備，`consolidate` 的來源） | [content-ownership.md](references/content-ownership.md) |
| `docs/CHANGELOG.md` | — | 已完成變更的唯一去處，只增不減（選備） | [consolidate.md](references/consolidate.md) |
| `docs/cli.md` `docs/development.md` | — | 細節下放的固定目的地（選備） | [content-ownership.md](references/content-ownership.md) |
| `docs/backlog/` | — | 待辦想法（選備，不整併） | — |
| `.geminiignore` | — | symlink → `.gitignore`（選備） | — |

`專案位址:` `~/projects/<project>/` 或 `~/projects/<category>/<project>/`。
分類目錄`不是專案`，不跑 `bootstrap`。歸屬確認用 `[[project-route]]`。

---

## 偵測一致性 (Consistency Detection)

文件在程式碼改動的當下就開始漂移。檢測兩個維度：`縱向`（文件 vs 程式碼）與
`橫向`（文件 vs 文件），逐條以 `doc says X → actually Y` / `A says X / B says Y` 回報。
檢查項清單與 audit 輸出格式見 [consistency.md](references/consistency.md)。

---

## 同步文件 (Syncing Documents)

檢測後根據模式決定寫入行為。程式碼永遠是真理來源；
橫向不一致以 `docs/terminology.md` 為裁決依據。

### 模式 (Modes)

由 Phase 0 自動判定；判定不明時預設 `audit`。

| 模式 | 觸發情境 | 寫入行為 |
| ---- | -------- | -------- |
| `audit` | 「doc sync」「文件同步」「does the README match」 | 無 — 只輸出一致性報告 |
| `refresh` | 「更新文件」「explore project」且文件已存在 | 僅改寫已證實不一致的段落 |
| `bootstrap` | `README.md` 或 `AGENTS.md` 缺漏、為空、僅剩標題 | 四份文件全產出 + symlinks |
| `business` | 「業務萃取」「extract business」「上下游分析」 | 僅 `README.business.md` |
| `consolidate` | 「文件整併」「整理 plans」「合併規格」 | 壓縮 `docs/specs/`／`plans/`，搬變更紀錄，`git rm` 來源 |
| `scope` | 「範疇清理」「文件瘦身」「這兩份文件不該有什麼」 | 刪除越界內容、下放高變動細節 |

`預設唯讀:` `audit` 是安全預設。其餘四個會寫檔的模式（`refresh`、`bootstrap`、
`consolidate`、`scope`）必須由使用者明確要求或由文件缺漏事實觸發；
`consolidate` 與 `scope` 另外會刪除內容，刪除前一律先列清單。

### 同步流程

#### Phase 0 — 模式判定

0. 確認目標是`專案根目錄`而非分類目錄。
1. 檢查 `README.md` 與 `AGENTS.md` 是否存在且非空（>10 行有效內容）。
2. 任一缺漏或為空 → `bootstrap`。
3. 兩者皆在 → 依觸發詞選 `audit`（預設）或 `refresh`。
4. 輸入不是完整 workspace，或請求純業務分析 → `business`。
5. 對象是`歷史文件的數量`（`docs/specs/`、`plans/`、變更紀錄章節）→ `consolidate`。
6. 對象是`正典文件的內容歸屬`（該不該寫在這裡）→ `scope`。

輸出判定結果：`Mode: audit（README.md 與 AGENTS.md 皆存在，未要求寫入）`。

#### Phase 1 — 掃描

`audit` 只需掃到足以驗證的程度；`refresh` / `bootstrap` 需完整掃描。

1. `Discover layout` — Glob 探索（排除 `.git`, `node_modules`, `vendor`, `dist`, `gen/` 等噪音）。
2. `Identify key files` — 讀取依賴檔、建置檔、entry points、現有文件。
3. `Read critical source` — skim 前 5-10 個高訊號原始檔。不逐檔閱讀。
4. `Identify business domains` — handler/service/module 分組成 3-7 個業務領域。

#### Phase 2 — 檢測一致性（`audit` 與 `refresh` 必做）

依上方「偵測一致性」章節執行縱向 + 橫向檢測。

`audit` 到此為止。要修，改跑 `refresh`：以程式碼為真理來源、`docs/terminology.md`
為用詞裁決依據更新文件，然後重跑 Phase 2 複驗。同步時不得夾帶範圍變更。

#### Phase 3 — 寫入（`refresh` / `bootstrap`）

`refresh` 只改寫 Phase 2 證實不一致的段落，保留其餘原文。
`bootstrap` 依樣板產出全文。寫入順序：

1. `docs/terminology.md` — 先建立用詞基準
2. `README.md` — 依術語表用詞撰寫
3. `AGENTS.md` — 模組對應與 README 領域對齊
4. `README.business.md` — 引用已確立的領域與術語
5. `README.todo` — 分析出的改善建議寫成待辦項，`不`寫進 `README.md`
6. `Symbolic links` — 執行 [setup-links.sh](scripts/setup-links.sh)：`CLAUDE.md` → `AGENTS.md`、
   `.geminiignore` → `.gitignore`；舊佈局（`CLAUDE.md` 為實體）就地反轉

各文件依 [references/](references/) 對應樣板產出。
連結已正確時跳過；兩邊皆為實體且無法安全反轉時 log `WARN` 並跳過。

#### Phase 4 — 報告

```text
✅ project-docs 完成 — Mode: <audit | refresh | bootstrap | business>

漂移 (Drift): <N> 項已修正 / <N> 項待處理
README.md: <line count> 行, <N> 個業務領域（改善建議 <N> 項 → README.todo）
AGENTS.md: <line count> 行, <N> 個核心模組
docs/terminology.md: <N> 筆術語, <N> 個狀態值, <N> 筆缺出處
README.business.md: <line count> 行, <N> 個業務約束, <N> 項風險

統一介面缺件 (Missing): <AGENTS.md | CLAUDE.md | README.todo | docs/memory/ | 無>

Symlinks:
- CLAUDE.md -> AGENTS.md ✅ (created | moved | already exists | skipped)
- .geminiignore -> .gitignore ✅ (created | already exists | skipped)

業務領域摘要:
- <Domain 1>: <1-sentence summary>
- <Domain 2>: <1-sentence summary>
```

`audit` 模式省略行數與 symlink 區塊，只輸出漂移清單。

---

### business 模式 — 純業務分析

輸入不是完整 workspace，或明確要求純業務分析時，跳過其餘文件與 symlinks，
僅產出 `README.business.md`。步驟見 [business-mode.md](references/business-mode.md)。

### consolidate 模式 — 歷史壓縮

把`兩週以前`的 `docs/specs/`／`plans/` 壓縮成每資料夾一份摘要表，把已完成的變更紀錄
（`AGENTS.md` 變更紀錄章節、`README.todo` 的 `## Archive`）搬到 `docs/CHANGELOG.md`。
`先寫後刪`：寫入失敗絕不進入刪除。範圍規則、各階段指令與樣板見
[consolidate.md](references/consolidate.md)。

### scope 模式 — 範疇清理

對象是 `README.md` 與 `AGENTS.md` `本身的內容`：判準是一個問題
「這句話會因為什麼而變成假的？」七種失效原因對應的歸屬、`消費端契約`的定義，
以及細節下放門檻，全部見 [content-ownership.md](references/content-ownership.md)。

| Phase | 動作 | 產物 |
| ----- | ---- | ---- |
| S0 `Audit` | 記錄 token 基線；逐段套歸屬判準，每筆標`失效原因`與`目的地` | 處置表（doc 說 X → 實際 Y） |
| S1 `Automate` | 可執行的斷言先落成測試／腳本，並`注入違規證明它會紅` | 測試檔／`scripts/` |
| S2 `Verify dst` | 確認要搬的內容`目的地已有`；已有則是`刪除`不是搬移 | 前提查核結果 |
| S3 `Cut` | 以 anchor 文字（非行號）逐段刪改 | 瘦身後的正典文件 |
| S4 `Sweep` | `重讀全文`找殘留：懸空引用、交叉連結、孤立表格列 | 殘留清單 |
| S5 `Lint` | 連結解析、機器路徑、外部細節、測試與腳本全綠 | 驗收輸出 |
| S6 `Report` | 逐類回報處置量、行數與 token 前後差異、未處理項 | 報告 |

`先自動化再刪除`：S1 未完成前不得進入 S3 —— 否則斷言會出現無人把關的空窗期。
各階段指令、規則、掃描樣板與常見錯誤見 [scope-cleanup.md](references/scope-cleanup.md)。

---

## 設計哲學 (Design Philosophy)

一個專案可能有幾十個 handler / service / module，但它們只屬於少數幾個
`業務領域 (Business Domains)`。README 應該以領域為單位組織，不是以檔案或
handler 為單位。

`專案位址:` 目標 repo 可能位於 `~/projects/<project>/`
或 `~/projects/<category>/<project>/`。分類目錄`本身不是專案`，
不要對它跑 `bootstrap`。要確認歸屬，先用 `[[project-route]]`。

---

## Rules

- 章節標題用繁體中文加英文括號；內文遵循輸入專案的原始語言慣例
- 寫入前先驗證：`refresh` 未經 Phase 2 佐證的段落不得改寫
- 四份文件用詞一律以 `docs/terminology.md` 為準；發現新名詞先入表再使用
- 八個業務章節缺一不可；查無資料的章節明寫「未偵測到 (Not detected)」
- 禁止技術實作細節進入 `README.business.md`
- 圖表一律 Mermaid；邊線文字必須雙引號包覆（`A -->|"文字"| B`）
- 狀態與名詞必須有程式或文件依據，禁止虛構
- 不使用粗體強調，改用 `backtick`
- 同步只對齊事實，範圍變更另案提出

`consolidate` 與 `scope` 的專屬規則分別見 [consolidate.md](references/consolidate.md)、
[scope-cleanup.md](references/scope-cleanup.md)；常見錯誤與失敗模式見
[pitfalls.md](references/pitfalls.md)。

## Related

- `[[project-route]]` 先解析路徑歸屬，再進入本技能
- `[[planner]]` 當真實目錄樹本身有問題、或程式碼對程式碼互相矛盾，而非文件寫錯；
  術語衝突則以 `docs/terminology.md` 為裁決依據
- `[[tutorial]]` 產出 `docs/tutorials/` 學習導向文件；其術語以 `docs/terminology.md` 為準
- `[[changelog]]` 從 git history `生成` CHANGELOG；`consolidate` 只`搬移`人工寫好的紀錄
- `[[universal-consolidate]]` 同類產物 N → 1 的通用凝聚算子，`consolidate` 是它在文件歷史上的具體化
- `ultra-explore` 處理跨 repo 知識庫；本技能綁定單一 repo 根目錄
