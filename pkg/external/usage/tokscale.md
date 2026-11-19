# tokscale

來源: [junhoyeo/tokscale](https://github.com/junhoyeo/tokscale) (MIT, ~5.9k stars, Rust core + Bun CLI, 最新 release `v4.23.0` 2026-11-15) · [website](https://tokscale.ai)

跨 60+ 個 AI coding agent 的 Token 用量與成本追蹤器. 預設開 TUI (Terminal UI, 終端互動介面), 附 GitHub 式 contribution graph, 並可選擇把用量上傳到公開 leaderboard. 定價來自 LiteLLM (1 小時磁碟快取, OpenRouter fallback).

## 功能

- 來源: Claude Code, Codex CLI, Gemini CLI, Antigravity (IDE / CLI), Cursor IDE, Copilot CLI, OpenCode, Hermes Agent, OpenClaw 等 60+ 個.
- TUI 10 個 view: Overview, Usage (訂閱配額), Models, Daily, Hourly, Monthly, Sessions, Projects, Stats, Agents.
- `--group-by workspace,model` / `machine,model`, `--merge-worktrees` 把 git worktree 併回主 repo.
- 分開統計 input / output / cache read / cache write / reasoning Token.
- `--light` (純表格), `--json`, `tokscale graph --output data.json`.
- 社群: `tokscale login` (GitHub auth) 後 `tokscale submit` 上傳; `autosubmit enable --interval 24h` 由 launchd 排程.

## 常用指令

```bash
npx tokscale@latest                     # TUI
npx tokscale@latest --light             # 表格輸出
tokscale antigravity status             # 偵測 Antigravity language server
tokscale antigravity sync               # 同步 Antigravity 用量到本機快取
tokscale submit --dry-run               # 預覽會上傳的內容
```

## 本 repo 的設定

TUI 偏好設定 [plugins/default/config/settings/tokscale.json](../../../plugins/default/config/settings/tokscale.json) 由 `scripts/run.sh` link 到 `~/.config/tokscale/settings.json` (green palette, 60 秒 auto refresh).

## 適用情境

- 同時使用多個 agent (Claude Code + Antigravity + Codex) 時, 在同一畫面比較用量.
- 依 project / workspace 看哪個 repo 最燒 Token.

## 注意

- Antigravity 同步只支援 macOS / Linux, 且編輯器必須開著 (讀其本機 language server).
- `submit` / `autosubmit` 會把用量公開到 tokscale.ai; API token 存在 `~/.config/tokscale/credentials.json`, 不要提交進 repo.
- Cursor 用量走 Cursor 的 API (會讀 desktop 的 access token), 非純本機.
