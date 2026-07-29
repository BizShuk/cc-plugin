# docs-consolidation — 執行程序 (Procedure)

## Phase 0 — Preflight

```bash
# 專案根目錄與 git 狀態
git rev-parse --show-toplevel && git status --short

# 兩週門檻（macOS / Linux 擇一）
date -v-14d +%F 2>/dev/null || date -d '14 days ago' +%F
```

1. 目標必須是`專案根目錄`（有 `README.md` 或 `.git`），不是分類目錄。
   不確定歸屬時先跑 `[[project-route]]`。
2. `未提交變更 (uncommitted changes)` 存在於 `docs/specs/`、`plans/`、`CLAUDE.md`、
   `README.todo` 或 `docs/CHANGELOG.md` 時，停下來回報，請使用者先提交 —
   本技能會刪檔並就地改寫正典文件，git 是唯一的還原路徑。
3. 非 git repo 時：不刪除任何檔案，改為把來源檔案移到 `docs/specs/archive/`
   並在報告中說明。變更紀錄搬移（Step 4.4）`不受此限`——它是搬移不是刪除，
   內容完整落在 `docs/CHANGELOG.md`。

## Phase 1 — Inventory

### Step 1.1 — 文件來源

```bash
ls -1 docs/specs/*.md plans/*.md 2>/dev/null
```

對每個資料夾建立來源清單：

| 檔案 | 文件日期 | 日期來源 | 早於門檻？ | 是舊摘要？ |
| ---- | -------- | -------- | ---------- | ---------- |

`資料夾不存在`或`門檻外檔案 < 2 份`（且無舊摘要）時，跳過該資料夾並在報告中註明
`跳過 (skipped): 不足以整併`。整併一份文件沒有意義。

### Step 1.2 — 變更紀錄來源

```bash
# CLAUDE.md 的變更紀錄章節（章節層級，非整檔）
rg -n '^#{2,4} .*(變更紀錄|更新紀錄|更新歷史|Changelog|Change Log|History|Release Notes)' CLAUDE.md

# README.todo 的 Archive 章節
rg -n '^## Archive' -A 200 README.todo

# 既有 CHANGELOG（決定是新建還是合併）
ls docs/CHANGELOG.md 2>/dev/null
```

逐條建立條目清單（`一列一個條目`，不是一列一個章節）：

| 來源檔 | 條目原文 | 日期 | 日期來源 | 早於門檻？ |
| ------ | -------- | ---- | -------- | ---------- |

`條目日期`判定順序：

1. 條目行內或其所屬子標題的 `YYYY-MM-DD`
2. 條目連結到的文件檔名 `YYYY-MM-DD-` 前綴
3. 該行首次進入 git 的日期：
   `git log -S '<條目原文>' --diff-filter=A --format=%ad --date=short -- <file> | tail -1`
4. 以上皆無 → 判為 `未定 (undated)`

`未定`的條目`留在原檔不動`並在報告標註 —— 無法證明它早於門檻，就不搬。
`CLAUDE.md` 找不到變更紀錄章節，或 `README.todo` 無 `## Archive`／該章節無已勾選項目時，
跳過該來源並在報告註明 `跳過 (skipped): 無變更紀錄條目`。

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

`plans/<today>-Refresh.md` 結構相同，標題改為
`# 計畫摘要 (Plans Refresh) — <YYYY-MM-DD>`。

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
- `CLAUDE.md`、`README.todo`、`docs/CHANGELOG.md` 永遠不在刪除清單內

### Step 4.4 — 變更紀錄搬移 (Changelog Migration)

`排在 4.3 之後`：搬移時要知道哪些來源文件已被吸收，才能把失效連結改指向摘要檔。

`docs/CHANGELOG.md`（不存在則以此樣板建立，存在則`合併`）：

```markdown
# 變更紀錄 (Changelog)

由 `docs-consolidation` 自 `CLAUDE.md` 與 `README.todo` 搬移彙整，新的在上。

| 日期 (Date) | 變更 (Change) | 來源 (Source) | 原始文件 (Reference) |
| ----------- | ------------- | ------------- | -------------------- |
| 2026-07-06 | `Skill metadata` 標準化 | `README.todo` | [specs/2026-07-06-skill-metadata-standardization.md](specs/2026-07-06-skill-metadata-standardization.md) |
| 2026-07-02 | `Viper` 預設值單一來源 | `README.todo` | [specs/2026-07-29-Summary.md](specs/2026-07-29-Summary.md)（原 `2026-07-02-architecture-config-externalization.md`，已整併） |
```

寫入規則：

- `變更 (Change)` 沿用條目原文，只去掉 `- [x]` 勾選前綴與 markdown 連結語法，`不重寫措辭`
- `來源 (Source)` 只有兩種值：`` `CLAUDE.md` `` 或 `` `README.todo` ``
- `原始文件` 用`相對於 docs/CHANGELOG.md` 的路徑；條目沒有連結時寫「無」
- 連結指向的文件已被 Step 4.3 吸收時，改指向 `<today>-Summary.md`／`<today>-Refresh.md`，
  並在括號內保留原檔名
- 合併後全表依`日期遞減`排序；`日期相同`時 `CLAUDE.md` 的條目排在 `README.todo` 之前
- `去重 (dedupe)`：`日期 + 變更` 已存在於既有表格時整列跳過，本步驟必須可重複執行
- 既有列`不得改寫或重排欄位`；只允許插入新列

### Step 4.5 — 清空原章節

搬移`寫入成功之後`才動原檔，且`只動已搬走的條目`：

- `CLAUDE.md`：移除已搬走的條目。章節被搬空時`保留標題`，內容替換為一行指標：

  ```markdown
  ## 變更紀錄 (Changelog)

  見 [docs/CHANGELOG.md](docs/CHANGELOG.md)。
  ```

- `README.todo`：移除已搬走的 `- [x]` 項目。`## Archive` 標題`一律保留`
  （`README.todo` 格式規範要求該章節存在），搬空時底下留同一行指標。
- 未定日期與門檻內的條目`原地保留`，不得順手一起清掉。
- 這兩個檔案是`就地編輯 (in-place edit)`，用 `Edit` 而非重寫整檔。

## Phase 5 — Report

```text
✅ docs-consolidation 完成 — <YYYY-MM-DD>

門檻 (Cutoff): <today - 14d>，門檻內 <N> 份文件未動

docs/specs/: <N> 份 → 2026-07-22-Summary.md（<M> 列，<K> 淘汰，<J> 待確認）
plans/:      <N> 份 → 2026-07-22-Refresh.md（<M> 列，<K> 淘汰）

變更紀錄 (Changelog) → docs/CHANGELOG.md（<新建 | 合併>）:
- CLAUDE.md:   搬出 <N> 條，章節 <保留指標 | 未發現變更紀錄章節>
- README.todo: 搬出 <N> 條已完成項目，Archive 保留 <M> 條
- 去重跳過 (skipped, 已存在): <N> 條
- 未定日期未搬 (undated): <逐條列出>

已刪除 (Removed):
- docs/specs/<file>.md
- docs/specs/<prev>-Summary.md（舊摘要，已吸收）

README.md 側記: 新增 <K> 筆淘汰記錄
待確認 (⚠️): <列出每筆及查無實證的理由>
```
