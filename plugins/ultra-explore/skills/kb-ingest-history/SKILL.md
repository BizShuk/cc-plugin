---
name: kb-ingest-history
description: >
    Use when ingesting a git repository's development history into the
    knowledge base — runs the bundled kb_history.py deterministic pipeline
    (commits jsonl, stats, weekly diffs; stdlib-only, no install) as evidence
    material, keeps a last-commit cursor so later runs only process new
    commits, distills the history into development phases, key decisions,
    and refactor events as captures, and fills the CHANGELOG.md weekly
    narratives as a by-product. Triggers on: "ingest git history",
    "understand the development history", "匯入開發歷史", "分析 commit 歷史",
    "generate changelog into kb".
version: "1.1.0"
allowed-tools: Read, Bash, Glob, Grep, Write, Edit
user-invocable: true
disable-model-invocation: true
effort: high
context: fork
---

# kb-ingest-history — 開發歷史入庫

以本技能內建的 `scripts/kb_history.py` 確定性管道產出素材
（commit 清單、統計、週分桶 diff），
以 `last-commit` 游標增量蒸餾成「開發史」captures：階段、關鍵決策、重構事件；
蒸餾過程順手回填 `CHANGELOG.md` 的每週敘事。
儲存格式、指紋去重、狀態追蹤依 `kb-spec`，先讀它再開始。
與 `kb-ingest-repo` 的分工：repo 技能讀「現在的程式」，本技能讀「如何走到現在」。

## Checklist（複製到回覆中，逐項執行並打勾）

```md
- [ ] Step 0: 檢查 _state/runs/ — 續跑或新建 run
- [ ] Step 1: 跑內建 kb_history.py 管道 → <proj>/ 素材落盤；讀游標定增量範圍
- [ ] Step 2: 新 commit 依 ISO 週分批（每批 ≤ 300 commits）→ 寫 manifest.json
- [ ] Step 3: 逐批蒸餾（commit 訊息 + 週 diff 佐證）→ 寫 captures
      + 回填該批各週的 CHANGELOG 敘事 → 更新 progress
- [ ] Step 4: 全批完成 → 前移 last-commit 游標 → status = done → 報告
```

## Step 0 — 續跑檢查 (Resume Check)

同 `kb-spec` 規則：有 `in-progress` 的 `kb-ingest-history` run 就續跑（跳到
Step 3），否則新建 `run-id = <date>-history-<name>`。

## Step 1 — 素材管道與游標 (Material Pipeline & Cursor)

### 1a — 跑確定性管道（本技能內建腳本，零安裝）

腳本在本技能目錄 `scripts/kb_history.py`（單檔、純 stdlib，
只需 `python3` 與 `git`）：

```bash
python3 "${CLAUDE_PLUGIN_ROOT}/skills/kb-ingest-history/scripts/kb_history.py" \
  run <repo>                     # 預設輸出 ~/projects/product/projects/<repo名>/
```

`CLAUDE_PLUGIN_ROOT` 未定義時，改用本 SKILL.md 所在目錄下的
`scripts/kb_history.py` 絕對路徑。`python3` 不可用 → 走「降級模式」（見下）。

產出（全部落在 `<proj>/`，是工具產物、不是 entity）：

| 產物                | 內容                                   | 在本技能的用途          |
| ------------------- | -------------------------------------- | ----------------------- |
| `_raw/commits.jsonl`| hash/date/author/subject + numstat     | commit 清單 + 規模訊號  |
| `stats.json`        | 週分桶、作者統計、top commits          | 分批依據、階段訊號      |
| `_diffs/<week>.diff`| 每 ISO 週過濾後的實際 diff             | 蒸餾時的第一手佐證      |
| `CHANGELOG.md`      | 骨架 + `<!-- LLM: ... -->` 佔位符      | Step 3 回填敘事         |

過濾規則：內建通用雜訊排除（vendor/generated/lockfile/test/資料檔/二進位）
+ 各 repo 自己的 `.gitignore` overlay（`git check-ignore --stdin` 批次判定）。
管道是冪等全量重跑（便宜、確定性）；增量控制交給游標 — LLM 只蒸餾新 commit。

降級模式（無法安裝 CLI 時）：
`git -C <repo> log --first-parent --reverse --date=short --pretty=format:'%h|%ad|%an|%s'`
追加到 `<kb>/_state/cache/history-<name>.log`，無週 diff 佐證，
蒸餾只憑 commit 訊息（truth 上限見 Step 3），不回填 CHANGELOG。

### 1b — 讀游標定範圍

```bash
grep '^last-commit:' <proj>/_sources/history-<name>.md
```

- 無游標（首次）→ 範圍 = 全部 commits
- 有游標 → 範圍 = `commits.jsonl` 中游標 hash 之後的 commits
- 新 commit 數為 0 → 回報「已是最新（游標 <sha>）」並結束

