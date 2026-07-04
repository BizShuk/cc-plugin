# 建立學習文件技能 (Learning Document Skill) 實作計畫 (Learning Document Skill Implementation Plan)

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

`Goal:` 在 `plugins/review/skills/` 底下建立 `learning-document` 技能，用以指導 AI 代理為當前專案生成逐步的學習文件，並分類存放於 `./docs/` 與 `./docs/tutorials/` 目錄中，且對關鍵術語進行特別高亮與解釋。

`Architecture:` 該技能符合 `agentskills.io` 規範。它將作為一個靜態指導文件 `SKILL.md`，並將其註冊到 `plugins/review/.claude-plugin/plugin.json`。我們將透過模擬測試腳本驗證 AI Agent 讀取該技能後，是否能正確分類並高亮術語。

`Tech Stack:` Markdown, `agentskills.io` frontmatter, Bash

---

### Task 1: 建立學習文件技能 (Learning Document Skill)

`Files:`
- Create: `plugins/review/skills/learning-document/SKILL.md`

- [ ] **Step 1: 建立技能檔案**

建立 `../plugins/review/skills/learning-document/SKILL.md`，內容須包含 `Use when...` 起頭的 YAML frontmatter 觸發條件描述，且以 Traditional Chinese 撰寫技能指引。

內容如下：
```markdown
---
name: learning-document
description: >
  Use when the user requests step-by-step guides, onboarding documents, or learning materials to understand the project workspace. Triggers on "create learning document", "generate tutorials", "學習文件", "學習指南".
version: "1.0.0"
allowed-tools: Read, Write, Glob
user-invocable: true
disable-model-invocation: false
effort: medium
metadata:
  type: reference
  platforms: [macos, linux]
---

# 學習文件 (Learning Document)

## 概述 (Overview)
本技能用以引導 AI 代理針對當前工作區 (Workspace) 產生結構化的逐步學習文件，協助讀者快速理解專案的領域知識與架構。

## 使用時機 (When to Use)
- 當使用者需要專案導覽、架構引導或新手入門文件時。
- 當需要建立特定技術模組的教學步驟時。

## 核心規範 (Core Specification)

1. `輸出路徑分流 (Output Directory Routing):`
   - 領域學習與概念理解文件：必須存放在 `./docs/tutorials/` 目錄。
   - 專案架構、開發指南與其他常規學習文件：存放在 `./docs/` 目錄。

2. `關鍵術語高亮與解釋 (Key Terms Highlight & Explanation):`
   - 對於工作區內的重要術語（如 `PSM`, `StateStore`, `Distiller`, `Observation` 等），在文件首次出現時必須使用 `backtick` 進行包裹高亮。
   - 必須提供專門的 `術語解釋` 區段或以側邊註解方式，詳細說明該術語在專案中的具體含意與業務邏輯。

3. `逐步引導 (Step-by-Step Guide):`
   - 教學或學習文件必須使用循序漸進的步驟格式（例如 `Step 1`, `Step 2`）。
   - 每個步驟必須有明確的目標、程式碼/設定範例與預期結果。
```

- [ ] **Step 2: 驗證檔案已正確建立**

確認檔案存在且內容無誤。

---

### Task 2: 註冊技能至插件定義 (Plugin Manifest)

`Files:`
- Modify: `plugins/review/.claude-plugin/plugin.json`

- [ ] **Step 1: 修改 plugin.json 註冊新技能**

修改 `../plugins/review/.claude-plugin/plugin.json`，將新技能的相對路徑 `./skills/learning-document` 加入 `skills` 陣列中。

修改處：
```json
  "skills": [
    "./skills/consistency",
    "./skills/business-improvement",
    "./skills/folder-structure",
    "./skills/naming-convention",
    "./skills/doc-sync",
    "./skills/dependency-hygiene",
    "./skills/test-coverage",
    "./skills/learning-document"
  ]
```

- [ ] **Step 2: 更新 plugins/review/README.md**

根據 CLAUDE.md 中「更新插件時也必須同步更新對應的 README.md」的規則，修改 `../plugins/review/README.md`，在技能列表中加入 `learning-document` 的描述。

修改處：
在目錄結構與技能對照表中加入 `learning-document` 相關內容。

---

### Task 3: 驗證技能安裝與運作 (Verify Skill Installation and Output)

`Files:`
- Modify: none (僅執行驗證命令與測試)

- [ ] **Step 1: 重新連結並載入技能**

於專案根目錄下執行安裝命令，將新技能註冊至 Agent：
Run: `npx skills add .`
Expected: 輸出中應包含 `learning-document` 技能的安裝與連結成功訊息。

- [ ] **Step 2: 模擬生成測試**

透過模擬指令或測試，請求 Agent 使用該技能為本專案的 `StateStore` 或是 `Memory Distillation` 業務生成一份學習指南。
Expected:
1. 在 `./docs/` 或 `./docs/tutorials/` 產生相應的 `.md` 文件。
2. 領域學習相關內容放置在 `./docs/tutorials/` 下。
3. 文件中有對 `StateStore` 或 `Distillation` 等術語使用 `backtick` 高亮，並附帶專屬的術語解釋區段。
