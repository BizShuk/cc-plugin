---
name: project-docs
description: >
    Establish, audit, and update a project's canonical docs — README.md
    (functional requirements), CLAUDE.md (technical context),
    docs/terminology.md (glossary), and README.business.md (business
    extraction). Audits first: every doc claim is verified against the repo and
    reported as `doc says X → actually Y`; writing happens only when the docs
    are missing or the user asks for an update. Use when onboarding to a new
    project, after major refactors, when README.md and CLAUDE.md are missing or
    outdated, or when extracting business value from any input (folder, repo,
    file, document, pasted text). Triggers on: "explore project", "summarize
    codebase", "doc sync", "docs out of date", "does the README match", "check
    the structure tree", "文件同步", "更新文件了嗎", "extract business",
    "business value", "業務萃取", "上下游分析", "術語表", "terminology",
    "glossary".
version: "4.0.0"
allowed-tools: Read, Bash, Glob, Grep, Write, Edit
user-invocable: true
disable-model-invocation: false
effort: high
context: fork
metadata:
    type: review
---

# project-docs

一個專案的正典文件由本技能統一負責建立、稽核與更新。
程式碼是真理來源；文件同步至程式碼，從不反向。

---

## 文件目標與格式 (Document Goals & Format)

規範來源：`cc-plugin/config/CLAUDE.global.md` 統一介面。
各文件的模板與寫入規則詳見 [references/](references/)。

| 檔案 | 問題 | 目標 | 樣板 |
| ---- | ---- | ---- | ---- |
| `README.md` | `WHAT` | 業務領域、領域流程、實體、使用情境 | [readme.md](references/readme.md) |
| `CLAUDE.md` | `HOW` | 專案結構、技術棧、模組對應、建置/部署、慣例 | [claude.md](references/claude.md) |
| `docs/terminology.md` | `WHICH` | 術語單一定義來源：領域名詞、縮寫、狀態值 | [docs-terminology.md](references/docs-terminology.md) |
| `README.business.md` | `WHY` | 業務價值：上下游、約束、風險、核心/非核心 | [readme-business.md](references/readme-business.md) |
| `AGENTS.md` | — | symlink → `CLAUDE.md`（必備） | — |
| `README.todo` | — | 待辦事項（必備，僅回報缺漏） | [readme-todo.md](references/readme-todo.md) |
| `docs/memory/` | — | 歷史決策（必備，僅回報缺漏） | — |
| `docs/tutorials/` | — | 領域知識導覽（選備，交由 `[[tutorial]]`） | [docs-structure.md](references/docs-structure.md) |
| `plans/` | — | 進行中計畫，`YYYY-MM-DD-<topic>.md`（選備） | — |
| `docs/specs/` | — | 既有設計與規格，`YYYY-MM-DD-<topic>.md`（選備） | — |
| `docs/backlog/` | — | 待辦想法（選備） | — |
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
- CLAUDE.md tree: omits plugins/god/ and plugins/superpower/ (both exist)
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

`預設唯讀:` `audit` 是安全預設。`refresh` 與 `bootstrap`
會寫檔，必須由使用者明確要求或由文件缺漏事實觸發。

### 同步流程

#### Phase 0 — 模式判定

0. 確認目標是`專案根目錄`而非分類目錄。
1. 檢查 `README.md` 與 `CLAUDE.md` 是否存在且非空（>10 行有效內容）。
2. 任一缺漏或為空 → `bootstrap`。
3. 兩者皆在 → 依觸發詞選 `audit`（預設）或 `refresh`。
4. 輸入不是完整 workspace，或請求純業務分析 → `business`。

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
5. `Symbolic links` — `AGENTS.md` → `CLAUDE.md`；`.geminiignore` → `.gitignore`

各文件依 [references/](references/) 對應樣板產出。
若連結已存在或目標是普通檔案（log `WARN`）就跳過。

#### Phase 4 — 報告

```text
✅ project-docs 完成 — Mode: <audit | refresh | bootstrap | business>

漂移 (Drift): <N> 項已修正 / <N> 項待處理
README.md: <line count> 行, <N> 個業務領域, <N> 項改善建議
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
| 對分類目錄跑 `bootstrap` | 分類不是專案；逐一處理其下的專案 |

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

## Related

- `[[project-route]]` 先解析路徑歸屬，再進入本技能
- `[[system-planner]]` 當真實目錄樹本身有問題，而非文件寫錯
- `[[consistency]]` 處理程式碼對程式碼的矛盾；術語衝突以 `docs/terminology.md` 為裁決依據
- `[[tutorial]]` 產出 `docs/tutorials/` 學習導向文件；其術語以 `docs/terminology.md` 為準
- `ultra-explore` 處理跨 repo 知識庫；本技能綁定單一 repo 根目錄
