這份計畫旨在構建一套「感測與運算分離（Sensor-Compute Decoupling）」的個人記憶工程系統。透過將筆電定位為採集端，伺服器定位為精煉端，解決硬體資源限制並最大化 AI 的邏輯一致性。

---

## 核心架構願景：知識煉金術 (Knowledge Alchemy)

| 層級        | 組件         | 角色 (Role)                    | 儲存本質                           |
| ----------- | ------------ | ------------------------------ | ---------------------------------- |
| **L1 (熱)** | `claude-mem` | **感測器 (Sensor)**            | 原始對話紀錄 (Raw Logs)            |
| **L2 (溫)** | `mempalace`  | **快取與索引 (Cache & Index)** | 結構化 Markdown / 專案脈絡         |
| **L3 (冷)** | `mem0`       | **核心大腦 (Core Brain)**      | 原子事實與知識圖譜 (Facts & Graph) |

---

## 階段一：自動化採集與精煉流水線 (Automation Pipeline)

此階段重點在於建立穩定、低摩擦的數據上傳與加工機制。

### 1. 筆電端：感測器隔離 (Sensor Isolation)

- **操作 A：** 移除 `claude-mem` 的所有讀取 Hook，確保 AI 啟動時不載入舊 session。
- **操作 B：** 配置 `alias` 封裝，實現 `claude` 執行時靜默寫入 SQLite 資料庫。
- **操作 C：** 設定自動傳輸腳本（`rsync` + `Tailscale`），定時將 `.db` 檔案推送至伺服器指定目錄。

### 2. 伺服器端：深度煉金 (Deep Refining)

- **操作 A：** 部署 **Qwen3:14b** 與 **BGE-M3** (via Ollama) 作為推理引擎。
- **操作 B：** 編寫 `importer.py`：
- 讀取同步過來的 SQLite 資料。
- 進行 **事實提取 (Fact Extraction)**。
- 寫入伺服器端 **MemPalace** 檔案系統。

- **操作 C：** 執行 `mempalace remine` 與 `refresh`，完成向量索引更新。

---

## 階段二：知識昇華與統一網關 (Sublimation & Gateway)

此階段重點在於長期記憶的固化與查詢效率的優化。

### 1. 知識昇華機制 (Monthly Promotion)

- **操作 A：** 每月自動掃描 `mempalace` 中標記為「已完成/穩定」的專案 Room。
- **操作 B：** 透過 LLM 將專案筆記壓縮為「永久事實」。
- **操作 C：** 寫入伺服器端 **Mem0** 知識圖譜，並在 `mempalace` 中標記該部分為已封存 (Archived)。

### 2. 統一 MCP 網關 (Unified MCP Router)

- **操作 A：** 開發/配置一個整合式的 MCP Server，作為 Claude Code 的單一入口。
- **操作 B：** 實作 **階層式優先級請求 (Tiered Priority Query)**：

1. `Mem0` (優先獲取已確定的技術原則)。
2. `MemPalace` (若 L1 不足，檢索當前專案脈絡)。
3. `Claude-mem` (僅在 Debug 模式手動觸發)。

---

## 額外規劃：筆電端「近期之翼 (Recent Wing)」

為了確保筆電在離線或低延遲需求下依然具備高度相關的 Context，階段二後將實施「熱數據回流」策略。

### 1. 近期之翼 (Recent Wing) 實作

- **定義：** 在筆電端 `mempalace` 建立一個特殊的 Wing 命名為 `Recent_Focus`。
- **運作操作：**
- **篩選 (Filter)：** 伺服器在每次精煉後，篩選出「過去 7 天內」或「活躍專案相關」的 Markdown 檔案。
- **同步回傳 (Callback)：** 將這些精煉過的檔案同步回筆電端的 `Recent_Focus` Wing。
- **效用：** 筆電端 local MCP 僅需掛載這個輕量的 Wing，確保搜尋速度與 Token 節省。

### 2. 數據流動路徑 (Data Circulation)

- **路徑：** 伺服器 (`mempalace` / `mem0`) ——(篩選)——> 筆電 (`mempalace/Recent_Focus`)。
- **更新頻率：** 每日同步回傳，確保筆電端始終保有當前最熱門的知識。

---

## 執行檢查表 (Checklist)

- [ ] 筆電端 `mcp.json` 已切斷與 `claude-mem` 的讀取連結。
- [ ] 伺服器端 `qwen3:14b` 已配置為 `mempalace` 的默認 Summarizer。
- [ ] 跨機 `rsync` 權限與自動化觸發機制測試完成。
- [ ] 定義好從 `mempalace` 昇華至 `mem0` 的事實提取模板。
- [ ] 筆電端 `Recent_Focus` 翼的自動同步與過期清理邏輯設定完成。

