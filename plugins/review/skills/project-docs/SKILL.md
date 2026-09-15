---
name: project-docs
description: >
    Owns a project's canonical docs end to end — README.md (business domains),
    CLAUDE.md (technical context), docs/terminology.md (glossary),
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
version: "5.0.0"
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

規範來源：`cc-plugin/config/CLAUDE.global.md` 統一介面。
各文件的模板與寫入規則詳見 [references/](references/)。
`哪句話該放哪個檔案`由 [content-ownership.md](references/content-ownership.md) 單一擁有
（含 `plans/`／`docs/specs/` 檔名規範與細節下放門檻）—— 寫或稽核文件前先讀它。

| 檔案 | 問題 | 目標 | 樣板 |
| ---- | ---- | ---- | ---- |
| `README.md` | `WHAT` | 業務領域、領域流程、實體、使用情境 | [readme.md](references/readme.md) |
| `CLAUDE.md` | `HOW` | 專案結構、技術棧、模組對應、建置/部署、慣例 | [claude.md](references/claude.md) |
| `docs/terminology.md` | `WHICH` | 術語單一定義來源：領域名詞、縮寫、狀態值 | [docs-terminology.md](references/docs-terminology.md) |
| `README.business.md` | `WHY` | 業務價值：上下游、約束、風險、核心/非核心 | [readme-business.md](references/readme-business.md) |
| `AGENTS.md` | — | symlink → `CLAUDE.md`（必備） | — |
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

文件在程式碼改動的當下就開始漂移。本技能檢測兩個維度的一致性：

- `縱向:` 每份文件 vs 程式碼（文件宣稱的事實是否仍為真）
- `橫向:` 文件 vs 文件（同一概念在不同文件中是否一致）

### 檢測流程

1. 從所有文件抽出可驗證的宣稱 — 路徑、指令、模組對應、術語、狀態值、領域名稱。
2. `縱向驗證:` 逐條對 repo 驗證，以 `doc says X → actually Y` 回報。
3. `橫向驗證:` 交叉比對文件間的引用，以 `A says X / B says Y` 回報。

### 縱向：文件 vs 程式碼

| 文件 | 檢查項 |
| ---- | ------ |
| `README.md` | 業務領域仍對應到真實 handler/module；CLI/API 指令仍可執行 |
| `CLAUDE.md` | 目錄樹 diff 真實目錄；模組對應 entry point 仍存在；build/deploy 指令有效；設定路徑一致 |
| `docs/terminology.md` | 術語出處路徑仍存在；狀態字面值與程式碼 enum/const 一致 |
| `README.business.md` | 狀態名稱能在程式中找到；上游服務在程式碼中有對應呼叫 |
| 統一介面 | 必備檔案存在：`AGENTS.md`、`README.todo`、`docs/memory/`、`docs/terminology.md` |

### 橫向：文件 vs 文件

| 比對對 | 檢查項 |
| ------ | ------ |
| `README.md` ↔ `CLAUDE.md` | README 的業務領域必須全數出現在 CLAUDE 的模組對應表；反之模組對應不得列出 README 未定義的領域 |
| `README.md` ↔ `docs/terminology.md` | README 中出現的領域名詞必須在術語表有定義；術語表的領域分節必須與 README 業務領域對齊 |
| `CLAUDE.md` ↔ `docs/terminology.md` | CLAUDE 使用的技術術語（若為領域詞）必須與術語表一致；不得出現同義詞漂移 |
| `README.business.md` ↔ `README.md` | business 的業務目的/常見操作必須與 README 業務領域對應；不得出現 README 未提及的領域 |
| `README.business.md` ↔ `docs/terminology.md` | 狀態值名稱必須與術語表的狀態值章節一致 |
| 所有文件 ↔ `docs/terminology.md` | 同一概念在所有文件中只能使用術語表定義的正名，不得有第二種說法 |

### audit 輸出格式

