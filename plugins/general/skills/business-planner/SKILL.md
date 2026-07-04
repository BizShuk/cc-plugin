---
name: business-planner
description: >
    Use when planning the business value of one feature — mining new
    unexplored business value or unlocking hidden value already inside an
    existing system: asset inventory, value mining, opportunity scoring,
    value proposition, MVP validation. One feature per run; writes
    plans/business-<feature_name>.md. Triggers on: "business expansion plan",
    "expand commercial value", "monetization", "hidden value",
    "商業價值", "業務規劃", "商業模式設計", "business-planner".
version: "2.0.0"
allowed-tools: Read, Bash, Glob, Grep, Write, WebSearch
effort: medium
context: fork
---

# business-planner

一次規劃 `一個 feature` 的商業價值：從既有系統挖掘 `尚未開發的新價值 (unexplored value)`，或解鎖 `已存在但被埋沒的價值 (hidden value)`，幫助使用者用現有資產達成更多。

> `Planning Only`：只產出計畫，不實作、不修改任何程式碼或設定。
> 唯一寫入：`<workspace>/plans/business-<feature_name>.md` 與 `<workspace>/README.todo`。

## 分工 (Division of Labor)

| 技能 (Skill)       | 角色 (Role)                              |
| ------------------ | ---------------------------------------- |
| `business-extract` | 萃取現況業務（系統現在做什麼）           |
| `business-planner` | 規劃一個 feature 的商業價值 — 本技能     |
| `system-planner`   | 規劃一個 feature 的系統架構              |

## 執行守則 (Execution Rules)

依序照做，不需要額外判斷：

1. 一次只規劃 `一個 feature`。使用者一次給多個時：只取第一個，
   其餘各以一行 `- [ ] business-<name>: pending` 記入 `README.todo`。
2. 資訊不足時不要中斷：採用明確假設繼續，寫入報告的
   `風險與假設 (Risks & Assumptions)` 章節。
3. 掃描上限：最多 Read `15` 個檔案。排除 `.git`, `node_modules`,
   `vendor`, `dist` 與產生碼。
4. 市場數字必須標註來源與查得日期（可用 `WebSearch`）；
   查不到就寫成 `假設待驗證`，不可虛構。

## 執行步驟 (Procedure)

### Step 1 — 界定 feature (Define the Feature)

1. 決定 `<feature_name>`：kebab-case，取自使用者需求的核心名詞。
2. 寫一句話目標：這個 feature 要為 `誰` 創造 `什麼價值`。
3. 列 2~3 條 `不做什麼 (out of scope)`。

產出：feature 名稱、一句話目標、範圍界線。

### Step 2 — 盤點既有資產 (Asset Inventory)

| 資產類型 (Asset Type) | 去哪裡找 (Where)                          |
| :-------------------- | :---------------------------------------- |
| 數據資產              | DB schema、log、匯出檔                    |
| 功能資產              | 已建功能、API、自動化管線                 |
| 用戶資產              | 使用者群、通路、整合的上下游              |
| 知識資產              | 演算邏輯、領域規則、文件                  |

每項資產一行：`資產` / `現在用途` / `使用程度（高/中/低）`。

產出：資產清單。

### Step 3 — 挖掘價值機會 (Value Mining)

兩個方向各列 `2~3` 個候選，共 `4~6` 個：

1. `隱藏價值 (Hidden Value)` — 已存在但未變現、未曝光的能力。提示問句：
   - 哪些數據只寫不讀？
   - 哪些功能只有內部在用？
   - 哪個副產品 (byproduct) 對別人有價值？
2. `未開發價值 (Unexplored Value)` — 核心能力延伸到新客群或新場景。提示問句：
   - 同樣的能力還有誰需要？
   - 價值鏈的上一步或下一步是什麼？

每個機會固定格式一句話：`為誰` × `解決什麼` × `槓桿哪項資產`。

產出：4~6 個候選機會。

### Step 4 — 評分與選定 (Score & Select)

用簡化 `RICE` 為每個候選評分，每項 `1~5` 分：

