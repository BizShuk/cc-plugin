---
name: business-planner
description: >
    Use when reviewing a codebase for business-value improvements (gaps, friction, missed revenue/retention) or when planning and expanding the commercial value of a feature (asset inventory, opportunity scoring, MVP design).
    Triggers on: "business value", "business improvement", "how to improve this", "業務改善", "增加價值", "business expansion plan", "expand commercial value", "monetization", "hidden value", "商業價值", "業務規劃", "商業模式設計", "business-planner".
version: "1.0.0"
allowed-tools: read_file, write_file, search_web
user-invocable: true
disable-model-invocation: false
effort: high
context: workspace
metadata:
    type: methodology
    platforms: [macos, linux]
---

# business-planner

## Overview

本技能整合了 `商業價值審查 (Business Value Review)` 與 `商業價值規劃 (Business Value Planning)`，旨在幫助診斷既存程式碼的商業缺陷，並規劃新功能的商業變現模式。

## When to Use

當您面臨以下情境時，應使用本技能：

- 審查功能流程，找出使用者摩擦點 (friction)、功能缺口 (gaps) 或未被衡量的成效時。
- 為特定功能 (feature) 進行商業定位、盤點系統資產、挖掘價值機會並規劃 MVP 與驗證指標時。

## Mode 1 — 商業價值審查 (Review Mode)

尋找產品未達成功商業目標的痛點，並提出依價值排序的具體修改建議。

### 審查步驟 (Procedure)

1. `理解目標`：閱讀 `README.md` 與 `CLAUDE.md` 以掌握專案業務範疇。
2. `分析核心流程`：追蹤用戶註冊、激活、核心操作與留存的完整路徑，尋找系統可以代勞的動作。
3. `評估發現`：依下方鏡頭清單評估問題，並依價值與努力程度 (Value-over-Effort) 排序。
4. `產出報告`：列出具體發現與推薦修復路徑。

### 審查鏡頭 (Review Lenses)

- `摩擦點 (Friction)`：多餘步驟、人工輸入或是不必要的等待。
- `缺口 (Gap)`：產品承諾了某項價值但流程中並未交付。
- `槓桿點 (Leverage)`：一項小改動即可顯著提升營收或留存。
- `風險點 (Risk)`：可能默默導致資料丟失、資金損耗或信任度下降的邊界狀況。
- `冗餘 (Redundancy)`：功能重複或是互相衝突的設計。
- `可衡量性 (Measurability)`：業務流程已在運作卻沒有被追蹤或衡量。

---

## Mode 2 — 商業價值規劃 (Planning Mode)

一次規劃 `一個 feature` 的商業模型。

### 規劃步驟 (Procedure)

1. `界定功能`：決定 `feature_name` (kebab-case)，說明該功能為誰創造何種價值，並列出 out of scope 項目。
2. `盤點資產`：找出專案既有的數據資產 (DB schema 等)、功能資產 (API 等)、用戶資產 (上下游整合) 與知識資產。
3. `挖掘機會`：找出 `隱藏價值 (Hidden Value)` (如未曝光數據) 與 `未開發價值 (Unexplored Value)` (延伸至新客群)。
4. `RICE 評分`：用簡化公式 `Score = Reach × Impact × Confidence ÷ Effort` 評分，決定 1 個主提案與 1 個備案。
5. `商業模式`：選擇 `訂閱` / `用量計費` / `授權` / `平台抽成` / `內部降本` 之一並說明定價假設。
6. `MVP 驗證`：規劃 2 週內可交付的 MVP 範圍、北極星指標 (具體數字門檻) 與假設驗證清單。
7. `寫入計畫`：輸出 `plans/business-<feature_name>.md` 並在 `README.todo` 追加待辦行。

---

## Output Examples

規劃結果寫入 `plans/<YYYY-MM-DD>-business-<feature_name>.md`

### Planning Mode Output Structure (plans/business-\*.md)

```markdown
# 商業價值計畫 — <feature_name> (Business Value Plan)

## 1. 目標與範圍 (Goal & Scope)

<!-- 一句話目標與 out of scope 項目 -->

## 2. 既有資產盤點 (Asset Inventory)

<!-- 數據、功能、用戶與知識資產清單 -->

## 3. 價值機會與評分 (Opportunities & Scoring)

<!-- RICE 評分表與選定的提案 -->

## 4. 價值主張與商業模式 (Value Proposition & Business Model)

<!-- 價值主張、商業模式、整合點 -->

## 5. MVP 與驗證計畫 (MVP & Validation)

<!-- 兩週 MVP 範圍與北極星指標 -->
```

## Common Mistakes

- `規劃範圍過大`：一次規劃多個功能。應一次專注於一個 feature。
- `缺乏數據來源`：引用市場數字卻未載明出處。無數據時應列為待驗證假設。
- `提案脫離既有資產`：提出與目前 codebase 無關的全新專案。每個機會都必須槓桿現有資產。
