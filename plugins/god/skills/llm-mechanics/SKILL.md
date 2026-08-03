---
name: llm-mechanics
description: >
  Use when deciding how much procedural detail to put into a prompt for a
  thinking-mode (extended reasoning) model, or when a highly-specified
  prompt is producing worse output than a looser one. Triggers on:
  "thinking mode prompt", "how much detail should the prompt have",
  "over-specified prompt", "reasoning model prompting".
version: "1.1.0"
metadata:
  type: reference
  tier: philosophy
---

# llm-mechanics

給具備延伸推理能力（Thinking Mode）的模型寫 prompt 時，`指定什麼`與`授權什麼`的分界。

## 前提

模型接收 prompt 後會自行展開推理鏈再產出答案。Prompt 寫得越細，被推理鏈佔用的空間越少 — 這是`交換`，不是免費的精確度。因此問題不是「要不要寫細」，而是「哪些部分值得用推理空間去換」。

## 分界表 (What to Specify vs What to Delegate)

| 應提供（模型推不出來的） | 應授權給模型（模型推得出來的） |
| :--- | :--- |
| 宏觀戰略目標 | 微觀邏輯推演 |
| 明確的系統邊界 | 極端案例處理 |
| 品質標準與驗收條件 | 實作路徑選擇 |

判準：這項資訊`只有你知道`（專案慣例、業務約束、驗收門檻）就寫進 prompt；模型`從目標＋邊界推導得出`（步驟順序、演算法選擇、edge case 列舉）就留空。

> 對 Thinking Mode 模型過度指定步驟，等於用你的推理取代它的推理 — 若你的推理較弱，輸出品質隨之下降。

## 應用時機

- 一份 prompt 寫得越詳細、輸出反而越差時，先刪步驟再看結果
- 為高階模型撰寫策略性 prompt（給目標與驗收，不給流程）
- 檢查既有 prompt：右欄的東西出現在 prompt 裡就是候選刪除項

## 邊界

本 skill 只涵蓋 prompt 的`詳略取捨`。模型內部機制（attention、embedding、tokenization）的正確描述請查該領域一手資料，不要從本 skill 推論。
