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
