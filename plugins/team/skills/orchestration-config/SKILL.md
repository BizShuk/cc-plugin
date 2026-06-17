---
name: orchestration-config
description: >
    Use when the user wants to set up, configure, or design the orchestration pattern, communication logic, or project-level shared rules for an AI agent team. Triggers on: "orchestration config", "team orchestration", "sequential pipeline", "manager pattern", "團隊編排", "編排設定", "流水線設定".
version: "1.0.0"
allowed-tools: []
user-invocable: true
disable-model-invocation: false
effort: medium
metadata:
    type: reference
    platforms: [macos, linux]
---

# 團隊編排與共享設定技能 (Team Orchestration & Shared Config Skill)

此技能協助使用者規劃多代理協作的 `編排模式 (Orchestration Patterns)`，並設定專案層級的 `共通規則 (Shared Rules)`，以提高整體團隊的維護效率與協作流暢度。

## 核心編排模式 (Core Orchestration Patterns)

### 1. 協調者模式 (Orchestrator Pattern)
- 結構：一個 `協調者代理 (Orchestrator Agent)` 負責拆解任務、指派工作、並彙整最終結果。
- 適用情境：複雜、需動態規劃或多步驟分析的任務。

### 2. 流水線模式 (Pipeline Pattern)
- 結構：多個代理以固定順序傳遞任務（例如：Scraper -> Validator -> Analyst）。
- 適用情境：步驟明確、工序標準化且依賴前一站產出的任務。

## 專案層級共享設定 (Project-Level Shared Config)
- 將編碼風格、特定名詞術語表、統一輸出語言限制及可用工具清單等共通規則提取出來。
- 避免在每個代理的 `系統提示 (System Prompt)` 中重複貼上相同的共通限制。

---

## 輸出格式範本 (Output Format Template)

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
