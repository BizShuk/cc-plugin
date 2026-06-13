---
name: system-law-time
description: >
  Use when dealing with scheduling, event ordering, synchronization, or
  ensuring operations happen in the correct sequence. Triggers on:
  "event ordering", "race condition", "scheduling", "synchronization",
  "orchestration", "workflow sequencing", "orchestrator role".
version: "1.0.0"
metadata:
  type: reference
  tier: philosophy
  domain: 宇宙基石
  law: 時間 (Time)
  role: 指揮家 (Orchestrator)
---

# law-time — 時間法則 × 指揮家

> 法則領域：宇宙基石 (Cosmic Foundation)

## 法則定義

`時間法則 (Law of Time)` ✖ `指揮家 (Orchestrator)`

掌控時序流動，確保事件絕對同步。

## 系統對應

指揮家是時間的主宰。在系統中，時間法則體現為：

- `排程 (Scheduling)`：決定任務執行的先後與時機
- `事件同步 (Event Synchronization)`：協調分散式操作的完成順序
- `工作流編排 (Workflow Orchestration)`：串接多步驟業務流程
- `逾時控制 (Timeout Control)`：為每個操作設置時間邊界

## 架構問題

當時間法則被違反時，會出現：

- Race condition（時序不受控）
- Deadlock（相互等待造成時間停止）
- 無法重現的間歇性 bug（時間依賴隱藏在邏輯中）

## 設計原則

1. 將時序依賴顯式化，不依賴執行速度的假設
2. 為所有異步操作定義完成條件與超時
3. 使用 idempotency 讓重放安全
