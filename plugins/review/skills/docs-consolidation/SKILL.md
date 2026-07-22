---
name: docs-consolidation
description: >
    Collapse a project's accumulated design and planning history into one file
    per folder — `docs/specs/*.md` into `docs/specs/<YYYY-MM-DD>-Summary.md`
    and `plans/*.md` into `plans/<YYYY-MM-DD>-Refresh.md` — as a
    date / feature / how-to-use / value table. Only documents older than two
    weeks are consolidated; features no longer present in the workspace are
    dropped from the table and recorded as a side note in `README.md`. The
    previous consolidated file is absorbed and removed when the next one is
    generated. Use when specs or plans have piled up, after shipping a batch of
    features, or before onboarding someone to a long-lived repo. Triggers on:
    "consolidate docs", "merge specs", "clean up plans", "文件整併",
    "合併規格", "整理 plans", "docs consolidation".
version: "1.0.0"
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
沒有人回頭讀。本技能把`兩週以前`的歷史壓縮成`每個資料夾一份`的摘要表，
只保留仍然存在於 workspace 的功能，其餘轉為 `README.md` 的淘汰註記。

| 來源 (Source)   | 產出 (Output)                         | 舊產物 (Previous)      |
| --------------- | ------------------------------------- | ---------------------- |
| `docs/specs/*.md` | `docs/specs/<YYYY-MM-DD>-Summary.md` | 被吸收後刪除           |
| `plans/*.md`      | `plans/<YYYY-MM-DD>-Refresh.md`      | 被吸收後刪除           |
| 已淘汰功能        | `README.md` 的側記章節                | 累加，不刪             |

`<YYYY-MM-DD>` 是`執行當天`的日期，不是來源文件的日期。

## When to Use

- `docs/specs/` 或 `plans/` 超過 5 份文件，且多數已實作完畢
- 交付一批功能之後，想把設計歷史收斂成一份可讀的清單
- 接手長期 repo，需要一張「有什麼功能、怎麼用、值多少」的表
- 定期維護（如每季）壓縮文件目錄

不適用：`docs/backlog/`（尚未實作的想法，沒有「是否仍存在」可驗證）、
`docs/memory/`（歷史決策本身就是保存目的）、`docs/tutorials/`（教學文件用
`[[tutorial]]`）、CHANGELOG 生成（用 `ultra-explore` 的 `changelog`）。

## Scope Rules

- `兩週門檻 (Two-week cutoff)`：只處理文件日期`早於` `today - 14d` 的檔案。
  門檻內的文件`原封不動`，它們還在活躍使用中。
- `文件日期 (Document date)` 判定順序：檔名 `YYYY-MM-DD-` 前綴 → frontmatter
  `date:` → `git log --diff-filter=A --format=%ad --date=short -- <file>`（首次提交）
  → 檔案 mtime。用到後兩者時在報告中標註。
- 上一份 `-Summary.md` / `-Refresh.md` `一律納入`來源（不受兩週門檻限制），
  其表格內容合併進新檔後刪除舊檔。
- `README.md`、`README.todo`、資料夾內的 `README.md` 不是來源，不刪。

---

## Phase 0 — Preflight

```bash
# 專案根目錄與 git 狀態
git rev-parse --show-toplevel && git status --short

# 兩週門檻（macOS / Linux 擇一）
date -v-14d +%F 2>/dev/null || date -d '14 days ago' +%F
```

1. 目標必須是`專案根目錄`（有 `README.md` 或 `.git`），不是分類目錄。
   不確定歸屬時先跑 `[[project-route]]`。
2. `未提交變更 (uncommitted changes)` 存在於 `docs/specs/` 或 `plans/` 時，
   停下來回報，請使用者先提交 — 本技能會刪檔，git 是唯一的還原路徑。
3. 非 git repo 時：不刪除任何檔案，改為把來源檔案移到 `docs/specs/archive/`
   並在報告中說明。

## Phase 1 — Inventory

```bash
ls -1 docs/specs/*.md plans/*.md 2>/dev/null
```

