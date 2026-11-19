# 外部工具與參考索引 (External Tools & References)

收錄與評估跨 AI Agent 生態系的外部工具、技能 (Skills)、記憶系統 (Memory)、桌面應用 (Apps) 與服務參考，依類別分置於各子目錄。

## 精選工具 (Featured Tools)

| 工具 / 專案 | 類別 | 定位與核心能力 | 詳細筆記 |
| :--- | :--- | :--- | :--- |
| **agent-desktop-macos** | `skill` | **macOS 桌面自動化 (Desktop Automation)**：透過原生 Accessibility Tree 觀察與操作桌面 App，支援 Menu Bar / Dock / 通知中心、Native Address Refs (`@ref`)、Headed/Headless 分流與 Chromium/Electron CDP 控制，並提供 TCC 權限引導與故障排除。 | [`skill/agent-desktop-macos.md`](skill/agent-desktop-macos.md) |
| **MarkItDown** | `tool` | **文件轉 Markdown (Document-to-Markdown)**：Microsoft 維護的 Python CLI / Library / MCP Server，將 PDF、Office（Word / PowerPoint / Excel）、EPUB、`.msg`、音訊與 YouTube 字幕轉為保留結構的 Markdown，作為 Agent 讀檔與 RAG 前處理管線（可選 `markitdown-ocr` 或 Azure 後端）。 | [`tool/markitdown.md`](tool/markitdown.md) |
| **Claude-Mem** | `memory` | **跨 Session 自動記憶 (Automatic Session Memory)**：以 Hooks + 本機 Worker + Observer Agent 自動記錄 Tool Call 為 Observation（SQLite FTS5 + Chroma），於下次 Session 開場注入精簡 Index，並透過 `search` → `timeline` → `get_observations` 漸進展開細節以節省 Token。 | [`memory/claude-mem.md`](memory/claude-mem.md) |
| **OpenWhispr** | `app` | **Privacy-First 語音轉文字與會議筆記 (Voice-to-Text & Meeting Notes)**：開源免費桌面 App，支援全域 Hotkey 語音聽寫、翻譯、Voice Assistant 與會議講者辨識 (Speaker Diarization)；支援本機離線 ASR (`whisper.cpp` / `sherpa-onnx`) 或 BYOK，並提供 Public API 與 MCP Server 整合。 | [`app/openwhispr.md`](app/openwhispr.md) |

## 目錄結構 (Directory Structure)

| 目錄 | 說明 |
| :--- | :--- |
| `app/` | 桌面與終端整合應用參考 |
| `memory/` | 跨 Session 記憶、知識圖譜與脈絡檢索工具 |
| `service/` | 可自架 (Self-hosted) 或雲端協作與生成式服務 |
| `skill/` | 外部 Agent Skills 評估與安裝指引 |
| `tool/` | 命令列與開發輔助工具 |
| `usage/` | Token 用量與花費追蹤工具 |
