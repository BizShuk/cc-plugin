---
name: kb-ingest-schema
description: >
    Use when ingesting a storage schema into the knowledge base — relational
    databases (tables/DDL), key-value stores (key patterns), message queues
    (topics/queues), or file storage (directory layouts/formats). Writes
    datastore captures with object inventories in resumable batches. Triggers
    on: "ingest this schema", "add database to knowledge base",
    "把資料庫結構加入知識庫", "匯入 schema", "記錄 MQ topics".
version: "1.0.0"
allowed-tools: Read, Bash, Glob, Grep, Write, Edit
user-invocable: true
disable-model-invocation: true
effort: medium
context: fork
---

# kb-ingest-schema — 儲存結構入庫

把四類儲存結構轉成 raw captures 寫入知識庫 `_inbox/`：
關聯式資料庫、key-value、message queue、檔案儲存。
儲存格式、指紋去重、狀態追蹤依 `kb-spec`，先讀它再開始。

## Checklist（複製到回覆中，逐項執行並打勾）

```md
- [ ] Step 0: 檢查 _state/runs/ — 續跑或新建 run
- [ ] Step 1: 盤點 schema 物件 → 寫 manifest.json（每批 30 物件）
- [ ] Step 2: 逐批處理：讀結構 → 萃取業務語意 → 寫 captures → 更新 progress
- [ ] Step 3: 全部批次完成 → status = done → 報告
```

## Step 0 — 續跑檢查 (Resume Check)

同 `kb-spec` 狀態追蹤規則：有同來源 `in-progress` run 就續跑，
否則新建 `run-id = <date>-schema-<name>`。

## Step 1 — 盤點與分批 (Plan)

按儲存類型取得物件清單（只讀，禁止任何寫入目標儲存的指令）：

| 類型       | 盤點方式                                                        | 物件單位       |
| ---------- | --------------------------------------------------------------- | -------------- |
| relational | DDL 檔 / migration 目錄 / `sqlite3 <db> .schema` / info_schema  | table, view    |
| key-value  | key 命名規則文件 / `redis-cli --scan --count 100`（抽樣即可）   | key pattern    |
| queue      | topic/queue 設定檔、broker 管理 API                             | topic, queue   |
| file       | 目錄佈局 + 檔名慣例 + 格式（`find <root> -maxdepth 3 -type d`） | 目錄/檔案格式  |

- 活體資料庫優先找 repo 內的 DDL/migration 檔，讀檔比連線可靠且可追溯
- 每 30 物件一批寫 `manifest.json` + 初始 `progress.json`，`STATUS.md` 加一列

## Step 2 — 逐批處理 (Process per Batch)

對批內每個物件萃取業務語意，不是照抄欄位清單：

1. `結構`：欄位/key 段/訊息格式 — 只列業務關鍵欄位（主鍵、外鍵、狀態欄、
   金額/數量欄），非全欄位傾印
2. `業務語意`：這個 table/topic/pattern 承載什麼業務物件？誰寫誰讀？
   （從 DDL 註解、命名、外鍵、repo 交叉引用推斷；推斷標 `candidate`）
3. `生命週期線索`：狀態欄 enum 值、TTL、retention 設定 — 逐字照抄，不虛構
4. `關聯線索`：外鍵、topic 的 producer/consumer、跨表 join —
   這些是未來 kb-connect 建邊的佐證，附定位（DDL 檔:行、設定鍵）

寫 captures（依 kb-spec 格式與指紋查重）：

- 每個儲存系統一個主 capture（datastore entity 候選）+ 物件群組 captures
  （同業務域的 tables/topics 合為一個 capture，勿一表一檔）
- `sources.type: schema`；DDL 逐字內容 → `firsthand`，語意推斷 → `candidate`
- `zone-hint` 依業務域
- 建立/更新 `_sources/schema-<name>.md`

整批做完 → 更新 `progress.json` 與 `STATUS.md`，才開始下一批。

## Step 3 — 收尾 (Finish)

- `status: done`；報告：物件數、capture 數、datastore entity 候選、
  發現的關聯線索數、建議下一步（`kb-distill`）

## Common Mistakes

| 錯誤                             | 修正                                     |
| -------------------------------- | ---------------------------------------- |
| 全欄位傾印當知識                 | 只收業務關鍵欄位與語意                   |
| 一張表一個 capture（粒度爆炸）   | 同業務域的物件群組合為一個 capture       |
| 對活體資料庫下寫入/掃全庫指令    | 只讀；KV 用抽樣；優先讀 DDL 檔           |
| 語意推斷標 `firsthand`           | DDL 逐字才是 firsthand；推斷 `candidate` |
| 狀態 enum 憑印象補全             | 逐字照抄，缺就寫「未偵測到」             |
| 全部處理完才寫進度               | 每批處理完立即更新 progress.json         |
