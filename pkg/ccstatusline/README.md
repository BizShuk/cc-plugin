# ccstatusline - Claude Code 自定義狀態欄 (Customizable Statusline)

`ccstatusline` 是一個為 `Claude Code` CLI 設計的高度可自定義狀態欄工具。它能美化終端機界面，並即時顯示模型資訊、Git 狀態及 Token 使用量等關鍵數據。

## 主要特色 (Features)

- **豐富的組件 (Widgets)**：支援顯示模型名稱、Token 速度、上下文長度 (Context Length)、Git 分支、PR 狀態及剩餘配額。
- **高度自定義 (Customizable)**：內建 TUI 配置介面，可自由排列組件、選擇顏色主題及自定義 Powerline 分隔符號。
- **支援 Nerd Fonts**：可啟用圖示顯示，讓狀態欄更具視覺吸引力。
- **Git 整合**：即時追蹤當前分支的變更狀態（暫存/未暫存檔案數量）與 PR/MR 連結。

## 快速開始 (Quick Start)

```bash
npx -y ccstatusline
# 設定檔路徑: `~/.config/ccstatusline/settings.json`
rm ~/.config/ccstatusline/settings.json
ln -s "$HOME/projects/cc-plugin/pkg/ccstatusline/settings.json" "$HOME/.config/ccstatusline/settings.json"
npx -y ccstatusline
```

> [!NOTE]
> 安裝完成後，請執行指令並依照 TUI 介面進行初始設定與更新。

## 相關連結 (Links)

- **GitHub**: <https://github.com/sirmalloc/ccstatusline>
