---
name: universal-review
description: >
  Use when evaluating any artifact for defects, gaps, or quality — code
  review, document review, plan validation, architecture review, skill
  audit — especially when needing systematic coverage instead of ad-hoc
  inspection. Triggers on: "review", "audit", "evaluate", "check quality",
  "find issues", "what's missing", "審查", "評估", "找問題".
version: "2.0.0"
metadata:
  type: pattern
  tier: philosophy
  operator: review
---

# universal-review — 審視算子

`簽名 (Signature)`：產物 (artifact) → 發現 (findings)，1 → N

以透鏡系統性審視任何產物。兩種缺陷同等重要：`錯的東西 (defect)` 與 `缺的東西 (gap)` — 後者靠負空間推理 (negative-space reasoning) 才抓得到。

## 四階段程序

### 1. 選透鏡 (Select Lenses)

從 `system-laws` 13 法則中挑選適用透鏡。全面審查用全部 13 條；聚焦審查挑 3-5 條最相關。

### 2. 正查 (Scan for Defects)

逐透鏡掃描「錯的東西」。每個發現記錄：

- `summary`：一句話陳述缺陷
- `failure scenario`：具體輸入/狀態 → 錯誤結果
- `severity`：critical / major / minor
- `confidence`：high / medium / low

### 3. 負查 (Scan for Gaps)

逐透鏡問「這法則對應的元件在哪？」空格即候選缺口，判定：

- (a) 真的不需要（記錄理由）
- (b) 存在但未提及（補對映）
- (c) 真正缺口（列入 findings）

### 4. 對抗驗證 (Adversarial Verify)

對每個發現嘗試反駁：「這真的會發生嗎？給出反例。」反駁成功的發現剔除。只交付驗證存活的發現，依嚴重度排序。

## 目標類型適配表

| 產物類型 | 首選透鏡 | 常見缺口 |
| :--- | :--- | :--- |
| 程式碼 | 因果、時間、生命 | 錯誤處理、timeout、可觀測性 |
| 架構 | 全部 13 條 | 缺光明（黑盒）、缺黑暗（機密裸奔） |
| 文件 | 空間、因果 | 讀者不明、結論埋在後面 |
| 計畫 | 時間、破壞、混沌 | 無回退路徑、未測最大不確定性 |
| 資料 schema | 空間、時間、冰霜 | 語意缺失、無 retention、0/1 值域不明 |

## 反模式

| 反模式 | 問題 | 修正 |
| :--- | :--- | :--- |
| 只正查不負查 | 缺的東西永遠找不到 | 負空間逐透鏡掃 |
| 發現不含失敗情境 | 無法行動、無法驗證 | 每個發現附具體 scenario |
| 跳過對抗驗證 | 似是而非的發現污染下游 | 先自我反駁再交付 |
| 審查者即生成者且不換透鏡 | 盲點相同 | 換一組透鏡或換視角重審 |

## 算子組合

`review` 的 findings 是 `universal-consolidate`（合併重複發現）與 `universal-evolve`（作為迭代的 feedback）的輸入。
