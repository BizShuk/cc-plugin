---
name: kb-ingest-repo
description: >
    Use when ingesting a git repository (local path or remote URL) into the
    knowledge base — scans entry points, domain models, and business rules in
    resumable batches, then writes raw captures and entity candidates into the
    KB inbox with per-batch progress persisted. Built for 1000+ file repos.
    Triggers on: "ingest this repo", "add repo to knowledge base",
    "把 repo 加入知識庫", "匯入程式庫知識".
version: "1.0.0"
allowed-tools: Read, Bash, Glob, Grep, Write, Edit
user-invocable: true
disable-model-invocation: true
effort: high
context: fork
---

# kb-ingest-repo — git repo 入庫

掃描 git repository，產出 raw captures 寫入知識庫 `_inbox/`。
本技能只負責「repo → captures」，不建 entity、不建邊。
儲存格式、指紋去重、狀態追蹤規則全部依 `kb-spec`，先讀它再開始。

## Checklist（複製到回覆中，逐項執行並打勾）

```md
- [ ] Step 0: 檢查 _state/runs/ — 續跑或新建 run
- [ ] Step 1: 取得 repo + 登記 _sources/
- [ ] Step 2: 列出高訊號檔案 → 寫 manifest.json（分批，每批 50 檔）
- [ ] Step 3: 逐批處理：讀檔 → 萃取候選事實 → 寫 captures → 更新 progress
- [ ] Step 4: 全部批次完成 → progress.status = done → 報告
```

## Step 0 — 續跑檢查 (Resume Check)

```bash
grep -rl '"skill": "kb-ingest-repo"' <kb>/_state/runs/*/manifest.json 2>/dev/null
```

- 找到同 `source` 且 `progress.json` 為 `in-progress` 的 run
  → 讀該 run 的 manifest 與 progress，從第一個未完成 batch 繼續，跳到 Step 3
- 沒有 → 建新 run：`run-id = <date>-repo-<name>`，繼續 Step 1

## Step 1 — 取得與登記 (Acquire)

- 本地路徑直接用；遠端 URL：`git clone --depth 50 <url> <scratchpad>/<name>`
- 建立或更新 `_sources/repo-<name>.md`（`type: repo`, `ref`, `last-seen`）

## Step 2 — 掃描與分批 (Plan)

1. 列出高訊號檔案（沿用 `project-docs` 排除清單：`.git`, `node_modules`,
   `vendor`, `dist`, `gen/`, 測試檔）：

    ```bash
    git -C <repo> ls-files | grep -vE '(_test\.|\.spec\.|^vendor/|^gen/|node_modules/|\.pb\.go$)' \
      | grep -E '\.(go|py|ts|js|java|rs|sql|proto|ya?ml|toml|md)$'
    ```

2. 按優先序排序：manifest 檔 → entry points（`main.*`, `cmd/`, route 註冊）→
   領域模型（`model/`, `entity/`, schema）→ 設定檔 → 其餘
3. 每 50 檔一批，寫入 `_state/runs/<run-id>/manifest.json`（格式見 kb-spec），
   同步建立 `progress.json`（`batches_done: 0`）並在 `STATUS.md` 加一列
4. 1000+ 檔的 repo：前 10 批（500 檔）通常已覆蓋高訊號區；之後的批次
   若連續 2 批未產生任何候選事實，可將剩餘批次標 `skipped-low-signal`
   並記入 `log.md`

## Step 3 — 逐批處理 (Process per Batch)

對每一批，依序做完以下四件事，才能開始下一批：

1. `讀檔`：Read 該批檔案（有 codegraph 時優先 `codegraph_explore`）
2. `萃取候選事實`，每條必附 `檔案:行` 佐證。只收四類：
    - 系統/模組目的與服務對象
    - 業務操作（業務動詞，非函數名）與業務約束（門檻、去重、retention）
    - 上下游依賴（datastore、MQ、外部 API）— 未來的 entity 候選
    - 狀態機（狀態名必須逐字來自程式）
3. `寫 captures`：每個主題一個 capture（依 kb-spec 格式與指紋查重）：
    - `sources.type: repo`；`truth`：程式碼直接可驗證 → `firsthand`，
      推論 → `candidate`
    - `zone-hint` 依業務領域給建議
    - 同主題已有 capture 時用 Edit 追加，不新建重複檔
4. `落盤進度`：更新 `progress.json`（batches_done、items_done、outputs）
   與 `STATUS.md`；跳過的檔在 `log.md` 寫原因

## Step 4 — 收尾 (Finish)

- `progress.json` 的 `status` 改 `done`，`STATUS.md` 同步
- 更新 `_sources/repo-<name>.md` 的 `## Captures` 清單
- 報告：批次數、capture 數、entity 候選清單（name/type/zone-hint）、
  跳過統計、建議下一步（`kb-distill`）

## Common Mistakes

| 錯誤                                    | 修正                                         |
| --------------------------------------- | -------------------------------------------- |
| 全部掃完才寫進度                        | 每批處理完立即更新 progress.json             |
| 中斷後從頭重掃                          | Step 0 先查 in-progress run，續跑            |
| 直接在 curated 區建 entity              | 本技能只寫 `_inbox/`；entity 由 kb-distill   |
| 推論寫成 `firsthand`                    | 程式碼可直接驗證才算；推論一律 `candidate`   |
| 把 logger/config/utils 列為 entity 候選 | 基礎設施雜訊排除（見 kb-spec）               |
| 候選事實不附佐證位置                    | 每條附 `檔案:行` 或設定鍵                    |

## Failure Modes

| 情境               | 動作                                                   |
| ------------------ | ------------------------------------------------------ |
| clone 失敗/私有庫  | 來源檔標 `status: unreachable`，run 標 `failed`，回報  |
| 單批處理中途出錯   | 該批不計入 batches_done，log.md 記錯誤，重試該批最多 2 次 |
| 既有同指紋 capture | 依 kb-spec：同來源跳過、異來源追加 sources             |