```text
Doc consistency review — <scope>

[縱向] 文件 vs 程式碼:
- CLAUDE.md tree: omits plugins/god/ and plugins/team/ (both exist)
- README.md: "uses MySQL" → code uses SQLite (model/store.go)
- README.business.md: state "archived" not found in code

[橫向] 文件 vs 文件:
- README ↔ CLAUDE: README 列「資料匯出」領域，CLAUDE 模組對應表無此領域
- README ↔ terminology: README 用「蒸餾管道」，CLAUDE 用「distiller pipeline」— 術語表未定義
- business ↔ terminology: business 狀態 "pending" vs 術語表 "waiting" — 應統一

[統一介面] 缺件:
- docs/terminology.md: missing (必備)
```

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
| `bootstrap` | `README.md` 或 `CLAUDE.md` 缺漏、為空、僅剩標題 | 四份文件全產出 + symlinks |
| `business` | 「業務萃取」「extract business」「上下游分析」 | 僅 `README.business.md` |
| `consolidate` | 「文件整併」「整理 plans」「合併規格」 | 壓縮 `docs/specs/`／`plans/`，搬變更紀錄，`git rm` 來源 |
| `scope` | 「範疇清理」「文件瘦身」「這兩份文件不該有什麼」 | 刪除越界內容、下放高變動細節 |

`預設唯讀:` `audit` 是安全預設。其餘四個會寫檔的模式（`refresh`、`bootstrap`、
`consolidate`、`scope`）必須由使用者明確要求或由文件缺漏事實觸發；
`consolidate` 與 `scope` 另外會刪除內容，刪除前一律先列清單。

### 同步流程

#### Phase 0 — 模式判定

0. 確認目標是`專案根目錄`而非分類目錄。
1. 檢查 `README.md` 與 `CLAUDE.md` 是否存在且非空（>10 行有效內容）。
2. 任一缺漏或為空 → `bootstrap`。
3. 兩者皆在 → 依觸發詞選 `audit`（預設）或 `refresh`。
4. 輸入不是完整 workspace，或請求純業務分析 → `business`。
5. 對象是`歷史文件的數量`（`docs/specs/`、`plans/`、變更紀錄章節）→ `consolidate`。
6. 對象是`正典文件的內容歸屬`（該不該寫在這裡）→ `scope`。

輸出判定結果：`Mode: audit（README.md 與 CLAUDE.md 皆存在，未要求寫入）`。

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
3. `CLAUDE.md` — 模組對應與 README 領域對齊
4. `README.business.md` — 引用已確立的領域與術語
5. `README.todo` — 分析出的改善建議寫成待辦項，`不`寫進 `README.md`
6. `Symbolic links` — `AGENTS.md` → `CLAUDE.md`；`.geminiignore` → `.gitignore`

各文件依 [references/](references/) 對應樣板產出。
若連結已存在或目標是普通檔案（log `WARN`）就跳過。

#### Phase 4 — 報告

```text
✅ project-docs 完成 — Mode: <audit | refresh | bootstrap | business>

漂移 (Drift): <N> 項已修正 / <N> 項待處理
README.md: <line count> 行, <N> 個業務領域（改善建議 <N> 項 → README.todo）
CLAUDE.md: <line count> 行, <N> 個核心模組
docs/terminology.md: <N> 筆術語, <N> 個狀態值, <N> 筆缺出處
README.business.md: <line count> 行, <N> 個業務約束, <N> 項風險

統一介面缺件 (Missing): <AGENTS.md | README.todo | docs/memory/ | 無>

Symlinks:
- AGENTS.md -> CLAUDE.md ✅ (created | already exists | skipped)
- .geminiignore -> .gitignore ✅ (created | already exists | skipped)

業務領域摘要:
- <Domain 1>: <1-sentence summary>
- <Domain 2>: <1-sentence summary>
```

`audit` 模式省略行數與 symlink 區塊，只輸出漂移清單。

---

### business 模式 — 純業務分析

輸入不是完整 workspace，或明確要求純業務分析時，
跳過 `README.md` / `CLAUDE.md` / symlinks，僅產出 `README.business.md`。

1. `Scope 界定` — folder/repo 用 Glob 鎖定 entry points；單一檔案直接讀取；
   純文字直接分析。
2. `Core Business & Operations` — 找出系統的存在理由、列出常見業務操作
   （業務動詞，不是函數名清單）。
3. `其餘章節` — 依 [readme-business.md](references/readme-business.md) 八章節完成。
4. `Write Report` — folder/repo 同步寫入兩個位置：
   - `<target>/README.business.md`
   - `~/projects/product/projects/<name>/README.business.md`
   既有檔案先讀取後合併，保留仍正確的內容。

