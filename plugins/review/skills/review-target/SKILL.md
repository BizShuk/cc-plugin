---
name: review-target
description: >
    Use when the user shares content about themselves, a social account, a
    trip, a sales pipeline, or a business and asks what to grow and how. Maps
    the content to growth categories (social media, life, travel, sales,
    business) and their sub-categories, picks the KPI (target attribute) to
    grow in each, and breaks every KPI into concrete, tool-scoped actions —
    e.g. Instagram Reel / Story / carousel post, a pull-up progression.
    Triggers on: "how do I grow", "growth strategy", "what should I focus on",
    "set my KPI", "grow my Instagram", "get more followers", "get fitter",
    "improve my sales", "grow my business", "成長目標", "成長策略",
    "要怎麼成長", "該專注什麼", "漲粉", "提升業績", "review-target".
version: "1.0.0"
allowed-tools: Read, Glob, Grep, WebFetch, AskUserQuestion
user-invocable: true
disable-model-invocation: false
effort: high
metadata:
    type: methodology
    platforms: [macos, linux]
---

# review-target

把使用者給的內容（自述、帳號、數據截圖、行程、業績表、營運資料）對應到`成長類別`，
找出每個類別要成長的 `KPI`，再把每個 KPI 拆成指定工具與動作的行動項目。

- `類別 (Category)`：高階領域，下分`子類別 (Sub-category)`。
- `KPI`：類別下要成長的可量測屬性，例如互動率、引體向上次數。
- `目標值 (Target value)`：KPI 在期限內要達到的數字，由使用者自身的現況推算。

## 類別索引 (Category Index)

| 類別 (Category)           | 子類別 (Sub-category)                                                                                                      | 內容訊號 (Signals)                           | 參考文件                                                 |
| ------------------------- | -------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------- | -------------------------------------------------------- |
| `社群媒體 (Social Media)` | Instagram、Threads、YouTube、TikTok、LinkedIn、自有受眾 (Owned Audience)                                                   | 帳號、貼文、粉絲、觸及、觀看、名單           | [references/social-media.md](references/social-media.md) |
| `生活 (Life)`             | 體能 (Fitness)、恢復 (Recovery)、心智 (Mind)、財務 (Finance)、人際 (Relationship)、學習與職涯 (Learning & Career)          | 身體、運動、睡眠、情緒、存錢、人際、學新技能 | [references/life.md](references/life.md)                 |
| `旅行 (Travel)`           | 頻率 (Frequency)、成本 (Cost)、體驗 (Experience)、旅行能力 (Travel Skills)、紀錄 (Documentation)                           | 行程、目的地、機票、住宿、假期               | [references/travel.md](references/travel.md)             |
| `銷售 (Sales)`            | 開發 (Prospecting)、成交 (Closing)、客單價 (Deal Size)、電商與門市 (E-commerce & Retail)、客戶擴展 (Account Growth)        | 業績、客戶、報價、成交率、訂單               | [references/sales.md](references/sales.md)               |
| `事業 (Business)`         | 產品市場契合 (PMF)、獲客 (Acquisition)、啟用 (Activation)、留存 (Retention)、營收 (Revenue)、推薦 (Referral)、單位經濟 (Unit Economics)、營運 (Operations) | 產品、用戶、營收、成本、團隊 | [references/business.md](references/business.md) |

每份 reference 都有同樣的三塊：`診斷順序`（決定主攻哪個 KPI）、各子類別的 `KPI 表`
（公式、量測位置、參考區間）、`行動目錄`（工具、具體動作、份量）。

分類規則：

- 一份內容可命中多個類別，例如「用 IG 賣手作課」= 社群媒體 + 銷售，每個命中的類別各自載入 reference。
- 只讀命中的 reference。
- 沒有任何類別命中時直接說明並停止。

## 流程 (Procedure)

1. `分類`：從內容抽出訊號，對照類別索引列出命中的類別與子類別，附上判斷依據。
2. `定現況`：每個候選 KPI 寫出現況值。內容裡找不到的標 `⚠️ 待量測` 並附 reference 的量測位置，
   數字一律來自使用者。缺的現況會改變主攻判斷時，用 AskUserQuestion 一次問完。
3. `找瓶頸`：依 reference 的診斷順序，每個類別選 1-2 個主攻 KPI，所有類別合計最多 3 個。
4. `訂目標`：每個主攻 KPI 給一個 30 或 90 天期限的目標值，寫成 `現況 → 目標`。
   reference 的參考區間只用來判斷強弱，目標值從現況往上推。
5. `拆行動`：每個主攻 KPI 拆成 2-4 項行動，從行動目錄挑選，再依使用者程度調整份量。

## 行動項目規格 (Action Item Spec)

每項行動同時具備四個欄位，缺任何一個就只是方向：

| 欄位      | 要求                         | 範例                                             |
| --------- | ---------------------------- | ------------------------------------------------ |
| 工具／範圍 | 平台格式、器材、App 或場景   | `Instagram Reel`、`單槓`、`Google Flights`        |
| 具體動作  | 可直接照做的動作或規格       | 第 1 秒放成品畫面，加上 8 字以內的大字標題       |
| 頻率／份量 | 次數、組數、時長或預算       | 每週 3 支；3 組 × 3-5 下離心引體                 |
| 完成判準  | 可觀察、可量測的結果         | 3 秒留存率高於自身近 30 天中位數                 |

`多發文`、`多運動`、`提升品牌` 這類缺工具與份量的句子不得輸出。

## 輸出格式 (Output Format)

```markdown
# 成長目標審查：<一句話內容摘要>

## 命中類別

- <類別> / <子類別>：<判斷依據>

## <類別> / <子類別>

| KPI   | 現況                  | 目標 (期限)          | 角色          |
| ----- | --------------------- | -------------------- | ------------- |
| <KPI> | <數值 或 ⚠️ 待量測>   | <數值> (<30/90> 天)  | 主攻 / 觀察   |

- [ ] <工具／範圍>：<具體動作>；<頻率／份量>；完成判準：<可觀察結果>

## 本週第一步

<所有行動中最小、最先做的一件事>
```

## 相關技能 (Related Skills)

- `[[planner]]` 的商業軸：目標是替軟體功能設計商業模式與 MVP，而非既有事業的成長指標。
- `[[content-summarizer]]`：使用者給的是 URL、影片或長文件時，先摘要再分類。
