# MemPalace

<https://mempalaceofficial.com/guide/getting-started.html>

### 簡介 (Introduction)

MemPalace 是一個開源且以本地優先 (Local-first) 的 AI 記憶系統，旨在為 AI 代理 (Agents) 提供持久的長期記憶能力。其靈感來自古老的「記憶宮殿 (Method of Loci)」技術，透過結構化的方式組織資訊，避免 AI 在對話結束後遺忘重要細節。

### 運作原理 (How It Works)

- **階層式架構：** 將記憶組織為 Wing (領域)、Room (主題)、Hall (類別) 與 Drawer (具體記錄)，實現結構化的導航與檢索。
- **逐字儲存 (Verbatim Storage)：** 不同於摘要式的記憶工具，MemPalace 傾向於逐字記錄對話內容，確保原始脈絡與細微細節不被遺漏。
- **本地技術棧：** 採用 ChromaDB (向量資料庫) 進行語義搜尋，並結合 SQLite 建立時間序列知識圖譜，所有處理均在本地完成。
- **MCP 整合：** 支援 Model Context Protocol (MCP)，可與 Claude、Cursor 等熱門 AI 工具無縫接軌。

### 核心優勢 (Benefits)

- **持久上下文 (Persistent Context)：** 消除 AI 的「健忘症」，讓代理能跨 Session 記住過去的決策與偏好。
- **高保真還原 (High Fidelity)：** 由於是逐字儲存，檢索結果比 AI 生成的摘要更具準確性與可追溯性。
- **隱私與安全：** 資料儲存在本地，無需上傳至雲端，適合處理敏感或私有數據。
- **零成本營運：** 身為本地優先的工具，無需支付昂貴的雲端記憶服務費用。
