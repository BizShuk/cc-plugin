# ccusage - Claude Code 使用量分析工具 (Claude Code Usage Analyzer)

`ccusage` 是一個強大的命令列工具 (CLI tool)，專門用於分析 `Claude Code` 或 `Codex CLI` 在本地產生的 JSONL 紀錄檔，幫助你精確掌握 AI 的使用成本與 Token 消耗狀況。

## 核心功能 (Key Features)

*   **多維度報告 (Multi-dimensional Reports)**：支援每日 (Daily)、每月 (Monthly) 以及對話會話 (Session) 等級別的聚合報告。
*   **成本估算 (Cost Estimation)**：自動根據最新的 LiteLLM 定價資料，計算 Input、Output 以及快取 (Cache) 的美金成本。
*   **模型追蹤 (Model Tracking)**：清晰顯示不同模型（如 Claude 3.5 Sonnet, Claude 3 Opus 等）的使用分布。
*   **快取分析 (Cache Metrics)**：追蹤快取寫入 (Cache Creation) 與快取讀取 (Cache Read) 的 Token，優化 Prompt Caching 策略。
*   **靈活過濾 (Filtering)**：支援日期範圍 (--since/--until)、特定專案 (--project) 以及多實例 (--instances) 過濾。
*   **狀態欄整合 (Statusline Integration)**：提供緊湊的狀態字串，可與 Claude Code 的 Hook 整合，在終端機顯示即時成本。

## 常用指令 (Common Commands)

```bash
npx ccusage          # 顯示每日摘要報告 (預設)
npx ccusage daily    # 顯示每日 Token 使用量與成本
npx ccusage monthly  # 顯示每月聚合報告
npx ccusage session  # 按對話會話分類顯示使用量
npx ccusage blocks   # 顯示 5 小時計費窗口的狀態
npx ccusage statusline  # 用於 Hook 的精簡狀態行 (Beta)
```

## 相關連結 (Links)

*   **GitHub**: <https://github.com/ryoppippi/ccusage>
*   **官方網站**: <https://ccusage.com/>
