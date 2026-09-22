---
name: tutorial
description: >
    Use when generating tutorials, step-by-step guides, domain knowledge documents, onboarding materials, or learning resources for the workspace. Triggers on "create learning document", "generate tutorials", "學習文件", "學習指南".
version: "1.2.0"
allowed-tools: Read, Write, Glob
user-invocable: true
disable-model-invocation: false
effort: high
metadata:
    type: reference
    platforms: [macos, linux]
---

# 教程 (Tutorial)

## 概述 (Overview)

本技能用以指導 AI 代理針對當前工作區 (Workspace) 產生逐步的學習與教學文件，確保產出的文件結構清晰，並正確分類存放。

## 使用時機 (When to Use)

- 使用者要求建立逐步引導教學 (step-by-step tutorials)。
- 建立引導新進成員理解專案的入門文件 (onboarding guide)。
- 當需要解釋專案核心業務領域知識與技術架構時。

## 核心規範 (Core Specification)

### 1. 輸出路徑分流 (Output Routing)

為了維護文件目錄的乾淨，產出檔案必須根據內容類型進行分流：

- 領域知識學習與概念導覽 (Domain tutorials)：必須存放在 `./docs/tutorials/` 目錄。
- 專案架構、開發流程、環境設定等常規文件 (Project guides)：存放在 `./docs/` 根目錄下。

`docs/tutorials/` 是統一介面 (Unified Interface) 的選備項目，規範見
`~/.claude/CLAUDE.md`。專案可能位於 `~/projects/<project>/` 或
`~/projects/<category>/<project>/`，路徑一律相對於`專案根目錄`，不寫進分類目錄。

### 2. 關鍵術語高亮與解釋 (Key Terminology)

專案特有的術語在文件中首次出現時，必須給予特別的高亮與說明：

- 必須使用 `backtick` 包裹該術語。
- 文件中必須有專屬的 `術語解釋 (Terminology)` 區段，用以說明該術語在當前專案中的意義與上下文（例如 `PSM`, `StateStore`, `Distiller`, `Observation`）。
- 定義以 `docs/terminology.md`（術語單一定義來源）為準：已收錄者`直接引用`不得改寫，
  未收錄者先補進術語表再於教學中使用，避免同一概念兩種說法。

### 3. 循序漸進引導 (Step-by-Step Flow)

教學文件應採用步驟引導格式：

- 使用明確的標題（例如 `步驟 1 (Step 1)`, `步驟 2 (Step 2)`）。
- 每個步驟必須包含：執行目的、詳細操作或程式碼區塊、預期結果。

### 4. 章節導覽列 (Chapter Navigation)

同一系列的教學檔之間必須能前後跳轉，每個檔案放兩次`完全相同`的導覽列：

- 位置：`檔首`緊接在 H1 標題之後；`檔尾`為檔案最後一行，前面以 `---` 分隔。
- 格式：上一章與下一章寫在`同一行`，以 ` | ` 分隔，連結文字必須包含目標檔案的`章節名稱`（取自對方 H1 標題，含編號），不得只寫「上一章」「下一章」或檔名。
- 邊界：系列第一篇的上一章寫 `無 (起點)`，最後一篇的下一章寫 `無 (終點)`，仍保留同一行格式。
- 順序以章節編號或系列索引 (`README.md`) 為準；新增、刪除或更名章節時，必須同步更新前後相鄰檔案的導覽列。

```markdown
# 1.2.0 台美市場交易機制與微觀規則對比

⬅️ 上一章: [1.1.0 總體經濟與景氣循環](1.1.0-macroeconomics-and-business-cycles.md) | 下一章: [1.3.0 產業價值鏈拆解方法論](1.3.0-industry-value-chain-mapping.md) ➡️

...（本文）...

---

⬅️ 上一章: [1.1.0 總體經濟與景氣循環](1.1.0-macroeconomics-and-business-cycles.md) | 下一章: [1.3.0 產業價值鏈拆解方法論](1.3.0-industry-value-chain-mapping.md) ➡️
```

## 常見錯誤 (Common Mistakes)

- 將領域知識教學隨意放在 `./docs/` 根目錄，未歸類至 `tutorials/`。
- 術語首次出現時未以 `backtick` 高亮，且缺乏專屬的術語對照說明。
- 自行改寫 `docs/terminology.md` 既有定義，造成同一概念兩種說法。
- 步驟式說明缺乏程式碼範例或具體指令。
- 只在檔首或只在檔尾放導覽列，或上一章與下一章拆成兩行。
- 導覽連結只寫「上一章」「下一章」或檔名，未顯示章節名稱；章節更名後相鄰檔案的導覽列未同步。
