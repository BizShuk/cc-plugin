---
name: universal-generate
description: >
  Use when creating any artifact from intent — code, document, skill,
  plan, schema, diagram, dataset — especially when starting from a vague
  requirement or blank page. Triggers on: "create", "build", "write a
  new", "generate", "draft", "scaffold", "from scratch", "產生", "建立",
  "從零開始".
version: "2.0.0"
metadata:
  type: pattern
  tier: philosophy
  operator: generate
---

# universal-generate — 創生算子

`簽名 (Signature)`：意圖 (intent) → 產物 (artifact)，1 → 1

從無到有創造任何產物的通用程序。產物類型無關 (target-agnostic)：程式碼、文件、技能、計畫、schema、圖表皆適用。

## 四階段程序

### 1. 探索 (Explore)

先感知，後創造。掃描既有地貌：已存在什麼？慣例是什麼？邊界在哪？

> 陌生領域時，先用 `domain-exploration` 的漸進式收斂三階段。

### 2. 約束 (Constrain)

把意圖轉為明確約束。用 `system-laws` 透鏡表問：這個產物需要哪些法則？

- 空間：放在哪？歸誰管？命名？
- 時間：生命週期？何時建立、更新、廢棄？
- 金句：約束不是限制創造力，約束就是設計本身

### 3. 生成 (Draft)

在約束內產出最小完整版本。原則：

- 一次只生成一個產物，完成再下一個
- 遵循目標所在地的既有慣例（風格同化，不引入異質風格）
- 生成物必須自帶可驗證性（測試、範例、驗收條件）

### 4. 自審 (Gate)

生成完立即以 `universal-review` 過閘，不把未審產物交付下游。

## 目標類型適配表

| 產物類型 | 探索重點 | 生成鐵則 |
| :--- | :--- | :--- |
| 程式碼 | 既有架構、慣例、測試方式 | 測試先行；風格同化 |
| 文件 | 讀者是誰、既有文件結構 | 結論先行；一份文件一個目的 |
| 技能 (skill) | 既有技能是否已覆蓋 | description 只寫觸發條件 |
| 計畫 | 現況與目標的差距 | 每步可驗證、可回退 |
| Schema / 資料 | 領域詞彙表、既有實體 | 語意欄位必填；標注 unknowns |

## 反模式

| 反模式 | 問題 | 修正 |
| :--- | :--- | :--- |
| 跳過探索直接生成 | 重複造輪、違反慣例 | 先掃描既有地貌 |
| 一次生成一批 | 錯誤放大 N 倍 | 逐一生成、逐一過閘 |
| 生成物無驗收條件 | 無法判斷完成 | 生成時同步定義驗證 |

## 算子組合

`generate` 的輸出是其他算子的輸入：generate → review（過閘）→ consolidate（與既有物去重）→ aggregate（納入整體）。整條迴圈見 `universal-evolve`。
