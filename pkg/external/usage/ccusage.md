# ccusage

來源: [ccusage/ccusage](https://github.com/ccusage/ccusage) (MIT, ~19k stars, Rust, 最新 release `v20.0.51` 2026-11-14) · [website](https://ccusage.com/)

讀取本機 coding agent CLI 的 session 紀錄 (Claude Code JSONL 等), 彙整成 daily / weekly / monthly / session 的 Token 用量與 USD 成本報表. 純 CLI, 不需登入, 定價來自 LiteLLM pricing data (`--offline` 用內建快照).

## 功能

- 支援 20 個來源: Claude Code, Codex, OpenCode, Amp, Gemini CLI, Antigravity, Grok Build CLI, Hermes Agent, Copilot CLI 等; 不指定來源時合併成一張報表.
- `ccusage blocks`: Claude Code 5 小時計費窗口 (billing window) 的用量與目前 active block.
- `ccusage statusline` (Beta): 給 Claude Code status line hook 用的精簡字串.
- 分開統計 cache creation / cache read Token; `--breakdown` 看各模型成本.
- 過濾: `--since` / `--until` / `--last 1`, `--timezone`, `claude daily --instances --project <name>`.
- 輸出: 彩色表格 (窄終端自動 compact), `--json`, `--no-cost`; `ccusage.json` 可覆寫模型定價.

## 常用指令

```bash
npx ccusage@latest              # 所有來源, 依日彙整 (預設)
npx ccusage claude daily        # 只看 Claude Code
npx ccusage session             # 依 session
npx ccusage blocks              # 5 小時計費窗口
npx ccusage daily --json --last 1
```

## 適用情境

- 快速回答「今天 / 這個月花了多少」, 或用 `--json` 餵給腳本與 dashboard.
- 訂閱制 (Max) 使用者追蹤 5 小時窗口, 預估何時撞到 limit.

## 注意

- 成本是依公開 API 定價換算的`估計值`, 訂閱制實際不按此計費.
- 只讀本機檔案; 預設會連網抓最新定價, 離線環境加 `--offline`.
