---
name: universal-consolidate
description: >
  Use when merging same-kind artifacts into one canonical version —
  deduplicating skills, unifying duplicate docs or configs, resolving
  conflicting designs, merging overlapping findings — or when two
  incompatible abstractions must coexist. Triggers on: "merge",
  "deduplicate", "unify", "canonical", "conflicting designs", "overlap",
  "design conflict", "how to combine", "合併", "去重", "整併".
version: "2.0.0"
metadata:
  type: pattern
  tier: philosophy
  operator: consolidate
---

# universal-consolidate — 凝聚算子

`簽名 (Signature)`：{同類產物 ×N} → 1 個典範 (canonical)，N → 1

把多個同類產物凝聚為單一真相來源。與 `universal-aggregate` 的分界：consolidate 消除`同類冗餘`（N 份變 1 份），aggregate 組合`異類部件`（部件變整體）。

## 三階段程序

### 1. 同構偵測 (Detect Isomorphism)

剝除術語表象，判斷候選物是否真的同義：

- 指紋比對：正規化後內容是否重疊？
- 語意比對：抽掉名字後，結構與行為是否一致？
- 警惕`假同義 (false synonym)`：同名不同義（auth.account ≠ banking.account）不可合併

### 2. 選融合法 (Select Fusion Method)

四大融合思維工具，依衝突類型選用：

| 衝突類型 | 融合法 | 做法 | 驗證 |
| :--- | :--- | :--- | :--- |
| 同時存在互不影響 | 正交性 (Orthogonality) | 讓概念垂直交叉，各撐一維 | 改 A 不迫使改 B → 正交 |
| 表面不同本質相同 | 同構性 (Isomorphism) | 剝除術語，收斂到底層共同形狀 | 抽掉術語後運算結構一致 |
| 直接對立無法共存 | 催化劑 (Catalyst) | 在兩者間插入輕量介面層 | 移除催化劑則兩端立即衝突 |
| 時序上交替出現 | 辯證循環 (Dialectical Cycle) | 對立概念拉長至時間軸成生命週期 | 能畫出 A → B → A 的循環 |

範例：`快取`、`快照`、`備份` 同構為「凍結狀態的副本」；`安全` 與 `可觀測` 以 Auth 層為催化劑共存；`建設` 與 `破壞` 辯證為 CI/CD 循環。

### 3. 無損合併 (Lossless Merge)

- 典範版本必須涵蓋所有來源的`不可丟棄資訊`（不變式、edge case、觸發條件）
- 被合併者原地刪除，留下指向典範的引用（不留分叉副本）
- 合併後以 `universal-review` 驗證：原本每個來源能回答的問題，典範版本仍能回答

## 反模式

| 反模式 | 問題 | 修正 |
| :--- | :--- | :--- |
| 合併假同義 | 語意壞掉、需翻譯 | 先過同構偵測 |
| 過度合併 | 產生上帝物件/上帝文件 | 不同關注點用正交性分開 |
| 有損合併 | 來源的不變式遺失 | 合併前列出不可丟棄清單 |
| 留下舊副本 | 定義漂移 (definition drift) | 刪除 + 引用重導 |

## 算子組合

`consolidate` 常接在 `generate`（新物與既有物去重）與 `review`（合併重複發現）之後；其產出的典範物是 `universal-aggregate` 的乾淨部件。
