---
name: domain-exploration
description: >
  Use when entering an unfamiliar domain, technology, or problem space and
  needing a structured approach to reduce uncertainty. Triggers on: "new
  territory", "unknown domain", "feasibility study", "PoC", "spike",
  "exploration", "where do I even start".
version: "1.0.0"
metadata:
  type: pattern
  tier: philosophy
---

# domain-exploration

面對未涉及的新領域，採用`漸進式收斂 (Progressive Convergence)` 策略，從廣域掃描逐步收斂至精確架構。

## 三階段收斂法

### 階段一：地貌掃描 (Domain Mapping)

釐清「是什麼」與「不是什麼」，建立粗略的概念地圖。

- 識別核心實體與關係
- 標記已知邊界與未知區域
- 產出：概念地圖 + 疑問清單

### 階段二：插旗測試 (PoC Testing)

透過最小可行性原型主動「撞牆」，測量真實的物理與邏輯邊界。

- 針對最大不確定性優先測試
- 每次測試只驗證一個假設
- 產出：邊界報告 + 限制清單

### 階段三：結構化收斂 (Structural Convergence)

依據碰撞出的限制，確立精確的依賴關係與系統架構。

- 將限制轉化為架構約束
- 建立依賴圖與介面定義
- 產出：架構設計 + 實作計畫

## 反模式

| 反模式               | 問題                         | 修正                     |
| :------------------- | :--------------------------- | :----------------------- |
| 直接跳到階段三       | 基於假設而非事實設計         | 先掃描 + 測試            |
| 階段二過度投入       | PoC 膨脹為產品               | 限定每次測試範圍與時間   |
| 跳過階段一           | 不知道自己不知道什麼         | 先建概念地圖             |

## 適用場景

- 評估全新技術棧的可行性
- 進入陌生業務領域
- 面對模糊需求的專案啟動
- 任何「不知道從哪裡開始」的情境
