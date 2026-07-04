---
name: kb-query
description: >
    Use when answering a question from the knowledge base — resolves the
    question to entities via the index, walks wikilink edges (2 hops max),
    and answers with truth-tier-labelled citations; unanswerable gaps are
    logged to the Frontier. Triggers on: "ask the knowledge base",
    "what does the kb say about", "查知識庫", "知識庫怎麼說", "kb query".
version: "1.0.0"
allowed-tools: Read, Bash, Glob, Grep
user-invocable: true
disable-model-invocation: true
effort: medium
---

# kb-query — 帶引用查詢

從知識庫回答問題。每個結論都必須指得出 entity 維度與其 truth tier；
知識庫答不了的部分明說，並把缺口記入 Frontier。本技能唯讀
（僅 `_index.md` 的 Frontier 可追加）。

## Checklist（複製到回覆中，逐項執行並打勾）

```md
- [ ] Step 1: 讀 _index.md 註冊表，定位相關 entities
- [ ] Step 2: 讀 entity 檔，沿邊擴展（最多 2 hop）
- [ ] Step 3: 組答案：結論先行 + 逐條引用 + truth tier
- [ ] Step 4: 答不了的缺口 → Frontier 登記
```

## Step 1 — 定位 (Locate)

1. 從問題抽關鍵詞（中英雙語，包含同義詞）
2. 先查 `_index.md` 註冊表；再全文搜：

    ```bash
    grep -rli '<keyword>' <proj>/ --include='*.md' \
      | grep -v '_state/' | head -20
    ```

3. curated 命中優先於 `_inbox/` 命中；inbox 內容引用時必須標注
   「未蒸餾 (raw capture)」

## Step 2 — 擴展 (Expand)

- Read 命中的 entity 檔，沿其 `References:` 與 `Backlinks` 擴展，
  最多 2 hop — 超過就停，避免整庫掃描
- 沿途注意 `supersedes` / `contradicts` 邊：被 supersede 的內容不作為
  現行結論；contradicts 存在時兩方都要呈現

## Step 3 — 回答 (Answer)

結論先行，然後逐條佐證：

```markdown
<一段結論>

依據 (Evidence):

- <事實一句> — [[service-a#Billing Cycle]]（truth: corroborated）
- <事實一句> — [[2026-07-04-payment-flow]]（raw capture，未蒸餾）

矛盾/保留 (Caveats):

- <兩方主張並列，指出 CONFLICT 待裁決>

知識庫未涵蓋 (Not in KB):

- <缺口一句>
```

硬規則：

- 沒有引用的句子不可出現在「依據」；一般常識推論放結論段並明說是推論
- `candidate` 或 raw capture 佐證必須標注，不可與 `confirmed` 混同呈現
- 找不到任何佐證 → 誠實回答「知識庫未涵蓋」，禁止編造

## Step 4 — 缺口登記 (Gap to Frontier)

答不了或佐證不足的問題，在 `_index.md` 的 `## Frontier` 追加一行：

```markdown
- GAP (<date>): <問題一句> — 建議來源：<repo/web/chat/schema>
```

這是知識庫的需求信號：下次 ingest 的優先清單。

## Common Mistakes

| 錯誤                             | 修正                                 |
| -------------------------------- | ------------------------------------ |
| 用模型自身知識冒充知識庫內容     | 每條依據必附 wikilink 引用           |
| candidate/raw 與 confirmed 混同  | truth tier 逐條標注                  |
| 無限沿邊擴展                     | 2 hop 上限                           |
| 被 supersede 的舊知識當現行結論  | 檢查 supersedes 邊                   |
| 答不了就沉默略過                 | 明列「未涵蓋」+ Frontier 登記        |
