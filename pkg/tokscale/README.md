# `tokscale` (Tokscale)

<https://github.com/junhoyeo/tokscale>

`tokscale` (Tokscale) 是一個高效能的`命令列工具` (CLI tool) 與`視覺化儀表板` (Visualization dashboard)，專為追蹤與分析多個 `AI 編碼代理` (AI coding agents) 的`標記使用量` (Token usage) 及`花費成本` (Costs) 而設計。

## `支援平台` (Supported Platforms)

`tokscale` 支援追蹤多個主流 `AI` 工具的標記使用量，包括但不限於：

- `Google 無重力` (Google Antigravity)
- `克勞德代碼` (Claude Code)
- `探索者` (Cursor IDE)
- `副駕駛` (Copilot CLI)
- `雙子座` (Gemini CLI)
- `開放爪` (OpenClaw)
- `赫米斯` (Hermes Agent)

## `安裝與執行` (Installation & Running)

您可以使用 `npx` 或 `bunx` 來直接執行：

```bash
npx tokscale@latest
# 或使用 bun
bunx tokscale@latest
```

## `常用指令` (Common Commands)

### `一般指令` (General Commands)

- `npx tokscale@latest`：啟動主要的命令列文字介面與儀表板。
- `bunx tokscale@latest submit`：將您的標記使用量數據提交至`全球排行榜` (Global Leaderboard) 並建立您的公開個人檔案。

### `無重力專屬指令` (Antigravity Commands)

`Google 無重力` (Google Antigravity) 的同步目前僅支援在 `macOS` 與 `Linux` 作業系統上運行。當啟用 `無重力` (Antigravity) 的編輯器開啟且其本地語言伺服器可用時，`tokscale` 會讀取該伺服器的使用數據並在本地快取 normalized artifacts。

- `npx tokscale@latest antigravity status`：檢查 `tokscale` 是否能偵測到正在運行的 `Google 無重力` (Google Antigravity) 語言伺服器。
- `npx tokscale@latest antigravity sync`：從本地的 `Google 無重力` (Google Antigravity) 語言伺服器同步標記使用數據至 `tokscale` 的快取中。
- `npx tokscale@latest antigravity purge-cache`：刪除本地已快取的 `Google 無重力` (Google Antigravity) 標記使用量快取。

## `快取位置` (Cache Location)

- `Google 無重力` (Google Antigravity) 的快取數據存放於：
  `~/.config/tokscale/antigravity-cache/sessions/*.jsonl`
- `tokscale` (Tokscale) 的通用快取存放於：
  `~/.config/tokscale/cache/`

## `設計靈感` (Inspiration)

本專案的靈感源自於`卡爾達肖夫指數` (Kardashev scale)，這是一種由天體物理學家 `尼古拉·卡爾達肖夫` (Nikolai Kardashev) 提出的方法，用來根據文明所消耗的能量來衡量其科技先進程度。
在 `AI` 輔助開發的時代，`標記` (Tokens) 就是新型態的能量。它們為我們的邏輯推理提供動力，推動我們的生產力與創造力輸出。正如`卡爾達肖夫指數`追蹤宇宙級的能量消耗，`tokscale` 則在您提升 `AI` 輔助開發能力的過程中，測量您的`標記`消耗量，幫助您視覺化這段從行星級開發者到銀河級代碼架構師的旅程。
