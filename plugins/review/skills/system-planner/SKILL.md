---
name: system-planner
description: >
    Use when planning a system architecture for a feature (placement, boundaries, interface, data flow, incremental landing steps) or when performing codebase quality reviews (code standards, directory layout, cross-file consistency, dependency hygiene).
    Triggers on: "system design", "architecture plan", "reduce coupling", "modularization", "系統架構", "架構規劃", "system-planner", "review code", "folder structure", "check consistency", "dependency audit", "unused packages", "are these deps needed".
version: "2.0.0"
allowed-tools: read_file, write_file, search_web
user-invocable: true
disable-model-invocation: false
effort: high
context: workspace
metadata:
    type: methodology
    platforms: [macos, linux]
---

# system-planner

## Overview

本技能整合了 `系統架構規劃 (System Architecture Planning)` 與 `系統與程式碼品質審查 (System & Code Quality Review)`，提供全方位的系統設計指引與靜態品質診斷規範。

- `1. Simpler (簡單化)`：專注於可讀性，避免過度設計。使用 guard clauses 減少巢狀，若實作需要繁瑣的文字解釋即代表過於複雜。
- `2. Concise (簡潔性)`：遵循 DRY 原則，移除死碼 (unused code/variables) 與冗餘邏輯。
- `3. Structural (結構化)`：嚴格遵守職責分離，確保每個模組與檔案只有單一且明確的用途。
- `4. Scalable (可擴展性)`：採用低耦合與模組化介面，讓新功能可以在不修改核心邏輯的情況下無縫擴充。
- `5. Consistent (一致性)`：遵循專案既有的設計、命名與架構慣例，確保跨檔案與模組的實作風格統一。

## When to Use

當您面臨以下情境時，應使用本技能：

- 規劃新功能 (feature) 的架構位置、邊界定義、資料流與漸進落地步驟
- 審查程式碼實作、檔案與資料夾路徑佈局、跨檔案邏輯一致性或第三方相依套件健康狀態時。

---

## Mode 1 — 系統架構規劃 (Architecture Planning Mode)

### 規劃步驟 (Procedure)

1. `界定功能`：決定 `feature_name` (kebab-case)，說明該功能為誰服務，並明列 out of scope 範圍。
2. `盤點現況`：透過尋找進入點 (handlers 等)、改動熱點 (git log) 畫出現況架構圖 (Mermaid flowchart，節點 ≤ 10)。
3. `位置與邊界`：依專案慣例放置新功能（如 Handler / Service / Repository），明訂單向依賴關係與資料邊界。
4. `介面與資料流`：定義跨模組互動的介面 (API Contract，數量 ≤ 5) 並用 Mermaid 畫出資料流。
5. `清晰與可擴充性檢查`：逐項檢查單一職責、依賴方向、可替換性、水平擴充以及擴充點。
6. `漸進落地步驟`：拆解為 3~7 個可獨立交付與回滾的實作步驟。

---

## Mode 2 — 系統與程式碼品質審查 (Quality Review Mode)

對 codebase 或變更進行多維度靜態審查，回報真實缺陷與優化建議。

### 2.1 跨檔案一致性審查 (Consistency)

檢驗修改是否與工作區其他地方產生矛盾。

- `業務規則`：確保沒有兩個地方對同一事件執行互斥的邏輯。
- `領域模型`：相同實體在不同檔案中的欄位與型別必須一致。
- `數據合約`：生產端 (producer) 的輸出格式必須符合消費端 (consumer) 的預期。
- `雙向參考`：實體連結、呼叫者/被呼叫者合約等必須雙向成立，單向引用視為不一致。

### 2.2 外部套件衛生審查 (Dependency Hygiene)

檢驗第三方相依套件，維持精簡且安全的依賴狀態。

- `未動用`：已宣告但在原始碼中完全沒有 import 或呼叫的套件。
- `重複或冗餘`：引入多個功能相近的套件，或該功能實則可用標準函式庫代替。
- `版本未釘定`：使用萬用字元或範圍版本限制，而非明確的版本號。
- `重度依賴輕度使用`：為了一個簡單的小函式而拉入整個巨大的套件。

### 2.3 目錄與檔案佈局審查 (Folder Structure)

檢驗實體目錄配置是否清晰且符合專案慣例。

- `職責混雜`：單一資料夾內含有多個不相干職責的檔案。
- `分層錯誤`：檔案被放置在其不該擁有的架構分層中。
- `孤立資料夾`：沒有任何檔案引用且無說明文檔提及的資料夾。
- `散落檔案`：將本應放入子套件的檔案直接散落在根目錄或頂層目錄。

## Output Examples

規劃結果寫入 `plans/architecture-<feature_name>.md`

### Output Structure

```markdown
# 架構計畫 — <feature_name>

## 1. 目標與範圍 (Goal & Scope)

<!-- 目標與 out of scope 項目 -->

## 2. 現況架構 (Current Architecture)

<!-- 現況 Mermaid 架構圖與模組清單 -->

## 3. 架構位置與邊界 (Placement & Boundaries)

<!-- 分層放置位置與依賴限制 -->

## 4. 介面與資料流 (Interfaces & Data Flow)

<!-- 介面定義表與資料流圖 -->

## 5. 清晰與可擴充性檢查 (Clarity & Scalability Check)

<!-- 檢查結果 -->

## 6. 漸進落地步驟 (Incremental Steps)

<!-- 3~7 步具體步驟與驗證/回滾方式 -->
```

## Common Mistakes

- `未依專案慣例放置`：新增功能時另創架構分層或引入全新框架。應跟隨專案既存分層與設計模式。
- `規劃過多介面`：規劃階段設計了超過 5 個介面，使系統複雜化。介面應保持精簡，必要時予以合併。
- `缺乏可回滾的落地步驟`：架構計畫僅描述終端狀態，卻缺乏可批次交付且安全的落地路徑。