對每個資料夾建立來源清單：

| 檔案 | 文件日期 | 日期來源 | 早於門檻？ | 是舊摘要？ |
| ---- | -------- | -------- | ---------- | ---------- |

`資料夾不存在`或`門檻外檔案 < 2 份`（且無舊摘要）時，跳過該資料夾並在報告中註明
`跳過 (skipped): 不足以整併`。整併一份文件沒有意義。

## Phase 2 — Extract

逐份 `Read` 來源文件，抽出四個欄位。不要憑檔名猜測，內容不足時明寫「未記載」。

| 欄位 | 抽取來源 | 規則 |
| ---- | -------- | ---- |
| `日期 (Date)` | Phase 1 判定的文件日期 | `YYYY-MM-DD`，一律用來源文件的日期，不是今天 |
| `功能 (Feature)` | 標題、目標章節 | 一個名詞短語 + `backtick` 標出實際識別符（指令、技能名、套件路徑） |
| `使用方式 (How to Use)` | 使用/介面/CLI 章節 | 可直接執行的指令或觸發詞；抽象描述不算 |
| `價值 (Value)` | 動機/問題/背景章節 | 一句話：解決什麼問題、省下什麼成本 |

舊摘要檔的表格`直接沿用其列`，不要重新推導 — 它的來源文件已經不存在了。

## Phase 3 — Verify Existence

每一列都必須驗證`功能是否仍存在於當前 workspace`。這是本技能唯一會刪除資訊的判斷，
不得憑印象。

```bash
# 依功能型態擇一驗證
ls <claimed-path>                          # 檔案/目錄宣稱
rg -n '<identifier>' --glob '!docs/**'     # 函數、指令、設定鍵
git log --oneline -1 -- <path>             # 最後異動
```

| 判定 | 條件 | 處置 |
| ---- | ---- | ---- |
| `存在 (live)` | 宣稱的路徑或識別符至少有一項在程式碼/設定中找得到 | 列入摘要表 |
| `淘汰 (deprecated)` | 全部宣稱都查無此物，且 git log 顯示曾被刪除 | 移出表格，寫入 `README.md` 側記 |
| `不確定 (unknown)` | 查無實證但也沒有刪除紀錄 | `保留在表格`並標註 `⚠️ 待確認`，不得逕行刪除 |

`不確定一律保留`。誤刪一筆真實功能的成本，遠高於多留一列待確認。

## Phase 4 — Write

### Step 4.1 — 摘要檔

`docs/specs/<today>-Summary.md`：

```markdown
# 規格摘要 (Specs Summary) — <YYYY-MM-DD>

整併自 `docs/specs/` 中 <N> 份文件（<最早日期> ~ <最晚日期>），
涵蓋 <today - 14d> 之前的設計歷史。

| 日期 (Date) | 功能 (Feature) | 使用方式 (How to Use) | 價值 (Value) |
| ----------- | -------------- | --------------------- | ------------ |
| 2026-05-25  | `markdownlint` 技能 | `/markdownlint <path>` | Markdown 規則子集 + 無粗體約束，避免逐檔手動檢查 |

## 已淘汰 (Deprecated)

<本次驗證為淘汰的項目，一行一筆；無則寫「無」>

## 來源檔案 (Source Files)

<刪除前的檔名清單，供 git history 回溯>
```

`plans/<today>-Refresh.md` 結構相同，標題改為 `# 計畫摘要 (Plans Refresh) — <YYYY-MM-DD>`。

### Step 4.2 — `README.md` 側記

在專案根目錄 `README.md` 末尾維護一個章節（不存在則建立，存在則`累加`）：

```markdown
## 已淘汰功能 (Deprecated Features)

| 淘汰日期 | 功能 | 原始文件 | 說明 |
| -------- | ---- | -------- | ---- |
| 2026-07-22 | `media` plugin | `2026-06-17-media-plugin.md` | 已從 `plugins/` 移除，無替代 |
```

`淘汰日期`是本次執行日期。既有列不得改寫或刪除。

### Step 4.3 — 刪除來源

