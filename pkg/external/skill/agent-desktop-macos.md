# agent-desktop-macos

來源: [vercel-labs/agent-desktop](https://github.com/vercel-labs/agent-desktop) (Apache-2.0) / [agent-browser.dev/desktop](https://agent-browser.dev/desktop)

`agent-desktop` CLI 在 macOS 上的桌面自動化（Desktop Automation）技能，透過原生 Accessibility Tree 觀察與操作桌面應用程式，涵蓋 macOS TCC 權限引導、專屬 Surfaces、Native Address Refs 與常見故障排除。

## 所需 CLI 與安裝

使用此技能必須安裝 **`agent-desktop`** CLI（支援 macOS 13+，Intel 或 Apple Silicon）：

```sh
# Homebrew（推薦）
brew install agent-desktop

# 或透過 npm / pnpm（需 Node.js 18+）
npm install -g agent-desktop
pnpm add -g agent-desktop
```

安裝後可透過 `agent-desktop skills get desktop` 讀取跨平台核心技能文件（Observe → Act → Verify 迴圈）。

## 如何下載此 Skill

```sh
# 透過 skills CLI 安裝
npx skills add vercel-labs/agent-desktop
```

（本 repo 已收錄於 `.agents/skills/agent-desktop-macos/`，包含 `SKILL.md` 與 `references/troubleshooting.md`。）

## 首次接觸與權限設定

1. `agent-desktop permissions`：檢查 Accessibility（macOS 27+ 為 Device Control and Data Access）、Screen Recording（僅 `screenshot` 與 `session start --screenshots` 需要）與 Automation（System Events）狀態。
2. `agent-desktop permissions --guide`：開啟 System Settings 引導面板協助授權；授權後下一個 CLI 指令立即生效，無須重啟終端機。
3. `agent-desktop snapshot -i`：驗證可讀取前景視窗的互動式元素。

## 核心能力與 macOS 特性

| 項目 | 說明 |
| ---- | ---- |
| **多 Surface 觀察** | 支援 `--surface menubar`、`menu-bar-extras`（系統或特定 App 狀態列圖示）、`dock`、`menu`；可透過 `--app "Control Center"` 觀察無一般視窗的背景程式。大型樹狀結構可用 `--skeleton -i --compact` 搭配 `--root @ref` 鑽取 |
| **Native Address Refs** | `@ref` 綁定原生元素記憶體位址，能區分多個同名元件（如每列各自的 `Delete` 按鈕）；元件消失或 App 重啟時明確回報 `STALE_REF` |
| **Headed / Headless 分流** | 語意操作（`expand`、`set-value`、`select`、`toggle`、`check`、`focus` 等）維持語意路徑；`click`、`right-click`、`type`、`clear`、`scroll` 在 `--headed` 下優先走實體輸入；`double-click`、`hover`、`drag` 需搭配 `--headed` |
| **App 與鍵盤控制** | `launch` / `close-app` 支援顯示名稱或 Bundle ID；`press` 支援 `cmd` 修飾鍵與 ANSI/JIS 標點符號；`cmd+q` 等危險快捷鍵受 `POLICY_DENIED` 保護（需 `--force` 或改用 `close-app`） |
| **通知中心** | 提供 `list-notifications`、`dismiss-notification`、`dismiss-all-notifications`、`wait --notification`，操作後自動回復原 App 焦點 |
| **Chromium / Electron** | 自動啟用 Web Accessibility，或透過 `launch "App" --cdp` 交由 CDP 控制；支援 `cursor-overlay enable` 繪製不搶焦點的虛擬游標 |
| **macOS VM Sandbox** | `agent-desktop sandbox`（需 Apple Silicon、macOS 13+ 與 Xcode Command Line Tools）可建立與管理隔離的 macOS 客體 VM |
| **故障排除** | `references/troubleshooting.md` 提供 `PERM_DENIED`、空/截斷樹、`ACTION_FAILED`、`POLICY_DENIED`、`APP_UNRESPONSIVE` 等標準解法 |