- `score = Reach × Impact × Confidence ÷ Effort`

| 機會 (Opportunity) | R | I | C | E | Score |
| :----------------- | - | - | - | - | ----- |

取最高分 `1` 個為主提案，第 2 名列為備案 (backup)。

產出：評分表 + 主提案。

### Step 5 — 價值主張與商業模式 (Value Proposition & Model)

1. 主提案的價值主張一段話：對誰、痛點、我們的差異。
2. 商業模式從清單擇一並說明定價假設：
   `訂閱` / `用量計費` / `授權` / `平台抽成` / `內部降本`。
3. 描述與現有系統的整合點；若需大幅架構改造，
   標註 `建議接力 system-planner`。

產出：價值主張 + 商業模式 + 整合點。

### Step 6 — MVP 與驗證 (MVP & Validation)

1. `MVP` 切片：`2 週` 內可驗證的最小功能範圍。
2. 定義 `北極星指標 (North Star Metric)` 與成功門檻（具體數字）。
3. 假設清單：每條標 `已驗證` 或 `待驗證`。

產出：MVP 範圍 + 指標 + 假設清單。

### Step 7 — 撰寫計畫 (Write Plan)

1. `mkdir -p plans`，寫入 `<workspace>/plans/business-<feature_name>.md`。
2. 在 `<workspace>/README.todo` 追加一行：
   `- [ ] business-<feature_name>: <一句話目標>`。

報告結構（每章填入對應 Step 的產出）：

```markdown
# 商業價值計畫 — <feature_name> (Business Value Plan)

## 1. 目標與範圍 (Goal & Scope)
<!-- Step 1：一句話目標 + out of scope -->

## 2. 既有資產盤點 (Asset Inventory)
<!-- Step 2：資產清單 -->

## 3. 價值機會與評分 (Opportunities & Scoring)
<!-- Step 3 + 4：候選機會、RICE 評分表、主提案與備案 -->

## 4. 價值主張與商業模式 (Value Proposition & Business Model)
<!-- Step 5：價值主張、商業模式、定價假設、整合點 -->

## 5. MVP 與驗證計畫 (MVP & Validation)
<!-- Step 6：MVP 範圍與驗證方式 -->

## 6. 成功指標 (Success Metrics)
<!-- Step 6：北極星指標與成功門檻 -->

## 7. 風險與假設 (Risks & Assumptions)
<!-- 執行中所有假設，各標已驗證/待驗證 -->
```

## 規則 (Rules)

- `僅規劃`：只產出計畫；唯一輸出是 `plans/` 報告與 `README.todo` 一行。
- 章節標題用繁體中文加英文括號；術語附英文與圓括號。
- 不使用粗體語法，一律以 `backtick` 強調。
- 圖表一律 Mermaid；邊線文字 (edge text) 必須用雙引號包覆。
- 每個機會都必須連回某項既有資產，且通過 `RICE` 評分排序。
- 外部市場數據必須標註來源與查得日期，不可虛構數字。

## 常見錯誤 (Common Mistakes)

| 錯誤 (Mistake)                   | 修正 (Fix)                             |
| -------------------------------- | -------------------------------------- |
| 一次規劃多個 feature             | 只取第一個，其餘記入 `README.todo`     |
| 提案與既有資產脫節成全新專案     | 每個機會必須槓桿 Step 2 的某項資產     |
| 端出多個機會卻不收斂             | `RICE` 評分後只留 1 主提案 + 1 備案    |
| 提案無法衡量成效                 | 北極星指標必須是具體數字門檻           |
| 引用市場數字卻無來源             | 標註來源與日期，或改寫為假設待驗證     |

## 完成前自檢 (Final Checklist)

- [ ] 只規劃了一個 feature
- [ ] 檔案位於 `plans/business-<feature_name>.md`
- [ ] 每個機會都連回一項既有資產
- [ ] 只有 1 個主提案（加 1 個備案）
- [ ] 指標有具體數字門檻；市場數字有來源或標為假設
- [ ] 沒有粗體；Mermaid 邊線文字有雙引號
- [ ] `README.todo` 已追加一行
