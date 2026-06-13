---
name: system-law-destruction
description: >
  Use when designing stress tests, load tests, deprecation strategies,
  technical debt removal, or deliberately breaking things to improve
  system robustness. Triggers on: "stress test", "load test", "deprecation",
  "technical debt", "remove legacy", "break things", "experimenter role".
version: "1.0.0"
metadata:
  type: reference
  tier: philosophy
  domain: 系統意志
  law: 破壞 (Destruction)
  role: 實驗家 (Experimenter)
---

# law-destruction — 破壞法則 × 實驗家

> 法則領域：系統意志 (System Will)

## 法則定義

`破壞法則 (Law of Destruction)` ✖ `實驗家 (Experimenter)`

摧毀冗餘，執行破壞性壓力測試。

## 系統對應

實驗家用破壞來淨化系統：

- `壓力測試 (Stress Test)`：找到系統在極限壓力下的真實行為
- `技術債清償 (Tech Debt Removal)`：主動消除腐化的舊設計
- `廢棄策略 (Deprecation Strategy)`：有序地終止過時的介面與功能
- `破壞性實驗 (Destructive Experiments)`：透過刻意破壞來驗證韌性

## 架構問題

當破壞法則被忽視時，系統積累腐朽：

- 技術債持續增加，最終無法維護
- 廢棄的 API 繼續被依賴，無法演進
- 不知道系統在壓力下的真實極限

## 設計原則

1. 破壞是演化的必要條件，不要恐懼刪除
2. 每次破壞性測試後，系統應該變得更強韌（而非只是修復）
3. 廢棄要有明確的時間表與遷移路徑，不能無限期保留
