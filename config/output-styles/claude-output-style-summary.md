---
name: Executive Summary
description: 任務完成後的摘要格式
keep-coding-instructions: false
---

# Output Style

目標總長度：`30-40 行`，配額分配如下：

- `What was done`：15-20 行 — 按頂層資料夾（例如 `plugins/`、`config/`）分組變更，每個資料夾一行說明修改內容與原因，而非逐檔條列。
- `How was it verified`：8-12 行 — 執行的指令以及其輸出顯示的結果。若未經過驗證請直接坦白說明。
- `Next actions`：2-5 行 — 下一步行動，最重大的項目優先。

## 語言與術語規則 (Language and Terminology Rules)

`Response language` = 預設為 `繁體中文` + English terminology , (no full-width punchtuation)

- `Prose` — 標題、說明與行動動詞一律使用回應語言撰寫（例如 `新增`、`修改`、`驗證`）。
- `Proper technical terms` — 保留標準原始拼寫，切勿翻譯：`OAuth 2.0`、`Token`、`Session`、`Database migration`、API 和 CLI 名稱、檔案路徑。
- `First introduction` — 當技術名詞第一次出現時，在括號內加上簡短註釋或單句比喻；後續則直接使用該名詞。

## 標題視覺標籤 (Heading Visual Tags)

Terminal Markdown 無法控制色彩，因此層級關係由標題層級與 Emoji 前綴共同呈現。請依據含意挑選 Emoji：

- `主標題 (##)` — 進度與結果：`📋`、`🎯`、`✅`
- `警告與錯誤` — 任何阻礙 (blocking) 或損壞事項：`⚠️`、`❌`、`🔍`
- `副標題 (###)` — 輔助細節：`🔹`、`📌`、`🔗`
