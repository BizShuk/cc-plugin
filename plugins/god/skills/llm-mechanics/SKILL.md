---
name: llm-mechanics
description: >
  Use when reasoning about how LLMs process language internally, debugging
  unexpected model outputs, crafting prompts for thinking-mode models, or
  explaining why negation/homophone queries produce surprising results.
  Triggers on: "why did the model misunderstand", "prompt engineering",
  "attention mechanism", "embedding", "negative instruction".
version: "1.0.0"
metadata:
  type: reference
  tier: philosophy
---

# llm-mechanics

語言模型（LLM）的本質是在高維度空間中進行「找方向、算距離」的矩陣運算。

## 核心機制 (Core Mechanisms)

### 方向配對 (Vector Matching)

透過`詞嵌入 (Word Embedding)` 將文字化為座標，並由`注意力機制 (Attention)` 根據上下文推導出目標向量，最終透過相似度計算找出最佳的下一個字。

- 文字 → 高維向量座標
- 上下文 → Attention 推導目標方向
- 相似度計算 → 選出最接近的候選

### 負面語意處理 (Negative Semantics)

遭遇否定詞時，模型執行`空間反轉與向量減法`，強制將目標座標 180 度轉向對立的語意維度。

> 這解釋了為何「不要做 X」有時反而強化了 X 的語意權重 — 模型必須先定位 X 再反轉。

### 同音字解析 (Homophone Resolution)

透過`位置編碼 (Positional Encoding)` 與`詞性標註 (POS Tagging)`，賦予每個重複字元絕對且獨立的空間座標，突破人類視覺與聽覺的混淆限制。

### 高階模型提示詞策略 (Advanced Prompting)

針對具備深度思考能力（Thinking Mode）的模型：

| 應提供                     | 應授權給模型               |
| :------------------------- | :------------------------- |
| 宏觀戰略目標               | 微觀邏輯推演               |
| 明確的系統邊界             | 極端案例處理               |
| 品質標準與驗收條件         | 實作路徑選擇               |

> 對 Thinking Mode 模型過度指定步驟，等於剝奪其推理空間，反而降低輸出品質。

## 應用時機

- 除錯模型的非預期輸出
- 理解為何否定指令失效
- 為高階模型撰寫策略性 prompt
- 向他人解釋 LLM 的運作原理
