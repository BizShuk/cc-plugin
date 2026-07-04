---
name: kb-ingest-web
description: >
    Use when ingesting web links (articles, docs sites, index pages, videos)
    into the knowledge base — fetches with markitdown (scrapling fallback),
    strips page chrome, and writes clean raw captures with provenance into the
    KB inbox. Handles large document sets (10000+) in resumable batches with
    per-batch progress persisted. Triggers on: "ingest this url",
    "add this link to knowledge base", "把這篇文章加入知識庫", "匯入網頁知識".
version: "1.0.0"
allowed-tools: Read, Bash, Write, Edit, WebFetch
user-invocable: true
disable-model-invocation: true
effort: medium
context: fork
---

# kb-ingest-web — 網頁入庫

把一個或一批 URL 轉成乾淨的 raw capture 寫入知識庫 `_inbox/`。
抓取與去噪沿用 `content-summarizer` / `markitdown` 的規則；
儲存格式、指紋去重、狀態追蹤依 `kb-spec`，先讀它再開始。

## Checklist（複製到回覆中，逐項執行並打勾）

```md
- [ ] Step 0: 檢查 _state/runs/ — 續跑或新建 run
- [ ] Step 1: 展開 URL 清單 → 寫 manifest.json（每批 20 份）
- [ ] Step 2: 逐批處理：抓取 → 去噪 → 寫 capture → 更新 progress
- [ ] Step 3: 全部批次完成 → status = done → 報告
```

單一 URL 也走同一流程（manifest 只有一批一項），狀態一樣要落盤 —
大量入庫時單一 URL 的紀錄就是去重與 corroboration 的依據。

## Step 0 — 續跑檢查 (Resume Check)

- `_state/runs/` 有同來源 `in-progress` 的 `kb-ingest-web` run
  → 從第一個未完成 batch 續跑，跳到 Step 2
- 沒有 → 新建 `run-id = <date>-web-<slug>`

## Step 1 — 展開與分批 (Plan)

1. 展開輸入為 URL 清單：
    - 單一文章 URL → 一項
    - Index/清單頁 → 先抓 index 本身，列出項目連結；預設取 top 20，
      更多需使用者確認。項目 URL 各自成為一項（只往下一層，不遞迴）
    - 使用者給的 URL 檔案/清單 → 逐行一項
2. 每 20 份一批寫 `manifest.json` + 初始 `progress.json`，`STATUS.md` 加一列

## Step 2 — 逐批處理 (Process per Batch)

對批內每個 URL 依序執行：

1. `指紋預查`：URL 已出現在 `_sources/` 或既有 capture → 跳過並記 log
2. `抓取`（fetch-first，禁止憑 URL 猜內容）：

    ```bash
    markitdown "<url>"        # 首選：靜態頁 / YouTube 逐字稿
    ```

    - 近空、login/paywall、cookie shell → 換 `scrapling`
      （`--ai-targeted` 或 `--css-selector "main"`）
    - 仍失敗 → 該項標 `failed` 記入 log.md，繼續下一項，不中斷整批
3. `去噪`：移除導覽列、頁尾、側欄、廣告、cookie、留言區
   （沿用 `content-summarizer` Step 2 清單）。Index 頁的項目連結是內容本體，不可剝除
4. `寫 capture`（依 kb-spec 格式與指紋查重）：
    - `sources.type: web`、`ref` 完整 URL、`truth: candidate`
      （網頁內容一律 candidate，第二個獨立來源出現才升 `corroborated`）
    - `# 內容 (Content)` 開頭補 `Source: [頁面標題](URL)`，保留原文連結
    - `# 候選事實 (Candidate Facts)`：每條一句，附原文引句或段落定位；
      作者主張記為候選事實，自己的推論不入庫
    - 建立/更新 `_sources/web-<domain>.md`（同網域共用一個來源檔）

整批做完 → 更新 `progress.json` 與 `STATUS.md`，才開始下一批。

## Step 3 — 收尾 (Finish)

- `status: done`；報告：成功/失敗/跳過數、capture 清單、
  corroboration 命中（指紋或主題與既有 capture 重疊者）、建議下一步（`kb-distill`）

## Common Mistakes

| 錯誤                          | 修正                                   |
| ----------------------------- | -------------------------------------- |
| 沒抓到內容就寫 capture        | fetch-first；失敗標 failed 並記 log    |
| 單一 URL 失敗就中斷整批       | 跳過該項續跑，收尾時一併回報           |
| 網頁主張直接標 `corroborated` | 單一網頁 = `candidate`                 |
| 把站台 chrome 一起入庫        | 去噪後才落檔                           |
| Index 頁無限遞迴              | 一層深、top 20，超過先確認             |
| 全部抓完才寫進度              | 每批處理完立即更新 progress.json       |
