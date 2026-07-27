# Claude Code 系統提示詞庫分析報告 (Claude Code System Prompts Analysis Report)

[claude-code-system-prompts](https://github.com/Piebald-AI/claude-code-system-prompts) 儲存庫的核心價值與功能。

## 1. 它能做什麼？ (What it can do)

該儲存庫的主要功能是**提取並追蹤** Anthropic 推出的 `Claude Code` 命令列工具中所有的系統提示詞 (System Prompts)。

- **全面提取 (Full Extraction)**：從 `Claude Code` 的編譯原始碼 (Compiled Source Code) 中直接提取超過 110 個提示詞字串。
- **工具說明 (Tool Descriptions)**：包含 24 個內建工具（如 `Bash`, `Write`, `TodoWrite` 等）的詳細功能描述。
- **子代理提示詞 (Sub-agent Prompts)**：揭露 `Explore`, `Plan`, `Task` 等子代理程式的運作邏輯。
- **公用程式提示詞 (Utility Prompts)**：包含自動生成 `CLAUDE.md`、安全審查 (`/security-review`)、批次處理 (`/batch`) 等功能的指令。
- **即時更新 (Real-time Updates)**：在 `Claude Code` 每次發布新版本後的幾分鐘內即完成更新，保持資訊最前線。
- **變更日誌 (Changelog)**：詳細記錄了從 v2.0.14 以來超過 180 個版本的提示詞演進過程。

## 2. 它的優點是什麼？ (What are the benefits)

- **透明度 (Transparency)**：讓開發者能夠確切知道 `Claude` 在後台收到的指令內容，消除「黑箱」作業。
- **學習與參考 (Educational Reference)**：是學習如何構建「代理式 AI」(Agentic AI) 與撰寫高品質「系統提示詞」的絕佳教科書。
- **行為偵錯 (Behavioral Debugging)**：當 `Claude Code` 出現非預期行為時，開發者可以透過檢查提示詞來理解其邏輯原因。
- **賦予自定義能力 (Enables Customization)**：配合 `tweakcc` 工具，開發者可以參考此庫的內容來修改自己本地端的 `Claude Code` 行為。
- **版本對比 (Historical Context)**：提供歷史紀錄，觀察 Anthropic 如何微調 AI 的安全性、語氣以及任務處理策略。

## 3. 為什麼需要它？ (Why)

- **結構複雜 (Complexity)**：`Claude Code` 並非使用單一提示詞，而是根據環境與配置動態組合上百個片段。
- **程式碼混淆 (Minified Code)**：提示詞隱藏在高度混淆與壓縮的 JavaScript 檔案中，一般使用者難以直接閱讀。
- **逆向工程的需求 (Reverse Engineering)**：為了優化與 AI 的協作，開發者需要了解 AI 被賦予的角色設定與限制條件。
- **持續演進 (Constant Evolution)**：AI 工具更新極快，需要一個自動化的機制來捕捉每一次細微的指令調整。

---

> [!TIP]
> 如果你想修改本地的 `Claude Code` 提示詞，建議搭配 [tweakcc](https://github.com/Piebald-AI/tweakcc) 使用。