---

### consolidate 模式 — 歷史壓縮

`docs/specs/` 與 `plans/` 會累積成一堆單次性文件，`CLAUDE.md` 的變更紀錄章節與
`README.todo` 的 `## Archive` 只增不減。本模式把`兩週以前`的歷史壓縮成`每個資料夾一份`
的摘要表，把已完成的變更紀錄搬到 `docs/CHANGELOG.md`，讓正典文件只描述`現況`。

| 來源 (Source) | 產出 (Output) | 舊產物 (Previous) |
| ------------- | ------------- | ----------------- |
| `docs/specs/*.md` | `docs/specs/<YYYY-MM-DD>-Summary.md` | 被吸收後刪除 |
| `plans/*.md` | `plans/<YYYY-MM-DD>-Refresh.md` | 被吸收後刪除 |
| `CLAUDE.md` 變更紀錄章節 | `docs/CHANGELOG.md` | 搬移後從原檔移除 |
| `README.todo` 的 `## Archive` | `docs/CHANGELOG.md` | 搬移後從原檔移除 |
| 已淘汰功能 | `README.md` 的側記章節 | 累加，不刪 |

`<YYYY-MM-DD>` 是`執行當天`的日期，不是來源文件的日期。
`docs/CHANGELOG.md` 只增不減，不隨執行日期換檔名。

範圍規則：

- `兩週門檻:` 只處理文件日期`早於` `today - 14d` 的檔案；門檻內的`原封不動`。
- 文件日期判定順序：檔名 `YYYY-MM-DD-` 前綴 → frontmatter `date:` →
  `git log --diff-filter=A` 首次提交 → mtime。用到後兩者時在報告標註。
- 上一份 `-Summary.md`／`-Refresh.md` `一律納入`來源，不受門檻限制。
- `CLAUDE.md` 與 `README.todo` 是`部分來源`：只搬出變更紀錄章節與 `## Archive` 的
  已勾選項目，檔案本身`絕不刪除`。`關鍵決策`、`技術棧`、`模組對應`、`慣例`
  描述`現況`，不是變更紀錄，一律不動。
- 不適用：`docs/backlog/`（未實作，無存在性可驗證）、`docs/memory/`（保存本身就是目的）、
  從 git history `生成`變更紀錄（那是 `[[changelog]]`；本模式只`搬移已經寫好`的）。

| Phase | 動作 | 產物 |
| ----- | ---- | ---- |
| 0 `Preflight` | 確認專案根目錄、git 乾淨、算出兩週門檻 | 門檻日期 |
| 1 `Inventory` | 列出文件來源與變更紀錄條目，判定各自日期 | 來源清單 + 條目清單 |
| 2 `Extract` | 逐份讀取，抽出日期／功能／使用方式／價值 | 四欄表格列 |
| 3 `Verify` | 驗證每列的功能是否仍存在於 workspace | live / deprecated / unknown |
| 4 `Write` | 寫摘要檔 → 補 README 側記 → `git rm` 來源 → 搬變更紀錄 → 清空原章節 | 摘要檔 + `docs/CHANGELOG.md` |
| 5 `Report` | 回報門檻、列數、刪除清單、搬移條數與待確認項 | 報告 |

`先寫後刪`：寫入失敗時絕不進入刪除步驟；`docs/CHANGELOG.md` 寫入成功之後
才清空 `CLAUDE.md` 與 `README.todo` 的原章節。
各階段指令、樣板與常見錯誤見 [consolidate.md](references/consolidate.md)。

---

### scope 模式 — 範疇清理

對象是 `README.md` 與 `CLAUDE.md` `本身的內容`：判準是一個問題
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
各階段指令、掃描樣板與常見錯誤見 [scope-cleanup.md](references/scope-cleanup.md)。

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

`consolidate` 另加：

- 摘要表`固定四欄`：日期、功能、使用方式、價值 — 不增欄不改序；
  `使用方式`必須是可執行的指令或明確觸發詞，不得寫「見原文件」
- 兩週門檻內的文件與條目`一律不動`，包含不得讀進摘要表
- 存在性`不確定`時一律保留並標 `⚠️ 待確認`，不得逕行刪除
- 一律 `git rm` 逐檔刪除，不用 `rm`、不用萬用字元
- 變更紀錄條目`原文照搬`，以 `日期 + 變更` 去重；日期`未定`的留在原檔

