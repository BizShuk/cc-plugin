# Tools — macOS 工具整合插件 (macOS Tool Integration Plugin)

提供 Apple 原生應用程式的 CLI 技能，統一由既有工具操作使用者資料：

| 技能 (Skill) | CLI | 用途 |
| :--- | :--- | :--- |
| `apple-calendar` | `accli` | 行事曆事件與空檔查詢 |
| `apple-email` | `email` | Apple Mail 郵件讀寫與整理 |
| `apple-notes` | `notes` | Apple Notes 筆記與資料夾管理 |
| `apple-reminders` | `remindctl` | 提醒事項與清單管理 |

## 安裝與使用 (Installation and Usage)

此插件由根目錄 `.claude-plugin/marketplace.json` 以 `./plugins/tools` 註冊。
各技能只在對應 Apple 應用操作或使用者明確要求時觸發；實際 CLI 的權限與安全守衛以各自 `SKILL.md` 為準。
