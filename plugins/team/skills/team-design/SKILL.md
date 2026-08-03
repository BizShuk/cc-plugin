---
name: team-design
description: >
    Use when the user wants to design, structure, or plan a cross-functional AI agent team for a specific project or software workflow, including its orchestration pattern and project-level shared rules. Triggers on: "team design", "team structuring", "team orchestration", "sequential pipeline", "manager pattern", "團隊設計", "規劃團隊", "團隊編排", "AI團隊規劃".
version: "2.0.0"
user-invocable: true
disable-model-invocation: false
effort: medium
metadata:
    type: reference
    platforms: [macos, linux]
---

# 團隊架構設計技能 (Team Architecture Design Skill)

此技能提供了一套系統化的框架，協助使用者根據專案需求規劃跨職能 AI `代理團隊 (Agent Team)` 的角色編制與主要職責，並決定代理之間的 `編排模式 (Orchestration Pattern)` 與專案層級 `共通規則 (Shared Rules)`。

產出兩份文件：團隊架構設計方案（誰做什麼）與團隊編排設定（怎麼串起來）。單一角色的 `系統提示 (System Prompt)` 撰寫交給 `role-generator`。

## 核心處理原則 (Core Principles)

### 1. 單一職責原則 (Single Responsibility Principle)

- 每個 AI 代理（Agent）應專注於單一核心任務，避免一個代理包山包海。
- 職責過於雜亂會導致模型輸出品質不穩定且難以進行單點優化。

### 2. 角色定位分析 (Role Definition Analysis)

- `專案交付物 (Deliverables)`：明確定義此專案最終需要產出什麼檔案或資料。
- `必備工作站 (Stations Required)`：將專案的生成管線拆解為獨立的工序站位。
- `代理編制 (Agent Staffing)`：為每個工作站分配一個最適合的 `角色 (Role)`。

## 核心編排模式 (Core Orchestration Patterns)

### 1. 協調者模式 (Orchestrator Pattern)

- 結構：一個 `協調者代理 (Orchestrator Agent)` 負責拆解任務、指派工作、並彙整最終結果。
- 適用情境：複雜、需動態規劃或多步驟分析的任務。

### 2. 流水線模式 (Pipeline Pattern)

- 結構：多個代理以固定順序傳遞任務（例如：Scraper -> Validator -> Analyst）。
- 適用情境：步驟明確、工序標準化且依賴前一站產出的任務。

### 專案層級共享設定 (Project-Level Shared Config)

- 將編碼風格、特定名詞術語表、統一輸出語言限制及可用工具清單等共通規則提取出來。
- 避免在每個代理的 `系統提示 (System Prompt)` 中重複貼上相同的共通限制。

---

## 輸出格式範本 A：團隊架構 (Output Template A — Team Architecture)

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

## 輸出格式範本 B：編排與共享設定 (Output Template B — Orchestration & Shared Config)

請依以下結構規劃編排與共享設定：

```markdown
# 團隊編排與共享設定：[專案名稱]

## 1. 協作模式選擇 (Orchestration Pattern Choice)
- 採用的模式：[協調者模式 (Orchestrator) 或 流水線模式 (Pipeline)]
- 選擇原因：[簡述為何此模式最符合專案]

## 2. 協作流程拓撲 (Collaboration Topology)
- [用 indented list 或 minimalist arrow 標示流程，例如：PM -> Backend -> QA -> SRE]

## 3. 專案層級共享限制 (Project-Level Shared Rules)
- 輸出語言限制：[例如 Always respond in Traditional Chinese, keeping technical terms in English.]
- 專案程式風格：[例如 Go 1.25 conventions, GORM SQLite]
- 工具調用限制：[例如 僅允許呼叫 read_file 與 write_to_file]
```
