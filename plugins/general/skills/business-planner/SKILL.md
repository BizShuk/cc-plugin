---
name: business-planner
description: >
    Use when analyzing the business model and functional scope of an existing
    system and planning how to expand its commercial value — core-value
    assessment, horizontal/vertical opportunities, market-fit, opportunity
    prioritization (RICE), business model and MVP design. Triggers on:
    "business expansion plan", "expand commercial value", "monetization",
    "go-to-market", "商業價值擴充", "業務規劃", "商業模式設計", "business-planner".
version: "1.1.0"
allowed-tools: Read, Bash, Glob, Grep, Write, WebSearch
effort: high
context: fork
---

# business-planner

指導 AI 代理 (AI Agent) 評估現有系統的核心價值，系統化盤點水平與垂直擴充機會，依 `投資報酬 (ROI)` 排序後，產出一份可執行、可衡量的商業價值擴充提案報告。

> `Planning Only`：本技能只產出規劃報告，不實作、不修改任何程式碼或設定，
> 唯一寫入的檔案是 `${cwd}/plans/` 下的提案報告。

## 概述 (Overview)

輸入可以是 folder、repo、單一檔案、既有業務文件或貼上的描述。輸出永遠回答：
這個系統「下一步該往哪裡長、為什麼、怎麼長、怎麼驗證」。

與相鄰技能的分工：

| 技能 (Skill)       | 角色 (Role)                                |
| ------------------ | ------------------------------------------ |
| `business-extract` | 萃取現況業務（系統現在做什麼）             |
| `business-planner` | 規劃業務未來（系統下一步該做什麼）— 本技能 |
| `system-planner`   | 確保架構撐得起擴充（系統有沒有能力長）     |

建議鏈路：`business-extract` → `business-planner` → `system-planner`。

## 使用時機 (When to Use)

- 使用者要求評估專案的商業發展潛力、變現方向或業務擴充路線。
- 專案需開發新業務子系統、進入新市場或設計新的收費模式。
- 不適用：僅需代碼重構或架構簡化時，改用 `system-planner`；
  僅需釐清現況業務時，改用 `business-extract`。

## 執行步驟 (Procedure)

### Step 1 — 蒐集脈絡與界定範圍 (Context Gathering & Scope)

| 輸入型態    | 處理方式                                                      |
| :---------- | :------------------------------------------------------------ |
| folder/repo | Glob 頂層結構，鎖定 entry points、handler/service/cmd、設定檔 |
| 單一檔案    | 直接 Read，必要時追蹤直接相依                                 |
| 文件/純文字 | 直接分析，不掃描檔案系統                                      |

排除噪音：`.git`, `node_modules`, `vendor`, `dist`, 產生碼。
若現有核心價值不明確，先呼叫 `business-extract` 取得業務基線再繼續。

### Step 2 — 評估現有核心價值 (Core Value Assessment)

1. 識別 `核心業務 (Core Business)` 與主要價值/獲利來源。
2. 盤點可被槓桿的關鍵資產：獨特數據、活躍使用者、網路效應、品牌、
   專利演算邏輯、已整合的上下游通路。
3. 標記現有業務的侷限性、瓶頸與尚未變現的能力。

### Step 3 — 水平與垂直擴展機會 (Horizontal & Vertical Expansion)

1. `水平擴展 (Horizontal Expansion)`：將核心能力套用到相似領域、產業或
   客群；評估跨平台、多租戶 (Multi-tenant)、白牌 (White-label) 可行性。
2. `垂直擴展 (Vertical Expansion)`：沿價值鏈上下游延伸（如資料 → 分析 →
   決策建議），為既有客戶提供更深度加值服務。
3. 每個機會以一句話描述：`為誰` × `解決什麼` × `槓桿哪項既有資產`。

### Step 4 — 用戶痛點與市場契合度 (Pain Point & Market Fit)

1. 蒐集目標客群痛點，聚焦現有系統「差一步就能滿足」的需求。
2. 競品分析 (Competitive Analysis)：必要時用 `WebSearch` 蒐集外部佐證，
   找出差異化定位 (Differentiation)，並標註資料來源與日期。
3. 評估每個機會的市場規模 (TAM/SAM/SOM) 量級與進入門檻。

