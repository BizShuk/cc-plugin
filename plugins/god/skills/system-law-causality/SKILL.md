---
name: system-law-causality
description: >
  Use when implementing tracing, debugging, root cause analysis, or building
  audit trails in a system. Triggers on: "tracing", "root cause analysis",
  "debugging", "audit log", "distributed tracing", "why did this happen",
  "cause and effect", "experimenter role".
version: "1.0.0"
metadata:
  type: reference
  tier: philosophy
  domain: 宇宙基石
  law: 因果 (Causality)
  role: 實驗家 (Experimenter)
---

# law-causality — 因果法則 × 實驗家

> 法則領域：宇宙基石 (Cosmic Foundation)

## 法則定義

`因果法則 (Law of Causality)` ✖ `實驗家 (Experimenter)`

建立因果驗證，追蹤並除錯。

## 系統對應

實驗家是因果鏈的追蹤者。在系統中，因果法則體現為：

- `分散式追蹤 (Distributed Tracing)`：跨服務追蹤請求的完整因果鏈
- `稽核日誌 (Audit Log)`：記錄「誰在何時做了什麼，結果是什麼」
- `根本原因分析 (Root Cause Analysis)`：從現象回溯到最初原因
- `可重現性 (Reproducibility)`：相同輸入必須得到相同輸出

## 架構問題

當因果法則被違反時，會出現：

- 「這個 bug 無法重現」（因果鏈斷裂）
- 無法確定某個狀態變化的起因（沒有 trace）
- 調試只能靠猜測（缺乏因果可觀測性）

## 設計原則

1. 為每個請求注入唯一的 trace ID，貫穿所有服務
2. 所有重要狀態變更都要有對應的原因紀錄
3. 測試要驗證因果關係，不只是驗證結果