`scope` 另加：

- 稽核先於改寫：每筆先寫成 `doc 說 X → 實際 Y`，經確認才動手
- 刪除前必須`實測目的地已有該內容`，否則是遺失不是搬移
- 文件裡的斷言`預設為錯`，一律先實跑；guard test 必須`注入違規證明它會紅`
- 編輯用 anchor 文字定位，`不用行號`；刪完`重讀全文`，不只看改動處
- 驗證工具必須唯讀；不預告最終行數

## Common Mistakes

| 錯誤 | 修正 |
| ---- | ---- |
| 沒驗證就整份覆寫既有 README | 先跑 Phase 2，只改證實漂移的段落 |
| 把環境初始化、設定同步當成業務領域 | 歸入非核心或直接排除 |
| 只描述流程、不畫狀態機 | 業務物件有狀態欄位就必須有 stateDiagram |
| 略過風險章節因為「看起來沒風險」 | 逐類別回報「無」也是結論 |
| 全部列為核心業務 | 強制二分，非核心需說明如何支撐核心 |
| 用函數名稱清單冒充業務操作 | 改寫為業務動詞 + 觸發者 + 結果 |
| 術語表塞入通用技術詞 (HTTP/JSON) | 只收領域名詞與專案自訂縮寫 |
| 把改善建議寫進 `README.md` | 待辦屬 `README.todo`，正典文件只描述現況 |
| README 列出 handler／型別名 | entry point 由 `CLAUDE.md` 模組對應表單一擁有 |
| 對分類目錄跑 `bootstrap` | 分類不是專案；逐一處理其下的專案 |
| 把 `關鍵決策` 當成變更紀錄搬走 | 它描述`現況`，是 `CLAUDE.md` 存在的理由 |
| 查不到實證就判定淘汰 | 無刪除紀錄一律標 `⚠️ 待確認`並保留 |
| 先刪文件再補測試 | 反序會有無人把關的空窗期；S1 未完成不得進 S3 |
| 沒查目的地就刪 | 先實測目的地已有該內容，否則是遺失不是搬移 |

## Failure Modes

| 情境 | 動作 |
| ---- | ---- |
| Workspace is empty | 寫最小 stub，註明「空專案 (Empty project)」 |
| Cannot detect language/framework | 在對應章節註明「未偵測到 (Not detected)」 |
| 既有 README/CLAUDE 有價值內容 | 合併 — 經 Phase 2 驗證後保留有效章節，只更新漂移部分 |
| 太多檔案無法全掃 | 聚焦頂層 + entry points，註明「僅掃描部分檔案 (Partial scan)」 |
| 找不到明確狀態機 | 改用 flowchart 描述業務流程並註明 |
| 業務邊界不明 | 依目錄/模組分組並註明「邊界不明確」 |
| 文件宣稱的指令無法執行 | 回報為漂移，不要自行發明替代指令 |
| 找不到足夠術語建表 | 寫最小表格並註明「待補 (To be extended)」 |
| 目標路徑是分類目錄 | 停止並列出其下專案，請使用者指定 |
| 非 git repo（`consolidate`） | 不刪檔，改 `mv` 到 `archive/` 並在報告註明 |
| 來源資料夾門檻外不足 2 份 | 回報「不足以整併」並跳過，不產生摘要檔 |
| `CLAUDE.md` 無變更紀錄章節 | 跳過該來源，`不得`自行從 git history 生成 |

## Related

- `[[project-route]]` 先解析路徑歸屬，再進入本技能
- `[[planner]]` 當真實目錄樹本身有問題、或程式碼對程式碼互相矛盾，而非文件寫錯；
  術語衝突則以 `docs/terminology.md` 為裁決依據
- `[[tutorial]]` 產出 `docs/tutorials/` 學習導向文件；其術語以 `docs/terminology.md` 為準
- `[[changelog]]` 從 git history `生成` CHANGELOG；`consolidate` 只`搬移`人工寫好的紀錄
- `[[universal-consolidate]]` 同類產物 N → 1 的通用凝聚算子，`consolidate` 是它在文件歷史上的具體化
- `ultra-explore` 處理跨 repo 知識庫；本技能綁定單一 repo 根目錄
