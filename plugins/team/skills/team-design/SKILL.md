---
name: team-design
description: >
    Use when the user wants to design, structure, or plan a cross-functional AI agent team for a specific project or software workflow. Triggers on: "team design", "team structuring", "團隊設計", "設計團隊", "規劃團隊", "AI團隊規劃".
version: "1.0.0"
allowed-tools: []
user-invocable: true
disable-model-invocation: false
effort: medium
metadata:
    type: reference
    platforms: [macos, linux]
---

# 團隊架構設計技能 (Team Architecture Design Skill)

此技能提供了一套系統化的框架，協助使用者根據專案需求規劃跨職能 AI `代理團隊 (Agent Team)` 的角色編制與主要職責。

## 核心處理原則 (Core Principles)

### 1. 單一職責原則 (Single Responsibility Principle)
- 每個 AI 代理（Agent）應專注於單一核心任務，避免一個代理包山包海。
- 職責過於雜亂會導致模型輸出品質不穩定且難以進行單點優化。

### 2. 角色定位分析 (Role Definition Analysis)
- `專案交付物 (Deliverables)`：明確定義此專案最終需要產出什麼檔案或資料。
- `必備工作站 (Stations Required)`：將專案的生成管線拆解為獨立的工序站位。
- `代理編制 (Agent Staffing)`：為每個工作站分配一個最適合的 `角色 (Role)`。

---

## 輸出格式範本 (Output Format Template)

請依以下結構規劃團隊架構：

```markdown
# 團隊架構設計方案：[專案名稱]

## 專案目標與交付物 (Project Goals & Deliverables)
- 最終目標：[描述專案要完成的任務]
- 主要交付物：[例如 代碼庫、報告、分析圖表]

## 團隊編制與站位規劃 (Team Staffing & Stations)

| # | 角色名稱 (Role Name) | 負責工作站位 (Station) | 核心任務 (Core Task) |
|---|---|---|---|
| 1 | [例如 產品經理] | 需求定義站 | 產出功能說明與 PRD |
| 2 | [例如 後端工程師] | 邏輯實作站 | 撰寫 API 與資料庫設計 |

## 角色交互關係 (Role Interaction)
- 流水線工序：[例如 產品經理 -> 後端工程師 -> 網站可靠性工程師]
- 關鍵交付邊界：[描述角色之間傳遞的資料格式與驗收標準]
```
