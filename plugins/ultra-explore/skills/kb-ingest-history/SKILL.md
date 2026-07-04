---
name: kb-ingest-history
description: >
    Use when ingesting a git repository's development history into the
    knowledge base — caches git log to disk, keeps a last-commit cursor so
    later runs only process new commits, and uses the LLM to distill the
    history into development phases, key decisions, and refactor events as
    captures. Triggers on: "ingest git history", "understand the development
    history", "匯入開發歷史", "分析 commit 歷史".
version: "1.0.0"
allowed-tools: Read, Bash, Glob, Grep, Write, Edit
user-invocable: true
disable-model-invocation: true
effort: high
context: fork
---

# kb-ingest-history — 開發歷史入庫

把 git log 快取到磁碟、以 `last-commit` 游標增量更新，並用 LLM 把
commit 流蒸餾成「開發史」captures：階段、關鍵決策、重構事件。
儲存格式、指紋去重、狀態追蹤依 `kb-spec`，先讀它再開始。
與 `kb-ingest-repo` 的分工：repo 技能讀「現在的程式」，本技能讀「如何走到現在」。

## Checklist（複製到回覆中，逐項執行並打勾）

```md
- [ ] Step 0: 檢查 _state/runs/ — 續跑或新建 run
- [ ] Step 1: 讀游標 → git log 增量抓取 → 追加快取檔
- [ ] Step 2: 分批（每批 300 commits）→ 寫 manifest.json
- [ ] Step 3: 逐批 LLM 蒸餾 → 寫 captures → 更新 progress
- [ ] Step 4: 全批完成 → 前移 last-commit 游標 → status = done → 報告
```

## Step 0 — 續跑檢查 (Resume Check)

同 `kb-spec` 規則：有 `in-progress` 的 `kb-ingest-history` run 就續跑（跳到
Step 3），否則新建 `run-id = <date>-history-<name>`。

## Step 1 — 游標與快取 (Cursor & Cache)

游標存在 `_sources/history-<name>.md` 的 frontmatter；快取存
`_state/cache/history-<name>.log`（只追加，不重寫）。

1. 讀游標：

    ```bash
    grep '^last-commit:' <proj>/_sources/history-<name>.md
    ```

2. 決定範圍並抓取（`--first-parent` 走主線、`--reverse` 由舊到新）：

    ```bash
    # 首次（無游標）
    git -C <repo> log --first-parent --reverse --date=short \
      --pretty=format:'%h|%ad|%an|%s' > <kb>/_state/cache/history-<name>.log

    # 增量（有游標）：只抓新 commit，追加到快取
    git -C <repo> log <last-commit>..HEAD --first-parent --reverse \
      --date=short --pretty=format:'%h|%ad|%an|%s' \
      >> <kb>/_state/cache/history-<name>.log
    ```

3. 新 commit 數為 0 → 直接回報「已是最新（游標 <sha>）」並結束
4. 需要變更規模訊號時，對重點 commit 個別補 `git show --stat <sha>`，
   不要對全歷史跑 `--numstat`（大 repo 會爆量）

## Step 2 — 分批 (Plan)

- 只分批「本次新抓的」commits（首次 = 全部），每批 300 個
- 寫 `manifest.json`（每批的 commit sha 範圍）+ 初始 `progress.json`，
  `STATUS.md` 加一列

## Step 3 — 逐批 LLM 蒸餾 (Distill per Batch)

讀該批快取內容，蒸餾為開發史知識。只收四類，每條必附 commit sha 佐證：

| 類別                | 判定訊號                                         | truth 初判  |
| ------------------- | ------------------------------------------------ | ----------- |
| 階段 (phase)        | 連續多個 commit 圍繞同一主題（如「導入插件化」）| `firsthand` |
| 決策 (decision)     | 訊息明示取捨：migrate X→Y、adopt、deprecate      | `firsthand` |
| 重構 (refactor)     | restructure/rename/extract 叢集                  | `firsthand` |
| 風險訊號 (signal)   | revert/hotfix/rollback 叢集 — 指向脆弱區         | `candidate` |

硬規則：

- commit 訊息記載的事（誰、何時、改了什麼）→ `firsthand`；
  對「為什麼」的推測 → `candidate`，並明寫是推論
- 逐 commit 流水帳禁止 — 產出是「事件與階段」，一批通常收斂為 3~10 條
- 版面雜訊（chore、bump、format、merge commit 本身）不收
- 寫 captures（依 kb-spec 格式與指紋查重）：`sources.type: history`、
  `ref` 含 repo 與 sha 範圍、`zone-hint` 依主題
- 整批做完 → 更新 `progress.json` 與 `STATUS.md`，才開始下一批

## Step 4 — 前移游標與收尾 (Advance Cursor & Finish)

順序不可顛倒 — captures 全部落盤後才前移游標：

1. 更新 `_sources/history-<name>.md`：`last-commit: <本次最新 sha>`、
   `last-seen: <date>`、`## Captures` 追加
2. `progress.json` 標 `done`，`STATUS.md` 同步
3. 報告：新 commit 數、批次數、各類別條數、游標新位置、
   建議下一步（`kb-distill`）

之後每次重跑本技能即為增量更新：游標之前的歷史不重讀。

## Common Mistakes

| 錯誤                             | 修正                                       |
| -------------------------------- | ------------------------------------------ |
| 每次重抓全歷史                   | 讀游標，只抓 `<last-commit>..HEAD`         |
| captures 未落盤就前移游標        | 先 captures、後游標；中斷時游標仍指舊位置  |
| 逐 commit 流水帳當知識           | 收斂為階段/決策/重構/風險訊號              |
| 推測動機標 `firsthand`           | commit 記載才是 firsthand；動機推測 candidate |
| chore/bump/merge 雜訊入庫        | 版面雜訊不收                               |
| 對全歷史跑 --numstat             | 只對重點 commit 補 `git show --stat`       |
| 全部處理完才寫進度               | 每批處理完立即更新 progress.json           |