### Step 5 — 機會優先排序 (Opportunity Prioritization)

以 `RICE` 為每個機會評分並排序，避免憑直覺挑選：

- `Reach`（影響範圍）× `Impact`（單位影響）× `Confidence`（信心）÷ `Effort`（投入）

低投入高影響者優先；高投入需附階段拆解。以 `quadrantChart` 視覺化：

```mermaid
quadrantChart
    title 機會優先矩陣 (Opportunity Matrix)
    x-axis "低投入 Low Effort" --> "高投入 High Effort"
    y-axis "低影響 Low Impact" --> "高影響 High Impact"
    quadrant-1 "重點投入 Bet"
    quadrant-2 "速贏 Quick Win"
    quadrant-3 "暫緩 Defer"
    quadrant-4 "謹慎評估 Question"
    "機會 A": [0.3, 0.8]
    "機會 B": [0.7, 0.7]
```

### Step 6 — 商業模式與 MVP 提案 (Business Model & MVP)

1. 為入選機會設計 `商業模式 (Business Model)`：訂閱、使用量計費、加值
   授權、平台抽成等，並說明定價假設與單位經濟 (Unit Economics)。
2. 定義 `最小可行性產品 (Minimum Viable Product, MVP)`：最小範圍、最快
   驗證的功能切片與其驗證假設。
3. 描繪與現有系統的整合點，標示需 `system-planner` 介入的架構風險。
4. 為每個提案定義 `北極星指標 (North Star Metric)` 與成功門檻。

### Step 7 — 撰寫提案報告 (Write Report)

`<feature_name>` 取與擴充計畫相關之名稱，或目標路徑最後一段（即專案名稱）。報告寫入
`./plans/business-<feature_name>.md`（目錄不存在先 `mkdir -p plans`）；
純文字輸入且未指定路徑時，輸出於對話並詢問是否落檔。結構如下：

```markdown
# 業務擴充提案 — <feature_name> (Business Expansion Proposal)

## 1. 現有核心價值分析 (Current Core Value Analysis)

## 2. 水平與垂直擴展機會 (Horizontal & Vertical Opportunities)

## 3. 市場痛點與契合度 (Market Pain Points & Fit)

## 4. 機會優先矩陣 (Opportunity Prioritization)

## 5. 擴充功能提案與商業模式 (Proposed Features & Business Model)

## 6. 成功指標 (Success Metrics)

## 7. 實作路徑與階段規劃 (Implementation Roadmap)

## 8. 風險與假設 (Risks & Assumptions)
```

## 規則 (Rules)

- `僅規劃`：只產出報告，不實作、不改碼、不改設定；唯一輸出是 `./plans/` 下的報告。
- 章節標題用繁體中文加英文括號；內文使用繁體中文，術語附英文與圓括號。
- 不使用粗體語法，一律以 `backtick` 強調。
- 圖表一律 Mermaid；邊線文字 (edge text) 必須用雙引號包覆。
- 每個機會都必須連回某項既有資產與某個明確痛點，且通過 `RICE` 排序。
- 外部市場數據必須標註來源與查得日期，不可虛構數字。

## 常見錯誤 (Common Mistakes)

| 錯誤                             | 修正                                   |
| -------------------------------- | -------------------------------------- |
| 擴充提案與核心價值脫節成全新專案 | 每個機會回扣既有資產與既有客群         |
| 純技術視角構想功能，缺市場痛點   | 先做 Step 4 痛點/契合度再提功能        |
| 一次端出大量機會卻不排序         | 強制 `RICE` 評分並用矩陣呈現           |
| 提案無法衡量成效                 | 每個提案定義北極星指標與成功門檻       |
| 引用市場數字卻無來源             | 標註來源與日期，或改寫為「假設待驗證」 |

## 失敗模式 (Failure Modes)

| 情境                 | 動作                                       |
| -------------------- | ------------------------------------------ |
| 現有核心價值不明確   | 先以 `business-extract` 提煉業務本質再規劃 |
| 輸入過大無法全讀     | 鎖定 entry points 與領域模型，註明部分掃描 |
| 擴充高度依賴架構改造 | 在報告標示並建議接力 `system-planner`      |
| 缺乏市場數據         | 以明確假設替代，列入「風險與假設」待驗證   |
