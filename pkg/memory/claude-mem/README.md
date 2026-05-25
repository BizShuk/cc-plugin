# claude-mem

<https://github.com/thedotmack/claude-mem>

### 工作原理 (How It Works)

`claude-mem` 是一個為 `Claude Code`、`OpenClaw` 和 `Gemini CLI` 等 `AI 代理 (AI Agents)` 設計的`持久化記憶壓縮系統 (Persistent Memory Compression System)`。

1. `生命週期鉤子 (Lifecycle Hooks)`：透過監聽 `SessionStart`、`PostToolUse` 等鉤子，自動擷取代理的操作紀錄。
2. `智慧摘要 (Smart Summarization)`：利用 `AI` 對擷取的觀察結果進行語義摘要，並儲存於 `SQLite` 資料庫中。
3. `上下文注入 (Context Injection)`：在未來的對話中，系統會根據當前需求自動注入相關的歷史背景。
4. `混合搜尋 (Hybrid Search)`：結合 `Chroma` 向量資料庫進行語義搜尋，並搭配關鍵字搜尋，精準檢索過往知識。
5. `三層工作流 (3-Layer Workflow)`：採用 `search` (搜尋索引)、`timeline` (時間軸上下文) 及 `get_observations` (獲取詳細觀察) 的分層模式，最大化節省 `Token` 成本。

### 主要優勢 (Key Benefits)

- `持久化上下文 (Persistent Context)`：知識可跨對話階段 (Sessions) 存續，代理能「記住」過去的工作內容。
- `提升連續性 (Enhanced Continuity)`：即使在對話中斷或重啟後，仍能維持專案知識的連續性。
- `節省 Token (Token Efficiency)`：分層檢索策略僅在需要時獲取詳細資訊，有效降低 `Token` 調用成本。
- `隱私保護 (Privacy Control)`：支援 `<private>` 標籤，可排除敏感內容不被儲存。
- `視覺化監控 (Visual Insights)`：提供 `Web 查看器 (Web Viewer UI)` (預設為 <http://localhost:37777)，方便即時監看記憶流。>