## Step 2 — 分批 (Plan)

- 只分批「範圍內」的 commits，以 ISO 週為邊界湊批：連續數週合為一批，
  每批 ≤ 300 commits（週界對齊才能整週讀 diff、整週回填敘事）
- 寫 `manifest.json`（每批的週清單 + commit hash 範圍）+ 初始
  `progress.json`，`STATUS.md` 加一列

## Step 3 — 逐批蒸餾 + 回填敘事 (Distill & Narrate per Batch)

對每一批依序完成三件事：

### 3a — 讀素材

該批各週：`_raw/commits.jsonl` 的 commit 訊息 + `_diffs/<week>.diff`。
單週 diff 過大（>100K tokens）→ 只讀 numstat 與檔名清單，
敘事前加 `<!-- warning: large diff -->`（沿用 changelog 技能規則）。

### 3b — 蒸餾 captures

只收四類，每條必附 commit sha 佐證：

| 類別                | 判定訊號                                         | truth 初判  |
| ------------------- | ------------------------------------------------ | ----------- |
| 階段 (phase)        | 連續多個 commit 圍繞同一主題（如「導入插件化」）| `firsthand` |
| 決策 (decision)     | 訊息明示取捨：migrate X→Y、adopt、deprecate      | `firsthand` |
| 重構 (refactor)     | restructure/rename/extract 叢集                  | `firsthand` |
| 風險訊號 (signal)   | revert/hotfix/rollback 叢集 — 指向脆弱區         | `candidate` |

truth 判定（diff 佐證是核心增益）：

- 「改了什麼」被週 diff 內容直接印證 → `firsthand`（可驗證事實），
  佐證欄寫 `sha + _diffs/<week>.diff`
- 訊息宣稱但 diff 看不到對應變更 → 降為 `candidate` 並註明落差
  （訊息與實作不符本身就是風險訊號）
- 「為什麼」的動機推測 → 一律 `candidate`；diff 與訊息同屬一個 repo，
  不構成獨立來源，禁止以此升 `corroborated`
- 逐 commit 流水帳禁止 — 產出是「事件與階段」，一批通常收斂為 3~10 條
- 版面雜訊（chore、bump、format、merge commit 本身）不收
- captures 依 kb-spec 格式與指紋查重：`sources.type: history`、
  `ref` 含 repo 與 sha 範圍、`zone-hint` 依主題

### 3c — 回填 CHANGELOG 敘事（`changelog` 技能的 Phase 2）

該批每個 ISO 週：把 `<proj>/CHANGELOG.md` 對應週的 `<!-- LLM: ... -->`
佔位符替換為 3~5 句過去式敘事（讀者視角的功能變化，非實作細節）。
規則：只動佔位符，不改 commit 清單與統計表；
該週 diff 為空 → 寫 `_No business-logic changes this week._`

### 3d — 落盤

更新 `progress.json`（captures 路徑 + 已回填的週）與 `STATUS.md`，
才開始下一批。

## Step 4 — 前移游標與收尾 (Advance Cursor & Finish)

順序不可顛倒 — captures 與敘事全部落盤後才前移游標：

1. 更新 `<proj>/_sources/history-<name>.md`：`last-commit: <本次最新 sha>`、
   `last-seen: <date>`、`## Captures` 追加
2. `progress.json` 標 `done`，`STATUS.md` 同步
3. 報告：新 commit 數、批次數、各類別條數、回填週數與剩餘佔位符數、
   訊息與 diff 不符清單、游標新位置、建議下一步（`kb-distill`）

之後每次重跑本技能即為增量更新：素材管道全量重跑（冪等），
LLM 蒸餾與敘事只處理游標之後的新 commit。

## Common Mistakes

| 錯誤                               | 修正                                          |
| ---------------------------------- | --------------------------------------------- |
| 每次全歷史重新蒸餾                 | 素材可全量重跑，LLM 只處理 `<last-commit>` 之後 |
| captures 未落盤就前移游標          | 先 captures、後游標；中斷時游標仍指舊位置     |
| 逐 commit 流水帳當知識             | 收斂為階段/決策/重構/風險訊號                 |
| diff 印證就升 `corroborated`       | 同 repo 非獨立來源；印證只到 `firsthand`      |
| 動機推測標 `firsthand`             | 「為什麼」一律 `candidate`                    |
| 訊息與 diff 不符仍照訊息入庫       | 降 `candidate` + 註明落差 + 列入報告          |
| 回填敘事時改動統計表或 commit 清單 | 只替換 `<!-- LLM: ... -->` 佔位符             |
| chore/bump/merge 雜訊入庫          | 版面雜訊不收                                  |
| 把 CHANGELOG.md 當 entity 建檔     | 工具產物非 entity（見 kb-spec）               |
| 全部處理完才寫進度                 | 每批處理完立即更新 progress.json              |
