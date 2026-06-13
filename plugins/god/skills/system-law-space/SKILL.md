---
name: system-law-space
description: >
  Use when designing namespaces, defining service boundaries, structuring
  package layouts, or assigning resource ownership in a system. Triggers on:
  "namespace", "service boundary", "domain boundary", "directory structure",
  "resource ownership", "architect role".
version: "1.0.0"
metadata:
  type: reference
  tier: philosophy
  domain: 宇宙基石
  law: 空間 (Space)
  role: 建築師 (Architect)
---

# law-space — 空間法則 × 建築師

> 法則領域：宇宙基石 (Cosmic Foundation)

## 法則定義

`空間法則 (Law of Space)` ✖ `建築師 (Architect)`

劃定領域，構築結構與命名空間。

## 系統對應

建築師是空間的主宰。在系統中，空間法則體現為：

- `命名空間 (Namespace)`：隔離不同業務域的資源，避免衝突
- `服務邊界 (Service Boundary)`：明確每個元件的責任範圍
- `目錄結構 (Directory Layout)`：以層次化結構反映業務邏輯
- `資源歸屬 (Resource Ownership)`：每份資料只有唯一的擁有者

## 架構問題

當空間法則被違反時，會出現：

- 不同服務共用同一資料庫 table（邊界模糊）
- 命名衝突導致的 bug（命名空間缺失）
- 「什麼都做」的上帝物件（職責無邊界）

## 設計原則

1. 為每個領域劃定清晰的物理或邏輯邊界
2. 邊界內的事物完全自治，邊界外的交互走明確介面
3. 命名即文件：好的名稱讓邊界自我解釋
