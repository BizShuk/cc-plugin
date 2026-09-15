---
name: planner
description: >
    Use when planning a feature before implementation, or reviewing an existing
    codebase for structural or commercial defects. Covers system architecture
    (placement, boundaries, interfaces, data flow, incremental landing steps),
    code quality review (folder layout, dependency hygiene, structural
    defects), business value review (friction, gaps, leverage, measurability),
    and business value planning (asset inventory, RICE scoring, business model,
    MVP validation).
    Triggers on: "system design", "architecture plan", "reduce coupling",
    "modularization", "review code", "folder structure", "dependency audit",
    "unused packages", "business value", "business
    improvement", "monetization", "hidden value", "MVP", "RICE",
    "系統架構", "架構規劃", "商業價值", "業務規劃", "業務改善", "商業模式設計",
    "planner".
version: "1.0.0"
allowed-tools: Read, Write, Edit, Bash, Glob, Grep
user-invocable: true
disable-model-invocation: false
effort: high
context: workspace
metadata:
    type: methodology
    platforms: [macos, linux]
---

# planner

規劃與審查是同一件事的兩個方向：`規劃`決定「該長成什麼樣」，`審查`檢查「現在長歪了哪裡」。
兩者共用同一組判準，因此收在同一個技能下，依`軸線 (axis)`與`方向 (direction)`選擇模式。

## 模式選擇 (Mode Selection)

| 軸線 (Axis)          | 規劃 (Plan)                 | 審查 (Review)               | 參考文件                                     |
| -------------------- | --------------------------- | --------------------------- | -------------------------------------------- |
| `系統 (System)`      | 架構位置、邊界、介面、資料流 | 目錄佈局、依賴衛生、結構缺陷 | [references/system.md](references/system.md)     |
| `商業 (Business)`    | 資產盤點、機會評分、商業模式、MVP | 摩擦點、缺口、槓桿點、可衡量性  | [references/business.md](references/business.md) |

選擇規則：

- 問題是`程式碼怎麼擺`→ 系統軸；問題是`使用者為何不買單`→ 商業軸。
- 要`新增`東西 → 規劃方向；要`檢視既有`東西 → 審查方向。
- 一次只走`一個模式`。同時需要兩軸時，先系統後商業，兩份計畫各自獨立輸出，不合併成一份。

## 共通原則 (Shared Principles)

無論走哪個模式，皆遵守下列五項判準：

- `1. Simpler (簡單化)`：專注於可讀性，避免過度設計。使用 guard clauses 減少巢狀，若實作需要繁瑣的文字解釋即代表過於複雜。
- `2. Concise (簡潔性)`：遵循 DRY 原則，移除死碼 (unused code/variables) 與冗餘邏輯。
- `3. Structural (結構化)`：嚴格遵守職責分離，確保每個模組與檔案只有單一且明確的用途。
- `4. Scalable (可擴展性)`：採用低耦合與模組化介面，讓新功能可以在不修改核心邏輯的情況下無縫擴充。
- `5. Consistent (一致性)`：遵循專案既有的設計、命名與架構慣例，確保跨檔案與模組的實作風格統一。

## 共通流程 (Shared Procedure)

1. `讀上下文`：先讀 `README.md`（業務）與 `CLAUDE.md`（技術）掌握專案範疇與慣例。
2. `選模式`：依上表決定軸線與方向，載入對應的 reference 檔案。
3. `界定範圍`：規劃方向一次只處理`一個 feature`，決定 `feature_name` (kebab-case) 並明列 out of scope。
4. `執行`：依 reference 內的步驟逐項完成。
5. `輸出`：規劃方向寫入 `plans/`；審查方向直接回報發現清單，不寫檔。

## 輸出規範 (Output Convention)

規劃結果一律寫入 `plans/YYYY-MM-DD-<topic>.md`，日期為 Asia/Taipei 本地日期，
`<topic>` 為 kebab-case 且需明確指涉本次變更：

| 模式     | 檔名                                     |
| -------- | ---------------------------------------- |
| 系統規劃 | `plans/YYYY-MM-DD-architecture-<topic>.md` |
| 商業規劃 | `plans/YYYY-MM-DD-business-<topic>.md`     |

審查方向`不寫檔`，只回報依價值與努力程度 (Value-over-Effort) 排序的發現清單。

各模式的完整步驟、檢查鏡頭與輸出樣板見
[references/system.md](references/system.md) 與 [references/business.md](references/business.md)。

## Common Mistakes

- `混用兩軸`：在一份架構計畫裡談定價假設，或在商業計畫裡畫依賴圖。兩軸各自獨立輸出。
- `未依專案慣例放置`：新增功能時另創架構分層或引入全新框架。應跟隨專案既存分層與設計模式。
- `規劃範圍過大`：一次規劃多個 feature。應一次專注於一個。
- `提案脫離既有資產`：提出與目前 codebase 無關的全新專案。每個機會都必須槓桿現有資產。
- `缺乏可回滾的落地步驟`：計畫僅描述終端狀態，卻缺乏可批次交付且安全的落地路徑。

## 相關技能 (Related Skills)

- `[[code-audit]]` 當問題是識別字命名、跨檔案矛盾或安全弱點，而非結構或邊界。
- `[[project-docs]]` 當問題是文件與程式碼不符，而非程式碼本身有缺陷。
- `[[project-docs]]` 的 `consolidate` 模式：當 `plans/` 累積過多歷史計畫需要壓縮。
