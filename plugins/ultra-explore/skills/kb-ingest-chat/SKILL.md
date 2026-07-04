---
name: kb-ingest-chat
description: >
    Use when ingesting chat histories (Slack/Discord exports, meeting notes,
    pasted conversations, or the current session) into the knowledge base —
    distills decisions, facts, and commitments from noisy dialogue into raw
    captures with speaker attribution, processed in resumable batches with
    progress persisted. Triggers on: "ingest this chat", "save this
    conversation to knowledge base", "把對話加入知識庫", "匯入聊天紀錄",
    "記住這段討論".
version: "1.0.0"
allowed-tools: Read, Bash, Glob, Write, Edit
user-invocable: true
disable-model-invocation: true
effort: medium
context: fork
---

# kb-ingest-chat — 對話紀錄入庫

從對話（匯出檔、貼上文字、或本次 session）蒸餾出決策、事實與承諾，
寫入知識庫 `_inbox/`。對話是最高噪音的來源 — 核心是分離「說了什麼」
與「成立什麼」。儲存格式、指紋去重、狀態追蹤依 `kb-spec`，先讀它再開始。

## Checklist（複製到回覆中，逐項執行並打勾）

```md
- [ ] Step 0: 檢查 _state/runs/ — 續跑或新建 run
- [ ] Step 1: 界定輸入 → 切段分批（每批 200 則訊息）→ 寫 manifest.json
- [ ] Step 2: 逐批處理：分類萃取 → 隱私過濾 → 寫 captures → 更新 progress
- [ ] Step 3: 全部批次完成 → status = done → 報告
```

## Step 0 — 續跑檢查 (Resume Check)

- `_state/runs/` 有同來源 `in-progress` 的 `kb-ingest-chat` run
  → 從第一個未完成 batch 續跑，跳到 Step 2
- 沒有 → 新建 `run-id = <date>-chat-<channel>`

## Step 1 — 界定與分批 (Plan)

| 輸入型態          | 處理方式                                   |
| ----------------- | ------------------------------------------ |
| 匯出檔 (json/txt) | 按 thread/主題切段，每批約 200 則訊息      |
| 貼上文字          | 單批處理                                   |
| 目前 session      | 只取已成立的結論，單批處理                 |

寫 `manifest.json`（每批的訊息範圍：thread id 或行號區間）+ 初始
`progress.json`，`STATUS.md` 加一列。

## Step 2 — 逐批處理 (Process per Batch)

### 2a — 分類萃取，只收四類，其他閒聊一律丟棄

| 類別              | 判定                                | truth 初判  |
| ----------------- | ----------------------------------- | ----------- |
| 決策 (decision)   | 明確拍板：「就用 X」「決定不做 Y」  | `confirmed` |
| 事實 (fact)       | 第一人稱親述：「我們的 QPS 是 3k」  | `firsthand` |
| 承諾 (commitment) | 誰、何時、做什麼                    | `firsthand` |
| 主張 (claim)      | 轉述/聽說/推測：「聽說 Z 要下線」   | `candidate` |

硬規則：

- 每條必附發言者 + 時間戳（或訊息定位）；無法歸屬發言者 → 不收
- 被後續訊息推翻的內容不收，或在既有 capture 註記 superseded
- 相對日期轉絕對日期（「下週五」→ `2026-07-10`）

### 2b — 隱私過濾 (Privacy Gate)

涉及個資、憑證、token、內部機密 → 不入庫，log.md 記 `EXCLUDED-PRIVACY` 一行
（只記類別，不記內容本身）。

### 2c — 寫 captures 並落盤

- 每個主題一個 capture（依 kb-spec 格式與指紋查重）：
  `sources.type: chat`、`ref` 含頻道與日期範圍
- 建立/更新 `_sources/chat-<channel>.md`
- 更新 `progress.json` 與 `STATUS.md`，才開始下一批

## Step 3 — 收尾 (Finish)

- `status: done`；報告：批次數、各類別條數、被排除的隱私項數、
  corroboration 命中、建議下一步（`kb-distill`）

## Common Mistakes

| 錯誤                      | 修正                                   |
| ------------------------- | -------------------------------------- |
| 把討論過程當結論入庫      | 只收拍板決策；未決事項進 Frontier      |
| 轉述/聽說標成 `firsthand` | 第一人稱親述才算；轉述一律 `candidate` |
| 遺漏發言者歸屬            | 附發言者 + 時間戳，無法歸屬則不收      |
| 相對日期原樣入庫          | 一律轉絕對日期                         |
| 憑證、token、個資照抄     | 隱私門檻過濾並回報                     |
| 全部處理完才寫進度        | 每批處理完立即更新 progress.json       |
