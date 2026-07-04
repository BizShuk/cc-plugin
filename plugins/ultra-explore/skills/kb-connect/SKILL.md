---
name: kb-connect
description: >
    Use when building edges between knowledge base entities — adds evidence-
    backed wikilink References between dimensions, rebuilds Backlinks from
    forward edges, and regenerates the _index.md registry with a Mermaid
    overview. Processes entities in resumable batches with progress persisted.
    Triggers on: "connect the knowledge base", "build kb edges", "rebuild kb
    index", "知識庫建邊", "重建知識庫索引", "kb connect".
version: "1.0.0"
allowed-tools: Read, Bash, Glob, Grep, Write, Edit
user-invocable: true
disable-model-invocation: true
effort: high
context: fork
---

# kb-connect — 建邊與索引

在 curated entities 之間建立有佐證的邊，重算 Backlinks，重建 `_index.md`。
邊規則完全沿用 `kb-spec`（方向、relation 動詞、edge grounding、雜訊過濾），
先讀它再開始。前置條件：`kb-distill` 已把維度落檔 — 邊只能指向磁碟上
既存的標題。

## Checklist（複製到回覆中，逐項執行並打勾）

```md
- [ ] Step 0: 檢查 _state/runs/ — 續跑或新建 run
- [ ] Step 1: 盤點 entities → 寫 manifest.json（每批 20 個 entity）
- [ ] Step 2: 逐批建邊：讀目標標題 → 寫 References → 更新 progress
- [ ] Step 3: 全圖重算 Backlinks（唯一全域步驟）
- [ ] Step 4: 重建 _index.md → status = done → 報告
```

## Step 0 — 續跑檢查 (Resume Check)

同 `kb-spec` 規則：有 `in-progress` 的 `kb-connect` run 就續跑，
否則新建 `run-id = <date>-connect-<slug>`。
注意：續跑只適用 Step 2；Step 3/4 是全域重算，必須在所有批次完成後執行。

## Step 1 — 盤點與分批 (Plan)

```bash
find <proj> -name '*.md' ! -path '*/_inbox/*' ! -path '*/_sources/*' \
  ! -path '*/_state/*' ! -name '_index.md'
```

每 20 個 entity 一批寫 `manifest.json` + 初始 `progress.json`，`STATUS.md` 加一列。

## Step 2 — 逐批建邊 (Process per Batch)

對批內每個 entity 的每個維度：

1. `找關係線索`：維度內文與其 `Sources:` 指向的 captures 中，
   提到其他 entity 的地方（呼叫、讀寫、依賴、討論）
2. `驗證佐證`：每條邊必須指得出發起方自身來源的直接佐證
   （`檔案:行`、DDL 外鍵、訊息時間戳、原文引句）。指不出 → 不建邊，
   弱訊號寫進維度內文一句話，或列入 Frontier
3. `逐字抄標題`：Read 目標 entity 檔，逐字複製 `##` 標題到 wikilink。
   禁止憑記憶拼寫 — 這是斷鏈的最大來源
4. `寫邊`：`- <relation> [[entity#Section]] — 佐證位置`；
   relation 動詞與方向規則依 kb-spec（含 `supersedes` / `contradicts`）；
   `mentions` 每維度 ≤ 2 條；logger/config/utils 不建邊
5. 整批做完 → 更新 `progress.json` 與 `STATUS.md`，才開始下一批

## Step 3 — 全圖重算 Backlinks（所有批次完成後）

1. 掃全圖正向邊（各維度 `References:` 清單）
2. 每個被指向的 entity：重寫其 `## Backlinks` 章節
   （保留 `<!-- auto-generated -->` 標記）
3. Backlinks 只能由正向邊推導；發現「無對應正向邊」的舊條目 → 刪除

## Step 4 — 重建 `_index.md`

依 kb-spec 的 `_index.md` 結構：

1. 註冊表：entity 清單（zone、type、維度數、最舊 truth tier）
2. Mermaid 總覽：zone 為 subgraph，邊聚合到 entity 層級，
   edge text 用雙引號包裹
3. `## Frontier`：合併既有條目 + 本輪新增（弱訊號關係、2-hop 外實體）
4. `## Unlinked`：無入邊/無出邊兩列，皆空寫「無 (None)」，
   兩列同現標 `(orphan)`

完成後 `progress.json` 標 `done`，報告：邊數、Backlinks 重算數、
orphan 清單、Frontier 新增數、建議下一步（`kb-verify`）。

## Common Mistakes

| 錯誤                             | 修正                                       |
| -------------------------------- | ------------------------------------------ |
| 憑記憶拼 wikilink 標題           | Read 目標檔逐字複製                        |
| 命名相似/分層慣例推斷邊          | 只憑發起方來源佐證建邊                     |
| 「被 X 呼叫」寫成正向邊          | 方向 = 發起者 → 接受者；反向交給 Backlinks |
| 手寫/補寫 Backlinks              | 一律由全圖正向邊重算                       |
| 批次未全完成就跑 Step 3/4        | 全域重算必須等所有批次 done                |
| `mentions` 浮濫                  | 每維度 ≤ 2 條，否則收斂為精確動詞或刪除    |
| 全部建完才寫進度                 | 每批處理完立即更新 progress.json           |
