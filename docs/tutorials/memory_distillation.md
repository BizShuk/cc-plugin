# 記憶蒸餾管道教學 (Memory Distillation Pipeline Tutorial)

本教學將引導您逐步了解 cc-plugin 中的記憶蒸餾機制，並說明如何運行它。

---

## 術語解釋 (Terminology)

在開始之前，請先理解以下本專案的核心概念：

- `觀察值 (Observation)`：來自不同 AI 記憶來源的原始紀錄。
- `候選記憶 (Candidate)`：由本地語言模型從觀察值中初步提取出的潛在事實或經驗。
- `事實 (Fact)`：通過真實性門檻，確認為真且寫入 `mempalace` 儲存庫的資訊。
- `狀態儲存器 (StateStore)`：使用 SQLite 記錄已處理的觀察值遊標與狀態，防止重複蒸餾。
- `蒸餾器 (Distiller)`：負責協調讀取、LLM 提取、分流寫入與狀態更新的完整管道。

---

## 學習步驟 (Step-by-Step Guide)

### 步驟 1 (Step 1): 啟動本地語言模型服務

記憶蒸餾依賴本地的 `Ollama` 進行事實提取。請確保您的本地 Ollama 服務已啟動，且已下載預設的模型（例如 `qwen3:14b-q4_K_M`）。

您可以透過以下指令確認 Ollama 服務是否正常運作：
Run: `curl http://localhost:11434/api/tags`
Expected: 回傳 JSON 格式的本地模型清單。

### 步驟 2 (Step 2): 執行蒸餾管道

在專案根目錄下，使用 cc-plugin CLI 執行 `distill` 命令。這會自動觸發整個蒸餾流程（讀取 `gbrain` 與 `claude-mem` 的觀察值、透過 Ollama 提取，並分流寫入）：
Run: `cc-plugin distill`
Expected: 終端機輸出顯示讀取的觀察值數量、提取的候選記憶，以及成功寫入的狀態。

### 步驟 3 (Step 3): 檢查狀態儲存器

蒸餾完成後，您可以檢查 SQLite 資料庫中的 `StateStore` 遊標以確認狀態已更新。
Run: `cc-plugin status` (或直接查詢 sqlite 狀態)
Expected: 顯示當前處理的遊標位置已向前推進。
