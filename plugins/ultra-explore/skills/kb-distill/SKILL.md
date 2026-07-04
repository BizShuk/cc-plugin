---
name: kb-distill
description: >
    Use when distilling raw inbox captures into curated knowledge base entities
    — resolves identity (dedupe/merge), applies the truth qualification gate,
    assigns zones, and writes entity dimension sections with sources cited.
    Processes captures in resumable batches with progress persisted. Triggers
    on: "distill the inbox", "promote captures to entities", "蒸餾知識庫",
    "整理 inbox", "kb distill".
version: "1.0.0"
allowed-tools: Read, Bash, Glob, Grep, Write, Edit
user-invocable: true
disable-model-invocation: true
effort: high
context: fork
---

# kb-distill — 蒸餾 inbox 為 curated entities

把 `_inbox/` 的 raw captures 蒸餾成 curated 區的 entity 檔。
這是知識庫的品質閘門：身分先於連結、真實性先於入庫。
儲存格式、truth tiers、狀態追蹤依 `kb-spec`，先讀它再開始。
本技能不建跨 entity 邊（那是 `kb-connect` 的事）。

## Checklist（複製到回覆中，逐項執行並打勾）

```md
- [ ] Step 0: 檢查 _state/runs/ — 續跑或新建 run
- [ ] Step 1: 盤點 status: raw 的 captures → 寫 manifest.json（每批 20 個）
- [ ] Step 2: 逐批處理：身分歸併 → 真實性閘門 → 寫 entity → 標記 capture → 更新 progress
- [ ] Step 3: 全部批次完成 → status = done → 報告
```

## Step 0 — 續跑檢查 (Resume Check)

同 `kb-spec` 規則：有 `in-progress` 的 `kb-distill` run 就續跑，
否則新建 `run-id = <date>-distill-<slug>`。

## Step 1 — 盤點與分批 (Plan)

```bash
grep -l '^status: raw' <proj>/_inbox/*.md
```

- 依 `zone-hint` 分組排序（同 zone 的 captures 相鄰，歸併判斷更準）
- 每 20 個 captures 一批寫 `manifest.json` + 初始 `progress.json`，
  `STATUS.md` 加一列

## Step 2 — 逐批處理 (Process per Batch)

對批內每個 capture 依序執行四個閘門：

### 2a — 身分歸併 (Identity Resolution)

1. 從候選事實找出它談論的 entity（服務、系統、資料庫、概念、決策…）
2. 查既有 entity：

    ```bash
    grep -rl 'name: <candidate>' <proj>/*/           # 檔名與 frontmatter
    grep -rli '<alias>' <proj>/*/ --include='*.md'   # aliases 也要查
    ```

3. 判定，只有三種結果：
    - `已存在` → 併入既有 entity（新維度或補充既有維度），別名進 `aliases`
    - `新 entity` → 依 kb-spec 判斷是否夠格：可獨立部署/維運/閱讀的單位才是
      entity；單一函數、單一表、單一訊息永遠不是（歸併為某 entity 的維度）
    - `不夠格也無歸屬` → 留在 inbox，log.md 記原因

### 2b — 真實性閘門 (Truth Gate)

逐條事實檢查 truth tier（規則見 kb-spec）：

- `confirmed` / `firsthand` → 可入 curated 區
- `candidate` → 查 `_sources/`：有 2+ 獨立來源指向同一事實 → 升
  `corroborated` 入庫；否則留在 inbox 並在 `_index.md` 的 Frontier 登記
- 與既有維度矛盾 → 兩邊都不動，Frontier 登記 `CONFLICT` 條目待人工裁決

### 2c — 寫 entity 檔 (Write)

- 依 kb-spec entity 格式：frontmatter（name/type/zone/tags/aliases/sources）、
  維度章節（`kind:` + `truth:` + 內文 + `Sources:` 行）
- zone 依業務域；新 zone 需 3+ entities 才成立，否則掛最近的既有 zone
- 每 entity 維度 2~12；只有 1 個維度 → 重新檢視是否該併入其他 entity
- `References:` 先留空或只放同檔 `[[#Section]]`；跨 entity 邊交給 kb-connect

### 2d — 標記與落盤 (Mark & Persist)

- 處理完的 capture：`status: raw` → `distilled`（Edit frontmatter）
- 更新 `progress.json`（含 `outputs` 的 entity 路徑）與 `STATUS.md`，
  才開始下一批

## Step 3 — 收尾 (Finish)

- `status: done`；報告：處理 capture 數、新建/更新 entity 數、
  留在 inbox 的 candidate 數、Frontier 新增條目（含 CONFLICT）、
  建議下一步（`kb-connect`）

## Common Mistakes

| 錯誤                               | 修正                                       |
| ---------------------------------- | ------------------------------------------ |
| candidate 事實直接進 curated 區    | 2+ 來源或人工確認才入庫，否則留 inbox      |
| 沒查 aliases 就新建 entity（重複） | 檔名 + frontmatter + aliases 三路都要查    |
| 函數/單表當 entity（粒度爆炸）     | 歸併為維度；entity 是可獨立維運的單位      |
| 矛盾內容直接覆寫舊維度             | 兩邊不動，Frontier 記 CONFLICT 待裁決      |
| 在這一步建跨 entity 邊             | 邊是 kb-connect 的職責                     |
| capture 處理完忘記改 status        | 每個 capture 落盤時同步標 distilled        |
| 全部處理完才寫進度                 | 每批處理完立即更新 progress.json           |