寫入成功`之後`才刪除，且刪除前把完整清單印給使用者確認：

```bash
git rm docs/specs/2026-05-14-feature-agent-design.md ...   # 逐檔列出，不用萬用字元
```

- 一律用 `git rm`，不用 `rm` — 保留 history 作為唯一還原路徑
- 只刪 Phase 1 清單中`早於門檻`的檔案與`舊摘要檔`
- 本次剛產生的摘要檔絕不在刪除清單內
- 非 git repo：改用 `mkdir -p docs/specs/archive && mv <files> docs/specs/archive/`

## Phase 5 — Report

```text
✅ docs-consolidation 完成 — <YYYY-MM-DD>

門檻 (Cutoff): <today - 14d>，門檻內 <N> 份文件未動

docs/specs/: <N> 份 → 2026-07-22-Summary.md（<M> 列，<K> 淘汰，<J> 待確認）
plans/:      <N> 份 → 2026-07-22-Refresh.md（<M> 列，<K> 淘汰）

已刪除 (Removed):
- docs/specs/<file>.md
- docs/specs/<prev>-Summary.md（舊摘要，已吸收）

README.md 側記: 新增 <K> 筆淘汰記錄
待確認 (⚠️): <列出每筆及查無實證的理由>
```

---

## Rules

- 章節標題用繁體中文加英文括號；表格內容跟隨來源文件的原始語言
- 表格`固定四欄`：日期、功能、使用方式、價值 — 不增欄不改序
- `使用方式`必須是可執行的指令或明確觸發詞，不得寫「見原文件」
- 兩週門檻內的文件`一律不動`，包含不得讀進摘要表
- 先寫後刪；寫入失敗時不得進入 Phase 4.3
- 不使用粗體強調，改用 `backtick`
- 摘要表不得虛構功能；來源文件沒寫的欄位寫「未記載」

## Common Mistakes

| 錯誤 | 修正 |
| ---- | ---- |
| 用今天的日期填表格的`日期`欄 | 日期欄是來源文件的日期，檔名的日期才是今天 |
| 舊摘要檔被當成一般來源重新解析 | 舊摘要的列直接沿用，其來源已不存在 |
| 查不到實證就判定淘汰 | 無刪除紀錄一律標 `⚠️ 待確認`並保留 |
| 用 `rm` 或萬用字元刪檔 | 逐檔 `git rm`，history 是唯一還原路徑 |
| 把 `docs/backlog/` 一起整併 | backlog 是未實作想法，無存在性可驗證 |
| 連兩週內的新文件一起吸收 | 門檻是硬規則，新文件還在活躍使用 |
| 改寫 `README.md` 既有的淘汰記錄 | 側記只累加 |
| 只有 1 份文件也照跑一次整併 | 回報「不足以整併」並跳過 |

## Failure Modes

| 情境 | 動作 |
| ---- | ---- |
| 非 git repo | 不刪檔，改 `mv` 到 `archive/` 並在報告註明 |
| `docs/specs/` 或 `plans/` 有未提交變更 | 停止，請使用者先提交 |
| 來源文件無日期可判定 | 用 mtime 並在報告標註 `日期來源: mtime` |
| 來源文件內容過短無法抽四欄 | 缺的欄位寫「未記載」，不猜 |
| 全部項目都判定為淘汰 | 仍產生摘要檔（只有淘汰章節），並在報告醒目提示 |
| 舊摘要檔格式不符（非四欄表格） | 整份內容以 `## 舊摘要 (Legacy)` 原文附在新檔末尾，不強轉 |
| 資料夾不存在 | 跳過，不建立空資料夾 |

## Related

- `[[project-docs]]` 正典文件（README / CLAUDE.md）的建立與稽核；本技能只碰歷史文件與 README 側記
- `[[project-route]]` 先確認目標是專案根目錄而非分類目錄
- `[[universal-consolidate]]` 同類產物 N → 1 的通用凝聚算子，本技能是它在文件歷史上的具體化
- `[[tutorial]]` 學習導向文件，不在本技能的整併範圍
