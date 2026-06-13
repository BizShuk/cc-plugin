---
name: system-law-chaos
description: >
  Use when exploring unknown territory, conducting feasibility studies,
  or deliberately introducing uncertainty to discover system limits.
  Triggers on: "unknown territory", "new technology", "proof of concept",
  "spike", "boundary testing", "what happens if", "explorer role".
version: "1.0.0"
metadata:
  type: reference
  tier: philosophy
  domain: 系統意志
  law: 混沌 (Chaos)
  role: 探險家 (Explorer)
---

# law-chaos — 混沌法則 × 探險家

> 法則領域：系統意志 (System Will)

## 法則定義

`混沌法則 (Law of Chaos)` ✖ `探險家 (Explorer)`

踏入無序，尋找邊界與新知。

## 系統對應

探險家是混沌的領航者。在系統中，混沌法則體現為：

- `混沌工程 (Chaos Engineering)`：主動注入故障，發現隱性弱點
- `概念驗證 (PoC)`：在未知領域的最小化試探
- `邊界測試 (Boundary Testing)`：找到系統能承受的極限
- `技術雷達 (Tech Radar)`：系統性地評估新技術的可行性

## 架構問題

當混沌法則被忽視時，會出現：

- 生產環境才首次遇到故障情境（沒有主動探索）
- 架構無法適應新需求（沒有探索新技術）
- 對系統的真實極限一無所知（缺乏邊界意識）

## 設計原則

1. 混沌不是敵人，是揭示真相的工具
2. 主動在非生產環境製造混沌，勝過被動在生產遭遇混沌
3. 每次探索都要設定「安全邊界」，防止混沌蔓延