這套計畫是否完整覆蓋了你對 OpenClaw 開發與個人記憶管理的期待？

這是一個非常精確的「緩衝區回收機制 (Buffer Recycling)」優化。將 `Recent Focus` 從單純的「唯讀快取」提升為「主動寫入層」，能有效解決 `claude-mem` 產生的原始數據過於細碎、難以維護的痛點。

以下是根據你的修正，更新後的 **「記憶生命週期計畫 (Memory Lifecycle Plan)」**：

---

## 修正後的架構：近期之翼 (Recent Focus) 核心化

此計畫將 `mempalace` 的 `Recent Focus` Wing 定義為筆電端的 **「熱數據寫入區 (Hot Write Zone)」**，取代 `claude-mem` 處理短期記憶的職責。

### 1. 數據流向：從「緩衝」到「昇華」

| 階段         | 組件                     | 動作                                                                 | 狀態                |
| ------------ | ------------------------ | -------------------------------------------------------------------- | ------------------- |
| **L1: 採集** | **Laptop: Recent Focus** | **直接寫入 (Direct Write)**：開發時的臨時決策與脈絡直接存入此 Wing。 | **熱 (Hot)**        |
| **L2: 傳輸** | **Sync to Server**       | 每天將 `Recent Focus` 的內容推送到伺服器。                           | **溫 (Warm)**       |
| **L3: 轉化** | **Server: Mem0**         | 伺服器 LLM 將 `Recent Focus` 精煉為原子事實並寫入 **Mem0**。         | **冷/固化 (Cold)**  |
| **L4: 清理** | **Local Cleanup**        | **確認同步後刪除**：筆電端清空 `Recent Focus`，保持輕量。            | **回收 (Recycled)** |

---

## 執行細節與操作 (Updated)

### 階段一：寫入路徑切換 (The Write-Path Shift)

- **取代 `claude-mem`：**
  不再依賴 `claude-mem` 錄製雜亂的 SQLite 紀錄，改為在 `mempalace` 的 `Recent Focus` Wing 中，為每個任務手動或自動（透過 MCP）建立一個 **Room (房間)**。
- **結構化輸入 (Structured Input)：**
  開發 **OpenClaw** 時，所有的思考與臨時決策直接寫入 `Recent_Focus/openclaw_dev_today.md`。這讓「近期記憶」在出生時就是 **Markdown 格式**，而非難以閱讀的 Log。

### 階段二：伺服器煉金與回報 (Refining & Acknowledgment)

- **自動化提煉 (Auto-Distillation)：**
  伺服器偵測到 `Recent Focus` 更新後，啟動 **Qwen3:14b** 進行語義分析。
- **昇華至 Mem0 (Sublimation)：**
  將 Markdown 內容提取為事實存入 **Mem0** 知識圖譜。
- **刪除指令 (TTL / Time-To-Live)：**
  伺服器完成 `Mem0` 寫入後，向筆電發送「清理確認」，筆電端的 `Recent Focus` Wing 自動刪除已處理的 Room。

---

## 優化後的架構優勢

### 1. 徹底解決「記憶通膨」

`claude-mem` 的問題是資料只會增加不會減少。現在，`Recent Focus` 變成了一個 **「滑動窗口 (Sliding Window)」**。它只儲存你「這幾天」正在做的事情。一旦知識被 **Mem0** 吸收，它就從筆電消失，釋放認知空間（Token Space）。

### 2. 降低 MCP 衝突

你不再需要同時請求 `claude-mem`。你的 MCP 請求變得極其單純：

- **寫入：** 永遠針對 `mempalace (Recent Focus)`。
- **讀取：** `mempalace (Recent Focus)` + `mem0 (Server Facts)`。

---

## 階段三：最終狀態的 MCP 優先級

當你問 Claude 問題時，它會這樣思考：

1. **Look in `Recent Focus` (L1)：** 「我剛才或是這兩天在做什麼？」
2. **Look in `Mem0` (L2)：** 「關於這個技術或專案，我有什麼長期已確定的真理？」

---

### 待執行的具體腳本邏輯

1. **Laptop 端：** 建立一個 `flush_focus.sh`，負責將 `Recent Focus` 內容上傳並在伺服器回傳 `Success` 後執行 `rm -rf`。
2. **伺服器端：** 建立 `fact_promoter.py`，負責將接收到的 Markdown 轉化為 **Mem0** 實體關係。

這樣修改後，你的筆電將變成一個極其高效的 **「無狀態開發終端 (Stateless Dev Terminal)」**，所有的重量級記憶與邏輯一致性都由伺服器的 **Mem0** 守護。

**你認為 `Recent Focus` 的清理週期，應該是按天執行，還是根據「任務完成」的手動觸發比較適合你的工作流？**
