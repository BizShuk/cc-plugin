---
name: apple-reminders
description: >
    Use when managing Apple 提醒事項 (Apple Reminders) on macOS via the `remindctl`
    CLI — listing, creating, editing, completing, or deleting reminders and lists.
    Triggers on: "提醒我", "add reminder", "show today reminders", "complete reminder",
    "delete reminder", "list Reminders", or any Reminders.app operation from the terminal.
version: "1.0.0"
allowed-tools: Bash
metadata:
    type: reference
    platforms: [macos]
    prerequisites:
        commands: [remindctl]
---

# Apple 提醒事項 (apple-reminders / remindctl)

使用 `remindctl` 可以直接從終端機 (terminal) 管理 Apple 提醒事項 (Apple Reminders)。

## 何時使用 (When to Use)

✅ **適用的情境：**

- 使用者明確提到「提醒事項 (reminder)」或「Reminders app」。
- 建立包含截止日期、並會同步至 iOS 裝置的個人待辦事項。
- 管理 Apple 提醒事項列表 (Reminders lists)。
- 希望任務顯示在 iPhone / iPad 的 Reminders 應用程式中。

## 何時避免使用 (When NOT to Use)

❌ **不適用的情境：**

- 排程內部警示或任務執行 → 請改用 `cron` 工具搭配系統事件 (systemEvent)。
- 行事曆事件 (Calendar events) 或會議約會 → 請使用 Apple Calendar。
- 專案與工作任務管理 → 請使用 Notion、GitHub Issues 或任務佇列 (task queue)。
- 單次通知 (One-time notifications) → 請改用 `cron` 工具進行定時警示。
- 當使用者說「提醒我 (remind me)」但其實是指內部系統的警示 → 請先與使用者釐清需求。

## 初始設定 (Setup)

- **透過 Homebrew 安裝 (Install via Homebrew)**:

    ```bash
    brew install steipete/tap/remindctl
    ```

- 僅支援 macOS，當系統跳出權限請求時，請允許存取「提醒事項」。
- 檢查狀態: `remindctl status`
- 請求權限存取: `remindctl authorize`

## 常用指令 (Common Commands)

### 檢視提醒事項 (View Reminders)

- `remindctl`：今天的提醒事項
- `remindctl today`：今天
- `remindctl tomorrow`：明天
- `remindctl week`：本週
- `remindctl overdue`：已過期 (Past due)
- `remindctl all`：全部 (Everything)
- `remindctl 2026-01-04`：指定特定日期

### 管理列表 (Manage Lists)

- `remindctl list`：列出所有列表
- `remindctl list Work`：顯示名為 Work 的特定列表
- `remindctl list Projects --create`：建立名為 Projects 的新列表
- `remindctl list Work --delete`：刪除名為 Work 的列表

### 建立提醒事項 (Create Reminders)

- `remindctl add "Buy milk"`：簡單新增事項
- `remindctl add --title "Call mom" --list Personal --due tomorrow`：指定列表為 Personal 並設定截止時間為明天
- `remindctl add --title "Meeting prep" --due "2026-02-15 09:00"`：設定特定日期與時間

### 完成與刪除 (Complete/Delete)

- `remindctl complete 1 2 3`：透過 ID 完成指定的提醒事項
- `remindctl delete 4A83 --force`：透過 ID 強制刪除提醒事項

### 輸出格式 (Output Formats)

- `remindctl today --json`：以 JSON 格式 (JSON format) 輸出 (適合腳本處理)
- `remindctl today --plain`：以 TSV 格式 (TSV format) 輸出 (純文字)
- `remindctl today --quiet`：僅顯示數量計數

## 日期與時間格式 (Date Formats)

`--due` 參數及日期過濾支援下列格式：

- `today`, `tomorrow`, `yesterday`
- `YYYY-MM-DD` (例如：`2026-01-04`)
- `YYYY-MM-DD HH:mm`
- ISO 8601 (例如：`2026-01-04T12:34:56Z`)
